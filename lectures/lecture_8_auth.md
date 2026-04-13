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

<!-- _class: lead -->

# Все про Auth

## 🐺☝🏻🐺☝🏻🐺☝🏻🐺☝🏻

---

<!-- _class: lead -->

# ID & AuthN & AuthZ

---

## Идентификация vs Аутентификация vs Авторизация

В чём разница?

---

# Идентификация vs Аутентификация vs Авторизация

<div class="columns-3">
<div>

## Идентификация
* Определяем ID пользователя

* (Опционально) Проверяем его существование

**Результат:** user with ID=132, (not)exists
</div>
<div>

## AuthN
Проверяем действительно ли пользователь тот за кого себя выдает

**Результат:** valid/invalid


</div>
<div>

## AuthZ
Проверяем права доступа пользователя для конкретной операции

**Результат:** разрешено/запрещено

</div>
</div>

---

# Аутентификация vs Авторизация

<div class="columns">
<div>

## Аналогия
```
Номер паспорта = Идентификация
Паспорт        = AuthN
Виза           = AuthZ
```

</div>
<div>

## Важно
- AuthN невозможна без идентификации
- AuthZ невозможна без AuthN
- AuthN может быть без AuthZ, AuthN может быть без идентификации
- Проверки обычно на разных слоях системы

</div>
</div>

---

<!-- _class: lead -->

# Механизмы аутентификации

---

# Basic Auth

```
Authorization: Basic dXNlcjpwYXNzd29yZA==
                     └── base64(login:password)
```

<div class="columns">
<div>

## Как работает
1. Клиент кодирует `login:password` в base64
2. Отправляет заголовок при **каждом** запросе
3. Сервер декодирует и проверяет

</div>
<div>

## Когда использовать
✅ Никогда

? Internal API между сервисами
? CLI-инструменты, Простые скрипты

❌ Публичные API, Браузерные клиенты

</div>
</div>

**Только HTTPS — иначе credentials в открытом виде**

---
# Basic Auth

![](https://habrastorage.org/r/w1560/files/c27/ac0/637/c27ac06373984352a1ebe2f6424cd9e9.png)

---

# Session-based Auth

```
POST /login
  → Set-Cookie: session_id=abc123

GET /profile
  Cookie: session_id=abc123
    → сервер ищет в Redis → OK
```

<div class="columns">
<div>

## Плюсы
- Легко инвалидировать сессию
- Сервер контролирует состояние
- Привычно для веба

</div>
<div>

## Проблемы
- **Stateful** — нужен shared storage
- При масштабировании → Redis/Memcached
- Уязвим к CSRF
- Плохо подходит для API

</div>
</div>

---

![](https://habrastorage.org/r/w1560/files/8e5/211/8bb/8e52118bbaa84a4286e2ef2a2a5ad36d.png)

---

# Безопасность токенов

* В Cookie или в localStorage?
* Если в Cookie, как передаем?

---

# Куки




```http
Set-Cookie: token=<token>; HttpOnly; Sequre; SameSite=Strict;
```

---

# JWT — JSON Web Token

```
eyJhbGciOiJSUzI1NiJ9  .  eyJzdWIiOiIxMjMiLCJleHAiOjE3fQ  .  signature
└── header (base64url)    └── payload (base64url)              └── подпись
```

Payload содержит claims:

```js
{
  "sub": "user_123",
  "exp": 1700000000,
  "iat": 1699999000,
  "iss": "auth.myapp.com",
  "roles": ["admin", "editor"]
}
```

Сервер проверяет подпись — без обращения к БД

---

# Структура JWT

header — содержит информацию об алгоритме шифрования и типе токена (JWT)

payload — данные токена. Стандартные поля:
* iss (Issuer) — издатель токена, ваше приложение
* sub (Subject) — собственник токена, id пользователя
* aud (Audience) — массив url серверов, для которых предназначен токен
* exp (Expiration Time) — время, в течение которого токен считается валидным.
* nbf (Not Before) — временная метка, до которй токен не считается валидным
* iat (Issued At) — время создания токена
* jti (JWT ID) — уникальный идентификатор токена

signature — подпись header + payload.

---

# JWT

<div class="columns">
<div>

## Нельзя инвалидировать
JWT валиден до истечения `exp`

Логаут = удалить из localStorage
Но токен продолжает работать

**Решения:**
- Короткий TTL + refresh tokens
- Blocklist в Redis

</div>
<div>

## Algorithm confusion
Аккуратно работает с алгоритмами подписи, запрещаем none

```go
// Всегда явно указывать алгоритм:
token, err := jwt.Parse(
  tokenString,
  keyFunc,
  jwt.WithValidMethods(
    []string{"RS256"},
  ),
)
```

</div>
</div>

---

![h:600](https://media2.dev.to/dynamic/image/width=800%2Cheight=%2Cfit=scale-down%2Cgravity=auto%2Cformat=auto/https%3A%2F%2Fdev-to-uploads.s3.amazonaws.com%2Fuploads%2Farticles%2F2b6yy74rcfqrrf2oz8zp.png)

---

# API Keys

```
X-API-Key: sk_live_abc123xyz...
Authorization: Bearer sk_live_abc123xyz...
```

<div class="columns">
<div>

## Когда использовать
- Machine-to-machine
- Разработчики подключают внешний API
- Нет интерактивного пользователя

</div>
<div>

## Как хранить
- В БД — **только хеш** (как пароль)
- Показывать пользователю один раз
- Никогда не логировать

</div>
</div>

---

# Как еще?

---

# Как еще?

* По сертификатам (mTLS)
* Биометрия
* OTP (SMS/email/push)
* TOTP
* Hardware Token (YubiKey)
* Passkeys


---

# OAuth 2.0

OAuth 2.0 — это протокол **авторизации**, не аутентификации

Цель — делегированный доступ: пользователь разрешает приложению A действовать от его имени в сервисе B

## Роли
- **Resource Owner** — пользователь
- **Client** — ваше приложение
- **Authorization Server** — выдаёт токены
- **Resource Server** — API с данными


---

## Flows
- **Authorization Code** — веб-приложения
- **Authorization Code + PKCE** — SPA, мобилки
- **Client Credentials** — сервис-к-сервису
- **Implicit** — устаревший 

---

# OAuth 2.0 — Authorization Code Flow

*Смотрим схему*

**PKCE** — для SPA и мобилок, где нельзя хранить `client_secret`:
вместо секрета — cryptographic challenge, генерируется на клиенте

---

# OIDC — OpenID Connect

OAuth 2.0 решает авторизацию, но не говорит ничего об **идентичности** пользователя

OIDC = OAuth 2.0 + слой аутентификации

<div class="columns">
<div>

## Что добавляет
- `id_token` — JWT с данными о пользователе
- `/userinfo` endpoint
- Стандартные claims: `name`, `email`, `picture`

</div>
<div>

## Когда встречаете
- "Login with Google"
- "Login with GitHub"
- "Login with Apple"

OAuth 2.0 там — транспорт
OIDC — то что даёт identity

</div>
</div>

---

![alt text](https://habrastorage.org/r/w1560/getpro/habr/upload_files/0ee/e27/4c7/0eee274c7f39827ff99ff431ad7651f2.png)

---

# SSO
Single-sign-on

1. Открываешь app1.com → нет сессии
   → редирект на idp.com/authorize
   → логинишься → IDP создаёт свою сессию (куку)
   → редирект обратно с code → app1 получает токены

2. Открываешь app2.com → нет сессии
   → редирект на idp.com/authorize
   → IDP видит свою куку → сессия есть, consent уже давал
   → сразу редирект с code → без повторного логина

Альтернативно: SAML, Kerberos

---

<!-- _class: lead -->

# Модели авторизации

---

# ACL — Access Control List

Самая простая модель. Список: кому что разрешено на конкретном объекте.
```
Файл report.pdf:
  alice  → read, write
  bob    → read
  charlie → (нет доступа)
```

<div class="columns">
<div>

## Плюсы
- Просто понять и реализовать
- Гибко — точный контроль на уровне объекта
- Привычно (Unix permissions, S3 bucket policies)

</div>
<div>

## Проблемы
- Плохо масштабируется — 10000 пользователей × 10000 объектов
- Сложно управлять — добавить нового пользователя = обновить тысячи списков

</div>
</div>

---

# RBAC — Role-Based Access Control

```
Роли:        admin        editor       viewer
               │            │            │
Права:       r+w+d         r+w           r
Пользователи:
  alice  → admin
  bob    → editor
  charlie → viewer
```

<div class="columns">
<div>

## Плюсы
- Легко управлять — меняешь роль, не права
- Стандарт де-факто для большинства приложений

</div>
<div>

## Проблемы
- Role explosion — со временем ролей становится слишком много
- Сложно выразить "пользователь видит только свои данные"

</div>
</div>

---

# RBAC — Иерархия ролей
```
        super-admin
             │
           admin
          /     \
       editor  moderator
          \     /
          viewer
```

Дочерняя роль наследует права родительской. `editor` автоматически имеет всё что есть у `viewer`.

---

# ABAC — Attribute-Based Access Control

Решение принимается на основе **атрибутов** субъекта, объекта и окружения
```
Разрешить если:
  user.department == document.department
  AND user.clearance >= document.classification
  AND request.time BETWEEN 09:00 AND 18:00
  AND request.ip IN corporate_network
```

---
# ABAC — Attribute-Based Access Control

<div class="columns">
<div>

## Плюсы
- Очень гибко — любые условия
- Решает проблему "пользователь видит только своё"
- Контекст запроса (время, локация, устройство)
- Не нужно создавать роль под каждый случай

</div>
<div>

## Проблемы
- Сложнее реализовать и отлаживать
- Политики трудно читать и аудировать
- Производительность — каждый запрос вычисляет политику

</div>
</div>

---

# ReBAC — Relationship-Based Access Control

Права определяются **отношениями** между объектами. Придумано в Google (Zanzibar, 2019) для Google Drive.
```
Пользователь имеет доступ к документу если:
  - он owner документа
  - он editor документа
  - он viewer документа
  - он member группы которая имеет доступ
  - у него есть доступ к папке которая содержит документ
```
```
alice  owner   document:report
bob    viewer  document:report
team   editor  document:report
alice  member  team
```

Google Drive, GitHub, Notion

---

# PBAC — Policy-Based Access Control

Права описываются декларативными политиками. AWS IAM — самый известный пример.
```json
{
  "Effect": "Allow",
  "Action": ["s3:GetObject", "s3:PutObject"],
  "Resource": "arn:aws:s3:::my-bucket/*",
  "Condition": {
    "StringEquals": { "aws:RequestedRegion": "eu-west-1" },
    "Bool": { "aws:MultiFactorAuthPresent": "true" }
  }
}
```

По сути ABAC но с явным языком политик. Политики можно версионировать, аудировать, переиспользовать.

---

<!-- _class: lead -->

# Готовые решения

Не пишите своё с нуля

---

# Self-hosted Identity Providers

<div class="columns-3">
<div>

## Keycloak
Enterprise-grade IdP
OAuth2, OIDC, SAML
User federation (LDAP/AD)
Social login, MFA

**Тяжёлый, но мощный**
Популярен в enterprise

</div>
<div>

## Zitadel
Современная альтернатива
Лучший UX и документация
OAuth2, OIDC, SAML
Есть cloud-версия

**Активно развивается**
Хороший выбор для новых проектов

</div>
<div>

## Ory Stack
Composable архитектура:
- **Kratos** — users
- **Hydra** — OAuth2/OIDC
- **Keto** — authorization
- **Oathkeeper** — proxy

**Для тех кто хочет контроль**

</div>
</div>
