# Практическое задание 15

## ЭФМО-02-25 Ильин Владислав Викторович

---

# Информация о проекте

**pz15-vps-deploy** — сборка бинарника **tasks** под Linux, unit **systemd**, переменные в `/etc/tasks/`, схема деплоя на **VPS** (без обязательного облака).

## Цели занятия

- Собрать статический бинарник Go для Linux и доставить на сервер.
- Оформить сервис как **systemd**-unit с перезапуском и логами journald.
- Разделить конфиг (env-файл) и секреты вне репозитория.

## Связь с курсом

Собирает **tasks** из предыдущих практик (pz7, pz10, pz13…) в **production-like** схему: бинарник на Linux, конфиг в `/etc/tasks/`, управление через `systemctl`.

## Файловая структура проекта

| Файл | Назначение |
|------|------------|
| `cmd/tasks/main.go` | Минимальный сервис с `GET /health` |
| `deploy/systemd/tasks.service` | Unit для systemd |
| `deploy/env/tasks.env.example` | Пример переменных |
| `build-linux.ps1` | Сборка `bin/tasks` под Linux |

## ВАЖНОЕ ПРИМЕЧАНИЕ

Локальная проверка на Windows — порт **8098** (чтобы не пересекаться с другими практиками). На VPS порты и пользователь (`tasks` vs `root`) задаются по политике сервера.

## Запуск (локальная проверка, Windows)

```powershell
$env:TASKS_PORT = "8098"
go run ./cmd/tasks
curl.exe http://localhost:8098/health
```

## Запуск без PowerShell

Из каталога `pz15-vps-deploy` (локально порт **8098**):

CMD:

```text
set TASKS_PORT=8098
go run ./cmd/tasks
```

Linux / macOS:

```text
export TASKS_PORT=8098
go run ./cmd/tasks
```

Проверка: http://localhost:8098/health

## Примеры запросов и ответов

### GET /health

```bash
curl http://localhost:8098/health
```

Ответ (`HTTP 200`):

```json
{"status":"ok"}
```

## Деплой на VPS (краткий чеклист)

### 1. Сборка на ПК

```powershell
.\build-linux.ps1
```

### 2. Пользователь и каталоги на VPS

```bash
sudo useradd --system --no-create-home --shell /usr/sbin/nologin tasksuser
sudo mkdir -p /opt/tasks /etc/tasks
sudo chown -R tasksuser:tasksuser /opt/tasks
```

### 3. Конфиг

```bash
sudo cp tasks.env.example /etc/tasks/tasks.env
sudo chmod 600 /etc/tasks/tasks.env
```

### 4. Загрузка бинарника

```powershell
scp bin/tasks user@<VPS_IP>:/tmp/tasks
```

На VPS:

```bash
sudo mv /tmp/tasks /opt/tasks/tasks
sudo chown tasksuser:tasksuser /opt/tasks/tasks
sudo chmod 755 /opt/tasks/tasks
```

### 5. systemd

```bash
sudo cp tasks.service /etc/systemd/system/tasks.service
sudo systemctl daemon-reload
sudo systemctl enable tasks
sudo systemctl start tasks
sudo systemctl status tasks
```

### 6. Логи

```bash
sudo journalctl -u tasks -f
```

### 7. Проверка

```bash
curl -i http://127.0.0.1:8082/health
```

## Обновление (rolling)

```bash
sudo systemctl stop tasks
sudo mv /opt/tasks/tasks /opt/tasks/tasks.old
sudo mv /tmp/tasks /opt/tasks/tasks
sudo chown tasksuser:tasksuser /opt/tasks/tasks
sudo systemctl start tasks
```

## Откат

```bash
sudo systemctl stop tasks
sudo mv /opt/tasks/tasks.old /opt/tasks/tasks
sudo systemctl start tasks
```

## Почему не root

Сервис работает от `tasksuser` — меньше рисков при компрометации процесса.

## Порт на VPS

По умолчанию в `tasks.env`: **8082** (как в методичке). Локально для теста можно **8098**.
