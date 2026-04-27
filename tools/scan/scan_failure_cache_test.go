package scan

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yumenaka/comigo/config"
)

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
