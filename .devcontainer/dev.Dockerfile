FROM golang:1.14 AS development

# https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md#walk-through
ENV GOBIN /app/bin
ENV PATH $GOBIN:$PATH

# postgresql-support: Add the official postgres repo to install the matching postgresql-client tools of your stack
# see https://wiki.postgresql.org/wiki/Apt
# run lsb_release -c inside the container to pick the proper repository flavor
# e.g. stretch=>stretch-pgdg, buster=>buster-pgdg
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ buster-pgdg main" \
    | tee /etc/apt/sources.list.d/pgdg.list \
    && wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc \
    | apt-key add -

# Install required system dependencies
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    locales \
    postgresql-client-12 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# vscode support: LANG must be supported, requires installing the locale package first
# see https://github.com/Microsoft/vscode/issues/58015
RUN sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
    dpkg-reconfigure --frontend=noninteractive locales && \
    update-locale LANG=en_US.UTF-8

ENV LANG en_US.UTF-8

# sql-formatting: Install the same version of pg_formatter as used in your editors, as of 2020-03 thats v4.2
# https://github.com/darold/pgFormatter/releases
# https://github.com/bradymholt/vscode-pgFormatter/commits/master
RUN wget https://github.com/darold/pgFormatter/archive/v4.2.tar.gz \
    && tar xzf v4.2.tar.gz \
    && cd pgFormatter-4.2 \
    && perl Makefile.PL \
    && make && make install

# go linting: (this package should NOT be installed via go get)
# https://github.com/golangci/golangci-lint#binary
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b $(go env GOPATH)/bin v1.24.0
