GO_OS = $(shell uname -o | sed 's/gnu\///i' | tr '[:upper:]' '[:lower:]')
GO_ARCH = $(shell uname -m | sed 's/aarch64/arm64/i' | sed 's/x86_64/amd64/i')

DOCKER_TAG = bill-splitter
HOST_PORT = 8000

build:
	./mvnw clean install

build-docker:
	docker build -t $(DOCKER_TAG) .

tests:
	./mvnw test

run:
	./mvnw -DskipTests clean install
	java -jar target/billsplitter-1.0-SNAPSHOT.jar server config.yml

run-docker:
	docker run --rm -p $(HOST_PORT):8000 $(DOCKER_TAG)
