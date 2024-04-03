#! /bin/bash

cd "$(dirname $0)"

export DATABASE_URL="postgresql://indr:indr@localhost:5432/palindr"
export PORT=3000
URL=$(./url.sh)
printf $URL > ../pal/indr.url
printf "URL: %s\nsaved in ../pal.indr.url\n" $URL
go run .
