package config

import "testing"

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

func TestSetConfigValueNormalizesBasePath(t *testing.T) {
	testCfg := newDefaultConfig()
	if err := testCfg.SetConfigValue("BasePath", "some/path/"); err != nil {
		t.Fatalf("SetConfigValue(BasePath) error = %v", err)
	}
	if testCfg.BasePath != "/some/path" {
		t.Fatalf("BasePath = %q, want /some/path", testCfg.BasePath)
	}
}

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
