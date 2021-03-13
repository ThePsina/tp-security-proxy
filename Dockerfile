FROM golang:1.16 AS build

ENV CGO_ENABLED=0

WORKDIR /opt/app
ADD . .
RUN go build ./cmd/main.go

FROM alpine:latest

MAINTAINER thepsina

RUN apt-get -y update && apt-get install -y tzdata

ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER

USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER thepsina WITH SUPERUSER PASSWORD 'postgres';" &&\
    createdb -O thepsina db &&\
    /etc/init.d/postgresql stop

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /usr/src/app

COPY . .
COPY --from=build /opt/app/main .

EXPOSE 8081 8082
ENV PGPASSWORD postgres
CMD service postgresql start &&  psql -h localhost -d db -U thepsina -p 5432 -a -q -f ./init/init.sql && ./main