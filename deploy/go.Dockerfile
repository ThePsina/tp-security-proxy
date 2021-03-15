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

RUN apk upgrade --update-cache --available && \
    apk add openssl && \
    rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=build /app/bin/proxy .
COPY /genCerts/cacert.pem /etc/ssl/certs/
ADD . .

ENTRYPOINT ["/app/proxy"]

EXPOSE 8081 8082
