# Multi-stage layout
FROM golang:1.14

#ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .


RUN GOOS=linux GOARCH=amd64 go build ./cmd/main.go

WORKDIR /app
ENTRYPOINT ["./main"]
