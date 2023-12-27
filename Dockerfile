FROM golang:latest AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o /app

FROM gcr.io/distroless/static-debian11
COPY --from=build /app /app
ENTRYPOINT [ "/app" ]