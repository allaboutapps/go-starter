### -----------------------
# --- Stage: development
# --- Purpose: Local development environment
# --- https://hub.docker.com/_/golang
# --- https://github.com/microsoft/vscode-remote-try-go/blob/master/.devcontainer/Dockerfile
### -----------------------
FROM golang:1.15.3 AS development

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

# Our Makefile / env fully supports parallel job execution
ENV MAKEFLAGS "-j 8 --no-print-directory"

# postgresql-support: Add the official postgres repo to install the matching postgresql-client tools of your stack
# https://wiki.postgresql.org/wiki/Apt
# run lsb_release -c inside the container to pick the proper repository flavor
# e.g. stretch=>stretch-pgdg, buster=>buster-pgdg
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ buster-pgdg main" \
    | tee /etc/apt/sources.list.d/pgdg.list \
    && wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc \
    | apt-key add -

# Install required system dependencies
RUN apt-get update \
    && apt-get install -y \
    #
    # Mandadory minimal linux packages
    # Installed at development stage and app stage
    # Do not forget to add mandadory linux packages to the final app Dockerfile stage below!
    # 
    # -- START MANDADORY --
    ca-certificates \
    # --- END MANDADORY ---
    # 
    # Development specific packages
    # Only installed at development stage and NOT available in the final Docker stage
    # based upon
    # https://github.com/microsoft/vscode-remote-try-go/blob/master/.devcontainer/Dockerfile
    # https://raw.githubusercontent.com/microsoft/vscode-dev-containers/master/script-library/common-debian.sh
    #
    # -- START DEVELOPMENT --
    apt-utils \
    dialog \
    openssh-client \
    less \
    iproute2 \
    procps \
    lsb-release \
    locales \
    sudo \
    bash-completion \
    bsdmainutils \
    graphviz \
    postgresql-client-12 \
    # --- END DEVELOPMENT ---
    # 
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# env/vscode support: LANG must be supported, requires installing the locale package first
# https://github.com/Microsoft/vscode/issues/58015
# https://stackoverflow.com/questions/28405902/how-to-set-the-locale-inside-a-debian-ubuntu-docker-container
RUN sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
    dpkg-reconfigure --frontend=noninteractive locales && \
    update-locale LANG=en_US.UTF-8

ENV LANG en_US.UTF-8

# sql pgFormatter: Install the same version of pg_formatter as used in your editors, as of 2020-04 thats v4.3
# requires perl to be installed
# https://github.com/bradymholt/vscode-pgFormatter/commits/master
# https://github.com/darold/pgFormatter/releases
RUN mkdir -p /tmp/pgFormatter \
    && cd /tmp/pgFormatter \
    && wget https://github.com/darold/pgFormatter/archive/v4.3.tar.gz \
    && tar xzf v4.3.tar.gz \
    && cd pgFormatter-4.3 \
    && perl Makefile.PL \
    && make && make install \
    && rm -rf /tmp/pgFormatter 

# go gotestsum: (this package should NOT be installed via go get)
# https://github.com/gotestyourself/gotestsum/releases
RUN mkdir -p /tmp/gotestsum \
    && cd /tmp/gotestsum \
    && wget https://github.com/gotestyourself/gotestsum/releases/download/v0.5.3/gotestsum_0.5.3_linux_amd64.tar.gz \
    && tar xzf gotestsum_0.5.3_linux_amd64.tar.gz \
    && cp gotestsum /usr/local/bin/gotestsum \
    && rm -rf /tmp/gotestsum 

# go linting: (this package should NOT be installed via go get)
# https://github.com/golangci/golangci-lint#binary
# https://github.com/golangci/golangci-lint/releases
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b $(go env GOPATH)/bin v1.31.0

# go swagger: (this package should NOT be installed via go get) 
# https://github.com/go-swagger/go-swagger/releases
RUN curl -o /usr/local/bin/swagger -L'#' \
    "https://github.com/go-swagger/go-swagger/releases/download/v0.25.0/swagger_linux_amd64" \
    && chmod +x /usr/local/bin/swagger

# golicense: (this package should NOT be installed via go get) 
# https://github.com/mitchellh/golicense/releases
RUN mkdir -p /tmp/golicense \
    && cd /tmp/golicense \
    && wget https://github.com/mitchellh/golicense/releases/download/v0.2.0/golicense_0.2.0_linux_x86_64.tar.gz \
    && tar xzf golicense_0.2.0_linux_x86_64.tar.gz \
    && cp golicense /usr/local/bin/golicense \
    && rm -rf /tmp/golicense

# linux permissions / vscode support: Add user to avoid linux file permission issues
# Detail: Inside the container, any mounted files/folders will have the exact same permissions
# as outside the container - including the owner user ID (UID) and group ID (GID). 
# Because of this, your container user will either need to have the same UID or be in a group with the same GID.
# The actual name of the user / group does not matter. The first user on a machine typically gets a UID of 1000,
# so most containers use this as the ID of the user to try to avoid this problem.
# 2020-04: docker-compose does not support passing id -u / id -g as part of its config, therefore we assume uid 1000
# https://code.visualstudio.com/docs/remote/containers-advanced#_adding-a-nonroot-user-to-your-dev-container
# https://code.visualstudio.com/docs/remote/containers-advanced#_creating-a-nonroot-user
ARG USERNAME=development
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd -s /bin/bash --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

# vscode support: cached extensions install directory
# https://code.visualstudio.com/docs/remote/containers-advanced#_avoiding-extension-reinstalls-on-container-rebuild
RUN mkdir -p /home/$USERNAME/.vscode-server/extensions \
    /home/$USERNAME/.vscode-server-insiders/extensions \
    && chown -R $USERNAME \
    /home/$USERNAME/.vscode-server \
    /home/$USERNAME/.vscode-server-insiders

# linux permissions / vscode support: chown $GOPATH so $USERNAME can directly work with it
# Note that this should be the final step after installing all build deps 
RUN mkdir -p /$GOPATH/pkg && chown -R $USERNAME /$GOPATH

# $GOBIN is where our own compiled binaries will live and other go.mod / VSCode binaries will be installed.
# It should always come AFTER our other $PATH segments and should be earliest targeted in stage "builder", 
# as /app/bin will the shadowed by a volume mount via docker-compose!
# E.g. "which golangci-lint" should report "/go/bin" not "/app/bin" (where VSCode will place it).
# https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md#walk-through
WORKDIR /app
ENV GOBIN /app/bin
ENV PATH $PATH:$GOBIN

### -----------------------
# --- Stage: builder
# --- Purpose: Statically built binaries and CI environment
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
RUN make go-build

### -----------------------
# --- Stage: app
# --- Purpose: Image for actual deployment
# --- Prefer https://github.com/GoogleContainerTools/distroless over
# --- debian:buster-slim https://hub.docker.com/_/debian (if you need apt-get).
### -----------------------

# Distroless images are minimal and lack shell access.
# https://github.com/GoogleContainerTools/distroless/blob/master/base/README.md
# The :debug image provides a busybox shell to enter (base-debian10 only, not static).
# https://github.com/GoogleContainerTools/distroless#debug-images
FROM gcr.io/distroless/base-debian10:debug as app

# FROM debian:buster-slim as app
# RUN apt-get update \
#     && apt-get install -y \
#     #
#     # Mandadory minimal linux packages
#     # Installed at development stage and app stage
#     # Do not forget to add mandadory linux packages to the base development Dockerfile stage above!
#     #
#     # -- START MANDADORY --
#     ca-certificates \
#     # --- END MANDADORY ---
#     #
#     && apt-get clean \
#     && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/app /app/
COPY --from=builder /app/api/swagger.yml /app/api/
COPY --from=builder /app/assets /app/assets/
COPY --from=builder /app/migrations /app/migrations/
COPY --from=builder /app/web /app/web/

WORKDIR /app

# Must comply to vector form
# https://github.com/GoogleContainerTools/distroless#entrypoints
# Sample usage of this image:
# docker run <image> help
# docker run <image> db migrate
# docker run <image> db seed
# docker run <image> env
# docker run <image> probe readiness
# docker run <image> probe liveness
# docker run <image> server
# docker run <image> server --migrate
ENTRYPOINT ["/app/app"]
CMD ["server", "--migrate"]