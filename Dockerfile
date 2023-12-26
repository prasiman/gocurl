FROM golang:latest

WORKDIR /app
COPY . .
COPY go.mod .

ENTRYPOINT [ "go", "run", "/app/main.go" ]