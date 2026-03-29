# Project Instructions

ComiGo 是一个漫画/图片阅读器，提供 Web 界面，支持多种压缩包格式和阅读模式。

## 关键目录
- `routers/` - 路由定义（Echo v4），`urls.go` 定义所有路由
- `mobile/` - 面向移动端宿主的导出包，供 `gomobile bind` 生成 Android AAR
- `model/` - 业务数据模型（Book, BookInfo, PageInfo, BookMark）。数据默认存储在内存中，并使用本地 json 文件做持久化（当前开发的默认数据存储方式）
- `sqlc/` - 数据库层 SQLite（sqlc 生成类型安全查询），未来扩展用，通过 model 包里的StoreInterface 与内存存储模式兼容，暂时不作为开发重点。`schema.sql` 定义表结构，`query.sql` 定义查询
- `config/` - 全局配置（Config 结构体在 `config.go`）
- `cmd/` - 命令行入口和启动逻辑
- `tools/` - 工具函数（扫描、图片处理、网络等）
- `store/` - 书库管理
- `assets/stores/` - 前端 Alpine.js store

## 前端目录
- `assets/script/` 放前端 JavaScript 代码 与  CSS 样式，通过 templ/common/html.templ 这个模板文件`common.Html()`的 insertScripts 这个参数来插入网页中。
- `assets/script/` 中，`assets/script/styles.css` 与 `assets/script/main.js` 是 `bun run dev` 编译生成的文件。其他文件静态导入，开发时可以直接修改源文件。
- `assets/locale/` 放国际化的 json 文件，前端通过 i18next 来加载和使用这些翻译字符串。国际化文件修改后也需要重新编译前端。

## Code Style
- 应该有必要的中文注释，尤其是函数和复杂逻辑部分
- 保持代码简洁易读

## 架构
- 后端：Go + Echo v4（路由分公开组和私有组，私有组需 JWT 认证）
- 嵌入式宿主：Android 侧通过 `gomobile bind` 将 `mobile/` 打包为 AAR，由宿主 App 启动本地 HTTP 服务
- 前端：bun + JavaScript + Alpine.js（+persist 插件）+ Flowbite + TailwindCSS
- 数据存储：默认内存+json持久化，`sqlc generate` 用于生成 SQLite 查询（未来扩展）
- 模板：templ，`*_templ.go` 是生成文件，修改 `*.templ` 后执行 `templ fmt ./templ && templ generate`
- 国际化：`assets/locale/` 下的 json 文件，log一般不会硬编码文字，而是修改（en_US.json, ja_JP.json, zh_CN.json）这三个文件来添加或修改文本内容。修改前端或后端时，尽量同步做好展示文字与log的国际化。
  - 后端使用翻译字符串：`locale.GetString("key")`
  - 前端使用翻译字符串：`i18next.t('key')`
- 前端状态：Alpine.js store 持久化键名格式 `模块.配置项`（如 `flip.autoHideToolbar`）
- 前端构建：`bun run dev`
- 运行指令：`templ fmt ./templ && templ generate && go run main.go`
- 嵌入式 Android 构建：由主仓库执行 `scripts/build_android_go.sh`，实际对 `cvgo/mobile` 运行 `gomobile bind`
- 开发时假设 CLI 工具与依赖已安装

## 移动端嵌入式约束
- `mobile/mobile.go` 提供宿主调用的 `Start`、`Stop`、`GetServerInfo` 等导出接口，接口签名必须保持 `gomobile bind` 兼容，优先使用基础类型。
- 路由层提供 `/healthz` 健康检查入口，供宿主等待本地 HTTP 服务就绪后再加载 WebView。
- `routers.StartWebServer()` / `StartEcho()` 在嵌入式模式下应返回错误，不能依赖 `os.Exit`、`logger.Fatalf` 这类直接结束整个宿主进程的路径。
- 嵌入式模式下的配置目录、缓存目录、书库路径由宿主显式传入，不依赖 CLI 的当前工作目录或 `os.Executable()` 推导。
- 涉及移动端接入时，优先保持 Tailscale、桌面托盘、Windows 特有逻辑与嵌入式启动链路解耦。

## 插件系统
- 内置插件：`templ/plugins/`
- 用户插件：`configDir/plugins/`
- 作用域：`global`（全局）、`shelf`/`flip`/`scroll`（特定页面）、`flip/{bookID}`（特定书籍）