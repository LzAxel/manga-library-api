.SILENT:

download:
	go mod download

build:
	go build -o ./cmd/app/main ./cmd/app/main.go

run: build
	./cmd/app/main

swag: 
	swag init -g ./cmd/app/main.go -o ./docs