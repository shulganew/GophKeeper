# GophKeeper
Password keeper - server (Yandex praktikum final project)

## Mock
```bash
mockgen -source=internal/services/user.go \
    -destination=internal/services/mocks/user.gen.go \
    -package=mocks
```

## Переменные окружения

Переменная окружения DSN:
```bash
export DATABASE_URI=postgresql://keeper:1@localhost/keeper?sslmode=disable
```
Переменная RUN_ADDRESS - адреc и порт сервиса Gophermart:
```bash
export RUN_ADDRESS=localhost:8088
```

## Запуск Postgres в контейнере

Для запуска и остановки Postgres в контейнере выполнятьются скрипты создания и миграции базы в make-файле:
* Инициализация
```bash
make pg
```
* Миграция goose
```bash
https://github.com/pressly/goose
go install github.com/pressly/goose/v3/cmd/goose@latest
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="postgresql://keeper:1@postgres:5438/keeper"
migrations up in make file during database startup
```
* Остановка и удаление контейнера
```bash
make pg-stop
```