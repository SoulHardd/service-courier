# Service Courier

Учебный микросервис для управления курьерами и доставками, реализованный в рамках backend-курса (Go).
Проект демонстрирует работу с чистой архитектурой, PostgreSQL, асинхронными событиями, метриками и CI.

## Возможности

- CRUD API для курьеров
- Назначение и снятие курьеров с заказов
- Автоматический расчёт времени доставки в зависимости от типа транспорта
- Фоновая обработка:
    - автоматическое освобождение курьеров по дедлайну с `time.Ticker`
- Интеграции:
    - gRPC gateway к service-order
    - Kafka consumer для обработки событий заказов
- Rate limiting и retry при обращении к внешним сервисам
- Метрики Prometheus + Grafana
- Логирование (zap)
- Профилирование через pprof
- Unit и интеграционные тесты
- CI (golangci-lint + go test)


## Архитектура

Проект следует принципам Clean Architecture:

cmd/ – точки входа
internal/
handlers/ – HTTP / queue handlers
useCase/ – бизнес-логика
repository/ – работа с БД
gateway/ – интеграции с внешними сервисами
domain/ – бизнес-сущности
factory/ – фабрики
middleware/ – логирование, метрики, rate limiter
worker/ – фоновые воркеры
pkg/ – переиспользуемые компоненты

Зависимости внедряются через конструкторы, слои общаются через интерфейсы.

## Технологии

- Go 1.21
- PostgreSQL
- goose (миграции)
- Kafka (sarama)
- Prometheus / Grafana
- Docker / docker-compose
- golangci-lint
- zap
- pprof