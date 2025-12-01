## macOS App 打包相关变量
## make app
MAC_APP_NAME    := Comigo
BUNDLE_ID   := xyz.comigo.comigo

# 最终 .app 输出目录
BIN_DIR     := bin
# 中间构建目录（放在 bin 下，编译完成后自动清理）
BUILD_DIR   := $(BIN_DIR)/build
APP_DIR     := $(BIN_DIR)/$(MAC_APP_NAME).app

# Info.plist / 图标文件在 assets 目录下
ASSETS_DIR  := assets
INFO_PLIST  := $(ASSETS_DIR)/Info.plist
INFO_PLIST_TMP := $(BUILD_DIR)/Info.plist
ICON_PNG    := $(ASSETS_DIR)/icon.png
ICON_ICNSET := $(BUILD_DIR)/AppIcon.iconset
APP_ICON    := $(ASSETS_DIR)/AppIcon.icns

# 从 config/version.go 提取版本号（不去掉 v 前缀）
# 注意：如果主 Makefile 已经定义了 VERSION，这里不会覆盖
VERSION_GO := config/version.go
ifndef VERSION
  VERSION := $(shell grep -o 'v[0-9]\+\.[0-9]\+\.[0-9]\+' $(VERSION_GO) | head -1)
  ifeq ($(VERSION),)
    VERSION := 1.1.5
  endif
endif

.PHONY: macos-app clean-app app universal icon version

# macOS app 打包的默认目标
macos-app: app

# 从 icon.png 生成 AppIcon.icns
$(APP_ICON): $(ICON_PNG)
	@echo "==> 从 $(ICON_PNG) 生成 $(APP_ICON)"
	@rm -rf $(ICON_ICNSET)
	@mkdir -p $(ICON_ICNSET)
	@echo "==> 生成不同尺寸的图标..."
	@sips -z 16 16     $(ICON_PNG) --out $(ICON_ICNSET)/icon_16x16.png > /dev/null 2>&1
	@sips -z 32 32     $(ICON_PNG) --out $(ICON_ICNSET)/icon_16x16@2x.png > /dev/null 2>&1
	@sips -z 32 32     $(ICON_PNG) --out $(ICON_ICNSET)/icon_32x32.png > /dev/null 2>&1
	@sips -z 64 64     $(ICON_PNG) --out $(ICON_ICNSET)/icon_32x32@2x.png > /dev/null 2>&1
	@sips -z 128 128   $(ICON_PNG) --out $(ICON_ICNSET)/icon_128x128.png > /dev/null 2>&1
	@sips -z 256 256   $(ICON_PNG) --out $(ICON_ICNSET)/icon_128x128@2x.png > /dev/null 2>&1
	@sips -z 256 256   $(ICON_PNG) --out $(ICON_ICNSET)/icon_256x256.png > /dev/null 2>&1
	@sips -z 512 512   $(ICON_PNG) --out $(ICON_ICNSET)/icon_256x256@2x.png > /dev/null 2>&1
	@sips -z 512 512   $(ICON_PNG) --out $(ICON_ICNSET)/icon_512x512.png > /dev/null 2>&1
	@sips -z 1024 1024 $(ICON_PNG) --out $(ICON_ICNSET)/icon_512x512@2x.png > /dev/null 2>&1
	@echo "==> 创建 Contents.json..."
	@echo '{' > $(ICON_ICNSET)/Contents.json
	@echo '  "images" : [' >> $(ICON_ICNSET)/Contents.json
	@echo '    { "size" : "16x16", "idiom" : "mac", "filename" : "icon_16x16.png", "scale" : "1x" },' >> $(ICON_ICNSET)/Contents.json
	@echo '    { "size" : "16x16", "idiom" : "mac", "filename" : "icon_16x16@2x.png", "scale" : "2x" },' >> $(ICON_ICNSET)/Contents.json
	@echo '    { "size" : "32x32", "idiom" : "mac", "filename" : "icon_32x32.png", "scale" : "1x" },' >> $(ICON_ICNSET)/Contents.json
	@echo '    { "size" : "32x32", "idiom" : "mac", "filename" : "icon_32x32@2x.png", "scale" : "2x" },' >> $(ICON_ICNSET)/Contents.json
	@echo '    { "size" : "128x128", "idiom" : "mac", "filename" : "icon_128x128.png", "scale" : "1x" },' >> $(ICON_ICNSET)/Contents.json
	@echo '    { "size" : "128x128", "idiom" : "mac", "filename" : "icon_128x128@2x.png", "scale" : "2x" },' >> $(ICON_ICNSET)/Contents.json
	@echo '    { "size" : "256x256", "idiom" : "mac", "filename" : "icon_256x256.png", "scale" : "1x" },' >> $(ICON_ICNSET)/Contents.json
	@echo '    { "size" : "256x256", "idiom" : "mac", "filename" : "icon_256x256@2x.png", "scale" : "2x" },' >> $(ICON_ICNSET)/Contents.json
	@echo '    { "size" : "512x512", "idiom" : "mac", "filename" : "icon_512x512.png", "scale" : "1x" },' >> $(ICON_ICNSET)/Contents.json
	@echo '    { "size" : "512x512", "idiom" : "mac", "filename" : "icon_512x512@2x.png", "scale" : "2x" }' >> $(ICON_ICNSET)/Contents.json
	@echo '  ],' >> $(ICON_ICNSET)/Contents.json
	@echo '  "info" : { "version" : 1, "author" : "xcode" }' >> $(ICON_ICNSET)/Contents.json
	@echo '}' >> $(ICON_ICNSET)/Contents.json
	@echo "==> 使用 iconutil 打包为 .icns 文件..."
	@iconutil -c icns $(ICON_ICNSET) -o $(APP_ICON)
	@rm -rf $(ICON_ICNSET)
	@echo "==> 已生成 $(APP_ICON)"

icon: $(APP_ICON)

# 从 config/version.go 提取版本号并更新 Info.plist
$(INFO_PLIST_TMP): $(INFO_PLIST) $(VERSION_GO)
	@echo "==> 从 $(VERSION_GO) 提取版本号: $(VERSION)"
	@mkdir -p $(BUILD_DIR)
	@cp $(INFO_PLIST) $(INFO_PLIST_TMP)
	@plutil -replace CFBundleVersion -string "$(VERSION)" $(INFO_PLIST_TMP)
	@plutil -replace CFBundleShortVersionString -string "$(VERSION)" $(INFO_PLIST_TMP)
	@echo "==> 已更新 Info.plist 版本号为 $(VERSION)"

version: $(INFO_PLIST_TMP)

# 构建 amd64 版本（使用 CGO 和版本号）
$(BUILD_DIR)/$(MAC_APP_NAME)_amd64: $(VERSION_GO)
	@echo "==> 构建 amd64 版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -X 'github.com/yumenaka/comigo/config.version=v$(VERSION)'" -o $@ .

# 构建 arm64 版本（使用 CGO 和版本号）
$(BUILD_DIR)/$(MAC_APP_NAME)_arm64: $(VERSION_GO)
	@echo "==> 构建 arm64 版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -trimpath -ldflags "-s -w -X 'github.com/yumenaka/comigo/config.version=v$(VERSION)'" -o $@ .

# 用 lipo 合成 universal binary
$(BUILD_DIR)/$(MAC_APP_NAME): $(BUILD_DIR)/$(MAC_APP_NAME)_amd64 $(BUILD_DIR)/$(MAC_APP_NAME)_arm64
	@echo "==> 使用 lipo 合成 universal binary..."
	@lipo -create -output $@ $^
	@lipo -info $@

universal: $(BUILD_DIR)/$(MAC_APP_NAME)

# 生成 .app 目录结构并塞进去
app: $(BUILD_DIR)/$(MAC_APP_NAME) $(INFO_PLIST_TMP) $(APP_ICON)
	@echo "==> 创建 $(APP_DIR)"
	@mkdir -p "$(APP_DIR)/Contents/MacOS"
	@mkdir -p "$(APP_DIR)/Contents/Resources"

	@echo "==> 拷贝通用二进制"
	@cp "$(BUILD_DIR)/$(MAC_APP_NAME)" "$(APP_DIR)/Contents/MacOS/$(MAC_APP_NAME)"
	@chmod +x "$(APP_DIR)/Contents/MacOS/$(MAC_APP_NAME)"

	@echo "==> 拷贝 Info.plist（版本号: $(VERSION)）"
	@cp "$(INFO_PLIST_TMP)" "$(APP_DIR)/Contents/Info.plist"

	@echo "==> 拷贝 App 图标"
	@cp "$(APP_ICON)" "$(APP_DIR)/Contents/Resources/AppIcon.icns"

	@echo "==> 已生成 $(APP_DIR)"
	@echo "==> 版本号: $(VERSION)"
	@echo "==> 清理中间文件..."
	@rm -rf "$(BUILD_DIR)"
	@echo "==> 构建完成！"

# macOS App 专用的清理目标（避免与跨平台编译的 clean 冲突）
clean-app:
	@rm -rf "$(BUILD_DIR)" "$(APP_DIR)"