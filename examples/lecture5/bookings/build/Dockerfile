FROM golang:latest AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/bookings ./cmd/

FROM alpine:latest

WORKDIR /build

COPY --from=builder /build/bin/bookings /build/bin/bookings

CMD ["./bin/bookings"]