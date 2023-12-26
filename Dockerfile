FROM golang:latest

COPY . /home/src
WORKDIR /home/src
ENTRYPOINT [ "go", "run", "main.go" ]