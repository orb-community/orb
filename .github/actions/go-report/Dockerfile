FROM golang:1.15.1-alpine3.12

LABEL author="Everton Taques<everton.taques@encora.com>"
LABEL maintainer="ns1labs"
LABEL version="1.0.0"

RUN mkdir -p /workspace

COPY ./entrypoint.sh /entrypoint.sh

RUN apk add git && \
wget https://github.com/cloudposse/github-commenter/releases/download/0.7.0/github-commenter_linux_amd64 -O /github-commenter && \
apk add jq && \
chmod +x /github-commenter /entrypoint.sh

RUN git clone https://github.com/gojp/goreportcard.git && \
cd goreportcard && \
make install && \
go install ./cmd/goreportcard-cli && \
mv /root/go/bin/goreportcard-cli /usr/local/bin/goreportcard-cli

WORKDIR /workspace

CMD ["/entrypoint.sh"]
