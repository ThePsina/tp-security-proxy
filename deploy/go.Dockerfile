# Build
FROM golang:1.16-buster AS build

WORKDIR /app
ADD . .

ENV CGO_ENABLED=0

RUN go build -o bin/proxy ./cmd/main.go

# Enviroment
FROM alpine:latest

WORKDIR /app
COPY --from=build /app/bin/proxy .
ADD config.yml .

ENTRYPOINT ["/app/proxy"]

EXPOSE 8081 8081
