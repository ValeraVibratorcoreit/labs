# tasks-api

Учебный Go-бэкенд из задания «Рабочий пример целиком»: HTTP API + PostgreSQL + pgxpool.

## Запуск

1. Запустить PostgreSQL и Adminer:

```bash
docker compose up -d
```

2. Установить зависимости:

```bash
go mod tidy
```

3. Запустить сервер:

```bash
go run .
```

Сервер: `http://localhost:8080`  
Adminer: `http://localhost:8081`

По умолчанию приложение использует:

```text
DATABASE_URL=postgres://app:secret@localhost:5432/app?sslmode=disable
PORT=8080
```

## API

### GET /tasks

Получить все задачи.

**200**

```json
[
  {"id":1,"title":"Сверстать карточку","done":false,"created_at":"2026-09-21T09:38:37+05:00"}
]
```

### GET /tasks/{id}

Получить одну задачу.

**200** — объект задачи.  
**400** — `id must be a number`.  
**404** — `task not found`.

### POST /tasks

Создать задачу.

Request:

```json
{"title":"Подключить бэк к фронту"}
```

**201** — созданный объект.  
**400** — некорректный JSON.  
**422** — `title is required`.

### PATCH /tasks/{id}

Изменить `title` и/или `done`.

Request:

```json
{"done":true}
```

**200** — обновлённый объект.  
**400** — некорректный id или JSON.  
**404** — задача не найдена.  
**422** — поля для изменения не переданы или title пустой.

### DELETE /tasks/{id}

Удалить задачу.

**204** — без тела.  
**400** — некорректный id.  
**404** — задача не найдена.

## Проверка через curl

```bash
curl -X POST http://localhost:8080/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"Подключить бэк к фронту"}'

curl http://localhost:8080/tasks
curl http://localhost:8080/tasks/1

curl -X PATCH http://localhost:8080/tasks/1 \
  -H 'Content-Type: application/json' \
  -d '{"done":true}'

curl -X DELETE http://localhost:8080/tasks/2
```
