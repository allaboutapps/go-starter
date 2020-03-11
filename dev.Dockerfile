# No build here, only dev environment
FROM golang:1.14
# EXPOSE 8080

# Yes no maybe. This is strange. Although all default shells are bash and bash has been set as the shell for yarn/npm to use, 
# it still runs everything as /bin/sh for some weird reason. Let's make sure it doesn't. Naughty yarn. 
# RUN rm /bin/sh \ 
#     && ln -s /bin/bash /bin/sh

# We use ssmtp (provides a sendmail binary) to relay emails to servers (must be custom configured)
# https://blog.philippklaus.de/2011/03/set-up-sending-emails-on-a-local-system-by-transfering-it-to-a-smtp-relay-server-smarthost
# Configure by mounting /etc/ssmtp/ssmtp.conf
# RUN apt-get update \
#     && apt-get install -y ssmtp \
#     && apt-get clean \
#     && rm -rf /var/lib/apt/lists/*

# Install required system dependencies
# E.g.
# RUN set -e \
#     && apt-get update \
#     && apt-get install -y --no-install-recommends \
#     imagemagick \
#     && rm -rf /var/lib/apt/lists/*

# Comment int, if using psql cli in tests
# TESTS_ONLY: We need an psql client (any version, fuckit) to execute tests that utilize the psql cli
# RUN set -e \
#     && apt-get update \
#     && apt-get install -y --no-install-recommends \
#     postgresql-client \
#     && rm -rf /var/lib/apt/lists/*
