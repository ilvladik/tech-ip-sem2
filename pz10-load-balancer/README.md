# Практическое задание 10

## ЭФМО-02-25 Ильин Владислав Викторович

---

# Информация о проекте

**pz10-load-balancer** — несколько реплик одного сервиса **tasks** и **NGINX** как балансировщик нагрузки (round-robin), заголовок `X-Instance-ID`.

## Цели занятия

- Запустить несколько реплик одного сервиса без локального состояния
- Настроить **upstream** в NGINX (распределение по кругу)
- Проверить балансировку по заголовку `X-Instance-ID`
- Реализовать `GET /health` для проверки живости

## ВАЖНОЕ ПРИМЕЧАНИЕ

| Компонент | Порт | Описание |
|-----------|------|----------|
| NGINX (вход) | **8090** | Единая точка входа |
| tasks-1 | 8091 | Локально без Docker |
| tasks-2 | 8092 | Локально без Docker |
| tasks-3 | 8093 | Локально без Docker |
| tasks в Docker | 8082 | Внутри сети compose |

Внешний порт **8090** не пересекается с pz3–pz5 (8083–8085).

## Файловая структура проекта

```
pz10-load-balancer/
├── services/tasks/          # API задач без локального состояния
├── deploy/lb/
│   ├── nginx.conf
│   └── docker-compose.yml   # tasks_1, tasks_2, tasks_3 + nginx
├── start-instances.ps1      # 3 процесса без Docker
├── start-compose.ps1        # NGINX + 3 контейнера
├── tests-instances.ps1
└── tests-lb.ps1
```

## Запуск без Docker (локально)

```powershell
cd pz10-load-balancer
.\start-instances.ps1
```

## Тесты (локальные реплики)

```powershell
cd pz10-load-balancer
.\tests-instances.ps1
```

В логах каждого окна — `instance=tasks-N method=...` (логирование запросов).

## Запуск с Docker + NGINX

```powershell
cd pz10-load-balancer
.\start-compose.ps1
```

## Тесты (через NGINX)

```powershell
cd pz10-load-balancer
.\tests-lb.ps1
```

Остановка:

```powershell
cd pz10-load-balancer
.\stop-compose.ps1
```

## Запуск без PowerShell

**Три копии tasks (8091–8093)** — три терминала, из `services/tasks`:

CMD: `set INSTANCE_ID=tasks-1` и `set APP_PORT=8091`, затем `go run ./cmd/server`
(для копий 2 и 3 — порты 8092, 8093).

Linux / macOS: `export INSTANCE_ID=tasks-1` и `export APP_PORT=8091`, затем `go run ./cmd/server`.

**NGINX + Docker:**

```text
docker compose -f deploy/lb/docker-compose.yml up --build
```

Вход: http://localhost:8090

Проверка распределения по кругу:

```powershell
1..10 | ForEach-Object { curl.exe -s -D - http://localhost:8090/whoami -o NUL | Select-String "X-Instance-ID" }
```

Должны чередоваться `tasks-1`, `tasks-2`, `tasks-3`.

## API (каждая реплика)

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/health` | `{"status":"ok","instance":"tasks-1"}` + `X-Instance-ID` |
| GET | `/whoami` | `{"instance":"tasks-1"}` |
| GET | `/v1/tasks` | Список задач + `X-Instance-ID` |

Переменные: `INSTANCE_ID`, `APP_PORT`.

## Примеры запросов и ответов

Через NGINX (**8090**) или напрямую на реплику (**8091**).

### GET /health

```bash
curl http://localhost:8090/health
```

Ответ (`HTTP 200`):

```json
{"status":"ok","instance":"tasks-1"}
```

Заголовок: `X-Instance-ID: tasks-1` (номер реплики меняется при балансировке).

### GET /whoami

Ответ (`HTTP 200`): `{"instance":"tasks-2"}`

### GET /v1/tasks

Ответ (`HTTP 200`):

```json
[
  {"id": 1, "title": "Изучить NGINX"},
  {"id": 2, "title": "Освоить load balancing"}
]
```

## Проверка отказоустойчивости

```powershell
cd deploy\lb
docker compose stop tasks_1
# снова tests-lb.ps1 — только tasks-2 и tasks-3
docker compose start tasks_1
```

## NGINX upstream

```nginx
upstream tasks_backend {
    server tasks_1:8082;
    server tasks_2:8082;
    server tasks_3:8082;
}
```

`proxy_pass` проксирует `Authorization`, `X-Request-ID`, `X-Forwarded-For`.

## Дополнительные задания

### 1. Третья реплика tasks-3

В `docker-compose.yml` и `nginx.conf` добавлен `tasks_3`.

### 2. Логирование входящих запросов

Промежуточный слой в `cmd/server/main.go` пишет в лог: инстанс, метод, путь, адрес клиента, длительность.

### 3. GET /whoami

Отдельный метод `/whoami` для проверки, какая копия ответила.

### 4. Общее состояние (конспект)

Для горизонтального масштабирования сервис должен быть **без локального состояния**:

- данные — в общей БД или Redis (см. pz9-redis-cache);
- состояние в памяти одной реплики ломает балансировку;
- сессии — во внешнем хранилище, не в памяти процесса.

## CI vs отчёт

- Скриншоты: `docker compose ps`, curl с разными `X-Instance-ID`
- Без локального состояния: одна и та же задача с любой реплики

## Контрольные вопросы (кратко)

| Вопрос | Ответ |
|--------|--------|
| Вертикальное vs горизонтальное | Больше CPU/RAM vs больше реплик |
| Зачем LB | Одна точка входа, распределение нагрузки |
| Upstream (NGINX) | Группа серверов приложения за балансировщиком |
| Без локального состояния | Нет сессии в памяти процесса — любая реплика может ответить |
| X-Instance-ID | Видно, какой инстанс обработал запрос |
