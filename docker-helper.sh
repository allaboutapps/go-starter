#!/bin/bash

if [ "$1" = "--up" ]; then
    # --rm removes container after exit
    # --service-ports maps ports to outside when using run
    docker-compose run --rm --service-ports service bash
fi

if [ "$1" = "--halt" ]; then
    echo
    echo "Stopping db container ..."
    docker stop go-mranftl-sample_postgres
    echo
fi

if [ "$1" = "--destroy" ]; then
    echo
    echo "Stopping db container ..."
    docker stop go-mranftl-sample_postgres
    echo "Removing db container ..."
    docker container rm go-mranftl-sample_postgres
    echo "Removing service image ..."
    docker image rm  go-mranftl-sample_service
    echo "Removing db volume ..."
    docker volume rm go-mranftl-sample_pgvolume
    echo
fi
