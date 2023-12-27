FROM golang:latest

WORKDIR /app
COPY . .

ENTRYPOINT [ "go", "run", "/app/main.go" ]