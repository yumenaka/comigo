#!/bin/bash

# ============================================================================
# Comigo 一键安装脚本
# ============================================================================
# 使用 curl：
#   bash <(curl -s https://raw.githubusercontent.com/yumenaka/comigo/master/get_comigo.sh)
# 使用 wget：
#   bash <(wget -qO- https://raw.githubusercontent.com/yumenaka/comigo/master/get_comigo.sh)
# ============================================================================

# 遇到错误时立即退出
set -e

# ============================================================================
# 初始化：获取系统语言
# ============================================================================
system_language=$(locale | grep -E '^LANG=' | cut -d= -f2 | cut -d. -f1)

# ============================================================================
# 国际化消息输出函数
# ============================================================================
# 根据系统语言输出多语言消息（支持中文、英文、日文）
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
        "error_file_not_gzip")
            case "$system_language" in
                zh_CN)
                    echo "错误：下载的文件不是有效的 gzip 压缩文件。"
                    ;;
                en_US)
                    echo "Error: Downloaded file is not a valid gzip archive."
                    ;;
                ja_JP)
                    echo "エラー：ダウンロードしたファイルは有効なgzip圧縮ファイルではありません。"
                    ;;
                *)
                    echo "Error: Downloaded file is not a valid gzip archive."
                    ;;
            esac
            ;;
        "error_file_not_found")
            case "$system_language" in
                zh_CN)
                    echo "错误：未找到文件 $param"
                    ;;
                en_US)
                    echo "Error: File not found: $param"
                    ;;
                ja_JP)
                    echo "エラー：ファイルが見つかりません：$param"
                    ;;
                *)
                    echo "Error: File not found: $param"
                    ;;
            esac
            ;;
        "error_cannot_execute")
            case "$system_language" in
                zh_CN)
                    echo "错误：文件 $param 无法执行。"
                    ;;
                en_US)
                    echo "Error: Cannot execute file: $param"
                    ;;
                ja_JP)
                    echo "エラー：ファイルを実行できません：$param"
                    ;;
                *)
                    echo "Error: Cannot execute file: $param"
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
            # param 是安装目录路径
            local install_dir="$param"
            # 将绝对路径转换为 ~ 格式（如果可能）
            local path_display="$install_dir"
            if [[ "$install_dir" == "$HOME"* ]]; then
                path_display="~${install_dir#$HOME}"
            fi
            case "$system_language" in
                zh_CN)
                    echo "添加执行权限并移动到 $path_display"
                    ;;
                en_US)
                    echo "Adding execution permissions and moving to $path_display"
                    ;;
                ja_JP)
                    echo "実行権限を追加し、$path_display に移動します"
                    ;;
                *)
                    echo "Adding execution permissions and moving to $path_display"
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
        "system_info")
            case "$system_language" in
                zh_CN)
                    echo "检测到系统：$param"
                    ;;
                en_US)
                    echo "Detected system: $param"
                    ;;
                ja_JP)
                    echo "検出されたシステム：$param"
                    ;;
                *)
                    echo "Detected system: $param"
                    ;;
            esac
            ;;
        "installing_version")
            case "$system_language" in
                zh_CN)
                    echo "准备安装版本：$param"
                    ;;
                en_US)
                    echo "Preparing to install version: $param"
                    ;;
                ja_JP)
                    echo "インストールするバージョン：$param"
                    ;;
                *)
                    echo "Preparing to install version: $param"
                    ;;
            esac
            ;;
        "extracting")
            case "$system_language" in
                zh_CN)
                    echo "正在解压 $param"
                    ;;
                en_US)
                    echo "Extracting $param"
                    ;;
                ja_JP)
                    echo "$param を展開しています"
                    ;;
                *)
                    echo "Extracting $param"
                    ;;
            esac
            ;;
        "verifying")
            case "$system_language" in
                zh_CN)
                    echo "正在验证安装..."
                    ;;
                en_US)
                    echo "Verifying installation..."
                    ;;
                ja_JP)
                    echo "インストールを確認しています..."
                    ;;
                *)
                    echo "Verifying installation..."
                    ;;
            esac
            ;;
        "selecting_install_path")
            case "$system_language" in
                zh_CN)
                    echo "选择安装路径：$param"
                    ;;
                en_US)
                    echo "Selected install path: $param"
                    ;;
                ja_JP)
                    echo "インストールパスを選択：$param"
                    ;;
                *)
                    echo "Selected install path: $param"
                    ;;
            esac
            ;;
        "path_not_in_path")
            case "$system_language" in
                zh_CN)
                    echo "警告：安装目录 $param 不在 PATH 环境变量中"
                    ;;
                en_US)
                    echo "Warning: Install directory $param is not in PATH"
                    ;;
                ja_JP)
                    echo "警告：インストールディレクトリ $param が PATH 環境変数に含まれていません"
                    ;;
                *)
                    echo "Warning: Install directory $param is not in PATH"
                    ;;
            esac
            ;;
        "path_config_hint")
            # param 格式：install_dir|shell_config
            local install_dir=$(echo "$param" | cut -d'|' -f1)
            local shell_config=$(echo "$param" | cut -d'|' -f2)
            # 将绝对路径转换为 ~ 格式（如果可能）
            local path_display="$install_dir"
            if [[ "$install_dir" == "$HOME"* ]]; then
                path_display="~${install_dir#$HOME}"
            fi
            case "$system_language" in
                zh_CN)
                    echo -e "\033[33m请将以下内容添加到您的 shell 配置文件中（$shell_config）：\033[0m"
                    echo -e "\033[36mexport PATH=\"$path_display:\$PATH\"\033[0m"
                    echo ""
                    echo "或者运行："
                    echo -e "\033[36mecho 'export PATH=\"$path_display:\$PATH\"' >> $shell_config\033[0m"
                    ;;
                en_US)
                    echo -e "\033[33mPlease add the following to your shell config file ($shell_config):\033[0m"
                    echo -e "\033[36mexport PATH=\"$path_display:\$PATH\"\033[0m"
                    echo ""
                    echo "Or run:"
                    echo -e "\033[36mecho 'export PATH=\"$path_display:\$PATH\"' >> $shell_config\033[0m"
                    ;;
                ja_JP)
                    echo -e "\033[33mシェル設定ファイル（$shell_config）に以下を追加してください：\033[0m"
                    echo -e "\033[36mexport PATH=\"$path_display:\$PATH\"\033[0m"
                    echo ""
                    echo "または実行："
                    echo -e "\033[36mecho 'export PATH=\"$path_display:\$PATH\"' >> $shell_config\033[0m"
                    ;;
                *)
                    echo -e "\033[33mPlease add the following to your shell config file ($shell_config):\033[0m"
                    echo -e "\033[36mexport PATH=\"$path_display:\$PATH\"\033[0m"
                    echo ""
                    echo "Or run:"
                    echo -e "\033[36mecho 'export PATH=\"$path_display:\$PATH\"' >> $shell_config\033[0m"
                    ;;
            esac
            ;;
        "file_in_current_dir")
            case "$system_language" in
                zh_CN)
                    echo -e "\033[33m文件已下载到当前目录：\033[0m"
                    echo -e "\033[36m$param\033[0m"
                    echo ""
                    echo "您可以使用以下方式运行："
                    echo -e "\033[36m./comi\033[0m"
                    echo ""
                    echo "或者将其移动到 PATH 中的目录，例如："
                    echo -e "\033[36mmv comi ~/.local/bin/\033[0m"
                    ;;
                en_US)
                    echo -e "\033[33mFile has been downloaded to the current directory:\033[0m"
                    echo -e "\033[36m$param\033[0m"
                    echo ""
                    echo "You can run it using:"
                    echo -e "\033[36m./comi\033[0m"
                    echo ""
                    echo "Or move it to a directory in your PATH, for example:"
                    echo -e "\033[36mmv comi ~/.local/bin/\033[0m"
                    ;;
                ja_JP)
                    echo -e "\033[33mファイルは現在のディレクトリにダウンロードされました：\033[0m"
                    echo -e "\033[36m$param\033[0m"
                    echo ""
                    echo "以下の方法で実行できます："
                    echo -e "\033[36m./comi\033[0m"
                    echo ""
                    echo "または、PATH 内のディレクトリに移動することもできます："
                    echo -e "\033[36mmv comi ~/.local/bin/\033[0m"
                    ;;
                *)
                    echo -e "\033[33mFile has been downloaded to the current directory:\033[0m"
                    echo -e "\033[36m$param\033[0m"
                    echo ""
                    echo "You can run it using:"
                    echo -e "\033[36m./comi\033[0m"
                    echo ""
                    echo "Or move it to a directory in your PATH, for example:"
                    echo -e "\033[36mmv comi ~/.local/bin/\033[0m"
                    ;;
            esac
            ;;
        *)
            echo "$key"
            ;;
    esac
}

# ============================================================================
# 依赖检查
# ============================================================================
dependencies=("tar")

# 检查是否有 curl 或 wget（用于下载文件）
if command -v curl &> /dev/null; then
    download_tool="curl"
elif command -v wget &> /dev/null; then
    download_tool="wget"
else
    print_message "error_no_curl_wget"
    exit 1
fi

# 检查必需的依赖工具
for cmd in "${dependencies[@]}"; do
    if ! command -v "$cmd" &> /dev/null; then
        print_message "error_cmd_not_found" "$cmd"
        exit 1
    fi
done

# ============================================================================
# 版本信息获取
# ============================================================================
# 从 GitHub API 获取最新版本标签
if [ "$download_tool" = "curl" ]; then
    latest_release=$(curl --silent "https://api.github.com/repos/yumenaka/comigo/releases/latest")
else
    latest_release=$(wget -qO- "https://api.github.com/repos/yumenaka/comigo/releases/latest")
fi

# 使用 sed 提取最新标签（MacOS 不支持 grep -P）
latest_tag=$(echo "$latest_release" | sed -n 's/.*"tag_name": *"\(v[^"]*\)".*/\1/p')
if [[ -z "$latest_tag" ]]; then
    print_message "error_cannot_get_latest_tag"
    exit 1
fi

# 从标签中提取版本号（支持两个或多个数字段的版本号，如 v1.2.3）
if [[ $latest_tag =~ ^v([0-9]+(\.[0-9]+)+) ]]; then
    Version="${BASH_REMATCH[1]}"
else
    print_message "error_cannot_parse_version" "$latest_tag"
    exit 1
fi

# ============================================================================
# 系统信息检测
# ============================================================================
# 获取操作系统和架构信息
OS=$(uname -s)
ARCH=$(uname -m)

# 显示检测到的系统信息
print_message "system_info" "${OS} (${ARCH})"

# 显示将要安装的版本
print_message "installing_version" "$latest_tag"

# 将系统标识符映射为发布包使用的名称
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

# ============================================================================
# 文件下载
# ============================================================================
# 构造下载文件名
file_name="comi_v${Version}_${OS_NAME}_${ARCH_NAME}.tar.gz"
url="https://github.com/yumenaka/comigo/releases/download/${latest_tag}/${file_name}"

print_message "downloading" "$file_name"

# 下载文件（包含 404 错误检测）
if [ "$download_tool" = "curl" ]; then
    # -f 遇到HTTP非200就返回错误码，不写进文件
    if ! curl -fSL -o "$file_name" "$url"; then
        print_message "error_download_failed"
        exit 1
    fi
else
    # --tries 和 --timeout 可酌情修改
    if ! wget --tries=3 --timeout=55 -O "$file_name" "$url"; then
        print_message "error_download_failed"
        exit 1
    fi
fi

# 验证下载的文件是否存在
if [[ ! -f "$file_name" ]]; then
    print_message "error_download_failed"
    exit 1
fi

# ============================================================================
# 文件验证
# ============================================================================
# 检查文件是否为有效的 gzip 格式（无效时可能是文本，比如 404 页面）
if ! file "$file_name" | grep -q "gzip compressed data"; then
    print_message "error_file_not_gzip"
    rm -f "$file_name"  # 清理无效文件
    exit 1
fi

# ============================================================================
# 文件解压
# ============================================================================
print_message "extracting" "$file_name"
tar -xzf "$file_name"

# 清理下载的压缩文件
print_message "cleaning" "$file_name"
rm "$file_name"

# ============================================================================
# 路径选择函数
# ============================================================================

# 检测目录是否在 PATH 环境变量中
is_in_path() {
    local dir="$1"
    # 规范化路径（去除尾部斜杠，转换为绝对路径）
    local normalized_dir=$(cd "$dir" 2>/dev/null && pwd || echo "$dir")
    # 检查 PATH 中是否包含该目录
    case ":$PATH:" in
        *":$normalized_dir:"*) return 0 ;;
        *":$normalized_dir/:"*) return 0 ;;
        *) return 1 ;;
    esac
}

# 获取 shell 配置文件路径
get_shell_config_file() {
    local shell_name=$(basename "$SHELL" 2>/dev/null || echo "bash")
    local home_dir="$HOME"
    
    case "$shell_name" in
        bash)
            if [ -f "$home_dir/.bash_profile" ]; then
                echo "$home_dir/.bash_profile"
            elif [ -f "$home_dir/.bashrc" ]; then
                echo "$home_dir/.bashrc"
            else
                echo "$home_dir/.bash_profile"
            fi
            ;;
        zsh)
            if [ -f "$home_dir/.zshrc" ]; then
                echo "$home_dir/.zshrc"
            else
                echo "$home_dir/.zshrc"
            fi
            ;;
        fish)
            if [ -d "$home_dir/.config/fish" ]; then
                echo "$home_dir/.config/fish/config.fish"
            else
                echo "$home_dir/.config/fish/config.fish"
            fi
            ;;
        *)
            # 默认使用 .profile
            echo "$home_dir/.profile"
            ;;
    esac
}

# 选择安装目录
select_install_directory() {
    local install_dir=""
    local need_sudo=false
    
    # 1. 检查环境变量
    if [ -n "$COMIGO_INSTALL_DIR" ]; then
        install_dir="$COMIGO_INSTALL_DIR"
        # 如果环境变量指定的目录需要 root 权限
        if [[ "$install_dir" == /usr/* ]] && [ "$EUID" -ne 0 ]; then
            need_sudo=true
        fi
    # 2. 优先使用 ~/.local/bin（现代标准，无需 root）
    elif [ -d "$HOME/.local/bin" ] || mkdir -p "$HOME/.local/bin" 2>/dev/null; then
        install_dir="$HOME/.local/bin"
    # 3. 备选 ~/bin（仅在 PATH 中时使用）
    elif is_in_path "$HOME/bin"; then
        # ~/bin 在 PATH 中，可以使用（如果不存在则创建）
        if [ -d "$HOME/bin" ] || mkdir -p "$HOME/bin" 2>/dev/null; then
            install_dir="$HOME/bin"
        fi
    fi
    # 4. 如果前面都没有选择，检查 /usr/local/bin
    if [ -z "$install_dir" ]; then
        # 检查 /usr/local/bin 是否在 PATH 中
        if is_in_path "/usr/local/bin"; then
            # 在 PATH 中，可以使用
            install_dir="/usr/local/bin"
            if [ "$EUID" -ne 0 ]; then
                need_sudo=true
            fi
        else
            # 不在 PATH 中，使用当前目录
            install_dir="."
            need_sudo=false
        fi
    fi
    
    # 确保目录存在
    if [ "$need_sudo" = true ]; then
        if [ ! -d "$install_dir" ]; then
            sudo mkdir -p "$install_dir" 2>/dev/null || {
                # 如果创建失败，回退到用户目录
                install_dir="$HOME/.local/bin"
                mkdir -p "$install_dir" 2>/dev/null || {
                    # 如果 ~/.local/bin 也失败，尝试 ~/bin（如果它在 PATH 中）
                    if is_in_path "$HOME/bin"; then
                        install_dir="$HOME/bin"
                        mkdir -p "$install_dir" 2>/dev/null
                    fi
                }
                need_sudo=false
            }
        fi
    else
        mkdir -p "$install_dir" 2>/dev/null || {
            # 如果创建失败，尝试使用备选路径
            if [ "$install_dir" != "$HOME/bin" ]; then
                # 尝试 ~/bin（如果它在 PATH 中）
                if is_in_path "$HOME/bin"; then
                    install_dir="$HOME/bin"
                    mkdir -p "$install_dir" 2>/dev/null || {
                        # 最后尝试系统目录（如果它在 PATH 中）
                        if is_in_path "/usr/local/bin"; then
                            install_dir="/usr/local/bin"
                            need_sudo=true
                            if [ "$EUID" -ne 0 ]; then
                                sudo mkdir -p "$install_dir" 2>/dev/null
                            else
                                mkdir -p "$install_dir" 2>/dev/null
                            fi
                        else
                            # /usr/local/bin 不在 PATH 中，使用当前目录
                            install_dir="."
                            need_sudo=false
                        fi
                    }
                else
                    # ~/bin 不在 PATH 中，检查 /usr/local/bin
                    if is_in_path "/usr/local/bin"; then
                        # /usr/local/bin 在 PATH 中，可以使用
                        install_dir="/usr/local/bin"
                        need_sudo=true
                        if [ "$EUID" -ne 0 ]; then
                            sudo mkdir -p "$install_dir" 2>/dev/null
                        else
                            mkdir -p "$install_dir" 2>/dev/null
                        fi
                    else
                        # /usr/local/bin 也不在 PATH 中，使用当前目录
                        install_dir="."
                        need_sudo=false
                    fi
                fi
            fi
        }
    fi
    
    # 输出选择的目录和是否需要 sudo
    echo "$install_dir|$need_sudo"
}

# ============================================================================
# 安装到系统路径
# ============================================================================
# 验证解压后的文件是否存在
if [[ ! -f "comi" ]]; then
    print_message "error_file_not_found" "comi"
    exit 1
fi

# 添加执行权限
chmod +x comi

# 选择安装目录
install_result=$(select_install_directory)
INSTALL_DIR=$(echo "$install_result" | cut -d'|' -f1)
NEED_SUDO=$(echo "$install_result" | cut -d'|' -f2)

# 显示选择的安装路径
if [ "$INSTALL_DIR" != "." ]; then
    print_message "selecting_install_path" "$INSTALL_DIR"
    
    # 移动到安装目录
    print_message "moving" "$INSTALL_DIR"
    if [ "$NEED_SUDO" = "true" ] && [ "$EUID" -ne 0 ]; then
        sudo mv comi "$INSTALL_DIR/"
    else
        mv comi "$INSTALL_DIR/"
    fi
    
    # ============================================================================
    # PATH 检测和提示
    # ============================================================================
    # 检测安装目录是否在 PATH 中
    if ! is_in_path "$INSTALL_DIR"; then
        print_message "path_not_in_path" "$INSTALL_DIR"
        shell_config=$(get_shell_config_file)
        # 传递安装目录和 shell 配置文件路径（用 | 分隔）
        print_message "path_config_hint" "$INSTALL_DIR|$shell_config"
        echo ""
    fi
else
    # 当前目录，不移动文件，提示用户
    current_dir=$(pwd)
    file_path="$current_dir/comi"
    print_message "file_in_current_dir" "$file_path"
fi

# ============================================================================
# 安装验证
# ============================================================================
if [ "$INSTALL_DIR" != "." ]; then
    print_message "verifying"
    # 检查文件是否存在且可执行
    if [ ! -f "$INSTALL_DIR/comi" ]; then
        print_message "error_file_not_found" "$INSTALL_DIR/comi"
        exit 1
    fi

    if [ ! -x "$INSTALL_DIR/comi" ]; then
        print_message "error_cannot_execute" "$INSTALL_DIR/comi"
        exit 1
    fi

    # 检查命令是否可用（需要重新加载 PATH 或新开终端）
    if command -v comi &> /dev/null; then
        print_message "installation_complete"
    else
        # 文件存在但命令不可用，可能是 PATH 未更新
        print_message "installation_complete"
        if ! is_in_path "$INSTALL_DIR"; then
            shell_config=$(get_shell_config_file)
            echo ""
            case "$system_language" in
                zh_CN)
                    echo -e "\033[33m提示：请重新打开终端或运行以下命令使 PATH 生效：\033[0m"
                    echo -e "\033[36msource $shell_config\033[0m"
                    ;;
                en_US)
                    echo -e "\033[33mNote: Please reopen your terminal or run the following to update PATH:\033[0m"
                    echo -e "\033[36msource $shell_config\033[0m"
                    ;;
                ja_JP)
                    echo -e "\033[33m注意：ターミナルを再起動するか、以下のコマンドを実行して PATH を更新してください：\033[0m"
                    echo -e "\033[36msource $shell_config\033[0m"
                    ;;
                *)
                    echo -e "\033[33mNote: Please reopen your terminal or run the following to update PATH:\033[0m"
                    echo -e "\033[36msource $shell_config\033[0m"
                    ;;
            esac
        fi
    fi
else
    # 当前目录，验证文件存在且可执行
    print_message "verifying"
    if [ ! -f "comi" ]; then
        print_message "error_file_not_found" "comi"
        exit 1
    fi

    if [ ! -x "comi" ]; then
        print_message "error_cannot_execute" "comi"
        exit 1
    fi
    
    print_message "installation_complete"
fi
