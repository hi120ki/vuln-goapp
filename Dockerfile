ARG GO_VERSION=1.19
ARG ALPINE_VERSION=3.16

# build-stage

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as build-stage

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build

# roduction-stage

FROM alpine:${ALPINE_VERSION} as production-stage

WORKDIR /app

COPY --from=build-stage /app/vuln-goapp .

CMD ./vuln-goapp
