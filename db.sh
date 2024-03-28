#! /bin/bash

cd "$(dirname $0)"
psql postgres -f ./palindr.sql

psql palindr -f ./norm.sql
psql palindr -f ./text.sql

psql palindr -f ./priv.sql
