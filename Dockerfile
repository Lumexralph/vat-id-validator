FROM golang:1.16 AS builder
WORKDIR /app/
COPY go.mod /
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o vat-id-validator ./cmd

FROM alpine:latest as certs
RUN apk update && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/vat-id-validator /
ENTRYPOINT [ "./vat-id-validator" ]
