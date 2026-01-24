package scan

import (
	"os"
	"path/filepath"
	"testing"
)

// testCfgScan 仅用于扫描相关的单元测试
type testCfgScan struct {
	excludePath       []string
	supportMediaType  []string
	supportFileType   []string
	supportTemplate   []string
	maxScanDepth      int
	minImageNum       int
	timeoutLimitForSc int
}

func (c *testCfgScan) GetStoreUrls() []string           { return nil }
func (c *testCfgScan) GetMaxScanDepth() int             { return c.maxScanDepth }
func (c *testCfgScan) GetMinImageNum() int              { return c.minImageNum }
func (c *testCfgScan) GetTimeoutLimitForScan() int      { return c.timeoutLimitForSc }
func (c *testCfgScan) GetExcludePath() []string         { return c.excludePath }
func (c *testCfgScan) GetSupportMediaType() []string    { return c.supportMediaType }
func (c *testCfgScan) GetSupportFileType() []string     { return c.supportFileType }
func (c *testCfgScan) GetSupportTemplateFile() []string { return c.supportTemplate }
func (c *testCfgScan) GetZipFileTextEncoding() string   { return "utf-8" }
func (c *testCfgScan) GetEnableDatabase() bool          { return false }
func (c *testCfgScan) GetClearDatabaseWhenExit() bool   { return false }
func (c *testCfgScan) GetDebug() bool                   { return false }

// TestHandleDirectory_ShouldCollectSupportedFiles
// 回归用例：目录递归扫描时，应该能收集到支持的文件（例如 .zip），否则 InitStore 的“处理文件”阶段会漏扫。
func TestHandleDirectory_ShouldCollectSupportedFiles(t *testing.T) {
	tmp := t.TempDir()
	root := filepath.Join(tmp, "test")
	if err := os.MkdirAll(filepath.Join(root, "TestDir"), 0o755); err != nil {
		t.Fatalf("mkdir TestDir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "TestDir 2"), 0o755); err != nil {
		t.Fatalf("mkdir TestDir 2: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "TestDir3"), 0o755); err != nil {
		t.Fatalf("mkdir TestDir3: %v", err)
	}
	// 放一个假的 zip 文件（不需要有效内容；这里只验证 HandleDirectory 是否能发现它）
	zipPath := filepath.Join(root, "TestDir", "a.zip")
	if err := os.WriteFile(zipPath, []byte("not-a-real-zip"), 0o644); err != nil {
		t.Fatalf("write a.zip: %v", err)
	}

	InitConfig(&testCfgScan{
		excludePath:       []string{},
		supportMediaType:  []string{".jpg", ".png", ".webp"},
		supportFileType:   []string{".zip", ".cbz", ".rar", ".cbr", ".tar", ".epub", ".pdf"},
		supportTemplate:   []string{".html"},
		maxScanDepth:      -1,
		minImageNum:       1,
		timeoutLimitForSc: 0,
	})

	_, _, foundFiles, err := HandleDirectory(root, 0)
	if err != nil {
		t.Fatalf("HandleDirectory error: %v", err)
	}
	found := false
	for _, f := range foundFiles {
		if f.Path == zipPath {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected to find %s in foundFiles, but it was missing", zipPath)
	}
}
