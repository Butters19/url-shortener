# URL Shortener

Сервис для сокращения ссылок на Go.

## Возможности

- Генерация уникального 10-символьного кода (a-z, A-Z, 0-9, _)
- Два хранилища: PostgreSQL и in-memory
- Потокобезопасная работа при высокой нагрузке

## Запуск

### С PostgreSQL (через Docker):
```bash
docker compose up --build
```

### Только in-memory (без Docker):
```bash
go run ./cmd/server -storage=memory
```

## API

### POST / — создать короткую ссылку

Request:
```json
{"url": "https://example.com"}
```

Response `201 Created`:
```json
{"short_code": "aB3_xKqZ1m"}
```

### GET /{code} — получить оригинальный URL

Response `200 OK`:
```json
{"url": "https://example.com"}
```

Response `404 Not Found`:
```json
{"error": "url not found"}
```

## Тесты
```bash
go test ./...
```

## Флаги запуска

| Флаг | По умолчанию | Описание |
|------|-------------|----------|
| `-storage` | `memory` | Тип хранилища: `memory` или `postgres` |
| `-dsn` | — | PostgreSQL DSN (обязателен для postgres) |
| `-addr` | `:8080` | Адрес HTTP сервера |