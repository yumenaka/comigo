package scan

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

const scanFailureCacheFileName = "scan_failures.json"

var scanFailureCacheMu sync.Mutex

// ScanFailureRecord 记录压缩文件扫描失败时的文件指纹。
// 后续扫描仅在文件变化，或版本跨度足够大时才会再次尝试。
type ScanFailureRecord struct {
	StoreURL         string    `json:"store_url"`
	FilePath         string    `json:"file_path"`
	FileSize         int64     `json:"file_size"`
	ModifiedUnixNano int64     `json:"modified_unix_nano"`
	CreatedByVersion string    `json:"created_by_version"`
	FailedAt         time.Time `json:"failed_at"`
	Error            string    `json:"error"`
	IsRemote         bool      `json:"is_remote"`
}

type scanFailureCache map[string]ScanFailureRecord

func isArchiveScanFailureTarget(filePath string) bool {
	switch model.GetBookTypeByFilename(filePath) {
	case model.TypeZip, model.TypeCbz, model.TypeRar, model.TypeCbr, model.TypeTar, model.TypeEpub:
		return true
	default:
		return false
	}
}

func scanFailureRecordKey(storeURL, filePath string, isRemote bool) string {
	return tools.Md5string(storeURL + "\x00" + filePath + "\x00" + boolString(isRemote))
}

func boolString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func scanFailureCachePath() (string, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "metadata", scanFailureCacheFileName), nil
}

func loadScanFailureCache() (scanFailureCache, error) {
	cachePath, err := scanFailureCachePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(cachePath)
	if os.IsNotExist(err) {
		return scanFailureCache{}, nil
	}
	if err != nil {
		return nil, err
	}
	var cache scanFailureCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	if cache == nil {
		cache = scanFailureCache{}
	}
	return cache, nil
}

func saveScanFailureCache(cache scanFailureCache) error {
	cachePath, err := scanFailureCachePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(cachePath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cachePath, data, 0o644)
}

func shouldSkipFailedArchiveFile(storeURL, filePath string, size int64, modTime time.Time, isRemote bool) bool {
	if !isArchiveScanFailureTarget(filePath) {
		return false
	}

	scanFailureCacheMu.Lock()
	defer scanFailureCacheMu.Unlock()

	cache, err := loadScanFailureCache()
	if err != nil {
		logger.Infof(locale.GetString("log_scan_failure_cache_load_failed"), err)
		return false
	}

	record, ok := cache[scanFailureRecordKey(storeURL, filePath, isRemote)]
	if !ok {
		return false
	}

	if shouldRetryFailedArchiveByVersion(record.CreatedByVersion, config.GetVersion()) ||
		record.FileSize != size ||
		record.ModifiedUnixNano != modTime.UnixNano() {
		logger.Infof(locale.GetString("log_scan_failure_cache_retry"), filePath)
		return false
	}

	logger.Infof(locale.GetString("log_scan_failure_cache_skip"), filePath, record.Error)
	return true
}

// shouldRetryFailedArchiveByVersion 判断旧版本失败缓存是否需要在当前版本重试。
// 同一主/次版本下，小幅 patch 升级继续沿用失败缓存；patch 跨度超过 10 才重新尝试。
func shouldRetryFailedArchiveByVersion(recordVersion, currentVersion string) bool {
	if recordVersion == currentVersion {
		return false
	}
	recordMajor, recordMinor, recordPatch, okRecord := parseSemanticVersion(recordVersion)
	currentMajor, currentMinor, currentPatch, okCurrent := parseSemanticVersion(currentVersion)
	if !okRecord || !okCurrent {
		return true
	}
	if recordMajor != currentMajor || recordMinor != currentMinor {
		return true
	}
	return currentPatch-recordPatch > 10 || currentPatch < recordPatch
}

// parseSemanticVersion 解析项目使用的 vX.Y.Z 版本号，只关心前三段数字。
func parseSemanticVersion(version string) (major int, minor int, patch int, ok bool) {
	version = strings.TrimPrefix(strings.TrimSpace(version), "v")
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return 0, 0, 0, false
	}
	values := make([]int, 3)
	for i, part := range parts {
		value, err := strconv.Atoi(part)
		if err != nil {
			return 0, 0, 0, false
		}
		values[i] = value
	}
	return values[0], values[1], values[2], true
}

func recordArchiveScanFailure(storeURL, filePath string, size int64, modTime time.Time, isRemote bool, scanErr error) {
	if scanErr == nil || !isArchiveScanFailureTarget(filePath) {
		return
	}

	scanFailureCacheMu.Lock()
	defer scanFailureCacheMu.Unlock()

	cache, err := loadScanFailureCache()
	if err != nil {
		logger.Infof(locale.GetString("log_scan_failure_cache_load_failed"), err)
		return
	}

	cache[scanFailureRecordKey(storeURL, filePath, isRemote)] = ScanFailureRecord{
		StoreURL:         storeURL,
		FilePath:         filePath,
		FileSize:         size,
		ModifiedUnixNano: modTime.UnixNano(),
		CreatedByVersion: config.GetVersion(),
		FailedAt:         time.Now(),
		Error:            scanErr.Error(),
		IsRemote:         isRemote,
	}

	if err := saveScanFailureCache(cache); err != nil {
		logger.Infof(locale.GetString("log_scan_failure_cache_save_failed"), err)
		return
	}
	logger.Infof(locale.GetString("log_scan_failure_cache_recorded"), filePath, scanErr)
}

func clearArchiveScanFailure(storeURL, filePath string, isRemote bool) {
	if !isArchiveScanFailureTarget(filePath) {
		return
	}

	scanFailureCacheMu.Lock()
	defer scanFailureCacheMu.Unlock()

	cache, err := loadScanFailureCache()
	if err != nil {
		logger.Infof(locale.GetString("log_scan_failure_cache_load_failed"), err)
		return
	}

	key := scanFailureRecordKey(storeURL, filePath, isRemote)
	if _, ok := cache[key]; !ok {
		return
	}
	delete(cache, key)
	if err := saveScanFailureCache(cache); err != nil {
		logger.Infof(locale.GetString("log_scan_failure_cache_save_failed"), err)
	}
}
