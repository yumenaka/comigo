//go:build !ios

package tools

import (
	"errors"
	"sync"
	"testing"
)

// 验证 CPU 使用率平均值会忽略无效输入并返回是否可用。
func TestAverageCPUPercent(t *testing.T) {
	tests := []struct {
		name   string
		input  []float64
		want   float64
		wantOK bool
	}{
		{
			name:   "average multiple logical cores",
			input:  []float64{10, 20, 30, 40},
			want:   25,
			wantOK: true,
		},
		{
			name:   "single value remains unchanged",
			input:  []float64{37.5},
			want:   37.5,
			wantOK: true,
		},
		{
			name:   "empty result is not usable",
			input:  nil,
			want:   0,
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := averageCPUPercent(tt.input)
			if ok != tt.wantOK {
				t.Fatalf("averageCPUPercent() ok = %v, want %v", ok, tt.wantOK)
			}
			if got != tt.want {
				t.Fatalf("averageCPUPercent() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 验证平台未实现的系统指标错误不会被记录为真实采集错误。
func TestLogSystemMetricErrorOnceIgnoresNotImplemented(t *testing.T) {
	systemMetricLoggedErrors = sync.Map{}

	logSystemMetricErrorOnce("cpu", errors.New("not implemented yet"))
	if hasLoggedSystemMetricErrors() {
		t.Fatal("not implemented metric errors should be ignored")
	}
}

func hasLoggedSystemMetricErrors() bool {
	hasAny := false
	systemMetricLoggedErrors.Range(func(_, _ any) bool {
		hasAny = true
		return false
	})
	return hasAny
}
