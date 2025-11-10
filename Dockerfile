# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pod-index .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# ⭐ 使用 /app 而不是 /root
WORKDIR /app

# ⭐ 复制文件到 /app
COPY --from=builder /app/pod-index .

# ⭐ 确保所有用户可执行
RUN chmod 755 pod-index && \
    chown 65534:65534 pod-index

EXPOSE 8080

# ⭐ 切换到非 root 用户
USER 65534

# ⭐ 使用绝对路径
CMD ["/app/pod-index"]
