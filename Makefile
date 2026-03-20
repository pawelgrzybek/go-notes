.PHONY: build build-server build-cli run test

build: build-server build-cli

build-server:
	go build -o bin/notes-server ./cmd/server

build-cli:
	go build -o bin/notes ./cmd/cli

run:
	./bin/notes-server

test:
	go test  ./...
