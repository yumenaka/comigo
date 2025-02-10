#!/bin/bash

# 一键安装脚本
# 使用 curl：
#   bash <(curl -s https://raw.githubusercontent.com/yumenaka/comigo/master/get_comigo.sh)
# 使用 wget：
#   bash <(wget -qO- https://raw.githubusercontent.com/yumenaka/comigo/master/get_comigo.sh)

# 遇到错误时立即退出
set -e

# 获取系统语言
system_language=$(locale | grep -E '^LANG=' | cut -d= -f2 | cut -d. -f1)

# 定义消息输出函数
function print_message() {
    local key="$1"
    local param="$2"
    case "$key" in
        "error_no_curl_wget")
            case "$system_language" in
                zh_CN)
                    echo "错误：需要安装 curl 或 wget"
                    ;;
                en_US)
                    echo "Error: curl or wget is required."
                    ;;
                ja_JP)
                    echo "エラー：curlまたはwgetが必要です。"
                    ;;
                *)
                    echo "Error: curl or wget is required."
                    ;;
            esac
            ;;
        "error_cmd_not_found")
            case "$system_language" in
                zh_CN)
                    echo "错误：未找到 $param，请先安装 $param"
                    ;;
                en_US)
                    echo "Error: $param not found. Please install $param."
                    ;;
                ja_JP)
                    echo "エラー：$param が見つかりません。$param をインストールしてください。"
                    ;;
                *)
                    echo "Error: $param not found. Please install $param."
                    ;;
            esac
            ;;
        "error_cannot_get_latest_tag")
            case "$system_language" in
                zh_CN)
                    echo "错误：无法获取最新版本标签。"
                    ;;
                en_US)
                    echo "Error: Unable to fetch the latest version tag."
                    ;;
                ja_JP)
                    echo "エラー：最新のバージョンタグを取得できません。"
                    ;;
                *)
                    echo "Error: Unable to fetch the latest version tag."
                    ;;
            esac
            ;;
        "error_cannot_parse_version")
            case "$system_language" in
                zh_CN)
                    echo "错误：无法从标签 $param 解析版本号。"
                    ;;
                en_US)
                    echo "Error: Unable to parse version number from tag $param."
                    ;;
                ja_JP)
                    echo "エラー：タグ $param からバージョン番号を解析できません。"
                    ;;
                *)
                    echo "Error: Unable to parse version number from tag $param."
                    ;;
            esac
            ;;
        "error_unsupported_os")
            case "$system_language" in
                zh_CN)
                    echo "错误：不支持的操作系统：$param"
                    ;;
                en_US)
                    echo "Error: Unsupported operating system: $param"
                    ;;
                ja_JP)
                    echo "エラー：サポートされていないオペレーティングシステム：$param"
                    ;;
                *)
                    echo "Error: Unsupported operating system: $param"
                    ;;
            esac
            ;;
        "error_unsupported_arch")
            case "$system_language" in
                zh_CN)
                    echo "错误：不支持的架构：$param"
                    ;;
                en_US)
                    echo "Error: Unsupported architecture: $param"
                    ;;
                ja_JP)
                    echo "エラー：サポートされていないアーキテクチャ：$param"
                    ;;
                *)
                    echo "Error: Unsupported architecture: $param"
                    ;;
            esac
            ;;
        "downloading")
            case "$system_language" in
                zh_CN)
                    echo "正在下载 $param"
                    ;;
                en_US)
                    echo "Downloading $param"
                    ;;
                ja_JP)
                    echo "$param をダウンロードしています"
                    ;;
                *)
                    echo "Downloading $param"
                    ;;
            esac
            ;;
        "error_download_failed")
            case "$system_language" in
                zh_CN)
                    echo "错误：下载失败。"
                    ;;
                en_US)
                    echo "Error: Download failed."
                    ;;
                ja_JP)
                    echo "エラー：ダウンロードに失敗しました。"
                    ;;
                *)
                    echo "Error: Download failed."
                    ;;
            esac
            ;;
        "cleaning")
            case "$system_language" in
                zh_CN)
                    echo "清理 $param"
                    ;;
                en_US)
                    echo "Cleaning up $param"
                    ;;
                ja_JP)
                    echo "$param をクリーンアップしています"
                    ;;
                *)
                    echo "Cleaning up $param"
                    ;;
            esac
            ;;
        "moving")
            case "$system_language" in
                zh_CN)
                    echo "添加执行权限并移动到 /usr/local/bin/"
                    ;;
                en_US)
                    echo "Adding execution permissions and moving to /usr/local/bin/"
                    ;;
                ja_JP)
                    echo "実行権限を追加し、/usr/local/bin/ に移動します"
                    ;;
                *)
                    echo "Adding execution permissions and moving to /usr/local/bin/"
                    ;;
            esac
            ;;
        "installation_complete")
            case "$system_language" in
                zh_CN)
                    echo -e "\033[34mComigo 安装完毕，可以在漫画目录下执行 'comi' 命令扫描漫画了。\033[0m"
                    ;;
                en_US)
                    echo -e "\033[34mComigo is installed. You can now run the 'comi' command in the comics directory to scan for comics.\033[0m"
                    ;;
                ja_JP)
                    echo -e "\033[34mComigoがインストールされました。コミックディレクトリで 'comi' コマンドを実行してコミックをスキャンできます。\033[0m"
                    ;;
                *)
                    echo -e "\033[34mComigo is installed. You can now run the 'comi' command in the comics directory to scan for comics.\033[0m"
                    ;;
            esac
            ;;
        *)
            echo "$key"
            ;;
    esac
}

# 检查依赖
dependencies=("tar")

# 检查是否有 curl 或 wget
if command -v curl &> /dev/null; then
    download_tool="curl"
elif command -v wget &> /dev/null; then
    download_tool="wget"
else
    print_message "error_no_curl_wget"
    exit 1
fi

for cmd in "${dependencies[@]}"; do
    if ! command -v "$cmd" &> /dev/null; then
        print_message "error_cmd_not_found" "$cmd"
        exit 1
    fi
done

# 获取最新版本标签
if [ "$download_tool" = "curl" ]; then
    latest_release=$(curl --silent "https://api.github.com/repos/yumenaka/comigo/releases/latest")
else
    latest_release=$(wget -qO- "https://api.github.com/repos/yumenaka/comigo/releases/latest")
fi

# 使用 sed 提取最新标签(MacOS不支持 grep -P)
latest_tag=$(echo "$latest_release" | sed -n 's/.*"tag_name": *"\(v[^"]*\)".*/\1/p')
if [[ -z "$latest_tag" ]]; then
    print_message "error_cannot_get_latest_tag"
    exit 1
fi

# 提取版本号，支持两个或多个数字段的版本号
if [[ $latest_tag =~ ^v([0-9]+(\.[0-9]+)+) ]]; then
    Version="${BASH_REMATCH[1]}"
else
    print_message "error_cannot_parse_version" "$latest_tag"
    exit 1
fi

# 获取操作系统和架构信息
OS=$(uname -s)
ARCH=$(uname -m)

# 映射操作系统名称
case "$OS" in
    Linux)
        OS_NAME="Linux"
        ;;
    Darwin)
        OS_NAME="MacOS"
        ;;
    *)
        print_message "error_unsupported_os" "$OS"
        exit 1
        ;;
esac

# 映射架构名称
case "$ARCH" in
    x86_64)
        ARCH_NAME="x86_64"
        ;;
    armv7l)
        ARCH_NAME="arm"
        ;;
    arm64|aarch64)
        ARCH_NAME="arm64"
        ;;
    *)
        print_message "error_unsupported_arch" "$ARCH"
        exit 1
        ;;
esac

# 构造下载文件名
file_name="comi_v${Version}_${OS_NAME}_${ARCH_NAME}.tar.gz"

# 下载并解压文件
url="https://github.com/yumenaka/comigo/releases/download/${latest_tag}/${file_name}"
print_message "downloading" "$url"


# =====================
# 下载404错误检测
# =====================
if [ "$download_tool" = "curl" ]; then
    # -f 遇到HTTP非200就返回错误码，不写进文件
    if ! curl -fSL -o "$file_name" "$url"; then
        print_message "error_download_failed"
        exit 1
    fi
else
    # --tries 和 --timeout 可酌情添加
    if ! wget --tries=3 --timeout=55 -O "$file_name" "$url"; then
        print_message "error_download_failed"
        exit 1
    fi
fi


if [ "$download_tool" = "curl" ]; then
    curl -L -o "$file_name" "$url"
else
    wget -O "$file_name" "$url"
fi

if [[ ! -f "$file_name" ]]; then
    print_message "error_download_failed"
    exit 1
fi

# =====================
# 检查是否真的是 gzip 格式 无效时可能是文本（比如404页面：ASCII text, with no line terminators）
# =====================
if ! file "$file_name" | grep -q "gzip compressed data"; then
    print_message "error_file_not_gzip"
    rm -f "$file_name"  # 清理无效文件
    exit 1
fi

tar -xzf "$file_name"

# 清理下载文件
print_message "cleaning" "$file_name"
rm "$file_name"

# 添加执行权限并移动到 /usr/local/bin/
print_message "moving"
chmod +x comi

if [ "$EUID" -ne 0 ]; then
    sudo mv comi /usr/local/bin/
else
    mv comi /usr/local/bin/
fi

# 提示用户
print_message "installation_complete"
