FROM golang:1.26-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /my_app

ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/scheduler.db

CMD ["/my_app"]