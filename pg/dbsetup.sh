#!/bin/sh
# Run as postgres
set -x

die () {
    echo >&2 "$@"
    exit 1
}

if [ "$#" -ne 3 ]; then
    die "Usage: $0 <dbuser> <dbpassword> <dbname>"
fi

dbusr=$1
dbpwd=$2
dbname=$3

/usr/pgsql-13/bin/pg_ctl start -D /var/lib/pgsql/13/data

psql <<END
CREATE USER $dbusr WITH PASSWORD '$dbpwd' CREATEROLE CREATEDB REPLICATION BYPASSRLS;
CREATE DATABASE $dbname OWNER=$dbusr
END

psql postgres://$dbusr:$dbpwd@localhost/testdb -c '\l+'
psql postgres://$dbusr:$dbpwd@localhost/testdb -c 'SELECT 1'

/usr/pgsql-13/bin/pg_ctl stop -D /var/lib/pgsql/13/data