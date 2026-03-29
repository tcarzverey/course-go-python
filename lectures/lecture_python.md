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

# Go и Python

В чём реальная разница?

---

## Что разберём сегодня

- Почему Go и Python решают одни задачи по-разному
- Почему Python медленный — и когда это не важно
- Concurrency: GIL, asyncio, goroutines
- Обработка ошибок: исключения vs явный return
- Типизация: type hints, mypy, pydantic vs компилятор Go
- Ловушки Python, которые стреляют в production
- Веб-фреймворки: FastAPI, Django, Flask
- Когда брать Go, когда Python

---

  ## История Python

  - **1989** — Гвидо ван Россум начал разработку как хобби-проект на Рождество
  - **1991** — Python 0.9.0: первый публичный релиз. Классы, исключения, функции
  - **2000** — Python 2.0: list comprehensions, cycle GC, unicode (частично)
  - **2008** — Python 3.0: намеренный разрыв обратной совместимости. Начало 12-летней миграции
  - **2010** — Python 2.7: последний релиз ветки 2.x, LTS до 2020
  - **2018** — Гвидо уходит с роли BDFL («Великодушный пожизненный диктатор»). Управление переходит к Steering Council
  - **2020** — Python 2 официально EOL. До этого момента — два несовместимых мира

Почему называется Python?

---

<!-- ## Ключевые версии Python 3

  | Версия | Год | Что добавили |
  |--------|-----|--------------|
  | **3.5** | 2015 | `async`/`await`, type hints (PEP 484) |
  | **3.6** | 2016 | f-строки, `__future__` annotations |
  | **3.7** | 2018 | `dataclasses`, `breakpoint()`, гарантированный порядок dict |
  | **3.8** | 2019 | Walrus-оператор `:=`, positional-only параметры |
  | **3.10** | 2021 | `match`/`case` (pattern matching), улучшенные сообщения об ошибках |
  | **3.11** | 2022 | +10–60% скорости, точные трейсбеки с указанием выражения |
  | **3.12** | 2023 | Новый синтаксис generics `[T]`, f-строки без ограничений |
  | **3.13** | 2024 | Экспериментальный **no-GIL режим** (PEP 703), free-threaded Python |

--- -->

## Контекст: оба языка в backend

<div class="columns">

```
Go — где используется:
• Системный софт, инфраструктура
• Kubernetes, Docker, Prometheus
• Высоконагруженные API
• Финтех, стриминг, телеком
• Микросервисы с жёсткими SLA
```

```
Python — где используется:
• ML/AI — де-факто стандарт
• Data Engineering (Airflow, Spark)
• Веб (Instagram, Pinterest, Dropbox)
• Скрипты, автоматизация, DevOps
• Прототипирование и исследования
```

</div>

---

# Go и Python

---

## Что Go решал

| Проблема | Решение в Go |
|----------|-------------|
| Долгая компиляция C++ | Компиляция за секунды |
| Сложная асинхронность | Goroutines + channels — конкурентность из коробки |
| Отсутствие единого тулинга | `go fmt`, `go test`, `go build` — один инструмент |
| Dependency hell | Go modules, vendor, воспроизводимые сборки |
| Разный кодстайл в команде | `gofmt` — единый стиль, нет споров |
| Читаемость больших кодовых баз | Минимализм языка, явность, нет магии |
| Управление памятью | GC без пауз (< 1ms), без ручного free |

---

## Как Python справляется с теми же задачами

| Проблема | Ситуация в Python |
|----------|------------------|
| Компиляция | Интерпретатор — нет шага компиляции, но нет и раннего обнаружения ошибок |
| Асинхронность | `asyncio` — мощный, но требует дисциплины. GIL ограничивает потоки |
| Тулинг | Зоопарк: pip/uv/poetry, black/ruff, pytest/unittest — выбор сложен |
| Зависимости | virtualenv, venv, conda — много решений, нет одного стандарта |
| Кодстайл | `black`, `ruff` — хорошие инструменты, но опциональны |
| Читаемость | Синтаксис ясный, но динамика усложняет понимание |
| Память | Reference counting + cyclic GC. Нет ручного контроля |

---

## Особенности Python


```
Python
• Duck typing без интерфейсов
  (Protocol появился в 3.8)
• Наследование есть (и MRO)
• Мультипарадигменный
  (ООП + функциональный)
• Нет указателей
• try/except — основной механизм работы с ошибками
• pytest — внешняя библиотека
• Generics через TypeVar (3.12+)
• Потоки дорогие, корутины дёшевы
```

---

## pip + requirements.txt

Стандарт из коробки. Просто список пакетов — никакой автоматики.

```text
# requirements.txt
fastapi==0.111.0
uvicorn==0.30.1
sqlalchemy>=2.0,<3.0
pydantic==2.7.1
```

Проблемы:
- pip freeze дампит транзитивные зависимости вместе с прямыми — не видно что откуда
- Нет разделения dev/prod в одном файле
- Нет встроенной изоляции (нужен venv отдельно)
- Воспроизводимость только если все версии точно закреплены

---
## pipenv

Первая попытка "всё в одном": менеджер зависимостей + виртуальное окружение.

```toml
# Pipfile — прямые зависимости, читаемый
[[source]]
url = "https://pypi.org/simple"

[packages]
fastapi = "==0.111.0"
sqlalchemy = ">=2.0,<3.0"

[dev-packages]
pytest = "*"
black = "*"

[requires]
python_version = "3.12"
```

Pipfile.lock — генерируется автоматически, хранит все транзитивные зависимости с SHA-хешами.

<!-- 
Был популярен 2018–2022. Сейчас теряет позиции из-за медленного резолвера и стагнации разработки. -->

---

## poetry

Современный стандарт для серьёзных проектов. Один файл для всего.

```toml
# pyproject.toml
[tool.poetry]
name = "my-service"
requires-python = "^3.12"

[tool.poetry.dependencies]
fastapi = "^0.111.1" # [0.111.1; 0.112.0)
pydantic = "^2.7" # [2.7; 3.0)

[tool.poetry.group.dev.dependencies]
fastapi = "~0.111.1" # [0.111.0; 0.112.0)
sqlalchemy = "~2.0" # [2.0; 2.1)
pydantic = "==2.7"
numpy = "*"
```

poetry.lock — полный граф зависимостей с хешами, коммитится в репозиторий.

---
## conda

Не только Python пакеты — умеет ставить системные бинарники: CUDA, MKL, компиляторы.

```yaml
# environment.yml
name: ml-project
channels:
- pytorch
- conda-forge
- defaults

dependencies:
- python=3.11
- cudatoolkit=11.8      # системная зависимость, не Python
- pytorch=2.3.0
- numpy=1.26
- pandas=2.2
- pip:
    - fastapi==0.111.0    # pip-пакеты внутри conda-окружения
```

Незаменим в ML/data science где нужна конкретная версия CUDA.

---

## uv

Написан на Rust. В 10–100x быстрее pip. Управляет и пакетами, и версиями Python.

```toml
# pyproject.toml — тот же формат что и poetry
[project]
name = "my-service"
version = "0.1.0"
requires-python = ">=3.12"
dependencies = [
    "fastapi==0.111.0",
    "sqlalchemy>=2.0,<3.0",
    "pydantic>=2.7",
]

[dependency-groups]
dev = [
    "pytest>=8.2",
    "black>=24.4",
    "mypy>=1.10",
]
```
uv.lock — автоматический lockfile. Тренд 2024–2025, многие проекты мигрируют с poetry.

---

## Python dependency hell

В одном окружении может быть установлена **только одна версия** пакета.

Проект зависит от:
```
library-a → требует requests>=2.28
library-b → требует requests<2.25
```

```bash
pip install library-a library-b
ERROR: Cannot install library-a and library-b
because these package versions have conflicting dependencies.
```

---

## Причины

- Авторы библиотек не соблюдают строго semver — breaking changes попадают в минорные версии
- pip до версии 20.3 (2020!) вообще не имел настоящего backtracking-резолвера — устанавливал первое попавшееся
- Нет механизма "пусть сосуществуют две версии одного пакета"

**Обходы:** виртуальные окружения изолируют проекты друг от друга, но внутри одного проекта конфликт неразрешим.

---

## Совместимость зависимостей с версий Python

* Хочу использовать новую фичу Python 3.12, но numpy 1.x не поддерживает 3.12 → жду пока numpy добавит поддержку

* Корпоративный сервер застрял на Python 3.8, FastAPI последних версий требует 3.9+ → либо старая версия FastAPI, либо обновляй Python
* torch, tensorflow, numpy — большие библиотеки с C-расширениями, (.whl файл) для Python 3.10 не установится на 3.11 — нужна пересборка.


---

## Go modules

Go использует **Minimum Version Selection (MVS)** — алгоритм Расса Кокса.

**Правило:** если два пакета требуют разные версии одной зависимости — берём **максимум из минимумов**.

lib-a требует: logger v1.2+
lib-b требует: logger v1.5+

Go выбирает: logger v1.5  ← минимальная версия, которая удовлетворяет всем

---

## Go semver

Это работает потому, что Go **настаивает** на соблюдении semver:

v1.x → обратная совместимость гарантирована компанией/инструментами
v2.x → это уже ДРУГОЙ модуль с другим import path:

```go
import "github.com/foo/bar"    // v1
import "github.com/foo/bar/v2" // v2 — живут рядом, не конфликтуют
```

**go.sum** —  каждая зависимость зафиксирована с хешем содержимого, гарантирована воспроизводимая сборка.


---

## Означает ли это, что в Go невозможно сломаться по нарушению обратной совместимости?


---

# Почему Python медленный

Или нет?

---

## Интерпретируемость: от кода до выполнения

```
Go:
  source.go  -->  [gc compiler]  -->  machine code  -->  CPU выполняет напрямую
                                      (ELF binary)

Python:
  source.py  -->  [py_compile]  -->  bytecode (.pyc)  -->  [CPython VM]  -->  CPU
                                                           интерпретирует
                                                           каждую инструкцию
```

**Что это означает:**
- Каждая операция в Python — вызов C-функции в CPython
- Go выполняет машинный код напрямую
- Типичный разрыв: Go в **10-100x** быстрее CPython на CPU-bound задачах
- На I/O-bound задачах (ждём сеть/диск) разрыв **гораздо меньше**

---

## Dynamic typing overhead

В Python каждая операция — поиск типа в runtime.

<div class="columns">

```go
// Go: компилятор знает типы
// a + b = одна машинная инструкция ADD
var a, b int = 3, 4
c := a + b
```

```python
# Python: runtime dispatch
a, b = 3, 4
c = a + b
# За кулисами:
# 1. type(a) -> int
# 2. int.__add__ существует?
# 3. вызываем int.__add__(a, b)
# 4. создаём новый объект int
```

</div>

**Каждое число в Python** — объект на куче: `ob_refcnt` + `ob_type` + `ob_digit`.
Число `42` в Go — 8 байт на стеке. Число `42` в Python — ~28 байт на куче.

---

## Альтернативные рантаймы Python

| Рантайм | Подход | Ускорение | Когда использовать |
|---------|--------|-----------|-------------------|
| **CPython** | Интерпретатор (эталон) | 1x | По умолчанию |
| **PyPy** | JIT-компиляция | 5-10x | Долгоживущие процессы, CPU-bound |
| **Cython** | Python → C | 10-100x | Горячие пути, есть C-экспертиза |
| **Numba** | JIT для NumPy | ~C скорость | Числовые вычисления, массивы |
| **GraalPython** | JVM, Truffle | 2-5x | Polyglot, Java экосистема |

**Ключевой вывод:** стандартный Python медленный на CPU. Но большинство backend-задач ограничены I/O, а не CPU — там разрыв незначителен.

---

# Конкурентность

GIL, asyncio, goroutines — три разных взгляда на одну задачу

---

## Проблема: один процессор, много задач

```
Задача: обработать 1000 HTTP-запросов одновременно

Подходы:
┌─────────────────────────────────────────────────────────┐
│  Процессы     │  Потоки (OS threads)  │  Корутины/Green │
│               │                       │  threads        │
│  Изоляция     │  Shared memory        │  Cooperative    │
│  Дорого (MB)  │  Дешевле (MB стека)   │  Дёшево (KB)    │
│  IPC сложно   │  Синхронизация сложна │  Один поток OS  │
└─────────────────────────────────────────────────────────┘

Go:    goroutines (~2KB) + планировщик GMP
Python: OS threads + GIL или asyncio корутины
```

---

## Python threading: OS потоки и GIL

```python
import threading

results = []

def fetch(url):
    response = requests.get(url)  # блокирующий вызов
    results.append(response.status_code)

threads = [threading.Thread(target=fetch, args=(url,)) for url in urls]
for t in threads: t.start()
for t in threads: t.join()
```

**GIL (Global Interpreter Lock)** — мьютекс, который позволяет только одному потоку выполнять байткод Python в каждый момент времени.

> Потоки создаются как OS threads, но Python код выполняется по очереди — не параллельно.

---

## GIL: почему он существует

```
Проблема без GIL:
  Thread 1: ob_refcnt(obj) = 5  →  читает: 5
  Thread 2: ob_refcnt(obj) = 5  →  читает: 5
  Thread 1: записывает: 6
  Thread 2: записывает: 6       ← должно быть 7!
  → утечка памяти / двойное освобождение

GIL решение:
  Только один поток выполняет Python байткод
  → reference counting всегда корректен
  → C-расширения безопасны без дополнительных блокировок
```

**Следствие:**
- I/O-bound: потоки работают — GIL освобождается во время syscall
- CPU-bound: потоки **не** дают ускорения — GIL не освобождается

---

## GIL workarounds: multiprocessing

```python
from concurrent.futures import ProcessPoolExecutor
import os

def cpu_intensive(n):
    # Считаем в отдельном процессе — у каждого свой GIL
    return sum(i * i for i in range(n))

# Каждый процесс = отдельный Python интерпретатор
with ProcessPoolExecutor(max_workers=os.cpu_count()) as executor:
    futures = [executor.submit(cpu_intensive, 10_000_000) for _ in range(8)]
    results = [f.result() for f in futures]
```

**Цена:**
- Каждый процесс ~50MB RAM (отдельный интерпретатор)
- Данные между процессами — через сериализацию (pickle)
- Запуск процесса: ~100ms (vs ~1μs для горутины)

---

## asyncio: event loop модель

```
asyncio event loop:
┌──────────────────────────────────────────────────┐
│                  Event Loop                      │
│                                                  │
│  ┌──────────┐   await   ┌──────────────────────┐ │
│  │ coroutine│ ────────> │   I/O selector       │ │
│  │ (task 1) │           │   (epoll/kqueue)     │ │
│  └──────────┘ <──────── │                      │ │
│                ready    │  Ждёт готовности:    │ │
│  ┌──────────┐           │  - сокетов           │ │
│  │ coroutine│           │  - таймеров          │ │
│  │ (task 2) │           │  - файлов            │ │
│  └──────────┘           └──────────────────────┘ │
└──────────────────────────────────────────────────┘
Один OS поток. Переключение — только в точках await.
```

---

## asyncio: async/await на практике

```python
import asyncio
import aiohttp

async def fetch(session, url):
    async with session.get(url) as response:
        return await response.text()

async def fetch_all(urls):
    async with aiohttp.ClientSession() as session:
        # gather запускает все корутины "одновременно"
        tasks = [fetch(session, url) for url in urls]
        results = await asyncio.gather(*tasks)
    return results

# Запуск event loop
results = asyncio.run(fetch_all(urls))
```

`asyncio.gather` — аналог `errgroup` в Go: запускает все задачи, ждёт завершения всех.

---

## Ограничения asyncio: один блокирующий вызов = всё стоит

```python
import asyncio
import time

async def fast_task():
    await asyncio.sleep(0.1)  # OK: освобождает event loop
    return "fast"

async def blocking_task():
    time.sleep(5)  # ПЛОХО: блокирует весь event loop
    return "blocked everything"

async def main():
    # blocking_task заморозит fast_task на 5 секунд
    results = await asyncio.gather(fast_task(), blocking_task())
```

<div class="warning">

Любой блокирующий синхронный вызов — блокирует весь event loop. 
Используйте `loop.run_in_executor()` для такого.

</div>

---

## Параллельный fetch N URL

<div class="columns">

```python
# Python asyncio
import asyncio
import aiohttp

async def fetch(session, url):
    async with session.get(url) as r:
        return r.status

async def main(urls):
    async with aiohttp.ClientSession() as s:
        tasks = [fetch(s, u) for u in urls]
        return await asyncio.gather(*tasks)

results = asyncio.run(main(urls))
```

```go
// Go goroutines
func fetch(url string) int {
    resp, err := http.Get(url)
    if err != nil { return 0 }
    defer resp.Body.Close()
    return resp.StatusCode
}

func main() {
    results := make([]int, len(urls))
    var wg sync.WaitGroup
    for i, url := range urls {
        wg.Add(1)
        go func(i int, url string) {
            defer wg.Done()
            results[i] = fetch(url)
        }(i, url)
    }
    wg.Wait()
}
```

</div>

---

## Channels vs asyncio.Queue

<div class="columns">


```python
# asyncio.Queue
import asyncio

q = asyncio.Queue(maxsize=10)

async def producer():
    for i in range(100):
        await q.put(i)

async def consumer():
    while True:
        v = await q.get()
        process(v)
        q.task_done()
```


```go
// Go channels
ch := make(chan int, 10)

// Producer
go func() {
    for i := 0; i < 100; i++ {
        ch <- i
    }
    close(ch)
}()

// Consumer
for v := range ch {
    process(v)
}
```

</div>

Концептуально схожи: буферизованная очередь с backpressure. Каналы Go — примитив языка с `select`. `asyncio.Queue` — библиотечный объект, только внутри одного event loop.

---


| | Python threads | Python asyncio | Go goroutines |
|--|--|--|--|
| **Параллелизм CPU** | Нет (GIL) | Нет (1 поток) | Да (GOMAXPROCS) |
| **I/O concurrency** | Да | Да | Да |
| **Стоимость единицы** | ~1MB (OS thread) | ~KB (корутина) | ~2KB (горутина) |
| **Модель** | Preemptive | Cooperative | Preemptive |
| **Блокирующий вызов** | Ок (другие потоки) | Блокирует всё | Ок (планировщик) |
| **Сложность** | Средняя | Высокая (заражение async) | Средняя |

---

# Обработка ошибок

try/except vs error

---

## Python: исключения и иерархия

```
BaseException
├── SystemExit
├── KeyboardInterrupt
├── GeneratorExit
└── Exception
    ├── ValueError, TypeError, AttributeError, ...
    ├── OSError (IOError, FileNotFoundError, ...)
    └── RuntimeError, NotImplementedError, ...
```

```python
def parse_config(path):
    with open(path) as f:        # может бросить FileNotFoundError
        data = json.load(f)      # может бросить JSONDecodeError
    return Config(**data)        # может бросить TypeError, ValueError

# Вызывающий код не знает — функция не объявляет исключения явно
config = parse_config("config.json")
```

---

## Проблема try/except: легко пропустить ошибку

```python
# Антипаттерн: поглощаем все ошибки
try:
    result = process_order(order)
except Exception:
    pass  # тихо игнорируем — order не обработан, никто не знает

# Слишком широко: ловим то, что не ожидали
try:
    user = db.get_user(user_id)
    send_email(user.email, message)
except Exception as e:
    log.error(f"failed: {e}")
    # KeyboardInterrupt? MemoryError? AttributeError в send_email?
    # Все попали сюда
```

<div class="warning">

**Правило:** ловите конкретные исключения, которые вы ожидаете и умеете обработать.

</div>

---

## Обработка ошибок: сравнение подходов

<div class="columns">

```python
# Плюсы Python исключений:
# + Меньше кода в happy path
# + Стектрейс из коробки
# + Легко пробросить наверх

def get_user(id):
    return db.query(id)  # бросает,
                         # если нет соединения

# Минусы:
# - Не видно из сигнатуры
# - Легко пропустить
# - Неожиданные исключения
```

```go
// Плюсы Go ошибок:
// + Явность: видно из сигнатуры
// + Компилятор напоминает обработать
// + errors.Is/As для типизации

func getUser(id int) (User, error) {
    return db.Query(id)
}

// Минусы:
// - Многословно (if err != nil)
// - Стектрейс не из коробки
//   (нужен pkg/errors или fmt.Errorf)
// - Паники всё равно есть
```

</div>

---

# Типизация

Динамика vs статика, и как Python эволюционирует

---

## Python: динамическая типизация + type hints

```python
# Без аннотаций — всё работало до Python 3.5
def greet(name):
    return "Hello, " + name

# С type hints (Python 3.5+) — опциональны, игнорируются рантаймом
def greet(name: str) -> str:
    return "Hello, " + name

# Python 3.10+: удобный синтаксис для Union
def process(value: int | str | None) -> str:
    ...

# Аннотации не мешают написать:
greet(42)  # рантайм не бросит ошибку!
           # нужен статический анализатор
```

---

## mypy, pyright, pydantic

**Статический анализ (проверка до запуска):**
```bash
# mypy — де-факто стандарт статического анализа
mypy service.py  # найдёт: Argument 1 to "greet" has incompatible type "int"

# pyright (Microsoft) — быстрее, строже, используется в Pylance/VS Code
pyright service.py
```

**Runtime валидация (проверка во время выполнения):**
```python
from pydantic import BaseModel

class Order(BaseModel):
    id: int
    amount: float
    status: str

# Pydantic проверит типы при создании объекта
order = Order(id="not-int", amount=10.0, status="new")
# ValidationError: id — value is not a valid integer
```

---

## Go: строгая статическая типизация

```go
// Ошибка типов = ошибка компиляции, не рантайм
func greet(name string) string {
    return "Hello, " + name
}

greet(42)  // compile error: cannot use 42 (type int) as type string

// Нет неявных преобразований
var x int = 42
var y float64 = x      // compile error
var y float64 = float64(x)  // явное приведение — OK

// Интерфейсы как абстракция
type Stringer interface {
    String() string
}
// Любой тип, реализующий String() string — автоматически Stringer
// Без явного implements
```

---

## Generics

<div class="columns">

```go
// Go 1.18+: generics
func Map[T, U any](
    slice []T,
    fn func(T) U,
) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// Использование
doubled := Map([]int{1, 2, 3},
    func(x int) int { return x * 2 })
```

```python
# Python 3.11+: generics
from typing import TypeVar

T = TypeVar('T')
U = TypeVar('U')

def map_list(
    lst: list[T],
    fn: Callable[[T], U]
) -> list[U]:
    return [fn(x) for x in lst]

# Python 3.12 новый синтаксис:
def map_list[T, U](
    lst: list[T],
    fn: Callable[[T], U]
) -> list[U]:
    return [fn(x) for x in lst]
```

</div>

---

# Особенности Python

Вещи, которые стреляют в production

---

## Управление памятью: reference counting + GC

```
Объект в Python:
  ob_refcnt  — счётчик ссылок
  ob_type    — указатель на тип
  ob_val     — значение

a = [1, 2, 3]   # refcnt([1,2,3]) = 1
b = a            # refcnt([1,2,3]) = 2
del a            # refcnt([1,2,3]) = 1
del b            # refcnt([1,2,3]) = 0 → освобождается немедленно

Проблема — циклические ссылки:
  a = {}
  a['self'] = a   # a → a: refcnt никогда не упадёт до 0
```

**Cyclic GC** (gc module) — периодически находит и удаляет циклы.
Запускается по порогам: после ~700 новых объектов.

---

## Ловушка: mutable default arguments

```python
# КАК ВЫГЛЯДИТ
def append_to(element, lst=[]):
    lst.append(element)
    return lst

# ЧТО ПРОИСХОДИТ
append_to(1)  # [1]
append_to(2)  # [1, 2]  ← не [2]!
append_to(3)  # [1, 2, 3]  ← не [3]!

# Дефолтный аргумент — объект, созданный ОДИН РАЗ при определении функции
# Все вызовы используют один и тот же список
```

```python
# ПРАВИЛЬНО: None как sentinel
def append_to(element, lst=None):
    if lst is None:
        lst = []
    lst.append(element)
    return lst
```

---

## Ловушка: late binding в closures

```python
# КАК ВЫГЛЯДИТ
funcs = [lambda: i for i in range(5)]

# ЧТО ОЖИДАЮТ
[f() for f in funcs]  # [0, 1, 2, 3, 4]?

# ЧТО ПОЛУЧАЮТ
[f() for f in funcs]  # [4, 4, 4, 4, 4]

# Closure захватывает переменную i, а не её значение.
# К моменту вызова цикл завершён, i == 4.
```

```python
# ПРАВИЛЬНО: захватить значение через default argument
funcs = [lambda i=i: i for i in range(5)]
[f() for f in funcs]  # [0, 1, 2, 3, 4]
```

---

## Импорты и cold start

```python
# Некоторые импорты занимают секунды:
import torch        # ~2-5 секунд, сотни MB
import transformers # ~1-3 секунды
import numpy        # ~200ms
import pandas       # ~500ms

# Для serverless/Lambda это критично:
# AWS Lambda cold start с torch = 10+ секунд
```

**Стратегии борьбы:**
- Ленивые импорты внутри функции (не на уровне модуля)
- Разделение зависимостей: inference-сервис отдельно от API
- Снапшоты процессов (AWS Lambda SnapStart, Google Cloud Run)
- uv — значительно быстрее pip для установки зависимостей

---

## Экосистема: зоопарк vs единый стандарт

<div class="columns">

```
Python — выбор на каждом шагу:
• Пакетный менеджер: pip / poetry /
  uv / conda / pipenv
• Форматтер: black / autopep8 /
  yapf / ruff
• Линтер: flake8 / pylint /
  ruff / pyflakes
• Тесты: pytest / unittest /
  nose2
• Type checker: mypy / pyright /
  pytype
→ Много хороших инструментов,
  нет одного стандарта
```

```
Go — один стандарт:
• Пакетный менеджер: go mod
• Форматтер: gofmt (один)
• Линтер: go vet + golangci-lint
• Тесты: go test (встроен)
• Type checker: компилятор (всегда)

+ Плюсы: новый разработчик
  сразу знает все инструменты
- Минусы: меньше гибкости

Python-экосистема богаче
для ML/data: numpy, pandas,
sklearn, torch, tf — без аналогов
```

</div>

---

# Практика и выводы

Веб-фреймворки, производительность, когда что выбирать

---

## Веб-фреймворки на Python

| Фреймворк | Модель | Когда брать |
|-----------|--------|-------------|
| **Django** | Sync, batteries included | Традиционный веб, CMS, admin-панель, ORM нужен |
| **Flask** | Sync, микрофреймворк | Простые API, обучение, полный контроль над стеком |
| **FastAPI** | Async-first, type hints | Современные API, нужен OpenAPI, production |
| **Starlette** | Async, минималистичный | Основа для FastAPI, когда нужен контроль |
| **Tornado** | Async, legacy | Веб-сокеты, старые проекты |



---

## FastAPI vs net/http: минимальный сервис

<div class="columns">

```python
# FastAPI
from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()

class Order(BaseModel):
    id: int
    amount: float

@app.get("/orders/{order_id}")
async def get_order(order_id: int) -> Order:
    # Pydantic автоматически валидирует
    # и сериализует ответ
    # OpenAPI docs: /docs
    return Order(id=order_id, amount=99.9)

# uvicorn main:app --reload
```

```go
// net/http
type Order struct {
    ID     int     `json:"id"`
    Amount float64 `json:"amount"`
}

func getOrder(w http.ResponseWriter,
              r *http.Request) {
    idStr := r.PathValue("order_id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "bad id", 400)
        return
    }
    order := Order{ID: id, Amount: 99.9}
    w.Header().Set("Content-Type",
        "application/json")
    json.NewEncoder(w).Encode(order)
}

func main() {
    http.HandleFunc("GET /orders/{order_id}",
        getOrder)
    http.ListenAndServe(":8080", nil)
}
```

</div>

---

## Производительность: когда разница важна

```
Типичные результаты (относительно CPython):

  CPython (baseline):    1x
  Go:                    10-50x  (зависит от задачи)
  PyPy:                  5-10x
  Node.js:               3-10x

НО: большинство backend-задач ограничены I/O, не CPU:

  Время запроса к PostgreSQL: 5-50ms
  Время выполнения Python кода: 0.1-1ms
  → Python код — не узкое место!

  Когда разница важна:
  • >10,000 RPS с тяжёлой бизнес-логикой
  • Обработка данных в реальном времени
  • Финансовые расчёты, low-latency системы
  • CPU-bound без NumPy/Cython

  Когда разница не важна:
  • CRUD API с базой данных
  • Оркестрация микросервисов
  • ML inference (узкое место — GPU/модель)
```

---

| Критерий | Go | Python |
|----------|----|--------|
| **Нагрузка >10k RPS** | Предпочтительно | Возможно (FastAPI + Gunicorn) |
| **ML/AI модели** | Нет экосистемы | Однозначно Python |
| **Data Engineering** | Нет экосистемы | Pandas, Spark, Airflow |
| **Микросервисы** | Отлично | Хорошо |
| **CLI утилиты** | Отлично (бинарь) | Хорошо (нужен Python) |
| **Прототип за день** | Многословно | Быстро |
| **Строгая типизация** | Из коробки | Требует дисциплины |
| **Командная кодовая база** | Легче поддерживать | Сложнее без type hints |
| **Системный уровень** | Да | Нет |
| **Latency < 10ms p99** | Предпочтительно | Сложно |

---

## Отталкиваемся от вводных

```

  ┌─────────────────────────────────────────────────┐
  │  Какая задача?                                  │
  │                                                 │
  │  Нагрузка?  Команда?  Экосистема?  Сроки?       │
  │                                                 │
  │  Что уже есть в стеке?                          │
  └─────────────────────────────────────────────────┘
                       │
            ┌──────────┴───────────┐
            │                      │
       ┌────┴─────┐           ┌────┴─────┐
       │   Go     │           │ Python   │
       │          │           │          │
       │Надёжность│           │Скорость  │
       │Скорость  │           │разработки│
       │Ресурсы   │           │Экосистема│
       └──────────┘           └──────────┘

```

---

## Материалы

1. PEP 703 — Making the GIL Optional: https://peps.python.org/pep-0703/
2. asyncio документация: https://docs.python.org/3/library/asyncio.html
3. FastAPI документация: https://fastapi.tiangolo.com/
4. Real Python — GIL объяснение: https://realpython.com/python-gil/
5. Go vs Python бенчмарки: https://benchmarksgame-team.pages.debian.net/benchmarksgame/
6. mypy документация: https://mypy.readthedocs.io/
7. Pydantic v2: https://docs.pydantic.dev/
8. Официальный Python memory model: https://devguide.python.org/internals/garbage-collector/
