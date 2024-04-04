#! /bin/bash

cd "$(dirname $0)"

ENV=$(cat ./ENV)
GCLOUD_PROJECT=$(jq -r .$ENV.project.gcloud ./settings.json)
gcloud --project=$GCLOUD_PROJECT app logs tail -s default
