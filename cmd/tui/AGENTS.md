# TUI Agent Notes

本文件记录 `cmd/tui` 终端图片渲染相关知识、当前现状和排查链接。后续修改 TUI 预览区或终端阅读时，先阅读并按实际结果更新本文件。

## 当前现状

- ANSI/Halfblocks：当前高度和布局相对稳定，但分辨率低，只适合作为通用兜底。
- iTerm2：可完整显示图片；预览区和终端阅读区的底部留白已通过“贴边轴固定、另一轴 auto”的 OSC 1337 参数策略修复/缓解。
- Kitty：预览区正常；终端阅读回到 Unicode placeholder 路径，让图片随 Bubble Tea 文本行刷新，避免页码和图像层不同步。最新现象不是“底部错位一行”，而是底部附近额外出现两条图像行，两条都来自同一图片第一行。
- Ghostty：支持 Kitty graphics protocol；终端阅读回到 Unicode placeholder 路径。预览区不走 placeholder，改用 Kitty overlay，避免裁切或尺寸过小。进入终端阅读时需要一次性清理预览 overlay，否则书架封面可能残留到阅读画面。
- WezTerm：对齐到 iTerm2 inline image 协议；Kitty/placeholder 路径在预览区和终端阅读都出现过布局错乱、页码与图片不同步或重复显示的问题。

## 实现要点

- `coverPreviewState.Lines` / `terminalReaderState.Lines` 只应保存随 TUI 文本布局输出的可见行。
- Kitty Unicode placeholder 的图像传输和 virtual placement 控制序列不是可见行，不能参与居中、裁切或宽度计算。
- iTerm2 协议是独立 inline image 层，使用 overlay 路径绘制；清理残留时只对 iTerm2 协议做 ECH/空格覆盖。
- WezTerm 走 iTerm2 协议时，尺寸估算仍使用 WezTerm 的字符格比例；OSC 1337 只固定贴边轴，另一轴传 `auto`，减少协议内部二次等比缩放带来的底部留白。
- Halfblocks 使用文本模式，不走图像层清理；它的尺寸计算保留半块字符的 1:2 逻辑。
- 涉及终端差异的特殊处理应按终端类型局部生效，不要把 Ghostty/iTerm2 的 workaround 扩散到 Kitty 或 ANSI。
- 可设置 `COMIGO_TUI_IMAGE_DEBUG=1` 打开 TUI 图片渲染诊断日志，观察协议、区域大小、源图大小、输出行数和控制序列字节数。
- Kitty placeholder 行包含 `U+10EEEE` 和组合符；未超宽时不要做通用 ANSI 截断重写，避免破坏行/列/ID 编码。
- Kitty 终端阅读的无效尝试需要及时撤回，避免刷新策略堆叠后掩盖真正原因。

## 终端图片协议资料快照（2026-05-17）

### 协议概览

| 协议 | 机制 | 支持程度与风险 | Comigo 当前用途 |
| --- | --- | --- | --- |
| Kitty graphics protocol | APC 图像协议，支持 image id、placement、z-index、删除、尺寸 `c/r`，也支持 Unicode placeholder。 | 现代终端支持度较高；能力最强，但不同实现对 overlay、placeholder、清理、动画的细节兼容性不同。 | Kitty/Ghostty 终端阅读优先尝试；Kitty 原生终端预览继续走 placeholder；Ghostty 预览临时走 overlay。 |
| Kitty Unicode placeholder | 用 `U+10EEEE` 加颜色和组合符表示图片网格；图片像普通文本一样被 TUI 布局移动。 | 适合 Bubble Tea 这类文本 TUI；但要求控制序列和可见 placeholder 行严格拆分，终端实现差异会导致裁切或错位。 | 终端阅读的 Kitty/Ghostty 实验路径；WezTerm 已改走 iTerm2 协议。 |
| iTerm2 Inline Images / OSC 1337 File | OSC 1337 `File`，base64 文件内容，`inline=1` 时显示在当前会话。支持 cell / px / percent / auto 尺寸参数。 | iTerm2 原生最稳定；WezTerm 也支持。multiplexer 场景可能不完整。 | iTerm2 与 WezTerm 当前主路径；保留 iTerm2 专属清理和 `doNotMoveCursor` 风格处理。 |
| Sixel | 传统 DEC bitmap 图形协议，需要把图片编码为 sixel 数据流。 | macOS 上 iTerm2/Rio 实测未能稳定显示：iTerm2 会输出可见字符，Rio 不显示图片。WezTerm 官方标为 experimental。 | 当前撤回手动补充支持，不作为 TUI 可用路径。 |
| ANSI/Halfblocks | 使用普通 ANSI 颜色和半块字符拼近似图像。 | 最稳、最低要求；分辨率低但不会产生独立图像层残留。 | 通用兜底；当前不作为 WezTerm 默认路径。 |

### 常见终端支持矩阵

| 终端 | 官方/当前资料支持 | 当前建议 |
| --- | --- | --- |
| iTerm2 | 官方支持 OSC 1337 inline images，参数支持 cell、px、percent、auto；任何 macOS 支持的图片格式可 inline 显示。 | 用 iTerm2 协议，不走 Kitty。底部留白问题改为只固定贴边轴，另一轴传 `auto`。 |
| Kitty | 官方 Kitty graphics 原生实现；支持 placement、`C=1`、删除命令、Unicode placeholder。 | 终端阅读优先 placeholder，预览区 placeholder。不要再把已失败的 overlay 清理逻辑扩散回来。 |
| Ghostty | 官方 Feature 列表明确支持 Kitty graphics protocol。未在官方资料中看到 iTerm2/Sixel 作为主能力。 | 终端阅读按 Kitty placeholder 实验；预览区目前用 Kitty overlay 避开 placeholder 裁切/过小。 |
| WezTerm | 官方 Feature 列表列出 iTerm2 compatible image protocol、Kitty graphics、Sixel experimental；`imgcat` 文档说明 iTerm2 协议在 multiplexer 中不完整；changelog 说明 Kitty Image Protocol 默认启用但动画未实现。 | 预览区和终端阅读都使用 iTerm2 协议，避开当前实测不稳定的 Kitty/placeholder 路径。 |
| Windows Terminal | 官方 1.22 Preview 先加入 Sixel image support，1.23 Preview 公告确认 stable 1.22 已包含该能力；显示仍需要 `libsixel`/`chafa` 等编码器输出 sixel。 | 非当前 macOS 主测试目标；如未来支持 Windows TUI 图片，优先按 Sixel 检测。 |
| xterm/foot/mlterm 等 Sixel 终端 | Sixel 是通用传统协议，具体终端常需编译/配置/版本支持；xterm 需确认是否启用 sixel。 | 可作为未来 Linux 终端扩展方向，不应影响当前 Kitty/iTerm2 路径。 |

### Go 库现状

- `github.com/blacktop/go-termimg` 提供协议检测函数和 `DetermineProtocols()`；文档给出的优先顺序是 Kitty -> iTerm2 -> Sixel -> Halfblocks，符合当前“原生协议优先、Halfblocks 兜底”的方向。
- 该库的自动检测只能作为入口，Comigo 已经验证过 WezTerm/Ghostty/Kitty 在不同区域的实际表现不一致，因此仍需要保留终端与区域维度的局部策略；WezTerm 当前主动覆盖为 iTerm2。

## 2026-05-17 处理记录

- 新增 `Setup` 字段保存 Kitty 图像传输/virtual placement 控制序列，`Lines` 只保存可见 placeholder 行。
- 拆分 Kitty 输出时，必须保留 placeholder 行开头的 SGR foreground color；该颜色编码 image id，不能放进 `Setup`。
- 封面预览和终端阅读的 Kitty placeholder 路径会在主文本前输出 `Setup`，再由正文中的 placeholder 行定位图片。
- 原生图像协议的尺寸计算改为按终端字符格像素比例估算，但不做交互式 CSI 查询，避免污染 Bubble Tea 输入流。
- Kitty overlay 路径在 Kitty/Ghostty/WezTerm 中出现页码已变但图片不刷新的问题；已撤回终端阅读 overlay、`C=1`、loading 清理等特殊处理。
- Kitty/Ghostty 终端阅读继续使用 Kitty 原生图像，但改回 Unicode placeholder。WezTerm 已改用 iTerm2 inline image 协议。
- iTerm2/WezTerm 底部留白问题已按实测线索处理：WezTerm 走 iTerm2 协议时仍使用 WezTerm 字符格比例估算；OSC 1337 输出只固定贴边轴，避免固定宽高框内二次等比缩放。
- Kitty 底部“多出两条第一行”的新判断：它不像整体文本错位，更像旧 row=0 placeholder/virtual placement 没有被完全删除，或当前 placeholder 输出在底部附近重复触发 row=0 解析。2 列安全边距、边框内嵌显示、翻页/缓存命中时强制清 Kitty 图层、synchronized output 都已实测无效并撤销。
- Ghostty 从书架进入终端阅读时，预览区 Kitty overlay 需要在阅读页第一帧前发送 Kitty delete-all 指令；该清理只在进入终端阅读前一次性触发，不扩散为 Kitty 翻页刷新策略。
- Sixel 手动补充模式已撤回：iTerm2 切换后输出可见字符，Rio 不显示图片；后续不要仅凭终端资料重新开放，必须先用实际 TUI 页面验证。

## 参考资料

- Kitty graphics protocol: https://sw.kovidgoyal.net/kitty/graphics-protocol/
- Kitty Unicode placeholders: https://sw.kovidgoyal.net/kitty/graphics-protocol/#unicode-placeholders
- iTerm2 Inline Images: https://iterm2.com/3.4/documentation-images.html
- WezTerm features: https://wezterm.org/features.html
- WezTerm iTerm image protocol: https://wezterm.org/imgcat.html
- WezTerm changelog image notes: https://wezterm.org/changelog.html
- Ghostty features: https://ghostty.org/docs/features
- Ghostty external protocols: https://ghostty.org/docs/vt/external
- Windows Terminal 1.22 Preview Sixel: https://devblogs.microsoft.com/commandline/windows-terminal-preview-1-22-release/
- Windows Terminal 1.23 Preview / stable 1.22 Sixel: https://devblogs.microsoft.com/commandline/windows-terminal-preview-1-23-release/
- go-termimg docs: https://pkg.go.dev/github.com/blacktop/go-termimg

## 排查路线

1. 先确认当前终端自动选择的协议是否符合预期。
2. 再看日志中的源图尺寸、目标区域、计算后的 `ImageW/ImageH` 是否过小。
3. Kitty/Ghostty 出现错行时，先确认控制序列没有进入可见行裁切。
4. iTerm2/WezTerm 留空时，先确认贴边轴是否正确传 `width=auto` 或 `height=auto`，再判断是否需要按终端局部微调字符格比例。
5. 只有在确认某终端协议路径不稳定后，才为该终端添加局部 fallback。
6. Kitty 终端阅读下一步只做无行为变更的诊断：记录每行 placeholder 数量和宽度、`Lines` 总数、首尾两行内容签名，并检查是否还有额外控制序列或宽度计算误差。
