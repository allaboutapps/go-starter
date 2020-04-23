FROM golang:1.14.2 AS development

# https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md#walk-through
ENV GOBIN /app/bin
ENV PATH $GOBIN:$PATH

# Our Makefile / env fully supports parallel job execution
ENV MAKEFLAGS "-j 8 --no-print-directory"

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
    sudo \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Add user to avoid linux file permission issues
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# Create the user
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

# ********************************************************
# * Anything else you want to do like clean up goes here *
# ********************************************************

# [Optional] Set the default user. Omit if you want to keep the default as root.
#USER $USERNAME

# vscode support: LANG must be supported, requires installing the locale package first
# see https://github.com/Microsoft/vscode/issues/58015
RUN sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
    dpkg-reconfigure --frontend=noninteractive locales && \
    update-locale LANG=en_US.UTF-8

ENV LANG en_US.UTF-8

# sql-formatting: Install the same version of pg_formatter as used in your editors, as of 2020-04 thats v4.3
# https://github.com/bradymholt/vscode-pgFormatter/commits/master
# https://github.com/darold/pgFormatter/releases
RUN wget https://github.com/darold/pgFormatter/archive/v4.3.tar.gz \
    && tar xzf v4.3.tar.gz \
    && cd pgFormatter-4.3 \
    && perl Makefile.PL \
    && make && make install

# go richgo: (this package should NOT be installed via go get)
# https://github.com/kyoh86/richgo/releases
RUN wget https://github.com/kyoh86/richgo/releases/download/v0.3.3/richgo_0.3.3_linux_amd64.tar.gz \
    && tar xzf richgo_0.3.3_linux_amd64.tar.gz \
    && cp richgo /usr/local/bin/richgo

# go linting: (this package should NOT be installed via go get)
# https://github.com/golangci/golangci-lint#binary
# https://github.com/golangci/golangci-lint/releases
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b $(go env GOPATH)/bin v1.24.0

# go swagger: (this package should NOT be installed via go get) 
# https://github.com/go-swagger/go-swagger/releases
RUN curl -o /usr/local/bin/swagger -L'#' \
    "https://github.com/go-swagger/go-swagger/releases/download/v0.23.0/swagger_linux_amd64" \
    && chmod +x /usr/local/bin/swagger

### -----------------------
# --- Stage: builder
### -----------------------

FROM development as builder
WORKDIR /app
COPY Makefile /app/Makefile
COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
RUN make modules
COPY tools.go /app/tools.go
RUN make tools
COPY . /app/
