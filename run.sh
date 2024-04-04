#! /bin/bash

cd "$(dirname $0)"

export DATABASE_URL="postgresql://indr:indr@localhost:5432/palindr"
export PORT=3000
go run .
