CLUSTER_NAME = "test-env"
CNPG_VERSION ?= "1.24.0"
PG_HOST ?= "http://127.0.0.1"
PG_USER_PASSWORD ?= "test_123"
SERVER_PORT ?= 8070

# common env setup
export PG_SERVER=$(PG_HOST):$(SERVER_PORT)
export PG_USER_PASSWORD
export PORT=$(SERVER_PORT)

codegen:
	go mod tidy
	mkdir -p pkg/api
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=configs/config.yaml api/source-score.yaml
	go mod tidy

lint: codegen
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

build: codegen
	go build

unit-tests:
	go run github.com/onsi/ginkgo/v2/ginkgo run --skip-package=acceptance --cover --coverprofile=coverage.out ./...

acceptance-tests: build
	chmod +x ./source-score
	( \
		./source-score & BG_PID=$$!; \
		trap "echo 'terminating the app'; kill $$BG_PID" EXIT; \
		echo "app running with PID $$BG_PID"; \
		go run github.com/onsi/ginkgo/v2/ginkgo run --cover --coverprofile=coverage.out acceptance/...; \
	)

tests: unit-tests acceptance-tests

start: codegen
	go run main.go

k3d-cleanup:
	@if k3d cluster list | grep -q "^$(CLUSTER_NAME) "; then \
		k3d cluster delete "$(CLUSTER_NAME)"; \
	fi

k3d-setup: k3d-cleanup
	k3d cluster create $(CLUSTER_NAME) --servers 1 --agents 1
	@echo -e "\n\n"

cnpg-controller-setup:
	kubectl apply --server-side -f \
		https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/main/releases/cnpg-$(CNPG_VERSION).yaml
	@echo -e "\n\e[0;32mInstalled CNPG controller on the cluster :)\n\e[0m"
	sleep 60
	kubectl get deployment -n cnpg-system cnpg-controller-manager
	@echo -e "\n\n"

pg-setup: cnpg-controller-setup
	helm upgrade --install cnpg-database --set cnpg_cluster.password=$(PG_USER_PASSWORD) helm/cnpg-database
	@echo -e "\n\e[0;32mCreated CNPG cluster :)\n\e[0m"
	sleep 240
	kubectl get pods -l cnpg.io/cluster=cnpg-cluster -n postgres-cluster

local-pg-setup: k3d-setup pg-setup

cloud-k8s-setup:
	chmod 400 configs/civo-kubeconfig
	cp -f configs/civo-kubeconfig ~/.kube/config

cloud-pg-setup: cloud-k8s-setup pg-setup