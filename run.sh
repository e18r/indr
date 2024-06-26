#! /bin/bash

cd "$(dirname $0)"

ENV=$(cat ./ENV)
if [ "$ENV" != "dev" ]; then
   echo "only run in dev env"
   exit 1
fi

export DATABASE_URL="postgresql://indr:indr@localhost:5433/palindr"
export PORT=3000
go run .
