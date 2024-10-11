# common env setup
export POSTGRES_VERSION=17

build:
	go mod tidy && \
	go generate ./... && \
	go build && \
	go mod tidy

acceptance-test:
	docker-compose -f acceptance/docker-compose.yaml up