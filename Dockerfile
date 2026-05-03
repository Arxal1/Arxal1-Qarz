FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o qarzi-app ./cmd/api
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/qarzi-app .
COPY --from=builder /app/.env .
COPY --from=builder /app/migration ./migration

EXPOSE 8080
CMD ["./qarzi-app"]