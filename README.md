<div align="center">

# ComiGo: Simple and Efficient Comic Reader  
[![Go Report](https://goreportcard.com/badge/github.com/yumenaka/comi?style=flat-square)](https://goreportcard.com/report/github.com/yumenaka/comi)
[![License](https://img.shields.io/github/license/yumenaka/comi?style=flat-square&color=blue)](https://github.com/yumenaka/comigo/blob/main/LICENSE)  
[中文文档](https://github.com/yumenaka/comigo/blob/master/README_ZH.md) | [日本語](https://github.com/yumenaka/comigo/blob/master/README_JP.md) | [English](https://github.com/yumenaka/comigo/blob/master/README.md)
<!--
[![Downloads](https://img.shields.io/github/downloads/yumenaka/comi/total?style=flat-square&color=success)](https://github.com/yumenaka/comigo/releases)
<img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo：Simple Comig & Manga Reader" width="200">
-->
</div>

![Windows Sample](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windows Sample")

## [Try it Online](https://comigo.xyz/reader)

## Features

- 📚 **Multiple Format Support**: Supports image folders and compressed files like `.rar`, `.zip`, `.tar`, `.cbz`, `.cbr`, `.epub`
- 📱 **Easy Access**: QR code scanning for mobile/tablet devices, drag-and-drop support for Windows
- 📴 **Offline Mode / PWA App**: Read ZIP/CBZ/RAR/CBR directly in your browser without uploading files
- 🐧 **Cross-Platform**: Compatibility with Windows, Linux, and macOS
- 📖 **Diverse Reading Modes**: Offers scroll, and page-turning modes
- ⚙️ **Flexible Configuration**: Command-line operation with `config.toml` library settings
- 🖼️ **Modern Image Formats**: In addition to `jpg` and `png`, it also supports next-gen formats like `heic` and `avif`
- ✂️ **Smart Optimization**: Automatic image cropping and compression for bandwidth saving
- 🔄 **Sync Reading**: Synchronized page-turning across different devices
- 🔌 **Plugin System**: Built-in plugins like auto page-turn and clock, with custom plugin support
- 🎬 **Media Playback**: Built-in audio and video player
- 📥 **Flexible Use**: Supports remote Comigo libraries, batch image-folder downloads, and EPUB conversion
- 📜 **Reading History**: Automatic reading history tracking for easy continuation

PC/Mobile Sync:   
![Mobile Sync Sample](https://www.yumenaka.net/wp-content/uploads/2026/04/scroll.gif "Mobile Sync Sample")

## Omarchy status bar integration

Comigo supports Omarchy through the [comigo-omarchy plugin](https://github.com/yumenaka/comigo-omarchy), which controls local or remote services from the status bar.

```bash
omarchy plugin add https://github.com/yumenaka/comigo-omarchy --enable
```

## Installation Guide

### GUI Version (Recommended for Beginners)

| Version | System | Download |
|---------|--------|----------|
| Tray | Windows 64-bit | [comigo-tray_latest_Windows_x86_64.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comigo-tray_latest_Windows_x86_64.zip) |
| Tray | Windows ARM64 | [comigo-tray_latest_Windows_arm64.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comigo-tray_latest_Windows_arm64.zip) |
| Tray | macOS (Intel/Apple Silicon) | [comigo-tray_latest_MacOS_universal.dmg](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comigo-tray_latest_MacOS_universal.dmg) |
| Tray | Linux 64-bit | [comigo-tray_latest_Linux_x86_64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comigo-tray_latest_Linux_x86_64.tar.gz) |
| Tray | Linux ARM64 | [comigo-tray_latest_Linux_arm64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comigo-tray_latest_Linux_arm64.tar.gz) |
| Desktop | Windows 64-bit | [comigo-desktop_latest_Windows_x86_64.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comigo-desktop_latest_Windows_x86_64.zip) |
| Desktop | Windows ARM64 | [comigo-desktop_latest_Windows_arm64.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comigo-desktop_latest_Windows_arm64.zip) |
| Desktop | macOS (Intel/Apple Silicon) | [comigo-desktop_latest_MacOS_universal.dmg](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comigo-desktop_latest_MacOS_universal.dmg) |
| Desktop | Linux 64-bit | [comigo-desktop_latest_Linux_x86_64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comigo-desktop_latest_Linux_x86_64.tar.gz) |
| Desktop | Linux ARM64 | [comigo-desktop_latest_Linux_arm64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comigo-desktop_latest_Linux_arm64.tar.gz) |

> 💡 **Which GUI package should I choose?** Choose Tray if you want ComiGo to keep running in the background and manage it from the system tray. The Tray version starts the local Web service and opens the reader in your browser; the tray menu can open the browser, copy the reading URL, open config/library folders, switch languages, check for updates, and toggle Tailscale remote access. On Windows, it can also register folder context menus, file associations, and desktop shortcuts.
>
> Choose Desktop if you prefer a regular app window. It wraps the same ComiGo reader in a Wails desktop shell and adds desktop-only actions such as selecting folders with the system picker and moving local source files to the system trash after confirmation. The current Desktop build does not include the tray menu or built-in Tailscale remote access.

> **macOS tip**: If macOS says the app is damaged and should be moved to Trash, the downloaded file is usually not actually broken. Run this after moving the app to Applications:
> ```bash
> xattr -dr com.apple.quarantine /Applications/comigo-desktop.app
> ```

### Quick Install for CLI

```bash
# Recommended:
installer=$(curl -fsSL https://raw.githubusercontent.com/yumenaka/comigo/master/get.sh) && bash -c "$installer" --

# For users in Mainland China:
installer=$(curl -fsSL https://comigo.xyz/get.sh) && bash -c "$installer" -- --cn

# If you have Golang (go 1.23 or higher):
go install github.com/yumenaka/comigo/cmd/comi@latest
```

The installer defaults to `$HOME/.local/bin` and does not edit shell configuration. It checks both the available `comi` command and the destination file. If the requested version is already installed, a terminal prompt offers skip (default) or overwrite at the displayed destination; without a terminal it skips. Use `--force` to overwrite the same version non-interactively.

For the examples below, first download the script with `curl -fSL -o get.sh https://raw.githubusercontent.com/yumenaka/comigo/master/get.sh`.

```bash
bash get.sh --help
bash get.sh --version v1.2.22
bash get.sh --system                 # /usr/local/bin; root-owned, sudo when needed
bash get.sh --install-dir "$HOME/bin" # overrides COMIGO_INSTALL_DIR
bash get.sh --force                  # overwrite the same version
```

The destination remains explicit: an existing command elsewhere is not automatically replaced. Use `--system` or `--install-dir` to select its directory. The installer reports if PATH selects another copy and prints shell-specific PATH instructions. Explicit directories never silently fall back; symlink targets are rejected. To uninstall, remove only `comi` from the chosen directory (use the package manager if it manages that copy).

Downloads require the release asset `checksums.txt`; for older releases without it, explicitly add `--skip-checksum` to disable verification. `make all` generates the SHA-256 manifest after building all packages; publish it alongside the release assets. To regenerate it separately, run `make checksums VERSION=vX.Y.Z`.

### CLI Version

| System | Download |
|--------|----------|
| Windows 64-bit | [comi_latest_Windows_x86_64.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Windows_x86_64.zip) |
| Windows ARM64 | [comi_latest_Windows_arm64.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Windows_arm64.zip) |
| Windows 32-bit | [comi_latest_Windows_i386.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Windows_i386.zip) |
| macOS Intel | [comi_latest_MacOS_x86_64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_MacOS_x86_64.tar.gz) |
| macOS Apple Silicon | [comi_latest_MacOS_arm64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_MacOS_arm64.tar.gz) |
| Linux 64-bit | [comi_latest_Linux_x86_64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Linux_x86_64.tar.gz) |
| Linux ARM64 | [comi_latest_Linux_arm64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Linux_arm64.tar.gz) |
| Linux ARM32 | [comi_latest_Linux_armv7.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Linux_armv7.tar.gz) |
| Linux 32-bit | [comi_latest_Linux_i386.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Linux_i386.tar.gz) |
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

1. Download the [`docker-compose.yml`](docs/docker/docker-compose.yml) file
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

For more details, see the complete [Docker documentation](docs/docker/README.md).

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
| `--debug` | - | false | Debug mode |

### Examples

```bash
# Open current directory
comi .

# Specify port and path
comi -p 8080 /path/to/manga

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

## Development with mise

The committed [`mise.toml`](mise.toml) pins Go, Bun, GNU Make, templ, Wails v2, Air, sqlc, Ent, and goversioninfo. It uses mise's documented core, `go:`, and `conda:` backends; GNU Make is provided by `conda:make`, whose mise backend supports macOS, Linux, and Windows x64 without a separate Conda installation.

The initialization and task flow below was verified on macOS ARM64 with mise 2026.8.14. The Linux and Windows commands and platform limits come from the linked official mise and Wails documentation.

### Install and activate mise

Follow the official [mise Getting Started](https://mise.jdx.dev/getting-started.html) and [installation](https://mise.jdx.dev/installing-mise.html) documentation.

macOS and Linux:

```bash
curl https://mise.run | sh

# Bash
echo 'eval "$(~/.local/bin/mise activate bash)"' >> ~/.bashrc

# Zsh (use this instead on the default macOS shell)
echo 'eval "$(~/.local/bin/mise activate zsh)"' >> ~/.zshrc
```

Restart the shell after adding one activation line. Homebrew is also supported with `brew install mise`, though the official installer is mise's recommended method.

Windows PowerShell:

```powershell
# Recommended; Scoop also adds mise shims to PATH.
scoop install mise

# Alternative
winget install jdx.mise
```

PowerShell activation is optional when the Scoop shims are available. To enable automatic environment switching, add the official activation command to your actual PowerShell profile:

```powershell
New-Item -ItemType Directory -Force (Split-Path $PROFILE)
Add-Content $PROFILE '(&mise activate pwsh) | Out-String | Invoke-Expression'
```

### Initialize the repository

The commands are the same in Bash, Zsh, and PowerShell:

```text
git clone https://github.com/yumenaka/comigo.git
cd comigo
mise trust
mise install
mise run setup
mise doctor
```

`mise install` installs the pinned tools. `mise run setup` downloads Go modules and runs `bun install`; mise tasks also install any missing configured tools automatically. Run `wails doctor` once before desktop development. Wails still needs native platform components that mise cannot supply: Xcode Command Line Tools on macOS, WebView2 on Windows, or GCC, GTK3, and WebKitGTK development packages on Linux. See the official [Wails v2 installation guide](https://wails.io/docs/gettingstarted/installation/).

On Windows x64, normal mise tasks work from PowerShell. The release-oriented Make targets additionally use POSIX utilities such as `sh`, `cp`, `rm`, `zip`, and `find`, so run those targets in MSYS2 with its toolchain on `PATH`. The mise Conda backend does not currently list Windows ARM64; on that host, install Make through MSYS2 and run the release targets there. Docker, platform SDKs, `dpkg`, UPX, and similar release dependencies remain system-managed.

### Common tasks and aliases

Run `mise tasks` to see the full list. Task aliases use the documented `mise run <alias>` form.

| Task | Alias | Purpose |
|------|-------|---------|
| `mise run setup` | `mise run i` | Install Go modules and frontend packages |
| `mise run generate` | `mise run g` / `mise run gen` | Format and generate templ output |
| `mise run frontend` | `mise run fe` | Build unminified frontend assets |
| `mise run dev` | `mise run d` | Start Air hot reload |
| `mise run run` | `mise run r` | Generate templates and run the Web server |
| `mise run test` | `mise run t` | Run all Go tests |
| `mise run sqlc` | `mise run db` | Regenerate sqlc output |
| `mise run wails-dev` | `mise run wd` | Start Wails v2 development mode |
| `mise run wails-build` | `mise run wb` | Build the current platform's Wails app |
| `mise run release` | `mise run rel` | Run the existing `make all` release matrix |
| `mise run clean` | `mise run c` | Clean Make build output |

Direct Make targets remain available inside the pinned environment, for example `mise exec -- make wails-prepare` or `mise exec -- make docker-help`.

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

## Star History

<a href="https://www.star-history.com/?repos=yumenaka%2Fcomigo&type=date&legend=top-left">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/chart?repos=yumenaka/comigo&type=date&theme=dark&legend=top-left&sealed_token=HWOAmg0U031iJgsEmeThimLrmGTYn9bIb0JtHUapnADRBmC_1ISYM46YGi28oBQ_f9q9VjtIlkis8AOOpPYxvE6SabcMQ-JQlEBlFtkE7BfBkJK7hWhJrLD2wE14VtpOifa54t5eIEJmN5OXz6HYbI9v0Dcxm01wSINkyz1DVjuccNQdQ9ljt5YRxF2f" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/chart?repos=yumenaka/comigo&type=date&legend=top-left&sealed_token=HWOAmg0U031iJgsEmeThimLrmGTYn9bIb0JtHUapnADRBmC_1ISYM46YGi28oBQ_f9q9VjtIlkis8AOOpPYxvE6SabcMQ-JQlEBlFtkE7BfBkJK7hWhJrLD2wE14VtpOifa54t5eIEJmN5OXz6HYbI9v0Dcxm01wSINkyz1DVjuccNQdQ9ljt5YRxF2f" />
   <img alt="Star History Chart" src="https://api.star-history.com/chart?repos=yumenaka/comigo&type=date&legend=top-left&sealed_token=HWOAmg0U031iJgsEmeThimLrmGTYn9bIb0JtHUapnADRBmC_1ISYM46YGi28oBQ_f9q9VjtIlkis8AOOpPYxvE6SabcMQ-JQlEBlFtkE7BfBkJK7hWhJrLD2wE14VtpOifa54t5eIEJmN5OXz6HYbI9v0Dcxm01wSINkyz1DVjuccNQdQ9ljt5YRxF2f" />
 </picture>
</a>

## License

This software is released under the MIT license.
