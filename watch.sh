#! /bin/bash

cd "$(dirname $0)"

ENV=$(cat ./environment)
GCLOUD_PROJECT=$(jq -r .gcloud.$ENV ./projects.json)
gcloud --project=$GCLOUD_PROJECT app logs tail -s default
