
ARG GO_VERSION=1.19

# build-stage

FROM golang:${GO_VERSION}-bullseye as build-stage

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

# roduction-stage

FROM gcr.io/distroless/static-debian11 as production-stage

WORKDIR /app

COPY --chown=nonroot:nonroot --from=build-stage /app/vuln-goapp .

USER nonroot

CMD ["/app/vuln-goapp"]
