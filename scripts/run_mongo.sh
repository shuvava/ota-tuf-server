#!/bin/bash
echo "running mongo docker image"

MONGO_USERNAME=${MONGO_USERNAME:-mongoadmin}
MONGO_PASSWD=${MONGO_PASSWD:-secret}

docker run -it --rm --name treehub-mongo \
  -e "MONGO_INITDB_ROOT_USERNAME=${MONGO_USERNAME}" \
  -e "MONGO_INITDB_ROOT_PASSWORD=${MONGO_PASSWD}" \
  -p 27017:27017 \
  mongo:5.0.11
