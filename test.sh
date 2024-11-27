#! /bin/bash

cd "$(dirname $0)"

ENV=$(cat ./ENV)

if [ "${ENV}" != "test" ]; then
    echo "not in test env"
    exit 1
fi

SLEEP=${1}

if [ -n "${SLEEP}" ] && [ "${SLEEP}" -eq "${SLEEP}" ]; then
    echo "Sending random POST requests to the API..."
    echo ""
else
    echo "usage ./test.sh INTERVAL"
    exit
fi

PROJECT=$(jq -r .${ENV}.project.gcloud ./settings.json)
URL=$(gcloud --project=$PROJECT app describe | \
          grep defaultHostname | cut -d' ' -f2)

while true; do
    SEQUENCE=$(cat /dev/random | tr -cd "[:alnum:]" | head -c 8)
    PAL=$(printf ${SEQUENCE}; echo ${SEQUENCE} | rev)
    echo "POST ${URL}/publish text=${PAL}"
    http ${URL}/publish text="${PAL}"
    sleep $SLEEP
done
