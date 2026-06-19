# Практическое задание 1

## ЭФМО-02-25 Ильин Владислав Викторович

---

# Информация о проекте

**pz1-microservices** — учебный проект из двух микросервисов на Go, связанных HTTP-запросами.
**user-service** хранит и отдаёт данные о пользователях, **order-service** — заказы и при необходимости запрашивает пользователя у первого сервиса. Хранение в памяти, без базы данных.

## Цели занятия

- Понять разницу между монолитом и микросервисной архитектурой.
- Разделить систему по зонам ответственности (пользователи / заказы).
- Реализовать HTTP-взаимодействие между сервисами.
- Научиться передавать и обрабатывать JSON.
- Отработать запуск, тестирование и проверку взаимодействия через curl/Postman.

## Файловая структура проекта

```
pz1-microservices/
├── start-servers.ps1      # запуск обоих сервисов
├── tests.ps1              # проверка всех эндпоинтов
├── user-service/
│   ├── cmd/server/main.go
│   ├── internal/user/
│   │   ├── model.go
│   │   ├── repo.go
│   │   └── handler.go
│   ├── pkg/middleware/logging.go
│   └── go.mod
└── order-service/
    ├── cmd/server/main.go
    ├── internal/order/
    │   ├── model.go
    │   ├── repo.go
    │   ├── handler.go
    │   └── client.go
    ├── pkg/middleware/logging.go
    └── go.mod
```

## ВАЖНОЕ ПРИМЕЧАНИЕ

Для локальной разработки используются порты **8081** (user-service) и **8082** (order-service).
Если практики запускаются на учебном сервере, порты могут отличаться — уточните актуальные значения у преподавателя.

Перед проверкой **order-service** с маршрутом `/orders/{id}/full` обязательно должен быть запущен **user-service**.

## Тестовые данные

| Сущность | id в памяти | Пример запроса |
|----------|-------------|----------------|
| Пользователь | 1, 2, 3 | `GET /users/1` → JSON пользователя |
| Заказ | 101, 102, 103 | `GET /orders/101` → JSON заказа |

Запросы с **id 999** в `tests.ps1` — **намеренная проверка ошибки**: такого id в памяти нет. Ответ **404** и JSON вида `{"error":"...","id":999}`. Если запросить **существующий** id (например `/users/1`), вернётся обычный JSON с данными.

## Запуск

Оба сервиса сразу (откроются 2 окна PowerShell):

```powershell
cd pz1-microservices
.\start-servers.ps1
```

Или вручную в двух терминалах:

```powershell
cd user-service
go run ./cmd/server
```

```powershell
cd order-service
go run ./cmd/server
```

Переменная окружения для order-service:

```powershell
$env:USER_SERVICE_URL = "http://localhost:8081"
go run ./cmd/server
```

## Запуск без PowerShell

Из каталога `pz1-microservices`. Нужны **два** терминала.

**Терминал 1 — user-service (порт 8081):**

```text
cd user-service
go run ./cmd/server
```

**Терминал 2 — order-service (порт 8082):**

CMD (Windows):

```text
cd order-service
set USER_SERVICE_URL=http://localhost:8081
go run ./cmd/server
```

Linux / macOS:

```text
cd order-service
export USER_SERVICE_URL=http://localhost:8081
go run ./cmd/server
```

## Тесты

```powershell
cd pz1-microservices
.\tests.ps1
```

## API

### user-service (:8081)

| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/users` | список пользователей |
| GET | `/users/{id}` | пользователь по ID |

### order-service (:8082)

| Метод | URL | Описание |
|-------|-----|----------|
| GET | `/orders/{id}` | заказ по ID |
| GET | `/orders/{id}/full` | заказ + данные пользователя |
| GET | `/orders/by-user/{userID}` | все заказы пользователя |

## Примеры запросов

### Список пользователей

```bash
curl http://localhost:8081/users
```

### Получение пользователя по ID

```bash
curl http://localhost:8081/users/1
```

Ответ:

```json
{"id":1,"name":"Иван Иванов","email":"ivan@example.com"}
```

### Пользователь не найден (404)

```bash
curl http://localhost:8081/users/999
```

Ответ (`HTTP 404`):

```json
{"error":"пользователь с таким id не найден в системе","id":999}
```

### Получение заказа по ID

```bash
curl http://localhost:8082/orders/101
```

### Заказы пользователя

```bash
curl http://localhost:8082/orders/by-user/1
```

### Заказ с данными пользователя

```bash
curl http://localhost:8082/orders/101/full
```

Ответ:

```json
{
  "order": {
    "id": 101,
    "user_id": 1,
    "item": "Ноутбук",
    "price": 79990
  },
  "user": {
    "id": 1,
    "name": "Иван Иванов",
    "email": "ivan@example.com"
  }
}
```

### user-service недоступен (502)

Остановите user-service и выполните:

```bash
curl http://localhost:8082/orders/101/full
```

## Тестовые данные

**Пользователи (user-service):**

| ID | Имя               | Email              |
|----|-------------------|--------------------|
| 1  | Иван Иванов       | ivan@example.com   |
| 2  | Мария Петрова     | maria@example.com  |
| 3  | Александр Сергеев | alex@example.com   |

**Заказы (order-service):**

| ID  | user_id | Товар      | Цена  |
|-----|---------|------------|-------|
| 101 | 1       | Ноутбук    | 79990 |
| 102 | 2       | Книга      | 2490  |
| 103 | 1       | Клавиатура | 5990  |

## Дополнительные задания (выполнено)

- **`GET /users`** — `ListAll()` обходит map в repo и отдаёт JSON-массив всех пользователей.
- **`GET /orders/by-user/{userID}`** — `GetByUserID()` фильтрует заказы по `user_id` без обращения к user-service.
- **Промежуточный слой логирования** — оборачивает `ServeMux`, после каждого запроса пишет в консоль метод, путь и время.
- **`USER_SERVICE_URL`** — order-service читает адрес user-service из переменных окружения при старте; если пусто — `http://localhost:8081`.

## Объяснение, как обрабатываются ошибки и коды ответа

Ошибки возвращаются через `http.Error()` с текстовым телом и соответствующими HTTP-статусами:

- **200 OK** — успешный GET-запрос, тело в формате JSON.
- **400 Bad Request** — не указан или неверный ID.
- **404 Not Found** — пользователь или заказ не найден.
- **405 Method Not Allowed** — метод отличен от GET.
- **502 Bad Gateway** — order-service не смог получить данные от user-service.

В консоли сервисов промежуточный слой выводит лог каждого запроса: метод, путь и время выполнения.
