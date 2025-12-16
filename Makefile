.PHONY: generate-seeds db-down db-seed db-up db-restart

generate-seeds:
	go run ./seeds/generate_seeds.go --players $${players-1000} --teams $${teams-50} --tournaments $${tournaments-25} --out $${out-seeds/generated_seed.sql}


db-seed: generate-seeds db-down
	docker compose up --build -d db
	docker compose exec -T db psql -U postgres -d cyber_tournament -f /seeds/generated_seed.sql

build:
	docker compose build

restart:
	docker compose down
	docker compose up --build

restard-d:
	docker compose down
	docker compose up --build -d

up:
	docker compose up

up-d:
	docker compose up -d

down:
	docker compose down 

down-v:
	docker compose down -v

swag:
	swag init -g cmd/api/main.go -o docs