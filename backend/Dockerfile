FROM golang:1.22.5 AS builder

WORKDIR /app

#COPY go.mod .

COPY ./backend .

RUN go mod download

#RUN go build -o main ./cmd/main.go
#RUN chmod +x main
#CMD ["./main"]
#"-db_url", "postgres://postgres:postgres@postgres:5432/soskot", "-log_level", "info", "-port", "postgres:5432"

CMD [ "go", "run", "cmd/main.go"]