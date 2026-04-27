//go:build ios

package tools

func populateSystemMetrics(sys *SystemStatus) {
	// iOS 下先使用基础 runtime 信息，避免引入依赖 IOKit 的桌面系统指标实现。
}
