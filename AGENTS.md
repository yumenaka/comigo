# Project Instructions

ComiGo 是漫画/图片阅读器，提供 Web 界面，支持压缩包、图片目录、PDF、音频和多种阅读模式。主要栈：Go + Echo v4、templ、bun、Alpine.js、TailwindCSS；默认数据存储为内存 + JSON 持久化；启用数据库后，书籍、页面、书签和扫描失败记录直接持久化到 SQLite/PostgreSQL，不读写 metadata JSON。

## 文档分工
- `README*.md` 面向用户，只介绍项目、安装和基础功能；除基础功能外，一般不增加详细说明。
- 开发细节、实现约束和维护规则写入项目根目录或对应模块的 `AGENTS.md`。开发备忘见 [TODO.md](TODO.md)。
- 用户相关的详细用法、配置、限制和故障处理写入自带手册 `assets/static/manual-content/*.json`；README 只保留必要的手册链接，三语同步。

## 目录与边界
- `routers/` 定义 Echo 路由，`urls.go` 管理公开/私有组，私有组走 JWT；`/healthz` 供宿主等待服务就绪。
- `cmd/` 是 CLI 与启动逻辑；`cmd/mobile/` 导出 `Start`、`Stop`、`GetServerInfo` 等 `gomobile bind` 接口，签名优先保持基础类型。
- `model/`、`store/` 管理 Book、BookInfo、PageInfo、BookMark 和书库；`config/` 管理全局配置，嵌入式模式下配置、缓存、书库路径由宿主显式传入。
- `tools/` 放扫描、图片处理、网络、VFS、系统工具；`sqlc/` 放数据库存储实现、SQLite/PostgreSQL schema/query 与生成代码；修改 schema/query 后执行 `sqlc generate`。
- `assets/frontend/` 是需编译的前端入口、样式、插件、stores、utils；`assets/static/` 是不参与主包编译的页面级 JS/CSS/WASM，通过 `common.Html(..., insertScripts)` 引入。
- `assets/dist/` 是 `bun run dev/build` 生成产物，由 `common.Html()` 自动插入；`assets/locale/` 放 `en_US.json`、`ja_JP.json`、`zh_CN.json`。

## 构建与生成产物
- `templ/**/*_templ.go` 和 `assets/dist/*` 是生成产物；不要手动编辑生成产物实现业务逻辑。
- Wails 绑定生成产物只允许由 Wails 生成器同步；当前前端走 HTTP 子路由，不依赖旧 `assets/wailsjs`。
- 修改 `*.templ` 后执行 `templ fmt ./templ && templ generate`；修改 `assets/frontend/*` 或 `assets/locale/*` 后执行 `bun run dev` 或等价构建命令。
- `templ/pages/manual/` 是内置手册页面，正文翻译独立放在 `assets/static/manual-content/*.json`；不需要单独构建，示例图片只使用 HTTPS 外链。
- 提交/review 前分开检查源码和生成产物；生成产物必须能被本轮源码变化解释，不能混入无关 churn。
- 汇报时说明哪些是源码修改，哪些只是同步生成产物。

## 代码与提交规则
- 当前开发环境按 macOS 处理，临时文件/缓存放 `$TMPDIR` 下；Go cache 权限问题可用 `GOCACHE="${TMPDIR%/}/cvgo-go-build-cache"`。
- 保持实现简洁，不做过度拆分抽象；除非明确要求，不需要额外向前兼容；一次修改中新增行数超过删除行数 1.5 倍时，额外 review 是否有不必要抽象、重复实现或啰嗦写法。
- 关键代码与函数写中文注释；前端、后端、测试代码都适用。
- 与长期功能维护无关的临时非 Go 测试、浏览器探针和安装验收脚本只放系统临时目录，不加入代码库；不恢复用户已清理的此类脚本。长期功能回归优先使用 `go test`，正常构建与安装功能脚本不受此限制。
- 临时修改方案及相关调研统一放在所属仓库根目录，命名为 `TEMP_中文描述.md`，并确保该仓库的 `.gitignore` 包含 `TEMP_*.md`；仅本地保存，不提交，不使用 `git add -f` 纳入版本控制。
- `AGENTS.md` 可以在独立的临时段落中引用当前计划；每次提交前删除该引用段落，检查暂存区不包含计划文件及其临时引用，再提交其余已授权的修改。计划文件保留在本地，无需删除；本规则本身保留。
- 其他分析、审计、调研类 Markdown 默认只保留在本地 `docs/`，不提交；长期保存需用户确认。
- 除非明确要求，代码和文档不要包含特定本地路径、文件或书籍 ID。
- 国际化展示文本和日志尽量使用 locale：Go/templ 用 `locale.GetString("key")`，前端用 `i18next.t("key")`。

## 架构约束
- `config.GetCfg()` 是只读快照；仅启动参数绑定可直接修改。运行时更新经 `config.ChangeConfig` 在独立副本上完成；需持久化时先写盘成功再原子发布，失败不应用副作用。重载／回滚用 `config.ReplaceConfig`；不要缓存字段指针或修改已发布快照的切片。
- `routers.StartWebServer()` / `StartEcho()` 在嵌入式模式下应返回错误，不能依赖 `os.Exit`、`logger.Fatalf` 结束宿主进程。
- 嵌入式接入优先保持 Tailscale、桌面托盘、Windows 特有逻辑与启动链路解耦。
- 远程 Comigo 书库使用 `http(s)` 主页 URL 配置；阅读与下载实时代理远端资源，本地只保存索引，不保留原始文件。
- 前端 Alpine.js store 持久化键名格式为 `模块.配置项`，如 `flip.autoHideToolbar`。

## EPUB 阅读
- `templ/pages/epub/`、`assets/static/js/epub_reader.js` 与对应 CSS 是独立阅读器，不复用 ZIP 漫画的 WASM、页码或翻页状态。
- `/epub` 只在浏览器读取本地文件；`/epub/:id` 复用认证后的原文件接口，支持远程书库。纯文字 EPUB 不受最少图片数过滤。
- `assets/static/lib/epub/` 固定分发 epub.js 0.3.93 和 JSZip 3.10.2 的官方构建与许可证，不能改为运行时 CDN。升级时同步文件与版本记录。
- 书籍章节在禁用脚本的 iframe 中展示；注入 CSP 必须发生在渲染前，只允许包内 blob/data 资源，不能允许书籍联网或执行脚本。
- EPUB 的 CFI 阅读位置、字号和配色仅保存在当前浏览器的 `epub.*` 键，不写入漫画页码书签。

## Wails v3
- Wails3 代码、编译与开发配置统一在 `dev` 维护，不再使用独立 `wails3` 分支。
- Wails 修改不要影响非 Wails 环境；优先用 build tag、Wails-only API 或前端运行时判断隔离。
- Wails 开发启动用 `bun run wails:dev`；Makefile 会按 `go.mod` 中的 Wails3 版本安装并优先使用 `GOBIN/wails3`。
- Wails 本机构建用 `bun run wails:build`；Wails3 dev 配置在 `wails3.yml`，构建入口在 `Taskfile.yml`。
- Wails 入口使用 `wails && !js && !bindings` build tag；绑定生成如需恢复，继续避免临时二进制启动 Web 服务、扫描书库或启动 Tailscale。
- Wails 构建使用 `tools/tailscale_plugin/ts_fake.go` 空实现，当前不支持内置 Tailscale；设置页不要在 Wails 环境暴露远程访问/Tailscale 配置。
- Wails3 页面资源由内嵌 Echo 服务提供；仅原生 WebView 插入 `/wails/runtime.js`，普通浏览器不加载 Wails 运行时。
- 独立 tray 使用 `tools/system_tray/` 的 `fyne.io/systray`，以 `!wails` 标签隔离，不得引入 Wails/GTK/WebKit 依赖；desktop 使用 `tools/wails_systray/` 的 Wails3 托盘与面板，仅在 `wails` 标签下编译。面板仅由 Wails 资产处理器提供，通过 Go 绑定调用本进程，不在公共 Echo 上注册控制接口。
- Linux 独立 tray 以 `CGO_ENABLED=0` 交叉编译；仅 desktop 依赖 CGO、GTK4、WebKitGTK 6.0，发布镜像见 `docs/docker/Dockerfile.wails-linux`。托盘需要 StatusNotifier/AppIndicator 宿主；无宿主时 desktop 正常关闭退出，tray 明确报错退出。
- macOS 最低版本为 13；两个入口复用 `docs/macos-app.mk` 打包，自签后生成 DMG。构建图标源为已跟踪的 `assets/appicon.png`。
- Wails3 开发版当前特性：桌面端保留托盘、外部链接、全屏和删除源文件；Android arm64 APK 可构建运行，主页支持整页下拉刷新、书卡长按删除、文件图标与桌面一致。
- Android Wails 暂不支持目录选择器；上传与设置页写操作当前等待 Wails3 正式能力补齐，不再用 fetch 拦截、HTML 重放或原生写文件桥做临时兼容。
- Android APK 构建用 `make wails-android-apk`；Go shared library 必须带 `production,wails,android` tags，少 `wails` 会编进 CLI 入口并导致首屏 `wails.localhost` 连接失败。

## 常用命令
- 模板生成：`templ fmt ./templ && templ generate`
- 前端构建：`bun run dev`
- 本地运行：`templ fmt ./templ && templ generate && go run main.go`
- Wails 开发：`bun run wails:dev`
- 热重载：`air -c .air.toml`；频繁改动导致失败时可杀掉 air 后重试，或改用 `go run main.go`。
- 测试：`GOCACHE="${TMPDIR%/}/cvgo-go-build-cache" go test ./...`；desktop 另执行 `go test -tags wails ./...`。
- 移动端构建由主仓库脚本触发：`scripts/build_android_go.sh`、`scripts/build_ios_go.sh`、`scripts/build_macos_go.sh`，实际对 `cvgo/cmd/mobile` 执行 `gomobile bind`。

## 插件与 TUI
- 内置插件在 `templ/plugins/`，用户插件在 `configDir/plugins/`；插件作用域包括 `global`、`shelf`、`flip`、`scroll`、`flip/{bookID}`。
- TUI 终端图片细节以 `cmd/tui/AGENTS.md` 为准；终端 workaround 不要扩散到全局。

## Omarchy 插件
- 修改 `comigo-omarchy/` 插件前，先阅读 [Comigo Omarchy 插件开发规则](COMIGO_OMARCHY.md)，实现约束与验证命令以该文件为准。

## CLI 安装脚本
- `get.sh` 保持兼容 macOS Bash 3.2 与 Linux Bash；未显式指定下载源、安装目录时分别通过 `/dev/tty` 交互选择，不设默认选择或按 PATH 推断。GitHub 不附地域推荐，comigo.xyz 标注“中国大陆推荐”；候选目录只保留 PATH 中的 `/usr/bin`、`/usr/local/bin`、`$HOME/.local/bin`，按此顺序展示并标注 root 权限要求；macOS 排除受保护的 `/usr/bin`。无终端时必须提供对应参数。
- 修改安装流程后执行 `bash -n get.sh`；需要行为验证时在系统临时目录使用下载替身与临时安装目录，不安装到系统目录，不把临时验证脚本加入代码库。
