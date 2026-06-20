# 查看 golang 支持哪些架构：
# go tool dist list

NAME=comi
TRAY_NAME := comigo-tray
DESKTOP_NAME := comigo-desktop
TRAY_DISPLAY_NAME := Comigo Tray
DESKTOP_DISPLAY_NAME := Comigo Desktop
TRAY_BUNDLE_ID := xyz.comigo.tray
DESKTOP_BUNDLE_ID := xyz.comigo.desktop

OS := $(shell uname -s)
HOST_ARCH := $(shell uname -m)
BINDIR := ./bin
MD5_TEXTFILE := $(BINDIR)/md5Sums.txt

# go: cannot install cross-compiled binaries when GOBIN is set
unexport GOBIN

# 当前 cmd/comi 的 Tailscale(tsnet) 与 modernc sqlite 依赖均可在 CGO_ENABLED=0 下跨平台编译。
# 发布目标统一禁用 CGO，避免为跨平台编译维护额外 Docker 镜像。
LDFLAGS := -s -w -X github.com/yumenaka/comigo/config.version=$(VERSION)
GOBUILD := go build -trimpath -ldflags "$(LDFLAGS)"
GOBUILD_CROSS := CGO_ENABLED=0 $(GOBUILD)
GOBUILD_WINDOWS_GUI := CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS) -H=windowsgui"
GOVERSIONINFO ?= goversioninfo
# Wails Linux 依赖 CGO 与 WebKitGTK，使用 Elastic 预置 crossbuild 镜像补齐编译环境。
WAILS_CLI_VERSION ?= $(shell awk '/github.com\/wailsapp\/wails\/v2/ {print $$2; exit}' go.mod)
WAILS_GO_VERSION ?= $(shell awk '/^go / {print $$2; exit}' go.mod)
WAILS_LINUX_IMAGE_AMD64 ?= docker.elastic.co/beats-dev/golang-crossbuild:$(WAILS_GO_VERSION)-main-debian12
WAILS_LINUX_IMAGE_ARM64 ?= docker.elastic.co/beats-dev/golang-crossbuild:$(WAILS_GO_VERSION)-base-arm-debian12
WAILS_LINUX_RUNTIME_IMAGE_AMD64 ?= comigo-wails-linux:$(WAILS_GO_VERSION)-$(WAILS_CLI_VERSION)-amd64
WAILS_LINUX_RUNTIME_IMAGE_ARM64 ?= comigo-wails-linux:$(WAILS_GO_VERSION)-$(WAILS_CLI_VERSION)-arm64
WAILS_LINUX_DEPS := ca-certificates build-essential pkg-config libgtk-3-dev libwebkit2gtk-4.0-dev
WAILS_BUILD_FLAGS := -m -nosyncgomod -skipembedcreate -trimpath
APP_VERSION := $(patsubst v%,%,$(VERSION))
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
TRAY_TARGETS := tray-Linux_x86_64 tray-Linux_arm64 tray-Windows_x86_64 tray-Windows_arm64 tray-MacOS_universal
DESKTOP_TARGETS := desktop-Linux_x86_64 desktop-Linux_arm64 desktop-Windows_x86_64 desktop-Windows_arm64 desktop-MacOS_universal

ifeq ($(OS), Darwin)
  DESKTOP_CURRENT_TARGET := desktop-MacOS_universal
else ifeq ($(OS), Linux)
  ifeq ($(HOST_ARCH), aarch64)
    DESKTOP_CURRENT_TARGET := desktop-Linux_arm64
  else ifeq ($(HOST_ARCH), arm64)
    DESKTOP_CURRENT_TARGET := desktop-Linux_arm64
  else
    DESKTOP_CURRENT_TARGET := desktop-Linux_x86_64
  endif
else ifneq (,$(findstring MINGW,$(OS)))
  ifeq ($(HOST_ARCH), aarch64)
    DESKTOP_CURRENT_TARGET := desktop-Windows_arm64
  else ifeq ($(HOST_ARCH), arm64)
    DESKTOP_CURRENT_TARGET := desktop-Windows_arm64
  else
    DESKTOP_CURRENT_TARGET := desktop-Windows_x86_64
  endif
endif

.PHONY: all tray-all desktop-all desktop-current compileAll android md5SumThemAll wails-frontend wails-linux-images wails-linux-image-amd64 wails-linux-image-arm64

all: compileAll deb-all tray-all desktop-all md5SumThemAll
tray-all: $(TRAY_TARGETS)
desktop-all: $(DESKTOP_TARGETS)
desktop-current:
	@if [ -z "$(DESKTOP_CURRENT_TARGET)" ]; then \
		echo "Unsupported host for desktop-current: $(OS)/$(HOST_ARCH)" >&2; \
		exit 1; \
	fi
	@$(MAKE) $(DESKTOP_CURRENT_TARGET)
compileAll: $(COMI_TARGETS)
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

wails-frontend: build-wasm
	bun run build

wails-linux-images: wails-linux-image-amd64 wails-linux-image-arm64

wails-linux-image-amd64:
	docker build --platform linux/amd64 \
		--build-arg BASE_IMAGE=$(WAILS_LINUX_IMAGE_AMD64) \
		--build-arg WAILS_CLI_VERSION=$(WAILS_CLI_VERSION) \
		--build-arg WAILS_LINUX_DEPS="$(WAILS_LINUX_DEPS)" \
		-t $(WAILS_LINUX_RUNTIME_IMAGE_AMD64) \
		-f sample/docker/Dockerfile.wails-linux .

wails-linux-image-arm64:
	docker build --platform linux/arm64 \
		--build-arg BASE_IMAGE=$(WAILS_LINUX_IMAGE_ARM64) \
		--build-arg WAILS_CLI_VERSION=$(WAILS_CLI_VERSION) \
		--build-arg WAILS_LINUX_DEPS="$(WAILS_LINUX_DEPS)" \
		-t $(WAILS_LINUX_RUNTIME_IMAGE_ARM64) \
		-f sample/docker/Dockerfile.wails-linux .

md5SumThemAll:
	@mkdir -p $(BINDIR)
	rm -f $(MD5_TEXTFILE)
	touch $(MD5_TEXTFILE)
	find $(BINDIR) -maxdepth 1 -type f \( -name "$(NAME)_$(VERSION)_*" -o -name "$(TRAY_NAME)_$(VERSION)_*" -o -name "$(DESKTOP_NAME)_$(VERSION)_*" \) -exec $(MD5_UTIL) {} >> $(MD5_TEXTFILE) \;
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

define build_tray_tar
$1: build-wasm
	@rm -rf $(BINDIR)/$(TRAY_NAME)_$(VERSION)_$2 $(BINDIR)/$(TRAY_NAME)_$(VERSION)_$2.tar.gz
	@mkdir -p $(BINDIR)/$(TRAY_NAME)_$(VERSION)_$2
	GOOS=$3 GOARCH=$4 $(GOBUILD_CROSS) -o $(BINDIR)/$(TRAY_NAME)_$(VERSION)_$2/$(TRAY_NAME) ./cmd/comigo
	tar --directory=$(BINDIR)/$(TRAY_NAME)_$(VERSION)_$2 -zcvf $(BINDIR)/$(TRAY_NAME)_$(VERSION)_$2.tar.gz $(TRAY_NAME)
	rm -rf $(BINDIR)/$(TRAY_NAME)_$(VERSION)_$2
endef

define build_tray_windows
$1: build-wasm
	cd cmd/comigo && $(GOVERSIONINFO) $4 -icon=../../icon.ico -manifest=goversioninfo.exe.manifest versioninfo.json
	@rm -rf $(BINDIR)/$(TRAY_NAME)_$(VERSION)_$2 $(BINDIR)/$(TRAY_NAME)_$(VERSION)_$2.zip
	@mkdir -p $(BINDIR)/$(TRAY_NAME)_$(VERSION)_$2
	GOOS=windows GOARCH=$3 $(GOBUILD_WINDOWS_GUI) -o $(BINDIR)/$(TRAY_NAME)_$(VERSION)_$2/$(TRAY_NAME).exe ./cmd/comigo
	zip -m -r -j -9 $(BINDIR)/$(TRAY_NAME)_$(VERSION)_$2.zip $(BINDIR)/$(TRAY_NAME)_$(VERSION)_$2
	rmdir $(BINDIR)/$(TRAY_NAME)_$(VERSION)_$2
	rm -f cmd/comigo/resource.syso
endef

define build_desktop_linux
$1: $5 wails-prepare wails-frontend
	@rm -rf build/bin $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_$2 $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_$2.tar.gz
	docker run --rm --platform linux/$3 \
		-v "$(CURDIR)":/workspace \
		-w /workspace \
		-e GOCACHE=/tmp/go-cache \
		-e GOMODCACHE=/go/pkg/mod \
		$4 \
		-p linux/$3 \
		-c "wails build $(WAILS_BUILD_FLAGS) -platform linux/$3 -s -nopackage -o $(DESKTOP_NAME) -ldflags \"$(LDFLAGS)\""
	@mkdir -p $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_$2
	cp build/bin/$(DESKTOP_NAME) $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_$2/$(DESKTOP_NAME)
	tar --directory=$(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_$2 -zcvf $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_$2.tar.gz $(DESKTOP_NAME)
	rm -rf $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_$2 build/bin
endef

define build_desktop_windows
$1: wails-prepare wails-frontend
	@rm -rf build/bin $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_$2 $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_$2.zip
	wails build $(WAILS_BUILD_FLAGS) -platform windows/$3 -s -o $(DESKTOP_NAME).exe -ldflags "$(LDFLAGS)"
	@mkdir -p $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_$2
	cp build/bin/$(DESKTOP_NAME).exe $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_$2/$(DESKTOP_NAME).exe
	zip -m -r -j -9 $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_$2.zip $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_$2
	rmdir $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_$2
	rm -rf build/bin
endef

tray-MacOS_universal: build-wasm
	@$(MAKE) clean-app MAC_APP_NAME=$(TRAY_NAME) MAC_DISPLAY_NAME="$(TRAY_DISPLAY_NAME)" BUNDLE_ID=$(TRAY_BUNDLE_ID) DMG_FILE=$(BINDIR)/$(TRAY_NAME)_$(VERSION)_MacOS_universal.dmg
	@$(MAKE) dmg MAC_APP_NAME=$(TRAY_NAME) MAC_DISPLAY_NAME="$(TRAY_DISPLAY_NAME)" BUNDLE_ID=$(TRAY_BUNDLE_ID) DMG_FILE=$(BINDIR)/$(TRAY_NAME)_$(VERSION)_MacOS_universal.dmg

desktop-MacOS_universal: wails-prepare wails-frontend
	@rm -rf build/bin $(BINDIR)/$(DESKTOP_NAME).app $(BINDIR)/dmg-$(DESKTOP_NAME) $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_MacOS_universal.dmg
	wails build $(WAILS_BUILD_FLAGS) -platform darwin/universal -s -o $(DESKTOP_NAME) -ldflags "$(LDFLAGS)"
	@cp -R build/bin/Comigo.app $(BINDIR)/$(DESKTOP_NAME).app
	@plutil -replace CFBundleName -string "$(DESKTOP_DISPLAY_NAME)" $(BINDIR)/$(DESKTOP_NAME).app/Contents/Info.plist
	@plutil -replace CFBundleDisplayName -string "$(DESKTOP_DISPLAY_NAME)" $(BINDIR)/$(DESKTOP_NAME).app/Contents/Info.plist
	@plutil -replace CFBundleIdentifier -string "$(DESKTOP_BUNDLE_ID)" $(BINDIR)/$(DESKTOP_NAME).app/Contents/Info.plist
	@plutil -replace CFBundleVersion -string "$(APP_VERSION)" $(BINDIR)/$(DESKTOP_NAME).app/Contents/Info.plist
	@plutil -replace CFBundleShortVersionString -string "$(APP_VERSION)" $(BINDIR)/$(DESKTOP_NAME).app/Contents/Info.plist
	@mkdir -p $(BINDIR)/dmg-$(DESKTOP_NAME)
	@cp -R $(BINDIR)/$(DESKTOP_NAME).app $(BINDIR)/dmg-$(DESKTOP_NAME)/$(DESKTOP_NAME).app
	@ln -s /Applications $(BINDIR)/dmg-$(DESKTOP_NAME)/Applications
	@# hdiutil 偶尔误报空间不足；删掉半成品后重试一次即可。
	hdiutil create -volname "$(DESKTOP_DISPLAY_NAME)" -srcfolder $(BINDIR)/dmg-$(DESKTOP_NAME) -ov -format UDZO $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_MacOS_universal.dmg > /dev/null || (rm -f $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_MacOS_universal.dmg && sleep 2 && hdiutil create -volname "$(DESKTOP_DISPLAY_NAME)" -srcfolder $(BINDIR)/dmg-$(DESKTOP_NAME) -ov -format UDZO $(BINDIR)/$(DESKTOP_NAME)_$(VERSION)_MacOS_universal.dmg > /dev/null)
	@rm -rf $(BINDIR)/dmg-$(DESKTOP_NAME) $(BINDIR)/$(DESKTOP_NAME).app build/bin

$(eval $(call build_windows_cli,Windows_x86_64,amd64,-64))
$(eval $(call build_windows_cli,Windows_i386,386,))
$(eval $(call build_windows_cli,Windows_arm64,arm64,-arm -64))
$(eval $(call build_tray_tar,tray-Linux_x86_64,Linux_x86_64,linux,amd64))
$(eval $(call build_tray_tar,tray-Linux_arm64,Linux_arm64,linux,arm64))
$(eval $(call build_tray_windows,tray-Windows_x86_64,Windows_x86_64,amd64,-64))
$(eval $(call build_tray_windows,tray-Windows_arm64,Windows_arm64,arm64,-arm -64))
$(eval $(call build_desktop_linux,desktop-Linux_x86_64,Linux_x86_64,amd64,$(WAILS_LINUX_RUNTIME_IMAGE_AMD64),wails-linux-image-amd64))
$(eval $(call build_desktop_linux,desktop-Linux_arm64,Linux_arm64,arm64,$(WAILS_LINUX_RUNTIME_IMAGE_ARM64),wails-linux-image-arm64))
$(eval $(call build_desktop_windows,desktop-Windows_x86_64,Windows_x86_64,amd64))
$(eval $(call build_desktop_windows,desktop-Windows_arm64,Windows_arm64,arm64))

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
