FROM golang:1.21.0-alpine

WORKDIR /app

COPY . ./src
WORKDIR ./src

RUN go build -o myapp cmd/main.go

CMD ["./myapp"]