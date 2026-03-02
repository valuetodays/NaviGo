# 构建阶段
FROM golang:1.23.4-alpine3.20 AS builder

# 安装 gcc + libc-dev
RUN apk add --no-cache gcc musl-dev

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制所有源代码
COPY . .

# 构建Go程序（启用CGO以支持SQLite）
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o navigo main.go

# 运行阶段
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 安装SQLite依赖
RUN apk add --no-cache sqlite-libs

# 从构建阶段复制编译好的程序
COPY --from=builder /app/navigo .
COPY --from=builder /app/index.html .

# 声明数据卷（用于挂载数据库文件）
VOLUME ["/app/navigo.db"]

# 暴露端口
EXPOSE 3000

# 启动命令
CMD ["./navigo"]
