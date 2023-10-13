BIN := "./bin"
DOCKER_IMG="image-preview-service:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

# GENERALLY

protoc:
	protoc --go-grpc_out=./core/storage_service/rpc/api --go_out=./core/storage_service/rpc/api ./core/storage_service/rpc/protofiles/storage.proto

install-lint-deps: tidy
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.53.3

lint: install-lint-deps
	golangci-lint run --out-format=github-actions ./core/...
	golangci-lint run --out-format=github-actions ./cmd/...
	golangci-lint run --out-format=github-actions ./docker.integration-test/...

build:
	go build -v -o $(BIN)/storage_service -ldflags "$(LDFLAGS)" ./cmd/storage_service
	go build -v -o $(BIN)/http_service -ldflags "$(LDFLAGS)" ./cmd/http_service

version: build
	$(BIN)/storage_service version
	$(BIN)/http_service version

up-storage_service: build
	$(BIN)/storage_service -config ./config/storage_service.json

up-http_service: build
	$(BIN)/http_service -config ./config/http_service.json

tidy:
	go mod tidy

# NATIVE RUN BY CODE

run-native-storage_service:
	go run ./cmd/storage_service/ -config ./config/storage_service.json

run-native-http_service:
	go run ./cmd/http_service/ -config ./config/http_service.json

run-int-test-storage_service:
	go run ./cmd/storage_service/ -config ./docker.integration-test/config/storage_service.json

run-int-test-http_service:
	go run ./cmd/http_service/ -config ./docker.integration-test/config/http_service.json

run-int-test-nginx_service-drop:
	docker stop previewer-nginx 2> /dev/null || true
	docker rm -f previewer-nginx 2> /dev/null || true

run-int-test-jpeg-generation:
	go test ./core/app/functions/ -run TestGenerateImage

run-int-test-nginx_service: run-int-test-nginx_service-drop run-int-test-jpeg-generation
	docker build -f ./docker.integration-test/Dockerfile.nginx_service -t int-test/previewer-nginx_service .
	docker run --name previewer-nginx -d -p 8082:80 int-test/previewer-nginx_service

run-int-test-exthttp_service:
	go test -v ./docker.integration-test/exthttp_service/ -run TestIntegration -config ../config/exthttp_service.json

run-int-test-fullstack_client:
	go test -v ./docker.integration-test/fullstack_client/ -run TestIntegration -config ../config/fullstack_client.json

# UNIT TESTS

testcache-clean:
	go clean -testcache

test-storage_service: testcache-clean
	go test -race ./core/storage_service/test/ -v

show-storage_service-coverpkg: testcache-clean
	go test ./core/storage_service/test/ -coverpkg=./.../server
	go test ./core/storage_service/test/ -coverpkg=./.../client
	go test ./core/storage_service/test/ -coverpkg=./.../rpc/api  
	go test ./core/storage_service/test/ -coverpkg=./.../common 
	go test ./core/storage_service/test/ -coverpkg=./.../models

test-app-functions-with-foreign-image-server: testcache-clean
	go test -race ./core/app/functions/ -v -run TestForeign

test-app-functions-with-internal-image-server: testcache-clean
	go test -race ./core/app/functions/ -v -run TestInternal

show-app-functions-cover: testcache-clean
	go test ./core/app/functions/ -cover

test-http_service: testcache-clean
	go test -race ./core/http_service/ -v

show-http_service-cover: 
	go test ./core/http_service/ -cover

show-http_service-coverpkg: testcache-clean
	go test ./core/http_service/ -coverpkg=./...

test-config: testcache-clean
	go test -race ./core/config/ -v
	go test -race ./docker.integration-test/exthttp_service/ -v -run TestExtHTTPServiceConfig
	go test -race ./docker.integration-test/fullstack_client/ -v -run TestFullstackClientConfig

show-config-cover: testcache-clean
	go test ./core/config/ -cover

test: \
	test-config \
	test-http_service \
	test-storage_service \
	test-app-functions-with-internal-image-server \
	test-app-functions-with-foreign-image-server

show-cover: \
	show-app-functions-cover \
	show-storage_service-coverpkg \
	show-http_service-cover \
	show-http_service-coverpkg \
	show-config-cover
	
# INTEGRATION TEST

integration-test-configure: run-int-test-jpeg-generation
	docker network create --subnet 10.10.0.0/24 int-test-previewer-network 2> /dev/null || true

integration-test-clear:
	docker rmi -f int-test-previewer-storage_service 2> /dev/null || true
	docker rmi -f int-test-previewer-http_service 2> /dev/null || true
	docker rmi -f int-test-previewer-nginx_service 2> /dev/null || true
	docker rmi -f int-test-previewer-exthttp_service 2> /dev/null || true
	docker rmi -f int-test-previewer-fullstack_client 2> /dev/null || true
	docker rm -f int-test-previewer-storage_service 2> /dev/null || true
	docker rm -f int-test-previewer-http_service 2> /dev/null || true
	docker rm -f int-test-previewer-nginx_service 2> /dev/null || true
	docker rm -f int-test-previewer-exthttp_service 2> /dev/null || true
	docker rm -f int-test-previewer-fullstack_client 2> /dev/null || true
	docker network rm int-test-previewer-network 2> /dev/null || true
	rm -f ./docker.integration-test/images/transformed.*

integration-test-build-d: \
	integration-test-clear \
	integration-test-configure
	BUILDKIT_PROGRESS=plain docker-compose -f ./docker.integration-test/docker-compose.yaml up -d --build --remove-orphans

integration-test-fullstack_client-run: run-int-test-fullstack_client

integration-test-d-up-infrastructure: \
	integration-test-clear \
	integration-test-configure
	BUILDKIT_PROGRESS=plain docker-compose -f ./docker.integration-test/docker-compose.yaml up --build -d --remove-orphans

integration-test-d-up-fullstack_client:
	BUILDKIT_PROGRESS=plain docker-compose -f ./docker.integration-test/docker-compose.fullstack_client.yaml up --build

integration-test-fullstack_client-exit-code:
	BUILDKIT_PROGRESS=plain docker-compose -f ./docker.integration-test/docker-compose.fullstack_client.yaml up --build --abort-on-container-exit --exit-code-from fullstackclient

integration-test-build:
	BUILDKIT_PROGRESS=plain docker-compose -f ./docker.integration-test/docker-compose.yaml -f ./docker.integration-test/docker-compose.fullstack_client.depended.yaml --build --remove-orphans

integration-test-exit-code:
	BUILDKIT_PROGRESS=plain docker-compose -f ./docker.integration-test/docker-compose.yaml -f ./docker.integration-test/docker-compose.fullstack_client.depended.yaml up --build --remove-orphans --exit-code-from fullstackclient --abort-on-container-exit

integration-test: \
	integration-test-clear \
	integration-test-configure \
	integration-test-exit-code

integration-test-autoclean: integration-test
	make integration-test-clear

# DOCKER COMPOSE

docker-compose:
	BUILDKIT_PROGRESS=plain docker-compose -f ./docker/docker-compose.yaml up -d --build --remove-orphans

# .PHONY

.PHONY: test integration-test lint docker-compose
