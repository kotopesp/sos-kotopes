ARG GO_VERSION=1.22.5
ARG GO_BASE_IMAGE=alpine

FROM golang:${GO_VERSION}-${GO_BASE_IMAGE} AS builder

WORKDIR /app

COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download

COPY ./backend ./

RUN go build -o ./backend ./cmd/main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/backend .

COPY --from=builder /app/internal/data ./internal/data
