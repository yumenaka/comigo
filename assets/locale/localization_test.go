package locale

import (
	"strings"
	"testing"
)

// 验证 WebP 设置保存完成的中文文案表达为成功而不是错误。
func TestWebPSettingSaveCompletedIsSuccessMessage(t *testing.T) {
	message := GetStringByLocal("webp_setting_save_completed", "zh-CN")
	if strings.Contains(message, "错误") || !strings.Contains(message, "成功") {
		t.Fatalf("webp_setting_save_completed zh-CN = %q, want a success message", message)
	}
}
