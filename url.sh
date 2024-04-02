#! /bin/bash

cd "$(dirname $0)"

ENV=$(cat ./environment)
PROJECT=$(jq -r .$ENV ./projects.json)
gcloud --project=$PROJECT app describe | grep defaultHostname | cut -d" " -f2
