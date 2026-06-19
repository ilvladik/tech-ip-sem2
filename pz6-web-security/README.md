# Практическое задание 6

## ЭФМО-02-25 Ильин Владислав Викторович

---

# Информация о проекте

**pz6-web-security** — веб-сервис на Go: защита от **CSRF** и **XSS**, безопасные **cookie** (`HttpOnly`, `SameSite=Lax`, опционально `Secure`), демонстрация небезопасного вывода HTML.

## Цели занятия

- Защита от CSRF через токен в форме и проверку на сервере
- Защита от XSS через `html/template` (автоэкранирование)
- Безопасные cookie: `HttpOnly`, `SameSite=Lax`, опционально `Secure`
- Демонстрация небезопасного вывода HTML для сравнения

## Файловая структура проекта

```
pz6-web-security/
├── cmd/server/main.go
├── internal/
│   ├── auth/          # CSRF-токены и session cookie
│   ├── config/        # SERVER_ADDR, SECURE_COOKIE
│   ├── httpapi/       # HTTP-обработчики
│   └── store/         # профили в памяти
├── templates/         # HTML-шаблоны
├── start-server.ps1
└── tests.ps1
```

## ВАЖНОЕ ПРИМЕЧАНИЕ

По умолчанию сервер слушает **http://localhost:8086** (чтобы не пересекаться с pz3–pz5). Для cookie с флагом `Secure` задайте `SECURE_COOKIE=true` (удобнее при HTTPS).

## Запуск

```powershell
cd pz6-web-security
.\start-server.ps1
```

Откройте в браузере: http://localhost:8086/login

Опционально — cookie с флагом `Secure` (удобнее при HTTPS):

```powershell
cd pz6-web-security
$env:SECURE_COOKIE = "true"
.\start-server.ps1
```

## Тесты

```powershell
cd pz6-web-security
.\tests.ps1
```

## Запуск без PowerShell

Из каталога `pz6-web-security` (порт **8086**):

```text
go run ./cmd/server
```

С флагом `Secure` для cookie (CMD):

```text
set SECURE_COOKIE=true
go run ./cmd/server
```

Linux / macOS:

```text
export SECURE_COOKIE=true
go run ./cmd/server
```

Браузер: http://localhost:8086/login

## API / маршруты

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/login` | Создать сессию, выдать cookie, редирект на `/profile` |
| GET | `/logout` | Удалить сессию, очистить cookie |
| GET | `/profile` | Форма редактирования имени |
| POST | `/profile` | Сохранить имя (проверка CSRF) |
| GET | `/hello` | Безопасное приветствие через шаблон |
| GET | `/comments` | Список комментариев |
| POST | `/comments` | Добавить комментарий (проверка CSRF) |
| GET | `/demo/unsafe` | **Учебный** небезопасный HTML (конкатенация строк) |

## Примеры запросов и ответов

Порт **8086**. Ответы в основном **HTML**, не JSON.

| Запрос | Код | Что в ответе |
|--------|-----|----------------|
| `GET /login` | 302 | Редирект на `/profile`, cookie `session_id` |
| `GET /profile` | 200 | HTML-форма с полем имени и скрытым `csrf_token` |
| `POST /profile` без CSRF | 403 | текст `invalid csrf token` |
| `POST /profile` с CSRF | 302 | Редирект, имя сохранено |
| `GET /hello` | 200 | HTML, имя из шаблона (экранировано) |
| `POST /comments` с `<script>...` | 302 | Комментарий в списке как **текст**, не выполняется |
| `GET /demo/unsafe?name=<b>XSS</b>` | 200 | HTML с **неэкранированным** вставленным именем (демо XSS) |

Пример небезопасного URL:

```text
http://localhost:8086/demo/unsafe?name=<script>alert('xss')</script>
```

## Дополнительные задания

### 1. Выход (`GET /logout`)

Удаляет профиль из памяти и сбрасывает cookie (`MaxAge: -1`), затем редирект на `/login`.

### 2. Флаг `Secure` для cookie

Переменная окружения `SECURE_COOKIE=true` включает `Secure` в session cookie. Для локальной работы по HTTP оставьте переменную пустой.

### 3. Ротация CSRF после POST

После успешного `POST /profile` и `POST /comments` генерируется новый CSRF-токен в хранилище — старый токен из формы больше не принимается.

### 4. Страница комментариев

`GET /comments` и `POST /comments` выводят и сохраняют комментарии через шаблон `comments.html`. Даже если в тексте есть `<script>`, браузер покажет его как текст, а не выполнит.

## Сравнение XSS

- **Безопасно:** `/hello`, `/comments` — `html/template`
- **Небезопасно (демо):** `/demo/unsafe?name=...` — конкатенация в HTML

Пример для демонстрации:

```
http://localhost:8086/demo/unsafe?name=<script>alert('xss')</script>
```

На `/hello` то же имя в профиле будет экранировано шаблоном.

## Коды ошибок

| Код | Когда |
|-----|--------|
| 400 | Пустое имя или комментарий, невалидная форма |
| 403 | Неверный CSRF-токен |
| 401 | Сессия не найдена после POST |
| 302 | Редиректы login / logout / успешные POST |
