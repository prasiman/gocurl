FROM golang:latest

COPY . /home/src
WORKDIR /home/src
ENTRYPOINT [ "go", "run", "/home/src/main.go" ]