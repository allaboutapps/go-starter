#!/bin/sh

IMAGE_NAME=integresql

gcloud auth configure-docker
docker build -t ${IMAGE_NAME} .
docker tag ${IMAGE_NAME} eu.gcr.io/a3cloud-192413/${IMAGE_NAME}
docker push eu.gcr.io/a3cloud-192413/${IMAGE_NAME}

