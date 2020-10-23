# Makefile for cross-compilation
# make all VERSION=v0.2.3
# mingw32-make all VERSION=v0.2.3
NAME=comi
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

compileThemAll: windows-amd64 windows-386 linux-armv7 linux-armv8 linux-amd64  darwin-amd64 

UPX := $(shell command -v upx 2> /dev/null)

md5SumThemAll:
	rm -f $(MD5_TEXTFILE)
	find $(BINDIR) -type f -name "$(NAME)-*" -exec $(MD5_UTIL) {} >> $(MD5_TEXTFILE) \;
	cat $(MD5_TEXTFILE)

#64位Windows	
windows-amd64:
	go generate
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@-$(VERSION)/$(NAME).exe 
ifdef UPX
	upx -9 $(BINDIR)/$(NAME)-$@-$(VERSION)/$(NAME).exe 
endif
	zip -m -r -j -9 $(BINDIR)/$(NAME)-$@-$(VERSION).zip $(BINDIR)/$(NAME)-$@-$(VERSION)
	rmdir $(BINDIR)/$(NAME)-$@-$(VERSION)
	rm  resource.syso 

#32位Windows	
windows-386:
	go generate
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@-$(VERSION)/$(NAME).exe 
ifdef UPX
	upx -9 $(BINDIR)/$(NAME)-$@-$(VERSION)/$(NAME).exe 
endif
	zip -m -r -j -9 $(BINDIR)/$(NAME)-$@-$(VERSION).zip $(BINDIR)/$(NAME)-$@-$(VERSION)
	rmdir $(BINDIR)/$(NAME)-$@-$(VERSION)
	rm   resource.syso	
	
#32位arm，比如树莓派	
linux-armv7:
	GOARCH=arm GOOS=linux GOARM=7 $(GOBUILD) -o $(BINDIR)/$(NAME)-$@-$(VERSION)/$(NAME) 
ifdef UPX
	upx -9 $(BINDIR)/$(NAME)-$@-$(VERSION)/$(NAME)
endif
	tar --directory=$(BINDIR)/$(NAME)-$@-$(VERSION)  -zcvf $(BINDIR)/$(NAME)-$@-$(VERSION).tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)-$@-$(VERSION)

#64位arm，如今大多数手机	
linux-armv8:
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@-$(VERSION)/$(NAME) 
ifdef UPX
	upx -9 $(BINDIR)/$(NAME)-$@-$(VERSION)/$(NAME)
endif
	tar --directory=$(BINDIR)/$(NAME)-$@-$(VERSION)  -zcvf $(BINDIR)/$(NAME)-$@-$(VERSION).tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)-$@-$(VERSION)

#64位Linux
linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@-$(VERSION)/$(NAME) 
ifdef UPX
	upx -9 $(BINDIR)/$(NAME)-$@-$(VERSION)/$(NAME)
endif
	tar --directory=$(BINDIR)/$(NAME)-$@-$(VERSION)  -zcvf $(BINDIR)/$(NAME)-$@-$(VERSION).tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)-$@-$(VERSION)
	
#64位MACOS
darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@-$(VERSION)/$(NAME)
ifdef UPX
	upx -9 $(BINDIR)/$(NAME)-$@-$(VERSION)/$(NAME)
endif
	tar --directory=$(BINDIR)/$(NAME)-$@-$(VERSION)  -zcvf $(BINDIR)/$(NAME)-$@-$(VERSION).tar.gz $(NAME)
	rm -rf $(BINDIR)/$(NAME)-$@-$(VERSION)

