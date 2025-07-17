# Project variables
APP_NAME := leader-elector
IMAGE_NAME := $(APP_NAME)
VERSION := latest

# Go-related
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GO_BUILD_FLAGS := -ldflags="-s -w" -trimpath
GOOS ?= linux
GOARCH ?= amd64

.PHONY: all build clean test docker docker-push deploy lint fmt

all: build

## Build the binary
build:
	@echo "🔧 Building $(APP_NAME)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build $(GO_BUILD_FLAGS) -o bin/$(APP_NAME) ./cmd/main.go

## Run tests
test:
	@echo "🧪 Running tests..."
	go test ./... -v

## Format Go code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

## Lint Go code (requires golangci-lint)
lint:
	@echo "🔍 Linting code..."
	golangci-lint run

## Build Docker image
docker:
	@echo "🐳 Building Docker image..."
	docker build -t $(IMAGE_NAME):$(VERSION) .

## Push Docker image
docker-push:
	@echo "📦 Pushing Docker image..."
	docker push $(IMAGE_NAME):$(VERSION)

## Deploy to Kubernetes
deploy:
	@echo "🚀 Deploying to Kubernetes..."
	kubectl apply -f deploy/rbac.yaml

## Clean built files
clean:
	@echo "🧹 Cleaning..."
	rm -rf bin

