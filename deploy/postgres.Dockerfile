FROM postgres:alpine

ADD init/init.sql /
ADD init/create-db.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/create-db.sh

EXPOSE 5432
