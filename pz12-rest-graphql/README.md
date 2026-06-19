# Практическое задание 12

## ЭФМО-02-25 Ильин Владислав Викторович

---

# Информация о проекте

**pz12-rest-graphql** — один процесс с **REST** и **GraphQL** над общим хранилищем задач: сравнение избыточности ответа, сценариев и ошибок.

## Цели занятия

Сравнить REST и GraphQL на **одних и тех же** сценариях:

1. Список задач (нужны только `id`, `title`, `done`)
2. Детали задачи (все поля)
3. Создание задачи
4. Обновление `done`
5. Ошибка «не найдено»

## Связь с другими практиками

| Практика | Роль |
|----------|------|
| **pz7 / pz10** | REST-сервис **tasks** |
| **pz11** | GraphQL + gqlgen для **Task** |
| **pz12** | **Оба API в одном процессе**, общее хранилище и сервисный слой |

Один источник данных — честное сравнение без расхождения между сервисами.

## ВАЖНОЕ ПРИМЕЧАНИЕ

**8095** — REST и GraphQL на одном сервере (см. таблицу URL ниже). На учебном сервере порт может отличаться.

## Порт и URL

**8095** — REST и GraphQL на одном сервере:

| API | URL |
|-----|-----|
| REST | `http://localhost:8095/v1/tasks` |
| GraphQL | `http://localhost:8095/graphql` |
| Песочница GraphQL | `http://localhost:8095/` |

## Запуск

```powershell
cd pz12-rest-graphql
.\start-server.ps1
```

## Тесты

```powershell
cd pz12-rest-graphql
.\tests.ps1
```

## Запуск без PowerShell

Из каталога `pz12-rest-graphql` (порт **8095**):

```text
go run ./cmd/server
```

REST: http://localhost:8095/v1/tasks
GraphQL / песочница: http://localhost:8095/

## Сценарии сравнения

### 1. Список (избыточные поля в ответе)

**REST** — всегда полный JSON (`description` тоже):

```powershell
curl.exe http://localhost:8095/v1/tasks
```

**GraphQL** — только нужные поля:

```graphql
query {
  tasks {
    id
    title
    done
  }
}
```

**Вывод:** для списка GraphQL отдаёт меньше лишних данных.

Ответ REST (`HTTP 200`):

```json
[
  {
    "id": "t_001",
    "title": "Первая задача",
    "description": "Учебный пример",
    "done": false
  },
  {
    "id": "t_002",
    "title": "Вторая задача",
    "description": "Проверка API",
    "done": true
  }
]
```

Ответ GraphQL (`HTTP 200`):

```json
{
  "data": {
    "tasks": [
      {"id": "t_001", "title": "Первая задача", "done": false},
      {"id": "t_002", "title": "Вторая задача", "done": true}
    ]
  }
}
```

### 2. Детали задачи

**REST:** `GET /v1/tasks/t_001`

**GraphQL:**

```graphql
query ($id: ID!) {
  task(id: $id) {
    id
    title
    description
    done
  }
}
```

### 3. Создание

**REST:**

```powershell
curl.exe -X POST http://localhost:8095/v1/tasks `
  -H "Content-Type: application/json" `
  -d "{\"title\":\"Сравнение REST и GraphQL\",\"description\":\"ПЗ12\"}"
```

**GraphQL:** мутация `createTask` в песочнице GraphQL.

Ответ REST (`HTTP 201`):

```json
{
  "id": "t_003",
  "title": "Сравнение REST и GraphQL",
  "description": "ПЗ12",
  "done": false
}
```

### 4. Обновление

**REST:** `PATCH /v1/tasks/t_001` с `{"done":true}`

**GraphQL:** мутация `updateTask`.

### 5. Ошибки

| Ситуация | REST | GraphQL |
|----------|------|---------|
| Неизвестный id | `404` + `{"error":"task not found"}` | `200` + `"task": null` |

Пример REST:

```bash
curl http://localhost:8095/v1/tasks/unknown
```

Ответ (`HTTP 404`): `{"error":"task not found"}`

Пример GraphQL: `task(id: "unknown")` → `"data": { "task": null }`

## Итоговая таблица (для отчёта)

| Критерий | REST | GraphQL |
|----------|------|---------|
| Точки входа | Несколько URL | Один `/graphql` |
| Выбор полей | Фиксированный ответ | Клиент выбирает поля |
| Лишние поля в ответе | Часто (список с `description`) | Меньше |
| Несколько запросов | Нужны отдельные вызовы при связях | Один запрос |
| Ошибки | HTTP-коды | Часто `200` и массив `errors` |
| Документация | OpenAPI / curl | Схема и песочница GraphQL |
| Кэширование | Удобно по URL | Сложнее |
| Когда удобнее | простой CRUD, публичные API | сложные клиенты, мобильные интерфейсы |

## Краткий вывод (шаблон)

- **REST** удобен для простого CRUD, привычен, хорошо кэшируется по URL.
- **GraphQL** удобен, когда клиенту нужны разные наборы полей без лишнего трафика.
- Для учебного **tasks**-сервиса REST проще; GraphQL выигрывает на сценарии «список без description».

## Файловая структура проекта

```
pz12-rest-graphql/
├── cmd/server/main.go      # REST + GraphQL
├── internal/
│   ├── store/              # общие данные
│   ├── service/            # общая логика
│   └── rest/               # обработчики REST
├── graph/                  # gqlgen (из pz11)
├── tests.ps1
└── README.md
```

## gqlgen

После изменения `graph/schema.graphqls`:

```powershell
.\generate.ps1
```
