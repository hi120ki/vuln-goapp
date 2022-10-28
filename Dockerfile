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

RUN addgroup appgroup && adduser --disabled-password --no-create-home appuser -G appgroup

WORKDIR /app

COPY --chown=appuser:appgroup --from=build-stage /app/vuln-goapp .

USER appuser

CMD ./vuln-goapp
