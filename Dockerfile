# 使用轻量级 Golang 构建镜像
FROM golang:1.24-alpine3.21 AS builder

# 配置 Golang 代理
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 以利用 Docker 缓存
COPY go.mod go.sum ./

# 预下载依赖，提高构建效率
RUN go mod download && go mod tidy

# 复制源代码
COPY . .

# 编译 Go 程序（去除符号表和调试信息以减小体积）
RUN go build -ldflags="-s -w" -o main

# 选择更小的运行时镜像
FROM alpine:3.21

WORKDIR /app

# 设置为生产环境
ENV APP_ENV=production

# 复制编译好的二进制文件
COPY --from=builder /app/main .

# 运行所需的依赖
RUN apk --no-cache add ca-certificates

# 开放端口
EXPOSE 9001

# 运行应用
CMD ["./main", "--mode", "production", "--config", "/app/config/stellux.production.yaml"]
