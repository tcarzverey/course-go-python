# Практикум: Observability в Go-сервисах

## Что такое Observability?

Observability («наблюдаемость») — способность понять внутреннее состояние системы
по её внешним проявлениям. Три кита:

| Сигнал   | Что даёт                                      | Инструмент       |
|----------|-----------------------------------------------|------------------|
| Метрики  | Агрегированные числа во времени               | Prometheus       |
| Логи     | Дискретные события с контекстом               | Loki             |
| Трейсы   | Путь запроса через сервисы                    | Tempo (OTel)     |

Мы будем добавлять их по одному к живому HTTP-сервису — **URL shortener**.

---

## Как работают step-комментарии

Весь итоговый код уже написан, но части, относящиеся к следующим шагам,
закомментированы по шаблону `// stepN код`.

Чтобы перейти к шагу N:
1. Откройте **Find and Replace** в IDE (Ctrl+Shift+H / Cmd+Shift+H).
2. Найдите `// stepN ` (с пробелом после цифры).
3. Замените на `` (пустая строка).
4. Применить ко всему проекту.

Например, для шага 3: найти `// step3 ` → заменить на ``.

---

## Сервисы

```
url-shortener (port 8080)   — основной сервис
stats-service  (port 8081)  — считает клики (step7)

Prometheus     (port 9090)
Grafana        (port 3000)  — логин admin / admin
Loki           (port 3100)
Tempo          (port 3200)
Jaeger         (port 16686) — альтернативный трейс-backend
OTel Collector (port 4317)
```

---

## Быстрый старт

```bash
# Запустить весь стек
docker compose up --build -d

# Проверить что сервис живёт
curl http://localhost:8080/health

# Создать короткую ссылку
curl -X POST http://localhost:8080/shorten \
     -H "Content-Type: application/json" \
     -d '{"url": "https://go.dev/doc/effective_go"}'
# → {"code":"aB3xYz","short_url":"http://localhost:8080/aB3xYz"}

# Редирект (откроет браузер или следуй -L)
curl -L http://localhost:8080/aB3xYz

# Статистика
curl http://localhost:8080/stats/aB3xYz
```

---

## Шаги практикума

---

### ✅ Шаг 0 — Базовый сервис (уже готово)

**Что есть:**
Простой HTTP URL shortener на стандартной библиотеке Go 1.22.
In-memory хранилище (sync.RWMutex + map).

**Эндпоинты:**
- `POST /shorten` — создать короткую ссылку
- `GET /{code}` — редирект на оригинальный URL
- `GET /stats/{code}` — количество кликов
- `GET /health` — хелсчек

**Запуск без Docker:**
```bash
cd url-shortener
go run ./main.go
```

**На что обратить внимание:**
Сервис работает, но мы ничего не знаем о том, что происходит внутри.
Сколько запросов? Насколько быстро? Есть ли ошибки?

---

### 🐳 Шаг 1 — Поднимаем инфраструктуру observability

**Цель:** запустить весь observability-стек через Docker Compose.

```bash
docker compose up -d
```

**Компоненты:**

```
┌─────────────────────────────────────────────┐
│  Go сервис                                  │
│  url-shortener :8080                        │
└──────────┬──────────────────────────────────┘
           │ OTLP/gRPC (traces + logs + metrics)
           ▼
    OTel Collector :4317
           │
           ├──► Tempo   :3200  (traces)
           ├──► Jaeger  :16686 (traces)
           ├──► Loki    :3100  (logs)
           └──► :8889          (metrics, Prometheus scrapes)
                │
           Prometheus :9090
                │
           Grafana :3000
```

**Откройте в браузере:**
- Grafana: http://localhost:3000 (admin/admin)
- Prometheus: http://localhost:9090
- Jaeger UI: http://localhost:16686

**На что смотреть:**
- В Grafana видны пустые datasources — Prometheus, Loki, Tempo
- Prometheus → Status → Targets: url-shortener пока `DOWN` (нет /metrics)
- Jaeger UI — пустой список сервисов, ждёт трейсов (появятся с step5)
- Готовая инфраструктура ждёт данных от нашего сервиса

**Jaeger vs Tempo:**
OTel Collector шлёт трейсы одновременно в оба backend'а.
Tempo — интегрирован с Grafana (correlation с логами).
Jaeger — standalone UI, удобен для быстрого просмотра и поиска трейсов без Grafana.
С step5 и далее трейсы видны в обоих местах.

---

### 📡 Шаг 2 — OTel SDK (все три сигнала)

**Активировать:** найти `// step2 ` → заменить на ``, перебилдить.

**Что добавляется:** `telemetry.Init` в `main.go` — одним вызовом поднимаются все три провайдера:

```go
shutdown, err := telemetry.Init(ctx, "url-shortener", otlpEndpoint)
defer func() { _ = shutdown(context.Background()) }()
```

Внутри `telemetry.Init` (`internal/telemetry/telemetry.go`):
```
gRPC conn → OTel Collector :4317
    ├── TracerProvider  → traces  → Tempo / Jaeger
    ├── MeterProvider   → metrics → Prometheus (через Collector :8889)
    └── LoggerProvider  → logs    → Loki
```

**На что смотреть:**
- Prometheus → Status → Targets → `otel-collector` = **UP** (метрики приложения придут на step3)
- Tempo / Jaeger — пустые (трейсы придут на step3 через otelhttp)
- Loki — пустой (логи придут на step4 через otelslog)

**Ключевое отличие от ручного подхода:**
Никакого `/metrics` эндпоинта — сервис сам *пушит* данные в коллектор по gRPC.

---

### ✏️ Шаг 3 — Ручные спаны (OTel Tracing API)

**Активировать:** найти `// step3 ` → заменить на ``, перебилдить.

**Что добавляется в `handler.go`:**
```go
func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
    ctx, span := otel.Tracer("url-shortener").Start(r.Context(), "handler.Shorten")
    defer span.End()

    // ... бизнес-логика ...

    span.SetAttributes(
        attribute.String("url.code", code),
        attribute.String("url.original", req.URL),
    )
}
```

Аналогично для `Redirect` — `handler.Redirect` со своими атрибутами.

**На что смотреть в Tempo / Jaeger:**
- Появились спаны: `handler.Shorten`, `handler.Redirect`
- Каждый span содержит атрибуты: `url.code`, `url.original`
- Ошибки: `span.RecordError(err)` + `span.SetStatus(codes.Error, "...")` → красный span

**Чего не хватает:**
Каждый span — изолированный. Нет HTTP-контекста (`http.method`, `status_code`).
Каждый хэндлер нужно инструментировать вручную.

---

### 🔧 Шаг 4 — OTel HTTP middleware (автоматические спаны и метрики)

**Активировать:** найти `// step4 ` → заменить на ``, перебилдить.

**Что добавляется в `main.go`:**
```go
root = otelhttp.NewHandler(root, "url-shortener-http",
    otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
)
```

`otelhttp` использует глобальные провайдеры (step2) и автоматически:
- Создаёт корневой **span** `url-shortener-http` для каждого запроса
- Записывает **метрики**: `http.server.request.duration`, `http.server.active_requests`

**На что смотреть в Tempo / Jaeger:**
```
url-shortener-http          [5ms]   ← добавил otelhttp (step4)
  └── handler.Shorten       [3ms]   ← наш ручной span (step3)
        attrs: url.code, url.original
```

Спаны теперь вложены — `otelhttp` создаёт родительский span, наши ручные становятся дочерними автоматически через context propagation.

**Prometheus:** метрика `otelcol_http_server_request_duration_seconds` из коллектора `:8889`.

**Ключевое отличие от step3:**
На step3 мы писали инструментацию руками и видели только бизнес-логику.
На step4 весь HTTP-слой закрыт одной строкой, а ручные спаны дополняют картину.

---

### 📝 Шаг 5 — Структурированные логи (slog + OTel bridge)

**Активировать:** найти `// step5 ` → заменить на ``, перебилдить.

**Что добавляется:**
1. В `main.go` — OTel-backed slog через bridge:
```go
logger = slog.New(otelslog.NewHandler("url-shortener"))
slog.SetDefault(logger)
```

2. В `middleware/middleware.go` — лог каждого запроса:
```go
logger.InfoContext(r.Context(), "http request",
    "method", r.Method,
    "path", r.URL.Path,
    "status", rw.status,
    "trace_id", traceID,  // заполнится на step6
)
```

3. В `handler.go` — контекстные логи событий:
```go
h.logger.InfoContext(r.Context(), "URL shortened", "code", code, "original_url", req.URL)
h.logger.WarnContext(ctx, "code not found", "code", code)
```

**Как это работает:**
```
slog.Logger (otelslog.NewHandler)
    → OTel LoggerProvider (global, set by InitLogProvider)
        → OTel Collector :4317 (OTLP/gRPC)
            → Loki
```

**Проверяем:**
```bash
# В Grafana → Explore → Loki:
{service_name="url-shortener"}
{service_name="url-shortener"} | json | severity="WARN"
```

**На что смотреть:**
- Логи приходят в Loki через OTel Collector, а не через Promtail
- Каждая запись содержит structured attributes: `method`, `path`, `status`
- `service_name` проставляется автоматически из OTel Resource

---

### 🔗 Шаг 6 — Корреляция логов и трейсов

**Активировать:** найти `// step6 ` → заменить на ``, перебилдить.

**Что добавляется в `middleware/middleware.go`:**
```go
if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
    traceID = span.SpanContext().TraceID().String()
}
logger.InfoContext(ctx, "http request", ..., "trace_id", traceID)
```

`otelhttp` (step4) уже добавил span в контекст каждого запроса.
step6 читает его `trace_id` и вкладывает в лог-запись — получается связь между Loki и Tempo.

**На что смотреть:**

**Loki + Tempo correlation:**
Grafana → Explore → Loki → `{service_name="url-shortener"}`
Кликните на лог-строку → в поле `trace_id` появится кнопка **"View Trace in Tempo"**
→ прямой переход к трейсу!

**Ключевой паттерн:**
`trace_id` в логах = мост между Loki и Tempo.
Один идентификатор связывает все три сигнала: метрики по времени → лог-строки → полный путь запроса.

---

### 🌐 Шаг 7 — Второй сервис + Distributed Tracing

**Активировать:** найти `// step7 ` → заменить на ``, перебилдить.

**Что добавляется:**

Новый сервис `stats-service` (порт 8081):
- `POST /track` — принимает `{"code": "abc"}`, увеличивает счётчик
- `GET /stats/{code}` — возвращает количество кликов
- Полная OTel инструментация

`url-shortener/handler.go` — вызов stats-service при редиректе:
```go
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
    // ...
    if err := h.trackClick(ctx, code); err != nil { ... }
    http.Redirect(w, r, record.OriginalURL, http.StatusFound)
}

func (h *Handler) trackClick(ctx context.Context, code string) error {
    ctx, span := otel.Tracer("url-shortener").Start(ctx, "handler.trackClick")
    defer span.End()

    req, _ := http.NewRequestWithContext(ctx, http.MethodPost, h.statsServiceURL+"/track", body)

    // ← Это ключевая строка для distributed tracing!
    otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

    http.DefaultClient.Do(req)
}
```

`stats-service/handler.go` — извлечение контекста:
```go
func (h *Handler) Track(w http.ResponseWriter, r *http.Request) {
    // ← Без этого — новый трейс вместо продолжения существующего
    ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
    ctx, span := otel.Tracer("stats-service").Start(ctx, "handler.Track")
    defer span.End()
    // ...
}
```

**Запускаем:**
```bash
docker compose up --build -d
curl -X POST http://localhost:8080/shorten -d '{"url":"https://go.dev"}' -H 'Content-Type: application/json'
curl http://localhost:8080/<code>
```

**На что смотреть в Tempo:**

```
url-shortener-http              [200ms]
  └── handler.Redirect          [180ms]
        └── handler.trackClick  [ 15ms]  ← url-shortener span
              └── handler.Track  [ 12ms]  ← stats-service span (другой процесс!)
```

Это **один трейс**, охватывающий **два разных сервиса**.
Связь обеспечивается HTTP-заголовком `traceparent` (W3C Trace Context):
```
traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
                  ↑ trace_id (16 байт)             ↑ span_id (8 байт)
```

**Distributed tracing решает:** Когда запрос медленный, но непонятно в каком сервисе —
трейс точно показывает, сколько времени провёл запрос в каждом звене цепочки.

---

### 🗄️ Шаг 8 — pgx Tracing (автоматическая инструментация БД)

**Активировать:** найти `// step8 ` → заменить на ``, затем `go mod tidy && docker compose up --build -d`.

**Что добавляется:**

1. PostgreSQL-хранилище (`pgx/v5`) вместо in-memory map
2. `otelpgx.NewTracer()` — автоматически создаёт span для каждого SQL-запроса

`internal/storage/postgres.go` — ключевые строки:
```go
config.ConnConfig.Tracer = otelpgx.NewTracer()  // ← вся "магия" в одной строке
pool, err := pgxpool.NewWithConfig(ctx, config)
```

`main.go`:
```go
pgStore, err := storage.NewPGX(ctx, dbDSN)
if err != nil { log.Fatalf("db: %v", err) }
store = pgStore
```

**Запускаем:**
```bash
docker compose up --build -d
curl -X POST http://localhost:8080/shorten -d '{"url":"https://go.dev"}' -H 'Content-Type: application/json'
curl http://localhost:8080/<code>
```

**На что смотреть в Tempo / Jaeger:**

```
url-shortener-http              [12ms]
  └── handler.Shorten           [ 8ms]
        └── db.Exec             [ 5ms]  ← автоматический span от otelpgx
              attrs: db.system=postgresql
                     db.statement="INSERT INTO urls..."
                     net.peer.name=postgres

url-shortener-http              [10ms]
  └── handler.Redirect          [ 7ms]
        ├── db.QueryRow         [ 3ms]  ← SELECT
        └── db.Exec             [ 1ms]  ← UPDATE clicks
```

**Ключевой момент:**
Нулевые изменения в бизнес-коде — `Save`, `Get`, `IncrementClicks` не тронуты.
Вся инструментация через единственный вызов `otelpgx.NewTracer()` в конфиге пула.

Это и есть **автоматическая инструментация** — библиотека знает о pgx internals
и внедряет spans через официальный hook-интерфейс pgx трассировщика.

---

## Итоговая архитектура observability

```
┌─────────────────────────────────────────────────────────┐
│                    Go Services                          │
│                                                         │
│  url-shortener ──HTTP──► stats-service                  │
│                                                         │
│  OTel SDK: traces + logs + metrics (OTLP/gRPC push)     │
└─────────────────────┬───────────────────────────────────┘
                      │ OTLP/gRPC
                      ▼
   ┌──────────────────────────────┐
   │      OTel Collector          │
   │  receiver: otlp              │
   │  processors: batch, memlimit │
   │  exporters:                  │
   │    → Tempo   (traces)        │
   │    → Jaeger  (traces)        │
   │    → Loki    (logs)          │
   │    → :8889   (metrics)       │
   └──────────────────────────────┘
          │                   ▲
          ▼                   │ scrape /metrics
   ┌──────────────────────────────┐
   │  Prometheus + Grafana        │
   │  ← OTel metrics (:8889)      │
   │  ← Loki datasource           │
   │  ← Tempo datasource          │
   │  ← Jaeger (standalone UI)    │
   └──────────────────────────────┘
```

---

## Полезные PromQL запросы

```promql
# RPS по методу и пути
rate(otelcol_http_server_request_duration_seconds_count{service_name="url-shortener"}[1m])

# p99 latency по пути
histogram_quantile(0.99, sum by (http_route, le) (
  rate(otelcol_http_server_request_duration_seconds_bucket{service_name="url-shortener"}[1m])
))

# Error rate (5xx)
rate(otelcol_http_server_request_duration_seconds_count{service_name="url-shortener", http_response_status_code=~"5.."}[5m])
/ rate(otelcol_http_server_request_duration_seconds_count{service_name="url-shortener"}[5m])
```

## Полезные LogQL запросы (Loki)

```logql
# Все логи сервиса
{service_name="url-shortener"}

# Только предупреждения и ошибки
{service_name="url-shortener"} | json | severity="WARN"
{service_name="url-shortener"} | json | severity="ERROR"

# Медленные запросы (если добавить duration в лог)
{service_name="url-shortener"} | json | duration > 100ms

# Найти лог по trace_id
{service_name="url-shortener"} | json | trace_id="4bf92f3577b34da6a3ce929d0e0e4736"
```
