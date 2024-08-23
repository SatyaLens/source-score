build:
	go mod tidy && \
	go generate ./... && \
	go build && \
	go mod tidy