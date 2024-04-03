#! /bin/bash

cd "$(dirname $0)"

ENV=$(cat ./environment)
printf "environment: %s\n" $ENV
GCLOUD_PROJECT=$(jq -r .gcloud.$ENV ./projects.json)
HEROKU_PROJECT=$(jq -r .heroku.$ENV ./projects.json)
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
URL=$(./url.sh)
echo $URL > ../pal/indr.url
printf "\nURL: %s\nsaved in ../pal/indr.url\n" $URL
