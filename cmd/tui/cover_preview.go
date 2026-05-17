package tui

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"strings"
	"time"

	termimg "github.com/blacktop/go-termimg"
	tea "github.com/charmbracelet/bubbletea"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	coverutil "github.com/yumenaka/comigo/tools/cover"
	"github.com/yumenaka/comigo/tools/logger"
)

const (
	previewTitleLines           = 1
	previewFooterLines          = 1
	coverPreviewMinWidth        = 18
	coverPreviewMinHeight       = 6
	coverPreviewMinResizeHeight = 1200
	halfblocksScalePercent      = 150
)

type tuiImageProtocol = termimg.Protocol

// coverPreviewState 保存当前封面预览的渲染结果。
// Lines 用于 ANSI/Kitty 占位符这类随主文本一起输出的协议，Overlay 用于 iTerm2/Sixel 独立图像层。
type coverPreviewState struct {
	BookID   string
	Title    string
	Width    int
	Height   int
	ImageW   int
	ImageH   int
	Protocol tuiImageProtocol
	Setup    string
	Lines    []string
	Overlay  string
	Loading  bool
	ErrText  string
}

type coverPreviewMsg struct {
	requestID int
	state     coverPreviewState
}

// detectTUIImageProtocol 自动选择 TUI 封面预览协议。
// COMIGO_TUI_IMAGE 可强制覆盖，用于排查不同终端协议的显示差异。
func detectTUIImageProtocol() tuiImageProtocol {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("COMIGO_TUI_IMAGE"))) {
	case "", "auto":
		// iTerm2 的环境变量比通用检测更稳定，优先识别，避免误走 ANSI 回退。
		if isITerm2Terminal() {
			return termimg.ITerm2
		}
		if isGhosttyTerminal() {
			// Ghostty 预览区仍使用 Kitty 协议，但后续渲染改走 overlay，避开 placeholder 裁切。
			return termimg.Kitty
		}
		if isWezTermTerminal() {
			// WezTerm 的 Kitty/placeholder 路径在预览区仍不稳定，改用实测问题更少的 iTerm2 inline image 协议。
			return termimg.ITerm2
		}
		protocol := termimg.DetectProtocol()
		// Kitty 与 iTerm2 保留原生图像能力；其它图像层协议默认回退 halfblocks，降低残留风险。
		if protocol == termimg.Kitty || protocol == termimg.ITerm2 {
			return protocol
		}
		return termimg.Halfblocks
	case "off", "none", "false", "0":
		return termimg.Unsupported
	case "ansi", "halfblock", "halfblocks":
		return termimg.Halfblocks
	case "kitty":
		return termimg.Kitty
	case "iterm", "iterm2":
		return termimg.ITerm2
	case "sixel":
		return termimg.Sixel
	default:
		return termimg.Halfblocks
	}
}

// protocolName 统一协议名，便于缓存 key 和状态展示使用。
func protocolName(protocol tuiImageProtocol) string {
	if protocol == termimg.Unsupported {
		return "off"
	}
	return protocol.String()
}

// isOverlayImageProtocol 判断是否需要在主文本渲染完成后额外移动光标输出图像。
func isOverlayImageProtocol(protocol tuiImageProtocol) bool {
	return protocol == termimg.ITerm2 || protocol == termimg.Sixel
}

// isCoverOverlayProtocol 判断预览区是否使用独立图像层。
func isCoverOverlayProtocol(protocol tuiImageProtocol) bool {
	return isOverlayImageProtocol(protocol) || (protocol == termimg.Kitty && isGhosttyTerminal())
}

// previewImageFrame 描述预览区内封面框的位置和内部可绘制区域。
type previewImageFrame struct {
	x      int
	y      int
	w      int
	h      int
	innerX int
	innerY int
	innerW int
	innerH int
}

// syncCoverPreviewCmd 同步当前选中书籍与封面预览状态，必要时发起异步加载命令。
func (m *appModel) syncCoverPreviewCmd() tea.Cmd {
	rect := m.layout().cover
	frame, ok := previewImageFrameFor(rect)
	width, height := 0, 0
	if ok {
		width, height = frame.innerW, frame.innerH
	}
	item := m.currentItem()
	if item == nil || item.BookID == "" || !ok {
		m.coverPreview = coverPreviewState{}
		return nil
	}
	if m.coverProtocol == termimg.Unsupported {
		m.coverPreview = coverPreviewState{
			BookID:  item.BookID,
			Title:   item.Title,
			Width:   width,
			Height:  height,
			ErrText: locale.GetString("tui_cover_disabled"),
		}
		return nil
	}
	if m.coverCache == nil {
		m.coverCache = make(map[string]coverPreviewState)
	}

	key := coverPreviewCacheKey(item.BookID, width, height, m.coverProtocol)
	if cached, ok := m.coverCache[key]; ok {
		if m.coverPreview.Loading && m.coverPreview.BookID != item.BookID {
			// 当前选择已命中缓存时，作废上一个异步封面请求，避免旧结果回写到新书籍上。
			m.coverRequestID++
		}
		m.coverPreview = cached
		return nil
	}
	if m.coverPreview.Loading &&
		m.coverPreview.BookID == item.BookID &&
		m.coverPreview.Width == width &&
		m.coverPreview.Height == height &&
		m.coverPreview.Protocol == m.coverProtocol {
		// 同一本书同一尺寸已有加载任务时不重复发起，避免滚动时堆积解码任务。
		return nil
	}

	m.coverRequestID++
	requestID := m.coverRequestID
	m.coverPreview = coverPreviewState{
		BookID:   item.BookID,
		Title:    item.Title,
		Width:    width,
		Height:   height,
		Protocol: m.coverProtocol,
		Loading:  true,
	}
	return loadCoverPreviewCmd(requestID, item.BookID, item.Title, width, height, m.coverProtocol)
}

// applyCoverPreviewMsg 只接收最新请求的结果，丢弃快速切换书籍时返回较晚的旧结果。
func (m *appModel) applyCoverPreviewMsg(msg coverPreviewMsg) {
	if msg.requestID != m.coverRequestID {
		return
	}
	m.coverPreview = msg.state
	if msg.state.ErrText == "" && !msg.state.Loading {
		m.coverCache[coverPreviewCacheKey(msg.state.BookID, msg.state.Width, msg.state.Height, msg.state.Protocol)] = msg.state
	}
}

// loadCoverPreviewCmd 在 Bubble Tea 命令中异步读取封面并渲染成当前终端协议需要的内容。
func loadCoverPreviewCmd(requestID int, bookID string, title string, width int, height int, protocol tuiImageProtocol) tea.Cmd {
	return func() tea.Msg {
		state := coverPreviewState{
			BookID:   bookID,
			Title:    title,
			Width:    width,
			Height:   height,
			Protocol: protocol,
		}
		result, err := coverutil.GetBookCover(coverutil.Request{BookID: bookID, ResizeHeight: coverResizeHeight(height, protocol)})
		if err != nil {
			state.ErrText = err.Error()
			return coverPreviewMsg{requestID: requestID, state: state}
		}
		img, _, err := image.Decode(bytes.NewReader(result.Data))
		if err != nil {
			state.ErrText = err.Error()
			return coverPreviewMsg{requestID: requestID, state: state}
		}
		imageW, imageH := fitImageCellsForProtocol(img.Bounds(), width, height, protocol)
		state.ImageW = imageW
		state.ImageH = imageH
		if protocol == termimg.ITerm2 {
			state.Overlay = renderITerm2InlineImage(result.Data, imageW, imageH, width, height)
			debugTUIImageRender("cover", protocol, img.Bounds(), width, height, imageW, imageH, 0, 0)
			return coverPreviewMsg{requestID: requestID, state: state}
		}
		renderW, renderH := termImageRenderSizeForProtocol(imageW, imageH, protocol)
		termImage := termimg.New(img).
			Protocol(protocol).
			Size(renderW, renderH).
			Scale(termimg.ScaleFit)
		if protocol == termimg.Kitty && !isCoverOverlayProtocol(protocol) {
			// Kitty 使用 Unicode placeholder，这样图片区域能跟随文本布局，不需要额外覆盖定位。
			termImage = termImage.UseUnicode(true)
		}
		rendered, err := termImage.Render()
		if err != nil {
			state.ErrText = err.Error()
			return coverPreviewMsg{requestID: requestID, state: state}
		}

		if isCoverOverlayProtocol(protocol) {
			state.Overlay = rendered
		} else {
			state.Setup, state.Lines = splitRenderedImageLines(rendered, protocol)
		}
		debugTUIImageRender("cover", protocol, img.Bounds(), width, height, imageW, imageH, len(state.Lines), len(state.Setup))
		return coverPreviewMsg{requestID: requestID, state: state}
	}
}

// renderITerm2InlineImage 按 iTerm2 官方协议用字符格单位指定尺寸，避免 Retina 像素单位导致封面偏小。
func renderITerm2InlineImage(data []byte, widthCells int, heightCells int, maxWidthCells int, maxHeightCells int) string {
	widthCells = max(1, widthCells)
	heightCells = max(1, heightCells)
	widthArg, heightArg := iterm2InlineImageSizeArgs(widthCells, heightCells, maxWidthCells, maxHeightCells)
	params := fmt.Sprintf("inline=1;doNotMoveCursor=1;size=%d;width=%s;height=%s;preserveAspectRatio=1", len(data), widthArg, heightArg)
	return iterm2OSCStart() + "1337;File=" + params + ":" + base64.StdEncoding.EncodeToString(data) + iterm2OSCEnd()
}

// iterm2InlineImageSizeArgs 只固定真正贴边的一条轴，避免终端在固定宽高框内二次等比缩放产生底部留白。
func iterm2InlineImageSizeArgs(widthCells int, heightCells int, maxWidthCells int, maxHeightCells int) (string, string) {
	switch {
	case maxHeightCells > 0 && heightCells >= maxHeightCells && widthCells < maxWidthCells:
		return "auto", fmt.Sprintf("%d", heightCells)
	case maxWidthCells > 0 && widthCells >= maxWidthCells && heightCells < maxHeightCells:
		return fmt.Sprintf("%d", widthCells), "auto"
	default:
		return fmt.Sprintf("%d", widthCells), fmt.Sprintf("%d", heightCells)
	}
}

// iterm2OSCStart 在 tmux/screen 下包一层 DCS passthrough，让 iTerm2 能收到 OSC 1337。
func iterm2OSCStart() string {
	if strings.HasPrefix(os.Getenv("TERM"), "screen") || strings.HasPrefix(os.Getenv("TERM"), "tmux") {
		return "\x1bPtmux;\x1b\x1b]"
	}
	return "\x1b]"
}

// iterm2OSCEnd 与 iterm2OSCStart 配套结束当前 OSC 序列。
func iterm2OSCEnd() string {
	if strings.HasPrefix(os.Getenv("TERM"), "screen") || strings.HasPrefix(os.Getenv("TERM"), "tmux") {
		return "\a\x1b\\"
	}
	return "\a"
}

// coverPreviewCacheKey 按书籍、尺寸和协议生成缓存键；同一封面在不同终端协议下不能混用。
func coverPreviewCacheKey(bookID string, width int, height int, protocol tuiImageProtocol) string {
	return fmt.Sprintf("%s:%d:%d:%s", bookID, width, height, protocolName(protocol))
}

// coverResizeHeight 计算封面解码目标高度，减少传给终端渲染器的数据量。
func coverResizeHeight(height int, protocol tuiImageProtocol) int {
	resizeHeight := height * 32
	if protocol == termimg.Halfblocks {
		// ANSI halfblocks 在高分屏终端里观感偏小，单独提高输入像素，不影响原生图像协议。
		resizeHeight = height * 32 * halfblocksScalePercent / 100
	}
	return max(coverPreviewMinResizeHeight, resizeHeight)
}

// previewImageFrameFor 在预览面板中计算封面框。
// 除标题行和底部版本行外，预览区剩余空间都交给封面框使用。
func previewImageFrameFor(rect panelRect) (previewImageFrame, bool) {
	w := rect.innerWidth()
	h := rect.innerHeight()
	if w < coverPreviewMinWidth || h < previewTitleLines+previewFooterLines+coverPreviewMinHeight {
		return previewImageFrame{}, false
	}
	frameW := w
	frameH := h - previewTitleLines - previewFooterLines
	if frameW < coverPreviewMinWidth || frameH < coverPreviewMinHeight {
		return previewImageFrame{}, false
	}
	y := previewTitleLines
	return previewImageFrame{
		x:      0,
		y:      y,
		w:      frameW,
		h:      frameH,
		innerX: 1,
		innerY: y + 1,
		innerW: frameW - 2,
		innerH: frameH - 2,
	}, true
}

// fitImageCells 根据源图比例和终端字符格比例，计算图片在封面框中的显示尺寸。
func fitImageCells(bounds image.Rectangle, maxW int, maxH int) (int, int) {
	srcW := bounds.Dx()
	srcH := bounds.Dy()
	if srcW <= 0 || srcH <= 0 || maxW <= 0 || maxH <= 0 {
		return 1, 1
	}
	// 终端半块字符一格约等于两行像素，按这个比例让封面至少贴近图片框的一条边。
	ratioW := float64(maxW) / float64(srcW)
	ratioH := float64(maxH*2) / float64(srcH)
	ratio := math.Min(ratioW, ratioH)
	imageW := max(1, int(math.Round(float64(srcW)*ratio)))
	imageH := max(1, int(math.Round(float64(srcH)*ratio/2)))
	return min(imageW, maxW), min(imageH, maxH)
}

// fitImageCellsForProtocol 计算实际显示字符格；不同协议保留独立入口，方便按终端局部修正。
func fitImageCellsForProtocol(bounds image.Rectangle, maxW int, maxH int, protocol tuiImageProtocol) (int, int) {
	if protocol == termimg.Halfblocks {
		return fitImageCells(bounds, maxW, maxH)
	}
	cellW, cellH := protocolCellPixels(protocol)
	return fitImageCellsWithCellPixels(bounds, maxW, maxH, cellW, cellH)
}

// fitImageCellsWithCellPixels 按终端字符格像素比例计算原生图像协议的字符格尺寸。
func fitImageCellsWithCellPixels(bounds image.Rectangle, maxW int, maxH int, cellW int, cellH int) (int, int) {
	srcW := bounds.Dx()
	srcH := bounds.Dy()
	if srcW <= 0 || srcH <= 0 || maxW <= 0 || maxH <= 0 {
		return 1, 1
	}
	if cellW <= 0 || cellH <= 0 {
		cellW, cellH = 8, 16
	}
	ratioW := float64(maxW*cellW) / float64(srcW)
	ratioH := float64(maxH*cellH) / float64(srcH)
	ratio := math.Min(ratioW, ratioH)
	imageW := max(1, int(math.Round(float64(srcW)*ratio/float64(cellW))))
	imageH := max(1, int(math.Round(float64(srcH)*ratio/float64(cellH))))
	return min(imageW, maxW), min(imageH, maxH)
}

// protocolCellPixels 返回无需交互查询的终端字符格比例，避免 TUI 输入流被 CSI 查询污染。
func protocolCellPixels(protocol tuiImageProtocol) (int, int) {
	if isGhosttyTerminal() {
		return termimg.GhosttyWidth, termimg.GhosttyHeight
	}
	if isWezTermTerminal() {
		return termimg.WezTermWidth, termimg.WezTermHeight
	}
	if protocol == termimg.ITerm2 || isITerm2Terminal() {
		return termimg.ITermWidth, termimg.ITermHeight
	}
	if protocol == termimg.Kitty || isKittyTerminal() {
		return 8, 16
	}
	return 8, 16
}

// termImageRenderSizeForProtocol 计算传给 termimg 的渲染尺寸。
// Halfblocks 底层 mosaic 以 2x2 像素块输出一个字符，尺寸需要加倍，最终文字区域才会贴住封面框。
func termImageRenderSizeForProtocol(width int, height int, protocol tuiImageProtocol) (int, int) {
	if protocol == termimg.Halfblocks {
		return max(1, width*2), max(1, height*2)
	}
	return width, height
}

// splitRenderedImageLines 将 termimg 输出拆成控制序列和可见行。
// Kitty Unicode placeholder 的传输/placement 序列不能参与行宽计算，否则会被裁切成坏图。
func splitRenderedImageLines(rendered string, protocol tuiImageProtocol) (string, []string) {
	rendered = strings.TrimRight(rendered, "\n")
	if rendered == "" {
		return "", nil
	}
	setup := ""
	if protocol == termimg.Kitty {
		if idx := strings.Index(rendered, termimg.PLACEHOLDER_CHAR); idx > 0 {
			if setupEnd := strings.LastIndex(rendered[:idx], "\x1b\\"); setupEnd >= 0 {
				setupEnd += len("\x1b\\")
				setup = rendered[:setupEnd]
				rendered = rendered[setupEnd:]
			}
		}
	}
	rawLines := strings.Split(rendered, "\n")
	lines := make([]string, 0, len(rawLines))
	for _, line := range rawLines {
		lines = append(lines, line)
	}
	return setup, lines
}

// renderCoverPreviewContent 渲染预览区：首行显示标题，中间全部给封面框，底部显示版本和当前时间。
func (m *appModel) renderCoverPreviewContent(rect panelRect) []string {
	w := rect.innerWidth()
	h := rect.innerHeight()
	if h <= 0 {
		return nil
	}
	item := m.currentItem()
	if item == nil || item.BookID == "" {
		return appendPreviewFooter(m, []string{locale.GetString("tui_cover_no_selection")}, w, h)
	}

	frame, ok := previewImageFrameFor(rect)
	lines := []string{shortenText(item.Title, w)}
	if !ok {
		lines = append(lines, locale.GetString("tui_cover_too_small"))
		return appendPreviewFooter(m, lines, w, h)
	}

	for len(lines) < frame.y {
		lines = append(lines, "")
	}
	lines = append(lines, m.renderCoverImageFrame(frame, item)...)
	return appendPreviewFooter(m, lines, w, h)
}

// renderCoverImageFrame 渲染预览区内封面边框，图片内容在边框内部居中显示。
func (m *appModel) renderCoverImageFrame(frame previewImageFrame, item *shelfItem) []string {
	imageLines := m.renderCoverImageLines(frame, item)
	prefix := strings.Repeat(" ", frame.x)
	rows := make([]string, 0, frame.h)
	rows = append(rows, prefix+singleBorder().topPlain(frame.w))
	for _, line := range imageLines {
		rows = append(rows, prefix+singleBorder().middleStyled(line, frame.w))
	}
	rows = append(rows, prefix+singleBorder().bottom(frame.w))
	return rows
}

// renderCoverImageLines 根据封面状态输出加载中、错误文本、占位空白或 ANSI/Kitty 图像行。
func (m *appModel) renderCoverImageLines(frame previewImageFrame, item *shelfItem) []string {
	switch {
	case m.coverPreview.Loading || m.coverPreview.BookID != item.BookID:
		return centeredPreviewText(locale.GetString("tui_cover_loading"), frame.innerW, frame.innerH)
	case m.coverPreview.ErrText != "":
		return centeredPreviewText(shortenText(m.coverPreview.ErrText, frame.innerW), frame.innerW, frame.innerH)
	case isCoverOverlayProtocol(m.coverPreview.Protocol) && m.coverPreview.Overlay != "":
		return fitStyledLines(nil, frame.innerW, frame.innerH)
	default:
		return centerPreviewImageLines(m.coverPreview.Lines, frame.innerW, frame.innerH)
	}
}

// appendPreviewFooter 将版本号和当前时间固定放到预览区右下角。
func appendPreviewFooter(_ *appModel, lines []string, width int, height int) []string {
	if height <= 0 {
		return nil
	}
	if len(lines) >= height {
		lines = lines[:height]
	}
	for len(lines) < height-1 {
		lines = append(lines, "")
	}
	lines = append(lines, rightAlignStyled(previewVersionLine(), width))
	return fitStyledLines(lines, width, height)
}

// previewVersionLine 生成预览区底部状态行。
func previewVersionLine() string {
	return time.Now().Format("2006-01-02 15:04:05") + "  Comigo " + config.GetVersion()
}

func (m *appModel) renderCoverClearPrefix(rect panelRect) string {
	if _, ok := previewImageFrameFor(rect); !ok || m.coverProtocol != termimg.ITerm2 {
		return ""
	}
	// iTerm2 是独立图像层，切换封面时先清一帧屏幕，避免上一次 inline image 残留。
	return "\x1b[2J\x1b[H"
}

// renderCoverSetupPrefix 输出 Kitty 图片传输控制序列；可见 placeholder 行仍在主文本里渲染。
func (m *appModel) renderCoverSetupPrefix() string {
	item := m.currentItem()
	if item == nil || item.BookID != m.coverPreview.BookID || m.coverPreview.Protocol != termimg.Kitty {
		return ""
	}
	return m.coverPreview.Setup
}

// renderCoverOverlay 在主 TUI 文本之后绘制独立图像层，并把图像放到封面框内部中心。
func (m *appModel) renderCoverOverlay(rect panelRect) string {
	frame, ok := previewImageFrameFor(rect)
	if !ok || !isCoverOverlayProtocol(m.coverProtocol) {
		return ""
	}
	absX := rect.x + 2 + frame.innerX
	absY := rect.y + 2 + frame.innerY
	var builder strings.Builder
	builder.WriteString(xansi.SaveCursorPosition)
	if m.coverProtocol == termimg.Kitty {
		// Ghostty 预览改走 Kitty overlay，每帧先清理旧 placement，避免滚动选择时残留。
		builder.WriteString(termimg.ClearAllString())
	}
	if m.coverProtocol == termimg.ITerm2 {
		// iTerm2 inline image 可能残留上一帧背景，先用 ECH 清空目标单元格区域。
		builder.WriteString(clearITerm2CellArea(absX, absY, frame.innerW, frame.innerH))
	}
	builder.WriteString(clearTerminalArea(absX, absY, frame.innerW, frame.innerH))

	item := m.currentItem()
	if item != nil && item.BookID == m.coverPreview.BookID && m.coverPreview.Overlay != "" {
		imageW := min(max(1, m.coverPreview.ImageW), frame.innerW)
		imageH := min(max(1, m.coverPreview.ImageH), frame.innerH)
		col := absX + max(0, (frame.innerW-imageW)/2)
		row := absY + max(0, (frame.innerH-imageH)/2)
		builder.WriteString(xansi.CursorPosition(col, row))
		builder.WriteString(m.coverPreview.Overlay)
	}
	builder.WriteString(xansi.RestoreCursorPosition)
	return builder.String()
}

// clearITerm2CellArea 使用 ECH 清除指定终端单元格区域，只用于 iTerm2 残留清理。
func clearITerm2CellArea(col int, row int, width int, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	var builder strings.Builder
	for y := 0; y < height; y++ {
		builder.WriteString(xansi.CursorPosition(col, row+y))
		builder.WriteString(fmt.Sprintf("\x1b[%dX", width))
	}
	return builder.String()
}

// clearTerminalArea 用空格覆盖普通文本区域，给 overlay 图像留出干净背景。
func clearTerminalArea(col int, row int, width int, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	var builder strings.Builder
	blank := strings.Repeat(" ", width)
	for y := 0; y < height; y++ {
		builder.WriteString(xansi.CursorPosition(col, row+y))
		builder.WriteString(blank)
	}
	return builder.String()
}

// makeStyledPanel 与 makePanel 类似，但按 ANSI escape 感知宽度裁切，避免图像占位符被截断。
func (m *appModel) makeStyledPanel(title string, content []string, rect panelRect, focused bool) []string {
	if rect.w <= 1 || rect.h <= 1 {
		return []string{""}
	}
	border := singleBorder()
	if focused {
		border = doubleBorder()
	}

	lines := make([]string, 0, rect.h)
	lines = append(lines, border.top(title, rect.w))
	content = fitStyledLines(content, rect.innerWidth(), rect.innerHeight())
	for _, line := range content {
		lines = append(lines, border.middleStyled(line, rect.w))
	}
	for len(lines) < rect.h-1 {
		lines = append(lines, border.middleStyled("", rect.w))
	}
	lines = append(lines, border.bottom(rect.w))
	return lines
}

// middleStyled 渲染可包含 ANSI escape 的面板中间行。
func (b boxBorder) middleStyled(line string, width int) string {
	inner := max(0, width-2)
	return b.vertical + clipAndPadStyled(line, inner) + b.vertical
}

// topPlain 渲染不带标题的内部封面框顶部边框。
func (b boxBorder) topPlain(width int) string {
	inner := max(0, width-2)
	return b.leftTop + strings.Repeat(b.horizontal, inner) + b.rightTop
}

// fitStyledLines 按 ANSI 感知宽度裁切和补齐多行内容。
func fitStyledLines(lines []string, width int, height int) []string {
	if height <= 0 {
		return nil
	}
	result := make([]string, 0, height)
	for _, line := range lines {
		if len(result) >= height {
			break
		}
		result = append(result, clipAndPadStyled(line, width))
	}
	for len(result) < height {
		result = append(result, clipAndPadStyled("", width))
	}
	return result
}

// clipAndPadStyled 裁切并补齐可能含 ANSI escape 的字符串。
func clipAndPadStyled(text string, width int) string {
	if width <= 0 {
		return ""
	}
	text = strings.ReplaceAll(text, "\t", "    ")
	if xansi.StringWidth(text) > width {
		text = xansi.Truncate(text, width, "")
	}
	return text + strings.Repeat(" ", max(0, width-xansi.StringWidth(text)))
}

// centerPreviewImageLines 将已渲染图片行在封面框内居中；超高时从中心裁切，避免只显示图片顶部。
func centerPreviewImageLines(lines []string, width int, height int) []string {
	if height <= 0 {
		return nil
	}
	if len(lines) > height {
		start := (len(lines) - height) / 2
		lines = lines[start : start+height]
	}
	topPad := max(0, (height-len(lines))/2)
	result := make([]string, 0, height)
	for len(result) < topPad {
		result = append(result, clipAndPadStyled("", width))
	}
	for _, line := range lines {
		if len(result) >= height {
			break
		}
		result = append(result, centerStyledLine(line, width))
	}
	for len(result) < height {
		result = append(result, clipAndPadStyled("", width))
	}
	return result
}

// centeredPreviewText 在封面框内居中显示加载或错误文本。
func centeredPreviewText(text string, width int, height int) []string {
	if height <= 0 {
		return nil
	}
	lineIndex := height / 2
	lines := make([]string, 0, height)
	for i := 0; i < height; i++ {
		if i == lineIndex {
			lines = append(lines, centerStyledLine(text, width))
		} else {
			lines = append(lines, clipAndPadStyled("", width))
		}
	}
	return lines
}

// centerStyledLine 将可能带 ANSI escape 的单行内容居中。
func centerStyledLine(text string, width int) string {
	if width <= 0 {
		return ""
	}
	text = xansi.Truncate(text, width, "")
	pad := max(0, (width-xansi.StringWidth(text))/2)
	return strings.Repeat(" ", pad) + clipAndPadStyled(text, width-pad)
}

// rightAlignStyled 将可能带 ANSI escape 的单行内容右对齐。
func rightAlignStyled(text string, width int) string {
	if width <= 0 {
		return ""
	}
	text = xansi.Truncate(text, width, "")
	return strings.Repeat(" ", max(0, width-xansi.StringWidth(text))) + text
}

// debugTUIImageRender 按需记录终端图片渲染诊断信息，默认关闭，避免污染普通 TUI 日志。
func debugTUIImageRender(scope string, protocol tuiImageProtocol, bounds image.Rectangle, areaW int, areaH int, imageW int, imageH int, lineCount int, setupBytes int) {
	if os.Getenv("COMIGO_TUI_IMAGE_DEBUG") == "" {
		return
	}
	logger.Infof(
		"TUI image debug scope=%s protocol=%s src=%dx%d area=%dx%d image=%dx%d lines=%d setupBytes=%d",
		scope,
		protocolName(protocol),
		bounds.Dx(),
		bounds.Dy(),
		areaW,
		areaH,
		imageW,
		imageH,
		lineCount,
		setupBytes,
	)
}
