# 查看 golang 支持哪些架构：
# go tool dist list

NAME=comi
FULL_NAME=comigo

OS := $(shell uname)
BINDIR := ./bin
MD5_TEXTFILE := $(BINDIR)/md5Sums.txt

# go: cannot install cross-compiled binaries when GOBIN is set
unexport GOBIN

# 当前 cmd/comi 的 Tailscale(tsnet) 与 modernc sqlite 依赖均可在 CGO_ENABLED=0 下跨平台编译。
# 发布目标统一禁用 CGO，避免为跨平台编译维护额外 Docker 镜像。
LDFLAGS := -s -w -X 'github.com/yumenaka/comigo/config.version=$(VERSION)'
GOBUILD := go build -trimpath -ldflags "$(LDFLAGS)"
GOBUILD_CROSS := CGO_ENABLED=0 $(GOBUILD)
GOBUILD_WINDOWS_GUI := CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS) -H=windowsgui"
GOVERSIONINFO ?= goversioninfo
WASM_DIR := assets/static/wasm
WASM_EXEC_JS := $(shell go env GOROOT 2>/dev/null)/lib/wasm/wasm_exec.js
WASM_EXEC_JS_LEGACY := $(shell go env GOROOT 2>/dev/null)/misc/wasm/wasm_exec.js

ifeq ($(OS), Darwin)
  MD5_UTIL = md5
  SED_INPLACE = sed -i ''
else
  MD5_UTIL = md5sum
  SED_INPLACE = sed -i
endif

COMI_TARGETS := Windows_x86_64 Windows_i386 Windows_arm64 Linux_x86_64 Linux_i386 Linux_armv7 Linux_arm64 MacOS_x86_64 MacOS_arm64
WINDOWS_FULL_TARGETS := Windows_x86_64_full Windows_i386_full Windows_arm64_full

all: dmg $(WINDOWS_FULL_TARGETS) compileAll deb-all md5SumThemAll
compileAll: $(COMI_TARGETS)
windows-full: $(WINDOWS_FULL_TARGETS)
android: Linux_arm_android Linux_arm64-android

gomobile:
	gomobile bind -target=android -o comigo.aar -androidapi 26

.PHONY: build-wasm
build-wasm:
	@echo "==> Building reader archive wasm"
	@mkdir -p $(WASM_DIR)
	CGO_ENABLED=0 GOOS=js GOARCH=wasm go build -trimpath -ldflags "-s -w" -o $(WASM_DIR)/archive.wasm ./cmd/archivewasm
	@if [ -f "$(WASM_EXEC_JS)" ]; then \
		cp "$(WASM_EXEC_JS)" "$(WASM_DIR)/wasm_exec.js"; \
	elif [ -f "$(WASM_EXEC_JS_LEGACY)" ]; then \
		cp "$(WASM_EXEC_JS_LEGACY)" "$(WASM_DIR)/wasm_exec.js"; \
	else \
		echo "wasm_exec.js not found in Go toolchain" >&2; \
		exit 1; \
	fi

md5SumThemAll:
	rm -f $(MD5_TEXTFILE)
	find $(BINDIR) -type f -name "$(NAME)_*" -exec $(MD5_UTIL) {} >> $(MD5_TEXTFILE) \;
	find $(BINDIR) -type f -name "$(FULL_NAME)_*" -exec $(MD5_UTIL) {} >> $(MD5_TEXTFILE) \;
	if [ -f "$(BINDIR)/Comigo_$(VERSION).dmg" ]; then \
		echo "==> 计算 Comigo_$(VERSION).dmg 的 md5"; \
		$(MD5_UTIL) "$(BINDIR)/Comigo_$(VERSION).dmg" >> $(MD5_TEXTFILE); \
	fi
	$(SED_INPLACE) 's|./bin/||g' $(MD5_TEXTFILE)
	cat $(MD5_TEXTFILE)

define build_comi_tar
$1: build-wasm
	@mkdir -p $(BINDIR)/$(NAME)_$(VERSION)_$1
	GOOS=$2 GOARCH=$3 $(if $4,GOARM=$4 )$(GOBUILD_CROSS) -o $(BINDIR)/$(NAME)_$(VERSION)_$1/$(NAME) ./cmd/comi
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$1 -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$1.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$1
endef

define build_windows_cli
$1: build-wasm
	cd cmd/comi && $(GOVERSIONINFO) $3 -icon=../../icon.ico -manifest=goversioninfo.exe.manifest versioninfo.json
	@mkdir -p $(BINDIR)/$(NAME)_$(VERSION)_$1
	GOOS=windows GOARCH=$2 $(GOBUILD_CROSS) -o $(BINDIR)/$(NAME)_$(VERSION)_$1/$(NAME).exe ./cmd/comi
	zip -m -r -j -9 $(BINDIR)/$(NAME)_$(VERSION)_$1.zip $(BINDIR)/$(NAME)_$(VERSION)_$1
	rmdir $(BINDIR)/$(NAME)_$(VERSION)_$1
	rm -f cmd/comi/resource.syso
endef

define build_windows_full
$1: build-wasm
	cd cmd/comigo && $(GOVERSIONINFO) $3 -icon=../../icon.ico -manifest=goversioninfo.exe.manifest versioninfo.json
	@mkdir -p $(BINDIR)/$(FULL_NAME)_$(VERSION)_$1
	GOOS=windows GOARCH=$2 $(GOBUILD_WINDOWS_GUI) -o $(BINDIR)/$(FULL_NAME)_$(VERSION)_$1/$(FULL_NAME).exe ./cmd/comigo
	zip -m -r -j -9 $(BINDIR)/$(FULL_NAME)_$(VERSION)_$1.zip $(BINDIR)/$(FULL_NAME)_$(VERSION)_$1
	rmdir $(BINDIR)/$(FULL_NAME)_$(VERSION)_$1
	rm -f cmd/comigo/resource.syso
endef

$(eval $(call build_windows_cli,Windows_x86_64,amd64,-64))
$(eval $(call build_windows_cli,Windows_i386,386,))
$(eval $(call build_windows_cli,Windows_arm64,arm64,-arm -64))
$(eval $(call build_windows_full,Windows_x86_64_full,amd64,-64))
$(eval $(call build_windows_full,Windows_i386_full,386,))
$(eval $(call build_windows_full,Windows_arm64_full,arm64,-arm -64))

$(eval $(call build_comi_tar,Linux_armv6,linux,arm,6))
$(eval $(call build_comi_tar,Linux_armv7,linux,arm,7))
$(eval $(call build_comi_tar,Linux_arm64,linux,arm64,))
$(eval $(call build_comi_tar,Linux_x86_64,linux,amd64,))
$(eval $(call build_comi_tar,Linux_i386,linux,386,))
$(eval $(call build_comi_tar,MacOS_x86_64,darwin,amd64,))
$(eval $(call build_comi_tar,MacOS_arm64,darwin,arm64,))

#Android，32位arm，Termux
Linux_arm_android:
	GOARCH=arm GOOS=android $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) cmd/comi/main.go
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@ -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Android，64位arm，Termux
Linux_arm64-android:
	GOARCH=arm64 GOOS=android $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) cmd/comi/main.go
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@ -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

## ============================================================================
## Debian Package (.deb) Build
## ============================================================================
## Usage:
##   make deb-amd64          - Build amd64 .deb package
##   make deb-arm64          - Build arm64 .deb package
##   make deb-all            - Build all .deb packages
## ============================================================================

.PHONY: deb-amd64 deb-arm64 deb-all deb-clean

DEB_NAME := comigo
DEB_MAINTAINER := Comigo Project <www.bailin.tv@gmail.com>
DEB_DESCRIPTION := Comic/Book reader server supporting ZIP, RAR, CBZ, EPUB, PDF formats

# Build amd64 .deb package
deb-amd64: build-wasm
	@echo "==> Building amd64 .deb package..."
	$(eval DEB_ARCH := amd64)
	$(eval DEB_DIR := $(BINDIR)/$(NAME)_$(VERSION)_$(DEB_ARCH))
	@mkdir -p $(DEB_DIR)/DEBIAN
	@mkdir -p $(DEB_DIR)/usr/bin
	@mkdir -p $(DEB_DIR)/lib/systemd/system
	@# Build binary
	GOARCH=amd64 GOOS=linux $(GOBUILD_CROSS) -o $(DEB_DIR)/usr/bin/$(NAME) cmd/comi/main.go
	@# Create control file
	@# Debian 包版本号必须是纯数字格式，去掉 v 前缀
	$(eval DEB_VERSION := $(patsubst v%,%,$(VERSION)))
	@echo "Package: $(DEB_NAME)" > $(DEB_DIR)/DEBIAN/control
	@echo "Version: $(DEB_VERSION)" >> $(DEB_DIR)/DEBIAN/control
	@echo "Section: web" >> $(DEB_DIR)/DEBIAN/control
	@echo "Priority: optional" >> $(DEB_DIR)/DEBIAN/control
	@echo "Architecture: $(DEB_ARCH)" >> $(DEB_DIR)/DEBIAN/control
	@echo "Maintainer: $(DEB_MAINTAINER)" >> $(DEB_DIR)/DEBIAN/control
	@echo "Description: $(DEB_DESCRIPTION)" >> $(DEB_DIR)/DEBIAN/control
	@echo "Homepage: https://github.com/yumenaka/comigo" >> $(DEB_DIR)/DEBIAN/control
	@echo "Depends: libc6" >> $(DEB_DIR)/DEBIAN/control
	@# Copy systemd service file
	@cp sample/systemd/comigo.service $(DEB_DIR)/lib/systemd/system/
	@# Build .deb package with root ownership
	dpkg-deb --root-owner-group --build $(DEB_DIR)
	@rm -rf $(DEB_DIR)
	@echo "==> Created $(BINDIR)/$(NAME)_$(VERSION)_$(DEB_ARCH).deb"

# Build arm64 .deb package
deb-arm64: build-wasm
	@echo "==> Building arm64 .deb package..."
	$(eval DEB_ARCH := arm64)
	$(eval DEB_DIR := $(BINDIR)/$(NAME)_$(VERSION)_$(DEB_ARCH))
	@mkdir -p $(DEB_DIR)/DEBIAN
	@mkdir -p $(DEB_DIR)/usr/bin
	@mkdir -p $(DEB_DIR)/lib/systemd/system
	@# Build binary
	GOARCH=arm64 GOOS=linux $(GOBUILD_CROSS) -o $(DEB_DIR)/usr/bin/$(NAME) cmd/comi/main.go
	@# Create control file
	@# Debian 包版本号必须是纯数字格式，去掉 v 前缀
	$(eval DEB_VERSION := $(patsubst v%,%,$(VERSION)))
	@echo "Package: $(DEB_NAME)" > $(DEB_DIR)/DEBIAN/control
	@echo "Version: $(DEB_VERSION)" >> $(DEB_DIR)/DEBIAN/control
	@echo "Section: web" >> $(DEB_DIR)/DEBIAN/control
	@echo "Priority: optional" >> $(DEB_DIR)/DEBIAN/control
	@echo "Architecture: $(DEB_ARCH)" >> $(DEB_DIR)/DEBIAN/control
	@echo "Maintainer: $(DEB_MAINTAINER)" >> $(DEB_DIR)/DEBIAN/control
	@echo "Description: $(DEB_DESCRIPTION)" >> $(DEB_DIR)/DEBIAN/control
	@echo "Homepage: https://github.com/yumenaka/comigo" >> $(DEB_DIR)/DEBIAN/control
	@echo "Depends: libc6" >> $(DEB_DIR)/DEBIAN/control
	@# Copy systemd service file
	@cp sample/systemd/comigo.service $(DEB_DIR)/lib/systemd/system/
	@# Build .deb package with root ownership
	dpkg-deb --root-owner-group --build $(DEB_DIR)
	@rm -rf $(DEB_DIR)
	@echo "==> Created $(BINDIR)/$(NAME)_$(VERSION)_$(DEB_ARCH).deb"

# Build all .deb packages
deb-all: deb-amd64 deb-arm64
	@echo "==> All .deb packages built successfully"

# Clean .deb build artifacts
deb-clean:
	@rm -f $(BINDIR)/*.deb
	@echo "==> Cleaned .deb packages"
