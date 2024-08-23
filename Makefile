export GOBIN := $(shell pwd)/bin
export GOPATH := $(shell pwd)/bin
export PATH := "$$GOBIN":$(PATH)

install-tools:
	@echo Installing tools from tools/tools.go && \
	cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go get % && \
	cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install % && \
	go mod tidy

build: install-tools
	@echo PATH is $(value PATH) && \
	OPATH := $(shell pwd)/bin go generate ./... && \
	go build