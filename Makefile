# Makefile for cross-compilation
# Window icon Need：go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo

##Release:
# make all VERSION=v1.0.3

## Windows Release(Need MSYS2 or mingw32 + find.exe make.exe zip.exe upx.exe):
# mingw32-make all VERSION=v0.9.9

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
#| **Linux**   | ARM 32 位（树莓派 2~4，安装 32 位系统）      | `Linux_arm.tar.gz`   |
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
# -ldflags 指定编译参数。-s 去掉符号信息。 -w去掉调试信息。
GOBUILD=CGO_ENABLED=0 go build -ldflags "-s -w -X config.Version=${VERSION}"

ifeq ($(OS), Darwin)
  MD5_UTIL = md5
else
  MD5_UTIL = md5sum
endif

all: compileThemAll md5SumThemAll

# 因为sqlite（ent）库的关系，部分架构（Windows_i386）无法正常运行，需要写条件编译代码。但是最近似乎都Pass了，或许可以不再分架构：
# ent库的编译检测状态： https://modern-c.appspot.com/-/builder/?importpath=modernc.org%2Fsqlite

compileThemAll: Windows_x86_64 Windows_i386  Windows_arm64 Linux_x86_64 Linux_i386 Linux_arm Linux_arm64 MacOS_x86_64 MacOS_arm64

android: Linux_arm_android Linux_arm64-android

UPX := $(shell command -v upx 2> /dev/null)

gomobile:
	export ANDROID_NDK_HOME=/Users/bai/Library/Android/sdk/ndk/26.1.10909125
	gomobile bind -target=android -o comigo.aar -androidapi 26

md5SumThemAll:
	rm -f $(MD5_TEXTFILE)
	find $(BINDIR) -type f -name "$(NAME)_*" -exec $(MD5_UTIL) {} >> $(MD5_TEXTFILE) \;
	# 删除 $(MD5_TEXTFILE)里面的 ./bin/ 字符串
	sed -i '' 's|./bin/||g' $(MD5_TEXTFILE)
	cat $(MD5_TEXTFILE)

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

# 看ARM处理器是否支持VFP功能:grep -i vfp /proc/cpuinfo
##Linux_armv5,GOARM=5：使用软件浮点； CPU 没有 VFP 协处理器
#Linux_armv5:
#	GOARCH=arm GOOS=linux GOARM=5 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
##ifdef UPX
##	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
##endif
#	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
#	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Linux_armv6 RaspberryPi1,2,zero,GOARM=6：仅使用 VFPv1；交叉编译时默认；通常是 ARM11 或更好的内核（也支持 VFPv2 或更好的内核）
Linux_armv6:
	GOARCH=arm GOOS=linux GOARM=6 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
#ifdef UPX
#	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
#endif
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Linux_armv7，RaspberryPi3 官方32位armv7l系统。GOARM=7：使用 VFPv3；通常是 Cortex-A 内核. 2012年发布的架构。
Linux_arm:
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
	
