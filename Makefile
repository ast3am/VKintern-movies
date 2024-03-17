.PHONY: run test build swag

all: build run

swag:
	swag init -g cmd/main.go

run:

	docker-compose up -d

test:
	docker-compose -f ./test/docker-compose.yml down -v
	docker-compose -f ./test/docker-compose.yml up -d
	go test -count=1 -cover ./...
	docker-compose -f ./test/docker-compose.yml down -v

build:
	docker-compose build
