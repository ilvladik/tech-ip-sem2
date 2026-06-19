# Практическое задание 13

## ЭФМО-02-25 Ильин Владислав Викторович

---

# Информация о проекте

**pz13-rabbitmq** — сервис **tasks** публикует событие `task.created` в **RabbitMQ**; отдельный **worker** потребляет сообщения и подтверждает **ack**.

## Цели занятия

- Издатель (tasks) публикует `task.created` в очередь `task_events`
- Потребитель (обработчик) читает, логирует, отправляет **подтверждение (ack)**
- **Постоянная** очередь и **постоянные** сообщения
- **Предварительная выборка** (prefetch, по умолчанию 1)

## Связь с другими практиками

| Практика | Связь |
|----------|--------|
| **pz7** | Сервис **tasks**, Bearer `demo-token`, `X-Request-ID` |
| **pz10** | Тот же домен задач |
| **pz13** | После `POST /v1/tasks` → событие в очередь → **обработчик** |

Поток: **HTTP синхронно** создаёт задачу → **асинхронно** обработчик обрабатывает событие.

## ВАЖНОЕ ПРИМЕЧАНИЕ

| Компонент | Порт |
|-----------|------|
| tasks HTTP | **8096** |
| RabbitMQ AMQP | **5672** |
| RabbitMQ UI | **15672** (guest / guest) |

Нужны **три процесса**: RabbitMQ (Docker), **worker**, **tasks** — порядок см. в «Запуск».

## Файловая структура проекта

```
pz13-rabbitmq/
├── deploy/rabbit/docker-compose.yml
├── pkg/events/              # JSON-события
├── internal/publisher/
├── services/tasks/          # издатель
├── services/worker/         # потребитель
├── start-rabbit.ps1
├── start-worker.ps1
├── start-tasks.ps1
└── tests.ps1
```

## Запуск (3 шага)

**1. RabbitMQ** (нужен Docker):

```powershell
.\start-rabbit.ps1
```

**2. Обработчик** (отдельное окно):

```powershell
.\start-worker.ps1
```

**3. Tasks** (ещё одно окно):

```powershell
.\start-tasks.ps1
```

## Тесты

Из каталога `pz13-rabbitmq` (после шагов 1–3):

```powershell
.\tests.ps1
```

## Запуск без PowerShell

Из каталога `pz13-rabbitmq`. Нужны **три** терминала.

**1. RabbitMQ:**

```text
docker compose -f deploy/rabbit/docker-compose.yml up -d
```

**2. Обработчик:**

```text
cd services/worker
go run ./cmd/worker
```

**3. Tasks (порт 8096):**

```text
cd services/tasks
go run ./cmd/tasks
```

В логах обработчика:

```
received event=task.created task_id=t_001 ts=... request_id=pz13-001
```

## Формат события (JSON)

```json
{
  "event": "task.created",
  "task_id": "t_001",
  "ts": "2026-05-17T12:00:00Z",
  "request_id": "pz13-001",
  "producer": "tasks-service",
  "version": "1"
}
```

## API tasks

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/health` | Проверка |
| POST | `/v1/tasks` | Создать задачу + опубликовать событие |

Заголовки: `Authorization: Bearer demo-token`, `X-Request-ID`.

## Примеры запросов и ответов

Порт **8096**.

### GET /health

```bash
curl http://localhost:8096/health
```

Ответ (`HTTP 200`):

```json
{"status":"ok","service":"tasks"}
```

### POST /v1/tasks

```bash
curl -X POST http://localhost:8096/v1/tasks \
  -H "Authorization: Bearer demo-token" \
  -H "X-Request-ID: pz13-001" \
  -H "Content-Type: application/json" \
  -d "{\"title\":\"Новая задача\",\"description\":\"Тест RabbitMQ\"}"
```

Ответ (`HTTP 201`):

```json
{
  "id": "t_001",
  "title": "Новая задача",
  "description": "Тест RabbitMQ"
}
```

В логах **worker** появится строка с событием `task.created` (см. формат JSON ниже).

Без токена → `HTTP 401`, `{"error":"unauthorized"}`.

## Дополнительно

### Режим публикации (`PUBLISH_MODE`)

| Значение | Поведение |
|----------|-----------|
| `best_effort` (по умолчанию) | Задача создана, ошибка публикации только в лог |
| `strict` | HTTP 500, если не удалось опубликовать в очередь |

```powershell
$env:PUBLISH_MODE = "strict"
.\start-tasks.ps1
```

### Предварительная выборка (prefetch)

```powershell
$env:PREFETCH = "5"
.\start-worker.ps1
```

## RabbitMQ Management

http://localhost:15672 — очередь `task_events`, сообщения, потребители.

## Переменные окружения

| Переменная | По умолчанию |
|------------|--------------|
| `SERVER_ADDR` | `:8096` |
| `RABBIT_URL` | `amqp://guest:guest@localhost:5672/` |
| `QUEUE_NAME` | `task_events` |
| `PUBLISH_MODE` | `best_effort` |
| `PREFETCH` | `1` |

## Отчёт (кратко)

- **Издатель** — tasks после успешного POST
- **Потребитель** — отдельный обработчик
- **Подтверждение (ack)** — сообщение не теряется при сбое до подтверждения
- **Постоянная очередь и сообщения** — переживают перезапуск RabbitMQ
