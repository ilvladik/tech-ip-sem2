# Практическое задание 16

## ЭФМО-02-25 Ильин Владислав Викторович

---

# Информация о проекте

**pz16-kubernetes** — образ **tasks**, манифесты **Deployment** / **Service** / **ConfigMap**, **probes** и масштабирование в кластере.

## Цели занятия

- Собрать Docker-образ и выкатить приложение в Kubernetes.
- Настроить **readiness** и **liveness** через `/health`.
- Понять **Service**, **port-forward** и горизонтальное масштабирование реплик.

## Связь с курсом

| Практика | Что даёт |
|----------|----------|
| pz7 | Docker-образ tasks |
| pz8 | конвейер CI |
| pz15 | systemd на VPS |
| **pz16** | **K8s: Deployment, Service, ConfigMap, проверки готовности** |

## ВАЖНОЕ ПРИМЕЧАНИЕ

По умолчанию в манифестах сервис слушает **8082** (как в примере методички); при `kubectl port-forward` проверка: `http://localhost:8082/health`. Нужны **Docker** и кластер (**minikube** / **kind** / **k3s**).

## Файловая структура проекта

```
pz16-kubernetes/
├── cmd/tasks/main.go
├── Dockerfile
├── deploy/k8s/
│   ├── configmap.yaml
│   ├── deployment.yaml
│   └── service.yaml
├── build-image.ps1
└── apply.ps1
```

## Требования

- Docker
- Kubernetes (minikube / kind / k3s)
- kubectl

## Шаги

### 1. Собрать образ

```powershell
.\build-image.ps1
```

### 2. Загрузить в kind (если kind)

```powershell
kind load docker-image techip-tasks:0.1
```

### 3. Применить манифесты

```powershell
.\apply.ps1
```

### 4. Проверить

```powershell
kubectl get pods
kubectl logs <pod-name>
kubectl port-forward svc/tasks 8082:8082
```

В другом окне:

```powershell
curl.exe -i http://localhost:8082/health
```

### 5. Масштабирование

```powershell
kubectl scale deployment tasks --replicas=2
kubectl get pods
kubectl scale deployment tasks --replicas=1
```

### 6. Удалить

```powershell
kubectl delete -f deploy/k8s/service.yaml
kubectl delete -f deploy/k8s/deployment.yaml
kubectl delete -f deploy/k8s/configmap.yaml
```

## ConfigMap

- `TASKS_PORT=8082`
- `AUTH_BASE_URL=http://auth:8081`
- `LOG_LEVEL=info`

## Probes

- **readiness** — `/health`, сервис получает трафик когда готов
- **liveness** — `/health`, перезапуск «зависшего» контейнера

## Образ

Тег **`techip-tasks:0.1`** (не `latest`) — как в методичке.

## Запуск без PowerShell

Из каталога `pz16-kubernetes`.

```text
docker build -t techip-tasks:0.1 .
kind load docker-image techip-tasks:0.1
kubectl apply -f deploy/k8s/
kubectl get pods
kubectl port-forward svc/tasks 8082:8082
```

Проверка: http://localhost:8082/health

## Примеры запросов и ответов

После `kubectl port-forward svc/tasks 8082:8082`:

```bash
curl http://localhost:8082/health
```

Ответ (`HTTP 200`):

```json
{"status":"ok"}
```

Масштабирование:

```text
kubectl scale deployment tasks --replicas=2
kubectl get pods
```

## Локально без K8s

```text
go run ./cmd/tasks
```

Порт по умолчанию — 8082 (или задать `TASKS_PORT`).
