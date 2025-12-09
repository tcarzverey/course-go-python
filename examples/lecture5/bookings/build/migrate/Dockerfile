FROM golang:latest as builder

WORKDIR /build

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# mount the app
ADD ./migrations /build/migrations

ENTRYPOINT ["/go/bin/goose", "-dir=migrations", "postgres"]
CMD ["up"]