ARG GO_VERSION=1.22.5
ARG GO_BASE_IMAGE=alpine

FROM golang:${GO_VERSION}-${GO_BASE_IMAGE} AS builder

WORKDIR /app

COPY ./backend .

RUN go mod download

RUN go build -o ./backend ./cmd/main.go
