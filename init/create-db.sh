#!/bin/bash

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
		CREATE USER thepsina WITH PASSWORD 'postgres';
		CREATE DATABASE db;
		GRANT ALL PRIVILEGES ON DATABASE db TO thepsina;
EOSQL

psql -U thepsina -d db -f /init.sql
