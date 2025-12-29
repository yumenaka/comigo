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
- **UI 组件**: [Flowbite](https://flowbite.com/docs/components/gallery/) - UI 组件库
- **交互增强**: 
  - [htmx](https://htmx.org/examples) - AJAX 库
  - [Alpine.js](https://alpinejs.dev) - 轻量 JavaScript 框架
- **样式框架**: [Tailwind CSS](https://tailwindcss.com/docs/flex) - 实用优先的 CSS 框架
- **国际化**: [i18next](https://www.i18next.com) - 国际化解决方案
- **图标库**: [Xicons](https://www.xicons.org/#/) - 图标集合
- **热重载**: [Air](https://github.com/air-verse/air) - Go 应用热重载工具
- **js运行时**: [Bun](https://github.com/oven-sh/bun) - 高性能 JavaScript 运行时

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
- [x] 阅读历史记录（跳转到上一次阅读的最远页）
- [x] 阅读历史持久化，meta文件，
- [x] 下载为单个html文件
- [x] 拆分内存存储与数据库存储，为支持samba与S3等文件系统做准备。
- [x] windows:系统托盘图标、后台运行 https://github.com/getlantern/systray
- [x] 系统托盘
- [x] 单例模式
- [x] Windows：可选择修改注册表为默认打开方式、右键文件夹菜单
- [x] 阅读历史功能

### 开发中

- [ ] 下载源文件，封面文件缓存
- [ ] 画漫画，当作示例漫画
- [ ] 性能优化，缩略图预生成与缓存。错误文件跳过。错误文件跳过。
- [ ] 多书库限制：每个库都有根文件夹，且任何库都不能共享其路径的任何部分
- [ ] 根据父文件夹的最后修改时间来决定是否重新扫描书库（似乎有系统差异）？
- [ ] 新官网,Discord频道,使用文档，参考 https://omarchy.org/ 。内置帮助文档?
- [ ] 编辑系列和书籍的元数据。支持解析Comicinfo.xml的CBZ,CBR文件。EPUB 文件的元数据。搜索系列和书籍，查看元数据、摘要等。
- [ ] 多设备访问：转换为epub？支持OPDS协议，方便用户在各种设备（如Kindle、Kobo等）上访问电子书。https://specs.opds.io/opds-1.2
- [ ] 帮助按钮与前置的简单菜单，优化翻页效果。
- [ ] 指定首页书籍功能
- [ ] 完整书签功能。
- [ ] 可选的计算文件哈希值，查找重复文件。
- [ ] SQLite 删除数据后不释放磁盘空间，只标记为空闲。建库时开启 auto_vacuum，可以防止文件持续膨胀
- [ ] 侧栏加返回书架，切换全屏按钮，方便操作。或者让Header的显示层级大于侧栏。翻页模式 header 合并到下部进度条？
- [ ] 自动更新，下载最新版本，替换当前程序。
- [ ] 下载文件到书库目录
- [ ] 书架按照最近阅读时间排序，无阅读进度的书籍，以文件修改时间排序
- [ ] 前后端的语言设定默认值同步？
- [ ] 上传页面挪到设置页面-书库设置。改造上传功能，可选上传到下拉框指定的书库。没有默认书库则不可上传。
- [ ] 官网自动探测浏览器平台，提供合适的平台版本（参考Audacity） 
- [ ] 支持smb、webdav文件系统
- [ ] 网页端日志查看 50%
- [ ] 注册为文件默认打开类型，简单托盘图标，gui界面。
- [ ] 手动或自动检测新版本提示，然后可以试着自动更新新版本（win与macos），最后是各种linux软件源
- [ ] 网页端：浏览器快捷键(50%)。
- [ ] 网页端：自动化测试，修改后自动测试基本功能。
- [ ] 网页端：卷轴模式页数同步体验优化。
- [ ] 翻页模式：滚轮滑动翻页
- [ ] 滑动模式：可快捷键调速的自动翻页
- [ ] cli 交互，tui支持
- [ ] 访问权限控制
- [ ] PWA 支持
- [ ] 系统监控（CPU、内存）
- [ ] wasm模式
- [ ] 自动更新
- [ ] 文件监控 https://github.com/sgtdi/fswatcher
- [ ] 用户系统增强
- [ ] Shell 交互
- [ ] 文件管理
- [ ] EPUB/PDF 阅读优化、
- [ ] 防剧透效果、回忆模式、特殊背景、背景音乐etc
- [ ] 第三方登录 https://github.com/markbates/goth
- [ ] 多语言：中文、英文、日文版toml配置文件注释
  
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
- [ ] 跨平台 GUI（Flutter + GoMobile或 Wails）PWA模式。wail debug https://github.com/wailsapp/wails/issues/3050
- [ ] 更新提示，自动更新自动更新(github.com/jpillora/overseer) 包管理支持。[self update 功能](https://github.com/minio/selfupdate)
- [ ] 文件夹监控(fsnotify)，https://github.com/helshabini/fsbroker/
- [ ] 用户系统、访问密码，流量限制 comigo后台：有几台设备在线，阅读文件，阅读页数，当前用户状态、阅读书籍、阅读进度、阅读时间、服务器状态 注册，阅读记录，
- [ ] shell 互动（<https://github.com/rivo/tview> ）
- [ ] 子命令，download rar2zip
- [ ] 支持rar压缩包密码。处理损坏文件，扩展名错误的文件，固实压缩文件（7z）。更准确的文件类型判断。
- [ ] 崩溃后恢复，恶意存档处理。
- [ ] 编写测试
- [ ] 命令行交互
- [ ] 调用第三方API
- [ ] 文件管理，删除。
- [ ] Debian，RPM包（<https://github.com/goreleaser/nfpm）>
- [ ] 优化epub与PDF阅读体验，支持图文混排（pdf.js与epub.js）
- [ ] 显示服务器log：[web终端示例](https://zenn.dev/ikedam/articles/2e078bfc2a4cb6)
- [ ] -start 参数，后台运行。-stop参数，停止后台运行的进程。
- [ ] 自定义js与css代码块功能。
- [ ] [使用 Go1.24 的 os.Root 类型](https://antonz.org/go-1-24/)，将[文件操作限制在特定目录](https://go.dev/blog/osroot)，以防止攻击者通过转义或相对路径非法访问文件
- [ ] 自动发版功能 [goreleaser](https://goreleaser.com/)  [github-action](https://dev.to/hadlow/how-to-release-to-homebrew-with-goreleaser-github-actions-and-semantic-release-2gbb)
- [ ] 添加[数据验证](https://dev.to/leapcell/validator-complex-structs-arrays-and-maps-validation-for-go-34ni)。
- [ ] SteamDeck支持（网页支持手柄操作）鼠标滚轮对应
- [ ] 同步翻页 -> 全局多端同步跟踪页面状态，除了不同id的书籍，其他页面状态都可以同步。
- [ ] 后台运行功能：unix：https://github.com/sevlyar/go-daemon 支持Windows但是最近没更新：https://github.com/takama/daemon
- [ ] 在终端显示图片 https://github.com/ploMP4/chafa-go
- [ ]  OpenID Connect 登录 https://github.com/zitadel/oidc  https://tailscale.com/community/community-projects/tsidp


## history
- 2025-11-12: v1.1.0 发布，支持下载为单个网页文件，Tailscale远程连接，多书架优化
- **新功能：**
1. 自动保存与恢复阅读进度，让阅读体验更加连贯。
2. 内置 Tailscale 远程连接功能，轻松实现跨设备访问。
3. 支持多个书库分批展示，加载更高效。

**优化：**
1. 同一浏览器的不同标签页之间也可同步翻页操作。
2. 自动忽略以「.」开头的隐藏文件。
3. 命令行模式下，打印帮助或版本信息后直接退出。
4. 修复网络较慢时书架与设置页面出现闪烁的问题。

**新機能：**
1. 読書進捗の自動保存・復元に対応し、より快適な読書体験を実現。
2. 内蔵の Tailscale リモート接続機能で、デバイス間アクセスがより簡単に。
3. 複数の書庫を分割して表示し、読み込み効率を向上。

**最適化：**
1. 同一ブラウザ内の別タブ間でもページめくり操作を同期。
2. 「.」で始まる隠しファイルを自動的に無視。
3. コマンドラインでヘルプまたはバージョン情報を表示後、自動的に終了。
4. ネットワーク速度が遅い場合に発生する書棚および設定ページのちらつき問題を修正。

**New Features:**
1. Automatically save and restore reading progress for a seamless experience.
2. Built-in Tailscale remote connection for easy cross-device access.
3. Support for displaying multiple libraries in batches for improved loading efficiency.

**Improvements:**
1. Page flipping is now synchronized across different tabs in the same browser.
2. Hidden files (starting with “.”) are now automatically ignored.
3. In command-line mode, the program now exits immediately after printing help or version info.
4. Fixed flickering issues on the bookshelf and settings pages when the network is slow.