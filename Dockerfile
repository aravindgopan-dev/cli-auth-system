FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auth-cli ./cmd/app/main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/auth-cli /app/auth-cli
COPY --from=builder /app/migration /app/migration

ENTRYPOINT ["/app/auth-cli"]
