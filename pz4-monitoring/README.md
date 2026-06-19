# Практическое задание 4

## ЭФМО-02-25 Ильин Владислав Викторович

---

# Информация о проекте

**pz4-monitoring** — HTTP-сервис на Go с экспортом метрик Prometheus и дашбордом Grafana.
Собираются счётчики запросов и ошибок, гистограммы длительности, метрики по студентам и датчик активных запросов.

## Цели занятия

- Подключить приложение к Prometheus: эндпоинт `/metrics`, типы метрик (counter, histogram, gauge).
- Считать HTTP-запросы, ошибки и задержки; разделять бизнес-метрики по студентам.
- Настроить сбор в Prometheus и визуализацию в Grafana (дашборд).
- Сопоставить подходы «логи» (pz3) и «метрики» (pz4).

## Файловая структура проекта

```
pz4-monitoring/
├── cmd/server/main.go
├── internal/httpapi/
├── internal/metrics/metrics.go
├── internal/student/
├── monitoring/
│   ├── prometheus.yml
│   ├── prometheus.docker.yml
│   └── grafana/
├── docker-compose.yml
├── start-server.ps1
├── start-monitoring.ps1
├── generate-traffic.ps1
└── tests.ps1
```

## ВАЖНОЕ ПРИМЕЧАНИЕ

| Сервис | Порт |
|--------|------|
| Go-приложение | **8084** |
| Prometheus | 9090 |
| Grafana | 3000 |

В `monitoring/prometheus.yml` цель (target) должна указывать на работающее приложение **:8084** (или имя сервиса в Docker-сети). На учебном сервере порты могут отличаться — уточните у преподавателя.

Для дашборда Grafana сначала поднимите стек мониторинга (`.\start-monitoring.ps1`), затем при необходимости сгенерируйте трафик (`.\generate-traffic.ps1`).

## Тестовые данные

Студенты в памяти: **id 1–3** (как в pz3). В `tests.ps1` запрос **`/students/999`** — проверка **404** (записи нет).

## Запуск

**1. Go-приложение**

```powershell
cd pz4-monitoring
.\start-server.ps1
```

**2. Prometheus + Grafana (Docker)**

```powershell
.\start-monitoring.ps1
.\generate-traffic.ps1
```

- Prometheus: http://localhost:9090 → раздел «Цели» → `go_app` в статусе «работает»
- Grafana: http://localhost:3000 (логин `admin` / `admin`)
- Дашборд **PZ4 Go App Monitoring** подхватывается автоматически

Метрики приложения: http://localhost:8084/metrics

## Запуск без PowerShell

Из каталога `pz4-monitoring`.

**Приложение Go (порт 8084):**

```text
go run ./cmd/server
```

**Prometheus и Grafana (Docker):**

```text
docker compose up -d
```

Перед этим в `monitoring/prometheus.yml` указан target `localhost:8084`. Метрики: http://localhost:8084/metrics

**Prometheus без Docker** (если установлен локально):

```powershell
prometheus --config.file=monitoring/prometheus.yml
```

## Тесты

```powershell
cd pz4-monitoring
.\tests.ps1
```

## API

| Метод | URL |
|-------|-----|
| GET | `/health` |
| GET | `/students/{id}` |
| GET | `/metrics` |

## Примеры запросов и ответов

Порт **8084**. Формат как в pz3 (студенты **1–3** в памяти).

### GET /health

Ответ (`HTTP 200`): `{"status":"ok"}`

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

### GET /students/999

Ответ (`HTTP 404`):

```json
{"error":"студент с таким id не найден в системе","id":999}
```

### GET /metrics

Ответ (`HTTP 200`): текст Prometheus, фрагмент:

```text
app_http_requests_total{method="GET",path="/health"} 1
app_http_errors_total{...} 0
```

## Метрики (база)

| Метрика | Тип | Назначение |
|---------|-----|------------|
| `app_http_requests_total` | Counter | все HTTP-запросы (method, path) |
| `app_http_errors_total` | Counter | ответы ≥ 400 (method, path, status_code) |
| `app_http_request_duration_seconds` | Histogram | время обработки запроса |

Пути нормализуются: `/students/1` → `/students/{id}` в labels.

## Дополнительные задания (выполнено)

### `app_student_requests_total`

**Как работает:** при каждом обращении к `GET /students/{id}` счётчик увеличивается с label `student_id` — видно, сколько раз запрашивали конкретного студента.

### `app_student_handler_duration_seconds`

**Как работает:** отдельная гистограмма только для обработчика студента — измеряет время бизнес-логики, не весь HTTP-стек.

### `app_active_requests` (Gauge)

**Как работает:** промежуточный слой увеличивает счётчик в начале запроса и уменьшает в конце — показывает, сколько запросов обрабатывается прямо сейчас.

### Дашборд ошибок в Grafana

**Как работает:** в `pz4-dashboard.json` панели: ошибки по `status_code` (pie), запросы по `path`, средняя latency, запросы по `student_id`, gauge активных запросов.

## PromQL (из методички)

```promql
sum(app_http_requests_total)
sum(app_http_errors_total)
sum by (path) (app_http_requests_total)
sum(rate(app_http_request_duration_seconds_sum[1m])) / sum(rate(app_http_request_duration_seconds_count[1m]))
sum by (status_code) (app_http_errors_total)
```

## Метрики vs логи (практика 3)

| Логи (pz3) | Метрики (pz4) |
|------------|---------------|
| отдельные события | агрегаты во времени |
| текст/JSON строка | числа для графиков |
| «что случилось» | «как часто и как долго» |
