FROM golang:1.15-alpine3.12 AS builder

COPY . /tmp/build
WORKDIR /tmp/build
RUN go build -o owm_exporter cmd/main.go

FROM alpine:3.12

WORKDIR /
COPY --from=builder /tmp/build/owm_exporter /owm_exporter

CMD ["/owm_exporter"]
