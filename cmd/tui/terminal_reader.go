package tui

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"strings"
	"time"

	termimg "github.com/blacktop/go-termimg"
	tea "github.com/charmbracelet/bubbletea"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	modelpkg "github.com/yumenaka/comigo/model"
	fileutil "github.com/yumenaka/comigo/tools/file"
)

const (
	defaultReaderAutoInterval = 10   // 自动翻页默认 10 秒，避免刚启用时翻页过快。
	minReaderAutoInterval     = 1    // 自动翻页最短 1 秒，保留快速浏览能力。
	maxReaderAutoInterval     = 120  // 自动翻页最长 120 秒，避免误操作调到过大值。
	readerMinResizeHeight     = 1600 // 阅读页取图最小高度，保证终端大图模式下有足够清晰度。
	readerResizePixelsPerCell = 48   // 每个终端字符格按约 48 像素估算取图高度，兼顾清晰度和解码成本。
	readerTitleBackPrefix     = "◀ "
)

// terminalReaderState 保存 TUI 终端阅读页的渲染结果。
// Lines 用于 ANSI/Kitty 占位符路径，Overlay 用于 iTerm2/Sixel 这类独立图像层。
type terminalReaderState struct {
	BookID    string
	Title     string
	PageIndex int
	PageCount int
	Width     int
	Height    int
	ImageW    int
	ImageH    int
	Protocol  tuiImageProtocol
	Setup     string
	Lines     []string
	Overlay   string
	Loading   bool
	ErrText   string
}

type terminalReaderPageMsg struct {
	requestID int
	state     terminalReaderState
}

type terminalReaderImageArea struct {
	x int
	y int
	w int
	h int
}

// detectTUIReaderImageProtocol 是终端阅读专用协议选择。
// WezTerm 避开已验证不稳定的 Kitty 路径，改走当前问题最少的 iTerm2 inline image 协议。
func detectTUIReaderImageProtocol() tuiImageProtocol {
	return detectTUIReaderImageProtocolWithKittyAuto(false)
}

// detectNativeTUIReaderImageProtocol 用于手动切换图片模式；自动启动仍让 Kitty 系终端默认走 ANSI。
func detectNativeTUIReaderImageProtocol() tuiImageProtocol {
	return detectTUIReaderImageProtocolWithKittyAuto(true)
}

func detectTUIReaderImageProtocolWithKittyAuto(allowKittyAuto bool) tuiImageProtocol {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("COMIGO_TUI_IMAGE"))) {
	case "", "auto":
		if isITerm2Terminal() {
			return termimg.ITerm2
		}
		if isWezTermTerminal() {
			return termimg.ITerm2
		}
		if isKittyTerminal() || isGhosttyTerminal() {
			// Kitty 系终端默认走 ANSI；用户按 c 切换图片模式时再启用 Kitty 原生路径。
			if allowKittyAuto {
				return termimg.Kitty
			}
			return termimg.Halfblocks
		}
		protocol := termimg.DetectProtocol()
		if protocol == termimg.ITerm2 || (allowKittyAuto && protocol == termimg.Kitty) {
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

func isITerm2Terminal() bool {
	return os.Getenv("TERM_PROGRAM") == "iTerm.app" ||
		strings.Contains(strings.ToLower(os.Getenv("LC_TERMINAL")), "iterm") ||
		os.Getenv("ITERM_SESSION_ID") != ""
}

func isGhosttyTerminal() bool {
	return strings.Contains(strings.ToLower(os.Getenv("TERM")), "ghostty") ||
		strings.Contains(strings.ToLower(os.Getenv("TERM_PROGRAM")), "ghostty") ||
		os.Getenv("GHOSTTY_RESOURCES_DIR") != ""
}

func isKittyTerminal() bool {
	return os.Getenv("KITTY_WINDOW_ID") != "" ||
		strings.Contains(strings.ToLower(os.Getenv("TERM")), "kitty") ||
		strings.Contains(strings.ToLower(os.Getenv("TERM_PROGRAM")), "kitty")
}

func isWezTermTerminal() bool {
	return strings.Contains(strings.ToLower(os.Getenv("TERM")), "wezterm") ||
		strings.Contains(strings.ToLower(os.Getenv("TERM_PROGRAM")), "wezterm") ||
		os.Getenv("WEZTERM_EXECUTABLE") != "" ||
		os.Getenv("WEZTERM_PANE") != ""
}

// syncActiveImageCmd 根据当前界面同步需要显示的图片：书架同步封面，阅读页同步当前页。
func (m *appModel) syncActiveImageCmd() tea.Cmd {
	if m.screen == screenReader {
		return m.syncTerminalReaderPageCmd()
	}
	return m.syncCoverPreviewCmd()
}

// startTerminalReader 从当前选中书籍的第一页进入终端阅读。
func (m *appModel) startTerminalReader(item *shelfItem) tea.Cmd {
	if item == nil || item.BookID == "" {
		m.setActionMsg(locale.GetString("tui_terminal_reader_no_book"))
		m.refreshStatus()
		return nil
	}
	if modelpkg.IStore == nil {
		m.setActionMsg(locale.GetString("tui_shelf_not_initialized"))
		m.refreshStatus()
		return nil
	}
	book, err := modelpkg.IStore.GetBook(item.BookID)
	if err != nil {
		m.setActionMsg(shortenText(err.Error(), maxActionMessage))
		m.refreshStatus()
		return nil
	}
	if len(book.PageInfos) == 0 {
		m.setActionMsg(locale.GetString("tui_terminal_reader_no_pages"))
		m.refreshStatus()
		return nil
	}

	if m.coverProtocol == termimg.Kitty && isCoverOverlayProtocol(m.coverProtocol) {
		// Ghostty 预览区使用 Kitty overlay，进入阅读页前先清掉书架预览的旧图层。
		m.clearKittyImagesNextFrame = true
	}
	m.screen = screenReader
	m.readerAutoFlip = false
	m.readerNextAutoAt = time.Time{}
	m.readerPendingPage = false
	m.readerPendingRequestKey = ""
	if m.readerAutoInterval <= 0 {
		m.readerAutoInterval = defaultReaderAutoInterval
	}
	m.terminalReader = terminalReaderState{
		BookID:    book.BookID,
		Title:     item.Title,
		PageIndex: 0,
		PageCount: len(book.PageInfos),
		Protocol:  m.readerProtocol,
		Loading:   true,
	}
	return m.syncTerminalReaderPageCmd()
}

// exitTerminalReader 返回书架界面，并重新同步书架封面预览。
func (m *appModel) exitTerminalReader() tea.Cmd {
	m.screen = screenShelf
	m.readerAutoFlip = false
	m.readerNextAutoAt = time.Time{}
	m.readerPendingPage = false
	m.readerPendingRequestKey = ""
	m.readerRequestID++
	return m.syncCoverPreviewCmd()
}

// handleTerminalReaderKey 处理终端阅读快捷键。
func (m *appModel) handleTerminalReaderKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "delete":
		return m, m.exitTerminalReader()
	case "up", "left", "pgup":
		m.moveTerminalReaderPage(-1)
	case "down", "right", "pgdown", "enter", " ":
		m.moveTerminalReaderPage(1)
	case "f":
		m.terminalReaderFullscreen = !m.terminalReaderFullscreen
	case "a":
		m.toggleTerminalReaderAutoFlip()
	case "+", "=":
		m.adjustTerminalReaderAutoInterval(1)
	case "-":
		m.adjustTerminalReaderAutoInterval(-1)
	}
	return m, m.syncTerminalReaderPageCmd()
}

// handleTerminalReaderMouse 处理终端阅读鼠标操作：标题返回书架，图片区域左右半屏翻页。
func (m *appModel) handleTerminalReaderMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Action != tea.MouseActionPress || msg.Button != tea.MouseButtonLeft {
		return m, nil
	}
	if !m.terminalReaderFullscreen && msg.Y == 0 && msg.X < m.terminalReaderTitleHitWidth(m.renderWidth()) {
		return m, m.exitTerminalReader()
	}
	area := m.terminalReaderImageArea()
	if msg.Y < area.y || msg.Y >= area.y+area.h || msg.X < area.x || msg.X >= area.x+area.w {
		return m, nil
	}
	if msg.X < area.x+area.w/2 {
		m.moveTerminalReaderPage(-1)
	} else {
		m.moveTerminalReaderPage(1)
	}
	return m, m.syncTerminalReaderPageCmd()
}

// moveTerminalReaderPage 翻到相邻页，越界时停在当前页并提示。
func (m *appModel) moveTerminalReaderPage(delta int) bool {
	next := m.terminalReaderTargetPageIndex() + delta
	if next < 0 {
		m.setActionMsg(locale.GetString("hint_first_page"))
		return false
	}
	if next >= m.terminalReader.PageCount {
		m.setActionMsg(locale.GetString("hint_last_page"))
		return false
	}
	if next == m.terminalReader.PageIndex {
		if m.readerPendingPage {
			// 回到当前可见页时取消后台请求，继续保留当前画面。
			m.readerPendingPage = false
			m.readerPendingRequestKey = ""
			m.readerRequestID++
			return true
		}
		return false
	}
	// 翻页时先只记录目标页，旧图和旧页码继续显示，等新图 ready 后一次性替换。
	m.readerRequestID++
	m.readerPendingPage = true
	m.readerPendingPageIndex = next
	m.readerPendingRequestKey = ""
	return true
}

// toggleTerminalReaderAutoFlip 开关自动翻页；开启后按当前间隔从现在重新计时。
func (m *appModel) toggleTerminalReaderAutoFlip() {
	m.readerAutoFlip = !m.readerAutoFlip
	if m.readerAutoFlip {
		m.scheduleNextReaderAutoFlip(time.Now())
		m.setActionMsg(fmt.Sprintf(locale.GetString("tui_terminal_reader_auto_started"), m.readerAutoInterval))
		return
	}
	m.readerNextAutoAt = time.Time{}
	m.setActionMsg(locale.GetString("tui_terminal_reader_auto_stopped"))
}

// adjustTerminalReaderAutoInterval 调整自动翻页秒数，开启状态下立即重置下一次触发时间。
func (m *appModel) adjustTerminalReaderAutoInterval(delta int) {
	if m.readerAutoInterval <= 0 {
		m.readerAutoInterval = defaultReaderAutoInterval
	}
	m.readerAutoInterval += delta
	if m.readerAutoInterval < minReaderAutoInterval {
		m.readerAutoInterval = minReaderAutoInterval
	}
	if m.readerAutoInterval > maxReaderAutoInterval {
		m.readerAutoInterval = maxReaderAutoInterval
	}
	if m.readerAutoFlip {
		m.scheduleNextReaderAutoFlip(time.Now())
	}
	m.setActionMsg(fmt.Sprintf(locale.GetString("tui_terminal_reader_auto_interval"), m.readerAutoInterval))
}

func (m *appModel) scheduleNextReaderAutoFlip(now time.Time) {
	m.readerNextAutoAt = now.Add(time.Duration(max(minReaderAutoInterval, m.readerAutoInterval)) * time.Second)
}

// handleReaderAutoFlip 在已有全局 tick 上驱动自动翻页，避免额外起定时器。
func (m *appModel) handleReaderAutoFlip(now time.Time) tea.Cmd {
	if m.screen != screenReader || !m.readerAutoFlip || m.readerNextAutoAt.IsZero() || now.Before(m.readerNextAutoAt) {
		return nil
	}
	if m.readerPendingPage {
		return nil
	}
	if !m.moveTerminalReaderPage(1) {
		m.readerAutoFlip = false
		m.readerNextAutoAt = time.Time{}
		m.setActionMsg(locale.GetString("tui_terminal_reader_auto_reached_end"))
		return nil
	}
	m.scheduleNextReaderAutoFlip(now)
	return m.syncTerminalReaderPageCmd()
}

// syncTerminalReaderPageCmd 同步当前页、尺寸和终端协议，必要时发起异步渲染。
func (m *appModel) syncTerminalReaderPageCmd() tea.Cmd {
	area := m.terminalReaderImageArea()
	if m.screen != screenReader || m.terminalReader.BookID == "" || area.w <= 0 || area.h <= 0 {
		return nil
	}
	if m.readerProtocol == termimg.Unsupported {
		m.readerPendingPage = false
		m.readerPendingRequestKey = ""
		m.terminalReader.ErrText = locale.GetString("tui_cover_disabled")
		m.terminalReader.Loading = false
		return nil
	}
	if m.terminalReaderCache == nil {
		m.terminalReaderCache = make(map[string]terminalReaderState)
	}

	targetPage := m.terminalReaderTargetPageIndex()
	key := terminalReaderCacheKey(m.terminalReader.BookID, targetPage, area.w, area.h, m.readerProtocol)
	if cached, ok := m.terminalReaderCache[key]; ok {
		wasPending := m.readerPendingPage
		m.readerPendingPage = false
		m.readerPendingRequestKey = ""
		m.terminalReader = cached
		if wasPending && m.readerAutoFlip {
			m.scheduleNextReaderAutoFlip(time.Now())
		}
		return nil
	}
	if m.readerPendingPage && m.readerPendingRequestKey == key {
		return nil
	}
	if m.terminalReader.Loading &&
		m.terminalReader.Width == area.w &&
		m.terminalReader.Height == area.h &&
		m.terminalReader.Protocol == m.readerProtocol {
		return nil
	}

	m.readerRequestID++
	requestID := m.readerRequestID
	state := m.terminalReader
	state.PageIndex = targetPage
	state.Width = area.w
	state.Height = area.h
	state.Protocol = m.readerProtocol
	state.Loading = true
	state.ErrText = ""
	state.Setup = ""
	state.Lines = nil
	state.Overlay = ""
	if m.readerPendingPage {
		m.readerPendingRequestKey = key
		return loadTerminalReaderPageCmd(requestID, state, m.readerProtocol)
	}
	m.terminalReader = state
	return loadTerminalReaderPageCmd(requestID, state, m.readerProtocol)
}

// applyTerminalReaderPageMsg 应用异步渲染结果；只接收最新请求，避免快速翻页时旧页面回写。
func (m *appModel) applyTerminalReaderPageMsg(msg terminalReaderPageMsg) {
	if msg.requestID != m.readerRequestID {
		return
	}
	wasPending := m.readerPendingPage
	m.readerPendingPage = false
	m.readerPendingRequestKey = ""
	m.terminalReader = msg.state
	if msg.state.ErrText == "" && !msg.state.Loading {
		m.terminalReaderCache[terminalReaderCacheKey(msg.state.BookID, msg.state.PageIndex, msg.state.Width, msg.state.Height, msg.state.Protocol)] = msg.state
	}
	if wasPending && m.readerAutoFlip {
		m.scheduleNextReaderAutoFlip(time.Now())
	}
}

func loadTerminalReaderPageCmd(requestID int, state terminalReaderState, protocol tuiImageProtocol) tea.Cmd {
	return func() tea.Msg {
		result := state
		result.Protocol = protocol
		result.Loading = false

		if modelpkg.IStore == nil {
			result.ErrText = locale.GetString("tui_shelf_not_initialized")
			return terminalReaderPageMsg{requestID: requestID, state: result}
		}
		book, err := modelpkg.IStore.GetBook(state.BookID)
		if err != nil {
			result.ErrText = err.Error()
			return terminalReaderPageMsg{requestID: requestID, state: result}
		}
		result.PageCount = len(book.PageInfos)
		if state.PageIndex < 0 || state.PageIndex >= len(book.PageInfos) {
			result.ErrText = locale.GetString("tui_terminal_reader_page_missing")
			return terminalReaderPageMsg{requestID: requestID, state: result}
		}

		page := book.PageInfos[state.PageIndex]
		imgData, _, err := fileutil.GetPictureData(buildTerminalReaderPictureOption(book, page.Name, terminalReaderResizeHeight(state.Height, protocol)))
		if err != nil {
			result.ErrText = err.Error()
			return terminalReaderPageMsg{requestID: requestID, state: result}
		}
		img, _, err := image.Decode(bytes.NewReader(imgData))
		if err != nil {
			result.ErrText = err.Error()
			return terminalReaderPageMsg{requestID: requestID, state: result}
		}
		renderTerminalReaderImage(&result, imgData, img, protocol)
		return terminalReaderPageMsg{requestID: requestID, state: result}
	}
}

func buildTerminalReaderPictureOption(book *modelpkg.Book, pageName string, resizeHeight int) fileutil.GetPictureDataOption {
	return fileutil.GetPictureDataOption{
		PictureName:      pageName,
		BookID:           book.BookID,
		BookIsPDF:        book.Type == modelpkg.TypePDF,
		BookIsDir:        book.Type == modelpkg.TypeDir,
		BookIsNonUTF8Zip: book.NonUTF8Zip,
		BookPath:         book.BookPath,
		Debug:            config.GetCfg().Debug,
		UseCache:         config.GetCfg().UseCache,
		ResizeHeight:     resizeHeight,
		IsRemote:         book.IsRemote,
		RemoteURL:        book.RemoteURL,
	}
}

// terminalReaderResizeHeight 根据终端可绘制高度推算取图高度，避免源图过小导致原生图片模式发虚。
func terminalReaderResizeHeight(height int, protocol tuiImageProtocol) int {
	resizeHeight := height * readerResizePixelsPerCell
	if protocol == termimg.Halfblocks {
		resizeHeight = height * readerResizePixelsPerCell * halfblocksScalePercent / 100
	}
	return max(readerMinResizeHeight, resizeHeight)
}

// renderTerminalReaderImage 将图片渲染成当前终端协议需要的文本行或 overlay 序列。
func renderTerminalReaderImage(state *terminalReaderState, data []byte, img image.Image, protocol tuiImageProtocol) {
	imageW, imageH := fitImageCellsForProtocol(img.Bounds(), state.Width, state.Height, protocol)
	state.ImageW = imageW
	state.ImageH = imageH
	if protocol == termimg.ITerm2 {
		state.Overlay = renderITerm2InlineImage(data, imageW, imageH, state.Width, state.Height)
		debugTUIImageRender("reader", protocol, img.Bounds(), state.Width, state.Height, imageW, imageH, 0, 0)
		return
	}
	renderW, renderH := termImageRenderSizeForProtocol(imageW, imageH, protocol)
	termImage := termimg.New(img).
		Protocol(protocol).
		Size(renderW, renderH).
		Scale(termimg.ScaleFit)
	if protocol == termimg.Kitty {
		// Kitty 系列终端用 Unicode placeholder，让图片跟随 Bubble Tea 文本行一起刷新，避免图层和页码不同步。
		termImage = termImage.UseUnicode(true)
		rendered, err := termImage.Render()
		if err != nil {
			state.ErrText = err.Error()
			return
		}
		state.Setup, state.Lines = splitRenderedImageLines(rendered, protocol)
		debugTUIImageRender("reader", protocol, img.Bounds(), state.Width, state.Height, imageW, imageH, len(state.Lines), len(state.Setup))
		return
	}
	rendered, err := termImage.Render()
	if err != nil {
		state.ErrText = err.Error()
		return
	}
	if isOverlayImageProtocol(protocol) {
		state.Overlay = rendered
		debugTUIImageRender("reader", protocol, img.Bounds(), state.Width, state.Height, imageW, imageH, 0, 0)
		return
	}
	state.Setup, state.Lines = splitRenderedImageLines(rendered, protocol)
	debugTUIImageRender("reader", protocol, img.Bounds(), state.Width, state.Height, imageW, imageH, len(state.Lines), len(state.Setup))
}

func terminalReaderCacheKey(bookID string, pageIndex int, width int, height int, protocol tuiImageProtocol) string {
	return fmt.Sprintf("%s:%d:%d:%d:%s", bookID, pageIndex, width, height, protocolName(protocol))
}

// terminalReaderTargetPageIndex 返回当前应该后台加载的页；没有 pending 时就是可见页。
func (m *appModel) terminalReaderTargetPageIndex() int {
	if m.readerPendingPage {
		return m.readerPendingPageIndex
	}
	return m.terminalReader.PageIndex
}

// terminalReaderImageArea 计算阅读页图片可用区域；普通模式保留顶部标题和底部状态各 1 行。
func (m *appModel) terminalReaderImageArea() terminalReaderImageArea {
	width := m.renderWidth()
	height := max(0, m.height)
	if m.terminalReaderFullscreen {
		return terminalReaderImageArea{w: width, h: height}
	}
	return terminalReaderImageArea{y: 1, w: width, h: max(0, height-2)}
}

func (m *appModel) renderTerminalReaderView() string {
	width := m.renderWidth()
	if width <= 0 || m.height <= 0 {
		return ""
	}
	area := m.terminalReaderImageArea()
	lines := make([]string, 0, m.height)
	if !m.terminalReaderFullscreen {
		lines = append(lines, m.renderTerminalReaderTitle(width))
	}
	lines = append(lines, m.renderTerminalReaderImageLines(area)...)
	if !m.terminalReaderFullscreen {
		lines = append(lines, m.renderTerminalReaderFooter(width))
	}
	lines = fitStyledLines(lines, width, m.height)
	return m.renderTerminalReaderClearPrefix() + m.renderTerminalReaderSetupPrefix() + strings.Join(lines, "\n") + m.renderTerminalReaderOverlay(area)
}

// renderTerminalReaderImageLines 渲染图片区域文本；overlay 协议只占位，实际图片由 renderTerminalReaderOverlay 输出。
func (m *appModel) renderTerminalReaderImageLines(area terminalReaderImageArea) []string {
	switch {
	case area.h <= 0:
		return nil
	case m.terminalReader.Loading:
		return fitStyledLines(nil, area.w, area.h)
	case m.terminalReader.ErrText != "":
		return centeredPreviewText(shortenText(m.terminalReader.ErrText, area.w), area.w, area.h)
	case isTerminalReaderOverlayProtocol(m.terminalReader.Protocol) && m.terminalReader.Overlay != "":
		return fitStyledLines(nil, area.w, area.h)
	default:
		return centerPreviewImageLines(m.terminalReader.Lines, area.w, area.h)
	}
}

// renderTerminalReaderTitle 在标题栏左侧显示书名，加载中时把状态提示放到整屏右上角。
func (m *appModel) renderTerminalReaderTitle(width int) string {
	title := m.renderTerminalReaderTitleText(width)
	if !m.terminalReaderIsLoading() {
		return clipAndPad(title, width)
	}
	loading := shortenText(locale.GetString("tui_terminal_reader_loading"), width)
	loadingWidth := runewidth.StringWidth(loading)
	title = shortenText(title, max(0, width-loadingWidth-1))
	line := title + strings.Repeat(" ", max(1, width-loadingWidth-runewidth.StringWidth(title))) + loading
	return clipAndPad(line, width)
}

// renderTerminalReaderTitleText 给标题加返回提示；点击这个范围会回到书架。
func (m *appModel) renderTerminalReaderTitleText(width int) string {
	if width <= 0 {
		return ""
	}
	return shortenText(readerTitleBackPrefix+m.terminalReader.Title, width)
}

func (m *appModel) terminalReaderIsLoading() bool {
	return m.terminalReader.Loading || m.readerPendingPage
}

// terminalReaderTitleHitWidth 返回标题栏可点击范围宽度，限制到实际显示的返回提示与书名。
func (m *appModel) terminalReaderTitleHitWidth(width int) int {
	return runewidth.StringWidth(m.renderTerminalReaderTitleText(width))
}

// topRightPreviewText 将短状态提示放在右上角，供无标题栏场景兜底使用。
func topRightPreviewText(text string, width int, height int) []string {
	if height <= 0 {
		return nil
	}
	lines := []string{rightAlignStyled(shortenText(text, width), width)}
	for len(lines) < height {
		lines = append(lines, clipAndPadStyled("", width))
	}
	return lines
}

func (m *appModel) renderTerminalReaderFooter(width int) string {
	left := locale.GetString("tui_terminal_reader_hint")
	if m.readerAutoFlip {
		left = fmt.Sprintf(locale.GetString("tui_terminal_reader_auto_hint"), m.readerAutoInterval)
	}
	center := fmt.Sprintf("%d/%d", min(m.terminalReader.PageIndex+1, max(1, m.terminalReader.PageCount)), m.terminalReader.PageCount)
	right := terminalReaderVersionLine()
	return formatThreePartStatusLine(left, center, right, width)
}

// terminalReaderVersionLine 只显示版本号；阅读页含图片层，时间变化会触发不必要的重绘。
func terminalReaderVersionLine() string {
	return "Comigo " + config.GetVersion()
}

// formatThreePartStatusLine 将快捷键、页码、版本号分别放在左中右，空间不足时优先截断左侧提示。
func formatThreePartStatusLine(left string, center string, right string, width int) string {
	if width <= 0 {
		return ""
	}
	centerWidth := runewidth.StringWidth(center)
	centerStart := max(0, (width-centerWidth)/2)
	left = shortenText(left, max(0, centerStart-1))
	rightStartMin := min(width, centerStart+centerWidth+1)
	right = shortenText(right, max(0, width-rightStartMin))

	line := left
	line += strings.Repeat(" ", max(1, centerStart-runewidth.StringWidth(line))) + center
	rightWidth := runewidth.StringWidth(right)
	line += strings.Repeat(" ", max(1, width-rightWidth-runewidth.StringWidth(line))) + right
	return clipAndPad(line, width)
}

// renderTerminalReaderClearPrefix 输出阅读页绘制前的清理序列；只在必要终端/协议上清旧图层。
func (m *appModel) renderTerminalReaderClearPrefix() string {
	var builder strings.Builder
	if m.clearKittyImagesNextFrame {
		builder.WriteString(termimg.ClearAllString())
		m.clearKittyImagesNextFrame = false
	}
	switch m.terminalReader.Protocol {
	case termimg.ITerm2:
		if m.terminalReader.Loading && m.terminalReader.Overlay == "" {
			// 初次进入或尺寸切换且没有旧图可保留时才整屏清理；普通翻页保留旧图，避免 iTerm2 先黑屏。
			builder.WriteString("\x1b[2J\x1b[H")
		}
	default:
	}
	return builder.String()
}

// renderTerminalReaderSetupPrefix 输出 Kitty 图片传输控制序列；可见 placeholder 行仍随正文布局。
func (m *appModel) renderTerminalReaderSetupPrefix() string {
	if m.terminalReader.Protocol != termimg.Kitty {
		return ""
	}
	return m.terminalReader.Setup
}

// isTerminalReaderOverlayProtocol 判断终端阅读页是否使用独立图像层。
func isTerminalReaderOverlayProtocol(protocol tuiImageProtocol) bool {
	return isOverlayImageProtocol(protocol)
}

func (m *appModel) renderTerminalReaderOverlay(area terminalReaderImageArea) string {
	if area.w <= 0 || area.h <= 0 || !isTerminalReaderOverlayProtocol(m.terminalReader.Protocol) || m.terminalReader.Overlay == "" {
		return ""
	}
	if m.readerPendingPage {
		return ""
	}
	imageW := min(max(1, m.terminalReader.ImageW), area.w)
	imageH := min(max(1, m.terminalReader.ImageH), area.h)
	col := 1 + area.x + max(0, (area.w-imageW)/2)
	row := 1 + area.y + max(0, (area.h-imageH)/2)
	clearCol := 1 + area.x
	clearRow := 1 + area.y

	var builder strings.Builder
	builder.WriteString(xansi.SaveCursorPosition)
	if m.terminalReader.Protocol == termimg.ITerm2 {
		builder.WriteString(clearITerm2CellArea(clearCol, clearRow, area.w, area.h))
	}
	builder.WriteString(clearTerminalArea(clearCol, clearRow, area.w, area.h))
	builder.WriteString(xansi.CursorPosition(col, row))
	builder.WriteString(m.terminalReader.Overlay)
	builder.WriteString(xansi.RestoreCursorPosition)
	return builder.String()
}
