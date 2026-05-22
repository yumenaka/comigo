package reverse_proxy

import (
	"strings"
	"testing"

	"github.com/yumenaka/comigo/config"
)

func TestShouldReplaceLatestWithVersionOnlyForComigoRelease(t *testing.T) {
	version := config.GetVersion()
	target := "https://github.com/yumenaka/comigo/releases/download/latest/comi_latest_MacOS_arm64.tar.gz"

	if !shouldReplaceLatestWithVersion(target) {
		t.Fatal("Comigo 官方 release 下载地址应该允许 latest 替换")
	}

	replaced := replaceLatestWithVersion(target)
	if !strings.Contains(replaced, "/releases/download/"+version+"/") {
		t.Fatalf("release 路径未替换为当前版本: %s", replaced)
	}
	if !strings.Contains(replaced, "comi_"+version+"_MacOS_arm64.tar.gz") {
		t.Fatalf("文件名未替换为当前版本: %s", replaced)
	}
}

func TestShouldNotReplaceLatestWithVersionForOtherRepositories(t *testing.T) {
	targets := []string{
		"https://github.com/yumenaka/other/releases/download/latest/comi_latest_MacOS_arm64.tar.gz",
		"https://github.com/other/comigo/releases/download/latest/comi_latest_MacOS_arm64.tar.gz",
		"https://raw.githubusercontent.com/yumenaka/comigo/master/latest.txt",
		"https://api.github.com/repos/yumenaka/comigo/releases/latest",
	}

	for _, target := range targets {
		if shouldReplaceLatestWithVersion(target) {
			t.Fatalf("非 github.com/yumenaka/comigo release 下载地址不应替换 latest: %s", target)
		}
	}
}
