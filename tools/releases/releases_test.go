package releases

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestCompare(t *testing.T) {
	for _, tc := range []struct {
		a, b, state string
		available   bool
	}{
		{"1.2.3", "v1.2.4", "checked", true}, {"v2.0.0", "v1.9.9", "checked", false},
		{"1.2.3-rc.1", "1.2.3", "checked", true}, {"dev", "v1.0.0", "unknown", false},
		{"1.2.3", "bad", "unknown", false},
	} {
		s := Compare(tc.a, &Release{TagName: tc.b})
		if s.State != tc.state || s.Available != tc.available {
			t.Fatalf("%+v: %+v", tc, s)
		}
	}
}

// 合并查询，同时保证普通状态读取不等待上游请求。
func TestCheckerCacheAndFailure(t *testing.T) {
	c := NewChecker()
	var calls atomic.Int32
	entered := make(chan struct{})
	release := make(chan struct{})
	c.load = func(context.Context, string) (*Release, error) {
		calls.Add(1)
		close(entered)
		<-release
		return &Release{TagName: "v2.0.0"}, nil
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); c.Check(context.Background(), "1.0.0") }()
	<-entered
	if c.Snapshot().State != "unchecked" {
		t.Fatal("snapshot")
	}
	close(release)
	wg.Wait()
	for range 5 {
		if _, err := c.Check(context.Background(), "1.0.0"); err != nil {
			t.Fatal(err)
		}
	}
	if calls.Load() != 1 {
		t.Fatal("cache miss")
	}
	c.mu.Lock()
	old := time.Now().Add(-2 * time.Hour)
	c.status.CheckedAt = &old
	c.mu.Unlock()
	c.load = func(context.Context, string) (*Release, error) { return nil, errors.New("offline") }
	s, err := c.Check(context.Background(), "1.0.0")
	if err == nil || s.State != "error" {
		t.Fatal("failure not exposed")
	}
}

type rewriteTransport struct{ url string }

func (r rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	target, _ := http.NewRequest("GET", r.url, nil)
	clone.URL = target.URL
	return http.DefaultTransport.RoundTrip(clone)
}
func TestLoadValidation(t *testing.T) {
	for _, body := range []string{`{"tag_name":"v1.2.3"}`, `{}`, `invalid`} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Accept-Encoding") != "identity" {
				t.Error("release proxy must receive an uncompressed request")
			}
			w.Write([]byte(body))
		}))
		rel, err := Load(context.Background(), &http.Client{Transport: rewriteTransport{server.URL}}, "1.0.0")
		server.Close()
		if body == `{"tag_name":"v1.2.3"}` {
			if err != nil || rel.TagName != "v1.2.3" {
				t.Fatal(rel, err)
			}
		} else if err == nil {
			t.Fatal("accepted malformed release")
		}
	}
}
