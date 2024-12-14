#! /bin/bash

cd "$(dirname $0)"

ENV=$(cat ./ENV)
if [ "$ENV" != "dev" ]; then
   echo "only run in dev env"
   exit 1
fi

export DATABASE_URL="postgresql://indr:indr@localhost:5432/palindr"
# export DATABASE_URL_2="postgresql://indr:indr@localhost:5432/palindr2"
export DATABASE_URL_2=""
export PORT=3000
go run .
