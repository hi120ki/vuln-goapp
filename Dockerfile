ARG GO_VERSION=1.19
ARG DEBIAN_VERSION=buster

# build-stage

FROM golang:${GO_VERSION}-${DEBIAN_VERSION} as build-stage

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build

# roduction-stage

FROM debian:${DEBIAN_VERSION} as production-stage

RUN groupadd -g 1000 appgroup && useradd -m -s /bin/bash -u 1000 -g 1000 appuser

WORKDIR /app

COPY --chown=appuser:appgroup --from=build-stage /app/vuln-goapp .

RUN chown -R appuser:appgroup /app

USER $USERNAME

CMD ./vuln-goapp
