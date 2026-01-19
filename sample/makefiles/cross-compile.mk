# 查看 golang 支持那些架构：
# go tool dist list

NAME=comi
FULL_NAME=comigo

OS := $(shell uname)
BINDIR := ./bin
MD5_TEXTFILE := $(BINDIR)/md5Sums.txt
#go: cannot install cross-compiled binaries when GOBIN is set
unexport GOBIN

MAIN_FILE_DIR := ./
# 覆盖版本号，只有 string 可以直接覆盖 包路径必须和 go list 输出的完全一致，需要带 module 路径
# -trimpath 去掉源码绝对路径，比如构建机目录
# -ldflags 指定编译参数。-s 去掉符号信息  -w去掉调试信息 减小二进制体积
GOBUILD=go build -trimpath -ldflags "-s -w -X 'github.com/yumenaka/comigo/config.version=${VERSION}'"

ifeq ($(OS), Darwin)
  MD5_UTIL = md5
else
  MD5_UTIL = md5sum
endif

# 跨平台编译的默认目标
all: app Windows_x86_64_full Windows_i386_full Windows_arm64_full compileAll_CGO deb-all md5SumThemAll

## windows 可能不需要CGO就能支持Tailscale？

# 过去因为sqlite（ent）库的关系，部分架构（Windows_i386）无法正常运行，需要写条件编译代码。但是最近似乎都Pass了，或许可以不再分架构：
# ent库的编译检测状态： https://modern-c.appspot.com/-/builder/?importpath=modernc.org%2Fsqlite
# 为了支持Tailscale，使用docker交叉编译
compileAll: Windows_x86_64 Windows_i386  Windows_arm64 Linux_x86_64 Linux_i386 Linux_armv7 Linux_arm64 MacOS_x86_64 MacOS_arm64
compileAll_CGO: Windows_x86_64 Windows_i386  Windows_arm64 Linux_x86_64_cgo Linux_i386_cgo Linux_armv7_cgo Linux_arm64_cgo MacOS_x86_64_cgo MacOS_arm64_cgo

android: Linux_arm_android Linux_arm64-android

UPX := $(shell command -v upx 2> /dev/null)
DOCKER := $(shell command -v docker 2> /dev/null)

gomobile:
	export ANDROID_NDK_HOME=/Users/bai/Library/Android/sdk/ndk/26.1.10909125
	gomobile bind -target=android -o comigo.aar -androidapi 26

md5SumThemAll:
	rm -f $(MD5_TEXTFILE)
	find $(BINDIR) -type f -name "$(NAME)_*" -exec $(MD5_UTIL) {} >> $(MD5_TEXTFILE) \;
	find $(BINDIR) -type f -name "$(FULL_NAME)_*" -exec $(MD5_UTIL) {} >> $(MD5_TEXTFILE) \;
	# 如果存在 Comigo.app 目录，则先打包为 zip 再计算 md5
	if [ -d "$(BINDIR)/Comigo.app" ]; then \
		echo "==> 打包 Comigo.app 为 Comigo.app.zip 用于计算 md5"; \
		rm -f "$(BINDIR)/Comigo.app.zip"; \
		cd "$(BINDIR)" && zip -r "Comigo.app.zip" "Comigo.app" > /dev/null; \
		cd .. && $(MD5_UTIL) "$(BINDIR)/Comigo.app.zip" >> $(MD5_TEXTFILE); \
	fi
	# 删除 $(MD5_TEXTFILE)里面的 ./bin/ 字符串
	sed -i '' 's|./bin/||g' $(MD5_TEXTFILE)
	cat $(MD5_TEXTFILE)

#docker run：启动一个新的容器执行命令。
#-it：交互模式 + TTY（可看输出、调试）。
#--rm：容器退出后自动删除，不留垃圾。
#-v 参数挂载宿主机目录到容器。左：本地的 Go 项目路径。右边：容器里的路径
# -w /go/src/github.com/user/go-project \
#指定容器的工作目录（go build 或 make 在这里执行）
#-e CGO_ENABLED=1 \
#  设置环境变量：启用 CGO 支持。
#  默认 golang-crossbuild 镜像内置了交叉编译链（gcc、musl、glibc 等），开启后就能编译依赖 C 代码的包。
#--build-cmd "make build" \
#镜像支持的参数：告诉它用什么命令来"构建"你的项目。
#这里执行的是 make build（你项目的 Makefile 中应定义了 build 目标）。

## 查看编译用Docker镜像信息：
# https://github.com/elastic/golang-crossbuild/releases

#  make Linux_armv7_cgo VERSION=v1.0.5
Linux_armv7_cgo:
ifndef DOCKER
	$(error "No docker found! Please install docker to build Linux_armv7_cgo")
endif
ifdef DOCKER
	docker run -it  \
	 -v "$$PWD":/go/src/github.com/user/go-project \
	 -w /go/src/github.com/user/go-project \
	 -e CGO_ENABLED=1 \
	 -e VERSION=$(VERSION) \
	 -e FILE_LABLE="Linux_armv7" \
	 docker.elastic.co/beats-dev/golang-crossbuild:1.25.5-armhf-debian9 \
	 --build-cmd "make linux_armv7_cgo_docker VERSION=$(VERSION)" \
	 -p "linux/armv7"
endif

linux_armv7_cgo_docker:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)/$(NAME) cmd/comi/main.go
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE).tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)

#  make MacOS_x86_64_cgo VERSION=v1.0.5
MacOS_x86_64_cgo:
ifndef DOCKER
	$(error "No docker found! Please install docker to build MacOS_x86_64_cgo")
endif
ifdef DOCKER
	docker run -it  \
	 -v "$$PWD":/go/src/github.com/user/go-project \
	 -w /go/src/github.com/user/go-project \
	 -e CGO_ENABLED=1 \
	 -e VERSION=$(VERSION) \
	 -e FILE_LABLE="MacOS_x86_64" \
	 docker.elastic.co/beats-dev/golang-crossbuild:1.25.5-darwin-debian12 \
	 --build-cmd "make darwin_x86_64_cgo_docker VERSION=$(VERSION)" \
	 -p "darwin/amd64"
endif

darwin_x86_64_cgo_docker:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)/$(NAME) cmd/comi/main.go
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE).tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)

#  make MacOS_arm64_cgo VERSION=v1.0.5
MacOS_arm64_cgo:
ifndef DOCKER
	$(error "No docker found! Please install docker to build MacOS_arm64_cgo")
endif
ifdef DOCKER
	docker run -it \
	 -v "$$PWD":/go/src/github.com/user/go-project \
	 -w /go/src/github.com/user/go-project \
	 -e CGO_ENABLED=1 \
	 -e VERSION=$(VERSION) \
	 -e FILE_LABLE="MacOS_arm64" \
	 docker.elastic.co/beats-dev/golang-crossbuild:1.25.5-darwin-arm64-debian12 \
	 --build-cmd "make darwin_arm64_cgo_docker VERSION=$(VERSION)" \
	 -p "darwin/arm64"
endif

darwin_arm64_cgo_docker:
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)/$(NAME) cmd/comi/main.go
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE).tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)

#  make Linux_arm64_cgo VERSION=v1.0.4
Linux_arm64_cgo:
ifndef DOCKER
	$(error "No docker found! Please install docker to build Linux_arm64_cgo")
endif
ifdef DOCKER
	docker run -it \
	 -v "$$PWD":/go/src/github.com/user/go-project \
	 -w /go/src/github.com/user/go-project \
	 -e CGO_ENABLED=1 \
	 -e VERSION=$(VERSION) \
	 -e FILE_LABLE="Linux_arm64" \
	 docker.elastic.co/beats-dev/golang-crossbuild:1.25.5-base-arm-debian9 \
	 --build-cmd "make linux_arm64_cgo_docker VERSION=$(VERSION)" \
	 -p "linux/arm64"
endif

linux_arm64_cgo_docker:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)/$(NAME) cmd/comi/main.go
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE).tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)

#  make Linux_x86_64_cgo VERSION=v1.0.4
Linux_x86_64_cgo:
ifndef DOCKER
	$(error "No docker found! Please install docker to build Linux_x86_64_cgo")
endif
ifdef DOCKER
	docker run -it  \
	 -v "$$PWD":/go/src/github.com/user/go-project \
	 -w /go/src/github.com/user/go-project \
	 -e CGO_ENABLED=1 \
	 -e VERSION=$(VERSION) \
	 -e FILE_LABLE="Linux_x86_64" \
	 docker.elastic.co/beats-dev/golang-crossbuild:1.25.5-main-debian7 \
	 --build-cmd "make linux_x86_64_cgo_docker VERSION=$(VERSION)" \
	 -p "linux/amd64"
endif

linux_x86_64_cgo_docker:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)/$(NAME) cmd/comi/main.go
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE).tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)

# Add Linux_i386_cgo mirroring the pattern
Linux_i386_cgo:
ifndef DOCKER
	$(error "No docker found! Please install docker to build Linux_i386_cgo")
endif
ifdef DOCKER
	docker run -it  \
	 -v "$$PWD":/go/src/github.com/user/go-project \
	 -w /go/src/github.com/user/go-project \
	 -e CGO_ENABLED=1 \
	 -e VERSION=$(VERSION) \
	 -e FILE_LABLE="Linux_i386" \
	 docker.elastic.co/beats-dev/golang-crossbuild:1.25.5-main-debian7 \
	 --build-cmd "make linux_i386_cgo_docker VERSION=$(VERSION)" \
	 -p "linux/386"
endif

linux_i386_cgo_docker:
	CGO_ENABLED=1 GOOS=linux GOARCH=386 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)/$(NAME) cmd/comi/main.go
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE).tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)

## No CGO build:

# 64位Windows	$(NAME)_$(VERSION)_$@
Windows_x86_64:
	go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo # Window icon Need
	cd cmd/comi && goversioninfo -64 -icon=../../icon.ico -manifest=goversioninfo.exe.manifest versioninfo.json
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME).exe ./cmd/comi 
	zip -m -r -j -9 $(BINDIR)/$(NAME)_$(VERSION)_$@.zip $(BINDIR)/$(NAME)_$(VERSION)_$@
	rmdir $(BINDIR)/$(NAME)_$(VERSION)_$@
	rm  cmd/comi/resource.syso 

# 64位Windows + system_tray	$(NAME)_$(VERSION)_$@
Windows_x86_64_full:
	go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo # Window icon Need
	cd cmd/comigo && goversioninfo -64 -icon=../../icon.ico -manifest=goversioninfo.exe.manifest versioninfo.json
	GOARCH=amd64 GOOS=windows $(GOBUILD) -ldflags -H=windowsgui -o $(BINDIR)/$(FULL_NAME)_$(VERSION)_$@/$(FULL_NAME).exe ./cmd/comigo
	zip -m -r -j -9 $(BINDIR)/$(FULL_NAME)_$(VERSION)_$@.zip $(BINDIR)/$(FULL_NAME)_$(VERSION)_$@
	rmdir $(BINDIR)/$(FULL_NAME)_$(VERSION)_$@
	rm cmd/comigo/resource.syso

# 32位Windows
Windows_i386:
	go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo # Window icon Need
	cd cmd/comi && goversioninfo -icon=../../icon.ico -manifest=goversioninfo.exe.manifest versioninfo.json
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME).exe ./cmd/comi 
	zip -m -r -j -9 $(BINDIR)/$(NAME)_$(VERSION)_$@.zip $(BINDIR)/$(NAME)_$(VERSION)_$@
	rmdir $(BINDIR)/$(NAME)_$(VERSION)_$@
	rm  cmd/comi/resource.syso

# 32位Windows + system_tray
Windows_i386_full:
	go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo # Window icon Need
	cd cmd/comigo && goversioninfo -icon=../../icon.ico -manifest=goversioninfo.exe.manifest versioninfo.json
	GOARCH=386 GOOS=windows $(GOBUILD) -ldflags -H=windowsgui  -o $(BINDIR)/$(FULL_NAME)_$(VERSION)_$@/$(FULL_NAME).exe ./cmd/comigo
	zip -m -r -j -9 $(BINDIR)/$(FULL_NAME)_$(VERSION)_$@.zip $(BINDIR)/$(FULL_NAME)_$(VERSION)_$@
	rmdir $(BINDIR)/$(FULL_NAME)_$(VERSION)_$@
	rm cmd/comigo/resource.syso

#windows arm64
Windows_arm64:
	go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo # Window icon Need
	cd cmd/comi && goversioninfo -arm -64 -icon=../../icon.ico -manifest=goversioninfo.exe.manifest versioninfo.json
	GOARCH=arm64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME).exe ./cmd/comi
	zip -m -r -j -9 $(BINDIR)/$(NAME)_$(VERSION)_$@.zip $(BINDIR)/$(NAME)_$(VERSION)_$@
	rmdir $(BINDIR)/$(NAME)_$(VERSION)_$@
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@
	rm cmd/comi/resource.syso

# windows arm64 + system_tray
Windows_arm64_full:
	go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo # Window icon Need
	cd cmd/comigo && goversioninfo -arm -64 -icon=../../icon.ico -manifest=goversioninfo.exe.manifest versioninfo.json
	GOARCH=arm64 GOOS=windows $(GOBUILD) -ldflags -H=windowsgui -o $(BINDIR)/$(FULL_NAME)_$(VERSION)_$@/$(FULL_NAME).exe ./cmd/comigo
	zip -m -r -j -9 $(BINDIR)/$(FULL_NAME)_$(VERSION)_$@.zip $(BINDIR)/$(FULL_NAME)_$(VERSION)_$@
	rmdir $(BINDIR)/$(FULL_NAME)_$(VERSION)_$@
	rm cmd/comigo/resource.syso

#Linux_armv6 RaspberryPi1,2,zero,GOARM=6：仅使用 VFPv1；交叉编译时默认；通常是 ARM11 或更好的内核（也支持 VFPv2 或更好的内核）
Linux_armv6:
	GOARCH=arm GOOS=linux GOARM=6 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) cmd/comi/main.go
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Linux_armv7，RaspberryPi3 官方32位armv7l系统。GOARM=7：使用 VFPv3；通常是 Cortex-A 内核. 2012年发布的架构。
Linux_armv7:
	GOARCH=arm GOOS=linux GOARM=7 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) cmd/comi/main.go
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#linux，64位arm。2012年发布的架构。
Linux_arm64:
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) cmd/comi/main.go
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Linux，x86_64
Linux_x86_64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) cmd/comi/main.go
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Linux，i386
Linux_i386:
	GOARCH=386 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) cmd/comi/main.go
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#MACOS x86_64
MacOS_x86_64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) cmd/comi/main.go
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@
	
#MACOS arm64 no upx
MacOS_arm64:
	GOARCH=arm64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) cmd/comi/main.go
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@
	
#Android，32位arm，Termux	
Linux_arm_android:
	GOARCH=arm GOOS=android $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) cmd/comi/main.go 
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Android，64位arm，Termux
Linux_arm64-android:
	GOARCH=arm64 GOOS=android $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) cmd/comi/main.go 
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
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
deb-amd64:
	@echo "==> Building amd64 .deb package..."
	$(eval DEB_ARCH := amd64)
	$(eval DEB_DIR := $(BINDIR)/$(NAME)_$(VERSION)_$(DEB_ARCH))
	@mkdir -p $(DEB_DIR)/DEBIAN
	@mkdir -p $(DEB_DIR)/usr/bin
	@mkdir -p $(DEB_DIR)/lib/systemd/system
	@# Build binary
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(DEB_DIR)/usr/bin/$(NAME) cmd/comi/main.go
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
	@mv $(DEB_DIR).deb $(BINDIR)/$(NAME)_$(VERSION)_$(DEB_ARCH).deb
	@rm -rf $(DEB_DIR)
	@echo "==> Created $(BINDIR)/$(NAME)_$(VERSION)_$(DEB_ARCH).deb"

# Build arm64 .deb package
deb-arm64:
	@echo "==> Building arm64 .deb package..."
	$(eval DEB_ARCH := arm64)
	$(eval DEB_DIR := $(BINDIR)/$(NAME)_$(VERSION)_$(DEB_ARCH))
	@mkdir -p $(DEB_DIR)/DEBIAN
	@mkdir -p $(DEB_DIR)/usr/bin
	@mkdir -p $(DEB_DIR)/lib/systemd/system
	@# Build binary
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(DEB_DIR)/usr/bin/$(NAME) cmd/comi/main.go
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
	@mv $(DEB_DIR).deb $(BINDIR)/$(NAME)_$(VERSION)_$(DEB_ARCH).deb
	@rm -rf $(DEB_DIR)
	@echo "==> Created $(BINDIR)/$(NAME)_$(VERSION)_$(DEB_ARCH).deb"

# Build all .deb packages
deb-all: deb-amd64 deb-arm64
	@echo "==> All .deb packages built successfully"

# Clean .deb build artifacts
deb-clean:
	@rm -f $(BINDIR)/*.deb
	@echo "==> Cleaned .deb packages"
