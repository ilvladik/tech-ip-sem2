# Практическое задание 9

## ЭФМО-02-25 Ильин Владислав Викторович

---

# Информация о проекте

**pz9-redis-cache** — HTTP API задач с **Redis** как внешним кэшем: стратегия «кэш сбоку», TTL, инвалидация при изменении данных.

## Цели занятия

- Использовать Redis как внешний кэш (не источник истины)
- Реализовать стратегию **кэш сбоку** (cache-aside)
- Настроить время жизни ключей (TTL) и разброс времени (jitter)
- Инвалидировать кэш при PATCH/DELETE
- Плавная деградация при недоступности Redis

## ВАЖНОЕ ПРИМЕЧАНИЕ

| Сервис | Порт |
|--------|------|
| HTTP API | **8089** |
| Redis | **6379** |

Для проверки нужен **Redis** (Docker из `deploy/redis` или локально). Поведение без Redis — в разделе **«Проверка без Redis»** ниже.

## Файловая структура проекта

```
pz9-redis-cache/
├── cmd/server/main.go
├── internal/
│   ├── cache/       # redis, keys, serializer, ttl
│   ├── config/
│   ├── httpapi/
│   ├── service/     # логика кэш сбоку
│   └── task/
├── deploy/redis/docker-compose.yml
├── start-redis.ps1
├── start-server.ps1
└── tests.ps1
```

## Запуск

**1. Redis** (нужен Docker):

```powershell
cd pz9-redis-cache
.\start-redis.ps1
```

**2. Сервер:**

```powershell
.\start-server.ps1
```

## Тесты

```powershell
cd pz9-redis-cache
.\tests.ps1
```

В логах сервера смотрите: попадание в кэш, промах, запись в кэш, сброс кэша.

## Запуск без PowerShell

Из каталога `pz9-redis-cache`.

**1. Redis:**

```text
docker compose -f deploy/redis/docker-compose.yml up -d
```

**2. Сервер (порт 8089):**

```text
go run ./cmd/server
```

## API

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/v1/tasks/{id}` | Задача по ID (кэш `tasks:task:<id>`) |
| GET | `/v1/tasks?page=1&limit=10` | Список (кэш `tasks:list:page=1:limit=10`) |
| PATCH | `/v1/tasks/{id}` | Обновить + инвалидация кэша |
| DELETE | `/v1/tasks/{id}` | Удалить + инвалидация кэша |

## Примеры запросов и ответов

Порт **8089**. В памяти задачи **id 1–2**.

### GET /v1/tasks/1

Первый вызов — промах кэша, чтение из хранилища. Второй — попадание в кэш (тот же JSON).

Ответ (`HTTP 200`):

```json
{
  "id": 1,
  "title": "Изучить Redis",
  "description": "Освоить cache-aside",
  "due_date": "2026-01-20T00:00:00Z"
}
```

### GET /v1/tasks?page=1&limit=10

Ответ (`HTTP 200`): JSON-массив задач.

### PATCH /v1/tasks/1

Тело: `{"id":1,"title":"Обновленная задача",...}`
Ответ (`HTTP 200`), тело пустое; кэш по id сбрасывается.

### DELETE /v1/tasks/1

Ответ (`HTTP 204`), тело пустое.

### GET /v1/tasks/1 после DELETE

Ответ (`HTTP 404`): текст `task not found`

## Кэш сбоку (кратко)

1. Читаем из Redis
2. **Попадание** — отдаём из кэша
3. **Промах** — читаем из хранилища в памяти
4. Кладём в Redis с TTL и jitter
5. При ошибке Redis — только хранилище (API не падает)

## Ключи Redis

| Ключ | Пример |
|------|--------|
| Задача | `tasks:task:1` |
| Список | `tasks:list:page=1:limit=10` |

## Переменные окружения

| Переменная | По умолчанию |
|------------|--------------|
| `SERVER_ADDR` | `:8089` |
| `REDIS_ADDR` | `localhost:6379` |
| `CACHE_TTL_SEC` | `120` |
| `CACHE_TTL_JITTER_SEC` | `30` |
| `LIST_CACHE_TTL_SEC` | `60` |
| `NEGATIVE_CACHE_TTL_SEC` | `30` |

## Дополнительные задания

### 1. Кэширование списка

`GET /v1/tasks` с query `page` и `limit`, ключ `tasks:list:page=1:limit=10`, TTL короче (60 с).

### 2. Явное логирование

В `task_service.go`: попадание в кэш, промах, запись, сброс кэша.

### 3. Построитель ключей и сериализатор

- `internal/cache/keys.go` — построитель ключей
- `internal/cache/serializer.go` — кодирование и декодирование JSON

### 4. Отрицательное кеширование

При 404 по ID в Redis кратко сохраняется маркер `__NOT_FOUND__` (TTL ~30 с + jitter), чтобы не обращаться к хранилищу повторно.

## Проверка без Redis

```powershell
cd deploy\redis
docker compose stop
curl.exe http://localhost:8089/v1/tasks/2
```

Сервис отвечает из репозитория, в логах — предупреждение при старте.

## Отчёт (кратко)

- Redis — ускоритель, не БД
- TTL — автоочистка устаревших данных
- Jitter — размазывает одновременное истечение ключей
- Инвалидация — консистентность после PATCH/DELETE
