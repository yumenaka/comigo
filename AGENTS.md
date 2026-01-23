# Project Instructions

ComiGo 是一个漫画/图片阅读器，提供 Web 界面，支持多种压缩包格式和阅读模式。

## 关键目录
- `routers/` - 路由定义（Echo v4），`urls.go` 定义所有路由
- `model/` - 业务数据模型（Book, BookInfo, PageInfo, BookMark）。数据默认存储在内存中，并使用本地 json 文件做持久化（当前开发的默认数据存储方式）
- `sqlc/` - 数据库层 SQLite（sqlc 生成类型安全查询），未来扩展用，通过 model 包里的StoreInterface 与内存存储模式兼容，暂时不作为开发重点。`schema.sql` 定义表结构，`query.sql` 定义查询
- `config/` - 全局配置（Config 结构体在 `config.go`）
- `cmd/` - 命令行入口和启动逻辑
- `tools/` - 工具函数（扫描、图片处理、网络等）
- `store/` - 书库管理
- `assets/stores/` - 前端 Alpine.js store

## Code Style
- 代码应该有必要的中文注释，尤其是函数和复杂逻辑部分
- 保持代码简洁易读

## 架构
- 后端：Go + Echo v4（路由分公开组和私有组，私有组需 JWT 认证）
- 前端：bun + JavaScript + Alpine.js（+persist 插件）+ Flowbite + TailwindCSS
- 数据存储：默认内存+json持久化，`sqlc generate` 用于生成 SQLite 查询（未来扩展）
- 模板：templ，`*_templ.go` 是生成文件，修改 `*.templ` 后执行 `templ fmt ./templ && templ generate`
- 国际化：`assets/locale/` 下的 json 文件（en_US.json, ja_JP.json, zh_CN.json）
  - 后端：`locale.GetString("key")`
  - 前端：`i18next.t('key')`
- 前端状态：Alpine.js store 持久化键名格式 `模块.配置项`（如 `flip.autoHideToolbar`）
- 前端构建：`bun run dev`
- 运行指令：`templ fmt ./templ && templ generate && go run main.go`
- 开发时假设 CLI 工具与依赖已安装

## 插件系统
- 内置插件：`templ/plugins/`
- 用户插件：`configDir/plugins/`
- 作用域：`global`（全局）、`shelf`/`flip`/`scroll`（特定页面）、`flip/{bookID}`（特定书籍）