FROM ubuntu:20.04
MAINTAINER Fabio Rehm "fgrehm@gmail.com"

RUN mkdir -p /app \
    && apt-get update \
    && apt-get install -y --no-install-recommends graphviz ca-certificates wget \
    && rm -rf /var/lib/apt/lists/*

ARG SLS_WEB_VERSION=0.1.0
RUN wget -O- "https://github.com/fgrehm/sls-web/releases/download/v${SLS_WEB_VERSION}/sls-web-linux-v${SLS_WEB_VERSION}.tgz" \
    | tar xvz --strip-components 1 -C /app \
    && adduser sls-web

USER sls-web
WORKDIR /app
ENTRYPOINT ["/app/sls-web"]
