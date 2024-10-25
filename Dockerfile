FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
COPY config ./config
COPY migrations ./migrations
RUN go mod download

COPY . .

RUN go build -o server ./cmd/servers

EXPOSE 8080

CMD ["./server"]
