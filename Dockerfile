FROM golang:1.19.3-alpine

WORKDIR /usr/src/app

COPY . .

RUN go mod download
RUN go build -o ./cmd/app/main ./cmd/app/main.go

CMD ["./cmd/app/main"]