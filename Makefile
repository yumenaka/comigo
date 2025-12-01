# Makefile for cross-compilation
# Window icon Need：go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo

## ============================================================================
## 使用说明
## ============================================================================
##
## 本 Makefile 已拆分为模块化结构，主要功能分为两部分：
## 1. macOS App 打包（assets/makefiles/macos-app.mk）
## 2. 跨平台编译（assets/makefiles/cross-compile.mk）
##
## 常用命令：
##   make all              - 编译所有平台（CGO 版本）+App 并生成 MD5 校验
##
## 【跨平台编译】
##   make compileAll_CGO   - 编译所有平台的 CGO 版本
##   make compileAll       - 编译所有平台的非 CGO 版本
##   make Windows_x86_64   - 编译 Windows 64 位版本
##   make Linux_x86_64     - 编译 Linux 64 位版本
##   make MacOS_x86_64     - 编译 macOS Intel 版本
##   make MacOS_arm64      - 编译 macOS Apple Silicon 版本
##
## 【macOS App 打包】
##   make app              - 构建 macOS .app 文件（支持 Intel 和 Apple Silicon）
##   make macos-app        - 同上，别名
##   make icon             - 仅生成 App 图标
##   make clean-app        - 清理 macOS App 构建文件
##
## 【其他】
##   make -n <target>      - 打印编译命令而不实际执行（用于调试）
##   make clean             - 清理所有构建文件
##
## ============================================================================

## Windows Release(Need MSYS2 or mingw32 + find.exe make.exe zip.exe upx.exe):
# mingw32-make all VERSION=v0.9.9

## 仅编译指定架构
# make Linux_x86_64_cgo VERSION=v1.1.5

# 应该下载哪个版本？
#
#| 操作系统    | 设备类型/芯片架构                            | 下载文件              |
#| ----------- | -------------------------------------------- | --------------------- |
#| **MacOS**   | Intel 芯片（2020 年以前的 Mac）              | `MacOS_x86_64.tar.gz` |
#|             | Apple 芯片（M 系列，2020 年以后）            | `MacOS_arm64.tar.gz`  |
#| **Linux**   | ARM 32 位（树莓派 2~4，安装 32 位系统）      | `Linux_armv7.tar.gz`   |
#|             | ARM 64 位（树莓派 4 或 5，安装了 64 位系统） | `Linux_arm64.tar.gz`  |
#| **Windows** | 64 位（大多数 Windows 设备）                 | `Windows_x86_64.zip`  |
#|             | 32 位（较老的 Windows 设备）                 | `Windows_i386.zip`    |
#|             | ARM 架构（如骁龙 Elite 本）                  | `Windows_arm64.zip`   |


## ============================================================================
## 公共变量定义
## ============================================================================

# 从 config/version.go 提取版本号（去掉 v 前缀）
# 如果通过命令行指定了 VERSION，则使用命令行指定的版本
VERSION_GO := config/version.go
ifndef VERSION
  VERSION := $(shell grep -o 'v[0-9]\+\.[0-9]\+\.[0-9]\+' $(VERSION_GO) | head -1 | sed 's/^v//')
  ifeq ($(VERSION),)
    VERSION := 1.0.0
  endif
endif

# 导出 VERSION 变量，供子 Makefile 使用
export VERSION

## ============================================================================
## 引入子 Makefile
## ============================================================================

# 引入 macOS App 打包相关规则
include assets/makefiles/macos-app.mk

# 引入跨平台编译相关规则
include assets/makefiles/cross-compile.mk

## ============================================================================
## 通用清理目标
## ============================================================================

.PHONY: clean

# 清理所有构建文件（包括 macOS App 和跨平台编译的产物）
clean: clean-app
	@echo "==> 清理完成"
