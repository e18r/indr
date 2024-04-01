#! /bin/bash

cd "$(dirname $0)"

ENV=$(cat ./env)
DATABASE_URL="$(pass palindr/$ENV/db)"
cat ./app.yaml.template | sed "s|\$DATABASE_URL|$DATABASE_URL|" > app.yaml
# gcloud app deploy
# shred -uzn99 ./app.yaml
