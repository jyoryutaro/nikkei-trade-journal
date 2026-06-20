.PHONY: up down seed server frontend

up:
	docker compose up -d
	@echo "Waiting for MySQL to be ready..."
	@until docker compose exec db mysqladmin ping -uroot -proot --silent 2>/dev/null; do sleep 1; done
	@echo "MySQL is ready"

down:
	docker compose down

seed:
	cd backend && go run ./cmd/seed

server:
	cd backend && go run ./cmd/server

frontend:
	cd frontend && npm run dev
