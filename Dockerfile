# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装git用于获取依赖
RUN apk add --no-cache git

# 复制go mod和sum文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建前端
RUN cd frontend && \
    npm install && \
    npm run build

# 构建Go应用
RUN go build -o game-news .

# 运行阶段
FROM alpine:latest

# 安装ca证书以便处理HTTPS请求
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/game-news .
COPY --from=builder /app/game_news.db .

# 复制前端构建文件
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/init.sql .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./game-news"]