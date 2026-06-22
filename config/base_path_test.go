package config

import "testing"

// 验证反向代理基础路径会被规范化为统一格式。
func TestNormalizeBasePath(t *testing.T) {
	tests := map[string]string{
		"":               "",
		"/":              "",
		"proxy":          "/proxy",
		"/proxy/":        "/proxy",
		"/some/path///":  "/some/path",
		"//some//path//": "/some/path",
	}
	for input, want := range tests {
		if got := NormalizeBasePath(input); got != want {
			t.Fatalf("NormalizeBasePath(%q) = %q, want %q", input, got, want)
		}
	}
}

// 验证路径加前缀和去前缀在根路径、子路径下都保持一致。
func TestPrefixAndStripBasePath(t *testing.T) {
	oldBasePath := cfg.BasePath
	t.Cleanup(func() { cfg.BasePath = oldBasePath })

	cfg.BasePath = "/some/path/"
	if got := PrefixPath("/api/get-book"); got != "/some/path/api/get-book" {
		t.Fatalf("PrefixPath API = %q", got)
	}
	if got := PrefixPath("/"); got != "/some/path/" {
		t.Fatalf("PrefixPath root = %q", got)
	}
	if got := PrefixPath("/some/path/"); got != "/some/path/" {
		t.Fatalf("PrefixPath prefixed root should be idempotent, got %q", got)
	}
	if got := PrefixPath("/some/path/settings"); got != "/some/path/settings" {
		t.Fatalf("PrefixPath prefixed nested path should be idempotent, got %q", got)
	}
	if got := StripBasePath("/some/path/flip/book-id"); got != "/flip/book-id" {
		t.Fatalf("StripBasePath nested = %q", got)
	}
	if got := StripBasePath("/some/path"); got != "/" {
		t.Fatalf("StripBasePath base = %q", got)
	}

	cfg.BasePath = ""
	if got := PrefixPath("/api/get-book"); got != "/api/get-book" {
		t.Fatalf("PrefixPath without base = %q", got)
	}
}

// 验证通过通用配置入口设置基础路径时也会自动规范化。
func TestSetConfigValueNormalizesBasePath(t *testing.T) {
	testCfg := newDefaultConfig()
	if err := testCfg.SetConfigValue("BasePath", "some/path/"); err != nil {
		t.Fatalf("SetConfigValue(BasePath) error = %v", err)
	}
	if testCfg.BasePath != "/some/path" {
		t.Fatalf("BasePath = %q, want /some/path", testCfg.BasePath)
	}
}

// 验证二维码基础地址包含反向代理基础路径。
func TestGetQrcodeURLIncludesBasePath(t *testing.T) {
	oldCfg := cfg
	t.Cleanup(func() { cfg = oldCfg })

	cfg = newDefaultConfig()
	cfg.Host = "example.com"
	cfg.Port = 1234
	cfg.BasePath = "/proxy"

	if got, want := GetQrcodeURL(), "http://example.com:1234/proxy/"; got != want {
		t.Fatalf("GetQrcodeURL() = %q, want %q", got, want)
	}
}

// 验证二维码公开地址只替换本机回环地址，不改写任意文本。
func TestToQrcodePublicURLRewritesLoopbackOnly(t *testing.T) {
	oldCfg := cfg
	t.Cleanup(func() { cfg = oldCfg })

	cfg = newDefaultConfig()
	cfg.Host = "example.com"
	cfg.Port = 1234

	got := ToQrcodePublicURL("http://127.0.0.1:1234/flip/book-id?page=2")
	want := "http://example.com:1234/flip/book-id?page=2"
	if got != want {
		t.Fatalf("ToQrcodePublicURL(loopback) = %q, want %q", got, want)
	}

	rawText := "just text"
	if got := ToQrcodePublicURL(rawText); got != rawText {
		t.Fatalf("ToQrcodePublicURL(text) = %q, want %q", got, rawText)
	}
}
