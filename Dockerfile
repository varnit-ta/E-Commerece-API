FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN ls -l /app

RUN CGO_ENABLED=0 GOOS=linux go build -o ecom-api ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/ecom-api .

EXPOSE 8080

CMD ["./ecom-api"]