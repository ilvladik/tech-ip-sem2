# Практическое задание 7

## ЭФМО-02-25 Ильин Владислав Викторович

---

# Информация о проекте

**pz7-docker** — два Go-сервиса (**auth**, **tasks**): сборка в **multi-stage** Docker-образы, запуск через **docker compose**, конфигурация через переменные окружения.

## Цели занятия

- Собрать Go-сервисы в Docker-образы (этап сборки и этап запуска)
- Передавать конфигурацию через переменные окружения
- Запускать связку **auth** + **tasks** через `docker compose`
- Понять разницу между image и container, роль `.dockerignore`

## ВАЖНОЕ ПРИМЕЧАНИЕ

| Сервис | Порт | Назначение |
|--------|------|------------|
| auth | **8087** | Проверка Bearer-токена |
| tasks | **8088** | Список задач (с вызовом auth) |

pz1 использует 8081/8082, pz6 — 8086; здесь отдельные порты, чтобы практики не мешали друг другу. Нужен [Docker Desktop](https://www.docker.com/products/docker-desktop/) для сценария с compose.

## Файловая структура проекта

```
pz7-docker/
├── services/
│   ├── auth/          # GET /health, GET /v1/validate
│   └── tasks/         # GET /health, GET /v1/tasks
├── deploy/
│   └── docker-compose.yml
├── start-compose.ps1
├── stop-compose.ps1
├── build-images.ps1
├── start-local.ps1
└── tests.ps1
```

## Запуск через Docker Compose (рекомендуется)

```powershell
cd pz7-docker
.\start-compose.ps1
```

## Тесты

```powershell
cd pz7-docker
.\tests.ps1
```

Остановка compose:

```powershell
cd pz7-docker
.\stop-compose.ps1
```

Проверка вручную:

```powershell
curl.exe -i http://localhost:8088/v1/tasks `
  -H "Authorization: Bearer demo-token" `
  -H "X-Request-ID: pz7-001"
```

## Локальный запуск без Docker

```powershell
cd pz7-docker
.\start-local.ps1
```

## Тесты (локально без Docker)

```powershell
cd pz7-docker
.\tests.ps1
```

## Запуск без PowerShell

Из каталога `pz7-docker`. Нужны **два** терминала.

**Терминал 1 — auth (8087):**

```text
cd services/auth
go run ./cmd/auth
```

**Терминал 2 — tasks (8088):**

CMD:

```text
cd services/tasks
set AUTH_BASE_URL=http://localhost:8087
go run ./cmd/tasks
```

Linux / macOS:

```text
cd services/tasks
export AUTH_BASE_URL=http://localhost:8087
go run ./cmd/tasks
```

**Docker Compose** (оба сервиса):

```text
docker compose -f deploy/docker-compose.yml up --build
```

## Сборка образов вручную

```powershell
.\build-images.ps1
docker run --rm -p 8087:8087 -e AUTH_PORT=8087 techip-auth:0.1
docker run --rm -p 8088:8088 -e TASKS_PORT=8088 -e AUTH_BASE_URL=http://host.docker.internal:8087 techip-tasks:0.1
```

На Windows `host.docker.internal` позволяет контейнеру tasks обращаться к auth на хосте.

## API

### auth (:8087)

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/health` | Проверка живости |
| GET | `/v1/validate` | Заголовок `Authorization: Bearer demo-token` → 200 |

Переменные: `AUTH_PORT`, `AUTH_VALID_TOKEN`, `LISTEN_ADDRESS`.

### tasks (:8088)

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/health` | Проверка живости |
| GET | `/v1/tasks` | Список задач; нужны `Authorization` и опционально `X-Request-ID` |

Переменные: `TASKS_PORT`, `AUTH_BASE_URL`, `LISTEN_ADDRESS`.

В compose: `AUTH_BASE_URL=http://auth:8087` — обращение по имени сервиса в docker-сети, не через `localhost`.

## Примеры запросов и ответов

### auth :8087 — GET /health

```bash
curl http://localhost:8087/health
```

Ответ (`HTTP 200`):

```json
{"status":"ok","service":"auth"}
```

### auth — GET /v1/validate

```bash
curl -H "Authorization: Bearer demo-token" http://localhost:8087/v1/validate
```

Ответ (`HTTP 200`): `{"status":"ok","subject":"demo-user"}`

Без токена или с неверным → `HTTP 401`, тело `{"error":"unauthorized"}`

### tasks :8088 — GET /v1/tasks

Без заголовка:

```bash
curl http://localhost:8088/v1/tasks
```

Ответ (`HTTP 401`): `{"error":"missing authorization"}`

С токеном:

```bash
curl -H "Authorization: Bearer demo-token" -H "X-Request-ID: pz7-001" http://localhost:8088/v1/tasks
```

Ответ (`HTTP 200`):

```json
{
  "request_id": "pz7-001",
  "tasks": [
    {"id": 1, "title": "Read Dockerfile guide", "status": "done"},
    {"id": 2, "title": "Build multi-stage image", "status": "in_progress"},
    {"id": 3, "title": "Run docker compose", "status": "todo"}
  ]
}
```

## Dockerfile (multi-stage)

1. **builder** — `golang:alpine`, `go mod download`, сборка бинарника
2. **Этап запуска** — `alpine`, только бинарник + CA-сертификаты, `CMD ["./app"]`

Секреты и URL **не** зашиваются в образ — только через `ENV` / compose.

## Дополнительно

- **Healthcheck** в `docker-compose.yml` для auth и tasks (`wget /health`)
- **`depends_on`** — tasks стартует после healthy auth
- **`.dockerignore`** — не копировать `.git`, `bin`, логи в контекст сборки
- **`LISTEN_ADDRESS=0.0.0.0`** — иначе порт не пробросится с хоста в контейнер

## Типичные проблемы

| Проблема | Решение |
|----------|---------|
| `connection refused` к auth из tasks | В compose использовать `http://auth:8087`, не `localhost` |
| Порт занят | Остановить другой сервис или сменить маппинг в compose |
| Долгая пересборка | Проверить `.dockerignore`, собирать из каталога сервиса |
| 401 на tasks | Токен `Bearer demo-token` (или `AUTH_VALID_TOKEN`) |

## Отчёт (кратко)

- **Image** — шаблон (слои файловой системы + метаданные)
- **Container** — запущенный экземпляр image
- **Multi-stage** — маленький финальный образ без компилятора Go
- **Сеть Docker** — сервисы видят друг друга по имени (`auth`, `tasks`)
