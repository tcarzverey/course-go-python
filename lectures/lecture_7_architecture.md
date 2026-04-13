---
marp: true
theme: default
paginate: true
backgroundColor: #fff
# backgroundImage: url('https://marp.app/assets/hero-background.svg')
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
    background: #1e1e1e;
    color: #d4d4d4;
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
---

<!-- _class: lead -->

# Архитектурные серверных приложений

## Устройство сервисов и межсистемное взаимодействие


---

<!-- _class: lead -->

# Архитектура

---

## Архитектура приложения vs Паттерны/Шаблоны проектирования

В чем разница?

---

# Архитектура приложения vs Паттерны/Шаблоны проектирования

<div class="columns">
<div>

## Архитектура
- набор решений о том, как модули приложения будут общаться друг с другом и с внешним миром
- решает проблемы проектирования всей системы в целом

</div>
<div>

## Паттерны
- шаблонное решение частой архитектурной проблемы
- более низкий уровень

</div>
</div>

---

# Зачем нужно думать про архитектуру?

<div class="columns">
<div>

## Проблемы
- Высокая связность кода
- Медленная разработка
- Сложность тестирования
- Страх что то сломать

</div>
<div>

## Решения
- Разделение ответственности
- Тестируемость
- Масштабируемость
- Понятность для команды

</div>
</div>
По сути платим сложностью сегодня, чтобы не платить хаосом завтра.

---

# Coupling & cohension

<div class="columns">
<div>

## Coupling
the degree of interdependence between software modules, a measure of how closely connected two routines or modules are, and the strength of the relationships between modules

</div>
<div>

## Cohension
the degree to which the elements inside a module belong together

</div>
</div>

---

# Coupling & cohension

<div class="columns-3">
<div>

![h:300](https://habrastorage.org/r/w1560/getpro/habr/upload_files/0d4/2f3/3f1/0d42f33f1aecaec248da08f3827d5509.png)

</div>
<div>

![h:300](https://habrastorage.org/r/w1560/getpro/habr/upload_files/569/df4/2c1/569df42c19542536b2d260e3c80cf51b.png)

</div>

<div>

![h:300](https://habrastorage.org/r/w1560/getpro/habr/upload_files/b52/1b6/c34/b521b6c34542baa53cc4c6f0ed3418a3.png)

</div>
</div>

---

# Coupling & cohension

![h:500](https://habrastorage.org/r/w1560/getpro/habr/upload_files/0fb/b88/5eb/0fbb885eb46596f0229f18752f7a52a1.png)

---

# Coupling & cohension

<div class="columns">
<div>

![h:500](https://habrastorage.org/r/w1560/getpro/habr/upload_files/2a5/c63/846/2a5c638460dd474a06a2e265fc7bc737.png)

</div>
<div>

![h:500](https://habrastorage.org/r/w1560/getpro/habr/upload_files/864/b7b/a9c/864b7ba9c89fe69bbba5918c317ee418.png)

</div>
</div>

---

# Эволюция архитектурных стилей

```
Monolith + MVC (1980+)
    ↓
Layered Architecture (1990+)
    ↓
DDD (2003+)
    ↓
Hexagonal/Onion/Clean (2005–2012)
    ↓
Microservices (2012+)
    ↓
Event-Driven / Serverless (2015+)
```

От tight coupling → к loose coupling
При этом растет сложность распределённых систем

---

<!-- _class: lead -->

# Архитектурные стили/подходы

---


# Монолит

<div class="columns">
<div>

## Плюсы
- Простота разработки, TTM
- Единая кодовая база, высокая связность кода
- Единая работа с данными, большой простор для переиспользования логики

</div>
<div>

## Проблемы (при росте)
- Высокая связность кода
- Медленная разработка
- Сложность тестирования
- Страх что-то сломать
- Масштабирование
- Сложности с деплоем

</div>
</div>

---

# MVC (Model-View-Controller)

![h:500](https://iq.opengenus.org/content/images/2020/06/MVC-3.png)

---

# MVC (Model-View-Controller)

<div class="columns">
<div>

## Когда использовать
- UI-heavy системы
- Быстрый прототип

</div>
<div>

## Проблемы
- Fat Controllers
- Tight coupling

</div>
</div>

---

# MVC coupling & cohension

<div class="columns-3">
<div>

```
app/
  models/
    user.go
    course.go
  controllers/
    user.go
    course.go
  views/
    user.go
    course.go
```

</div>
<div>

```
app/
  user/
    controller.go
    model.go
    store.go
    view.go
  course/
    controller.go
    model.go
    store.go
    view.go
```

</div>
<div>

```
app/
  user/
    controller/
      controller.go
    model/
      model.go
      store.go
    view/
      view.go
  course/
    controller/
      controller.go
    model/
      model.go
      store.go
    view/
      view.go
```

</div>
</div>

---

# MVC: Пример на Go

```go
// Model
type User struct {
    ID    int
    Name  string
    Email string
}

func (u *User) Validate() error { /* ... */ }
func GetUserByID(id int) (*User, error) { /* ... */ }

// Controller
func UserHandler(w http.ResponseWriter, r *http.Request) {
    id := getIDFromRequest(r)
    user, err := GetUserByID(id)
    if err != nil {
        http.Error(w, "Not found", 404)
        return
    }
    renderUserView(w, user)  // View
}
```

**Проблема:** бизнес-логика размазана между Model и Controller

---

# Layered Architecture (N-tier)

<div class="columns">
<div>

## Три классических слоя
```
┌─────────────────────┐
│  Presentation       │ HTTP, UI
├─────────────────────┤
│  Business Logic     │ Domain rules
├─────────────────────┤
│  Data Access        │ Database
└─────────────────────┘
```

**Правило:** Зависимости только вниз ↓

</div>
<div>

## Когда использовать
- CRUD приложения
- Простые бизнес-правила
- Типичные Enterprise системы

## Проблемы
- Coupling к БД
- Сложно тестировать
- Business logic "утекает" вверх

</div>
</div>

---

# Layered Architecture: Структура проекта

```go
project/
├── presentation/        // HTTP handlers, REST API
│   ├── handlers/
│   │   └── user_handler.go
│   └── dto/            // Data Transfer Objects
│       └── user_dto.go
├── business/           // Бизнес-логика
│   ├── services/
│   │   └── user_service.go
│   └── validators/
│       └── user_validator.go
└── data/              // Работа с БД
    ├── repositories/
    │   └── user_repository.go
    └── models/
        └── user_model.go
```

---

# Layered: Пример кода

```go
// Data Layer
type UserRepository struct { db *sql.DB }
func (r *UserRepository) FindByID(id int) (*User, error) { /* SQL */ }

// Business Layer
type UserService struct { repo *UserRepository }
func (s *UserService) GetUser(id int) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil { return nil, err }
    
    // Бизнес-логика
    if !user.IsActive { return nil, errors.New("user inactive") }
    return user, nil
}

// Presentation Layer
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    user, err := userService.GetUser(getID(r))
    json.NewEncoder(w).Encode(user)
}
```

---

# Проблема классических подходов

```go
// В Layered Architecture бизнес-логика часто зависит от БД:

type UserService struct {
    db *sql.DB  // ❌ Прямая зависимость от инфраструктуры!
}

func (s *UserService) CreateUser(name, email string) error {
    // Бизнес-правило смешано с SQL
    if len(name) < 3 {
        return errors.New("name too short")
    }
    
    _, err := s.db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", 
                        name, email)
    return err
}
```

**Проблема:** Невозможно протестировать без БД
**Решение:** Инверсия зависимостей → DDD, Hexagonal, Clean

---

#  Onion, Hexagonal, Clean Architecture

---

# Onion Architecture

<div class="columns">
<div>


![h:450](https://www.hosting.work/wp-content/uploads/2021/06/what-is-onion-architecture.png)

</div>
<div>

![h:450](https://habrastorage.org/r/w1560/getpro/habr/post_images/106/969/a45/106969a450369ea13161b15b438ec6f8.png)
</div>
</div>

**Jeffrey Palermo, 2008**

---

# Hexagonal Architecture
![h:500](https://habrastorage.org/r/w1560/getpro/habr/post_images/e41/7a6/cdb/e417a6cdb598c25f81851adf736f4006.png)

---

# Hexagonal Architecture (Ports & Adapters)

<div class="columns">
<div>

## Концепция
```
        ┌─────────────┐
    ┌───│   Adapter   │
    │   └─────────────┘
    │         ↓
┌───▼────────────────┐
│                    │
│   Domain Logic     │
│   (Hexagon)        │
│                    │
└───▲────────────────┘
    │         ↑
    │   ┌─────────────┐
    └───│   Adapter   │
        └─────────────┘
```

</div>
<div>

## Идея
- **Hexagon** = Domain Logic
- **Ports** = Интерфейсы
- **Adapters** = Реализации

**Inbound Ports:** HTTP, gRPC, CLI
**Outbound Ports:** DB, Queue, API

</div>
</div>

**Alistair Cockburn, 2005**

---

# Ports & Adapters: Детали

<div class="columns">
<div>

## Inbound Port
"Как извне вызвать бизнес-логику"

```go
// Port (интерфейс)
type OrderService interface {
    PlaceOrder(dto OrderDTO) error
}

// Adapter
type HTTPAdapter struct {
    service OrderService
}
func (h *HTTPAdapter) HandlePost(
    w http.ResponseWriter, 
    r *http.Request,
) { /* ... */ }
```

</div>
<div>

## Outbound Port
"Как бизнес-логика вызывает внешнее"

```go
// Port (интерфейс в domain)
type OrderRepository interface {
    Save(order Order) error
}

// Adapter (в infrastructure)
type PostgresOrderRepo struct {
    db *sql.DB
}
func (r *PostgresOrderRepo) 
    Save(order Order) error {
    // SQL
}
```

</div>
</div>

---

# Hexagonal
![](https://www.happycoders.eu/wp-content/uploads/2023/01/hexagonal-architecture-and-microservices.png)

---

# Hexagonal: Структура проекта

```go
project/
├── domain/                    // Hexagon (центр)
│   ├── order.go              // Entities
│   ├── order_service.go      // Business logic
│   └── ports/                // Интерфейсы
│       ├── order_repository.go    // Outbound port
│       └── order_use_case.go      // Inbound port
│
├── adapters/
│   ├── inbound/              // Входящие адаптеры
│   │   ├── http/
│   │   │   └── order_handler.go
│   │   └── grpc/
│   │       └── order_server.go
│   │
│   └── outbound/             // Исходящие адаптеры
│       ├── postgres/
│       │   └── order_repository.go
│       └── kafka/
│           └── event_publisher.go
```

---

# Hexagonal: Пример кода

```go
// domain/ports/order_repository.go (Outbound Port)
type OrderRepository interface {
    Save(order *Order) error
    FindByID(id string) (*Order, error)
}

// domain/order_service.go (Business Logic)
type OrderService struct {
    repo OrderRepository  // Зависимость от интерфейса
}

func (s *OrderService) PlaceOrder(items []Item) error {
    order := NewOrder(items)
    if err := order.Validate(); err != nil {
        return err
    }
    return s.repo.Save(order)  // Вызов через интерфейс
}
```

---

# Hexagonal: Адаптеры

```go
// adapters/inbound/http/handler.go (Inbound Adapter)
type OrderHandler struct {
    service *domain.OrderService  // Используем domain service
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    var dto OrderDTO
    json.NewDecoder(r.Body).Decode(&dto)
    
    err := h.service.PlaceOrder(dto.Items)  // Вызов domain
    // ...
}

// adapters/outbound/postgres/repository.go (Outbound Adapter)
type PostgresOrderRepository struct { db *sql.DB }

func (r *PostgresOrderRepository) Save(order *domain.Order) error {
    // SQL implementation
}
```

---

# Clean Architecture (Uncle Bob)

```
┌─────────────────────────────────────┐
│      Frameworks & Drivers           │  ← Web, DB, UI
│  ┌───────────────────────────────┐  │
│  │   Interface Adapters          │  │  ← Controllers, Gateways
│  │  ┌─────────────────────────┐  │  │
│  │  │   Use Cases             │  │  │  ← Application Logic
│  │  │  ┌───────────────────┐  │  │  │
│  │  │  │    Entities       │  │  │  │  ← Business Rules
│  │  │  └───────────────────┘  │  │  │
│  │  └─────────────────────────┘  │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘

Dependency Rule: зависимости направлены только внутрь →
```

**Robert Martin (Uncle Bob), 2012**

---

# Clean Architecture
![](https://habrastorage.org/r/w1560/web/fe8/c82/a32/fe8c82a32b1548b1a297187e24ae755a.png)

---

# Clean Architecture: Слои

<div class="columns">
<div>

## Entities (Центр)
Бизнес-правила уровня предприятия
```go
type Order struct {
    ID     string
    Total  Money
    Status Status
}

func (o *Order) Submit() error {
    // Правила бизнес-логики
}
```

</div>
<div>

## Use Cases
Application-специфичные правила
```go
type PlaceOrderUseCase struct {
    repo OrderRepository
}

func (uc *PlaceOrderUseCase) 
    Execute(req Request) error {
    order := NewOrder(req)
    order.Submit()
    return uc.repo.Save(order)
}
```

</div>
</div>

---

# Clean Architecture: Внешние слои

<div class="columns">
<div>

## Interface Adapters
Конвертация данных
```go
// Controller
type OrderController struct {
    useCase PlaceOrderUseCase
}

func (c *OrderController) 
    Post(w, r) {
    dto := parseDTO(r)
    request := toRequest(dto)
    c.useCase.Execute(request)
}
```

</div>
<div>

## Frameworks & Drivers
Технические детали
```go
// Repository implementation
type SqlOrderRepo struct {
    db *sql.DB
}

func (r *SqlOrderRepo) 
    Save(order Order) error {
    // SQL
}
```

</div>
</div>

**Dependency Rule:** Внешние слои зависят от внутренних, но не наоборот

---

# Clean Architechture data flow
![h:400](https://habrastorage.org/r/w1560/web/531/04c/89d/53104c89d9cf44a59c95e351b7485574.png)

---

# Clean Architecture: Структура

```go
project/
├── entities/                 // Самый внутренний слой
│   └── order.go
│
├── usecases/                // Application business rules
│   ├── place_order.go
│   └── ports/               // Interfaces
│       └── order_repository.go
│
├── adapters/                // Interface adapters
│   ├── controllers/
│   │   └── order_controller.go
│   ├── presenters/
│   │   └── order_presenter.go
│   └── gateways/
│       └── postgres_order_repo.go
│
└── frameworks/              // Frameworks & drivers
    └── web/
        └── server.go
```

---

# Domain-Driven Design (DDD)

---

# Domain-Driven Design: Философия

<div class="columns">
<div>

## Основная идея
**Бизнес-логика — это центр системы**

Код должен отражать:
- Язык бизнеса
- Реальные процессы
- Доменные правила

**Eric Evans, 2003**

</div>
<div>

## Ключевые принципы
1. **Ubiquitous Language**
   Единый язык для dev + business
   
2. **Bounded Context**
   Явные границы моделей
   
3. **Domain Model**
   Rich objects с поведением

</div>
</div>

---

# Слои DDD
![h:500](https://habrastorage.org/r/w1560/getpro/habr/upload_files/92a/df4/105/92adf41052435c6e65d8b8b58975adc8.jpg)

---

# DI
![h:500](https://habrastorage.org/r/w1560/getpro/habr/upload_files/dc7/3ca/a9a/dc73caa9a5afbe04723aee40aaa3b8aa.jpg)

---

# Ubiquitous Language

<div class="columns">
<div>

## Проблема
```go
// Разработчик:
type Rec struct {
    Amt float64
    Dt  time.Time
}

// Бизнес говорит:
"Invoice", "Amount", 
"Due Date"
```

</div>
<div>

## Решение
```go
// Используем язык домена:
type Invoice struct {
    Amount  Money
    DueDate time.Time
}

```

</div>
</div>

**Правило:** Если в коде другие термины, чем в разговоре — вы делаете это неправильно

---

# Bounded Context

```
E-commerce система:

┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   Catalog       │  │   Orders        │  │   Shipping      │
│                 │  │                 │  │                 │
│  Product:       │  │  Product:       │  │  Product:       │
│  - Name         │  │  - SKU          │  │  - Weight       │
│  - Description  │  │  - Price        │  │  - Dimensions   │
│  - Images       │  │  - Quantity     │  │  - Address      │
└─────────────────┘  └─────────────────┘  └─────────────────┘

     Разные модели "Product" в разных контекстах!
```

**Bounded Context** = граница, внутри которой модель имеет одно значение

---

# Bounded Contexts
![h:500](https://habrastorage.org/r/w1560/getpro/habr/upload_files/b36/003/1d9/b360031d9eaa620c01a3a0aedd8def2a.png)

---

# DDD: Building Blocks (Tactical Patterns)

<div class="columns">
<div>

### Entities
Объекты с идентичностью
```go
type Order struct {
    ID uuid.UUID // ← Identity
    CustomerID uuid.UUID
    Items []OrderItem
    Total Money
}
```
`Order{ID: 1} ≠ Order{ID: 2}`

</div>
<div>

### Value Objects
Объекты без идентичности
```go
type Money struct {
    Amount   decimal.Decimal
    Currency string
}
```
`Money{100, "USD"} == Money{100, "USD"}`

</div>
</div>

**Правило:** Если объекты равны по значениям → Value Object
Если нужно отслеживать изменения → Entity

---

# Entities vs Value Objects: Пример

```go
// Entity - имеет идентичность, изменяется
type Customer struct {
    ID      uuid.UUID  // Главное - ID
    Name    string     // Может меняться
    Email   Email      // Может меняться
    Address Address    // Может меняться
}

// Value Object - неизменяемый, без ID
type Email struct {
    value string
}

func NewEmail(s string) (Email, error) {
    if !isValidEmail(s) {
        return Email{}, errors.New("invalid email")
    }
    return Email{value: s}, nil
}

// Email сравниваются по значению, не по ссылке
```

---

# Aggregates и Aggregate Root

<div class="columns">
<div>

## Aggregate
Кластер связанных объектов, которые меняются вместе

```
Order (Root)
  ├── OrderItem[]
  │   ├── Item 1: Product A, qty 2
  │   └── Item 2: Product B, qty 1
  └── ShippingAddress
```

</div>
<div>

## Aggregate Root
Точка входа в aggregate
1. Внешние объекты держат ссылку только на Root
2. Root обеспечивает инварианты
3. Транзакция = 1 aggregate

</div>
</div>

```go
// ✅ ПРАВИЛЬНО: через Root
order.AddItem(product, quantity)

// ❌ НЕПРАВИЛЬНО: прямой доступ
order.Items = append(order.Items, item)
```

---

# Aggregate: Пример кода

```go
// Aggregate Root
type Order struct {
    id         OrderID
    customerID CustomerID
    items      []OrderItem    // private!
    status     OrderStatus
    total      Money
}

// Публичный метод для изменения
func (o *Order) AddItem(productID ProductID, quantity int, price Money) error {
    // Проверка инвариантов
    if o.status != OrderStatusDraft {
        return errors.New("cannot modify completed order")
    }
    
    item := OrderItem{ProductID: productID, Quantity: quantity, Price: price}
    o.items = append(o.items, item)
    o.recalculateTotal()  // Поддержка консистентности
    return nil
}

// Нельзя получить items напрямую - только через методы Root
```

---

# Domain Services vs Application Services

<div class="columns">
<div>

## Domain Service
Логика, которая не принадлежит Entity

```go
type PricingService struct{}

func (s *PricingService) 
    CalculateDiscount(
        order *Order, 
        customer *Customer,
    ) Money {
    // Сложная логика скидок
    // между Order и Customer
}
```

</div>
<div>

## Application Service
Оркестрация use case

```go
type OrderService struct {
    repo OrderRepository
    pricing PricingService
}

func (s *OrderService) 
    PlaceOrder(dto OrderDTO) error {
    order := s.buildOrder(dto)
    discount := s.pricing
        .CalculateDiscount(order)
    order.ApplyDiscount(discount)
    return s.repo.Save(order)
}
```

</div>
</div>

---

# Repositories

**Концепция:** Абстракция для получения/сохранения Aggregates

```go
// Интерфейс в Domain слое
type OrderRepository interface {
    FindByID(id OrderID) (*Order, error)
    Save(order *Order) error
    FindByCustomer(customerID CustomerID) ([]*Order, error)
}

// Реализация в Infrastructure слое
type PostgresOrderRepository struct {
    db *sql.DB
}

func (r *PostgresOrderRepository) FindByID(id OrderID) (*Order, error) {
    // SQL запрос, маппинг в domain объект
}
```

**Ключевое:** Domain не знает про SQL, только про интерфейс

---

# DDD: Структура проекта

```go
project/
├── domain/                    // Ядро системы - бизнес-логика
│   ├── order/
│   │   ├── order.go          // Aggregate Root
│   │   ├── order_item.go     // Entity
│   │   ├── money.go          // Value Object
│   │   ├── repository.go     // Interface
│   │   └── service.go        // Domain Service
│   └── customer/
│       └── ...
├── application/              // Use cases, оркестрация
│   └── order_service.go
├── infrastructure/           // Реализация технических деталей
│   ├── persistence/
│   │   └── postgres_order_repo.go
│   └── messaging/
└── interfaces/              // API, Controllers
    └── http/
        └── order_handler.go
```

---

# DDD: Полный пример

```go
// domain/order/order.go
type Order struct {
    id     OrderID
    items  []OrderItem
    status OrderStatus
}

func (o *Order) Submit() error {
    if len(o.items) == 0 {
        return errors.New("cannot submit empty order")
    }
    o.status = Submitted
    return nil
}

// domain/order/repository.go
type Repository interface {
    Save(order *Order) error
    FindByID(id OrderID) (*Order, error)
}
```

---

# DDD: Полный пример (продолжение)

```go
// application/order_service.go
type OrderService struct {
    repo domain.Repository
}

func (s *OrderService) PlaceOrder(items []Item) error {
    order := domain.NewOrder()
    for _, item := range items {
        order.AddItem(item.ProductID, item.Quantity)
    }
    
    if err := order.Submit(); err != nil {
        return err
    }
    
    return s.repo.Save(order)
}

// infrastructure/postgres_order_repo.go
type PostgresOrderRepository struct { db *sql.DB }
func (r *PostgresOrderRepository) Save(order *Order) error { /* SQL */ }
```

---

# Когда использовать DDD?

<div class="columns">
<div>

## Подходит для:
- Сложная бизнес-логика
- Много бизнес-правил
- Большая команда
- Долгоживущий проект
- Домен часто меняется

**Примеры:**
Финтех, страхование, e-commerce, logistics

</div>
<div>

## Не подходит для:
- CRUD приложения
- Простые правила
- Маленькая команда (1-2 dev)
- Прототип/MVP
- Данные > логика

**Примеры:**
Admin панели, content management, простые REST API

</div>
</div>

---

# Проблемы DDD

1. **Высокий порог входа**
   Много новых концепций, нужно время на освоение

2. **Overengineering риск**
   Легко переусложнить простую систему

3. **Требует опыта**
   Правильное выделение Bounded Contexts сложно

4. **Больше кода**
   Value Objects, Services, Repositories вместо простых CRUD

**Совет:** Начинайте с tactical patterns (Entity, VO, Repository), добавляйте strategic patterns (Bounded Context) по мере роста

---

# Практический выбор архитектуры

| Сложность домена | Размер проекта | Рекомендация |
|------------------|----------------|--------------|
| Простая (CRUD) | Маленький | **Layered** |
| Простая | Средний | **Layered** или **Hexagonal** |
| Средняя | Средний | **Hexagonal** |
| Средняя | Большой | **Clean** |
| Высокая | Любой | **DDD + Hexagonal/Clean** |
| Очень высокая | Большой | **DDD + Clean + CQRS** |

**Правило:** Не переусложняйте! Начните с простого, эволюционируйте по мере роста.

---

# Рефакторинг без боли

Начинаем постепенно
- Выносим use-cases из контроллеров
- Изолируем инфраструктуру через интерфейсы, DI
- Вводим слой маппинга DTO <-> домен
- По мере необходимости выделяем крупные сущности, агрегаты и инварианты
- По мере большой необходимости выделяем bounded contexts


---

# Микросервисы и распределённые системы

---

# От монолита к микросервисам

<div class="columns">
<div>

## Монолит
```
┌─────────────────────┐
│                     │
│   Orders            │
│   Products          │
│   Users             │
│   Payments          │
│   Shipping          │
│                     │
│   Single DB         │
└─────────────────────┘
```
Всё в одном процессе

</div>
<div>

## Микросервисы
```
┌────────┐  ┌─────────┐
│Orders  │  │Products │
│   DB   │  │   DB    │
└────────┘  └─────────┘

┌────────┐  ┌─────────┐
│Users   │  │Payments │
│   DB   │  │   DB    │
└────────┘  └─────────┘
```
Отдельные сервисы

</div>
</div>

---

# Почему все стали делать микросервисы?

<div class="columns">
<div>

## Проблемы монолита
- Долгий деплой всего приложения
- Scaling только вертикально
- Одна команда на всё
- Tight coupling между модулями
- Страшно что-то менять
- Technology lock-in

</div>
<div>

## Обещания микросервисов
- Независимый деплой
- Horizontal scaling
- Команды по сервисам
- Loose coupling
- Изолированные изменения
- Выбор технологий

</div>
</div>

**Netflix, Amazon, Uber** показали что это работает

---

# Реальные причины перехода

<div class="columns-3">
<div>

### Бизнес
- Быстрее фичи
- Масштабирование команд
- Отдельные релизы
- Изоляция рисков

</div>
<div>

### Технические
- Scaling по нагрузке
- Разные стеки
- Изоляция падений
- Независимое обновление

</div>
<div>

### Организационные
- Conway's Law
- Ownership сервисов
- Автономные команды
- Быстрая доставка

</div>
</div>

**Conway's Law:** Структура системы отражает структуру коммуникаций организации

---

# Как микросервисы влияют на архитектуру?

<div class="columns">
<div>

## DDD

</div>
<div>

## Hexagonal/Clean

</div>
</div>


---

# Как микросервисы влияют на архитектуру?

<div class="columns">
<div>

## DDD становится критичным
- Bounded Context = Микросервис
- Aggregate границы важны
- Ubiquitous Language внутри
- Контекст-маппинг между сервисами

</div>
<div>

## Hexagonal/Clean менее актуальны
- Сервисы уже изолированы
- Меньше coupling внутри
- Фокус на внешнюю интеграцию
- Но принципы остаются

</div>
</div>

---

# Проблемы распределённых систем

<div class="columns">
<div>

## Новые вызовы
- Network unreliability
- Latency
- Partial failures
- Data consistency
- Distributed transactions
- Monitoring/Debugging

</div>
<div>

## Старые проблемы × 10
- Сложность деплоя
- Версионирование API
- Testing интеграций
- Service discovery
- Load balancing
- Security

</div>
</div>

**Distributed systems are hard!**

---

# CAP теорема

<div class="columns">
<div>

![h:400](https://habrastorage.org/r/w1560/getpro/habr/upload_files/adf/b42/f13/adfb42f131f9b6109a581ebcf67a9cbd.png)

</div>
<div>

* consistency — данные актуальны и одинаковы в любой момент времени на каждом узле системы;

* availability — любой запрос к распределённой системе завершается откликом;

* partition tolerance — потеря связи между узлами не приводит к некорректности отклика от каждого из узлов.

</div>
</div>


---

# CAP: Практические примеры

<div class="columns">
<div>

## CP (Consistency + Partition)
Жертвуем Availability

**Примеры:**
- Банковские транзакции
- Inventory management
- Booking системы

**Если сеть упала:**
❌ Система недоступна
✅ Данные консистентны

</div>
<div>

## AP (Availability + Partition)
Жертвуем Consistency

**Примеры:**
- Social media feeds
- Shopping carts
- Analytics

**Если сеть упала:**
✅ Система доступна
❌ Eventual consistency

</div>
</div>

---

# BASE vs ACID

<div class="columns">
<div>

## ACID (Монолит)
- **A**tomicity
- **C**onsistency
- **I**solation
- **D**urability


</div>
<div>

## BASE (Микросервисы)
- **B**asically **A**vailable
- **S**oft state
- **E**ventual consistency


</div>
</div>

**В микросервисах:** Часто приходится принимать eventual consistency

---

<!-- _class: lead -->

# Паттерны взаимодействия

---

# Синхронное vs Асинхронное взаимодействие

<div class="columns">
<div>

```
Client → [Request]  → Service
       ← [Response] ←
```

**Характеристики:**
- Блокирующий вызов
- Immediate response
- Tight coupling
- Simple to understand

**Примеры:** REST, gRPC

</div>
<div>

```
Producer → [Message] → Queue
                         ↓
                     Consumer
```

**Характеристики:**
- Non-blocking
- Eventual processing
- Loose coupling
- Harder to debug

**Примеры:** Kafka, RabbitMQ

</div>
</div>

---

# Когда использовать синхронное?

<div class="columns">
<div>

## Подходит для
- User-facing операции
- Нужен немедленный ответ
- Simple workflows
- Read операции
- Низкая latency критична

**Пример:**
```
GET /users/123
→ Сразу нужен User
```

</div>
<div>

## Проблемы
- Cascading failures
- Higher latency
- Tight coupling
- Сложнее scaling
- Timeouts

**Антипаттерн:**
```
API → Service1 → Service2 →
    → Service3 → Service4
```


</div>
</div>

---

# Когда использовать асинхронное?

<div class="columns">
<div>

## Подходит для
- Background processing
- Event notifications
- High throughput
- Decoupling сервисов
- Retry logic нужна

**Пример:**
```
OrderCreated event
→ Email service
→ Warehouse service
→ Analytics service
```

</div>
<div>

## Проблемы
- Eventual consistency
- Сложная отладка
- Message ordering
- Duplicate messages
- Dead letter queues

**Паттерн:**
Idempotent consumers

</div>
</div>

---

# Event-Driven Architecture (EDA)

```
Service A                  Event Bus                 Service B
   │                           │                         │
   │──[UserCreated event]─────>│                         │
   │                           │────[UserCreated]───────>│
   │                           │                         │
   │                           │                     Service C
   │                           │                         │
   │                           │────[UserCreated]───────>│
```

**Принцип:** Сервисы реагируют на события, не вызывают друг друга напрямую

---

# Event-Driven Architecture: Преимущества

<div class="columns">
<div>

**Плюсы**
- Loose coupling
- Easy to add consumers
- Natural audit log
- Реактивность
- Масштабируемость

**Пример:**
Новый сервис просто подписывается

</div>
<div>

**Минусы**
- Eventual consistency
- Сложная отладка
- Event versioning
- Ordering проблемы

**Решения:**
- Event schema registry
- Correlation IDs
- Idempotency
- Dead letter queues

</div>
</div>

---

# Pub/Sub паттерн

<div class="columns">
<div>

## Концепция
```
Publisher     Topic      Subscriber
   │            │            │
   │─[event]───>│            │
   │            │──[event]──>│
   │            │            │
   │            │      Subscriber 2
   │            │            │
   │            │──[event]──>│
```

**Publisher** не знает о subscribers
**Subscribers** не знают друг о друге

</div>
<div>

## Варианты
**Topic-based:**
События по темам
`orders.created`, `users.updated`

**Content-based:**
Фильтры по содержимому
```json
{
  "type": "order",
  "amount": { "$gt": 1000 }
}
```

</div>
</div>

---

# Event Sourcing

<div class="columns">
<div>

## Традиционный подход
```sql
UPDATE accounts 
SET balance = 100 
WHERE id = 1;
```
Храним только текущее состояние

</div>
<div>

## Event Sourcing
```
Event 1: AccountOpened(100)
Event 2: MoneyDeposited(50)
Event 3: MoneyWithdrawn(30)
Event 4: MoneyDeposited(20)
---
Current state: 140
```
Храним все события

</div>
</div>

**Идея:** Состояние = сумма всех событий

---

# Event Sourcing: Преимущества и проблемы

<div class="columns">
<div>

**Плюсы**
- Полная история изменений
- Audit log из коробки
- Time travel (состояние на момент)
- Легко добавить проекции
- Event replay

**Use cases:**
- Финансы (аудит)
- Collaboration tools
- Analytics

</div>
<div>

**Минусы**
- Сложность queries
- Storage growth
- Schema evolution
- Eventually consistent
- Learning curve

**Решения:**
- Snapshots (периодические)
- CQRS (отдельные read модели)
- Event upcasting

</div>
</div>

---

# CQRS (Command Query Responsibility Segregation)

```
           Write Side                      Read Side
┌────────────────────────┐      ┌───────────────────────┐
│   Command              │      │   Query               │
│   ↓                    │      │   ↓                   │
│   Domain Model         │─────>│   Read Models         │
│   ↓                    │Event │   (Denormalized)      │
│   Write DB             │      │   ↓                   │
│   (Normalized)         │      │   Read DBs            │
└────────────────────────┘      └───────────────────────┘
```

**Идея:** Разные модели для записи и чтения

---

# CQRS: Зачем нужен?

<div class="columns">
<div>

## Проблемы одной модели
- Сложные запросы vs простая запись
- Нормализация vs денормализация
- Разная нагрузка (чтение >> запись)
- Разные требования к consistency

**Пример:**
Dashboard с агрегациями
vs
Simple CRUD операции

</div>
<div>

## CQRS решает
- Оптимизация чтения независимо
- Разные БД для read/write
- Scaling независимо
- Простые модели

**Архитектура:**
- Write: DDD Aggregates
- Read: Denormalized views
- Sync через события

</div>
</div>

---

# CQRS: Пример

```go
// Command Side
type CreateOrderCommand struct {
    CustomerID string
    Items      []Item
}

func (h *OrderCommandHandler) Handle(cmd CreateOrderCommand) error {
    order := domain.NewOrder(cmd.CustomerID, cmd.Items)
    h.repo.Save(order)
    h.events.Publish(OrderCreated{...})
    return nil
}

// Query Side
type OrderListQuery struct {
    CustomerID string
}

func (h *OrderQueryHandler) Handle(q OrderListQuery) ([]OrderDTO, error) {
    // Читаем из денормализованной read БД
    return h.readDB.FindOrdersByCustomer(q.CustomerID)
}
```

---

# Transactional Outbox

<div class="columns">
<div>

## Проблема
```go
func CreateOrder(order Order) error {
    tx.Begin()
    db.Save(order)
    tx.Commit()
    
    // ❌ Что если упадет?
    kafka.Publish(OrderCreated{...})
}
```

БД сохранилась, событие потерялось
или наоборот

</div>
<div>

## Решение: Outbox
```go
func CreateOrder(order Order) {
    tx.Begin()
    db.Save(order)
    db.SaveEvent(OutboxEvent{
        Type: "OrderCreated",
        Payload: {...}
    })
    tx.Commit()
}

// Отдельный процесс
func OutboxPublisher() {
    events := db.FetchUnpublished()
    for _, e := range events {
        kafka.Publish(e)
        db.MarkPublished(e.ID)
    }
}
```

</div>
</div>

---

# Transactional Outbox

```
┌─────────────────────────────────┐
│         Database (TX)           │
│  ┌──────────┐  ┌─────────────┐  │
│  │ Orders   │  │ Outbox      │  │
│  │ Table    │  │ Table       │  │
│  └──────────┘  └─────────────┘  │
└────────────────────┬────────────┘
                     │
              ┌──────▼──────┐
              │   Poller    │ (или CDC)
              └──────┬──────┘
                     │
              ┌──────▼──────┐
              │ Message Bus │
              └─────────────┘
```

**Гарантия:** Если запись в БД, то событие будет опубликовано

---

# API Gateway и Service Mesh

---

# API Gateway

![h:450](https://docs.solo.io/gloo-mesh-gateway/img/api-gateway.png)

**Единая точка входа** для всех клиентов

---

# API Gateway: Ответственности

<div class="columns">
<div>

## Что делает
- Authentication/Authorization
- Rate limiting
- Request/Response transformation
- Routing
- Load balancing
- Caching
- Logging/Monitoring
- SSL termination

</div>
<div>

## Паттерны
**Backend for Frontend (BFF)**
```
Mobile App  → Mobile Gateway
              → Service A, B, C

Web App     → Web Gateway
              → Service A, D, E
```

Разные gateway для разных клиентов

</div>
</div>

---

# API Gateway: Проблемы

<div class="columns">
<div>

## Single Point of Failure
Gateway падает = всё недоступно

**Решение:**
- Multiple instances
- Circuit breakers
- Fallbacks

</div>
<div>

## Latency
Дополнительный hop

**Решение:**
- Caching
- Connection pooling
- Async где возможно

</div>
</div>

### God Gateway антипаттерн
Не пихайте всю логику в gateway!

---

# Service Mesh

![h:400](https://securitypatterns.io/images/04-service-mesh/overview-service-mesh.png)

**Sidecar pattern:** Прокси рядом с каждым сервисом

---

# Service Mesh: Возможности

<div class="columns-3">
<div>

## Traffic Management
- Load balancing
- Circuit breaking
- Retries/Timeouts
- Traffic splitting (A/B)
- Fault injection

</div>
<div>

## Observability
- Distributed tracing
- Metrics collection
- Access logs
- Service graph

</div>
<div>

## Security
- mTLS between services
- Authorization policies
- Certificate management


</div>
</div>


**Популярные:** Istio, Linkerd, Consul Connect

---

<!-- _class: lead -->

# Распределённые транзакции

---

# Проблема распределённых транзакций

```
Service A        Service B        Service C
   │                 │                │
   │─[Create Order]──┤                │
   │                 │                │
   │                 │─[Reserve]──────┤
   │                 │                │
   │                 │◄──❌ Failed────┤
   │                 │                │
   │◄──❌ Rollback?──┤                │
```

**Проблема:** Нет общей транзакции, как откатить?

---

# Two-Phase Commit (2PC)

```
Coordinator
    │
    │──[Prepare]─────>│ Service A: ✅ OK
    │──[Prepare]─────>│ Service B: ✅ OK
    │──[Prepare]─────>│ Service C: ❌ FAIL
    │
    │──[Rollback]────>│ Service A
    │──[Rollback]────>│ Service B
```

**Фазы:**
1. **Prepare** - все готовы?
2. **Commit** (если все OK) или **Rollback**

---

# 2PC: Проблемы

<div class="columns">
<div>

## Проблемы
- Blocking protocol
- Coordinator SPOF
- Long-lived locks
- Не работает при partition
- Сложная реализация

**Coordinator упал?**
Участники ждут вечно

</div>
<div>

## Когда использовать
- Критичная consistency
- Небольшое число участников
- Контролируемая среда
- Низкая latency сети

**Примеры:**
- Банковские переводы
- Резервирование билетов

</div>
</div>

**В микросервисах:** Обычно избегают 2PC

---

# Saga Pattern

```
Choreography:
Service A ──[OrderCreated]──> Service B ──[PaymentProcessed]──>
Service C ──[ItemsReserved]──> Service D

Если ошибка:
Service D ──[ReservationFailed]──> Service C ──[CancelPayment]──>
Service B ──[CancelOrder]──> Service A
```

**Идея:** Последовательность локальных транзакций с compensating actions

---

# Saga: Choreography vs Orchestration

<div class="columns">
<div>

**Choreography**
Сервисы реагируют на события
```
Order ──[Created]──> Payment
Payment ──[Processed]──> Shipping
```

**Плюсы:**
- Loose coupling
- Простая добавка сервисов

**Минусы:**
- Сложно отследить flow
- Cyclic dependencies риск

</div>
<div>

**Orchestration**
Центральный координатор
```
Orchestrator:
  1. Create Order
  2. Process Payment
  3. Reserve Items
  4. Ship Order
```

**Плюсы:**
- Явный flow
- Легче debugging

**Минусы:**
- Coordinator coupling, SPOF риск

</div>
</div>

---

# Saga: Compensating Transactions

```go
// Happy path
CreateOrder()
ChargePayment()
ReserveInventory()
SendShipment()

// Failure на любом шаге
SendShipment() // FAILED
↓
ReleaseInventory()  // Compensate
RefundPayment()     // Compensate
CancelOrder()       // Compensate
```

**Важно:** Compensating != Rollback
- Rollback откатывает транзакцию
- Compensate создает новую компенсирующую транзакцию

---

# Saga: Пример кода (Orchestration)

```go
type OrderSaga struct {
    orderService     OrderService
    paymentService   PaymentService
    inventoryService InventoryService
}

func (s *OrderSaga) Execute(order Order) error {
    // Step 1
    orderID, err := s.orderService.CreateOrder(order)
    if err != nil {
        return err
    }
    
    // Step 2
    paymentID, err := s.paymentService.Charge(order.Total)
    if err != nil {
        s.orderService.CancelOrder(orderID)  // Compensate
        return err
    }
    
    // Step 3
    err = s.inventoryService.Reserve(order.Items)
    if err != nil {
        s.paymentService.Refund(paymentID)   // Compensate
        s.orderService.CancelOrder(orderID)  // Compensate
        return err
    }
    
    return nil
}
```

---

# Saga: Когда использовать?

<div class="columns">
<div>

## Подходит для
- Долгие бизнес-процессы
- Eventual consistency OK
- Распределённые транзакции
- Compensating actions возможны

**Use cases:**
- Order processing
- Booking systems
- Multi-step workflows

</div>
<div>

## Не подходит для
- Строгая consistency нужна
- Compensating невозможны
- Короткие транзакции

**Проблемы:**
- Сложность отладки
- Partial failures
- Нет isolation
- Compensating logic сложна

</div>
</div>

---

<!-- _class: lead -->

# Антипаттерны

---

# Distributed Monolith

```
Service A ─────> Service B ─────> Service C
    ↓                ↓                ↓
    └────────────────┴────────────────┘
              Shared Database
```

**Проблема:** Разделили код, но не данные и зависимости

**Признаки:**
- Shared database
- Синхронные вызовы цепочкой
- Deплой только вместе
- Изменение требует обновления всех

---

# Chatty APIs

```
Client                  Services
  │
  │────[GET /users/1]───────>│
  │◄──────[User data]────────│
  │
  │────[GET /orders/user/1]──>│
  │◄──────[Orders]────────────│
  │
  │────[GET /order/123]───────>│
  │◄──────[Order details]─────│
  │
  │────[GET /order/124]───────>│
  │◄──────[Order details]─────│
```

**Проблема:** N+1 queries в распределённой системе

**Решение:** Aggregation endpoints, GraphQL, BFF (apico)

---

# Timeout Cascade

```
Client [30s timeout]
   │
   └──> API Gateway [25s timeout]
           │
           └──> Service A [20s timeout]
                   │
                   └──> Service B [15s timeout]
                           │
                           └──> Database [10s timeout]
```

**Проблема:** Таймауты не согласованы, cascade failures

**Решение:** Уменьшать timeout на каждом уровне

---

# Retry Storm

```
Service A ──[Request]──> Service B (перегружен)
          ◄─[500 Error]─

Service A ──[Retry 1]──> Service B (ещё хуже)
Service A ──[Retry 2]──> Service B (умирает)
Service A ──[Retry 3]──> Service B (💀)
```

**Проблема:** Retries усугубляют проблему

**Решение:**
- Exponential backoff
- Circuit breaker
- Rate limiting

---

# God Service / God Gateway

```
        ┌────────────────┐
        │  Gateway       │
        │  - Auth        │
        │  - Validation  │
        │  - Business    │ ← ❌ НЕТ!
        │  - Transform   │
        │  - Aggregation │
        │  - Caching     │
        │  - etc...      │
        └────────────────┘
```

**Проблема:** Вся логика в одном месте

**Решение:** Thin gateway, логика в сервисах

---

# Shared Libraries Hell

```
Service A, B, C, D, E, F
        │
        └──> shared-library v1.2.3
               ├── models/
               ├── utils/
               ├── validators/
               └── business-logic/ ← ❌
```

**Проблема:** Coupling через библиотеки

**Что можно:**
- Utils, helpers
- DTO definitions
- Clients

---

# Ключевые выводы

## Микросервисы
- Не серебряная пуля
- Распределённость сложна
- CAP/BASE vs ACID
- Eventual consistency
- DDD помогает с границами
