run-servers:
	go run cmd/servers/main.go

run-central-auth:
	go run cmd/central-auth/main.go

run-worker:
	go run cmd/background-tasks/main.go

create-migration:
	migrate create -ext sql -dir migrations/migrations -seq $(name)

docker/up:
	docker compose up --build

docker/down:
	docker compose down

EXCLUDE_DIRS="(proto|vendor|doc|mock|tool|scripts|model|gen|test|filter)"

test:
	./scripts/go-test.sh