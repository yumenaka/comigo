// Package releases 为 CLI、托盘及状态 API 共享只读的版本查询逻辑。
package releases

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/mod/semver"
)

const LatestURL = "https://comigo.xyz/yumenaka/api.github.com/repos/yumenaka/comigo/releases/latest"

type Asset struct {
	Name string `json:"name"`
}
type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

// Canonical 只接受合法语义版本，开发构建不能误判为需要升级。
func Canonical(v string) string {
	v = strings.TrimSpace(v)
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	if !semver.IsValid(v) {
		return ""
	}
	return v
}
func Load(ctx context.Context, client *http.Client, current string) (*Release, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, LatestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Comigo-Upgrade/"+current)
	req.Header.Set("Accept", "application/vnd.github+json")
	// 发布代理可能丢失压缩响应头；小型 JSON 请求直接使用原始编码。
	req.Header.Set("Accept-Encoding", "identity")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("release HTTP %d", resp.StatusCode)
	}
	var rel Release
	if err := json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(&rel); err != nil {
		return nil, err
	}
	if strings.TrimSpace(rel.TagName) == "" {
		return nil, fmt.Errorf("empty release tag")
	}
	return &rel, nil
}

type Status struct {
	State         string     `json:"state"`
	LatestVersion string     `json:"latestVersion,omitempty"`
	Available     bool       `json:"available"`
	CheckedAt     *time.Time `json:"checkedAt,omitempty"`
	ReleaseURL    string     `json:"releaseURL,omitempty"`
}

func Compare(current string, rel *Release) Status {
	now := time.Now()
	s := Status{State: "checked", LatestVersion: rel.TagName, CheckedAt: &now,
		ReleaseURL: "https://github.com/yumenaka/comigo/releases/tag/" + url.PathEscape(rel.TagName)}
	a, b := Canonical(current), Canonical(rel.TagName)
	if a == "" || b == "" {
		s.State = "unknown"
	} else {
		s.Available = semver.Compare(b, a) > 0
	}
	return s
}

// Checker 的锁只保护快照；慢速联网不会阻塞普通 /api/server 请求。
type Checker struct {
	mu      sync.Mutex
	status  Status
	current string
	pending chan struct{}
	load    func(context.Context, string) (*Release, error)
}

func NewChecker() *Checker {
	return &Checker{status: Status{State: "unchecked"}, load: func(ctx context.Context, current string) (*Release, error) {
		return Load(ctx, &http.Client{Timeout: 15 * time.Second}, current)
	}}
}

var Default = NewChecker()

func (c *Checker) Snapshot() Status { c.mu.Lock(); defer c.mu.Unlock(); return c.status }
func (c *Checker) Check(ctx context.Context, current string) (Status, error) {
	c.mu.Lock()
	if c.current == current && c.status.CheckedAt != nil && time.Since(*c.status.CheckedAt) < time.Hour && c.status.State != "error" {
		s := c.status
		c.mu.Unlock()
		return s, nil
	}
	if pending := c.pending; pending != nil {
		c.mu.Unlock()
		select {
		case <-ctx.Done():
			return c.Snapshot(), ctx.Err()
		case <-pending:
		}
		s := c.Snapshot()
		if s.State == "error" {
			return s, fmt.Errorf("release check failed")
		}
		return s, nil
	}
	c.pending = make(chan struct{})
	c.mu.Unlock()
	rel, err := c.load(ctx, current)
	c.mu.Lock()
	defer c.mu.Unlock()
	if err != nil {
		now := time.Now()
		c.status.State = "error"
		c.status.CheckedAt = &now
	} else {
		c.status = Compare(current, rel)
		c.current = current
	}
	close(c.pending)
	c.pending = nil
	return c.status, err
}
