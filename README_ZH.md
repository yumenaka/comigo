<div align="center">

# ComiGo：简单高效的漫画阅读器
[![Go Report](https://goreportcard.com/badge/github.com/yumenaka/comi?style=flat-square)](https://goreportcard.com/report/github.com/yumenaka/comi)
[![License](https://img.shields.io/github/license/yumenaka/comi?style=flat-square&color=blue)](https://github.com/yumenaka/comigo/blob/main/LICENSE)     
[English](https://github.com/yumenaka/comigo/blob/master/README.md) | [日本語](https://github.com/yumenaka/comigo/blob/master/README_JP.md) | [中文文档](https://github.com/yumenaka/comigo/blob/master/README_ZH.md)  
<!--
[![Downloads](https://img.shields.io/github/downloads/yumenaka/comi/total?style=flat-square&color=success)](https://github.com/yumenaka/comigo/releases)
<img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo：简单方便的漫画阅读器" width="200">
-->
</div>

![Windows示例](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windows示例")

## 功能特点

- 📚 **多格式支持**：支持图片文件夹与 `.rar`、`.zip`、`.tar`、`.cbz`、`.cbr`、`.epub` 等压缩包格式
- 📱 **便捷访问**：支持手机/平板扫描二维码阅读，Windows 支持拖拽打开
- 🐧 **跨平台支持**：适配 Windows、Linux、MacOS 系统
- 📖 **多样化阅读模式**：提供卷轴、翻页等多种阅读模式
- ⚙️ **灵活配置**：支持命令行操作，可通过 `config.toml` 配置文件设定书库
- 🖼️ **现代图片格式**：除了常见的`jpg`、`png`，还支持 `heic`、`avif` 等新一代图片格式
- ✂️ **智能优化**：支持图片自动裁边，压缩图片节省流量
- 🔄 **同步阅读**：支持不同设备间同步翻页进度
- 🔌 **插件系统**：支持自动翻页、时钟等插件，可扩展自定义插件
- 🎬 **媒体播放**：内置音频、视频播放器
- 📥 **灵活下载**：支持打包下载图片文件夹，支持转换并下载为 EPUB 格式
- 📜 **阅读历史**：自动记录阅读历史，方便续读

## 安装指南

### 托盘版（推荐）

| 系统 | 下载链接 |
|------|---------|
| Windows 64位 | [comigo_latest_Windows_x86_64_full.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comigo_latest_Windows_x86_64_full.zip) |
| macOS (Intel/Apple芯片) | [Comigo_latest.dmg](https://comigo.xyz/yumenaka/comigo/releases/download/latest/Comigo_latest.dmg) |

> 💡 **说明**：托盘版提供系统托盘图标，可最小化到后台运行。
>  Windows 用户双击运行即可，macOS 用户需要将APP拖入 应用程序 文件夹。

### 一键安装(命令行版)

```bash
# 中国大陆用户推荐使用中转脚本：
bash <(curl -s https://comigo.xyz/get.sh) --cn

# 从 GitHub下载：
bash <(curl -s https://raw.githubusercontent.com/yumenaka/comigo/master/get.sh)

# 如果您已设置 Golang 环境：
go install github.com/yumenaka/comigo/cmd/comi@latest
```

### 下载命令行版

| 系统类型 | 下载链接 |
|---------|---------|
| Windows 64位 | [comi_latest_Windows_x86_64.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Windows_x86_64.zip) |
| Windows ARM | [comi_latest_Windows_arm64.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Windows_arm64.zip) |
| macOS Intel | [comi_latest_MacOS_x86_64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_MacOS_x86_64.tar.gz) |
| macOS Apple芯片 | [comi_latest_MacOS_arm64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_MacOS_arm64.tar.gz) |
| Linux 64位 | [comi_latest_Linux_x86_64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Linux_x86_64.tar.gz) |
| Linux ARM64 | [comi_latest_Linux_arm64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Linux_arm64.tar.gz) |
| Linux ARM32 | [comi_latest_Linux_armv7.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Linux_armv7.tar.gz) |
| Debian/Ubuntu 64位 | [comi_latest_amd64.deb](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_amd64.deb) |
| Debian/Ubuntu ARM64 | [comi_latest_arm64.deb](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_arm64.deb) |

> 💡 **说明**：命令行版适合服务器部署和高级用户。下载后需手动添加到系统 PATH 环境变量。

### 发布页

也可以在 [Releases 页面](https://github.com/yumenaka/comigo/releases) 下载最新版本，并将可执行文件添加到系统的 `PATH` 环境变量中。


## 使用Docker 部署

### 快速开始

```bash
# 拉取并运行最新镜像
docker run -d \
  --name comigo \
  -p 1234:1234 \
  -v /path/to/your/books:/data \
  yumenaka/comigo:latest
```

访问 `http://localhost:1234` 即可使用。

### 使用 Docker Compose

1. 下载 [`docker-compose.yml`](sample/docker/docker-compose.yml) 文件
2. 根据需要编辑配置
3. 启动服务：

```bash
docker-compose up -d
```

更多详细说明请查看完整的 [Docker 使用文档](sample/docker/README.md)。

## 使用方法

```bash
comi [flags] file_or_dir
```

### 命令行参数

| 参数 | 简写 | 默认值 | 说明                   |
|------|------|--------|----------------------|
| `--config` | `-c` | - | 指定配置文件路径             |
| `--port` | `-p` | 1234 | 服务端口                 |
| `--host` | - | - | 自定义主机名/域名            |
| `--local` | - | false | 仅本地访问，不对局域网开放        |
| `--max-depth` | `-m` | 5 | 文件扫描最大深度             |
| `--open-browser` | `-o` | false | 启动后自动打开浏览器           |
| `--enable-upload` | - | true | 启用文件上传功能             |
| `--read-only` | - | false | 只读模式，禁止网页端修改配置       |
| `--login-protection` | - | false | 显式启用登录保护           |
| `--username` | - | - | 登录用户名                |
| `--password` | - | - | 登录密码                 |
| `--lang` | - | auto | CLI语言（auto/zh/en/ja） |
| `--debug` | - | false | 启用调试模式               |

<details>
<summary>更多参数（点击展开）</summary>

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--tls` | false | 启用 HTTPS |
| `--auto-tls` | false | 自动申请 Let's Encrypt 证书 |
| `--tls-crt` | - | TLS 证书文件路径 |
| `--tls-key` | - | TLS 密钥文件路径 |
| `--use-cache` | false | 启用本地图片缓存 |
| `--cache-dir` | - | 缓存目录路径 |
| `--cache-clean` | false | 退出时清除缓存 |
| `--database` | false | 启用本地数据库存储 |
| `--auto-rescan-min` | 0 | 自动扫描间隔（分钟，0为禁用） |
| `--min-image` | 1 | 最少图片数量才认定为漫画 |
| `--zip-encode` | gbk | 非UTF-8 ZIP文件的编码 |
| `--log-file` | false | 输出日志到文件 |
| `--print-all` | false | 打印所有网卡的访问地址 |
| `--single-instance` | false | 单实例模式 |
| `--plugin` | true | 启用插件系统 |
| `--tailscale` | false | 启用 Tailscale 网络 |
| `--tailscale-hostname` | comigo | Tailscale 主机名 |
| `--tailscale-funnel` | false | 启用 Tailscale Funnel |

</details>

### 使用示例

```bash
# 打开当前目录
comi .

# 指定端口和书库路径
comi -p 8080 /path/to/manga

# 仅本地访问，启用账号密码登录保护
comi --local --login-protection --username admin --password 123456 /path/to/manga

# 使用配置文件
comi -c /path/to/comigo.toml
```

## 配置文件说明

Comigo 支持多种配置文件位置：

1. **用户主目录**  
   - Windows: `C:\Users\用户名\.config\comigo.toml`
   - Linux/MacOS: `/home/用户名/.config/comigo.toml`
   - 程序启动时默认读取此位置

2. **程序目录**  
   - 将 `comigo.toml` 放在可执行文件同目录
   - 适合作为绿色软件使用

3. **当前运行目录**  
   - 在启动命令的当前目录下查找配置文件

4. **自定义位置**  
   - 使用 `--config` 参数指定配置文件路径

### 配置文件示例

```toml
# comigo.toml 配置示例

# 服务设置
Port = 1234                    # 服务端口
Host = ""                      # 自定义主机名（留空自动检测）
DisableLAN = false             # 仅本地访问
OpenBrowser = false            # 启动后打开浏览器
Language = "auto"              # 界面语言 (auto/zh/en/ja)

# 书库设置
StoreUrls = ["/path/to/manga", "/path/to/comics"]  # 书库路径列表
MaxScanDepth = 5               # 扫描深度
MinImageNum = 1                # 最少图片数量
AutoRescanIntervalMinutes = 0  # 自动扫描间隔（0为禁用）

# 登录保护
LoginProtection = false        # 显式启用登录保护
Username = ""                  # 本地账号密码登录的用户名
Password = ""                  # 密码
EnableOAuthLogin = false       # 启用 OAuth 登录
Timeout = 43200                # Cookie过期时间（分钟）

# 功能开关
EnableUpload = true            # 启用上传
ReadOnlyMode = false           # 只读模式
EnablePlugin = true            # 启用插件
Debug = false                  # 调试模式

# 缓存设置
UseCache = false               # 启用图片缓存
CacheDir = ""                  # 缓存目录（留空使用系统临时目录）
ClearCacheExit = false         # 退出时清除缓存

# ZIP文件设置
ZipFileTextEncoding = "gbk"    # 非UTF-8编码ZIP的解析编码
```

## 反馈与支持

如果您有任何建议或遇到问题，欢迎：
- 提交 [Issue](https://github.com/yumenaka/comigo/issues)
- 通过 [Twitter](https://x.com/yumenaka7) 联系我
- Discord讨论群 [Discord](https://discord.gg/c5q6d3dM8r)
## 开发与TODO
- [开发备忘](https://github.com/yumenaka/comigo/blob/master/TODO.md)

## 特别鸣谢

感谢以下开源项目及其贡献者：
- [mholt](https://github.com/mholt)
- [spf13](https://github.com/spf13)
- [disintegration](https://github.com/disintegration)
- [Baozisoftware](https://github.com/Baozisoftware)
- 以及更多贡献者

## 项目统计

[![Stargazers over time](https://starchart.cc/yumenaka/comigo.svg?variant=adaptive)](https://starchart.cc/yumenaka/comigo)

## 开源协议

本项目采用 MIT 协议开源。
