FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest

# Install required packages including wget for health check
RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app
COPY --from=builder /app/main .

# Create a non-root user and switch to it
RUN adduser -D -g '' appuser
USER appuser

# Expose port 8080 for App Engine
EXPOSE 8080

# Add health check
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

CMD ["./main"] 