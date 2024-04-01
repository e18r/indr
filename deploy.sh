#! /bin/bash

cd "$(dirname $0)"

ENV=$(cat ./environment)
PROJECT=$(jq -r .$ENV ./projects.json)
DATABASE_URL="$(pass palindr/$ENV/db)"
cat ./app.yaml.template | sed "s|\$DATABASE_URL|$DATABASE_URL|" > app.yaml
gcloud --project=$PROJECT app deploy
shred -uzn99 ./app.yaml
gcloud app logs tail -s default
