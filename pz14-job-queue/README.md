# Практическое задание 14

## ЭФМО-02-25 Ильин Владислав Викторович

---

# Информация о проекте

**pz14-job-queue** — постановка фоновых **jobs** в RabbitMQ: **202 Accepted** сразу, повторы при ошибке, **DLQ**, идемпотентность по `message_id`.

## Цели занятия

- Разделить HTTP-приём заявки и фактическую обработку в воркере.
- Настроить повторные попытки и мёртвую очередь (**DLQ**).
- Обеспечить идемпотентность обработки по идентификатору сообщения.

## Связь с pz13

| pz13 | pz14 |
|------|------|
| `task_events` — уведомление | `task_jobs` — **задача на обработку** |
| событие без ожидания результата | повторы + DLQ |
| — | `message_id` — идемпотентность |

## ВАЖНОЕ ПРИМЕЧАНИЕ

| Компонент | Порт |
|-----------|------|
| tasks API | **8097** |
| RabbitMQ | 5672 / UI **15672** |

Нужны RabbitMQ (Docker), процесс **worker** и **tasks** — как в pz13.

## Очереди

- `task_jobs` — основная (durable, DLX → dlq)
- `task_jobs_dlq` — сообщения после 3 неудачных попыток

## Файловая структура проекта

```
pz14-job-queue/
├── deploy/rabbit/docker-compose.yml
├── internal/
│   ├── publisher/
│   ├── rabbit/
│   └── store/
├── pkg/jobs/
├── services/tasks/
├── services/worker/
├── start-rabbit.ps1
├── start-worker.ps1
├── start-tasks.ps1
├── go.mod
└── tests.ps1
```

## Запуск

```powershell
cd pz14-job-queue
.\start-rabbit.ps1
```

В отдельных окнах:

```powershell
cd pz14-job-queue
.\start-worker.ps1
```

```powershell
cd pz14-job-queue
.\start-tasks.ps1
```

## Тесты

```powershell
cd pz14-job-queue
.\tests.ps1
```

## API

```http
POST /v1/jobs/process-task
Authorization: Bearer demo-token
{"task_id": "t_001"}
```

## Примеры запросов и ответов

Порт **8097**.

### POST /v1/jobs/process-task

```bash
curl -i -X POST http://localhost:8097/v1/jobs/process-task \
  -H "Authorization: Bearer demo-token" \
  -H "Content-Type: application/json" \
  -d "{\"task_id\":\"t_001\"}"
```

Ответ (`HTTP 202`):

```json
{"status":"accepted","task_id":"t_001"}
```

В логах worker — успешная обработка, **ack**.

### task_id=t_fail

Тот же запрос с `"task_id":"t_fail"` → после **3** попыток сообщение уходит в очередь **DLQ** (`task_jobs_dlq`).

Без токена → `HTTP 401`, `{"error":"unauthorized"}`.

## Запуск без PowerShell

Из каталога `pz14-job-queue`. Три терминала — как в pz13.

```text
docker compose -f deploy/rabbit/docker-compose.yml up -d
```

```text
cd services/worker
go run ./cmd/worker
```

```text
cd services/tasks
go run ./cmd/tasks
```

API: http://localhost:8097

Сценарии в `tests.ps1`:

- **`t_001`** — успех, **ack**
- **`t_fail`** — три попытки, затем **DLQ**

## JSON job

```json
{
  "job": "process_task",
  "task_id": "t_001",
  "attempt": 1,
  "message_id": "uuid"
}
```

## Дополнительно

- **Предварительная выборка** — `PREFETCH=1` (по умолчанию)
- **Идемпотентность** — повтор с тем же `message_id` пропускается
- **DLQ** — смотреть в http://localhost:15672
