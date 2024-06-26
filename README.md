# GophKeeper
Password keeper - server (Yandex praktikum final project)


## minIO
install
https://min.io/docs/minio/container/index.html




Get and set data : habr example
https://habr.com/ru/companies/ozontech/articles/586024/?code=9040949a58b797539b7c7d5b3a3462e2&state=Zfbrcwb5DSFDKXKFDgu6bkBC&hl=ru


## Mock generate 

```bash
go install github.com/golang/mock/mockgen@v1.6.0
go get github.com/golang/mock/gomock

```

```bash
mockgen -source=internal/services/keeper.go \
    -destination=internal/services/mocks/keeper.gen.go \
    -package=mocks


```


## Generate oapi
Use make or bash command or //TODO build generate
```
make oapi
```
```bash
go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=internal/rest/oapi/cfg.yaml --package oapi internal/rest/oapi/keeper.yaml
```
/home/igor/Desktop/code/GophKeeper/internal/rest/oapi/keeper.yaml
## Mock
```bash
mockgen -source=internal/services/keeper.go \
    -destination=internal/services/mocks/keeper.gen.go \
    -package=mocks
```
-destination mock-interfaces.go github.com/yourhandle/worker Doer
```

## ENV variables


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

## Git - remove file with git ignore, when it already added.
git rm --cached internal/api/oapi/gokeeper.gen.go
//git reset internal/api/oapi/gokeeper.gen.go

## Curl tests


## Private key for JWT token
PrivateKey is an ECDSA private key which was generated with the following
command:

```bash
openssl ecparam -name prime256v1 -genkey -noout -out cert/jwtpkey.pem
```
We are using a hard coded key here in this example, but in real applications,
you would never do this. Your JWT signing key must never be in your application,
only the public key.

## Create certificates TLS
Generate private key:
```bash
openssl genrsa -out pkey.pem 2048
```
Generate CSR: (In the "Common Name" set the domain of your service provider app)
```bash
openssl req -new -key pkey.pem -out server.csr
```
Generate Self Signed Cert
```bash
openssl x509 -req -days 365 -in server.csr -signkey pkey.pem -out pkey.crt
rm pkey.pem
```