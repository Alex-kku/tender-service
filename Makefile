include .env
export

build:
	docker-compose build tender
.PHONY: build

run:
	docker-compose up  tender
.PHONY: run

migrate-up:
	migrate -path migrations -database '$(POSTGRES_CONN)?sslmode=disable' up
.PHONY: migrate-up

migrate-down:
	echo "y" | migrate -path migrations -database '$(POSTGRES_CONN)?sslmode=disable' down
.PHONY: migrate-down


