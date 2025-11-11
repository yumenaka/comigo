# Makefile for cross-compilation
# Window icon Need：go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo

##Release:
# make all VERSION=v1.1.0

## Windows Release(Need MSYS2 or mingw32 + find.exe make.exe zip.exe upx.exe):
# mingw32-make all VERSION=v0.9.9

## 打印编译命令，而不实际执行
# make -n Windows_i386_tailscale Windows_arm64_tailscale Linux_i386_tailscale

# make Windows_x86_64
# 版本号含义
# vMAJOR.MINOR.PATCH, 如 v1.2.3
# MAJOR：主版本号，不兼容的 API 修改
# MINOR：次版本号，向下兼容的功能性新增
# PATCH：修订号，向下兼容的问题修正

# 下一个版本是 v0.10.0

# 我应该下载哪个版本？
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

# 查看 golang 支持那些架构：
# go tool dist list

NAME=comi

OS := $(shell uname)
BINDIR := ./bin
MD5_TEXTFILE := $(BINDIR)/md5Sums.txt
#go: cannot install cross-compiled binaries when GOBIN is set
unexport GOBIN

MAIN_FILE_DIR := ./
# 覆盖版本号，只有 string 可以直接覆盖 包路径必须和 go list 输出的完全一致，需要带 module 路径
# -trimpath 去掉源码绝对路径，比如构建机目录
# -ldflags 指定编译参数。-s 去掉符号信息  -w去掉调试信息 减小二进制体积
GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X 'github.com/yumenaka/comigo/config.version=${VERSION}'"

ifeq ($(OS), Darwin)
  MD5_UTIL = md5
else
  MD5_UTIL = md5sum
endif

all: compileAll_CGO md5SumThemAll

## windows 可能不需要CGO就能支持Tailscale？

# 过去因为sqlite（ent）库的关系，部分架构（Windows_i386）无法正常运行，需要写条件编译代码。但是最近似乎都Pass了，或许可以不再分架构：
# ent库的编译检测状态： https://modern-c.appspot.com/-/builder/?importpath=modernc.org%2Fsqlite
# 为了支持Tailscale，使用docker交叉编译
compileAll: Windows_x86_64 Windows_i386  Windows_arm64 Linux_x86_64 Linux_i386 Linux_armv7 Linux_arm64 MacOS_x86_64 MacOS_arm64
compileAll_CGO: Windows_x86_64_cgo Windows_i386_cgo  Windows_arm64_cgo Linux_x86_64_cgo Linux_i386_cgo Linux_armv7_cgo Linux_arm64_cgo MacOS_x86_64_cgo MacOS_arm64_cgo

android: Linux_arm_android Linux_arm64-android

UPX := $(shell command -v upx 2> /dev/null)
DOCKER := $(shell command -v docker 2> /dev/null)

gomobile:
	export ANDROID_NDK_HOME=/Users/bai/Library/Android/sdk/ndk/26.1.10909125
	gomobile bind -target=android -o comigo.aar -androidapi 26

md5SumThemAll:
	rm -f $(MD5_TEXTFILE)
	find $(BINDIR) -type f -name "$(NAME)_*" -exec $(MD5_UTIL) {} >> $(MD5_TEXTFILE) \;
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
#镜像支持的参数：告诉它用什么命令来“构建”你的项目。
#这里执行的是 make build（你项目的 Makefile 中应定义了 build 目标）。

## 查看编译用Docker镜像信息：
# https://github.com/elastic/golang-crossbuild/releases

# make Windows_x86_64_cgo VERSION=v1.0.5
Windows_x86_64_cgo:
ifndef DOCKER
	$(error "No docker found! Please install docker to build Windows_x86_64_cgo")
endif
ifdef DOCKER
	docker run -it  \
	 -v "$$PWD":/go/src/github.com/user/go-project \
	 -w /go/src/github.com/user/go-project \
	 -e CGO_ENABLED=1 \
	 -e VERSION=$(VERSION) \
	 -e FILE_LABLE="Windows_x86_64" \
	 docker.elastic.co/beats-dev/golang-crossbuild:1.25.4-main-debian7 \
	 --build-cmd "make windows_x86_64_cgo_docker VERSION=$(VERSION)" \
	 -p "windows/amd64"
endif

windows_x86_64_cgo_docker:
	cp resource.syso.windows_amd64 resource.syso
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w -X 'github.com/yumenaka/comigo/config.version=${VERSION}'" -o $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)/$(NAME).exe ./cmd/comi
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE).tar.gz $(NAME).exe
	rm -rf  $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)
	rm   resource.syso

#make Windows_i386_cgo VERSION=v1.0.5
Windows_i386_cgo:
ifndef DOCKER
	$(error "No docker found! Please install docker to build Windows_i386_cgo")
endif
ifdef DOCKER
	docker run -it  \
	 -v "$$PWD":/go/src/github.com/user/go-project \
	 -w /go/src/github.com/user/go-project \
	 -e CGO_ENABLED=1 \
	 -e VERSION=$(VERSION) \
	 -e FILE_LABLE="Windows_i386" \
	 docker.elastic.co/beats-dev/golang-crossbuild:1.25.4-main-debian7 \
	 --build-cmd "make windows_i386_cgo_docker VERSION=$(VERSION)" \
	 -p "windows/386"
endif

windows_i386_cgo_docker:
	cp resource.syso.windows_386 resource.syso
	CGO_ENABLED=1 GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)/$(NAME).exe ./cmd/comi
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE).tar.gz $(NAME).exe
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)
	rm   resource.syso

#make Windows_arm64_cgo VERSION=v1.0.5
Windows_arm64_cgo:
ifndef DOCKER
	$(error "No docker found! Please install docker to build Windows_arm64_cgo")
endif
ifdef DOCKER
	docker run -it  \
	 -v "$$PWD":/go/src/github.com/user/go-project \
	 -w /go/src/github.com/user/go-project \
	 -e CGO_ENABLED=1 \
	 -e VERSION=$(VERSION) \
	 -e FILE_LABLE="Windows_arm64" \
	 docker.elastic.co/beats-dev/golang-crossbuild:1.25.4-windows-arm64-debian12 \
	 --build-cmd "make windows_arm64_cgo_docker VERSION=$(VERSION)" \
	 -p "windows/arm64"
endif

windows_arm64_cgo_docker:
	CGO_ENABLED=1 GOOS=windows GOARCH=arm64 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)/$(NAME).exe ./cmd/comi
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE).tar.gz $(NAME).exe
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)

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
	 docker.elastic.co/beats-dev/golang-crossbuild:1.25.4-armhf-debian9 \
	 --build-cmd "make linux_armv7_cgo_docker VERSION=$(VERSION)" \
	 -p "linux/armv7"
endif

linux_armv7_cgo_docker:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)/$(NAME) ./cmd/comi
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
	 docker.elastic.co/beats-dev/golang-crossbuild:1.25.4-darwin \
	 --build-cmd "make darwin_x86_64_cgo_docker VERSION=$(VERSION)" \
	 -p "darwin/amd64"
endif

darwin_x86_64_cgo_docker:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)/$(NAME) ./cmd/comi
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
	 docker.elastic.co/beats-dev/golang-crossbuild:1.25.4-darwin-arm64-debian10 \
	 --build-cmd "make darwin_arm64_cgo_docker VERSION=$(VERSION)" \
	 -p "darwin/arm64"
endif

darwin_arm64_cgo_docker:
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)/$(NAME) ./cmd/comi
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
	 docker.elastic.co/beats-dev/golang-crossbuild:1.25.4-base-arm-debian9 \
	 --build-cmd "make linux_arm64_cgo_docker VERSION=$(VERSION)" \
	 -p "linux/arm64"
endif

linux_arm64_cgo_docker:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)/$(NAME) ./cmd/comi
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
	 docker.elastic.co/beats-dev/golang-crossbuild:1.25.4-main-debian7 \
	 --build-cmd "make linux_x86_64_cgo_docker VERSION=$(VERSION)" \
	 -p "linux/amd64"
endif

linux_x86_64_cgo_docker:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)/$(NAME) ./cmd/comi
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
	 docker.elastic.co/beats-dev/golang-crossbuild:1.25.4-main-debian7 \
	 --build-cmd "make linux_i386_cgo_docker VERSION=$(VERSION)" \
	 -p "linux/386"
endif

linux_i386_cgo_docker:
	CGO_ENABLED=1 GOOS=linux GOARCH=386 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)/$(NAME) ./cmd/comi
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE).tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$(FILE_LABLE)


## No CGO build:

# upx可能导致报毒，取消windows平台的upx压缩
# 换行用TAB而不是空格
#64位Windows	$(NAME)_$(VERSION)_$@   
Windows_x86_64:
	go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo # Window icon Need
	GOARCH=amd64 GOOS=windows go generate #go: cannot install cross-compiled binaries when GOBIN is set
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME).exe 
	zip -m -r -j -9 $(BINDIR)/$(NAME)_$(VERSION)_$@.zip $(BINDIR)/$(NAME)_$(VERSION)_$@
	rmdir $(BINDIR)/$(NAME)_$(VERSION)_$@
	rm  resource.syso 

#32位Windows	
Windows_i386:
	go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo # Window icon Need
	GOARCH=386 GOOS=windows go generate
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME).exe 
	zip -m -r -j -9 $(BINDIR)/$(NAME)_$(VERSION)_$@.zip $(BINDIR)/$(NAME)_$(VERSION)_$@
	rmdir $(BINDIR)/$(NAME)_$(VERSION)_$@
	rm   resource.syso	

#windows arm64 no upx
Windows_arm64:
	GOARCH=arm64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME).exe 
	zip -m -r -j -9 $(BINDIR)/$(NAME)_$(VERSION)_$@.zip $(BINDIR)/$(NAME)_$(VERSION)_$@
	rmdir $(BINDIR)/$(NAME)_$(VERSION)_$@
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@


#Linux_armv6 RaspberryPi1,2,zero,GOARM=6：仅使用 VFPv1；交叉编译时默认；通常是 ARM11 或更好的内核（也支持 VFPv2 或更好的内核）
Linux_armv6:
	GOARCH=arm GOOS=linux GOARM=6 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
#ifdef UPX
#	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
#endif
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Linux_armv7，RaspberryPi3 官方32位armv7l系统。GOARM=7：使用 VFPv3；通常是 Cortex-A 内核. 2012年发布的架构。
Linux_armv7:
	GOARCH=arm GOOS=linux GOARM=7 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
#ifdef UPX
#	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
#endif
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#linux，64位arm。2012年发布的架构。
Linux_arm64:
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
#ifdef UPX
#	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
#endif
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Linux，x86_64
Linux_x86_64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
#ifdef UPX
#	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
#endif
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Linux，i386
Linux_i386:
	GOARCH=386 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
#ifdef UPX
#	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
#endif
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#MACOS x86_64
MacOS_x86_64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
#ifdef UPX
#	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
#endif
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@
	
#MACOS arm64 no upx
MacOS_arm64:
	GOARCH=arm64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@
	
#Android，32位arm，Termux	
Linux_arm_android:
	GOARCH=arm GOOS=android $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Android，64位arm，Termux
Linux_arm64-android:
	GOARCH=arm64 GOOS=android $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@
	
