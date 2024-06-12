# 使用 Go 1.22.4 作为构建阶段的基础镜像
FROM golang:1.22.4 AS builder

# 设置工作目录
WORKDIR /app

# 复制所有文件到工作目录
COPY . .

# 安装依赖
RUN go mod tidy

# 构建 Go 程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# 使用 Alpine 作为运行阶段的基础镜像
FROM alpine:latest

# 安装 ca-certificates 以便于进行 HTTPS 请求
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制可执行文件到运行阶段
COPY --from=builder /app/app .

# 暴露应用的端口
EXPOSE 8080

# 设置容器启动时执行的命令
CMD ["./app"]