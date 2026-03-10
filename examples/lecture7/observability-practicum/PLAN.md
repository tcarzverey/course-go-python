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
           │
    ┌──────┴──────┐
    │             │
    ▼             ▼
Prometheus    OTel Collector :4317
 :9090        (step5: traces/logs)
    │             │
    ▼             ├──► Tempo :3200  (traces)
  Grafana         └──► Loki  :3100  (logs)
  :3000  ◄──────────────────────────────────
```

**Откройте в браузере:**
- Grafana: http://localhost:3000 (admin/admin)
- Prometheus: http://localhost:9090

**На что смотреть:**
- В Grafana видны пустые datasources — Prometheus, Loki, Tempo
- Prometheus → Status → Targets: url-shortener пока `DOWN` (нет /metrics)
- Готовая инфраструктура ждёт данных от нашего сервиса

---

### 📊 Шаг 2 — Ручные метрики Prometheus

**Активировать:** найти `// step2 ` → заменить на ``, затем `docker compose up --build -d`.

**Что добавляется в `main.go`:**
- Эндпоинт `GET /metrics` через `promhttp.Handler()`

**Что добавляется в `handler.go` (ручная инструментация):**
```go
// В handler.Shorten — считаем созданные ссылки
shortenTotal.Inc()

// В handler.Redirect — считаем клики и измеряем время
redirectsTotal.Inc()
timer := prometheus.NewTimer(redirectDuration)
defer timer.ObserveDuration()
```

**Проверяем:**
```bash
curl http://localhost:8080/metrics | grep http_
# Prometheus: http://localhost:9090 → Graph → http_requests_total
```

**На что смотреть:**
- Prometheus → Status → Targets → url-shortener = **UP**
- `http_requests_total` — счётчик по методам и статусам
- `http_request_duration_seconds` — гистограмма latency
- В Grafana → Dashboards → Practicum → URL Shortener: первые графики оживают

**Проблема ручной инструментации:**
Каждый новый хэндлер нужно инструментировать вручную. Легко забыть.

---

### 🔧 Шаг 3 — Метрики через middleware (автоматически)

**Активировать:** найти `// step3 ` → заменить на ``, перебилдить.

**Что добавляется в `middleware/middleware.go`:**
```go
func Metrics(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        rw := newResponseWriter(w)    // захватываем status code
        next.ServeHTTP(rw, r)

        httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rw.status)).Inc()
        httpRequestDurationSeconds.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
    })
}
```

**Middleware подключён в `main.go`:**
```go
root = middleware.Metrics(root)  // оборачивает весь mux
```

**На что смотреть:**
- Метрики те же, но теперь **любой новый эндпоинт** получает их автоматически
- Попробуйте добавить хэндлер без единой строки метрик — данные всё равно появятся
- Grafana: `rate(http_requests_total[1m])` — живой трафик
- Grafana: `histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[1m]))` — p99 latency

**Паттерн:**
Middleware = единственное место для сквозной логики.
Всё что нужно сделать для всех запросов — делаем здесь.

---

### 📝 Шаг 4 — Структурированные логи (slog + Loki)

**Активировать:** найти `// step4 ` → заменить на ``, перебилдить.

**Что добавляется:**
1. В `main.go` — JSON handler для slog:
```go
logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
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

**Проверяем:**
```bash
# Логи в stdout контейнера
docker compose logs url-shortener -f

# В Grafana → Explore → Loki:
{service="url-shortener"}
{service="url-shortener"} | json | level="warn"
{service="url-shortener"} | json | path="/shorten"
```

**На что смотреть:**
- Каждая строка лога — валидный JSON с полями `time`, `level`, `msg`, `method`, `path`, `status`
- В Loki можно фильтровать по любому полю: `| json | status=500`
- Promtail автоматически читает логи из Docker и шлёт в Loki

**Почему JSON?**
Структурированные логи можно парсить, индексировать и искать.
Неструктурированный текст (`log.Printf`) в Loki — просто строка.

---

### 🔍 Шаг 5 — OpenTelemetry Tracing (ручные спаны)

**Активировать:** найти `// step5 ` → заменить на ``, перебилдить.

**Что добавляется:**

`internal/telemetry/telemetry.go` — инициализация TracerProvider:
```go
func InitTracer(ctx context.Context, serviceName, otlpEndpoint string) (func(context.Context) error, error) {
    // gRPC соединение к OTel Collector
    // OTLP exporter
    // Resource с именем сервиса
    // TracerProvider с BatchSpanProcessor
    // Регистрация глобального propagator (W3C TraceContext)
}
```

`main.go`:
```go
shutdown, err := telemetry.InitTracer(ctx, "url-shortener", otlpEndpoint)
defer func() { _ = shutdown(context.Background()) }()
```

`handler.go` — ручные спаны:
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

**Проверяем:**
```bash
# Создаём ссылку и делаем редирект
curl -X POST http://localhost:8080/shorten -d '{"url":"https://go.dev"}' -H 'Content-Type: application/json'
curl http://localhost:8080/<code>
```

В Grafana → Explore → Tempo:
- Выбрать `Search` → Service Name = `url-shortener`
- Видим трейс: один спан `handler.Redirect` с атрибутами

**На что смотреть:**
- Каждый span имеет `trace_id`, `span_id`, `parent_span_id`
- Атрибуты `url.code`, `url.original` видны в Tempo
- Ошибки в хэндлере → `span.RecordError(err)` → красный span в Tempo

**Вопрос студентам:** Что сейчас не хватает в трейсе? (Нет HTTP-уровня: method, url, status)

---

### ⚡ Шаг 6 — OTel HTTP Middleware + корреляция логов

**Активировать:** найти `// step6 ` → заменить на ``, перебилдить.

**Что добавляется:**

`main.go` — OTel HTTP middleware от contrib:
```go
root = otelhttp.NewHandler(root, "url-shortener-http",
    otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
)
```

`middleware/middleware.go` — добавляем trace_id в каждый лог:
```go
if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
    traceID = span.SpanContext().TraceID().String()
}
logger.InfoContext(ctx, "http request", ..., "trace_id", traceID)
```

**На что смотреть:**

1. **Tempo:** Теперь у каждого трейса есть корневой span `url-shortener-http`
   с атрибутами `http.method`, `http.url`, `http.status_code`, `http.flavor`

2. **Loki + Tempo correlation:**
   Grafana → Explore → Loki → `{service="url-shortener"} | json`
   Кликните на лог-строку → в поле `trace_id` появится кнопка **"View Trace in Tempo"**
   → прямой переход к трейсу!

3. **Grafana Dashboard:**
   Метрики + логи + трейсы в одном месте для одного запроса.

**Ключевой паттерн:**
`trace_id` в логах = мост между Loki и Tempo.
Один идентификатор связывает всё: метрики по времени, лог-строки, полный путь запроса.

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

## Итоговая архитектура observability

```
┌─────────────────────────────────────────────────────────┐
│                    Go Services                          │
│                                                         │
│  url-shortener ──HTTP──► stats-service                  │
│       │                       │                         │
│   Prometheus              Prometheus                    │
│   client_golang           client_golang                 │
│       │                       │                         │
│  OTel SDK (traces+logs)   OTel SDK (traces+logs)        │
└─────────┬─────────────────────┬───────────────────────--┘
          │ OTLP/gRPC           │ OTLP/gRPC
          ▼                     ▼
   ┌──────────────────────────────┐
   │      OTel Collector          │
   │  receiver: otlp              │
   │  processors: batch, memlimit │
   │  exporters:                  │
   │    → Tempo  (traces)         │
   │    → Loki   (logs)           │
   │    → Prometheus (metrics)    │
   └──────────────────────────────┘
          │
          ▼
   ┌──────────────────────────────┐
   │         Grafana              │
   │  ← Prometheus datasource     │
   │  ← Loki datasource           │
   │  ← Tempo datasource          │
   │                              │
   │  Dashboards + Explore +      │
   │  Cross-datasource linking    │
   └──────────────────────────────┘
```

---

## Полезные PromQL запросы

```promql
# RPS по методу и пути
rate(http_requests_total{job="url-shortener"}[1m])

# p99 latency по пути
histogram_quantile(0.99, sum by (path, le) (
  rate(http_request_duration_seconds_bucket{job="url-shortener"}[1m])
))

# Error rate (5xx)
rate(http_requests_total{job="url-shortener", status=~"5.."}[5m])
/ rate(http_requests_total{job="url-shortener"}[5m])
```

## Полезные LogQL запросы (Loki)

```logql
# Все логи сервиса
{service="url-shortener"}

# Только ошибки
{service="url-shortener"} | json | level="error"

# Медленные запросы (если добавить duration в лог)
{service="url-shortener"} | json | duration > 100ms

# Найти лог по trace_id
{service="url-shortener"} | json | trace_id="4bf92f3577b34da6a3ce929d0e0e4736"
```
