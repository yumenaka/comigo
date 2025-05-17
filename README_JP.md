<div align="center">

# ComiGo：シンプルで使いやすい漫画リーダー

[![Go Report](https://goreportcard.com/badge/github.com/yumenaka/comi?style=flat-square)](https://goreportcard.com/report/github.com/yumenaka/comi)
[![License](https://img.shields.io/github/license/yumenaka/comi?style=flat-square&color=blue)](https://github.com/yumenaka/comigo/blob/main/LICENSE)

</div>

![Windowsサンプル](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windowsサンプル")

[English](https://github.com/yumenaka/comigo/blob/master/README.md) | [中文文档](https://github.com/yumenaka/comigo/blob/master/README_CN.md) |  [日本語](https://github.com/yumenaka/comigo/blob/master/README_JP.md)

## 主な機能

- 📚 **多様なフォーマット対応**：画像フォルダ、`.rar`、`.zip`、`.tar`、`.cbz`、`.cbr`、`.epub` などの圧縮ファイルをサポート
- 🔄 **簡単アクセス**：スマートフォン/タブレットでのQRコードスキャン、Windowsでのドラッグ＆ドロップ対応
- 🐧 **マルチプラットフォーム**：Windows、Linux、MacOS に対応
- 📖 **多彩な閲覧モード**：スクロール、ページめくりなど閲覧モードを提供
- ⚙️ **柔軟な設定**：コマンドライン操作、`config.toml` 設定ファイルによるライブラリ設定
- 🖼️ **最新画像フォーマット**：`heic`、`avif` などの最新画像フォーマットをサポート
- ✂️ **スマート最適化**：画像の自動トリミング、トラフィック節約のための画像圧縮
- 🔄 **同期閲覧**：異なるデバイス間でのページめくり同期

## インストール方法

### ワンクリックインストール（推奨）

```bash
# curlを使用：
bash <(curl -s https://raw.githubusercontent.com/yumenaka/comigo/master/get_comigo.sh)

# wgetを使用：
bash <(wget -qO- https://raw.githubusercontent.com/yumenaka/comigo/master/get_comigo.sh)

# Golang環境（go 1.23以上）が設定済みの場合：
go install github.com/yumenaka/comigo/cmd/comi@latest
```

### 手動インストール

[Releases ページ](https://github.com/yumenaka/comigo/releases) から最新バージョンをダウンロードし、実行ファイルをシステムの `PATH` 環境変数に追加してください。

### バージョン選択ガイド

| システム            | ダウンロードバージョン         |
|-----------------|---------------------|
| Windows 64bit   | Windows_x86_64.zip  |
| Windows ARM版    | Windows_arm64.zip   |
| MacOS Appleチップ  | MacOS_arm64.tar.gz  |
| MacOS Intelチップ  | MacOS_x86_64.tar.gz |
| Linux 64bit     | Linux_x86_64.tar.gz |
| Linux ARM 32bit | Linux_arm.tar.gz    |
| Linux ARM 64bit | Linux_arm64.tar.gz  |

## 使用方法

```bash
comi [flags] file_or_dir
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
