# НИС Построение сервисов на Golang и Python

## Материалы
### Занятие 1: Основы Go
Темы: концепция и основные характеристики языка, решаемые проблемы, сравнение с другими языками, базовые конструкции языка, массивы, слайсы, мапы, интерфейсы, структуры.  
* [Презентация](https://docs.google.com/presentation/d/1ikFbcWoQrE7Je3kHW3RIEpvAP1cwElM_IV1fqCIkbVg/edit?usp=sharing)
* [Полезные материалы](./materials/lecture1.md)
* Домашнее задание: пройти [tour of Go](https://go.dev/tour/welcome/1), кроме concurrency и generics

### Занятие 2: Продолжение Go, concurrency
Темы: обработка ошибок, асинхронность vs конкурентность vs параллельность vs многопоточность, concurrency: горутины, каналы, примитивы синхронизации, go concurrency patterns, кратко про scheduler
* [Презентация](https://docs.google.com/presentation/d/1e_xCM6JOcFC5n87-sU_yl5rCoX1foQHASNJWqXIZYLY/edit?usp=sharing)
* [Полезные материалы](./materials/lecture2.md)
* [Домашнее задание](./homeworks/hw1/README.md)

### Занятие 3: API, Протоколы
Темы: TCP/IP, SOAP, HTTP, REST, gRPC, GraphQL
* [Презентация](https://docs.google.com/presentation/d/1BtnP6m6C5pz7rfcYV-9K9gSf0hR-57EajhOYKyMdGco/edit?usp=drive_link)
* [Полезные материалы](./materials/lecture3.md)
* [Запись](https://disk.yandex.ru/i/45n1XxvskJgG2Q)
* [Домашнее задание](./homeworks/hw2/README.md)

## Занятие 4: БД и хранилища данных, go-библиотеки
Темы: RmDBS, NoSQL, Redis, Brokers & Queues
* [Презентация](https://docs.google.com/presentation/d/11_gxZ5K6qqGjN_RsiUYLDfl9a27IT2AfpuPuoTyYubo/edit?usp=sharing)
* [Запись](https://disk.yandex.ru/i/Bah7_8ZQY1dXkA)
* [Полезные материалы](./materials/lecture4.md)

## Занятие 5-6: Практикум: строим сервис с нуля
Темы: OpenAPI, кодогенерация, sqlc, unit & integration тестирование, моки, (чуть-чуть) pprof
* [Проект](./examples/lecture5)
* Записи: [1](https://disk.yandex.ru/i/ngUgsks6U5BZhw) и [2](https://disk.yandex.ru/i/OI3LW3UqHdHBRA)
* Полезные материалы: [1](./materials/lecture5.md) и [2](./materials/lecture6.md)

## Занятие 7: Архитектура & Паттерны
Темы: Coupling & cohension, Виды слоистых архитектур, DDD, Микросервисы vs Монолит, CAP, sync/async, CQRS, Event Sourcing, TxOutbox, ApiGW, Mesh, 2PC vs Saga
* [Презентация](./lectures/lecture_7_architecture.md)
* [Запись](https://disk.yandex.ru/i/F_WFO6f8ouOzVw)

## Занятие 8: AuthN & AuthZ + Тест
Темы: Идентификация, Авторизация, Аутентификация, JWT, OAuth2, SSO, ACL/RBAC/ABAC/PBAC, IDC
* [Презентация](./lectures/lecture_8_auth.md)
* [Запись](https://disk.yandex.ru/i/2KHUKngEApUPUg)

## Занятие 9: Observability + Практикум
Темы: Logging, Tracing, Metrics, OTEL, SLI/SLO/SLA
* [Презентация](./lectures/lecture_9_observability.md)
* [Проект](./examples/lecture9)
* [Запись](https://disk.yandex.ru/i/4nguNQPJbkNnEA)

## Занятие 10: Go vs Python
Темы: История Python, Пакетные менеджеры, GIL, Asyncio, Exceptions, Типизация, Django/FastAPI и другие фреймворки
* [Презентация](./lectures/lecture_10_python.md)
* [Полезные материалы](./materials/lecture10.md)
* [Запись](https://disk.yandex.ru/i/DEfRfgOcwu-GFQ)