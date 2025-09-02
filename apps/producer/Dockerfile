FROM golang:1.24-alpine as builder

WORKDIR /app

RUN echo 'precedence ::ffff:0:0/96  100' >> /etc/gai.conf

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o producer ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/producer .

CMD ["./producer"]