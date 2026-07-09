FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:3.20
RUN apk add --no-cache openssl
WORKDIR /app
COPY --from=builder /app/server .
RUN mkdir -p certs && \
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
      -keyout certs/server.key \
      -out certs/server.crt \
      -subj "/C=AR/ST=Cordoba/L=Cordoba/O=Ceibo/CN=localhost"
EXPOSE 8443
CMD ["./server"]
