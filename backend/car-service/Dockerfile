FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /bin/car-service ./cmd

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /bin/car-service /app/car-service

EXPOSE 8082

CMD ["/app/car-service"]
