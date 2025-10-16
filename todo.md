# Comigo 开发文档

## 目录
- [技术栈](#技术栈)
- [项目特性](#项目特性)
- [开发环境搭建](#开发环境搭建)
- [项目结构](#项目结构)
- [开发指南](#开发指南)

## 技术栈

### 前端技术
- **模板引擎**: [Templ](https://templ.guide/) - 高性能 Go 模板引擎
- **UI 组件**: [Flowbite](https://flowbite.com/docs/components/gallery/) - 现代化 UI 组件库
- **交互增强**: 
  - [htmx](https://htmx.org/examples) - 现代化 AJAX 库
  - [Alpine.js](https://alpinejs.dev) - 轻量级 JavaScript 框架
- **样式框架**: [Tailwind CSS](https://tailwindcss.com/docs/flex) - 实用优先的 CSS 框架
- **国际化**: [i18next](https://www.i18next.com) - 强大的国际化解决方案
- **图标库**: [Xicons](https://www.xicons.org/#/) - 丰富的图标集合

### 后端技术
- **Web 框架**: [Gin](https://gin-gonic.com/zh-cn/docs/examples/) - 高性能 Go Web 框架
- **HTMX 集成**: [htmx-go](https://github.com/angelofallars/htmx-go) - Go 语言 HTMX 集成

### 开发工具
- **热重载**: [Air](https://github.com/air-verse/air) - Go 应用热重载工具
- **运行时**: [Bun](https://github.com/oven-sh/bun) - 高性能 JavaScript 运行时
- **代码质量**: [golangci-lint](https://github.com/golangci/golangci-lint) - Go 代码质量检查工具
- **代码格式化**: [Prettier](https://prettier.io/docs/en/index.html) - 代码格式化工具

## 项目特性

### 已实现
- [x] 章节快速导航
- [x] 日志记录
- [x] 卷轴模式分页
- [x] 配置文件（TOML 格式）
- [x] 预定义主题与颜色
- [x] 浏览器快捷键
- [x] 多文件支持
- [x] 网页书架
- [x] 新一代图片格式支持（heic avif）。
- [x] 图片自动裁边，分割、拼接单双页。
- [x] 网页端：分享功能
- [x] 网页端：显示QRCode
- [x] 网页端：多种展示模式
- [x] 网页端：服务器设置
- [x] 网页端：HTTPS加密
- [x] 网页端：显示服务器信息
- [x] 网页端：上一章、下一章,快速跳转。
- [x] websocket通信（[参考](https://github.com/Unrud/remote-touchpad)）
- [x] 访问权限设置，账号系统
- [x] log记录
- [x] 设置中心，设置热重载
- [x] CPU、内存占用、状态监控
- [x] 网页端：二维码界面优化。
- [x] gin -> echo。
- [x] 尝试wails, https://v3alpha.wails.io/getting-started/installation/
- [x] 优化配置文件 （[参考](<https://toml.io/cn/v1.0.0）> (better config file formart).
- [x] 服务器设置页面 v1.0
- [x] tailscale 集成
- [x] 解决tailscale带来的cgo编译问题：https://github.com/elastic/golang-crossbuild 或 https://github.com/goreleaser/goreleaser-cross

### 开发中
- [ ]网页端日志查看 50%
- [ ] “阅读历史记录（是否跳转到上一次阅读的最远页）”与“按照最近阅读时间排序”
- [ ] 文件持久化，meta文件，阅读历史与统计。
- [ ] 新官网,Discord频道,使用文档，参考 https://omarchy.org/ 。内置帮助文档?
- [ ] 官网自动探测浏览器平台，提供合适的平台版本（参考Audacity） 
- [ ] 支持smb、webdav文件系统
- [ ] 注册为文件默认打开类型，简单托盘图标，gui界面。
- [ ] 手动或自动检测新版本提示，然后可以试着自动更新新版本（win与macos），最后是各种linux软件源
- [ ] 网页端：浏览器快捷键(50%)。
- [ ] 下载为单个html文件（50%）
- [ ] 示例漫画
- [ ] 翻页模式：滚轮滑动翻页
- [ ] 滑动模式：可快捷键调速的自动翻页
- [ ] cli 交互，tui支持
- [ ] 访问权限控制
- [ ] PWA 支持
- [ ] 系统监控（CPU、内存）
- [ ] 嵌入 HTML
- [ ] 网页端：
- [ ] wasm模式
- [ ] 拆分内存存储与数据库存储，并为支持samba与S3等文件系统做准备。
- [ ] 自动更新
- [ ]文件监控
- [ ]用户系统增强
- [ ]Shell 交互
- [ ]文件管理
- [ ]EPUB/PDF 阅读优化、
- [ ] 高级阅读体验
    - 防剧透效果、回忆模式、特殊背景、背景音乐etc

  
## 开发环境搭建

### 前置要求
- Go 1.24 或更高版本
- Node.js 16 或更高版本
- Bun 运行时

### 安装步骤

1. 安装必要的工具：
```bash
# 安装 Templ
go install github.com/a-h/templ/cmd/templ@latest

# 安装 Air（热重载工具）
go install github.com/air-verse/air@latest

# 安装 Bun
curl -fsSL https://bun.sh/install | bash
```

2. 启动开发服务器：
```bash
#使用 Air
air
```

## 开发指南

### 后端开发
- 后端代码位于项目根目录的 `*.go` 文件中
- 使用 Echo 框架处理 HTTP 请求
- 默认服务器端口：1234

### 前端开发
- 模板文件位于 `./templ` 目录
- 样式和脚本文件位于 `./assets` 目录
- 静态资源位于 `./static` 目录

### 开发提示
- 使用 Air 实现热重载
- 使用 Templ 生成 HTML

### 长期计划
- [ ] 嵌入html，防剧透效果。回忆模式，特殊背景，音乐etc
- [ ] 网页端：优化图片预加载，长图片支持。
- [ ] 网页端：网页前端查看log
- [ ] 跨平台 GUI（Flutter + GoMobile或 Wails）PWA模式。
- [ ] 更新提示，自动更新。包管理支持。
- [ ] 文件夹监控(fsnotify)，自动更新(github.com/jpillora/overseer)
- [ ] 用户系统、访问密码，流量限制等
- [ ] shell 互动（<https://github.com/rivo/tview> ）
- [ ] 子命令，download rar2zip
- [ ] 支持rar压缩包密码。处理损坏文件，扩展名错误的文件，固实压缩文件（7z）。更准确的文件类型判断。
- [ ] 崩溃后恢复，恶意存档处理。
- [ ] 使用新版压缩包处理库（https://github.com/mholt/archives）
- [ ] 编写测试
- [ ] 命令行交互
- [ ] 调用第三方API
- [ ] 文件管理，删除。
- [ ] Debian，RPM包（<https://github.com/goreleaser/nfpm）>
- [ ] 优化epub与PDF阅读体验，支持图文混排（pdf.js与epub.js）
- [ ]  [self update 功能](https://github.com/minio/selfupdate)
- [ ] 显示服务器log：[web终端示例](https://zenn.dev/ikedam/articles/2e078bfc2a4cb6)
- [ ] -start 参数，后台运行。-stop参数，停止后台运行的进程。
- [ ] 自定义js与css代码块功能。
- [ ] [使用 Go1.24 的 os.Root 类型](https://antonz.org/go-1-24/)，将[文件操作限制在特定目录](https://go.dev/blog/osroot)，以防止攻击者通过转义或相对路径非法访问文件
- [ ] 自动发版功能 [goreleaser](https://goreleaser.com/)  [github-action](https://dev.to/hadlow/how-to-release-to-homebrew-with-goreleaser-github-actions-and-semantic-release-2gbb)
- [ ] comigo后台：有几台设备在线，阅读文件，阅读页数，当前用户状态、阅读书籍、阅读进度、阅读时间、服务器状态 注册，阅读记录，
- [ ] 合并htmx代码，参考[pagoda](https://github.com/mikestefanello/pagoda)
  ，重新规划项目结构。我用的许多组件，最终都换成和这个模板一样的了，估计从这个项目里可以学到很多东西。[go-blueprint](https://docs.go-blueprint.dev/)
  也是一个不错的参考，可以看看怎么集成websockets与templ。
- [ ] 添加[数据验证](https://dev.to/leapcell/validator-complex-structs-arrays-and-maps-validation-for-go-34ni)。
- [ ] 优化打开浏览器与扫描逻辑，减少等待时间。可以使用[端口检测包](https://github.com/wait4x/wait4x)。
- [ ] SteamDeck支持（网页支持手柄操作）鼠标滚轮对应
- [ ] 同步翻页 -> 全局多端同步跟踪页面状态，除了不同id的书籍，其他页面状态都可以同步。
- [ ] 后台运行功能：unix：https://github.com/sevlyar/go-daemon 支持Windows但是最近没更新：https://github.com/takama/daemon
- [ ] 在终端显示图片 https://github.com/ploMP4/chafa-go
- [ ] 文件监视器 https://github.com/helshabini/fsbroker/
