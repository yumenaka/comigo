package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yumenaka/comigo/tools"
)

// 验证保存配置会创建目标文件，并记录实际使用的配置路径。
func TestUpdateConfigFileCreatesTargetAndTracksConfigFile(t *testing.T) {
	oldCfg := cfg
	t.Cleanup(func() {
		cfg = oldCfg
	})

	cfg = newDefaultConfig()
	targetDir := t.TempDir()
	targetFile := filepath.Join(targetDir, "nested", "config.toml")
	cfg.ConfigFile = targetFile
	cfg.Debug = true

	if err := UpdateConfigFile(); err != nil {
		t.Fatalf("UpdateConfigFile 返回错误: %v", err)
	}

	if cfg.ConfigFile != targetFile {
		t.Fatalf("ConfigFile 未记录实际写入路径: got %q want %q", cfg.ConfigFile, targetFile)
	}

	content, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("读取生成的配置文件失败: %v", err)
	}
	if len(content) == 0 {
		t.Fatal("生成的配置文件内容为空")
	}
}

// 验证启用插件列表会写入配置文件。
func TestUpdateConfigFilePersistsEnabledPluginList(t *testing.T) {
	oldCfg := cfg
	t.Cleanup(func() {
		cfg = oldCfg
	})

	cfg = newDefaultConfig()
	cfg.ConfigFile = filepath.Join(t.TempDir(), "config.toml")
	cfg.EnablePlugin = true
	cfg.EnabledPluginList = []string{"clock"}

	if err := UpdateConfigFile(); err != nil {
		t.Fatalf("UpdateConfigFile 返回错误: %v", err)
	}

	content, err := os.ReadFile(cfg.ConfigFile)
	if err != nil {
		t.Fatalf("读取生成的配置文件失败: %v", err)
	}
	if !strings.Contains(string(content), "EnabledPluginList") || !strings.Contains(string(content), "clock") {
		t.Fatalf("启用插件列表未写入配置文件:\n%s", string(content))
	}
}

// 验证临时阅读模式不会创建、保存或删除配置文件。
func TestTemporaryReaderModeSkipsConfigFileWrites(t *testing.T) {
	oldCfg := cfg
	t.Cleanup(func() {
		cfg = oldCfg
	})

	cfg = newDefaultConfig()
	cfg.TemporaryReaderMode = true
	dir := t.TempDir()
	t.Chdir(dir)

	explicitConfig := filepath.Join(dir, "explicit.toml")
	cfg.ConfigFile = explicitConfig
	if err := UpdateConfigFile(); err != nil {
		t.Fatalf("UpdateConfigFile 返回错误: %v", err)
	}
	if _, err := os.Stat(explicitConfig); !os.IsNotExist(err) {
		t.Fatalf("临时模式不应写入显式配置文件: %v", err)
	}

	if err := SaveConfig(WorkingDirectory); err != nil {
		t.Fatalf("SaveConfig 返回错误: %v", err)
	}
	workingConfig := filepath.Join(dir, PlatformConfigFilename())
	if _, err := os.Stat(workingConfig); !os.IsNotExist(err) {
		t.Fatalf("临时模式不应保存工作目录配置: %v", err)
	}

	if err := os.WriteFile(workingConfig, []byte("Port = 4321\n"), 0o644); err != nil {
		t.Fatalf("写入测试配置失败: %v", err)
	}
	if err := DeleteConfigIn(WorkingDirectory); err != nil {
		t.Fatalf("DeleteConfigIn 返回错误: %v", err)
	}
	if _, err := os.Stat(workingConfig); err != nil {
		t.Fatalf("临时模式不应删除已有配置文件: %v", err)
	}
}

// 验证默认保存使用当前启动壳的配置文件名。
func TestSaveConfigUsesPlatformFilename(t *testing.T) {
	oldCfg := cfg
	t.Cleanup(func() {
		cfg = oldCfg
	})

	cfg = newDefaultConfig()
	dir := t.TempDir()
	t.Setenv("HOME", t.TempDir())
	t.Chdir(dir)

	if err := SaveConfig(WorkingDirectory); err != nil {
		t.Fatalf("SaveConfig 默认保存返回错误: %v", err)
	}
	platformConfig := filepath.Join(dir, PlatformConfigFilename())
	if _, err := os.Stat(platformConfig); err != nil {
		t.Fatalf("默认保存未写入平台配置文件: %v", err)
	}
	content, err := os.ReadFile(platformConfig)
	if err != nil {
		t.Fatalf("读取平台配置失败: %v", err)
	}
	if strings.Contains(string(content), "UseUnifiedConfig") {
		t.Fatalf("不应再写入全平台共享配置字段:\n%s", string(content))
	}
}

// 验证只读取当前启动壳的配置文件名。
func TestFindConfigFileUsesOnlyPlatformFilename(t *testing.T) {
	oldCfg := cfg
	t.Cleanup(func() {
		cfg = oldCfg
	})

	cfg = newDefaultConfig()
	home := t.TempDir()
	work := t.TempDir()
	t.Setenv("HOME", home)
	t.Chdir(work)

	otherConfig := filepath.Join(work, "cli.toml")
	if otherConfig == filepath.Join(work, PlatformConfigFilename()) {
		otherConfig = filepath.Join(work, "unused.toml")
	}
	if err := os.WriteFile(otherConfig, []byte("Port = 4321\n"), 0o644); err != nil {
		t.Fatalf("写入非当前壳配置失败: %v", err)
	}

	if location, filePath := FindConfigFile(); location != "" || filePath != "" {
		t.Fatalf("非当前壳配置不应生效: got %q %q", location, filePath)
	}

	platformConfig := filepath.Join(work, PlatformConfigFilename())
	if err := os.WriteFile(platformConfig, []byte("Port = 1235\n"), 0o644); err != nil {
		t.Fatalf("写入平台配置失败: %v", err)
	}
	location, filePath := FindConfigFile()
	if location != WorkingDirectory || filePath != platformConfig {
		t.Fatalf("平台配置读取不正确: got %q %q want %q %q", location, filePath, WorkingDirectory, platformConfig)
	}
}

// 验证已有生效配置时，配置管理器只能保存回同一个目录。
func TestSaveConfigOnlyAllowsActiveLocation(t *testing.T) {
	oldCfg := cfg
	t.Cleanup(func() {
		cfg = oldCfg
	})

	cfg = newDefaultConfig()
	home := t.TempDir()
	work := t.TempDir()
	t.Setenv("HOME", home)
	t.Chdir(work)

	if !canSaveConfigTo(HomeDirectory) || !canSaveConfigTo(WorkingDirectory) || !canSaveConfigTo(ProgramDirectory) {
		t.Fatal("没有生效配置时应允许保存到任一配置目录")
	}

	homeConfig := filepath.Join(home, ".config", "comigo", PlatformConfigFilename())
	if err := os.MkdirAll(filepath.Dir(homeConfig), 0o755); err != nil {
		t.Fatalf("创建配置目录失败: %v", err)
	}
	if err := os.WriteFile(homeConfig, []byte("Port = 4321\n"), 0o644); err != nil {
		t.Fatalf("写入 Home 配置失败: %v", err)
	}
	if !canSaveConfigTo(HomeDirectory) {
		t.Fatal("应允许保存回当前生效配置目录")
	}
	if canSaveConfigTo(WorkingDirectory) || canSaveConfigTo(ProgramDirectory) {
		t.Fatal("已有生效配置时不应允许保存到其他目录")
	}
	if err := SaveConfig(WorkingDirectory); err == nil {
		t.Fatal("SaveConfig 应拒绝保存到非当前生效目录")
	}
}

// 验证显式配置文件尚不存在时也会使用其所在目录作为配置目录。
func TestGetConfigDirUsesExplicitConfigPathWhenFileDoesNotExist(t *testing.T) {
	oldCfg := cfg
	t.Cleanup(func() {
		cfg = oldCfg
	})

	cfg = newDefaultConfig()
	targetDir := filepath.Join(t.TempDir(), "nested")
	cfg.ConfigFile = filepath.Join(targetDir, "config.toml")

	got, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir 返回错误: %v", err)
	}
	if got != targetDir {
		t.Fatalf("配置目录不正确: got %q want %q", got, targetDir)
	}
	if _, err := os.Stat(targetDir); err != nil {
		t.Fatalf("显式配置目录未创建: %v", err)
	}
}

// 验证 JSON 配置更新会写入全局配置对象。
func TestUpdateConfigByJsonUpdatesGlobalConfig(t *testing.T) {
	oldCfg := cfg
	t.Cleanup(func() {
		cfg = oldCfg
	})

	cfg = newDefaultConfig()
	cfg.Debug = false
	cfg.MinImageNum = 3

	if err := UpdateConfigByJson(`{"Debug":true,"MinImageNum":1}`); err != nil {
		t.Fatalf("UpdateConfigByJson 返回错误: %v", err)
	}
	if !cfg.Debug {
		t.Fatal("Debug 未按 JSON 更新为 true")
	}
	if cfg.MinImageNum != 1 {
		t.Fatalf("MinImageNum 未按 JSON 更新: got %d want 1", cfg.MinImageNum)
	}
}

// 确认远程书库和扫描共享的超时有明确默认值。
func TestDefaultConfigHasRemoteTimeout(t *testing.T) {
	c := newDefaultConfig()
	if c.TimeoutLimitForScan != 20 {
		t.Fatalf("默认扫描/远程书库超时不正确: got %d want 20", c.TimeoutLimitForScan)
	}
}

// 验证书库路径重叠检测能拦截父子目录重复添加。
func TestIsPathOverlapping(t *testing.T) {
	// 创建一个测试配置
	c := &Config{
		StoreUrls: []string{},
	}

	// 获取临时目录作为测试基础路径
	tmpDir := os.TempDir()

	// 测试场景1: 添加第一个路径
	testPath1 := filepath.Join(tmpDir, "test_store1")
	err := c.AddStoreUrl(testPath1)
	if err != nil {
		t.Errorf("添加第一个路径失败: %v", err)
	}

	// 验证路径已经转换为绝对路径
	if len(c.StoreUrls) != 1 {
		t.Errorf("期望有1个书库，实际有 %d 个", len(c.StoreUrls))
	}
	if !filepath.IsAbs(c.StoreUrls[0]) {
		t.Errorf("路径未转换为绝对路径: %s", c.StoreUrls[0])
	}

	// 测试场景2: 尝试添加相同路径（应该失败）
	err = c.AddStoreUrl(testPath1)
	if err == nil {
		t.Error("添加相同路径应该失败")
	}

	// 测试场景3: 尝试添加子目录（应该失败）
	testPath2 := filepath.Join(testPath1, "subdir")
	err = c.AddStoreUrl(testPath2)
	if err == nil {
		t.Error("添加子目录应该失败")
	}

	// 测试场景4: 尝试添加父目录（应该失败）
	testPath3 := filepath.Dir(testPath1)
	err = c.AddStoreUrl(testPath3)
	if err == nil {
		t.Error("添加父目录应该失败")
	}

	// 测试场景5: 添加不相关的路径（应该成功）
	testPath4 := filepath.Join(tmpDir, "test_store2")
	err = c.AddStoreUrl(testPath4)
	if err != nil {
		t.Errorf("添加不相关的路径失败: %v", err)
	}

	if len(c.StoreUrls) != 2 {
		t.Errorf("期望有2个书库，实际有 %d 个", len(c.StoreUrls))
	}
}

// 验证本地路径从属关系判断。
func TestIsSubPath(t *testing.T) {
	tests := []struct {
		parent   string
		child    string
		expected bool
	}{
		{"/a/b", "/a/b/c", true},
		{"/a/b", "/a/b/c/d", true},
		{"/a/b", "/a/c", false},
		{"/a/b", "/a", false},
		{"/a/b", "/a/b", false},       // 相同路径不算子路径
		{"/x/y", "/x/y/z/../w", true}, // 清理后仍是子路径
	}

	for _, tt := range tests {
		result := tools.IsSubPath(filepath.Clean(tt.parent), filepath.Clean(tt.child))
		if result != tt.expected {
			t.Errorf("isSubPath(%s, %s) = %v, 期望 %v", tt.parent, tt.child, result, tt.expected)
		}
	}
}

// 验证添加书库地址时会做路径重合校验。
func TestAddStringArrayConfigWithStoreUrls(t *testing.T) {
	c := &Config{
		StoreUrls: []string{},
	}

	// 测试添加相对路径
	relPath := "."
	_, err := c.AddStringArrayConfig("StoreUrls", relPath)
	if err != nil {
		t.Errorf("添加相对路径失败: %v", err)
	}

	// 验证已转换为绝对路径
	if len(c.StoreUrls) != 1 {
		t.Errorf("期望有1个书库，实际有 %d 个", len(c.StoreUrls))
	}
	if !filepath.IsAbs(c.StoreUrls[0]) {
		t.Errorf("路径未转换为绝对路径: %s", c.StoreUrls[0])
	}
}

// 验证删除书库地址会更新书库列表。
func TestDeleteStringArrayConfigWithStoreUrls(t *testing.T) {
	c := &Config{
		StoreUrls: []string{},
	}

	tmpDir := os.TempDir()
	testPath := filepath.Join(tmpDir, "test_store_delete")

	// 添加一个路径
	err := c.AddStoreUrl(testPath)
	if err != nil {
		t.Errorf("添加路径失败: %v", err)
	}

	if len(c.StoreUrls) != 1 {
		t.Errorf("期望有1个书库，实际有 %d 个", len(c.StoreUrls))
	}

	// 使用绝对路径删除
	absPath, _ := filepath.Abs(testPath)
	_, err = c.DeleteStringArrayConfig("StoreUrls", absPath)
	if err != nil {
		t.Errorf("删除路径失败: %v", err)
	}

	if len(c.StoreUrls) != 0 {
		t.Errorf("期望有0个书库，实际有 %d 个", len(c.StoreUrls))
	}
}
