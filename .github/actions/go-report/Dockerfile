FROM golang:1.20-alpine

LABEL author="Everton Taques<everton.taques@encora.com>"
LABEL maintainer="ns1labs"
LABEL version="1.0.0"

RUN mkdir -p /workspace

COPY ./entrypoint.sh /entrypoint.sh

RUN apk add git make curl jq && \
wget https://github.com/cloudposse/github-commenter/releases/download/0.7.0/github-commenter_linux_amd64 -O /github-commenter && \
chmod +x /github-commenter /entrypoint.sh

RUN git clone https://github.com/gojp/goreportcard.git && \
cd goreportcard && \
make install && \
go install ./cmd/goreportcard-cli && \
mv /go/bin/goreportcard-cli /usr/local/bin/goreportcard-cli

WORKDIR /workspace

CMD ["/entrypoint.sh"]
