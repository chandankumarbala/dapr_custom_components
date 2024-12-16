# Variables
IMAGE_NAME := echo-app
VERSION := 1.0
DOCKER_REGISTRY := alicejgibbons
FULL_IMAGE_NAME := $(DOCKER_REGISTRY)/$(IMAGE_NAME)

.PHONY: build
build:
	@echo "Building Docker image..."
	docker buildx build --no-cache --platform=linux/amd64,linux/arm64 -t ${DOCKER_REGISTRY}/${IMAGE_NAME}:${VERSION} --push .

.PHONY: run
run:
	@echo "Running service with Dapr..."
	npm install
	dapr run \
		--app-id echo-app \
		--app-port 3000 \
		--dapr-http-port 3500 \
		--log-level debug \
		--config ./components/config.yaml \
		--resources-path ./components \
		-- node echoapp.js

.PHONY: clean
clean:
	@echo "Cleaning node_modules..."
	rm -rf node_modules/