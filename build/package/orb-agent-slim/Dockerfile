FROM golang:latest AS gobuild
COPY . /src/
WORKDIR /src/
RUN go build -o orb-agent cmd/orb-agent/main.go

FROM debian:bullseye-slim AS runtime

COPY --from=gobuild /src/orb-agent /usr/local/bin/orb-agent

ENTRYPOINT /usr/local/bin/orb-agent