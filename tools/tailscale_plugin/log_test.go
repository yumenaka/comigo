package tailscale_plugin

import "testing"

// 验证 Tailscale 用户日志中的登录提示会被降噪过滤。
func TestShouldSuppressTailscaleUserLog(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    bool
	}{
		{
			name:    "auth loop running state",
			message: "AuthLoop: state is Running; done",
			want:    true,
		},
		{
			name:    "auth loop other state",
			message: "AuthLoop: state is Stopped; done",
			want:    true,
		},
		{
			name:    "auth url remains visible",
			message: "To start this tsnet server, restart with TS_AUTHKEY set, or go to: https://login.tailscale.com/a/example",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldSuppressTailscaleUserLog(tt.message); got != tt.want {
				t.Fatalf("shouldSuppressTailscaleUserLog(%q) = %v, want %v", tt.message, got, tt.want)
			}
		})
	}
}
