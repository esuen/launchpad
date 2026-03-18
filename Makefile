APP_NAME := launchpad
BIN_DIR := bin

DOCKER_IMAGE := ghcr.io/esuen/launchpad
DOCKER_TAG := latest

HELM_CHART := deploy/helm/launchpad
HELM_RELEASE := launchpad

.PHONY: build run test lint clean docker-build docker-run helm-lint helm-template helm-install helm-uninstall

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
