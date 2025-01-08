SERVICE_NAME=sos-kotopes

TEST_FLAGS=-count=1
BUILD_FLAGS=

.PHONY: run, build, lint, test, coverage, migrate-new, migrate-up, migrate-down

# TODO define your envs, switch log_level to `debug` during developing
PG_URL=postgres://postgres:postgres@localhost:5432/kotopes

run: ### run app
	go run cmd/main.go -db_url '$(PG_URL)' -port localhost:8080 -log_level debug -vk_client_id 52010687 -vk_client_secret M6SWkM8KmIIJA60hTVx1

build: ### build app
	go build ${BUILD_FLAGS} -o ${SERVICE_NAME} cmd/main.go

lint: ### run linter
	@golangci-lint --timeout=2m run

test: ### run test
	go test ${TEST_FLAGS} ./...

coverage: ### generate coverage report
	go test ${TEST_FLAGS} -coverprofile=coverage.out ./...
	go tool cover -html="coverage.out"

MIGRATION_NAME=vets

migrate-new: ### create a new migration
	migrate create -ext sql -dir ./internal/data -seq ${MIGRATION_NAME}

migrate-up: ### apply all migrations
	migrate -path ./internal/data -database '$(PG_URL)?sslmode=disable' up

migrate-down: ### migration down
	migrate -path ./internal/data -database '$(PG_URL)?sslmode=disable' down

swagger: ### generate swagger docs
	swag fmt
	swag init -g cmd/main.go

mock: ### generate mocks
	mockery
