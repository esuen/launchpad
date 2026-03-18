APP_NAME := launchpad
BIN_DIR := bin

DOCKER_IMAGE := ghcr.io/esuen/launchpad
DOCKER_TAG := latest

HELM_CHART := deploy/helm/launchpad
HELM_RELEASE := launchpad

.PHONY: build run test lint clean docker-build docker-run helm-lint helm-template helm-install helm-uninstall grafana prometheus

build:
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./... -v -race -count=1

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BIN_DIR)

docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	docker run --rm -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

helm-lint:
	helm lint $(HELM_CHART)

helm-template:
	helm template $(HELM_RELEASE) $(HELM_CHART)

helm-install:
	helm upgrade --install $(HELM_RELEASE) $(HELM_CHART)

helm-uninstall:
	helm uninstall $(HELM_RELEASE)

grafana:
	@echo "Grafana: http://localhost:3000 (admin/admin)"
	@echo "Dashboard: http://localhost:3000/d/launchpad/launchpad"
	kubectl port-forward svc/grafana 3000:80

prometheus:
	@echo "Prometheus: http://localhost:9090"
	kubectl port-forward svc/prometheus-server 9090:80
