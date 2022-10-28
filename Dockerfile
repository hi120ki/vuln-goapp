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

ARG USERNAME=user
ARG GROUPNAME=user
ARG UID=1000
ARG GID=1000
RUN groupadd -g $GID $GROUPNAME && \
    useradd -m -s /bin/bash -u $UID -g $GID $USERNAME

WORKDIR /app

COPY --chown=appuser:appgroup --from=build-stage /app/vuln-goapp .

RUN chown -R appuser:appgroup /app

USER $USERNAME

CMD ./vuln-goapp
