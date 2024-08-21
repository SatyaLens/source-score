export GOBIN := $(shell pwd)/bin
export PATH := "$$GOBIN":$(PATH)

install-tools:
	@echo $(value PATH) && \
	echo Installing tools from tools/tools.go && \
	cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go get % && \
	cat ./tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install % && \
	go mod tidy

build: install-tools
	go generate ./... && \
	go build