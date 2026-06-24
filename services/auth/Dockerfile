FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY sdk/go.mod sdk/go.sum* ./sdk/
COPY services/auth/go.mod services/auth/go.sum ./services/auth/

WORKDIR /app/services/auth
RUN go mod download

WORKDIR /app
COPY sdk ./sdk
COPY services/auth ./services/auth

WORKDIR /app/services/auth
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.21

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /server .

EXPOSE 5001

CMD ["./server"]
