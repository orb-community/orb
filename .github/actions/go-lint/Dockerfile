FROM golang:1.15.1-buster

LABEL author="Rodrigo Pavan<rpoliveira@daitan.com>"
LABEL maintainer="Daitan Digital Solutions"
LABEL version="1.0.0"

RUN mkdir -p /workspace

WORKDIR /workspace

COPY ./entrypoint.sh /entrypoint.sh

RUN apt-get update \
    && apt-get install jq git -y \
    && wget https://github.com/cloudposse/github-commenter/releases/download/0.7.0/github-commenter_linux_amd64 -O /github-commenter \
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /golangci-lint v1.38.0 \
    && chmod +x /github-commenter /entrypoint.sh /golangci-lint

ENTRYPOINT ["/entrypoint.sh"]