package common

import (
	"net/url"
	"testing"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
)

// 两种下载入口都必须携带 BasePath 与远程书库参数。
func TestDownloadBookURL(t *testing.T) {
	saved := config.CopyCfg()
	defer config.ReplaceConfig(saved)
	config.GetCfg().BasePath = "/reader"
	for _, endpoint := range []string{"/api/download-zip", "/api/download-epub"} {
		for _, key := range []string{"", "remote/key+"} {
			book := &model.Book{BookInfo: model.BookInfo{BookID: "fixture", RemoteStoreKey: key}}
			u, err := url.Parse(downloadBookURL(book, endpoint))
			if err != nil {
				t.Fatal(err)
			}
			if u.Path != "/reader"+endpoint || u.Query().Get("id") != "fixture" || u.Query().Get("remote_store") != key {
				t.Fatalf("download URL=%s", u)
			}
		}
	}
}
