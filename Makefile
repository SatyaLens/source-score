# common env setup
export POSTGRES_VERSION=17

build:
	go mod tidy && \
	go generate ./... && \
	go build && \
	go mod tidy

acceptance-test:
	docker-compose -f acceptance/docker-compose.yaml up

minikube-cleanup:
	minikube stop

minikube-setup: minikube-cleanup
	minikube start --cpus 3 --memory 4096

pg-setup: minikube-setup
	kubectl apply --server-side -f \
		https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/main/releases/cnpg-1.24.0.yaml