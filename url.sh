#! /bin/bash

cd "$(dirname $0)"

ENV=$(cat ./environment)
GCLOUD_PROJECT=$(jq -r .gcloud.$ENV ./projects.json)
gcloud --project=$GCLOUD_PROJECT app describe \
    | grep defaultHostname | cut -d" " -f2
