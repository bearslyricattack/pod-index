# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pod-index .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# ⭐ Use /app instead of /root
WORKDIR /app

# ⭐ Copy files to /app
COPY --from=builder /app/pod-index .

# ⭐ Ensure executable by all users
RUN chmod 755 pod-index && \
    chown 65534:65534 pod-index

EXPOSE 8080

# ⭐ Switch to non-root user
USER 65534

# ⭐ Use absolute path
CMD ["/app/pod-index"]
