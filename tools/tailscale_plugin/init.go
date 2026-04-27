package tailscale_plugin

// Wails 分支禁用真实 Tailscale/tsnet 初始化；由 ts_fake.go 提供空实现，
// 避免引入与 Wails 打包不兼容的 Tailscale 运行时依赖。
