# Project Instructions

ComiGo 是漫画/图片阅读器，提供 Web 界面，支持压缩包、图片目录、PDF、音频和多种阅读模式。主要栈：Go + Echo v4、templ、bun、Alpine.js、TailwindCSS；默认数据存储为内存 + JSON 持久化，SQLite/sqlc 为后续扩展。

## 关键目录
- `routers/`：Echo 路由，`urls.go` 定义路由；公开/私有组分离，私有组走 JWT。
- `cmd/`：CLI 与启动逻辑；`cmd/mobile/` 供 `gomobile bind` 导出 Android/iOS/macOS 原生库。
- `model/`、`store/`：Book/BookInfo/PageInfo/BookMark 与书库管理，默认内存 + JSON 持久化。
- `config/`：全局配置；嵌入式模式下配置、缓存、书库路径由宿主显式传入。
- `tools/`：扫描、图片处理、网络、VFS、系统工具。
- `sqlc/`：SQLite schema/query 与生成代码，暂非开发重点。
- `assets/frontend/`：需编译的前端入口、样式、插件、stores、utils。
- `assets/static/`：不参与主包编译的页面级 JS/CSS/WASM；页面脚本通过 `common.Html(..., insertScripts)` 引入。
- `assets/dist/`：`bun run dev/build` 生成的主包产物，由 `common.Html()` 自动插入。
- `assets/locale/`：`en_US.json`、`ja_JP.json`、`zh_CN.json` 国际化文本。

## 生成产物与构建边界
- `templ/**/*_templ.go` 是 `*.templ` 的生成产物；`assets/dist/*` 是 `assets/frontend/*` 的编译产物；不要手动编辑生成产物实现业务逻辑。
- 修改 `*.templ` 后执行 `templ fmt ./templ && templ generate`。
- 修改 `assets/frontend/*` 或 `assets/locale/*` 后执行 `bun run dev` 或等价构建命令。
- 提交/review 前分开检查源码和生成产物；生成产物必须能被本轮源码变化解释，不能混入无关 churn。
- 汇报时说明哪些是源码修改，哪些只是同步生成产物。

## 代码与提交规则
- 当前开发环境按 macOS 处理，临时文件/缓存放 `$TMPDIR` 下；保持代码简洁，不做过度复杂的拆分抽象；除非明确要求，不需要额外向前兼容。
- 关键代码与函数写中文注释；前端、后端都适用。
- 一次修改中新增行数超过删除行数 1.5 倍时，额外 review 是否引入不必要抽象、重复实现或啰嗦写法。
- 分析、审计、调研、临时计划类 Markdown 默认只保留在本地 `docs/`，不提交；长期保存需用户确认。
- 除非明确要求，代码和文档不要包含特定本地路径、文件或书籍 ID。
- 国际化展示文本和日志尽量使用 locale：Go/templ 用 `locale.GetString("key")`，前端用 `i18next.t("key")`。

## 架构约束
- 后端：Go + Echo v4；路由层提供 `/healthz` 供宿主等待服务就绪。
- 前端状态：Alpine.js store 持久化键名格式 `模块.配置项`，如 `flip.autoHideToolbar`。
- `cmd/mobile/mobile.go` 导出 `Start`、`Stop`、`GetServerInfo` 等接口，签名保持 `gomobile bind` 兼容，优先使用基础类型。
- `routers.StartWebServer()` / `StartEcho()` 在嵌入式模式下应返回错误，不能依赖 `os.Exit`、`logger.Fatalf` 结束宿主进程。
- 嵌入式接入优先保持 Tailscale、桌面托盘、Windows 特有逻辑与启动链路解耦。

## 常用命令
- 模板生成：`templ fmt ./templ && templ generate`
- 前端构建：`bun run dev`
- 本地运行：`templ fmt ./templ && templ generate && go run main.go`
- 热重载：`air -c .air.toml`；频繁改动导致失败时可杀掉 air 后重试，或改用 `go run main.go`。
- Go cache 权限问题可用：`GOCACHE="${TMPDIR%/}/cvgo-go-build-cache" go test ./...`
- 移动端构建由主仓库脚本触发：`scripts/build_android_go.sh`、`scripts/build_ios_go.sh`、`scripts/build_macos_go.sh`，实际对 `cvgo/cmd/mobile` 执行 `gomobile bind`。

## 插件与 TUI
- 内置插件：`templ/plugins/`；用户插件：`configDir/plugins/`。
- 插件作用域：`global`、`shelf`、`flip`、`scroll`、`flip/{bookID}`。
- TUI 终端图片细节以 `cmd/tui/AGENTS.md` 为准；终端 workaround 不要扩散到全局。
