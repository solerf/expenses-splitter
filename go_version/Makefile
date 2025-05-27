GO_OS = $(shell uname -o | sed 's/gnu\///i' | tr '[:upper:]' '[:lower:]')
GO_ARCH = $(shell uname -m | sed 's/aarch64/arm64/i' | sed 's/x86_64/amd64/i')

DOCKER_TAG = bill-splitter
HOST_PORT = 8000

build:
	env GOOS=$(GO_OS) GOARCH=$(GO_ARCH) CGO_ENABLED=0 go build -o ./bin/bill-splitter .

build-docker:
	docker build -t $(DOCKER_TAG) .

tests:
	go test ./...

run:
	go run .

run-docker:
	docker run --rm -p $(HOST_PORT):8000 $(DOCKER_TAG)
