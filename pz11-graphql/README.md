# Практическое задание 11

## ЭФМО-02-25 Ильин Владислав Викторович

---

# Информация о проекте

**pz11-graphql** — GraphQL API на Go (**gqlgen**): схема `Task`, запросы и мутации, песочница, Bearer для мутаций.

## Цели занятия

- Схема GraphQL: `Task`, `Query`, `Mutation`
- Генерация сервера через **gqlgen**
- Песочница GraphQL для проверки запросов
- Отличие от REST: клиент сам выбирает поля ответа

## Связь с другими практиками

| Практика | Что общего |
|----------|------------|
| **pz7-docker** | Домен **tasks** — список задач по API |
| **pz10-load-balancer** | Тот же сервис tasks, но REST; здесь — **GraphQL**-фасад |
| **pz9-redis-cache** | Кэш задач по ID; GraphQL читает тот же тип сущности |

GraphQL **не дублирует** REST-обработчики: используется общий сервисный слой и хранилище (как в методичке — одна бизнес-логика, разные протоколы).

## ВАЖНОЕ ПРИМЕЧАНИЕ

**8094** — песочница `http://localhost:8094/`, точка входа API: `http://localhost:8094/query`. На учебном сервере порт может отличаться.

## Файловая структура проекта

```
pz11-graphql/
├── graph/
│   ├── schema.graphqls
│   ├── schema.resolvers.go
│   ├── generated.go
│   └── model/
├── internal/
│   ├── store/          # в памяти (источник истины для учебки)
│   ├── service/        # общая бизнес-логика
│   └── auth/           # Bearer для mutations
├── cmd/graphql/main.go
├── gqlgen.yml
├── generate.ps1
├── start-server.ps1
└── tests.ps1
```

## Запуск

```powershell
cd pz11-graphql
.\generate.ps1      # после изменения schema.graphqls
.\start-server.ps1
```

Откройте песочницу: http://localhost:8094/

## Запуск без PowerShell

Из каталога `pz11-graphql` (порт **8094**).

После изменения `graph/schema.graphqls`:

```text
go run github.com/99designs/gqlgen generate
```

Сервер:

```text
go run ./cmd/graphql
```

Песочница: http://localhost:8094/

## Схема

```graphql
type Task {
  id: ID!
  title: String!
  description: String
  done: Boolean!
}

type Query {
  tasks: [Task!]!
  task(id: ID!): Task
}

type Mutation {
  createTask(input: CreateTaskInput!): Task!
  updateTask(id: ID!, input: UpdateTaskInput!): Task!
  deleteTask(id: ID!): Boolean!
}
```

## Песочница — примеры

**Список:**

```graphql
query {
  tasks {
    id
    title
    done
  }
}
```

**Одна задача:**

```graphql
query GetTask($id: ID!) {
  task(id: $id) {
    id
    title
    description
    done
  }
}
```

Variables: `{ "id": "t_001" }`

**Создание** (нужен заголовок `Authorization: Bearer demo-token` в песочнице):

```graphql
mutation Create($input: CreateTaskInput!) {
  createTask(input: $input) {
    id
    title
    done
  }
}
```

## Примеры ответов (GraphQL)

Порт **8094**, POST на `/query` или песочница. В хранилище задачи **t_001**, **t_002**.

### query { tasks { id title done } }

Ответ (`HTTP 200`):

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

### query { task(id: "t_001") { id title description done } }

Ответ (`HTTP 200`): объект `task` с полями или `"task": null`, если id нет.

### mutation createTask (с Bearer demo-token)

Ответ (`HTTP 200`): `{"data":{"createTask":{"id":"t_003","title":"...","done":false}}}`

Без токена на mutation → ошибка в поле `errors` GraphQL.

## Дополнительно: авторизация mutations

- **Query** — без токена
- **Мутация** — заголовок `Authorization: Bearer demo-token`

В песочнице: вкладка **Headers** (заголовки):

```json
{
  "Authorization": "Bearer demo-token"
}
```

## Тесты

```powershell
.\tests.ps1
```

## Порядок работы с gqlgen

1. Редактировать `graph/schema.graphqls`
2. `.\generate.ps1`
3. Дописать логику в `graph/schema.resolvers.go` (если generate перезаписал — восстановить вызовы service)
4. `go run ./cmd/graphql`

## GraphQL vs REST (кратко)

| | REST (pz7/pz10) | GraphQL |
|---|-----------------|---------|
| Точка входа | много URL | один `/query` |
| Данные | фиксированный JSON | клиент выбирает поля |
| Лишние поля в ответе | возможны | меньше лишних полей |

## Отчёт

- Скрин песочницы: `tasks`, `createTask`, `updateTask`
- Схема + сгенерированные файлы
- Пояснение связи с tasks из pz7/pz10
