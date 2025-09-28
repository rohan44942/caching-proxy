APP_NAME=caching-proxy

.PHONY: build run test docker

build:
	go build -o $(APP_NAME) ./cmd/caching-proxy

run:
	go run ./cmd/caching-proxy --port 3000 --origin http://dummyjson.com

test:
	go test ./...

docker:
	docker build -t $(APP_NAME):latest .