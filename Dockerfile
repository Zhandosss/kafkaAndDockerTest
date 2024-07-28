FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN  go mod download
COPY db ./db
COPY producer ./producer

EXPOSE 8000

RUN go build -o producer/main ./producer/main.go

FROM alpine:3 AS production
COPY --from=builder /build/producer/main /bin/main
COPY producer/configs /configs

ENTRYPOINT ["/bin/main"]