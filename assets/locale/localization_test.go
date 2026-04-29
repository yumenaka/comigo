package locale

import (
	"strings"
	"testing"
)

func TestWebPSettingSaveCompletedIsSuccessMessage(t *testing.T) {
	message := GetStringByLocal("webp_setting_save_completed", "zh-CN")
	if strings.Contains(message, "错误") || !strings.Contains(message, "成功") {
		t.Fatalf("webp_setting_save_completed zh-CN = %q, want a success message", message)
	}
}
