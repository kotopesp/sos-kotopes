FROM golang:1.22.5 AS builder

WORKDIR /app


COPY ./backend .

RUN go mod download

CMD [ "go", "run", "cmd/main.go", "-db_url", "postgres://postgres:postgres@postgres:5432/soskot", "-log_level", "info", "-port", ":80"]