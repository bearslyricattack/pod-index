# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pod-index .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/pod-index .

# ⭐ 添加执行权限 (关键修复)
RUN chmod +x /root/pod-index

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./pod-index"]
