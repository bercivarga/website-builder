.PHONY: up down logs frontend backend db migrate migrate-create migrate-up migrate-down migrate-status migrate-reset test clean

up:
	docker-compose up -d
	docker-compose logs -f frontend backend

down:
	docker-compose down

logs:
	docker-compose logs -f frontend backend

logs-all:
	docker-compose logs -f

frontend:
	cd frontend && npm start

backend:
	cd backend && go run cmd/server/main.go

db:
	docker-compose up -d db

migrate: migrate-up

migrate-create:
	@read -p "Enter migration name: " name; \
	docker-compose run --rm backend goose -dir ./migrations create $$name sql

migrate-up:
	docker-compose run --rm backend goose -dir ./migrations up

migrate-down:
	docker-compose run --rm backend goose -dir ./migrations down

migrate-status:
	docker-compose run --rm backend goose -dir ./migrations status

migrate-reset:
	docker-compose run --rm backend goose -dir ./migrations reset

test:
	cd backend && go test ./...
	cd frontend && npm test

clean:
	docker-compose down -v
	rm -rf backend/tmp
	rm -rf frontend/node_modules