#! /bin/bash

cd "$(dirname $0)"

if [ "$ENV" != "test" -a "$ENV" != "prod" ]; then
    echo "only watch on test or prod envs"
    exit 1
fi

ENV=$(cat ./ENV)
GCLOUD_PROJECT=$(jq -r .$ENV.project.gcloud ./settings.json)
gcloud --project=$GCLOUD_PROJECT app logs tail -s default
