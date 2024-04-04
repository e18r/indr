#! /bin/bash

cd "$(dirname $0)"

if [ ! -e ./ENV ]; then
    echo "./ENV does not exist"
    exit 1
fi

ENV=$(cat ./ENV)
printf "ENV: %s\n" $ENV
if [ "$ENV" != "test" -o "$ENV" != "prod" ]; then
    echo "only deploy to test or prod envs"
    exit 1
fi
GCLOUD_PROJECT=$(jq -r .$ENV.project.gcloud ./settings.json)
HEROKU_PROJECT=$(jq -r .$ENV.project.heroku ./settings.json)
echo "obtaining database url..."
DATABASE_URL="$(heroku pg:credentials:url -a $HEROKU_PROJECT \
                       | tail -n1 | xargs)"
cat ./app.yaml.template | sed "s|\$DATABASE_URL|$DATABASE_URL|" > ./app.yaml
if [ "$1" = "dry" ]; then
    cat ./app.yaml
else
    gcloud --project=$GCLOUD_PROJECT app deploy
fi
shred -uzn99 ./app.yaml
