#! /bin/bash

cd "$(dirname $0)"

ENV=$(cat ./environment)
if [ $ENV = "dev" ]; then
    PROTOCOL="http://"
    URL=$(ip addr | grep 192 | tail -n1 | xargs | cut -d" " -f2 | sed "s|/.*||")
    PORT=":3000"
else
    GCLOUD_PROJECT=$(jq -r .gcloud.$ENV ./projects.json)
    PROTOCOL="https://"
    URL=$(gcloud --project=$GCLOUD_PROJECT app describe \
              | grep defaultHostname | cut -d" " -f2)
    PORT=""
fi

printf "%s%s%s" $PROTOCOL $URL $PORT