#!/bin/bash

# Create databases
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE DATABASE userdb;
	CREATE DATABASE orderdb;
EOSQL
