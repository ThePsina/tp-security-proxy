# Build
FROM golang:1.16-buster AS build

WORKDIR /app
ADD . .

ENV CGO_ENABLED=0

RUN go mod tidy
RUN go mod vendor
RUN go build -o bin/proxy ./cmd/main.go

# Enviroment
FROM alpine:latest

WORKDIR /app
COPY --from=build /app/bin/proxy .
ADD config.yml .
ADD params .

ENTRYPOINT ["/app/proxy"]

EXPOSE 8081 8082
