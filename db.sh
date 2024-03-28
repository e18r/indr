#! /bin/bash

psql postgres -f ./palindr.sql
psql palindr -f ./norm.sql
