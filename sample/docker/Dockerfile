# ============================================================================
# 多阶段构建 Dockerfile - Comigo 漫画/书籍浏览服务器
# ============================================================================
# 支持多平台构建: linux/amd64, linux/arm64, linux/arm/v7
# 使用方法：
#   docker build -t yumenaka/comigo:latest .
#   docker build --build-arg VERSION=v1.2.5 -t yumenaka/comigo:v1.2.5 .
# ============================================================================

# ============================================================================
# 第一阶段: 前端构建
# ============================================================================
FROM alpine:3.23 AS frontend-builder

WORKDIR /build

# 安装 bun（JavaScript 运行时和包管理器）
RUN apk add --no-cache curl unzip bash && \
    curl -fsSL https://bun.sh/install | bash && \
    ln -s /root/.bun/bin/bun /usr/local/bin/bun

# 复制前端相关文件
COPY package.json bun.lock* ./
COPY assets/ ./assets/
COPY tailwind.config.js ./

# 安装依赖并构建前端资源
RUN bun install && \
    bun run dev || true

# ============================================================================
# 第二阶段: Go 构建
# ============================================================================
FROM golang:alpine AS go-builder

# 构建参数：版本号
ARG VERSION
ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

# 安装构建依赖
RUN apk add --no-cache git make

# 复制 go.mod 和 go.sum，利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 从前端构建阶段复制构建产物
COPY --from=frontend-builder /build/assets/ ./assets/

# 生成 templ 模板文件
RUN go run github.com/a-h/templ/cmd/templ@latest generate

# 编译 Go 程序
# -trimpath: 去掉源码绝对路径
# -ldflags: -s 去掉符号信息，-w 去掉调试信息，减小二进制体积
RUN if [ -z "$VERSION" ]; then \
        VERSION=$(grep -o 'v[0-9]\+\.[0-9]\+\.[0-9]\+' config/version.go | head -1); \
    fi && \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath \
    -ldflags "-s -w -X 'github.com/yumenaka/comigo/config.version=${VERSION}'" \
    -o /build/comigo \
    .

# ============================================================================
# 第三阶段: 运行时镜像
# ============================================================================
FROM alpine:latest

# 构建参数：版本号（从构建阶段传递）
ARG VERSION

# 镜像元数据
LABEL maintainer="Comigo Project" \
      description="Comigo - 漫画/书籍浏览服务器，支持多种格式" \
      version="${VERSION}"

# 安装运行时依赖
# ca-certificates: HTTPS 支持
# tzdata: 时区数据
# wget: 健康检查使用
RUN apk add --no-cache ca-certificates tzdata wget && \
    rm -rf /var/cache/apk/*

# 创建必要的目录
# 配置文件目录使用标准的 HomeDirectory 路径
RUN mkdir -p /data /root/.config/comigo

WORKDIR /app

# 从构建阶段复制二进制文件
# 注意：资源文件已通过 Go embed 内嵌到二进制文件中，无需单独复制 assets 目录
COPY --from=go-builder /build/comigo /app/comigo

# 设置可执行权限
RUN chmod +x /app/comigo

# 暴露端口
EXPOSE 1234

# 设置环境变量
ENV COMIGO_PORT=1234 \
    COMIGO_ENABLE_UPLOAD=true \
    COMIGO_LANGUAGE=auto

# 数据卷
VOLUME ["/data", "/root/.config/comigo"]

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${COMIGO_PORT}/ || exit 1

# 启动命令
# 默认扫描 /data 目录
CMD ["/app/comigo", "/data"]
