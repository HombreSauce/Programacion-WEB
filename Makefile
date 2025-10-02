APP_NAME = turnero

up: 
	docker compose up -d

down: 
	docker compose down

reset: 
	docker compose down -v
	rm -r ./db/sqlc

generate: 
	sqlc generate

run: build
	go run main.go


build:
	docker compose up -d
	@echo "Esperando a que Postgres esté listo..."
	@docker exec -i db_turnero bash -c "until pg_isready -U postgres -d base_turnero; do sleep 1; done"
	sqlc generate