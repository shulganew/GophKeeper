# MinIO
.PHONY: minio

minio: 
	docker run --rm \
		-d \
   		-p 9000:9000 \
   		-p 9001:9001 \
   		--name minio_v1 \
   		-v ~/minio/data:/data \
   		-e "MINIO_ROOT_USER=admin" \
   		-e "MINIO_ROOT_PASSWORD=12345678" \
   		quay.io/minio/minio server /data --console-address ":9001"

.PHONY: minio-init
minio-init: 
	go install github.com/minio/mc@latest
	mc alias set minio http://localhost:9000 admin 12345678
	mc mb minio/gohpkeeper
	mc ls minio > [0B] gohpkeeper


.PHONY: minio-stop
minio-stop:
	docker stop minio_v1

# oapi generate files
.PHONY: oapi

oapi: 
	go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=internal/app/config/oapi.yaml --package oapi api/api.yaml

#Migrations

.PHONY: pg

pg: 
	docker run --rm \
		--name=keeper_v1 \
		-v $(abspath ./docker/init/):/docker-entrypoint-initdb.d \
		-e POSTGRES_PASSWORD="postgres" \
		-d \
		-p 5438:5432 \
		postgres:15.3
	sleep 5
	
	goose -dir ./migrations  up

.PHONY: pg-stop
pg-stop:
	docker stop keeper_v1


# Linter 

GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache

.PHONY: golangci-lint-run
golangci-lint-run: _golangci-lint-run

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-lint-run: _golangci-lint-reports-mkdir
	-docker run --rm \
    -v $(shell pwd):/app \
    -v $(GOLANGCI_LINT_CACHE):/root/.cache \
    -w /app \
    golangci/golangci-lint:v1.58.1 \
        golangci-lint run \
            -c .golangci.yml \
	| tee ./golangci-lint/report.json


.PHONY: golangci-lint-clean
golangci-lint-clean:
	sudo rm -rf ./golangci-lint 