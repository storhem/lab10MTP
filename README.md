# Лабораторная работа №10
## Веб-разработка: FastAPI (Python) vs Gin (Go)

**Студент:** Евланичев Максим Юрьевич  
**Группа:** 221131  
**Вариант:** 6  

---

## Задания

### Средняя сложность
- **М2** — Middleware для логирования в Go (Gin)
- **М6** — Сравнение скорости FastAPI vs Gin под нагрузкой (wrk/ab)
- **М8** — Swagger-документация для FastAPI и OpenAPI для Gin

### Повышенная сложность
- **В2** — API-шлюз на Go: маршрутизация к Python и Go микросервисам
- **В6** — Тесты производительности, сравнение потребления памяти

---

## Описание

Проект демонстрирует создание и сравнение веб-сервисов на двух языках:

- **Go-сервис** (Gin) — REST API с кастомным middleware для логирования запросов: метод, путь, статус и время ответа.
- **Python-сервис** (FastAPI) — будет добавлен в следующих заданиях.

---

## Структура проекта

```
.
├── src/
│   └── go-service/       # Go-сервис (Gin)
│       ├── app/
│       │   └── router.go # Роутер и эндпоинты
│       ├── middleware/
│       │   └── logger.go # Кастомный Logger middleware
│       ├── main.go
│       ├── go.mod
│       └── go.sum
├── tests/
│   └── go-service/       # Тесты Go-сервиса
│       ├── router_test.go
│       └── go.mod
├── .gitignore
├── PROMPT_LOG.md
└── README.md
```

---

## Технологии

- **Go** 1.22 + [Gin](https://gin-gonic.com/)
- **Python** (будет добавлен) + FastAPI

---

## Запуск Go-сервиса

### Требования
- Go 1.22+

### Установка зависимостей и запуск

```bash
cd src/go-service
go mod download
go run .
```

Сервис запустится на `http://localhost:8080`.

### Эндпоинты

| Метод | Путь         | Описание              |
|-------|--------------|-----------------------|
| GET   | /ping        | Проверка работы       |
| GET   | /items       | Список товаров        |
| GET   | /items/:id   | Товар по ID           |

### Примеры запросов

```bash
curl http://localhost:8080/ping
# {"message":"pong"}

curl http://localhost:8080/items
# [{"id":1,"name":"Apple","price":1.5},{"id":2,"name":"Banana","price":0.75}]

curl http://localhost:8080/items/1
# {"id":1,"name":"Apple","price":1.5}

curl http://localhost:8080/items/99
# {"error":"item not found"}
```

### Пример лога middleware

```
2026/04/04 00:20:32 [GET] /ping | status=200 | duration=312µs | ip=127.0.0.1
```

---

## Запуск тестов

```bash
cd tests/go-service
go test ./... -v
```
