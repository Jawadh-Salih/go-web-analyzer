# Makefile for Go Web Analyzer
APP_NAME := go-web-analyzer
BUILD_DIR := ./build
SRC_DIR := ./cmd
GO_FILES := $(shell find $(SRC_DIR) -type f -name '*.go')

# Define the Docker image name
DOCKER_IMAGE := $(APP_NAME):latest
DOCKER_PORT := 8080


build: $(GO_FILES)
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)

run: build
	@echo "Running the application..."
	@./$(BUILD_DIR)/$(APP_NAME)

test:
	@echo "Running tests..."
	@go test ./... -v

coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out

coverage-html:
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

docker-run: docker-build
	@echo "Running Docker container..."
	@docker run -p $(DOCKER_PORT):$(DOCKER_PORT) $(DOCKER_IMAGE)
