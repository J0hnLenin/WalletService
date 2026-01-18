# WalletService

Микросервис для управления балансами кошельков с поддержкой горизонтального шардирования PostgreSQL. Сервис предоставляет REST и gRPC API для операций с кошельками.

![Go Version](https://img.shields.io/badge/Go-1.25.6-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Latest-blue)
![Docker](https://img.shields.io/badge/Docker-Compose-blue)
![Coverage Status](https://coveralls.io/repos/github/J0hnLenin/WalletService/badge.svg?branch=main)
![Tests Status](https://github.com/J0hnLenin/WalletService/actions/workflows/coverage.yml/badge.svg)

## Cтек

- **Backend**: Golang 1.25.6
- **Database**: PostgreSQL
- **API**: gRPC + REST (HTTP Gateway)
- **Контейнеризация**: Docker Compose
- **Тестирование**: Testify
- **Документация API**: Swagger/OpenAPI

## Структура репозитория

```
WalletService/
├── api/                   # API файлы
├── cmd/app/               # Точка входа приложения
├── config/                # Конфигурация приложения
├── internal/              # Внутренние пакеты
│   ├── services/          # Бизнес-логика
│   ├── storage/           # Шардированное хранилище
│   │   └── pgstorage/
│   │       └── migrations/# Миграции для шардов
│   └── pb/                # Protobuf модели и API
├── scripts/               # Вспомогательные скрипты
└── .github/workflows/     # Автозапуск тестов
```

## База данных

### Концепция

Используем шардирование по 64 бакетам с распределением по 32 бакета на 2 инстанса Postgres:

- **Шард 0**: бакеты 0-31
- **Шард 1**: бакеты 32-63

Алгоритм распределения по бакетам построен на основе хэша от uuid кошелька

### Структура базы данных

Каждый бакет хранится в отдельной схеме PostgreSQL:

```sql
-- Пример для bucket_0
CREATE SCHEMA bucket_0;
CREATE TABLE bucket_0.wallets (
    id UUID PRIMARY KEY,
    balance BIGINT NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
```

## Запуск

1. Клонировать репозиторий:

   ```bash
   git clone https://github.com/J0hnLenin/WalletService
   cd ./WalletService
   ```

2. Собрать и запустить контейнеры:

   ```bash
   docker compose up -d
   ```

3. Проверить статус:
   ```bash
   docker compose ps
   ```

Все настройки содержатся в `config.env`

## API

### REST API (HTTP Gateway)

- **Порт**: `8080`

| Метод    | Endpoint                     | Описание                                  |
| -------- | ---------------------------- | ----------------------------------------- |
| **POST** | `/api/v1/wallet`             | Применение операции (пополнение/списание) |
| **GET**  | `/api/v1/wallets/{walletId}` | Получение баланса кошелька                |

### gRPC API

- **Порт**: `50001`
- **Сервис**: `WalletsService`
- **Методы**:
  - `ApplyOperation(Operation) returns (Wallet)`
  - `GetBalance(GetBalanceRequest) returns (Wallet)`

### Примеры запросов

#### 1. Получение баланса

```bash
curl -X GET "http://localhost:8080/api/v1/wallets/550e8400-e29b-41d4-a716-446655440000"
```

**Ответ:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "balance": "1000"
}
```

#### 2. Пополнение баланса

```bash
curl -X POST "http://localhost:8080/api/v1/wallet" \
  -H "Content-Type: application/json" \
  -d '{
    "walletId": "550e8400-e29b-41d4-a716-446655440000",
    "operationType": "DEPOSIT",
    "amount": "500"
  }'
```

#### 3. Списание средств

```bash
curl -X POST "http://localhost:8080/api/v1/wallet" \
  -H "Content-Type: application/json" \
  -d '{
    "walletId": "550e8400-e29b-41d4-a716-446655440000",
    "operationType": "WITHDRAW",
    "amount": "200"
  }'
```

### Модели данных

#### Wallet

```json
{
  "id": "string (UUID)",
  "balance": "string (int64)"
}
```

#### Operation

```json
{
  "walletId": "string (UUID)",
  "operationType": "DEPOSIT | WITHDRAW",
  "amount": "string (int64)"
}
```

## Тестирование

### Запуск всех тестов

```bash
make run-tests
```

### Генерация моков

```bash
make gen-moks
```

### Генерация API кода

```bash
make gen-api
```