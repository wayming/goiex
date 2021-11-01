#!/bin/bash
# Run as postgres
set -x
dbusr=test
dbpwd=test
dbname=testdb

/usr/pgsql-13/bin/pg_ctl start -D /var/lib/pgsql/13/data

psql <<END
CREATE USER $dbusr WITH PASSWORD '$dbpwd' CREATEROLE CREATEDB REPLICATION BYPASSRLS;
CREATE DATABASE $dbname OWNER=$dbusr
END

psql postgres://$dbusr:$dbpwd@localhost/testdb -c '\l+'
psql postgres://$dbusr:$dbpwd@localhost/testdb -c 'SELECT 1'

/usr/pgsql-13/bin/pg_ctl stop -D /var/lib/pgsql/13/data