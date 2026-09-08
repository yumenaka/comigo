#!/bin/bash
# Comigo CLI 安装器；兼容 macOS Bash 3.2 和 Linux Bash。
# 先下载脚本并检查下载状态，再执行：bash get.sh --help
# Bash 3.2 的 nounset 不支持空数组展开；可选变量在读取时显式给默认值。
set -eo pipefail

LANGUAGE=${LC_ALL:-${LC_MESSAGES:-${LANG:-en}}}
WORK_DIR=""
STAGE_DIR=""
PRIVILEGE=()

# 独立脚本不依赖应用资源；每条消息保留中英日翻译。
message() {
    local format="$1"
    case "$LANGUAGE" in zh*) format="$2" ;; ja*) format="$3" ;; esac
    shift 3
    printf "$format\n" "$@"
}

die() { message "$@" >&2; exit 1; }

# 暂存文件与目标位于同一文件系统；退出或中断时清理，旧程序始终保留到替换成功。
cleanup() {
    if [[ -n "$STAGE_DIR" ]]; then "${PRIVILEGE[@]}" rm -rf -- "$STAGE_DIR" || :; fi
    if [[ -n "$WORK_DIR" ]]; then rm -rf -- "$WORK_DIR" || :; fi
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

usage() {
    message 'Comigo CLI installer' 'Comigo CLI 安装器' 'Comigo CLI インストーラー'
    message 'Usage: bash get.sh [options]' '用法：bash get.sh [选项]' '使い方：bash get.sh [オプション]'
    message '  --version, -V TAG   Install a release (default: latest stable)' '  --version, -V TAG   指定版本（默认最新稳定版）' '  --version, -V TAG   バージョン指定（既定：最新安定版）'
    message '  --install-dir DIR  Destination (overrides COMIGO_INSTALL_DIR)' '  --install-dir DIR  安装目录（优先于 COMIGO_INSTALL_DIR）' '  --install-dir DIR  インストール先（COMIGO_INSTALL_DIR より優先）'
    message '  --system           Install to /usr/local/bin; sudo only when needed' '  --system           安装到 /usr/local/bin；仅必要时使用 sudo' '  --system           /usr/local/bin にインストール；必要時のみ sudo'
    message '  --force            Overwrite without asking when versions match' '  --force            同版本时直接覆盖，不询问' '  --force            同じバージョンでも確認せず上書き'
    message '  --cn, --proxy      Download through https://comigo.xyz' '  --cn, --proxy      通过 https://comigo.xyz 下载' '  --cn, --proxy      https://comigo.xyz 経由でダウンロード'
    message '  --proxy-base URL   Custom HTTPS proxy base' '  --proxy-base URL   自定义 HTTPS 代理地址' '  --proxy-base URL   HTTPS プロキシを指定'
    message '  --arch ARCH        x86_64, arm64, armv7 or i386' '  --arch ARCH        x86_64、arm64、armv7 或 i386' '  --arch ARCH        x86_64、arm64、armv7 または i386'
    message '  --skip-checksum    Explicitly allow releases without SHA-256 verification' '  --skip-checksum    显式跳过 SHA-256 校验（用于没有清单的旧版本）' '  --skip-checksum    SHA-256 検証を明示的に省略（一覧のない旧リリース用）'
    message '  --help, -h         Show help without installing' '  --help, -h         显示帮助，不安装' '  --help, -h         ヘルプを表示して終了'
    message 'Default: $HOME/.local/bin. Same version: ask on a terminal, otherwise skip.' '默认目录：$HOME/.local/bin。同版本：有终端时询问，否则跳过。' '既定：$HOME/.local/bin。同じバージョン：端末で確認、それ以外はスキップ。'
}

require_value() {
    if [[ -z "${2:-}" || "$2" == -* ]]; then
        die '%s requires a value.' '%s 需要参数值。' '%s には値が必要です。' "$1"
    fi
}

# 只转换实际使用的 API 与资源下载地址，避免维护无关的 raw URL 分支。
release_url() {
    if [[ "$USE_PROXY" == false ]]; then printf '%s\n' "$1"; return; fi
    case "$1" in
        https://api.github.com/*) printf '%s/yumenaka/api.github.com/%s\n' "$PROXY_BASE" "${1#https://api.github.com/}" ;;
        https://github.com/yumenaka/*) printf '%s/yumenaka/%s\n' "$PROXY_BASE" "${1#https://github.com/yumenaka/}" ;;
        *) return 1 ;;
    esac
}

# API 和发布文件共用下载策略，HTTP 错误、重定向、重试和超时行为一致。
download() {
    local url
    url=$(release_url "$1")
    if command -v curl >/dev/null; then
        curl --fail --silent --show-error --location --retry 3 --connect-timeout 15 --max-time 300 \
            --proto '=https' --proto-redir '=https' --output "$2" "$url"
    else
        command -v wget >/dev/null || die 'curl or wget is required.' '需要 curl 或 wget。' 'curl または wget が必要です。'
        wget --quiet --https-only --tries=3 --timeout=30 -O "$2" "$url"
    fi
}

# 只识别 Comigo 的版本行，不能把其他同名命令的任意数字当作版本。
binary_version() {
    local output
    output=$("$1" --version) || return 1
    printf '%s\n' "$output" | sed -n 's/^Comigo \(v[0-9][^[:space:]]*\)[[:space:]]*$/\1/p'
}

# 从终端读取，不消耗 curl | bash 的脚本输入；无终端时安全地跳过同版本。
confirm_reinstall() {
    local answer
    if [[ "$FORCE" == true ]]; then return 0; fi
    if { true </dev/tty; } 2>/dev/null; then
        message 'Same version. [s] Skip (default), [u] overwrite %s:' '已安装同版本。[s] 跳过（默认），[u] 更新并覆盖 %s：' '同じバージョンです。[s] スキップ（既定）、[u] %s を上書き：' "$TARGET" >/dev/tty
        while IFS= read -r answer </dev/tty; do
            case "$answer" in
                u|U|update|更新|覆盖) return 0 ;;
                ''|s|S|skip|跳过) return 1 ;;
                *) message 'Enter s or u:' '请输入 s 或 u：' 's または u を入力してください：' >/dev/tty ;;
            esac
        done
    fi
    return 1
}

# 输出适合当前 shell 的命令；%q 保证空格等字符不会破坏可复制的命令。
path_hint() {
    local shell_name config quoted_dir quoted_config
    shell_name=${SHELL##*/}
    printf -v quoted_dir '%q' "$INSTALL_DIR"
    case "$shell_name" in
        fish)
            message 'Add to PATH with:' '运行以下命令加入 PATH：' 'PATH に追加するには：'
            # Fish 的单引号规则不同于 Bash，使用字面路径并转义反斜杠和单引号。
            local fish_dir=${INSTALL_DIR//\\/\\\\}
            fish_dir=${fish_dir//\'/\\\'}
            printf "fish_add_path '%s'\n" "$fish_dir"
            return ;;
        zsh) config="${ZDOTDIR:-$HOME}/.zshrc" ;;
        bash)
            if [[ "$OS" == Darwin ]]; then config="$HOME/.bash_profile"; else config="$HOME/.bashrc"; fi ;;
        *) config="$HOME/.profile" ;;
    esac
    printf -v quoted_config '%q' "$config"
    message 'Run now, and add this line to %s for future terminals:' '现在运行，并将此行添加到 %s 以供新终端使用：' '今すぐ実行し、新しい端末用に %s にも追加してください：' "$quoted_config"
    printf 'export PATH=%s:"$PATH"\n' "$quoted_dir"
}

main() {
    USE_PROXY=false
    PROXY_BASE=https://comigo.xyz
    VERSION=""
    INSTALL_DIR=${COMIGO_INSTALL_DIR:-}
    SYSTEM=false
    FORCE=false
    SKIP_CHECKSUM=false
    ARCH=""
    SHELL=${SHELL:-/bin/bash}
    while [[ $# -gt 0 ]]; do
        case "$1" in
            -h|--help) usage; return ;;
            --version|-V) require_value "$@"; VERSION="$2"; shift ;;
            --install-dir) require_value "$@"; INSTALL_DIR="$2"; shift ;;
            --system) SYSTEM=true ;;
            --force) FORCE=true ;;
            --cn|--proxy|--use-proxy) USE_PROXY=true ;;
            --proxy-base) require_value "$@"; PROXY_BASE=${2%/}; USE_PROXY=true; shift ;;
            --arch) require_value "$@"; ARCH="$2"; shift ;;
            --skip-checksum) SKIP_CHECKSUM=true ;;
            *) die 'Unknown option: %s' '未知参数：%s' '不明なオプション：%s' "$1" ;;
        esac
        shift
    done
    if [[ "$SYSTEM" == true && -n "$INSTALL_DIR" ]]; then
        die '--system conflicts with --install-dir / COMIGO_INSTALL_DIR.' '--system 不能与 --install-dir / COMIGO_INSTALL_DIR 同用。' '--system と --install-dir / COMIGO_INSTALL_DIR は併用できません。'
    fi
    if [[ "$PROXY_BASE" != https://?* || "$PROXY_BASE" == *[[:space:]\?#]* ]]; then
        die 'Invalid HTTPS proxy base: %s' '无效的 HTTPS 代理地址：%s' '無効な HTTPS プロキシ：%s' "$PROXY_BASE"
    fi
    if [[ -n "$VERSION" && "$VERSION" != v* ]]; then VERSION="v$VERSION"; fi
    if [[ -n "$VERSION" && ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$ ]]; then
        die 'Invalid version: %s' '版本号无效：%s' '無効なバージョン：%s' "$VERSION"
    fi

    OS=$(uname -s)
    if [[ -z "$ARCH" ]]; then
        ARCH=$(uname -m)
        # Rosetta 中优先下载 Apple Silicon 原生程序。
        if [[ "$OS" == Darwin && "$ARCH" == x86_64 ]] && [[ "$(sysctl -in sysctl.proc_translated 2>/dev/null || :)" == 1 ]]; then ARCH=arm64; fi
    fi
    case "$OS" in Linux) OS_NAME=Linux ;; Darwin) OS_NAME=MacOS ;; *) die 'Unsupported OS: %s' '不支持的系统：%s' '未対応の OS：%s' "$OS" ;; esac
    case "$ARCH" in aarch64) ARCH=arm64 ;; armv7l) ARCH=armv7 ;; i[3-6]86) ARCH=i386 ;; esac
    case "$OS/$ARCH" in Linux/x86_64|Linux/arm64|Linux/armv7|Linux/i386|Darwin/x86_64|Darwin/arm64) ;; *) die 'Unsupported platform: %s' '不支持的平台：%s' '未対応のプラットフォーム：%s' "$OS/$ARCH" ;; esac
    local dep
    for dep in sed mktemp; do
        command -v "$dep" >/dev/null || die 'Missing command: %s' '缺少命令：%s' '必要なコマンドがありません：%s' "$dep"
    done

    # 固定默认目录；显式指定目录失败时不回退，也不自动修改 shell 配置。
    if [[ "$SYSTEM" == true ]]; then INSTALL_DIR=/usr/local/bin
    elif [[ -z "$INSTALL_DIR" ]]; then
        [[ -n "${HOME:-}" && "$HOME" == /* ]] || die 'HOME must be an absolute path.' 'HOME 必须是绝对路径。' 'HOME は絶対パスである必要があります。'
        INSTALL_DIR="$HOME/.local/bin"
    fi
    [[ "$INSTALL_DIR" == /* ]] || INSTALL_DIR="$PWD/$INSTALL_DIR"
    local ancestor="$INSTALL_DIR"
    while [[ ! -e "$ancestor" && ! -L "$ancestor" ]]; do ancestor=${ancestor%/*}; [[ -n "$ancestor" ]] || ancestor=/; done
    [[ -d "$ancestor" ]] || die 'Not a directory: %s' '不是目录：%s' 'ディレクトリではありません：%s' "$ancestor"
    TARGET="${INSTALL_DIR%/}/comi"
    if [[ -L "$TARGET" || ( -e "$TARGET" && ! -f "$TARGET" ) ]]; then
        die 'Refusing to replace a symlink or non-file: %s' '不能覆盖符号链接或非普通文件：%s' 'シンボリックリンクまたは通常ファイル以外は上書きできません：%s' "$TARGET"
    fi
    message 'Platform: %s; destination: %s' '平台：%s；安装位置：%s' 'プラットフォーム：%s；インストール先：%s' "$OS_NAME/$ARCH" "$TARGET"

    WORK_DIR=$(mktemp -d "${TMPDIR:-/tmp}/comigo-install.XXXXXX")
    if [[ -z "$VERSION" ]]; then
        download https://api.github.com/repos/yumenaka/comigo/releases/latest "$WORK_DIR/release.json" || die 'Cannot fetch latest release.' '获取最新版本失败。' '最新リリースを取得できません。'
        VERSION=$(sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$WORK_DIR/release.json")
    fi
    [[ "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$ ]] || die 'Invalid release tag: %s' '发布标签无效：%s' '無効なリリースタグ：%s' "$VERSION"

    # 同时检查 PATH 中的可用命令和目标文件，避免目标尚未加入 PATH 时重复安装。
    local existing existing_version candidate same=false
    existing=$(type -P comi || :)
    for candidate in "$existing" "$TARGET"; do
        [[ -n "$candidate" && -f "$candidate" && -x "$candidate" ]] || continue
        existing_version=$(binary_version "$candidate" 2>/dev/null) || existing_version=""
        message 'Existing command: %s (%s)' '已有命令：%s（%s）' '既存のコマンド：%s（%s）' "$candidate" "${existing_version:-unknown}"
        [[ "$existing_version" != "$VERSION" ]] || same=true
        [[ "$candidate" != "$TARGET" ]] || break
        [[ "$existing" != "$TARGET" ]] || break
    done
    if [[ "$same" == true ]] && ! confirm_reinstall; then
        message 'Skipped %s. Use --force to overwrite non-interactively.' '已跳过 %s。非交互覆盖可使用 --force。' '%s をスキップしました。非対話で上書きするには --force を使用してください。' "$VERSION"
        return
    fi

    # 只有确认安装后才检查写权限和解包、校验工具，同版本跳过不需要它们。
    for dep in tar awk install; do
        command -v "$dep" >/dev/null || die 'Missing command: %s' '缺少命令：%s' '必要なコマンドがありません：%s' "$dep"
    done
    # 全局文件必须归 root 所有，即使目录本身可写也需要管理员权限。
    if [[ "$SYSTEM" == true && "$EUID" -ne 0 ]]; then
        command -v sudo >/dev/null || die 'System install requires sudo.' '系统安装需要 sudo。' 'システムへのインストールには sudo が必要です。'
        PRIVILEGE=(sudo)
    elif [[ ! -w "$ancestor" || ! -x "$ancestor" ]]; then
        die 'Directory is not writable: %s' '目录不可写：%s' 'ディレクトリに書き込めません：%s' "$ancestor"
    fi
    if [[ "$SKIP_CHECKSUM" == false ]]; then
        if command -v sha256sum >/dev/null; then HASH_TOOL=(sha256sum)
        elif command -v shasum >/dev/null; then HASH_TOOL=(shasum -a 256)
        else die 'sha256sum or shasum is required.' '需要 sha256sum 或 shasum。' 'sha256sum または shasum が必要です。'; fi
    fi

    local name="comi_${VERSION}_${OS_NAME}_${ARCH}.tar.gz"
    local base="https://github.com/yumenaka/comigo/releases/download/$VERSION"
    message 'Downloading %s' '正在下载 %s' '%s をダウンロード中' "$name"
    download "$base/$name" "$WORK_DIR/archive.tar.gz" || die 'Download failed: %s' '下载失败：%s' 'ダウンロード失敗：%s' "$name"
    if [[ "$SKIP_CHECKSUM" == true ]]; then
        message 'SHA-256 verification explicitly disabled.' '已按参数要求跳过 SHA-256 校验。' '指定により SHA-256 検証を省略します。' >&2
    else
        download "$base/checksums.txt" "$WORK_DIR/checksums" || die 'Cannot fetch checksums.txt. Old releases require explicit --skip-checksum.' '无法获取 checksums.txt。旧版本需显式使用 --skip-checksum。' 'checksums.txt を取得できません。旧リリースには --skip-checksum の明示指定が必要です。'
        local expected actual
        expected=$(awk -v name="$name" '$2 == name {print tolower($1)}' "$WORK_DIR/checksums")
        [[ "$expected" =~ ^[0-9a-fA-F]{64}$ ]] || die 'Missing or invalid SHA-256 for %s' '%s 的 SHA-256 缺失或无效' '%s の SHA-256 がないか無効です' "$name"
        actual=$("${HASH_TOOL[@]}" "$WORK_DIR/archive.tar.gz")
        actual=${actual%% *}
        [[ "$actual" == "$expected" ]] || die 'SHA-256 mismatch: %s' 'SHA-256 不匹配：%s' 'SHA-256 不一致：%s' "$name"
    fi

    # 发布包只允许一个名为 comi 的普通文件，拒绝额外成员和链接。
    local members details
    members=$(tar -tzf "$WORK_DIR/archive.tar.gz") || die 'Invalid archive.' '压缩包无效。' '無効なアーカイブです。'
    details=$(LC_ALL=C tar -tvzf "$WORK_DIR/archive.tar.gz") || return 1
    [[ "$members" == comi && "$details" == -* ]] || die 'Archive must contain only a regular comi file.' '压缩包必须仅包含普通文件 comi。' 'アーカイブには通常ファイル comi のみを含めてください。'
    tar -xzf "$WORK_DIR/archive.tar.gz" -C "$WORK_DIR" comi
    [[ -f "$WORK_DIR/comi" && ! -L "$WORK_DIR/comi" ]] || return 1

    "${PRIVILEGE[@]}" mkdir -p -- "$INSTALL_DIR"
    INSTALL_DIR=$(cd "$INSTALL_DIR" && pwd -P)
    TARGET="$INSTALL_DIR/comi"
    STAGE_DIR=$("${PRIVILEGE[@]}" mktemp -d "$INSTALL_DIR/.comigo-install.XXXXXX")
    # 系统暂存目录允许当前用户执行验证，但只有目录所有者能修改暂存文件。
    "${PRIVILEGE[@]}" chmod 755 "$STAGE_DIR"
    "${PRIVILEGE[@]}" install -m 0755 "$WORK_DIR/comi" "$STAGE_DIR/comi"
    if [[ "$SYSTEM" == true ]]; then "${PRIVILEGE[@]}" chown 0:0 "$STAGE_DIR/comi"; fi
    local installed_version
    installed_version=$(binary_version "$STAGE_DIR/comi") || die 'New binary cannot run; existing installation preserved.' '新程序无法运行，已保留原安装。' '新しいプログラムを実行できません。既存のインストールは保持されます。'
    [[ "$installed_version" == "$VERSION" ]] || die 'Version mismatch: expected %s, got %s; existing installation preserved.' '版本不匹配：期望 %s，实际 %s；已保留原安装。' 'バージョン不一致：期待値 %s、実際 %s。既存のインストールは保持されます。' "$VERSION" "$installed_version"
    "${PRIVILEGE[@]}" mv -f -- "$STAGE_DIR/comi" "$TARGET"
    message 'Installed %s: %s' '已安装 %s：%s' '%s をインストールしました：%s' "$VERSION" "$TARGET"
    printf '  %q --version\n' "$TARGET"
    local resolved
    resolved=$(type -P comi || :)
    if [[ -z "$resolved" || ! "$resolved" -ef "$TARGET" ]]; then
        if [[ -n "$resolved" ]]; then message 'PATH currently selects another command: %s' 'PATH 当前优先使用另一个命令：%s' 'PATH は現在別のコマンドを選択しています：%s' "$resolved"; fi
        path_hint
    fi
}

main "$@"
