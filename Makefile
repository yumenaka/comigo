# Makefile for cross-compilation
# make all VERSION=v0.5.0
# make md5SumThemAll VERSION=v0.5.0

# mingw32-make all VERSION=v0.4.5
# need MSYS2 or mingw32 or find.exe make.exe zip.exe upx.exe
NAME=comi
SKETCH_NAME=sketch_66seconds
OS := $(shell uname)
BINDIR := ./bin
MD5_TEXTFILE := $(BINDIR)/md5Sums.txt

MAIN_FILE_DIR := ./
# -ldflags 指定编译参数。-s 去掉符号信息。 -w去掉调试信息。
GOBUILD=CGO_ENABLED=0 go build -ldflags "-s -w -X $controllers.Version=$VERSION" 

ifeq ($(OS), Darwin)
  MD5_UTIL = md5
else
  MD5_UTIL = md5sum
endif

all: compileThemAll md5SumThemAll

compileThemAll: Windows_x86_64 Windows_i386 Linux_x86_64 Linux_i386 MacOS_x86_64 MacOS_arm64 Linux-armv5 Linux-armv7 Linux-armv6 Linux-armv8 

android: Linux-arm-android Linux-arm64-android

UPX := $(shell command -v upx 2> /dev/null)

md5SumThemAll:
	rm -f $(MD5_TEXTFILE)
	find $(BINDIR) -type f -name "$(NAME)_*" -exec $(MD5_UTIL) {} >> $(MD5_TEXTFILE) \;
	cat $(MD5_TEXTFILE)

#64位Windows	$(NAME)_$(VERSION)_$@   
Windows_x86_64:
	go generate
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME).exe 
ifdef UPX
	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME).exe 
endif
	zip -m -r -j -9 $(BINDIR)/$(NAME)_$(VERSION)_$@.zip $(BINDIR)/$(NAME)_$(VERSION)_$@
	rmdir $(BINDIR)/$(NAME)_$(VERSION)_$@
	rm  resource.syso 

#32位Windows	
Windows_i386:
	go generate
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME).exe 
ifdef UPX
	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME).exe 
endif
	zip -m -r -j -9 $(BINDIR)/$(NAME)_$(VERSION)_$@.zip $(BINDIR)/$(NAME)_$(VERSION)_$@
	rmdir $(BINDIR)/$(NAME)_$(VERSION)_$@
	rm   resource.syso	

# golang支持的交叉编译架构的列表，参见 go tool dist list
# 看ARM处理器是否支持VFP功能:grep -i vfp /proc/cpuinfo

#Linux-armv5,GOARM=5：使用软件浮点；当 CPU 没有 VFP 协处理器时
Linux-armv5:
	GOARCH=arm GOOS=linux GOARM=5 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
ifdef UPX
	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
endif
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Linux-armv6 RaspberryPi1,2,zero,GOARM=6：仅使用 VFPv1；交叉编译时默认；通常是 ARM11 或更好的内核（也支持 VFPv2 或更好的内核）
Linux-armv6:
	GOARCH=arm GOOS=linux GOARM=6 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
ifdef UPX
	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
endif
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Linux-armv7，RaspberryPi3 官方32位armv7l系统。GOARM=7：使用 VFPv3；通常是 Cortex-A 内核.
Linux-armv7:
	GOARCH=arm GOOS=linux GOARM=7 $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
ifdef UPX
	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
endif
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#linux，64位arm
Linux-armv8:
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
ifdef UPX
	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
endif
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Linux，x86_64
Linux_x86_64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
ifdef UPX
	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
endif
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Linux，i386
Linux_i386:
	GOARCH=386 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
ifdef UPX
	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
endif
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#MACOS x86_64
MacOS_x86_64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
ifdef UPX
	upx -9 $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
endif
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@
	
#MACOS arm64 no upx
MacOS_arm64:
	GOARCH=arm64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME)
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Android，32位arm，Termux	
Linux-arm-android:
	GOARCH=arm GOOS=android $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@

#Android，64位arm，Termux
Linux-arm64-android:
	GOARCH=arm64 GOOS=android $(GOBUILD) -o $(BINDIR)/$(NAME)_$(VERSION)_$@/$(NAME) 
	tar --directory=$(BINDIR)/$(NAME)_$(VERSION)_$@  -zcvf $(BINDIR)/$(NAME)_$(VERSION)_$@.tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)_$(VERSION)_$@
	
