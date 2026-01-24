<div align="center">

# ComiGo：シンプルで使いやすい漫画リーダー  
[![Go Report](https://goreportcard.com/badge/github.com/yumenaka/comi?style=flat-square)](https://goreportcard.com/report/github.com/yumenaka/comi)
[![License](https://img.shields.io/github/license/yumenaka/comi?style=flat-square&color=blue)](https://github.com/yumenaka/comigo/blob/main/LICENSE)  
[中文文档](https://github.com/yumenaka/comigo/blob/master/README.md) | [English](https://github.com/yumenaka/comigo/blob/master/README_EN.md) | [日本語](https://github.com/yumenaka/comigo/blob/master/README_JP.md)  
![Windowsサンプル](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windowsサンプル") 
<!--
[![Downloads](https://img.shields.io/github/downloads/yumenaka/comi/total?style=flat-square&color=success)](https://github.com/yumenaka/comigo/releases)
<img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo：Simple Comig & Manga Reader" width="200">
-->
</div>



## 主な機能

- 📚 **多様なフォーマット対応**：画像フォルダ、`.rar`、`.zip`、`.tar`、`.cbz`、`.cbr`、`.epub` などの圧縮ファイルをサポート
- 🔄 **簡単アクセス**：スマートフォン/タブレットでのQRコードスキャン、Windowsでのドラッグ＆ドロップ対応
- 🐧 **マルチプラットフォーム**：Windows、Linux、MacOS に対応
- 📖 **多彩な閲覧モード**：スクロール、ページめくりなど閲覧モードを提供
- ⚙️ **柔軟な設定**：コマンドライン操作、`config.toml` 設定ファイルによるライブラリ設定
- 🖼️ **最新画像フォーマット**：jpg や png に加えて、heic や avif などの次世代画像フォーマットにも対応しています。
- ✂️ **スマート最適化**：画像の自動トリミング、トラフィック節約のための画像圧縮
- 🔄 **同期閲覧**：異なるデバイス間でのページめくり同期
- 🔌 **プラグインシステム**：自動ページめくり、時計などの内蔵プラグイン、カスタムプラグインの拡張に対応
- 🎬 **メディア再生**：内蔵オーディオ・ビデオプレーヤー
- 📥 **柔軟なダウンロード**：画像フォルダの一括ダウンロード、EPUBフォーマットへの変換とダウンロードに対応
- 📜 **閲覧履歴**：閲覧履歴を自動記録、続きから読める

## インストール方法

### GUI版（初心者におすすめ）

| システム | ダウンロード |
|----------|-------------|
| Windows 64bit | [comigo_latest_Windows_x86_64_full.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comigo_latest_Windows_x86_64_full.zip) |
| macOS (Intel/Apple Silicon) | [Comigo.app.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/Comigo.app.zip) |

> 💡 **説明**：GUI版はシステムトレイアイコンを提供し、バックグラウンドで実行できます。Windows: ダブルクリックで実行; macOS: Applicationsフォルダにドラッグ。

### CLI版のワンクリックインストール

```bash
# 推奨：
bash <(curl -s https://raw.githubusercontent.com/yumenaka/comigo/master/get.sh)

# 中国本土のユーザー向け：
bash <(curl -s https://comigo.xyz/get.sh) --cn

# Golang環境（go 1.23以上）が設定済みの場合：
go install github.com/yumenaka/comigo/cmd/comi@latest
```

### CLI版

| システム | ダウンロード |
|---------|-------------|
| Windows 64bit | [comi_latest_Windows_x86_64.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Windows_x86_64.zip) |
| Windows ARM | [comi_latest_Windows_arm64.zip](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Windows_arm64.zip) |
| macOS Intel | [comi_latest_MacOS_x86_64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_MacOS_x86_64.tar.gz) |
| macOS Apple Silicon | [comi_latest_MacOS_arm64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_MacOS_arm64.tar.gz) |
| Linux 64bit | [comi_latest_Linux_x86_64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Linux_x86_64.tar.gz) |
| Linux ARM64 | [comi_latest_Linux_arm64.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Linux_arm64.tar.gz) |
| Linux ARM32 | [comi_latest_Linux_armv7.tar.gz](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_Linux_armv7.tar.gz) |
| Debian/Ubuntu 64bit | [comi_latest_amd64.deb](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_amd64.deb) |
| Debian/Ubuntu ARM64 | [comi_latest_arm64.deb](https://comigo.xyz/yumenaka/comigo/releases/download/latest/comi_latest_arm64.deb) |

> 💡 **説明**：CLI版はサーバーデプロイと上級ユーザーに適しています。ダウンロード後、手動でシステムPATHに追加する必要があります。

### 手動インストール

[Releases ページ](https://github.com/yumenaka/comigo/releases) から最新バージョンをダウンロードし、実行ファイルをシステムの `PATH` 環境変数に追加してください。

## Docker デプロイ

### クイックスタート

```bash
# 最新イメージをプルして実行
docker run -d \
  --name comigo \
  -p 1234:1234 \
  -v /path/to/your/books:/data \
  yumenaka/comigo:latest
```

`http://localhost:1234` にアクセスして使用を開始します。

### Docker Compose を使用

1. [`docker-compose.yml`](sample/docker/docker-compose.yml) ファイルをダウンロード
2. 必要に応じて設定を編集
3. サービスを開始：

```bash
docker-compose up -d
```

### サポートされているプラットフォーム

- `linux/amd64` - 標準的な x86_64 サーバー
- `linux/arm64` - ARM64 サーバー（Raspberry Pi 4/5）
- `linux/arm/v7` - ARMv7 デバイス（Raspberry Pi 2-4）

### 環境変数

| 変数名 | 説明 | デフォルト値 |
|--------|------|-------------|
| `COMIGO_PORT` | サービスポート | `1234` |
| `COMIGO_USERNAME` | ログインユーザー名（オプション） | - |
| `COMIGO_PASSWORD` | ログインパスワード（オプション） | - |
| `COMIGO_ENABLE_UPLOAD` | ファイルアップロードを有効化 | `true` |

詳細については、完全な [Docker ドキュメント](sample/docker/README.md) をご覧ください。

## 使用方法

```bash
comi [flags] file_or_dir
```

### コマンドラインオプション

| オプション | 短縮形 | デフォルト | 説明 |
|------------|--------|-----------|------|
| `--config` | `-c` | - | 設定ファイルのパス |
| `--port` | `-p` | 1234 | サービスポート |
| `--host` | - | - | カスタムホスト名 |
| `--local` | - | false | ローカルアクセスのみ |
| `--max-depth` | `-m` | 5 | 最大スキャン深度 |
| `--open-browser` | `-o` | false | 起動時にブラウザを開く |
| `--enable-upload` | - | true | アップロード機能を有効化 |
| `--read-only` | - | false | 読み取り専用モード |
| `--username` | - | - | ログインユーザー名 |
| `--password` | - | - | ログインパスワード |
| `--lang` | - | auto | 言語設定（auto/zh/en/ja） |
| `--debug` | - | false | デバッグモード |

### 使用例

```bash
# カレントディレクトリを開く
comi .

# ポートとパスを指定
comi -p 8080 /path/to/manga

# ローカルのみ、ログイン保護付き
comi --local --username admin --password 123456 /path/to/manga
```

## 設定ファイルについて

Comigo は複数の設定ファイルの場所をサポートしています：

1. **ユーザーホームディレクトリ**  
   - Windows: `C:\Users\ユーザー名\.config\comigo.toml`
   - Linux/MacOS: `/home/ユーザー名/.config/comigo.toml`
   - プログラム起動時にデフォルトで読み込まれます

2. **プログラムディレクトリ**  
   - 実行ファイルと同じディレクトリに `comigo.toml` を配置
   - ポータブルアプリケーションとして使用する場合に適しています

3. **現在の実行ディレクトリ**  
   - コマンド実行時のカレントディレクトリで設定ファイルを検索

4. **カスタムロケーション**  
   - `--config` パラメータで設定ファイルのパスを指定可能

## フィードバックとサポート

ご意見や問題がございましたら、以下からお気軽にご連絡ください：
- [Issue](https://github.com/yumenaka/comigo/issues) を投稿
- [Twitter](https://x.com/yumenaka7) でメッセージを送信
- [Discord](https://discord.gg/c5q6d3dM8r) でディスカッションに参加
## 特別な感謝

以下のオープンソースプロジェクトとその貢献者に感謝いたします：
- [mholt](https://github.com/mholt)
- [spf13](https://github.com/spf13)
- [disintegration](https://github.com/disintegration)
- [Baozisoftware](https://github.com/Baozisoftware)
- その他の貢献者の皆様

## プロジェクト統計

[![Stargazers over time](https://starchart.cc/yumenaka/comigo.svg?variant=adaptive)](https://starchart.cc/yumenaka/comigo)

## ライセンス

このソフトウェアは MIT ライセンスの下で公開されています。
