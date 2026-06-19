# Практические работы по Go (2 семестр)

**Автор:** Ильин Владислав Викторович · группа ЭФМО-02-25

---

## О репозитории

Репозиторий содержит решения практических занятий №1–№16 по курсу разработки на Go.
Каждое занятие выделено в отдельный каталог формата `pzN-<тема>`, внутри которого находятся:

- исходный код;
- формулировка задания и примечания по реализации;
- инструкция по запуску и проверке (`README.md`);
- служебные скрипты (`start-*.ps1`, `tests.ps1`), если они предусмотрены.

---

## Список практических работ

| №  | Каталог                         | Краткое описание |
|----|---------------------------------|------------------|
| 1  | [pz1-microservices](pz1-microservices) | Взаимодействие HTTP-микросервисов |
| 2  | [pz2-grpc](pz2-grpc)                   | gRPC, Protocol Buffers |
| 3  | [pz3-logging](pz3-logging)             | Структурированное логирование (zap) |
| 4  | [pz4-monitoring](pz4-monitoring)       | Метрики Prometheus, дашборды Grafana |
| 5  | [pz5-security](pz5-security)           | HTTPS, безопасная работа с СУБД |
| 6  | [pz6-web-security](pz6-web-security)   | Защита веб-приложений (CSRF, XSS, cookie) |
| 7  | [pz7-docker](pz7-docker)               | Контейнеризация и docker-compose |
| 8  | [pz8-ci-cd](pz8-ci-cd)                 | CI/CD (GitHub Actions, GitLab CI) |
| 9  | [pz9-redis-cache](pz9-redis-cache)     | Распределённый кэш на Redis |
| 10 | [pz10-load-balancer](pz10-load-balancer) | Балансировка нагрузки (NGINX) |
| 11 | [pz11-graphql](pz11-graphql)           | GraphQL (gqlgen) |
| 12 | [pz12-rest-graphql](pz12-rest-graphql) | REST против GraphQL |
| 13 | [pz13-rabbitmq](pz13-rabbitmq)         | Очереди сообщений (RabbitMQ) |
| 14 | [pz14-job-queue](pz14-job-queue)       | Фоновые задачи, повторная обработка, DLQ |
| 15 | [pz15-vps-deploy](pz15-vps-deploy)     | Деплой на VPS (systemd) |
| 16 | [pz16-kubernetes](pz16-kubernetes)     | Запуск в Kubernetes |

---

## Сетевые порты

Для избежания конфликтов порты заранее разведены по работам.
Диапазон **8081–8098** резервируется под HTTP-сервисы, отдельные порты выделены под gRPC и внешнюю инфраструктуру.

> Не запускайте одновременно несколько работ, использующих один и тот же порт.

### Go‑приложения (HTTP / gRPC)

| №  | Каталог             | Компонент                         | Порт  |
|----|---------------------|-----------------------------------|-------|
| 1  | pz1-microservices   | user-service                      | **8081** |
| 1  | pz1-microservices   | order-service                     | **8082** |
| 2  | pz2-grpc            | gRPC-сервер                       | **50051** |
| 3  | pz3-logging         | HTTP API                          | **8083** |
| 4  | pz4-monitoring      | HTTP API                          | **8084** |
| 5  | pz5-security        | HTTP (редирект на HTTPS)          | **8085** |
| 5  | pz5-security        | HTTPS API                         | **8443** |
| 6  | pz6-web-security    | HTTP (веб-интерфейс)              | **8086** |
| 7  | pz7-docker          | auth                              | **8087** |
| 7  | pz7-docker          | tasks                             | **8088** |
| 8  | pz8-ci-cd           | —                                 | собственных портов нет (проверяются сервисы из pz7) |
| 9  | pz9-redis-cache     | HTTP API                          | **8089** |
| 10 | pz10-load-balancer  | NGINX (входной HTTP)              | **8090** |
| 10 | pz10-load-balancer  | tasks-1 (локальный экземпляр)     | **8091** |
| 10 | pz10-load-balancer  | tasks-2 (локальный экземпляр)     | **8092** |
| 10 | pz10-load-balancer  | tasks-3 (локальный экземпляр)     | **8093** |
| 11 | pz11-graphql        | GraphQL / песочница               | **8094** |
| 12 | pz12-rest-graphql   | REST + GraphQL                    | **8095** |
| 13 | pz13-rabbitmq       | tasks HTTP                        | **8096** |
| 14 | pz14-job-queue      | tasks HTTP                        | **8097** |
| 15 | pz15-vps-deploy     | tasks (локальный запуск)          | **8098** |
| 15 | pz15-vps-deploy     | tasks (по умолчанию на VPS)       | **8082** |
| 16 | pz16-kubernetes     | tasks (Service / port-forward)    | **8082** |

### Внешние сервисы (СУБД, брокеры, мониторинг)

| №  | Каталог            | Сервис                    | Порт            |
|----|--------------------|---------------------------|-----------------|
| 4  | pz4-monitoring     | Prometheus                | **9090**        |
| 4  | pz4-monitoring     | Grafana                   | **3000**        |
| 5  | pz5-security       | PostgreSQL                | **5432**        |
| 9  | pz9-redis-cache    | Redis                     | **6379**        |
| 10 | pz10-load-balancer | tasks (внутри Docker-сети)| **8082**        |
| 13 | pz13-rabbitmq      | RabbitMQ AMQP             | **5672**        |
| 13 | pz13-rabbitmq      | RabbitMQ Management UI    | **15672**       |
| 14 | pz14-job-queue     | RabbitMQ AMQP / UI        | **5672** / **15672** |

### Конфликтующие порты

| Порт              | Работы                                      | Комментарий |
|-------------------|---------------------------------------------|------------|
| **8082**          | pz1 (order-service), pz15 (VPS), pz16 (K8s), pz10 (Docker) | Используются в разных сценариях; при локальном тестировании включайте только одну из работ |
| **5672**, **15672** | pz13, pz14                                 | Достаточно одного экземпляра RabbitMQ на машину |

Точные настройки портов и переменных окружения указаны в `README.md` соответствующей работы.

---

## Как получить материалы

### Весь репозиторий целиком

Рекомендуемый способ — работать с полной копией репозитория:

- `git clone ...`
  или
- загрузка архива целиком с Git-сервера.

## Структура корня репозитория

```text
Go_Practice_2_sem/
├── README.md
├── .github/workflows/
├── .gitlab-ci.yml
├── pz1-microservices/
├── pz2-grpc/
└── …
```

Каталог `.github/workflows` и файл `.gitlab-ci.yml` относятся к работе №8 и обеспечивают запуск конвейеров CI при отправке изменений на GitHub и GitLab.

---

## Общая схема запуска и проверки

1. Перейти в каталог нужной работы (`pzN-...`).
2. Открыть и прочитать местный `README.md`.
3. Для автоматизированного запуска использовать скрипт `start-*.ps1` (PowerShell).
   Либо следовать разделу **«Запуск без PowerShell»** (CMD / bash).
4. При наличии — выполнить `tests.ps1` или команды проверки, описанные в `README`.

---
