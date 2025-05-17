<div align="center">

# ComiGo：简单高效的漫画阅读器

[![Go Report](https://goreportcard.com/badge/github.com/yumenaka/comi?style=flat-square)](https://goreportcard.com/report/github.com/yumenaka/comi)
[![License](https://img.shields.io/github/license/yumenaka/comi?style=flat-square&color=blue)](https://github.com/yumenaka/comigo/blob/main/LICENSE)

<!--
[![Downloads](https://img.shields.io/github/downloads/yumenaka/comi/total?style=flat-square&color=success)](https://github.com/yumenaka/comigo/releases)
<img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo：简单粗暴的漫画阅读器" width="200">
-->
</div>

![Windows示例](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windows示例")

[English](https://github.com/yumenaka/comigo/blob/master/README.md) | [日本語](https://github.com/yumenaka/comigo/blob/master/README_JP.md) | [中文文档](https://github.com/yumenaka/comigo/blob/master/README_CN.md) 

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

### 一键安装（推荐）

```bash
# 使用 curl：
bash <(curl -s https://raw.githubusercontent.com/yumenaka/comigo/master/get_comigo.sh)

# 使用 wget：
bash <(wget -qO- https://raw.githubusercontent.com/yumenaka/comigo/master/get_comigo.sh)

# 如果您已设置 Golang 环境（go 1.23 或更高版本）：
go install github.com/yumenaka/comigo/cmd/comi@latest
```

### 手动安装

在 [Releases 页面](https://github.com/yumenaka/comigo/releases) 下载最新版本，并将可执行文件添加到系统的 `PATH` 环境变量中。

### 版本选择指南

| 系统类型          | 下载版本                |
|---------------|---------------------|
| Windows 64位   | Windows_x86_64.zip  |
| Windows ARM版  | Windows_arm64.zip   |
| MacOS Apple芯片 | MacOS_arm64.tar.gz  |
| MacOS Intel芯片 | MacOS_x86_64.tar.gz |
| Linux 64位     | Linux_x86_64.tar.gz |
| Linux ARM 32位 | Linux_arm.tar.gz    |
| Linux ARM 64位 | Linux_arm64.tar.gz  |

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

## 开发与TODO
- [开发备忘](https://github.com/yumenaka/comigo/blob/master/Develop_Todo.md)

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
