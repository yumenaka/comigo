<div align="center">

# ComiGo：简单高效的漫画阅读器

[![Go Report](https://goreportcard.com/badge/github.com/yumenaka/comi?style=flat-square)](https://goreportcard.com/report/github.com/yumenaka/comi)
[![License](https://img.shields.io/github/license/yumenaka/comi?style=flat-square&color=blue)](https://github.com/yumenaka/comigo/blob/main/LICENSE)

<!--
[![Downloads](https://img.shields.io/github/downloads/yumenaka/comi/total?style=flat-square&color=success)](https://github.com/yumenaka/comigo/releases)
<img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo：简单方便的漫画阅读器" width="200">
-->
</div>

![Windows示例](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windows示例")

[English](https://github.com/yumenaka/comigo/blob/master/README_EN.md) | [日本語](https://github.com/yumenaka/comigo/blob/master/README_JP.md) | [中文文档](https://github.com/yumenaka/comigo/blob/master/README.md) 

## 功能特点

- 📚 **多格式支持**：支持图片文件夹与 `.rar`、`.zip`、`.tar`、`.cbz`、`.cbr`、`.epub` 等压缩包格式
- 📱 **便捷访问**：支持手机/平板扫描二维码阅读，Windows 支持拖拽打开
- 🐧 **跨平台支持**：适配 Windows、Linux、MacOS 系统
- 📖 **多样化阅读模式**：提供卷轴、翻页等多种阅读模式
- ⚙️ **灵活配置**：支持命令行操作，可通过 `config.toml` 配置文件设定书库
- 🖼️ **现代图片格式**：除了常见的`jpg`、`png`，还支持 `heic`、`avif` 等新一代图片格式
- ✂️ **智能优化**：支持图片自动裁边，压缩图片节省流量
- 🔄 **同步阅读**：支持不同设备间同步翻页进度

## 安装指南

### 托盘版（推荐）

| 系统 | 下载链接 |
|------|---------|
| Windows 64位 | [comigo_latest_Windows_x86_64_full.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comigo_latest_Windows_x86_64_full.zip) |
| macOS (Intel/Apple芯片) | [Comigo.app.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/Comigo.app.zip) |

> 💡 **说明**：托盘版提供系统托盘图标，可最小化到后台运行。Windows 用户双击运行，macOS 用户拖入 Applications 文件夹。

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

## 反馈与支持

如果您有任何建议或遇到问题，欢迎：
- 提交 [Issue](https://github.com/yumenaka/comigo/issues)
- 通过 [Twitter](https://x.com/yumenaka7) 联系我
- Discord讨论群 [Discord](https://discord.gg/c5q6d3dM8r)
## 开发与TODO
- [开发备忘](https://github.com/yumenaka/comigo/blob/master/todo.md)

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
