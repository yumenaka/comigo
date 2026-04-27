package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	toolsfile "github.com/yumenaka/comigo/tools/file"
	"github.com/yumenaka/comigo/tools/logger"
	"golang.org/x/mod/semver"
)

// pickTrayAssetName 托盘版发布包：Windows 为 full.zip，macOS 为 Comigo_*.dmg，Linux 为 comigo_*_Linux_*.tar.gz。
func pickTrayAssetName(assets []ghAsset, tag string) string {
	switch runtime.GOOS {
	case "windows":
		want := fmt.Sprintf("comigo_%s_Windows_x86_64_full.zip", tag)
		for _, a := range assets {
			if a.Name == want {
				return want
			}
		}
		for _, a := range assets {
			n := a.Name
			if strings.HasPrefix(n, "comigo_") && strings.Contains(n, "Windows") && strings.HasSuffix(n, "_full.zip") {
				return n
			}
		}
	case "darwin":
		want := fmt.Sprintf("Comigo_%s.dmg", tag)
		for _, a := range assets {
			if a.Name == want {
				return want
			}
		}
		for _, a := range assets {
			n := a.Name
			if strings.HasSuffix(strings.ToLower(n), ".dmg") && strings.HasPrefix(n, "Comigo_") {
				return n
			}
		}
	case "linux":
		_, archLabel, ok := platformReleaseLabels()
		if !ok {
			return ""
		}
		want := fmt.Sprintf("comigo_%s_Linux_%s.tar.gz", tag, archLabel)
		for _, a := range assets {
			if a.Name == want {
				return want
			}
		}
		suffix := fmt.Sprintf("_Linux_%s.tar.gz", archLabel)
		for _, a := range assets {
			n := a.Name
			if strings.HasPrefix(n, "comigo_") && strings.Contains(n, "Linux") && strings.HasSuffix(n, suffix) {
				return n
			}
		}
	}
	return ""
}

// RunTraySelfUpgrade 托盘程序用：经 comigo.xyz 检查版本，若有更新则下载对应平台安装包并替换当前进程镜像。
// 返回 upgraded==true 表示已成功替换，调用方应 PrepareTrayUpgradeRestart 后退出；已是最新则 upgraded==false 且 err==nil。
func RunTraySelfUpgrade() (upgraded bool, err error) {
	logger.Info(locale.GetString("upgrade_checking_release"))
	current := strings.TrimSpace(config.GetVersion())

	rel, err := loadLatestRelease()
	if err != nil {
		return false, err
	}
	remoteTag := strings.TrimSpace(rel.TagName)

	rc, cc := canonicalSemverTag(remoteTag), canonicalSemverTag(current)
	if rc == "" || cc == "" {
		logger.Infof(locale.GetString("upgrade_invalid_version_compare"), current, remoteTag)
		return false, nil
	}
	if semver.Compare(rc, cc) <= 0 {
		logger.Infof(locale.GetString("upgrade_already_latest"), current, remoteTag)
		return false, nil
	}

	assetName := pickTrayAssetName(rel.Assets, remoteTag)
	if assetName == "" {
		return false, fmt.Errorf("%s", locale.GetString("upgrade_tray_no_asset"))
	}

	logger.Infof(locale.GetString("upgrade_new_version"), remoteTag, current)
	dlURL := buildDownloadURL(remoteTag, assetName)
	logger.Infof(locale.GetString("upgrade_downloading"), dlURL)

	archivePath, err := downloadToTempFile(dlURL, current)
	if err != nil {
		return false, fmt.Errorf(locale.GetString("upgrade_download_failed"), err)
	}
	defer os.Remove(archivePath)

	var newBin string
	switch {
	case strings.HasSuffix(strings.ToLower(assetName), ".dmg"):
		newBin, err = extractBinaryFromDmg(archivePath)
		if err != nil {
			return false, fmt.Errorf(locale.GetString("upgrade_tray_dmg_failed"), err)
		}
		defer os.Remove(newBin)
	default:
		extractDir, e := os.MkdirTemp("", "comigo-tray-upgrade-*")
		if e != nil {
			return false, fmt.Errorf(locale.GetString("upgrade_extract_failed"), e)
		}
		defer os.RemoveAll(extractDir)

		cfg := config.GetCfg()
		if e := toolsfile.UnArchiveAuto(archivePath, extractDir, cfg.ZipFileTextEncoding); e != nil {
			return false, fmt.Errorf(locale.GetString("upgrade_extract_failed"), e)
		}
		newBin, err = findTrayBinary(extractDir)
		if err != nil {
			return false, err
		}
	}

	execPath, err := os.Executable()
	if err != nil {
		return false, fmt.Errorf(locale.GetString("upgrade_replace_failed"), err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return false, fmt.Errorf(locale.GetString("upgrade_replace_failed"), err)
	}

	if err := replaceExecutable(execPath, newBin); err != nil {
		return false, fmt.Errorf(locale.GetString("upgrade_replace_failed"), err)
	}

	logger.Info(locale.GetString("upgrade_success"))
	return true, nil
}

// extractBinaryFromDmg 挂载 DMG，从 .app/Contents/MacOS 复制主可执行到临时文件后卸载。
func extractBinaryFromDmg(dmgPath string) (string, error) {
	if runtime.GOOS != "darwin" {
		return "", fmt.Errorf("DMG only on darwin")
	}
	mount, err := os.MkdirTemp("", "comigo-dmg-mount-*")
	if err != nil {
		return "", err
	}
	attached := false
	defer func() {
		if attached {
			_ = exec.Command("hdiutil", "detach", mount, "-force").Run()
		}
		_ = os.RemoveAll(mount)
	}()

	out, err := exec.Command("hdiutil", "attach", "-nobrowse", "-readonly", "-mountpoint", mount, dmgPath).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	attached = true

	src, err := findTrayBinary(mount)
	if err != nil {
		return "", err
	}

	dst, err := os.CreateTemp("", "comigo-tray-bin-*")
	if err != nil {
		return "", err
	}
	dstPath := dst.Name()
	_ = dst.Close()

	if err := copyFile(src, dstPath); err != nil {
		_ = os.Remove(dstPath)
		return "", err
	}
	if err := os.Chmod(dstPath, 0o755); err != nil {
		_ = os.Remove(dstPath)
		return "", err
	}
	return dstPath, nil
}

// findTrayBinary 在解压目录或 DMG 挂载点内查找托盘主程序（与当前进程名或 comigo/Comigo 匹配）。
func findTrayBinary(root string) (string, error) {
	selfPath, err := os.Executable()
	if err != nil {
		selfPath = ""
	}
	want := strings.TrimSuffix(filepath.Base(selfPath), ".exe")
	if want == "" {
		want = "comigo"
	}

	type scored struct {
		path  string
		score int
	}
	var best *scored

	errWalk := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		relLower := strings.ToLower(rel)
		// 跳过 macOS 辅助进程
		if strings.Contains(relLower, "helper") || strings.Contains(relLower, "crash") {
			return nil
		}
		base := filepath.Base(path)
		baseLower := strings.ToLower(base)

		info, err := d.Info()
		if err != nil || info.IsDir() {
			return nil
		}
		// 仅考虑常规文件；忽略明显非二进制
		if !info.Mode().IsRegular() {
			return nil
		}
		if runtime.GOOS == "darwin" && !strings.Contains(rel, ".app/Contents/MacOS/") {
			// DMG 内只认主 MacOS 目录；zip 解压可能在根目录直接有 comigo
			if !strings.HasSuffix(baseLower, ".exe") && baseLower != "comigo" && !strings.EqualFold(base, want) {
				return nil
			}
		}

		score := 0
		if strings.EqualFold(base, want) {
			score += 100
		}
		switch baseLower {
		case "comigo", "comigo.exe":
			score += 80
		case "comi", "comi.exe":
			score += 60
		}
		if strings.EqualFold(base, "Comigo") {
			score += 75
		}
		if strings.Contains(rel, ".app/Contents/MacOS/") {
			score += 10
		}
		if score == 0 {
			return nil
		}
		if best == nil || score > best.score {
			best = &scored{path: path, score: score}
		}
		return nil
	})
	if errWalk != nil {
		return "", fmt.Errorf("%s: %w", locale.GetString("upgrade_binary_not_found"), errWalk)
	}
	if best == nil {
		return "", fmt.Errorf("%s", locale.GetString("upgrade_binary_not_found"))
	}
	return best.path, nil
}
