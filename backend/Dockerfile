ARG GO_VERSION=1.22.5
ARG DB_URL='postgres://postgres:postgres@postgres:5432/soskot'
ARG LOG_LEVEl=info
ARG PORT=:8080

FROM golang:$GO_VERSION AS builder

WORKDIR /app

COPY ./backend .

RUN go mod download

RUN go build ./cmd/main.go 

CMD ./main -db_url 'postgres://postgres:postgres@postgres:5432/soskot' -log_level info -port :8080
