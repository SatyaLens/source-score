CNPG_VERSION ?= "1.24.0"
PG_USER_PASSWORD ?= "test_123"
POSTGRES_VERSION ?= "17"

# common env setup
export PG_USER_PASSWORD
export POSTGRES_VERSION

codegen:
	go mod tidy && \
	go generate ./...

build: codegen
	go build

minikube-cleanup:
	@if minikube status > /dev/null 2>&1; then \
		minikube stop; \
		minikube delete; \
	fi

minikube-setup: minikube-cleanup
	minikube start --cpus 4 --memory 6144
	@echo -e "\n\n"

cnpg-controller-setup:
	kubectl apply --server-side -f \
		https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/main/releases/cnpg-$(CNPG_VERSION).yaml
	@echo -e "\n\e[0;32mInstalled CNPG controller on the cluster :)\n\e[0m"
	sleep 60
	kubectl get deployment -n cnpg-system cnpg-controller-manager
	@echo -e "\n\n"

pg-setup: cnpg-controller-setup
	helm install --set cnpg_cluster.password=$(PG_USER_PASSWORD) cnpg-database helm/cnpg-database
	@echo -e "\n\e[0;32mCreated CNPG cluster :)\n\e[0m"
	sleep 120
	kubectl get pods -l cnpg.io/cluster=cnpg-cluster -n postgres-cluster

local-pg-setup: minikube-setup pg-setup