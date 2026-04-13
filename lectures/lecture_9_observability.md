---
marp: true
theme: default
paginate: true
backgroundColor: #fff
style: |
  section {
    font-size: 28px;
  }
  h1 {
    color: #1E3A8A;
    font-size: 52px;
  }
  h2 {
    color: #3B82F6;
    font-size: 42px;
  }
  h3 {
    color: #60A5FA;
    font-size: 32px;
  }
  code {
    background: #f4f4f4;
    padding: 2px 6px;
    border-radius: 3px;
  }
  pre {
    background: #f8f8f8;
    color: #333;
    padding: 20px;
    border-radius: 8px;
    font-size: 18px;
  }
  .columns {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 1rem;
  }
  .columns-3 {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 1rem;
  }
  .practice {
    background: #f0fdf4;
    border-left: 4px solid #22c55e;
    padding: 16px 20px;
    border-radius: 6px;
    font-size: 26px;
  }
  .warning {
    background: #fff7ed;
    border-left: 4px solid #f97316;
    padding: 16px 20px;
    border-radius: 6px;
  }
---

<!-- _class: title-slide -->
<!-- _paginate: false -->

# Observability в распределённых системах

Мониторинг, трассировка, логирование — и почему это всё важнее, чем кажется

---

## Просыпаетесь в 3 ночи от звонка железной женщины, прод упал. Что делать?

---


## Observability

Observability решает именно это: способность понять внутреннее состояние системы по её внешним проявлениям — без необходимости знать заранее, что именно сломается.
Дополнительно - позволяет понять **связи** между информацией - логами, трейсами, метриками

---

## Monitoring vs Observability

| | Monitoring | Observability |
|---|---|---|
| Подход | Мы знаем, что может сломаться | Мы можем исследовать любое поведение |
| Данные | Предопределённые метрики | Метрики + трейсы + логи + контекст |
| Вопрос | Система работает? | Почему система ведёт себя именно так? ||

> Monitoring — это знать, что сервер упал. Observability — это понять, что запрос от пользователя X проходил через 7 сервисов, завис на 2 секунды в очереди RabbitMQ и упал именно на шаге валидации платежа.

---

## Откуда термин и почему он важен сейчас

**Control Theory (1960-е):** система наблюдаема, если её внутреннее состояние можно определить по внешним выходам.

**Почему сейчас это критично:**

- Монолит: один процесс, один лог, один стектрейс
- Микросервисы: запрос проходит через 10-20 сервисов, каждый пишет свои логи, у каждого своя база

```
User Request
    |
    v
[API Gateway] --> [Auth Service] --> [User Service]
                                          |
                                    [Order Service] --> [Payment Service]
                                          |                    |
                                    [Inventory]         [Notification]
```

---

# Три столпа Observability

???

---

# Три столпа Observability

## Metrics
## Logs
## Traces

---


## Три столпа: обзор

```
          OBSERVABILITY
       /       |       \
      /        |        \
  Metrics     Logs     Traces
    |          |          |
  "Что?"     "Детали"   "Где?"
  "Сколько?" "Почему?"  "Как долго?"

  CPU 85%   ERROR: DB    [req-123]
  RPS 1200  connection   api -> auth
  p99 450ms failed at    -> db (450ms)
            line 47      -> cache
```

Алерт по метрике говорит "что-то не так" → трейс показывает "где" → лог объясняет "почему".

---

## Технологии


| | Метрики | Логи | Трейсы |
|---|---|---|---|
| **Сбор / SDK** | OTel Metrics API | OTel Logs API, zap, slog | OTel Tracing API |
| **Хранение** | Prometheus, VictoriaMetrics | Loki, Elasticsearch | Jaeger, Grafana Tempo |
| **Визуализация** | Grafana | Grafana, Kibana | Grafana, Jaeger UI |
| **Протокол** | Prometheus scrape, OTLP | OTLP, Promtail | OTLP, Zipkin |
| **Запросы** | PromQL | LogQL, KQL | TraceQL |

**Почему так много вариантов?** Исторически каждый столп развивался независимо, каждый со своим стеком. OpenTelemetry стандартизирует сбор данных — а куда их хранить, выбираете сами.

---

## Метрики: что это и типы

Метрика — числовое измерение состояния системы в момент времени.

### Четыре базовых типа (Prometheus-модель)

| Тип | Описание | Пример |
|-----|----------|--------|
| **Counter** | Монотонно возрастает, никогда не убывает | Количество запросов, ошибок |
| **Gauge** | Произвольное значение вверх/вниз | Текущий CPU, RAM, активные соединения |
| **Histogram** | Распределение значений по бакетам | Latency запросов, размер payload |
| **Summary** | Перцентили на стороне клиента | p50, p95, p99 latency |

---

## Counter

```go
meter := otel.Meter("order-service")

requestCounter, _ := meter.Int64Counter(
    "http.requests.total",
    metric.WithDescription("Total HTTP requests"),
)

// Использование
func handleOrder(w http.ResponseWriter, r *http.Request) {
    requestCounter.Add(ctx, 1, metric.WithAttributes(
        attribute.String("method", r.Method),
        attribute.String("path", r.URL.Path),
        attribute.Int("status", 200),
    ))
}
```
```
Время:    10:00  10:01  10:02  restart  10:03  10:04
Counter:  1000   1200   1400     →       50     300
```
---

## Gauge

```go
activeConnections, _ := meter.Int64UpDownCounter(
    // UpDownCounter — это Gauge в терминах OTel
    "db.connections.active",
    metric.WithDescription("Active DB connections"),
)

func acquireConnection() {
    activeConnections.Add(ctx, +1)
}

func releaseConnection() {
    activeConnections.Add(ctx, -1)
}
```
---

## Observable gauge

```go
// Регистрируем callback — он вызывается при каждом сборе метрик
meter.Int64ObservableGauge(
    "process.memory.bytes",
    metric.WithInt64Callback(func(ctx context.Context, o metric.Int64Observer) error {
        var ms runtime.MemStats
        runtime.ReadMemStats(&ms)
        o.Observe(int64(ms.Alloc))
        return nil
    }),
)
```
---

## Histogram

```go
latencyHistogram, _ := meter.Float64Histogram(
    "http.request.duration",
    metric.WithUnit("ms"),
    metric.WithExplicitBucketBoundaries(
        // Настраивай бакеты под свои нужды!
        // Для API latency в ms:
        5, 10, 25, 50, 100, 200, 500, 1000, 2000,
    ),
)

func handleOrder(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    // ... обработка ...
    duration := float64(time.Since(start).Milliseconds())
    
    latencyHistogram.Record(ctx, duration, metric.WithAttributes(
        attribute.String("path", r.URL.Path),
        attribute.Int("status", statusCode),
    ))
}
```

---

## Histogram

```
Бакеты по умолчанию (в миллисекундах): [5, 10, 25, 50, 100, 250, 500, 1000, ∞]

Запрос 1: 37ms  → попадает в бакеты ≥50, ≥100, ≥250 ... ≥∞
Запрос 2: 82ms  → попадает в бакеты ≥100, ≥250 ... ≥∞  
Запрос 3: 6ms   → попадает в бакеты ≥10, ≥25 ... ≥∞

Итого в Prometheus хранится:
  latency_bucket{le="5"}    = 0
  latency_bucket{le="10"}   = 1   (только 6ms)
  latency_bucket{le="25"}   = 1
  latency_bucket{le="50"}   = 2   (6ms и 37ms)
  latency_bucket{le="100"}  = 3   (все три)
  latency_bucket{le="+Inf"} = 3
  latency_sum               = 125  (6+37+82)
  latency_count             = 3
```

---
## Summary

Автоматически считает на клиенте перцентили, как в histogram

```
3 реплики сервиса.
Каждая посчитала свой p99: 120ms, 95ms, 200ms

Histogram: Prometheus агрегирует бакеты со всех реплик → честный p99
Summary:   avg(120, 95, 200) = 138ms → неверно
```

---

## Метрики: RED и USE методологии

### RED (для сервисов — API, микросервисы)

- **Rate** — сколько запросов в секунду обрабатывает сервис
- **Errors** — какой процент запросов завершается ошибкой
- **Duration** — сколько времени занимает обработка

### USE (для ресурсов — CPU, память, диск, сеть)

- **Utilization** — процент времени, когда ресурс занят
- **Saturation** — насколько ресурс перегружен (очередь)
- **Errors** — количество ошибок ресурса

**Практика:** для каждого нового сервиса первым делом добавляют RED-метрики. Для каждого нового инфраструктурного компонента — USE.

--- 

## Cardinality

**Cardinality** — количество уникальных комбинаций значений labels.

```go
// ПРАВИЛЬНО
httpRequestsTotal.WithLabelValues("GET", "/api/orders", "200")
// Labels: method (5 значений) * endpoint (50 путей) * status (10) = 2500 time series

// ОПАСНО
httpRequestsTotal.WithLabelValues("GET", "/api/orders/12345", "200")
//                                                   ^^^^^ user ID в URL
// Каждый уникальный ID = новый time series
// 1M пользователей = 1M time series = OOM в Prometheus
```

**Правило:** labels должны быть из ограниченного множества, не используйте в labels: user_id, request_id, IP-адрес, произвольные строки.

**Нормализация URL:** `/api/orders/12345` → `/api/orders/:id`

---

## Логи

```
2024-01-15 10:23:45 ERROR Failed to process order 12345 for user john@example.com: timeout
```
---

## Логи: структурированное логирование

**Структурированный лог**:
```json
{
  "timestamp": "2024-01-15T10:23:45Z",
  "level": "error",
  "message": "Failed to process order",
  "order_id": "12345",
  "user_id": "usr_abc123",
  "error": "context deadline exceeded",
  "duration_ms": 5023,
  "service": "order-service",
  "trace_id": "4bf92f3577b34da6"
}
```

Структурированные логи можно фильтровать, агрегировать, индексировать и связывать с трейсами по `trace_id`.

---

## Логи в Go: slog

```go
// Стандартная библиотека Go 1.21+: slog
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))

// Добавление контекстных полей
logger.With(
    slog.String("service", "order-service"),
    slog.String("version", "1.2.3"),
).Error("Failed to process order",
    slog.String("order_id", orderID),
    slog.String("user_id", userID),
    slog.Duration("duration", elapsed),
    slog.Any("error", err),
)
```

---

## Логи в Go: zap

```go
// uber-go/zap (быстрее в ~10x), либо zerolog
import "go.uber.org/zap"

logger, _ := zap.NewProduction()
logger.Error("Failed to process order",
    zap.String("order_id", orderID),
    zap.String("trace_id", traceIDFromContext(ctx)),
    zap.Error(err),
)
```

---

## Log Levels: когда что использовать

| Уровень | Когда использовать | Пример |
|---------|-------------------|--------|
| **DEBUG** | Детали для разработчика, отключены в prod | Значения переменных, SQL-запросы |
| **INFO** | Нормальные бизнес-события | Заказ создан, пользователь вошёл |
| **WARN** | Необычно, но не критично | Retry 3 из 5, deprecated API |
| **ERROR** | Ошибка, требует внимания | DB недоступна, внешний API вернул 500 |
| **FATAL** | Сервис не может продолжить работу | Нет конфига, нет подключения к БД при старте |

---

## log.Error and ignore
```go
// логируем и продолжаем, как будто ничего не случилось
if err != nil {
    log.Printf("error: %v", err) // и дальше работаем с nil data
}
```

Каждый ERROR-лог должен либо завершать операцию, либо иметь явную стратегию восстановления.

---

## Distributed Tracing: как это работает

**Trace** — полный путь одного запроса через систему.
**Span** — одна операция внутри трейса (HTTP вызов, DB запрос, обработка сообщения).

```
Trace ID: 4bf92f3577b34da6
|
+-- [Span] api-gateway: POST /checkout          0ms - 523ms
    |
    +-- [Span] auth-service: ValidateToken      2ms - 15ms
    |
    +-- [Span] order-service: CreateOrder       20ms - 480ms
        |
        +-- [Span] postgres: INSERT orders      22ms - 45ms
        |
        +-- [Span] inventory-service: Reserve   50ms - 420ms  <-- ПРОБЛЕМА
            |
            +-- [Span] redis: GET stock         52ms - 53ms
            +-- [Span] postgres: UPDATE stock   55ms - 415ms  <-- узкое место
```

**Без трейсинга:** "checkout медленный, 523ms"
**С трейсингом:** "проблема в UPDATE stock в inventory-service, занимает 360ms"

--- 

## Span: анатомия

```
Span {
    TraceID:    4bf92f3577b34da6   // Общий для всего запроса
    SpanID:     a2fb4a1d1a96d312   // Уникальный для этого span
    ParentID:   b9c7c989f97918e1   // ID родительского span
    
    Name:       "order-service.CreateOrder"
    StartTime:  2024-01-15T10:23:45.020Z
    EndTime:    2024-01-15T10:23:45.480Z
    Duration:   460ms
    
    Attributes: {
        "http.method":      "POST",
        "http.url":         "/api/orders",
        "user.id":          "usr_abc123",
        "order.id":         "ord_xyz789",
        "db.system":        "postgresql",
    }
    
    Events: [                           // Аннотации внутри span
        {time: "+10ms", name: "validation_passed"},
        {time: "+20ms", name: "payment_initiated"},
    ]
    
    Status:     ERROR
    StatusMsg:  "context deadline exceeded"
}
```

---

# OpenTelemetry

Современный стандарт observability

---

## OpenTelemetry

**История проблемы:**
- OpenTracing (2016) — стандарт API для трейсинга
- OpenCensus (Google, 2018) — трейсинг + метрики
- Два несовместимых стандарта, библиотеки поддерживали оба

**Решение: OpenTelemetry (OTel, 2019)**
- Слияние OpenTracing + OpenCensus
- CNCF-проект (Cloud Native Computing Foundation)
- **Единый стандарт** для метрик, трейсов и логов
- Vendor-agnostic: один SDK, любой бэкенд

---

## OTEL 

```
Ваш сервис
    |
    | OTLP (OpenTelemetry Protocol)
    |
    v
[OTel Collector] --> Jaeger (трейсы)
                --> Prometheus (метрики)
                --> Loki (логи)
                --> Datadog (если нужен SaaS)
                --> Tempo (трейсы Grafana)
```

Vendor locking уходит на уровень коллектора, а не кода.

---

## OTel архитектура: SDK + Collector

```
+------------------+        +--------------------+        +-----------+
|   Наш сервис     |        |   OTel Collector   |        |  Backend  |
|                  |        |                    |        |           |
| [OTel SDK]       | OTLP   | [Receiver]         |        | Jaeger    |
|  - Tracer        |------->| [Processor]        |------->| Tempo     |
|  - Meter         | gRPC   |  - Batching        | OTLP   | Prometheus|
|  - Logger        | /HTTP  |  - Sampling        |        | Loki      |
|                  |        |  - Enrichment      |        | Datadog   |
+------------------+        | [Exporter]         |        +-----------+
                             +--------------------+
```

**Collector не обязателен** — SDK может экспортировать напрямую в бэкенд. Но коллектор даёт удобное использование:
- Батчинг и буферизацию (не теряем данные при пиках)
- Sampling (отбрасываем N% "скучных" трейсов)
- Обогащение данных (добавляем k8s labels, hostname)
- Маршрутизацию в несколько бэкендов одновременно

---

# Grafana Stack

Prometheus — Loki — Tempo — Grafana

---

<!-- _footer: "Observability | Grafana Stack" -->

## Grafana Stack: полный цикл

```
Метрики         Логи            Трейсы
    |               |               |
    v               v               v
[Prometheus]    [Loki]         [Tempo / Jaeger]
    |               |               |
    +-------+-------+---------------+
            |
            v
        [Grafana]
            |
    +-------+-------+
    |               |
[Dashboards]   [Alerting]
```

**Ключевое преимущество Grafana:** корреляция между сигналами. Нашёл аномалию на графике метрик → кликнул → видишь логи за тот же период → кликнул на trace_id → видишь весь путь запроса.

---

## Prometheus: архитектура и pull-модель

**Pull vs Push:**
- **Pull** (Prometheus по умолчанию): Prometheus сам ходит к сервисам за метриками (`/metrics`)
- **Push** (Pushgateway, OTLP): сервис сам отправляет метрики

---

## Loki: логи без индексации

**Elasticsearch:** индексирует всё содержимое → дорого, быстрый поиск по любому полю

**Loki:** индексирует только labels, хранит текст как есть → дёшево, поиск медленнее

```
Log entry:
  Labels (индексируются): {app="order-service", env="prod", level="error"}
  Content (не индексируется): {"order_id":"123","error":"timeout"}
```

**Loki + Promtail/Alloy:** агент на каждом узле собирает логи из файлов/docker и отправляет в Loki с нужными labels.

---

## Jaeger и Grafana Tempo: трейсы

**Jaeger** (CNCF, Uber):
- Самостоятельное хранилище трейсов
- UI для просмотра и поиска трейсов
- Поддерживает OTLP, Zipkin, Jaeger протоколы
- Бэкенды: in-memory (dev), Cassandra, Elasticsearch, Badger

**Grafana Tempo:**
- Интегрирован в Grafana экосистему
- Дешевле: хранит трейсы в object storage (S3, GCS)
- Нет собственного поиска — используется через Grafana
- Прямая корреляция с метриками и логами

---

## Sampling: не хранить всё подряд

В production с 10,000 RPS хранить каждый трейс слишком дорого.

### Стратегии sampling

| Стратегия | Описание | Когда использовать |
|-----------|----------|-------------------|
| **Head-based** | Решение принимается на входе трейса (apigw) | Простота, низкая задержка |
| **Tail-based** | Решение после завершения трейса | Можно оставить все ошибки |
| **Rate limiting** | N трейсов в секунду, не более | Предсказуемая нагрузка |
| **Probabilistic** | Случайный % трейсов | Общее представление о системе |

---

# SLO, SLA, Alerting

Как измерить надёжность и что делать когда всё плохо

---


## SLI, SLO, SLA: иерархия надёжности

**SLI (Service Level Indicator)** — конкретная измеримая метрика качества сервиса.
> "95-й перцентиль latency для /api/orders за последние 5 минут"

**SLO (Service Level Objective)** — целевое значение для SLI.
> "p99 latency < 500ms в 99.9% времени за последние 30 дней"

**SLA (Service Level Agreement)** — юридически обязывающий договор с последствиями за нарушение.
> "Если uptime < 99.9%, возвращаем 10% стоимости за месяц"

---

### Error Budget

```
SLO: 99.9% запросов успешны за 30 дней
Total requests: 10,000,000
Error budget = 10,000,000 * 0.001 = 10,000 failed requests
```

Когда error budget исчерпан:
- Останавливаем новые фичи
- Фокусируемся только на надёжности

---

## Alerting: алертить на симптомы, не на причины

**Плохой алерт:** "CPU > 80%"
- Причина: высокая нагрузка, это может быть нормально
- Создаёт alert fatigue — команда перестаёт реагировать

**Хороший алерт:** "Error rate > 1% в течение 5 минут"
- Симптом: пользователи получают ошибки прямо сейчас
- Требует немедленного действия

---

# Production практики

Что отличает хорошую observability от плохой

---

<!-- _footer: "Observability | Production" -->

## Correlation: связываем данные

```go
// Правильный подход: trace_id в каждом логе
func handler(w http.ResponseWriter, r *http.Request) {
    ctx, span := tracer.Start(r.Context(), "HandleRequest")
    defer span.End()

    // Получаем trace ID для логирования
    traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()

    logger.With(
        slog.String("trace_id", traceID),
        slog.String("span_id", span.SpanContext().SpanID().String()),
    ).Info("Processing request")

    // Добавляем trace ID в заголовок ответа — полезно для debugging
    w.Header().Set("X-Trace-Id", traceID)
}
```

**В Grafana:** когда видите аномалию на метриках → переходите к логам за этот период → находите `trace_id` → открываете трейс и видите весь путь запроса.

---

## Типичные ошибки

### Логировать слишком много или слишком мало
```go
log.Info("entering function")
log.Info("validating input")
log.Info("calling database")
```

### Игнорировать cardinality
```go
httpRequests.WithLabelValues(userID)
```

### Не добавлять trace_id в логи
```go
log.Error("payment failed", zap.Error(err))
```

---

## Observability-driven Development

**При проектировании новой фичи задайте вопросы:**
- Как я узнаю, что фича работает правильно?
- Как я узнаю, что фича не работает?
- Как я найду проблему, если она возникнет в 3 ночи?

**Чеклист перед деплоем:**
- Метрики для RED-показателей добавлены?
- Критические операции обёрнуты в spans?
- Логи структурированы, trace_id прокидывается?
- Алерты на error rate и latency настроены?
- Dashboard в Grafana создан?

---

# Практика