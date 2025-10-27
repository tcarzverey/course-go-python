# wscat

`wscat` — это утилита командной строки, реализующая простой WebSocket клиент и сервер, аналог [npm-пакета
`wscat`](https://www.npmjs.com/package/wscat).

WebSocket — это двусторонний канал поверх TCP. `wscat` позволяет отправлять и получать сообщения через этот канал, читая
данные из `stdin` и выводя все входящие сообщения в `stdout`.

---

## Логика

Программа работает в одном из двух режимов (взаимоисключающих):

```bash
wscat (--listen <port> | --connect <url>)
```

### Режим 1: Клиент

```bash
wscat --connect ws://example.com/socket
```

* Подключается к WebSocket серверу по указанному адресу.
    * Если не удается подключиться - завершает работу с ошибкой
* После подключения:
    * Читает строки из `stdin` и отправляет их на сервер.
    * Все входящие сообщения от сервера печатает в `stdout`.
* При закрытии соединения сервером (или обрыве сети) клиент **завершается автоматически**.
* При получении сигналов `SIGINT` / `SIGTERM` (например при нажатии `Ctrl+C`) соединение **gracefully** закрывается,
  и программа выходит.

---

### Режим 2: Сервер

```bash
wscat --listen 8080
```

* Запускает WebSocket сервер, слушающий указанный порт (`ws://localhost:<port>`).
* Сервер принимает **только одно соединение**.
* До подключения клиента введённые в stdin строки не сохраняются и не отправляются.
* После подключения:
    * Сообщения из `stdin` пересылаются клиенту.
    * Все сообщения от клиента печатаются в `stdout`.
* При разрыве соединения или сигналах `SIGINT` / `SIGTERM` сервер **gracefully** закрывает WebSocket и останавливается.

---

## Общие моменты

* Все переходы между состояниями, ошибки и некорректные действия должны логгироваться, например:
    * старт/завершение сервера/клиента
    * ввод текста на сервере до подключения клиента
    * ошибка при соединении
    * и так далее
* Завершение **graceful**:
    * Закрываются все соединения.
    * Можно использовать `context.WithCancel` и перехват сигналов из `os/signal`.
* После выполнения надо ответить на вопросы в [questions.md](questions.md)

---

## Примеры

Информационные логи могут быть в любом читаемом формате, в примерах ниже приведен один из вариантов, но не обязательно делать
именно так

### Клиент:

```bash
$ wscat --connect ws://localhost:8080
2025/01/01 11:11:14 connecting to  ws://localhost:7777
2025/01/01 11:11:14 connected
hi from server # получаем от сервера
hi from client # пишем сами
2025/01/01 11:11:15 websocket connection closed
2025/01/01 11:11:16 client exiting
```

### Сервер:

```bash
$ wscat --listen 8080
2025/01/01 11:11:11 starting websocket server at port=8080
2025/01/01 11:11:12 waiting for client connection...
hi from server # пишем сами
2025/01/01 11:11:13 error: client is not connected yet 
2025/01/01 11:11:14 client connected
hi from server # пишем сами
hi from client  # получаем от клиента
2025/01/01 11:11:15 received signal interrupt
2025/01/01 11:11:16 server exiting
```

---

## Полезные ссылки

* [WebSocket (Wikipedia)](https://en.wikipedia.org/wiki/WebSocket)
* [Gorilla WebSocket (Go package)](https://pkg.go.dev/github.com/gorilla/websocket)
* [Graceful shutdown and signal handling in Go (Habr)](https://habr.com/ru/articles/908344/)
