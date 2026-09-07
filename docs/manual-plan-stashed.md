# ComiGo 内置 Manual 计划

状态：调研与讨论，未批准实施。更新：2026-09-05。
仅修改本计划；用户明确表示「可以开始改」后，才修改实现、安装依赖或迁移内容。不提交。

## 1. 已确认要求

- `/manual/` 公开访问，作为 ComiGo 内置用户手册。
- `docs/manual` 为内容唯一来源；初期整理、补全 `docs/wiki` 的半成品 Markdown，上线确认后再移除或精简旧 Wiki。
- 中、日、英三种语言；默认跟随 ComiGo 当前语言，其他语言默认英文，允许手动修改 ComiGo / Manual 语言。正文翻译独立维护，不塞进主应用翻译文件。
- 继承 ComiGo 主题，而非仅支持手册自己的亮色 / 暗色。
- 宽屏显示左侧目录与搜索、中央正文、右侧页内导航；手机左侧目录改为 Menu 按钮打开 drawer，页内导航改为可展开的快速跳转；移动端仍须能搜索。
- 代码块提供复制按钮。示例图片只用 HTTPS 外链，不下载、内联或打包进二进制。
- 日志不纳入 Manual 搜索。搜索范围限定为正式手册内容。
- 允许独立技术栈，减少与阅读器业务代码耦合。
- 视觉参考：[Omarchy Manual](https://omarchy.org/manual/)；[mise 安装文档](https://mise.jdx.dev/installing-mise.html)。借鉴布局，不复制正文、品牌或无关功能。

## 2. 当前代码与需求的差异

- 项目已有 Echo、templ、Alpine.js、i18next、CSS 主题变量及 `go:embed` 静态资源机制；无需新增运行时文档服务。
- 当前工作区已存在 `templ/pages/manual/`、页面级 JS/CSS 和 `assets/static/manual-content/*.json`。正文按语言加载 JSON，界面复用主题与 i18next，搜索是本地字符串匹配。
- 当前未找到根目录原计划、`docs/wiki` 或 `docs/manual`。不能将现有 JSON 擅自认定为最终内容源，也不能据此删除既有实现。
- `AGENTS.md` 描述的是现有 templ + JSON、无需单独构建的方案；若最终选择独立构建或 Markdown 内容源，实施时须同步修订这项约定。
- 本轮保留全部已有源码及未提交修改，只记录可选方向。Wiki 来源位置需用户确认。

## 3. 新增方案比较

以下维护成本和适配程度是结合项目需求的判断，不是实测性能排名。

| 方案 | 适合之处 | 主要代价 / 风险 | 定位 |
| --- | --- | --- | --- |
| 保留现有 templ + 页面级 JS，增加 Markdown 内容生成 | 最容易复用完整主题与语言；不新增前端框架；现有 UI 可保留 | 仍需维护目录、搜索、复制等代码；必须补齐 Markdown → 页面数据流程，避免手工双写 MD 和 JSON | 当前代码的最小改动选项 |
| VitePress | mise 实际使用；有多语言目录、浏览器本地搜索，适合接近参考站的文档体验 | 增加独立构建与 Vue 前端资源；主题和 App 语言需适配；中日文搜索分词需验证 | 优先追求 mise 体验时选 |
| Astro Starlight | 内置多语言路由、缺译回退、Pagefind 本地全文搜索；适合持续整理和翻译 Markdown | 增加 Astro 构建；需将内容加载路径接到 `docs/manual`；主题 / 语言适配不是自动完成 | 独立文档方案优先候选 |
| Hugo Book | Markdown 静态输出、多语言、移动布局；主要功能不依赖 JS；构建工具可与 Go 项目并存 | 仍需安装 Hugo 和管理主题版本；搜索中日文效果、复制和指定手机交互要验证，不能视为全套满足 | 偏好 Hugo 工具链时备选 |
| Docsify | 直接在浏览器渲染 Markdown，无文档构建步骤；有搜索及复制插件 | 正文依赖 JS；搜索需读取文档建立 / 缓存索引；双侧导航和主题仍要适配，所有脚本必须本地随包提供 | 明确拒绝新增构建时备选 |
| Docusaurus | 有多语言体系，适合大型文档站 | 新增 React 文档栈；官方重点支持 Algolia，本地搜索依赖社区方案；当前没有需要它的额外需求 | 暂不优先 |

依据：[VitePress 国际化](https://vitepress.dev/guide/i18n)、[本地搜索](https://vitepress.dev/reference/default-theme-search)、[mise 配置源码](https://github.com/jdx/mise/blob/main/docs/.vitepress/config.ts)；[Starlight 国际化](https://starlight.astro.build/guides/i18n/)、[搜索](https://starlight.astro.build/guides/site-search/)；[Hugo Book](https://github.com/alex-shpak/hugo-book)；[Docsify](https://github.com/docsifyjs/docsify)、[插件说明](https://github.com/docsifyjs/docsify/blob/develop/docs/plugins.md)；[Docusaurus 国际化](https://docusaurus.io/docs/i18n/introduction)、[搜索](https://docusaurus.io/docs/search)。

## 4. 影响选型的几个要点

### 搜索：本地优先，单独检查中日文

- `/manual/` 公开不等于每个 ComiGo 实例都能被互联网爬取。建议搜索随程序分发，不依赖外部搜索服务；mise 当前使用 Algolia，不能直接照搬这部分配置。
- VitePress 内置 MiniSearch，支持自定义分词；搜索界面能翻译，不代表无空格的中文、日文能达到预期检索效果。
- Pagefind 按 HTML `lang` 分别构建 / 加载索引；中日文分词需要 extended 版本，官方说明 `npx pagefind` 默认使用它。Starlight 最终锁定的版本与构建依赖仍须核实。
- Pagefind 也可单独用于其他静态 HTML 方案，不必为搜索更换整个页面框架；但不能直接把当前浏览器加载 JSON 后才出现的正文当作构建期 HTML 索引。
- 索引只包含正式手册；排除日志、草稿和开发文档。先按当前语言搜索，不增加跨语言混合检索。

依据：[MiniSearch](https://lucaong.github.io/minisearch/)、[Pagefind 多语言与分词](https://pagefind.app/docs/multilingual/)。

### 主题与语言：共享偏好，不引入阅读器业务

- 主题包括 `retro` 等多套配色、随机主题解析和自定义色值，不只是 dark / light；独立文档框架需要映射背景、文字、边框、强调色和代码块颜色。
- 优先复用现有主题变量和实际生效的主题值；避免在文档框架中复制维护一套完整配色表。独立页面不应为取主题而加载书架、阅读器或 SSE 业务。
- 正文按 `docs/manual/{zh-CN,ja-JP,en-US}/` 分文件；同一章节使用相同 slug。界面少量文案另放 Manual 自己的语言资源。
- 不建议将三种正文嵌入同一 Markdown / 模板；编辑和翻译审阅更困难。JSON 如保留，应为生成产物，不再手工维护第二份正文。
- 需确认手动修改 Manual 语言是否同步修改 App，以及语言深链接与已保存偏好的优先级；不能只依靠框架默认行为决定。

依据：[VitePress 主题扩展](https://vitepress.dev/guide/extending-default-theme)、[Starlight 样式定制](https://starlight.astro.build/guides/css-and-tailwind/)，以及当前项目主题与 i18next 代码。

### 打包：区分构建依赖和运行依赖

- VitePress / Starlight / Hugo 都可采用构建期生成、Go 嵌入静态产物的方式；用户运行 ComiGo 不需要另装这些工具或启动服务。
- 只嵌入手册 HTML、必要 CSS/JS、搜索索引等发布产物；不嵌入依赖目录、源码映射或示例图片。优先沿用系统字体，不增加字体包。
- 图片保留远程 URL、说明文字及 alt；离线无法显示图片是明确限制。建议正文与本地搜索在 ComiGo 服务可用但无外网时仍可用。
- 尚未构建候选方案，不给出未经测量的 KB / MB 数字。批准后的验证比较产物总大小、二进制增量、首次页面与首次搜索传输量。

依据：[VitePress 静态部署](https://vitepress.dev/guide/deploy)、[VitePress 移除默认字体](https://vitepress.dev/guide/extending-default-theme#using-different-fonts)。

## 5. 建议与实施顺序（待批准）

建议只在三个方向中选：最小改动选现有实现；重视 mise 体验选 VitePress；重视独立 Markdown 文档、多语言和本地搜索的完整配套，优先验证 Starlight。不并行维护两套手册，也不先写通用文档引擎。

1. 确认技术方向、Wiki 来源、语言规则和缺译处理，再冻结计划。
2. 获准后，仅用同一篇含三级标题、长代码块、表格、外链图片的中日英样例验证所选方案；不先迁移全部内容。
3. 验证公开路由、手机 drawer / 下拉目录、代码复制、主题继承、语言切换、当前语言搜索和静态嵌入；评估体积后再决定是否继续该方案。
4. 整理 `docs/manual`：安装与快速开始、阅读操作、书库管理、桌面使用、部署与远程访问、常见问题。逐项核对现有代码，不把 Wiki 草稿当作已确认功能。
5. 完善三语内容及链接检查，接入发布构建；新手册上线验收后才精简或移除旧 Wiki。

验收重点：未登录可访问；320px 至宽屏无横向溢出；键盘可操作搜索及 drawer；HTTP 局域网环境复制失败有可用提示 / 手动选择；中文无空格词组、日文复合词、英文及 CLI 参数可检索；深链接刷新正确；主题无明显闪白；无外网时正文与搜索可用；产物不含示例图片。

## 6. 待确认

- 选择现有实现、VitePress 还是 Starlight？Hugo Book / Docsify / Docusaurus 作为参考保留，不同时试做。
- `docs/wiki` 原文现在位于哪个目录或分支？现有 Manual JSON 是否也是应保留的内容来源？
- Manual 手动切换语言是否与 App 双向同步？明确语言的链接是否优先于已保存偏好？
- 未翻译章节是显示默认正文并标注缺译，还是首版必须三语齐全？「不支持的界面语言回退英文」与「章节缺译回退」是两个独立规则。
