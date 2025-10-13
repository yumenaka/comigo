<div align="center">

# ComiGo: Simple and Efficient Comic Reader

[![Go Report](https://goreportcard.com/badge/github.com/yumenaka/comi?style=flat-square)](https://goreportcard.com/report/github.com/yumenaka/comi)
[![License](https://img.shields.io/github/license/yumenaka/comi?style=flat-square&color=blue)](https://github.com/yumenaka/comigo/blob/main/LICENSE)

<!--
[![Downloads](https://img.shields.io/github/downloads/yumenaka/comi/total?style=flat-square&color=success)](https://github.com/yumenaka/comigo/releases)
<img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo：Simple Comig & Manga Reader" width="200">
-->
</div>

![Windows Sample](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windows Sample")

[中文文档](https://github.com/yumenaka/comigo/blob/master/README_CN.md) |[日本語](https://github.com/yumenaka/comigo/blob/master/README_JP.md) | [English](https://github.com/yumenaka/comigo/blob/master/README.md)

## Features

- 📚 **Multiple Format Support**: Supports image folders and compressed files like `.rar`, `.zip`, `.tar`, `.cbz`, `.cbr`, `.epub`
- 📱 **Easy Access**: QR code scanning for mobile/tablet devices, drag-and-drop support for Windows
- 🐧 **Cross-Platform**: Compatibility with Windows, Linux, and macOS
- 📖 **Diverse Reading Modes**: Offers scroll, and page-turning modes
- ⚙️ **Flexible Configuration**: Command-line operation with `config.toml` library settings
- 🖼️ **Modern Image Formats**: In addition to `jpg` and `png`, it also supports next-gen formats like `heic` and `avif`
- ✂️ **Smart Optimization**: Automatic image cropping and compression for bandwidth saving
- 🔄 **Sync Reading**: Synchronized page-turning across different devices

## Installation Guide

### Installation Script (Recommended)

```bash
# Using curl:
bash <(curl -s https://raw.githubusercontent.com/yumenaka/comigo/master/get_comigo.sh)

# Using wget:
bash <(wget -qO- https://raw.githubusercontent.com/yumenaka/comigo/master/get_comigo.sh)

# If you have Golang  (go 1.23 or higher):
go install github.com/yumenaka/comigo/cmd/comi@latest
```

### Manual Installation

Download the latest version from the [Releases page](https://github.com/yumenaka/comigo/releases) and add the executable to your system's `PATH` environment variable.

### Version Selection Guide

| System              | Download            |
|---------------------|---------------------|
| Windows 64-bit      | Windows_x86_64.zip  |
| Windows ARM         | Windows_arm64.zip   |
| MacOS Apple Silicon | MacOS_arm64.tar.gz  |
| MacOS Intel         | MacOS_x86_64.tar.gz |
| Linux 64-bit        | Linux_x86_64.tar.gz |
| Linux ARM 32-bit    | Linux_arm.tar.gz    |
| Linux ARM 64-bit    | Linux_arm64.tar.gz  |

## Usage

```bash
comi [flags] file_or_dir
```

## Configuration File

Comigo supports  configuration file locations:

1. **User Home Directory**  
   - Windows: `C:\Users\username\.config\comigo.toml`
   - Linux/MacOS: `/home/username/.config/comigo.toml`
   - Default location read at startup

2. **Program Directory**  
   - Place `comigo.toml` in the same directory as the executable
   - Suitable for portable usage

3. **Current Working Directory**  
   - Searches for configuration file in the current directory when running commands

4. **Custom Location**  
   - Specify configuration file path using the `--config` parameter

## Feedback and Support

If you have any suggestions or encounter issues, feel free to:
- Submit an [Issue](https://github.com/yumenaka/comigo/issues)
- Contact me via [Twitter](https://x.com/yumenaka7)
- Join the discussion on [Discord](https://discord.gg/brBfSExJPn)
## Special Thanks

Thanks to the following open-source projects and their contributors:
- [mholt](https://github.com/mholt)
- [spf13](https://github.com/spf13)
- [disintegration](https://github.com/disintegration)
- [Baozisoftware](https://github.com/Baozisoftware)
- And many more contributors

## Project Statistics

[![Stargazers over time](https://starchart.cc/yumenaka/comigo.svg?variant=adaptive)](https://starchart.cc/yumenaka/comigo)

## License

This software is released under the MIT license.
