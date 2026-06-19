# Практическое задание 5

## ЭФМО-02-25 Ильин Владислав Викторович

---

# Информация о проекте

**pz5-security** — HTTPS-сервис на Go с PostgreSQL. Защита транспорта (TLS) и данных (параметризованные SQL-запросы, подготовленные выражения).

## Цели занятия

- Настроить HTTPS (TLS) и редирект с HTTP на HTTPS.
- Подключить PostgreSQL и безопасные запросы (prepared statements).
- Разделить безопасный и небезопасный (демо) доступ к данным.
- Понять риски SQL-инъекций и роль валидации входа.

## Файловая структура проекта

```
pz5-security/
├── certs/
├── cmd/
│   ├── server/main.go
│   └── gencert/main.go
├── internal/
│   ├── config/
│   ├── httpapi/
│   ├── httpserver/
│   └── student/
├── sql/init.sql
├── docker-compose.yml
├── generate-certs.ps1
├── start-db.ps1
├── start-server.ps1
├── tests.ps1
├── go.mod
└── go.sum
```

## ВАЖНОЕ ПРИМЕЧАНИЕ

| Сервис | Порт |
|--------|------|
| HTTP → HTTPS redirect | **8085** |
| HTTPS API | **8443** |
| PostgreSQL | **5432** |

Перед API нужны **PostgreSQL** (`.\start-db.ps1` или `docker compose`) и **сертификаты** в `certs/` (`.\generate-certs.ps1` или `go run ./cmd/gencert`). Запросы к API: `curl -k` из‑за самоподписанного сертификата. На учебном сервере порты могут отличаться.

## Тестовые данные

После `sql/init.sql` в БД есть студенты **id 1–3**. В `tests.ps1` **`id=999`** — проверка **404** (нет строки в БД).

## Запуск

```powershell
cd pz5-security

# 1. База данных
.\start-db.ps1

# 2. TLS-сертификаты (самоподписанные, через Go)
.\generate-certs.ps1

# 3. Сервер
.\start-server.ps1
```

## Тесты

```powershell
cd pz5-security
.\tests.ps1
```

## Запуск без PowerShell

Из каталога `pz5-security`.

**1. PostgreSQL:**

```text
docker compose up -d
```

**2. Сертификаты** (если папки `certs/` ещё нет — см. `generate-certs.ps1` или методичку):

```text
go run ./cmd/gencert
```

**3. Сервер** — HTTPS **8443**, редирект с HTTP **8085**:

```text
go run ./cmd/server
```

Проверка (curl с игнором самоподписанного сертификата):

```text
curl -k https://localhost:8443/health
curl -s -o NUL -w "%{http_code}" http://localhost:8085/health
```

## API (HTTPS)

| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/health` | проверка сервиса |
| GET | `/students?id=1` | студент по ID (безопасно) |
| GET | `/students/by-email?email=...` | студент по email |
| GET | `/students/unsafe?id=...` | **демо** небезопасного SQL |

Все запросы: `curl -k` (самоподписанный сертификат).

## Примеры запросов и ответов

HTTPS **8443**. В БД после `init.sql` студенты **id 1–3**.

### GET /health

```bash
curl -k https://localhost:8443/health
```

Ответ (`HTTP 200`):

```json
{"status":"ok","scheme":"https"}
```

### HTTP → HTTPS (порт 8085)

```bash
curl -s -o NUL -w "%{http_code}" http://localhost:8085/health
```

Ответ: **`301`** → `https://localhost:8443/health`

### GET /students?id=1

Ответ (`HTTP 200`):

```json
{
  "id": 1,
  "full_name": "Иванов Иван Иванович",
  "study_group": "ИТТ-01-25",
  "email": "ivanov@example.com"
}
```

### GET /students?id=999

Ответ (`HTTP 404`):

```json
{"error":"студент с таким id не найден в базе","id":999}
```

### GET /students/by-email?email=ivanov@example.com

Ответ (`HTTP 200`): тот же JSON студента.

### GET /students/unsafe?id=1 (демо)

Ответ (`HTTP 200`):

```json
{
  "warning": "unsafe SQL concatenation — do not use in production",
  "query": "SELECT ... WHERE id = 1",
  "student": { "id": 1, "full_name": "...", "study_group": "...", "email": "..." }
}
```

## Переменные окружения

| Переменная | По умолчанию |
|------------|--------------|
| `HTTPS_ADDR` | `:8443` |
| `HTTP_REDIRECT_ADDR` | `:8085` |
| `TLS_CERT_FILE` | `certs/server.crt` |
| `TLS_KEY_FILE` | `certs/server.key` |
| `DB_DSN` | `postgres://postgres:postgres@localhost:5432/study_security?sslmode=disable` |

## Дополнительные задания (выполнено)

### HTTP → HTTPS редирект (:8085)

**Как работает:** отдельный HTTP-сервер на `:8085` отвечает `301 Moved Permanently` на `https://localhost:8443` + тот же путь. Любой запрос `http://localhost:8085/...` перенаправляется на HTTPS.

### Конфигурация из переменных окружения

**Как работает:** `internal/config` читает адрес, DSN и пути к сертификатам из переменных окружения. Если не заданы — значения по умолчанию из таблицы выше.

### GET /students/by-email

**Как работает:** email передаётся как query-параметр, проверяется по белому списку формата, затем выборка через **подготовленное выражение** с `$1` — инъекция в SQL невозможна.

### Валидация по белому списку

**Как работает:** `id` — только положительное целое ≤ 1e9; `email` — проверка regexp и длины ≤ 254. Некорректный ввод → `400 Bad Request` до обращения к БД.

### Демо небезопасного SQL

**Как работает:** `/students/unsafe` склеивает `id` в строку SQL (`UnsafeGetByID`). В ответе видно сформированный запрос — для сравнения с безопасным `/students?id=1`. **В production так делать нельзя.**

## Безопасный vs небезопасный SQL

```go
// ОПАСНО — конкатенация
query := "SELECT ... WHERE id = " + rawID

// БЕЗОПАСНО — плейсхолдер + подготовленное выражение
stmt.QueryRow(id)  // SQL: ... WHERE id = $1
```

## HTTPS vs HTTP

| HTTP | HTTPS |
|------|-------|
| данные открыты | шифрование TLS |
| порт 8085 (редирект) | порт 8443 (основной API) |
| нет проверки сервера | сертификат server.crt/key |

Самоподписанный сертификат подходит для учёбы; браузер и curl покажут предупреждение — для curl используйте `-k`.
