APP_USER_PASSWORD ?= "sourcescore"
SERVER_HOST ?= "localhost"
SUPER_USER_PASSWORD ?= "test_123"
TEST_CLUSTER_NAME = "test-env"

# common env setup
export APP_USER_PASSWORD
export PG_HOST=$(SERVER_HOST)
export SUPER_USER_PASSWORD

codegen:
	go mod tidy
	go generate ./...
	mkdir -p pkg/api
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=configs/config.yaml api/source-score.yaml
	go mod tidy

lint: codegen
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

build: codegen
	go build -o ./source-score ./cmd/app
	chmod +x ./source-score

unit-tests:
	go run github.com/onsi/ginkgo/v2/ginkgo run --skip-package=acceptance --cover --coverprofile=coverage.out ./...

start-containers:
	docker compose -f acceptance/compose.yaml up -d

acceptance-tests: start-containers
	sleep 20 && cd acceptance && go run github.com/onsi/ginkgo/v2/ginkgo -r ./... && cd -

tests: unit-tests acceptance-tests

cleanup-containers:
	docker compose -f acceptance/compose.yaml down -v
	docker rmi acceptance-app:latest