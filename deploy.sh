#! /bin/bash

cd "$(dirname $0)"

ENV=$(cat ./environment)
GCLOUD_PROJECT=$(jq -r .gcloud.$ENV ./projects.json)
HEROKU_PROJECT=$(jq -r .heroku.$ENV ./projects.json)
DATABASE_URL="$(heroku pg:credentials:url -a $HEROKU_PROJECT \
                       | tail -n1 | xargs)"
cat ./app.yaml.template | sed "s|\$DATABASE_URL|$DATABASE_URL|" > ./app.yaml
if [ $1 = "dry" ]; then
    cat ./app.yaml
else
    gcloud --project=$GCLOUD_PROJECT app deploy
fi
shred -uzn99 ./app.yaml
