#!/bin/sh
# Run as postgres
set -x

die () {
    echo >&2 "$@"
    exit 1
}

if [ "$#" -ne 4 ]; then
    die "Usage: $0 <dbuser> <dbpassword> <dbname> <dbversion>"
fi

dbusr=$1
dbpwd=$2
dbname=$3
dbver=$4
/usr/pgsql-$dbver/bin/pg_ctl start -D /var/lib/pgsql/$dbver/data


psql <<END
CREATE USER $dbusr WITH PASSWORD '$dbpwd' CREATEROLE CREATEDB REPLICATION BYPASSRLS;
CREATE DATABASE $dbname OWNER=$dbusr
END

if ! psql postgres://$dbusr:$dbpwd@localhost/$dbname -c '\l+'; then
    echo "Failed to list database\n"
    exit -1
fi

if ! psql postgres://$dbusr:$dbpwd@localhost/$dbname -c 'SELECT 1'; then
    echo "Failed to query database\n"
    exit -1
fi

if ! psql postgres://$dbusr:$dbpwd@localhost/$dbname -c 'DROP SCHEMA IF EXISTS iex CASCADE'; then
    echo "Failed to drop schema\n"
    exit -1
fi

if ! DB_USER=$dbusr envsubst < iex.yaml.template > iex.yaml; then
    echo "Failed to generate schemal file from template\n"
    exit -1
fi

if ! yamltodb $dbname iex.yaml > iex.sql; then
    echo "Failed to create schema\n"
    exit -1
fi

if ! psql postgres://$dbusr:$dbpwd@localhost/$dbname -f iex.sql; then
    echo "Failed to create schema\n"
    exit -1
fi

/usr/pgsql-$dbver/bin/pg_ctl stop -D /var/lib/pgsql/$dbver/data
if [ $? -ne 0 ]
then
   echo "Failed to stop database\n"
   exit -1;
fi