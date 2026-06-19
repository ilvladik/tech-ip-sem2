# Практическое задание 8

## ЭФМО-02-25 Ильин Владислав Викторович

---

# Информация о проекте

**pz8-ci-cd** — настройка **CI/CD**: GitHub Actions и GitLab CI для тестов, сборки Go и образов Docker сервисов из **pz7-docker**. Собственного HTTP-сервиса в этой папке нет.

## Цели занятия

- Понять разницу **CI** (тесты, сборка) и **CD** (доставка, деплой)
- Настроить автоматический конвейер при отправке кода в репозиторий
- Собирать Docker-образы в CI (на исполнителях GitHub)
- Хранить секреты вне репозитория

## Файловая структура проекта

| Расположение | Содержимое |
|--------------|------------|
| Корень `Go_Practice_2_sem/` | `.github/workflows/ci.yml`, `.github/workflows/docker-publish.yml`, `.gitlab-ci.yml` |
| Каталог `pz8-ci-cd/` | `verify-ci.ps1`, `README.md` |

## На чём строится конвейер CI

CI настроен для сервисов из **pz7-docker**:

| Сервис | Путь | Порт (runtime) |
|--------|------|----------------|
| auth | `pz7-docker/services/auth` | 8087 |
| tasks | `pz7-docker/services/tasks` | 8088 |

## Файлы CI в корне репозитория

Для работы проекта (практика 8, отправка на GitHub/GitLab) в **корне** `Go_Practice_2_sem/` должны лежать:

| Файл |
|------|
| `.github/workflows/ci.yml` |
| `.github/workflows/docker-publish.yml` |
| `.gitlab-ci.yml` |

Если их нет — скачай из репозитория и положи в **корень**, не в папку `pz8-ci-cd/`.

```
Go_Practice_2_sem/
├── .github/workflows/ci.yml
├── .github/workflows/docker-publish.yml
├── .gitlab-ci.yml
└── pz8-ci-cd/
```

## Что делает CI (`ci.yml`)

**Задания `test-auth` и `test-tasks`** (параллельно):

1. загрузка кода из репозитория
2. установка Go 1.22
3. `go mod tidy`
4. `go test ./...`
5. `go build ./...`

**Задание `docker-build`** (после успешных тестов):

- `docker build` для auth и tasks
- тег образа: `${{ github.sha }}`

## Дополнительные задания

### Публикация в реестр образов (`docker-publish.yml`)

Запуск вручную из GitHub → **Actions** → **Docker Publish**.

Секреты в **Settings → Secrets**:

| Secret | Назначение |
|--------|------------|
| `REGISTRY_USERNAME` | логин GHCR/Docker Hub |
| `REGISTRY_PASSWORD` | токен |

Образы: `ghcr.io/<owner>/techip-auth:<sha>`, `ghcr.io/<owner>/techip-tasks:<sha>`.

### GitLab CI (`.gitlab-ci.yml`)

Стадии: `test` → `docker`. Тег: `$CI_COMMIT_SHORT_SHA`.

### Деплой на VPS (описание)

После отправки образа в реестр на сервере:

```bash
docker pull ghcr.io/my-org/techip-tasks:<tag>
cd deploy && docker compose up -d
```

Для учебной сдачи достаточно скриншота успешного задания в GitHub Actions.

## Запуск (локальная имитация CI)

```powershell
cd pz8-ci-cd
.\verify-ci.ps1
```

Повторяет шаги test/build из CI за несколько секунд (аналог раздела **«Тесты»** в других практиках).

## Запуск без PowerShell

Собственных сервисов нет. Локальная проверка — те же шаги, что в CI, для **pz7-docker**:

```text
cd pz7-docker/services/auth
go mod tidy
go test ./...
go build ./...

cd ../tasks
go mod tidy
go test ./...
go build ./...
```

Файлы `.github/workflows/` и `.gitlab-ci.yml` должны лежать в **корне** репозитория (см. выше).

## Запуск на GitHub

1. Закоммитьте и запушьте репозиторий
2. Откройте вкладку **Actions**
3. Сценарий **CI Pipeline** запускается на ветках `main` / `master` и при запросе на слияние (pull request)

## Примеры запросов и ответов

Собственного HTTP API нет. Форматы ответов сервисов **auth** / **tasks** — как в [pz7-docker](../pz7-docker/README.md#примеры-запросов-и-ответов) (порты **8087** / **8088** после `.\start-compose.ps1` в pz7).

## CI vs CD

| | CI | CD |
|---|----|----|
| Когда | каждый коммит / запрос на слияние | после успешного CI |
| Что | тесты, сборка, docker build | публикация в реестр, деплой |
| В этой работе | `ci.yml` | `docker-publish.yml` + compose на VPS |

## Секреты — правила

- не коммитить пароли и токены в YAML
- использовать GitHub Secrets / GitLab Variables
- не хранить `.env` с секретами в git

## Типичные ошибки

| Ошибка | Решение |
|--------|---------|
| `go test` падает | добавить модульные тесты, проверить `verify-ci.ps1` |
| неверный путь | `working-directory: pz7-docker/services/tasks` |
| ошибка Docker build | Dockerfile в каталоге сервиса, контекст сборки = этот каталог |
| конвейер красный без тестов | временно оставить только `go build`, потом вернуть тесты |

## Контрольные вопросы (кратко)

- **Образ и контейнер** — шаблон и запущенный экземпляр
- **Зачем CI** — раннее обнаружение поломок, единый стандарт сборки
- **Секреты** — доступ конвейера к реестру/SSH без утечки в код
- **Тег образа** — `github.sha` / `CI_COMMIT_SHORT_SHA` для трассировки версии
