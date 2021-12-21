FROM golang:1.17-alpine AS builder
ARG SVC
ARG GOARCH
ARG GOARM

WORKDIR /go/src/github.com/ns1labs/orb
COPY . .
RUN apk update \
    && apk add make
RUN make $SVC \
    && mv build/orb-$SVC /exe

FROM scratch
# Certificates are needed so that mailing util can work.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
