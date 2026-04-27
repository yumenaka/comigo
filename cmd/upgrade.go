package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	toolsfile "github.com/yumenaka/comigo/tools/file"
	"golang.org/x/mod/semver"
)

const (
	upgradeProxyHost   = "https://comigo.xyz"
	latestReleasePath  = "/yumenaka/api.github.com/repos/yumenaka/comigo/releases/latest"
	downloadPathPrefix = "/yumenaka/comigo/releases/download/"
	upgradeUserAgent   = "Comigo-Upgrade/"
)

// ghAsset release 资源条目（仅解析名称）
type ghAsset struct {
	Name string `json:"name"`
}

// ghRelease 解析 GitHub releases/latest JSON 所需字段
type ghRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []ghAsset `json:"assets"`
}

// loadLatestRelease 从 comigo.xyz 反代拉取 GitHub latest release JSON。
func loadLatestRelease() (*ghRelease, error) {
	current := strings.TrimSpace(config.GetVersion())
	clientAPI := newUpgradeHTTPClient(45 * time.Second)

	req, err := http.NewRequest(http.MethodGet, upgradeProxyHost+latestReleasePath, nil)
	if err != nil {
		return nil, fmt.Errorf(locale.GetString("upgrade_fetch_release_failed"), err)
	}
	req.Header.Set("User-Agent", upgradeUserAgent+current)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := clientAPI.Do(req)
	if err != nil {
		return nil, fmt.Errorf(locale.GetString("upgrade_fetch_release_failed"), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(locale.GetString("upgrade_http_status"), resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return nil, fmt.Errorf(locale.GetString("upgrade_fetch_release_failed"), err)
	}

	var rel ghRelease
	if err := json.Unmarshal(body, &rel); err != nil {
		return nil, fmt.Errorf(locale.GetString("upgrade_fetch_release_failed"), err)
	}
	if strings.TrimSpace(rel.TagName) == "" {
		return nil, fmt.Errorf(locale.GetString("upgrade_fetch_release_failed"), fmt.Errorf("empty tag_name"))
	}
	return &rel, nil
}

// runSelfUpgrade 经 comigo.xyz 反代检查版本、下载 CLI 包并替换当前可执行文件。
func runSelfUpgrade() error {
	cfg := config.GetCfg()
	fmt.Println(locale.GetString("upgrade_checking_release"))

	current := strings.TrimSpace(config.GetVersion())
	rel, err := loadLatestRelease()
	if err != nil {
		return err
	}

	remoteTag := strings.TrimSpace(rel.TagName)

	rc, cc := canonicalSemverTag(remoteTag), canonicalSemverTag(current)
	if rc == "" || cc == "" {
		fmt.Printf("%s\n", fmt.Sprintf(locale.GetString("upgrade_invalid_version_compare"), current, remoteTag))
		return nil
	}
	if semver.Compare(rc, cc) <= 0 {
		fmt.Printf("%s\n", fmt.Sprintf(locale.GetString("upgrade_already_latest"), current, remoteTag))
		return nil
	}

	osLabel, archLabel, ok := platformReleaseLabels()
	if !ok {
		return fmt.Errorf(locale.GetString("upgrade_unsupported_arch"), runtime.GOOS, runtime.GOARCH)
	}

	assetName := pickCLIAssetName(rel.Assets, remoteTag, osLabel, archLabel)
	if assetName == "" {
		return fmt.Errorf(locale.GetString("upgrade_no_matching_asset"), osLabel+"_"+archLabel)
	}

	fmt.Printf("%s\n", fmt.Sprintf(locale.GetString("upgrade_new_version"), remoteTag, current))
	dlURL := buildDownloadURL(remoteTag, assetName)
	fmt.Printf("%s\n", fmt.Sprintf(locale.GetString("upgrade_downloading"), dlURL))

	archivePath, err := downloadToTempFile(dlURL, current)
	if err != nil {
		return fmt.Errorf(locale.GetString("upgrade_download_failed"), err)
	}
	defer os.Remove(archivePath)

	extractDir, err := os.MkdirTemp("", "comigo-upgrade-*")
	if err != nil {
		return fmt.Errorf(locale.GetString("upgrade_extract_failed"), err)
	}
	defer os.RemoveAll(extractDir)

	if err := toolsfile.UnArchiveAuto(archivePath, extractDir, cfg.ZipFileTextEncoding); err != nil {
		return fmt.Errorf(locale.GetString("upgrade_extract_failed"), err)
	}

	newBin, err := findComiBinary(extractDir)
	if err != nil {
		return err
	}

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf(locale.GetString("upgrade_replace_failed"), err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf(locale.GetString("upgrade_replace_failed"), err)
	}

	if err := replaceExecutable(execPath, newBin); err != nil {
		return fmt.Errorf(locale.GetString("upgrade_replace_failed"), err)
	}

	fmt.Println(locale.GetString("upgrade_success"))
	return nil
}

func canonicalSemverTag(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	if semver.IsValid(v) {
		return v
	}
	return ""
}

// platformReleaseLabels 与 comigo_xyz 插件中 CLI 包命名一致
func platformReleaseLabels() (osLabel, archLabel string, ok bool) {
	switch runtime.GOOS {
	case "windows":
		osLabel = "Windows"
	case "darwin":
		osLabel = "MacOS"
	case "linux":
		osLabel = "Linux"
	default:
		return "", "", false
	}
	switch runtime.GOARCH {
	case "amd64":
		archLabel = "x86_64"
	case "arm64":
		archLabel = "arm64"
	case "arm":
		archLabel = "armv7"
	default:
		return "", "", false
	}
	return osLabel, archLabel, true
}

func pickCLIAssetName(assets []ghAsset, tag, osLabel, archLabel string) string {
	want := fmt.Sprintf("comi_%s_%s_%s.tar.gz", tag, osLabel, archLabel)
	for _, a := range assets {
		if a.Name == want {
			return want
		}
	}
	suffix := fmt.Sprintf("_%s_%s.tar.gz", osLabel, archLabel)
	for _, a := range assets {
		if strings.HasPrefix(a.Name, "comi_") && strings.HasSuffix(a.Name, suffix) {
			return a.Name
		}
	}
	return ""
}

func buildDownloadURL(tag, fileName string) string {
	return upgradeProxyHost + downloadPathPrefix + url.PathEscape(tag) + "/" + url.PathEscape(fileName)
}

func downloadToTempFile(rawURL, version string) (string, error) {
	client := newUpgradeHTTPClient(45 * time.Minute)
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", upgradeUserAgent+version)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s", resp.Status)
	}

	f, err := os.CreateTemp("", "comigo-dl-*")
	if err != nil {
		return "", err
	}
	tmpPath := f.Name()
	ok := false
	defer func() {
		if !ok {
			_ = f.Close()
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}

	fi, err := os.Stat(tmpPath)
	if err != nil || fi.Size() == 0 {
		_ = os.Remove(tmpPath)
		if err == nil {
			err = fmt.Errorf("empty file")
		}
		return "", err
	}
	ok = true
	return tmpPath, nil
}

func findComiBinary(root string) (string, error) {
	var found string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		if strings.EqualFold(base, "comi") || strings.EqualFold(base, "comi.exe") {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", locale.GetString("upgrade_binary_not_found"), err)
	}
	if found == "" {
		return "", fmt.Errorf("%s", locale.GetString("upgrade_binary_not_found"))
	}
	return found, nil
}

func replaceExecutable(execPath, newBin string) error {
	dir := filepath.Dir(execPath)
	tmpNext := filepath.Join(dir, ".comigo-upgrade-next"+randomSuffix())
	if err := copyFile(newBin, tmpNext); err != nil {
		return err
	}
	if runtime.GOOS != "windows" {
		if err := os.Chmod(tmpNext, 0o755); err != nil {
			os.Remove(tmpNext)
			return err
		}
		if err := os.Rename(tmpNext, execPath); err != nil {
			os.Remove(tmpNext)
			return err
		}
		return nil
	}
	// Windows：先改名正在运行的 exe，再把新文件改为目标名
	oldBackup := execPath + ".upgrade_bak"
	_ = os.Remove(oldBackup)
	if err := os.Rename(execPath, oldBackup); err != nil {
		os.Remove(tmpNext)
		return fmt.Errorf("rename running exe: %w", err)
	}
	if err := os.Rename(tmpNext, execPath); err != nil {
		_ = os.Rename(oldBackup, execPath)
		os.Remove(tmpNext)
		return fmt.Errorf("install new exe: %w", err)
	}
	_ = os.Remove(oldBackup)
	return nil
}

func randomSuffix() string {
	return fmt.Sprintf("-%d", time.Now().UnixNano())
}

// newUpgradeHTTPClient 升级流程专用客户端：关闭 HTTP/2，减轻部分反代与 Content-Length 不一致的兼容问题。
func newUpgradeHTTPClient(timeout time.Duration) *http.Client {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		// 不向服务端声明 gzip，避免部分 CDN/反代返回压缩体却未带 Content-Encoding，导致 net/http 不解压、JSON 解析失败
		DisableCompression:    true,
		ForceAttemptHTTP2:     false,
		TLSHandshakeTimeout:   15 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	return &http.Client{
		Timeout:   timeout,
		Transport: tr,
		CheckRedirect: func(_ *http.Request, via []*http.Request) error {
			if len(via) >= 16 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		os.Remove(dst)
		return err
	}
	return out.Close()
}
