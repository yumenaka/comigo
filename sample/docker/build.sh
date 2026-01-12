#!/bin/bash

# ============================================================================
# Docker 镜像构建脚本 - Comigo
# ============================================================================
# 功能：
#   - 自动提取版本号
#   - 构建单平台或多平台镜像
#   - 推送到 Docker Hub / ghcr.io
#   - 自动打标签（版本号、latest）
# ============================================================================

set -e  # 遇到错误立即退出

# ============================================================================
# 配置区
# ============================================================================

# 镜像仓库（根据需要修改）
DOCKER_HUB_REPO="yumenaka/comigo"     # Docker Hub 仓库
GHCR_REPO="ghcr.io/yumenaka/comigo"   # GitHub Container Registry

# 支持的平台
PLATFORMS="linux/amd64,linux/arm64,linux/arm/v7"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# ============================================================================
# 辅助函数
# ============================================================================

# 打印信息
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

# 打印警告
warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# 打印错误
error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 显示使用帮助
usage() {
    cat << EOF
使用方法: $0 [选项]

选项:
    -v, --version <version>    指定版本号（例如：v1.2.5）
    -p, --platform <platform>  指定平台（默认：$PLATFORMS）
    -r, --repo <repo>          指定镜像仓库（默认：$DOCKER_HUB_REPO）
    --push                     构建后推送到远程仓库
    --latest                   同时打上 latest 标签
    --ghcr                     推送到 GitHub Container Registry
    --local                    仅构建本地镜像（单平台，不推送）
    --no-cache                 不使用缓存构建
    -h, --help                 显示此帮助信息

示例:
    # 构建本地测试镜像
    $0 --local

    # 构建并推送多平台镜像
    $0 --version v1.2.5 --push --latest

    # 推送到 GitHub Container Registry
    $0 --version v1.2.5 --push --ghcr

    # 构建指定平台
    $0 --platform linux/amd64 --push

EOF
}

# 提取版本号
get_version() {
    # 从 config/version.go 提取版本号
    # 注意：此函数在切换到项目根目录后调用
    if [ -f "config/version.go" ]; then
        VERSION=$(grep -o 'v[0-9]\+\.[0-9]\+\.[0-9]\+' config/version.go | head -1)
    fi
    
    # 如果未找到版本号，使用默认值
    if [ -z "$VERSION" ]; then
        VERSION="v1.0.0"
        warn "未找到版本号，使用默认版本: $VERSION"
    fi
    
    echo "$VERSION"
}

# 检查 Docker Buildx
check_buildx() {
    if ! docker buildx version > /dev/null 2>&1; then
        error "Docker Buildx 未安装或不可用"
        error "请运行: docker buildx install"
        exit 1
    fi
}

# 创建/使用 buildx builder
setup_builder() {
    BUILDER_NAME="comigo-builder"
    
    if ! docker buildx inspect "$BUILDER_NAME" > /dev/null 2>&1; then
        info "创建新的 builder: $BUILDER_NAME"
        docker buildx create --name "$BUILDER_NAME" --use --bootstrap
    else
        info "使用现有的 builder: $BUILDER_NAME"
        docker buildx use "$BUILDER_NAME"
    fi
}

# ============================================================================
# 主函数
# ============================================================================

main() {
    # 默认参数
    VERSION=""
    PLATFORM=""
    REPO=""
    PUSH=false
    TAG_LATEST=false
    USE_GHCR=false
    LOCAL_ONLY=false
    NO_CACHE=""
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--version)
                VERSION="$2"
                shift 2
                ;;
            -p|--platform)
                PLATFORM="$2"
                shift 2
                ;;
            -r|--repo)
                REPO="$2"
                shift 2
                ;;
            --push)
                PUSH=true
                shift
                ;;
            --latest)
                TAG_LATEST=true
                shift
                ;;
            --ghcr)
                USE_GHCR=true
                shift
                ;;
            --local)
                LOCAL_ONLY=true
                shift
                ;;
            --no-cache)
                NO_CACHE="--no-cache"
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                error "未知选项: $1"
                usage
                exit 1
                ;;
        esac
    done
    
    # 切换到项目根目录
    cd "$(dirname "$0")/../.." || exit 1
    info "当前目录: $(pwd)"
    
    # 获取版本号
    if [ -z "$VERSION" ]; then
        VERSION=$(get_version)
    fi
    info "版本号: $VERSION"
    
    # 确定镜像仓库
    if [ -z "$REPO" ]; then
        if [ "$USE_GHCR" = true ]; then
            REPO="$GHCR_REPO"
        else
            REPO="$DOCKER_HUB_REPO"
        fi
    fi
    info "镜像仓库: $REPO"
    
    # 显示标签信息
    if [ "$TAG_LATEST" = true ]; then
        info "将同时打上 latest 标签"
    fi
    
    # 本地构建模式（单平台，不推送）
    if [ "$LOCAL_ONLY" = true ]; then
        info "本地构建模式（仅构建当前平台）"
        docker build \
            --build-arg VERSION="$VERSION" \
            $NO_CACHE \
            -f sample/docker/Dockerfile \
            -t ${REPO}:${VERSION} \
            -t ${REPO}:latest \
            .
        
        info "✅ 本地镜像构建完成"
        info "运行测试: docker run --rm -p 1234:1234 -v \$(pwd)/test:/data ${REPO}:${VERSION}"
        exit 0
    fi
    
    # 多平台构建模式
    check_buildx
    setup_builder
    
    # 确定平台
    if [ -z "$PLATFORM" ]; then
        PLATFORM="$PLATFORMS"
    fi
    info "目标平台: $PLATFORM"
    
    # 构建参数
    BUILD_ARGS="--platform $PLATFORM"
    if [ "$PUSH" = true ]; then
        BUILD_ARGS="$BUILD_ARGS --push"
        info "构建后将推送到远程仓库"
    else
        BUILD_ARGS="$BUILD_ARGS --load"
        warn "未指定 --push，镜像将仅保存到本地"
    fi
    
    # 执行构建
    info "开始构建 Docker 镜像..."
    if [ "$TAG_LATEST" = true ]; then
        docker buildx build \
            $BUILD_ARGS \
            --build-arg VERSION="$VERSION" \
            $NO_CACHE \
            -f sample/docker/Dockerfile \
            -t ${REPO}:${VERSION} \
            -t ${REPO}:latest \
            .
    else
        docker buildx build \
            $BUILD_ARGS \
            --build-arg VERSION="$VERSION" \
            $NO_CACHE \
            -f sample/docker/Dockerfile \
            -t ${REPO}:${VERSION} \
            .
    fi
    
    # 完成
    info "✅ 镜像构建完成"
    
    if [ "$PUSH" = true ]; then
        info "✅ 镜像已推送到: $REPO"
        echo ""
        echo "用户可以使用以下命令拉取镜像："
        echo "  docker pull ${REPO}:${VERSION}"
        if [ "$TAG_LATEST" = true ]; then
            echo "  docker pull ${REPO}:latest"
        fi
    else
        echo ""
        echo "本地测试命令："
        echo "  docker run --rm -p 1234:1234 -v \$(pwd)/test:/data ${REPO}:${VERSION}"
    fi
}

# 执行主函数
main "$@"
