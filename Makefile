.PHONY: up down server frontend db-console db-migrate test

up:
	docker compose up -d
	@echo "Waiting for MySQL to be ready..."
	@until docker compose exec db mysqladmin ping -uroot -proot --silent 2>/dev/null; do sleep 1; done
	@echo "MySQL is ready"

down:
	docker compose down

server:
	cd backend && INTERNAL_SECRET=$${INTERNAL_SECRET:-dev-secret} go run ./cmd/server

frontend:
	cd frontend && npm run dev

db-console:
	docker compose exec db mysql -uapp -papp nikkei_trade

db-migrate:
	@for f in db/init/*.sql; do \
		echo "Applying $$f ..."; \
		docker compose exec -T db mysql -uapp -papp nikkei_trade < $$f; \
	done
	@echo "Migration done"

test:
	cd backend && go test ./...
