package routers

import (
	"slices"
	"testing"

	"golang.org/x/crypto/acme"
)

func TestAutoTLSNextProtosIncludesHTTPSProtocols(t *testing.T) {
	t.Helper()

	got := autoTLSNextProtos()
	want := []string{acme.ALPNProto, "h2", "http/1.1"}

	for _, proto := range want {
		if !slices.Contains(got, proto) {
			t.Fatalf("autoTLSNextProtos() 缺少协议 %q，当前值: %v", proto, got)
		}
	}
}
