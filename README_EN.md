<div align="center">

# ComiGo: Simple and Efficient Comic Reader

[![Go Report](https://goreportcard.com/badge/github.com/yumenaka/comi?style=flat-square)](https://goreportcard.com/report/github.com/yumenaka/comi)
[![License](https://img.shields.io/github/license/yumenaka/comi?style=flat-square&color=blue)](https://github.com/yumenaka/comigo/blob/main/LICENSE)

[中文文档](https://github.com/yumenaka/comigo/blob/master/README.md) |[日本語](https://github.com/yumenaka/comigo/blob/master/README_JP.md) | [English](https://github.com/yumenaka/comigo/blob/master/README_EN.md)

<!--
[![Downloads](https://img.shields.io/github/downloads/yumenaka/comi/total?style=flat-square&color=success)](https://github.com/yumenaka/comigo/releases)
<img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo：Simple Comig & Manga Reader" width="200">
-->
</div>

![Windows Sample](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windows Sample")


## Features

- 📚 **Multiple Format Support**: Supports image folders and compressed files like `.rar`, `.zip`, `.tar`, `.cbz`, `.cbr`, `.epub`
- 📱 **Easy Access**: QR code scanning for mobile/tablet devices, drag-and-drop support for Windows
- 🐧 **Cross-Platform**: Compatibility with Windows, Linux, and macOS
- 📖 **Diverse Reading Modes**: Offers scroll, and page-turning modes
- ⚙️ **Flexible Configuration**: Command-line operation with `config.toml` library settings
- 🖼️ **Modern Image Formats**: In addition to `jpg` and `png`, it also supports next-gen formats like `heic` and `avif`
- ✂️ **Smart Optimization**: Automatic image cropping and compression for bandwidth saving
- 🔄 **Sync Reading**: Synchronized page-turning across different devices
- 🔌 **Plugin System**: Built-in plugins like auto page-turn and clock, with custom plugin support
- 🎬 **Media Playback**: Built-in audio and video player
- 📥 **Flexible Download**: Batch download image folders, convert and download as EPUB format
- 📜 **Reading History**: Automatic reading history tracking for easy continuation

## Installation Guide

### GUI Version (Recommended for Beginners)

| System | Download |
|--------|----------|
| Windows 64-bit | [comigo_latest_Windows_x86_64_full.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comigo_latest_Windows_x86_64_full.zip) |
| macOS (Intel/Apple Silicon) | [Comigo.app.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/Comigo.app.zip) |

> 💡 **Note**: GUI version provides system tray icon, can run minimized in background. Windows: Double-click to run; macOS: Drag to Applications folder.

### Quick Install for CLI

```bash
# Recommended:
bash <(curl -s https://raw.githubusercontent.com/yumenaka/comigo/master/get.sh)

# For users in Mainland China:
bash <(curl -s https://comigo.xyz/get.sh) --cn

# If you have Golang (go 1.23 or higher):
go install github.com/yumenaka/comigo/cmd/comi@latest
```

### CLI Version

| System | Download |
|--------|----------|
| Windows 64-bit | [comi_latest_Windows_x86_64.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Windows_x86_64.zip) |
| Windows ARM | [comi_latest_Windows_arm64.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Windows_arm64.zip) |
| macOS Intel | [comi_latest_MacOS_x86_64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_MacOS_x86_64.tar.gz) |
| macOS Apple Silicon | [comi_latest_MacOS_arm64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_MacOS_arm64.tar.gz) |
| Linux 64-bit | [comi_latest_Linux_x86_64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Linux_x86_64.tar.gz) |
| Linux ARM64 | [comi_latest_Linux_arm64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Linux_arm64.tar.gz) |
| Linux ARM32 | [comi_latest_Linux_armv7.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Linux_armv7.tar.gz) |
| Debian/Ubuntu 64-bit | [comi_latest_amd64.deb](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_amd64.deb) |
| Debian/Ubuntu ARM64 | [comi_latest_arm64.deb](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_arm64.deb) |

> 💡 **Note**: CLI version suitable for server deployment and advanced users. Manual PATH configuration required after download.

### Manual Installation

Download the latest version from the [Releases page](https://github.com/yumenaka/comigo/releases) and add the executable to your system's `PATH` environment variable.

## Docker Deployment

### Quick Start with Docker

```bash
# Pull and run the latest image
docker run -d \
  --name comigo \
  -p 1234:1234 \
  -v /path/to/your/books:/data \
  yumenaka/comigo:latest
```

Visit `http://localhost:1234` to access your library.

### Using Docker Compose

1. Download the [`docker-compose.yml`](sample/docker/docker-compose.yml) file
2. Edit the configuration as needed
3. Start the service:

```bash
docker-compose up -d
```

### Supported Platforms

- `linux/amd64` - Standard x86_64 servers
- `linux/arm64` - ARM64 servers (Raspberry Pi 4/5)
- `linux/arm/v7` - ARMv7 devices (Raspberry Pi 2-4)

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `COMIGO_PORT` | Service port | `1234` |
| `COMIGO_USERNAME` | Login username (optional) | - |
| `COMIGO_PASSWORD` | Login password (optional) | - |
| `COMIGO_ENABLE_UPLOAD` | Enable file upload | `true` |

For more details, see the complete [Docker documentation](sample/docker/README.md).

## Usage

```bash
comi [flags] file_or_dir
```

### Command Line Options

| Option | Short | Default | Description |
|--------|-------|---------|-------------|
| `--config` | `-c` | - | Specify config file path |
| `--port` | `-p` | 1234 | Service port |
| `--host` | - | - | Custom hostname/domain |
| `--local` | - | false | Local access only |
| `--max-depth` | `-m` | 5 | Max scan depth |
| `--open-browser` | `-o` | false | Open browser on start |
| `--enable-upload` | - | true | Enable file upload |
| `--read-only` | - | false | Read-only mode |
| `--username` | - | - | Login username |
| `--password` | - | - | Login password |
| `--lang` | - | auto | Language (auto/zh/en/ja) |
| `--debug` | - | false | Debug mode |

### Examples

```bash
# Open current directory
comi .

# Specify port and path
comi -p 8080 /path/to/manga

# Local only with login protection
comi --local --username admin --password 123456 /path/to/manga
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
- Join the discussion on [Discord](https://discord.gg/c5q6d3dM8r)
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
