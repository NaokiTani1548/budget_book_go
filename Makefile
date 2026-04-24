DB_URL=postgres://budget_book_user:budget_book_pass@localhost:5432/budget_book_db?sslmode=disable

.PHONY: up down migrate migrate-down seed sqlc run

## Docker
up:
	docker compose up -d

down:
	docker compose down

## マイグレーション
migrate:
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path db/migrations -database "$(DB_URL)" down 1

## シードデータ投入
seed:
	docker exec -i budget-book-db psql -U budget_book_user -d budget_book_db < db/seed.sql

## sqlcコード生成
sqlc:
	sqlc generate

## サーバー起動
run:
	go run cmd/api/main.go