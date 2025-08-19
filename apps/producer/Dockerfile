FROM golang:1.24-alpine as builder

WORKDIR /app

copy . .

RUN go build -o producer ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/producer .

CMD ["./producer"]