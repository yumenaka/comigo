package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIsPathOverlapping 测试路径重合检测
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

// TestIsSubPath 测试子路径检测
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
		result := isSubPath(filepath.Clean(tt.parent), filepath.Clean(tt.child))
		if result != tt.expected {
			t.Errorf("isSubPath(%s, %s) = %v, 期望 %v", tt.parent, tt.child, result, tt.expected)
		}
	}
}

// TestAddStringArrayConfigWithStoreUrls 测试 AddStringArrayConfig 对 StoreUrls 的特殊处理
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

// TestDeleteStringArrayConfigWithStoreUrls 测试删除 StoreUrls
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
