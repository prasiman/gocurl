FROM golang:latest AS build
COPY . /home/src
WORKDIR /home/src
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o /bin/action ./

FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=build /bin/action /app
USER nonroot:nonroot
ENTRYPOINT [ "/app" ]