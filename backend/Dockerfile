FROM golang:1.24.3-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main ./cmd

EXPOSE 8080

ENTRYPOINT ["./main"]
