.PHONY: proto sqlc migrate-create build build-server build-cli run test tools

proto:
	protoc --proto_path=proto --go_out=gen --go_opt=paths=source_relative --go-grpc_out=gen --go-grpc_opt=paths=source_relative  notes/v1/notes.proto

sqlc:
	sqlc generate

migrate-create:
	migrate create -ext sql -dir internal/sql/migrations -seq $(name)

build: build-server build-client

build-server:
	go build -o bin/notes-server ./cmd/server

build-client:
	go build -o bin/notes ./cmd/client

run:
	./bin/notes-server

test:
	go test  ./...

tools:
	brew install protobuf golang-migrate
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
