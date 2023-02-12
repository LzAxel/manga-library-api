.SILENT:

download:
	go mod download

build:
	go build -o ./cmd/app/main ./cmd/app/main.go

run: build
	./cmd/app/main

swag: 
	./tools/swag init -g ./cmd/app/main.go -o ./docs

test:
	./tools/gotest -v ./...

gen-mock:
	./tools/mockgen -source=internal/storage/storage.go \
	-destination=internal/storage/mocks/mock_storage.go