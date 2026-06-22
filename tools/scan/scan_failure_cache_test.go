package scan

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
)

// 验证扫描失败缓存只跳过同版本且指纹一致的压缩包。
func TestScanFailureCacheSkipsOnlySameVersionAndFingerprint(t *testing.T) {
	t.Setenv("COMIGO_CONFIG_DIR", t.TempDir())

	storeURL := "/books"
	filePath := filepath.Join(storeURL, "bad.zip")
	modTime := time.Unix(1700000000, 123)
	size := int64(42)

	recordArchiveScanFailure(storeURL, filePath, size, modTime, false, errors.New("broken archive"))

	if !shouldSkipFailedArchiveFile(storeURL, filePath, size, modTime, false) {
		t.Fatalf("expected unchanged failed archive to be skipped")
	}
	if shouldSkipFailedArchiveFile(storeURL, filePath, size+1, modTime, false) {
		t.Fatalf("expected changed file size to allow retry")
	}
	if shouldSkipFailedArchiveFile(storeURL, filePath, size, modTime.Add(time.Second), false) {
		t.Fatalf("expected changed modified time to allow retry")
	}
}

// 验证版本变化后会重新尝试曾经失败的压缩包。
func TestScanFailureCacheVersionChangeAllowsRetry(t *testing.T) {
	t.Setenv("COMIGO_CONFIG_DIR", t.TempDir())

	storeURL := "/books"
	filePath := filepath.Join(storeURL, "bad.cbz")
	modTime := time.Unix(1700000000, 0)
	size := int64(99)

	cache := scanFailureCache{
		scanFailureRecordKey(storeURL, filePath, false): {
			StoreURL:         storeURL,
			FilePath:         filePath,
			FileSize:         size,
			ModifiedUnixNano: modTime.UnixNano(),
			CreatedByVersion: config.GetVersion() + "-old",
			FailedAt:         time.Now(),
			Error:            "old failure",
			IsRemote:         false,
		},
	}
	if err := saveScanFailureCache(cache); err != nil {
		t.Fatalf("saveScanFailureCache: %v", err)
	}

	if shouldSkipFailedArchiveFile(storeURL, filePath, size, modTime, false) {
		t.Fatalf("expected version mismatch to allow retry")
	}
}

// 验证小补丁版本变化不会绕过扫描失败缓存。
func TestScanFailureCacheSmallPatchVersionChangeStillSkips(t *testing.T) {
	t.Setenv("COMIGO_CONFIG_DIR", t.TempDir())

	storeURL := "/books"
	filePath := filepath.Join(storeURL, "bad.cbz")
	modTime := time.Unix(1700000000, 0)
	size := int64(99)

	cache := scanFailureCache{
		scanFailureRecordKey(storeURL, filePath, false): {
			StoreURL:         storeURL,
			FilePath:         filePath,
			FileSize:         size,
			ModifiedUnixNano: modTime.UnixNano(),
			CreatedByVersion: versionWithPatchDeltaForTest(t, -1),
			FailedAt:         time.Now(),
			Error:            "old patch failure",
			IsRemote:         false,
		},
	}
	if err := saveScanFailureCache(cache); err != nil {
		t.Fatalf("saveScanFailureCache: %v", err)
	}

	if !shouldSkipFailedArchiveFile(storeURL, filePath, size, modTime, false) {
		t.Fatalf("expected small patch version change to keep skipping failed archive")
	}
}

// 验证失败压缩包是否重试取决于记录版本和当前版本。
func TestShouldRetryFailedArchiveByVersion(t *testing.T) {
	tests := []struct {
		name           string
		recordVersion  string
		currentVersion string
		wantRetry      bool
	}{
		{name: "same version", recordVersion: "v1.2.36", currentVersion: "v1.2.36", wantRetry: false},
		{name: "small patch increase", recordVersion: "v1.2.35", currentVersion: "v1.2.36", wantRetry: false},
		{name: "patch increase equal ten", recordVersion: "v1.2.26", currentVersion: "v1.2.36", wantRetry: false},
		{name: "patch increase greater than ten", recordVersion: "v1.2.25", currentVersion: "v1.2.36", wantRetry: true},
		{name: "minor changes", recordVersion: "v1.1.36", currentVersion: "v1.2.36", wantRetry: true},
		{name: "current patch lower", recordVersion: "v1.2.40", currentVersion: "v1.2.36", wantRetry: true},
		{name: "unparseable old version", recordVersion: "v1.2.36-old", currentVersion: "v1.2.36", wantRetry: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldRetryFailedArchiveByVersion(tt.recordVersion, tt.currentVersion); got != tt.wantRetry {
				t.Fatalf("shouldRetryFailedArchiveByVersion(%q, %q) = %v, want %v", tt.recordVersion, tt.currentVersion, got, tt.wantRetry)
			}
		})
	}
}

// 验证失败缓存只作用于压缩包，并会在成功后清除记录。
func TestScanFailureCacheIgnoresNonArchiveFilesAndClearsOnSuccess(t *testing.T) {
	t.Setenv("COMIGO_CONFIG_DIR", t.TempDir())

	storeURL := "/books"
	archivePath := filepath.Join(storeURL, "bad.rar")
	pdfPath := filepath.Join(storeURL, "bad.pdf")
	modTime := time.Unix(1700000000, 0)

	recordArchiveScanFailure(storeURL, pdfPath, 12, modTime, false, errors.New("pdf failure"))
	cachePath, err := scanFailureCachePath()
	if err != nil {
		t.Fatalf("scanFailureCachePath: %v", err)
	}
	if _, err := os.Stat(cachePath); !os.IsNotExist(err) {
		t.Fatalf("expected non-archive failure not to create cache file, got err=%v", err)
	}

	recordArchiveScanFailure(storeURL, archivePath, 12, modTime, false, errors.New("archive failure"))
	if !shouldSkipFailedArchiveFile(storeURL, archivePath, 12, modTime, false) {
		t.Fatalf("expected archive failure to be cached")
	}

	clearArchiveScanFailure(storeURL, archivePath, false)
	if shouldSkipFailedArchiveFile(storeURL, archivePath, 12, modTime, false) {
		t.Fatalf("expected cleared archive failure to allow retry")
	}
}

// 验证初始化书库时会跳过已记录失败且无需重试的压缩包。
func TestInitStoreSkipsPreviouslyFailedArchiveFiles(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("COMIGO_CONFIG_DIR", filepath.Join(tmp, "config"))

	library := filepath.Join(tmp, "library")
	if err := os.MkdirAll(library, 0o755); err != nil {
		t.Fatalf("mkdir library: %v", err)
	}
	// 这里动态创建损坏压缩包，避免依赖仓库 test/ 下的提示文件。
	if err := os.WriteFile(filepath.Join(library, "broken.rar"), []byte("not a rar archive"), 0o644); err != nil {
		t.Fatalf("write broken.rar: %v", err)
	}
	if err := os.WriteFile(filepath.Join(library, "fake.zip"), []byte("not a zip archive"), 0o644); err != nil {
		t.Fatalf("write fake.zip: %v", err)
	}

	oldStore := model.IStore
	oldCfg := config.CopyCfg()
	model.IStore = &store.StoreInRam{}
	*config.GetCfg() = oldCfg
	config.GetCfg().ConfigFile = ""
	config.GetCfg().MinImageNum = 1
	t.Cleanup(func() {
		model.IStore = oldStore
		*config.GetCfg() = oldCfg
	})

	scanCfg := &testCfgScan{
		excludePath:       []string{},
		supportMediaType:  []string{".jpg", ".png", ".webp"},
		supportFileType:   []string{".zip", ".cbz", ".rar", ".cbr", ".tar", ".epub", ".pdf"},
		supportTemplate:   []string{".html"},
		maxScanDepth:      -1,
		minImageNum:       1,
		timeoutLimitForSc: 0,
	}
	if err := InitStore(library, scanCfg); err != nil {
		t.Fatalf("initial InitStore: %v", err)
	}

	firstCache := loadScanFailureCacheForTest(t)
	if len(firstCache) != 2 {
		t.Fatalf("first scan failure cache length = %d, want 2: %+v", len(firstCache), firstCache)
	}
	firstFailedAt := failedAtByBaseName(firstCache)

	if err := InitStore(library, scanCfg); err != nil {
		t.Fatalf("second InitStore: %v", err)
	}

	secondCache := loadScanFailureCacheForTest(t)
	if len(secondCache) != 2 {
		t.Fatalf("second scan failure cache length = %d, want 2: %+v", len(secondCache), secondCache)
	}
	secondFailedAt := failedAtByBaseName(secondCache)
	for name, firstTime := range firstFailedAt {
		secondTime, ok := secondFailedAt[name]
		if !ok {
			t.Fatalf("missing failure record for %s after second scan", name)
		}
		if !secondTime.Equal(firstTime) {
			t.Fatalf("failure record for %s was rewritten: first=%v second=%v", name, firstTime, secondTime)
		}
	}
}

func loadScanFailureCacheForTest(t *testing.T) scanFailureCache {
	t.Helper()
	cache, err := loadScanFailureCache()
	if err != nil {
		t.Fatalf("loadScanFailureCache: %v", err)
	}
	return cache
}

func failedAtByBaseName(cache scanFailureCache) map[string]time.Time {
	result := make(map[string]time.Time, len(cache))
	for _, record := range cache {
		result[filepath.Base(record.FilePath)] = record.FailedAt
	}
	return result
}

func versionWithPatchDeltaForTest(t *testing.T, delta int) string {
	t.Helper()
	major, minor, patch, ok := parseSemanticVersion(config.GetVersion())
	if !ok {
		t.Fatalf("cannot parse current version %q", config.GetVersion())
	}
	return "v" + strconv.Itoa(major) + "." + strconv.Itoa(minor) + "." + strconv.Itoa(patch+delta)
}
