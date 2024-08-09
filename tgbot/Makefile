SERVICE_NAME=tgbot

TEST_FLAGS=-count=1
BUILD_FLAGS=

.PHONY: run, build, lint, test, coverage, migrate-new, migrate-up, migrate-down

# TODO define your envs, switch log_level to `debug` during developing

# jwt_secret flag should be identical for /backend and /tgbot

run: ### run app
	go run cmd/main.go -telegram_token 7272065857:AAEPlwtxddJwIwtLYqBhX4OuBndvowJPNV4

build: ### build app
	go build ${BUILD_FLAGS} -o ${SERVICE_NAME} cmd/main.go

lint: ### run linter
	@golangci-lint --timeout=2m run

test: ### run test
	go test ${TEST_FLAGS} ./...

coverage: ### generate coverage report
	go test ${TEST_FLAGS} -coverprofile=coverage.out ./...
	go tool cover -html="coverage.out"
