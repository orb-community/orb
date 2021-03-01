FROM golang:latest AS gobuild
COPY . /src/
WORKDIR /src/
RUN go build -o orb-api cmd/orb-api/main.go

FROM debian:bullseye-slim AS runtime

COPY --from=gobuild /src/orb-api /usr/local/bin/orb-api

ENTRYPOINT /usr/local/bin/orb-api
