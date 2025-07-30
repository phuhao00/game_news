# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装git和nodejs用于获取依赖和构建前端
RUN apk add --no-cache git nodejs npm

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
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o game-news .

# 运行阶段
FROM alpine:latest

# 安装ca证书和时区数据
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/game-news .

# 复制前端构建文件
COPY --from=builder /app/dist ./dist

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./game-news"]
