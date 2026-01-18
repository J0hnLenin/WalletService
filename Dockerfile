FROM golang:1.25.6-alpine

WORKDIR /app

COPY ./internal ./internal
COPY ./config ./config
COPY ./cmd ./cmd

COPY go.mod go.sum ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

EXPOSE 8080

CMD ["./main"]