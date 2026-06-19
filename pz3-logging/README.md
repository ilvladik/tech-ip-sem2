# Практическое задание 3

## ЭФМО-02-25 Ильин Владислав Викторович

---

# Информация о проекте

**pz3-logging** — HTTP-сервис на Go со структурированным логированием через **zap**.
Эндпоинты: проверка здоровья и работа со студентами. Все запросы логируются с полями method, path, status_code, duration, request_id.

## Цели занятия

- Настроить structured logging (zap).
- Логировать HTTP-запросы через промежуточный слой.
- Использовать уровни debug, info, warn, error.
- Добавлять контекст в логи: путь, метод, статус, время, request_id.

## Файловая структура проекта

```
pz3-logging/
├── cmd/server/main.go
├── internal/
│   ├── httpapi/
│   │   ├── handler.go
│   │   ├── middleware.go
│   │   └── response_writer.go
│   └── student/
│       ├── model.go
│       └── repo.go
├── pkg/logger/logger.go
├── logs/app.log          # создаётся при запуске
├── start-server.ps1
├── tests.ps1
├── go.mod
└── go.sum
```

## ВАЖНОЕ ПРИМЕЧАНИЕ

HTTP API слушает порт **8083** (чтобы не пересекаться с другими практиками). На учебном сервере порт может отличаться — уточните у преподавателя.

## Тестовые данные

В памяти студенты **id 1–3**. В `tests.ps1` запросы **`/students/999`** и дубликат email — **намеренные** проверки ошибок (404, 409).

## Запуск

```powershell
cd pz3-logging
.\start-server.ps1
```

## Тесты

```powershell
cd pz3-logging
.\tests.ps1
```

## Запуск без PowerShell

Из каталога `pz3-logging` (порт **8083**):

```text
go run ./cmd/server
```

Переменные (по желанию), CMD:

```text
set LOG_FILE=logs/app.log
set LOG_LEVEL=info
go run ./cmd/server
```

Linux / macOS:

```text
export LOG_FILE=logs/app.log
export LOG_LEVEL=info
go run ./cmd/server
```

Переменные окружения:

| Переменная | По умолчанию | Назначение |
|------------|--------------|------------|
| `LOG_FILE` | `logs/app.log` | путь к файлу логов |
| `LOG_LEVEL` | `info` | debug / info / warn / error |

## API

| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/health` | статус сервиса |
| GET | `/students/{id}` | студент по ID |
| POST | `/students` | создать студента (JSON body) |

## Примеры запросов и ответов

Порт **8083**. В памяти студенты **id 1–3**; **999** — проверка 404.

### GET /health

```bash
curl http://localhost:8083/health
```

Ответ (`HTTP 200`):

```json
{"status":"ok"}
```

### GET /students/1

Ответ (`HTTP 200`):

```json
{
  "id": 1,
  "full_name": "Иванов Иван Иванович",
  "group": "ИТТ-01-25",
  "email": "ivanov@example.com"
}
```

Заголовок ответа: `X-Request-ID: <uuid>`.

### GET /students/999

Ответ (`HTTP 404`):

```json
{"error":"студент с таким id не найден в системе","id":999}
```

### POST /students

```bash
curl -X POST http://localhost:8083/students \
  -H "Content-Type: application/json" \
  -d '{"full_name":"Козлов Дмитрий","group":"ИТТ-04-25","email":"kozlov@example.com"}'
```

Ответ (`HTTP 201`): JSON созданного студента с новым `id`.

Дубликат email → `HTTP 409`, текст `email already exists`.

## Дополнительные задания (выполнено)

### Логи в файл

**Как работает:** `pkg/logger` через `zapcore.NewTee` пишет JSON одновременно в **stdout** и в файл `logs/app.log` (путь задаётся `LOG_FILE`). Каталог создаётся автоматически.

### request_id в ответе

**Как работает:** промежуточный слой генерирует ID, кладёт в контекст запроса и в заголовок ответа **`X-Request-ID`**. Тот же ID есть во всех связанных логах — удобно сопоставить запрос и строку в логе.

### Debug в бизнес-логике

**Как работает:** в `repo.GetByID` и `repo.Create` в начале операции пишется уровень отладки с `student_id` или полями входа — видно, что репозиторий реально вызван, до проверки ошибок.

### POST /students

**Как работает:** принимает JSON `{full_name, group, email}`, логирует тело на `Debug`, при пустых полях — `Warn` + 400, при дубликате email — `Warn` + 409, при успехе — `Info` + 201 и JSON созданного студента.

## Уровни логов в проекте

| Уровень | Примеры |
|---------|---------|
| **Отладка** | health, начало операций в репозитории, тело POST |
| **Info** | старт сервера, входящий/завершённый запрос, успешная выдача/создание |
| **Warn** | неверный метод, bad id, валидация, дубликат email |
| **Error** | студент не найден, ошибки чтения тела |

## Сравнение с обычным log

Текстовый `log.Println` даёт одну строку без полей. Zap пишет JSON с ключами — в Kibana/Grafana/Loki можно фильтровать по `student_id`, `status_code`, `request_id` без парсинга текста.
