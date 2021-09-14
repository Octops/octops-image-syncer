FROM golang:1.17 as build-env

WORKDIR /go/src/github.com/Octops/octops-image-syncer
ADD . /go/src/github.com/Octops/octops-image-syncer

RUN go get -d -v ./...

RUN go build -o /go/bin/octops-image-syncer

FROM gcr.io/distroless/base

COPY --from=build-env /go/bin/octops-image-syncer /

CMD ["/octops-image-syncer"]