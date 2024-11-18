CNPG_VERSION ?= "1.24.0"
PG_USER_PASSWORD ?= "test_123"

# common env setup
export PG_USER_PASSWORD

codegen:
	go mod tidy
	mkdir -p pkg/api
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=configs/config.yaml api/source-score.yaml
	go mod tidy

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
	helm upgrade --install cnpg-database --set cnpg_cluster.password=$(PG_USER_PASSWORD) helm/cnpg-database
	@echo -e "\n\e[0;32mCreated CNPG cluster :)\n\e[0m"
	sleep 240
	kubectl get pods -l cnpg.io/cluster=cnpg-cluster -n postgres-cluster

local-pg-setup: minikube-setup pg-setup

cloud-k8s-setup:
	chmod 400 configs/civo-kubeconfig
	cp -f configs/civo-kubeconfig ~/.kube/config

cloud-pg-setup: cloud-k8s-setup pg-setup