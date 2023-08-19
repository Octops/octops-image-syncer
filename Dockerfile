FROM golang:1.20 as build-env

WORKDIR /go/src/github.com/Octops/octops-image-syncer

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV APP_BIN /go/bin/octops-image-syncer
ENV VERSION v0.1.1

RUN make build

FROM gcr.io/distroless/static:nonroot

COPY --from=build-env /go/bin/octops-image-syncer /

ENTRYPOINT ["/octops-image-syncer"]
