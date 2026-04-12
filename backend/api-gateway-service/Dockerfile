FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api-gateway ./cmd/main.go

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /bin/api-gateway /usr/local/bin/api-gateway

EXPOSE 8080

CMD ["api-gateway"]
