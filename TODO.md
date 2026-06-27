# ComiGo TODO

更新日期：2026-06-18

本文件只记录当前仍需要推进的事项。已完成或已失效的旧条目不继续放在待办正文里，避免后续按过期信息排期。

## 当前现状

- 最新 GitHub Release：`v1.2.36`，标题为“优化书架扫描功能”，发布时间为 2026-06-07，非 draft、非 prerelease。
- 当前开发版本：`v1.2.37`。`v1.2.36` 之后的提交已经包含远程 Comigo 书库基础能力、运行时显示状态统一、扫描失败缓存改进、Tailscale 状态接口、设置页远程书库体验优化和依赖更新。
- 已具备：Web 书架、卷轴阅读、翻页阅读、本地 `/reader`、便携 HTML、PWA、Tailscale、账号/权限、阅读历史、书签、上传到书库、服务器设置页、日志页、OPDS、TUI 预览、Windows 托盘/文件关联、SQLite 兼容层、Postgres 适配雏形、SMB/WebDAV 等 VFS 基础能力。
- 近期已完成：Kitty TUI 图片显示优化、随机主题与主题刷新、内部路径输出收紧、图片参数边界校验、Postgres 后端与重扫控制、书架扫描优化、远程 Comigo 书库索引与实时代理基础链路。

## 下一版优先处理

### P0：`v1.2.37` 发布前回归

- [ ] 远程 Comigo 书库端到端回归。覆盖添加 `http(s)` 主页 URL、认证失败、分组/子书架、封面、scroll/flip/player 阅读、原文件下载、EPUB 下载、阅读历史、书签新增/删除和远端不可达错误提示。
- [ ] 确认远程 Comigo 书库只保存索引，不在本地保存原始文件；阅读和下载必须实时代理远端资源。
- [ ] 复查远程书库配置变更后的状态清理。更换 URL、用户名、密码、只读模式后，旧连接、旧分组和旧权限状态不能继续被复用。
- [ ] 设置页与 SSE 重载回归。修改书库、Tailscale、显示相关配置后，前端状态、按钮可用性和刷新提示应保持一致，不出现重复提示或旧状态残留。
- [ ] 扫描/重扫真实书库回归。覆盖本地目录、压缩包、失败缓存、重新扫描单书/目录、远程书库索引刷新和空目录删除逻辑。
- [ ] 发布前 smoke 脚本补齐并运行：首页、书架、scroll、flip、player、settings、reader、OPDS、TUI 启动、静态导出页面。

### P1：安全与稳定性

- [ ] 统一前端文本写入边界。优先处理 toast、上传文件名、日志面板、远程错误信息等用户或远端可控文本，避免直接拼入 `innerHTML`。
- [ ] 把普通读接口里的全局清理副作用移出请求路径。`GetBook`、首页书架、子书架不应在普通 GET 中触发 `ClearBookWhenStoreUrlNotExist` / `ClearBookNotExist`。
- [ ] 扫描流程去掉包级全局状态。`tools/scan` 中的 `cfg/currentFS` 应收敛为显式扫描上下文，短期至少串行化扫描入口，避免上传扫描、手动重扫和远程索引刷新互相覆盖。
- [ ] VFS registry key 加入认证身份边界，或在远程书库配置变更后注销旧连接，避免同 host/path 更换用户或密码仍复用旧连接。
- [ ] 文件访问边界继续收紧。评估 Go `os.Root` 或等价封装，把文件操作限制在书库、缓存、配置目录内。
- [ ] 远程代理请求继续收紧。限制可代理路径、超时、响应大小和错误透传内容，避免远端异常拖垮本地服务或暴露敏感信息。

### P1：阅读体验

- [ ] 实现卷轴阅读加载策略重构。用户可见阅读模式只保留 `卷轴阅读 / 翻页阅读`，`无限卷轴 / 延迟加载 / 分页加载` 下沉为卷轴专用加载策略。
- [ ] 在线卷轴阅读增加延迟加载模式。默认首批加载 `scroll.pageLimit = 32` 页，其余页接近视口再设置图片 `src`。
- [ ] EPUB/PDF 阅读体验优化。PDF 可评估 `pdf.js`，EPUB 可评估 `epub.js` 或继续增强现有解析，目标是支持图文混排和更稳定的页码/目录。
- [ ] 本地 `/reader` 和在线 scroll/flip 的模式切换文案保持一致，但本地 reader 不引入服务器卷轴分页/延迟加载设置。
- [ ] 全局多端同步翻页继续整理。除了不同书籍 ID 的阅读内容，页面状态、页码和控制状态应能跨标签/设备同步。

### P1：书库与元数据

- [ ] 编辑系列和书籍元数据。支持 `ComicInfo.xml`、CBZ/CBR 元数据、EPUB 元数据，提供搜索、摘要、系列信息编辑入口。
- [ ] 可选计算文件哈希，用于查找重复文件；默认不影响普通扫描速度。
- [ ] 文件管理功能：删除、移动、重命名、重新扫描单书/单目录，并明确只读模式下的禁用行为。
- [ ] 支持 rar 密码、损坏文件、扩展名错误、固实压缩包和 7z 的更准确错误提示与跳过策略。
- [ ] 崩溃恢复和恶意存档防护。重点限制解压炸弹、极端图片尺寸、异常目录结构和超大 metadata。

## Wiki 与文档 TODO

- [ ] `sample/wiki/11-Changelog-ZH.md`、`12-Changelog-JA.md`、`13-Changelog-EN.md`：同步 `v1.2.32` 到 `v1.2.36`，包含 Kitty TUI、主题刷新、随机主题、Postgres 适配、重扫控制和书架扫描优化。
- [ ] 为 `v1.2.37` 准备 changelog 草稿：远程 Comigo 书库、运行时状态统一、远程书库设置页、Tailscale 状态接口、扫描失败缓存和依赖更新。
- [ ] `sample/wiki/08-Reading-and-Library.md`：补充本地书库、远程 Comigo 书库、本地 `/reader` 三者差异，以及远程阅读/下载实时代理、不保存原始文件的边界。
- [ ] `sample/wiki/09-Remote-and-Security.md`：补充 Tailscale、ZeroTier、局域网 IP、只读模式、上传开关、远程 Comigo 认证和文件不会上传到在线 reader 的边界说明。
- [ ] `sample/wiki/03-Try-it-Online.md`：根据最新 `/reader` 行为更新离线/PWA/便携 HTML 说明，特别说明 `file://` 便携 HTML 不连接服务器。
- [ ] `README.md`、`README_ZH.md`、`README_JP.md`：同步稳定的远程 Comigo 书库、PWA、OPDS、TUI 和 reader 说明。
- [ ] 新官网与使用文档继续推进，可参考 Omarchy 的“先解释场景，再给入口”结构。Discord/社区入口可以放到官网而不是核心 App。

## 平台与发布

- [ ] 自动发版：Goreleaser、GitHub Actions、校验和、签名、Changelog 自动生成。
- [ ] 发布流程整理：`dev` 发版后同步 `master`，补齐 tag、本地/远端 release 校验和 release notes 检查清单。
- [ ] 更新提示与自动更新。桌面端优先支持提示新版本；自动替换当前程序需分平台评估，不能影响包管理安装方式。
- [ ] Linux 软件源与包格式：Debian、RPM、Homebrew/Linuxbrew、AUR 或其他发行渠道。
- [ ] 后台运行命令：`-start` 后台启动、`-stop` 停止进程。需要区分 CLI、桌面托盘和嵌入式宿主。
- [ ] 跨平台 GUI：继续评估 Flutter + GoMobile、Wails、Electron 宿主。移动端和桌面端启动链路必须与 Tailscale、托盘、Windows 注册表逻辑解耦。
- [ ] Wails3 alpha 分支继续跟进正式版 API 变化，发布前补齐打包、签名和跨平台回归。
- [ ] desktop 版支持系统托盘，复用现有托盘能力时保持 Wails、tray、CLI 启动链路解耦。
- [ ] 编译并验证 Wails desktop 的 iOS 与 Android 版本，确认移动端启动、资源加载和书库路径选择边界。
- [ ] 继续测试 Wails3 是否兼容内置 Tailscale；兼容前继续隐藏 Wails 环境的远程访问/Tailscale 配置。
- [ ] 终端图片显示继续增强。TUI 侧可继续评估 Kitty/iTerm/Sixel/chafa 等协议，但要保持非图片终端可用。

## 产品方向

- [ ] AI 编写、程序控制、可互动日式漫画创作工具。不要继续把作者工具塞进 `BookInfo` / `PageInfo`；应新增 work/project 模型、manifest/DSL、`/studio`、`/work/:id` 和可导出的交互式 HTML runtime。
- [ ] 交互漫画示例作品。先做一个短篇示例，覆盖分镜、对白、选择、音效/背景音乐、局部动画、状态变量和单文件 HTML 导出。
- [ ] 防剧透效果、回忆模式、特殊背景、背景音乐等阅读增强能力，可以作为交互漫画 runtime 的可选效果，而不是普通阅读器主流程的硬依赖。
- [ ] 调用第三方或本地 AI/API 处理图片压缩、放大、翻译、文字识别时，应先定义离线/隐私边界和用户确认流程。

## 后续工程债

- [ ] Settings API 去重复。string/bool/number、数组配置、插件启停可以抽统一更新流程，减少只读检查、保存、SSE 通知的重复代码。
- [ ] 配置字段 setter 统一化。反射赋值、数组字段、路径字段校验应集中，方便 Settings 和 CLI 共享。
- [ ] 扫描本地/远程流程合并。候选收集、单文件入库、目录入库、失败缓存、封面提取可以共享一套流程。
- [ ] `GetPictureData` 拆分为读取源数据和图片变换两层，降低远程、本地、压缩包、PDF、目录混在一个函数里的维护成本。
- [ ] 内存 Store 与 SQLite/Postgres Store 的书组生成逻辑抽成无副作用服务，避免多套存储实现行为漂移。
- [ ] SQLite/Postgres 如果进入主路径，需要减少 `ListBooks` 的 N+1 查询，补批量查询和大书库性能测试。
- [ ] OPDS 继续补兼容性测试，覆盖常见客户端、认证、封面、分页和远程书库。
- [ ] 前端 store 拆分。`global_store.js` 里主题、设备、播放器、URL、书签 API 等职责应继续拆小。

## 已完成但保留索引

- [x] OPDS 协议基础支持。
- [x] 上传页面迁移到设置页书库设置。
- [x] 书架按照最近阅读时间排序，无阅读进度时按文件修改时间排序。
- [x] PWA 支持与 reader 专用 Service Worker。
- [x] 便携 HTML 导出。
- [x] 阅读历史和书签持久化。
- [x] Tailscale 集成。
- [x] 账号、访问权限与只读模式基础能力。
- [x] Windows 托盘、文件关联和右键菜单基础能力。
- [x] TUI 预览、终端阅读和 Kitty 图片模式基础能力。
- [x] SQLite 兼容层和 `auto_vacuum` 初始化。
- [x] Postgres 后端与 sqlc 适配雏形。
- [x] SMB/WebDAV 远程书库基础支持。
- [x] 图片生成和图片变换参数基础边界校验。
- [x] 内部路径从公开 JSON 中收紧。
- [x] 随机主题与主题刷新。
- [x] 远程 Comigo 书库基础索引与实时代理链路。
