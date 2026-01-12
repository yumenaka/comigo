# Comigo Docker 镜像使用指南

[English](#english) | [中文](#中文) | [日本語](#日本語)

---

## 中文

### 📖 简介

Comigo 是一个功能强大的漫画/书籍浏览服务器，支持多种格式（ZIP、RAR、PDF 等）。本 Docker 镜像提供了便捷的部署方式，支持多架构平台。

### 🚀 快速开始

#### 使用 Docker 命令

```bash
# 拉取镜像
docker pull yumenaka/comigo:latest

# 运行容器
docker run -d \
  --name comigo \
  -p 1234:1234 \
  -v /path/to/your/books:/data \
  yumenaka/comigo:latest
```

访问 `http://localhost:1234` 开始使用。

#### 使用 Docker Compose

1. 下载 `docker-compose.yml` 文件
2. 编辑配置（可选）
3. 启动服务：

```bash
cd sample/docker
docker-compose up -d
```

### 📋 支持的架构

- `linux/amd64` - 标准 x86_64 服务器
- `linux/arm64` - ARM64 服务器（树莓派 4/5 64位系统）
- `linux/arm/v7` - ARMv7 设备（树莓派 2-4 32位系统）

### ⚙️ 配置说明

#### 环境变量

| 变量名 | 说明 | 默认值      |
|--------|------|----------|
| `COMIGO_PORT` | 服务端口 | `1234`   |
| `COMIGO_USERNAME` | 登录用户名 | 空（不启用认证） |
| `COMIGO_PASSWORD` | 登录密码 | 空（不启用认证） |
| `COMIGO_ENABLE_UPLOAD` | 启用上传功能 | `true`   |
| `COMIGO_LANGUAGE` | 终端输出语言 | `auto`   |
| `COMIGO_DEBUG` | 调试模式 | `false`  |
| `COMIGO_MAX_DEPTH` | 最大扫描深度 | `5`      |
| `COMIGO_MIN_IMAGE` | 最少图片数量 | `1`      |

更多环境变量请参考 `docker-compose.yml` 中的注释。

#### 数据卷

| 容器路径 | 说明 | 必需 |
|----------|------|------|
| `/data` | 书库数据目录 | ✅ |
| `/root/.config/comigo` | 配置文件目录 | ❌ |

### 📝 使用示例

#### 基础使用

```bash
docker run -d \
  --name comigo \
  -p 1234:1234 \
  -v /home/user/manga:/data \
  yumenaka/comigo:latest
```

#### 启用认证

```bash
docker run -d \
  --name comigo \
  -p 1234:1234 \
  -v /home/user/manga:/data \
  -e COMIGO_USERNAME=admin \
  -e COMIGO_PASSWORD=your_secure_password \
  yumenaka/comigo:latest
```

#### 自定义端口

```bash
docker run -d \
  --name comigo \
  -p 8080:8080 \
  -v /home/user/manga:/data \
  -e COMIGO_PORT=8080 \
  yumenaka/comigo:latest
```

#### 中文界面

```bash
docker run -d \
  --name comigo \
  -p 1234:1234 \
  -v /home/user/manga:/data \
  -e COMIGO_LANGUAGE=zh \
  yumenaka/comigo:latest
```

### 🔧 进阶配置

#### 使用配置文件

创建 `config.toml` 文件：

```toml
Port = 1234
Host = ""
EnableUpload = true
Language = "zh"
MaxScanDepth = 5
MinImageNum = 1
```

挂载配置文件：

```bash
docker run -d \
  --name comigo \
  -p 1234:1234 \
  -v /home/user/manga:/data \
  -v /home/user/comigo/config.toml:/root/.config/comigo/config.toml \
  yumenaka/comigo:latest
```

#### 反向代理（Nginx）

```nginx
server {
    listen 80;
    server_name comics.example.com;

    location / {
        proxy_pass http://localhost:1234;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

#### 资源限制

```bash
docker run -d \
  --name comigo \
  -p 1234:1234 \
  -v /home/user/manga:/data \
  --memory="512m" \
  --cpus="1" \
  yumenaka/comigo:latest
```

### 🛠️ 构建镜像

#### 使用构建脚本

```bash
# 本地测试构建
cd sample/docker
chmod +x build.sh
./build.sh --local

# 构建并推送多平台镜像
./build.sh --version v1.2.5 --push --latest
```

#### 使用 Makefile

```bash
# 查看所有可用命令
make -f sample/docker/Makefile.docker help

# 构建本地镜像
make -f sample/docker/Makefile.docker docker-build

# 构建多平台镜像
make -f sample/docker/Makefile.docker docker-build-all

# 推送到 Docker Hub
make -f sample/docker/Makefile.docker docker-push

# 本地测试
make -f sample/docker/Makefile.docker docker-test
```

#### 手动构建

```bash
# 切换到项目根目录
cd /path/to/comigo

# 构建单平台镜像
docker build \
  --build-arg VERSION=v1.2.5 \
  -t yumenaka/comigo:v1.2.5 \
  -f sample/docker/Dockerfile \
  .

# 构建多平台镜像（需要 Docker Buildx）
docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  --build-arg VERSION=v1.2.5 \
  -t yumenaka/comigo:v1.2.5 \
  -t yumenaka/comigo:latest \
  -f sample/docker/Dockerfile \
  --push \
  .
```

### 🐛 故障排除

#### 容器无法启动

```bash
# 查看日志
docker logs comigo

# 检查容器状态
docker ps -a | grep comigo
```

#### 端口被占用

```bash
# 修改映射端口
docker run -d \
  --name comigo \
  -p 8080:1234 \
  -v /home/user/manga:/data \
  yumenaka/comigo:latest
```

#### 权限问题

```bash
# 检查挂载目录权限
ls -ld /path/to/your/books

# 修改目录所有者（UID 1000）
sudo chown -R 1000:1000 /path/to/your/books
```

#### 数据库问题

```bash
# 清理数据重新扫描
docker exec comigo rm -f /root/.config/comigo/*
docker restart comigo
```

### 📚 常用命令

```bash
# 查看日志
docker logs -f comigo

# 进入容器
docker exec -it comigo sh

# 重启容器
docker restart comigo

# 停止容器
docker stop comigo

# 删除容器
docker rm comigo

# 更新镜像
docker pull yumenaka/comigo:latest
docker stop comigo
docker rm comigo
# 重新运行容器（使用新镜像）
```

### 🔗 相关链接

- [项目主页](https://github.com/yumenaka/comigo)
- [问题反馈](https://github.com/yumenaka/comigo/issues)
- [Docker Hub](https://hub.docker.com/r/yumenaka/comigo)

---

## English

### 📖 Introduction

Comigo is a powerful manga/book browsing server supporting multiple formats (ZIP, RAR, PDF, etc.). This Docker image provides a convenient deployment method with multi-architecture support.

### 🚀 Quick Start

#### Using Docker Command

```bash
# Pull the image
docker pull yumenaka/comigo:latest

# Run the container
docker run -d \
  --name comigo \
  -p 1234:1234 \
  -v /path/to/your/books:/data \
  yumenaka/comigo:latest
```

Visit `http://localhost:1234` to get started.

#### Using Docker Compose

1. Download the `docker-compose.yml` file
2. Edit configuration (optional)
3. Start the service:

```bash
cd sample/docker
docker-compose up -d
```

### 📋 Supported Architectures

- `linux/amd64` - Standard x86_64 servers
- `linux/arm64` - ARM64 servers (Raspberry Pi 4/5 64-bit)
- `linux/arm/v7` - ARMv7 devices (Raspberry Pi 2-4 32-bit)

### ⚙️ Configuration

#### Environment Variables

| Variable | Description | Default         |
|----------|-------------|-----------------|
| `COMIGO_PORT` | Service port | `1234`          |
| `COMIGO_USERNAME` | Login username | Empty (no auth) |
| `COMIGO_PASSWORD` | Login password | Empty (no auth) |
| `COMIGO_ENABLE_UPLOAD` | Enable upload | `true`          |
| `COMIGO_LANGUAGE` | UI language | `auto`          |
| `COMIGO_DEBUG` | Debug mode | `false`         |

See `docker-compose.yml` for more environment variables.

#### Volumes

| Container Path | Description | Required |
|----------------|-------------|----------|
| `/data` | Book library directory | ✅ |
| `/root/.config/comigo` | Configuration directory | ❌ |

### 📝 Usage Examples

#### Basic Usage

```bash
docker run -d \
  --name comigo \
  -p 1234:1234 \
  -v /home/user/manga:/data \
  yumenaka/comigo:latest
```

#### With Authentication

```bash
docker run -d \
  --name comigo \
  -p 1234:1234 \
  -v /home/user/manga:/data \
  -e COMIGO_USERNAME=admin \
  -e COMIGO_PASSWORD=your_secure_password \
  yumenaka/comigo:latest
```

#### With Cache and Database

```bash
docker run -d \
  --name comigo \
  -p 1234:1234 \
  -v /home/user/manga:/data \
  -v /home/user/comigo/config:/root/.config/comigo \
  yumenaka/comigo:latest
```

---

## 日本語

### 📖 概要

Comigo は、複数のフォーマット（ZIP、RAR、PDF など）をサポートする強力な漫画/書籍閲覧サーバーです。この Docker イメージは、マルチアーキテクチャをサポートする便利なデプロイ方法を提供します。

### 🚀 クイックスタート

#### Docker コマンドを使用

```bash
# イメージをプル
docker pull yumenaka/comigo:latest

# コンテナを実行
docker run -d \
  --name comigo \
  -p 1234:1234 \
  -v /path/to/your/books:/data \
  yumenaka/comigo:latest
```

`http://localhost:1234` にアクセスして使用を開始します。

#### Docker Compose を使用

1. `docker-compose.yml` ファイルをダウンロード
2. 設定を編集（オプション）
3. サービスを開始：

```bash
cd sample/docker
docker-compose up -d
```

### 📋 サポートされているアーキテクチャ

- `linux/amd64` - 標準的な x86_64 サーバー
- `linux/arm64` - ARM64 サーバー（Raspberry Pi 4/5 64ビット）
- `linux/arm/v7` - ARMv7 デバイス（Raspberry Pi 2-4 32ビット）

### ⚙️ 設定

#### 環境変数

| 変数名 | 説明 | デフォルト値  |
|--------|------|---------|
| `COMIGO_PORT` | サービスポート | `1234`  |
| `COMIGO_USERNAME` | ログインユーザー名 | 空（認証なし） |
| `COMIGO_PASSWORD` | ログインパスワード | 空（認証なし） |
| `COMIGO_ENABLE_UPLOAD` | アップロードを有効化 | `true`  |
| `COMIGO_LANGUAGE` | UI言語 | `auto`  |
| `COMIGO_DEBUG` | デバッグモード | `false` |

その他の環境変数については `docker-compose.yml` を参照してください。

---

### 📄 ライセンス

このプロジェクトのライセンスについては、プロジェクトのルートディレクトリにある LICENSE ファイルを参照してください。

### 🤝 貢献

問題の報告、機能のリクエスト、プルリクエストを歓迎します！

### ⭐ サポート

このプロジェクトが役に立った場合は、GitHub でスターを付けてください！
