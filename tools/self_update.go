//go:build !js

package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/minio/selfupdate"
	"golang.org/x/mod/semver"
)

// latestReleaseResponse 用于解析 GitHub Releases API 的响应
type latestReleaseResponse struct {
	TagName string `json:"tag_name"`
}

// UpdateHandler 提供一个 API，用于：
// 1. 检查当前版本与 GitHub 上最新版本号
// 2. 若有新版本，则执行更新
// 3. 返回更新结果(JSON)
func UpdateHandler(c echo.Context) error {
	owner := "yumenaka"
	repo := "comigo"

	// 获取当前版本
	curVer := "v0.0.0" // TODO: 从config获取当前版本号

	// 1. 拉取最新版本信息
	latestTag, err := fetchLatestTagName(owner, repo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to fetch latest release info",
			"details": err.Error(),
		})
	}

	// 2. 规范化对比(tag_name 类似 "v1.2.3"，可使用 x/mod/semver)
	if !strings.HasPrefix(latestTag, "v") {
		latestTag = "v" + latestTag
	}
	if semver.Compare(curVer, latestTag) >= 0 {
		// 当前版本不低于远程版本 => 无需更新
		return c.JSON(http.StatusOK, map[string]interface{}{
			"current_version": curVer,
			"latest_version":  latestTag,
			"updated":         false,
			"message":         "Already up-to-date",
		})
	}

	// 3. 构造下载链接
	downloadURL := buildDownloadURL(owner, repo, latestTag)

	// 4. 解压并执行更新
	// 如果发布的是单二进制，可以直接 Apply
	// 为了减小文件大小，comigo用的是 tar.gz或zip，需先下载并解压再 Apply
	if err := doSelfUpdate(downloadURL); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Update failed",
			"details": err.Error(),
		})
	}
	// 更新成功
	return c.JSON(http.StatusOK, map[string]interface{}{
		"current_version": curVer,
		"latest_version":  latestTag,
		"updated":         true,
		"message":         "Update successful. Please restart the application if needed.",
	})
}

// fetchLatestTagName 调用 GitHub Releases 最新版接口并返回 tag_name
func fetchLatestTagName(owner, repo string) (string, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var release latestReleaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	if release.TagName == "" {
		return "", fmt.Errorf("tag_name not found in GitHub response")
	}

	return release.TagName, nil
}

// buildDownloadURL 根据 OS/ARCH 和最新版本号 构造下载地址
// 文件名类似：comi_v1.0.0_Linux_x86_64.tar.gz
func buildDownloadURL(owner, repo, latestTag string) string {
	osName := runtime.GOOS                        // "linux", "darwin", ...
	arch := runtime.GOARCH                        // "amd64", "arm64", ...
	version := strings.TrimPrefix(latestTag, "v") // "1.2.3"

	//  GitHub Release 实际命名是 comi_v{Version}_{OS}_{ARCH}.tar.gz
	// 可做类似如下处理：
	// 注意：如果你OS/ARCH 映射和 Go runtime 不同(如 x86_64 => amd64)，
	// 则需要自行转换
	fileName := fmt.Sprintf("comi_v%s_%s_%s", version, osName, arch)

	// 如果不是windows平台，发布文件是 .tar.gz, 这里再补上 .tar.gz
	if osName != "windows" {
		fileName += ".tar.gz"
	} else {
		fileName += ".zip"
	}
	downloadURL := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
		owner, repo, latestTag, fileName)

	return downloadURL
}

// doSelfUpdate 使用 selfupdate.Apply 更新自身
func doSelfUpdate(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download update binary: %w", err)
	}
	defer resp.Body.Close()

	// 写入  tar.gz 或 zip文件到 tmpFile
	if runtime.GOOS == "windows" {
		// 解压 zip文件
	} else {
		// 解压 tar.gz文件
	}

	// 取出可执行文件 comi
	// 调用 selfupdate.Apply(可执行文件的io.Reader, ...)

	// minio/selfupdate 直接对二进制进行替换
	// TODO： checksum 用于校验下载文件的完整性
	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	if err != nil {
		return fmt.Errorf("failed to apply update: %w", err)
	}
	return nil
}

// func main() {
//	e := echo.New()
//
//	// 路径仅供参考，如需安全控制可加权限校验
//	e.GET("/update", UpdateHandler)
//
//	e.Logger.Fatal(e.Start(":8080"))
// }
