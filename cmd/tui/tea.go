package tui

import (
	"fmt"
	"io"
	stdlog "log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
	"golang.org/x/term"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	modelpkg "github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

// TUI 全局常量
const (
	logBufferLimit   = 5000                   // 日志缓冲区最大行数
	tickInterval     = 500 * time.Millisecond // 定时刷新间隔
	minTUIWidth      = 40                     // TUI 最小终端宽度
	minTUIHeight     = 24                     // TUI 最小终端高度
	narrowThreshold  = 100                    // 窄屏/宽屏布局切换阈值（宽度小于此值为窄屏单栏模式）
	layoutGap        = 1                      // 面板之间的间距列数
	maxActionMessage = 120                    // 操作消息最大显示宽度
	doubleClickGap   = 500 * time.Millisecond // 终端没有原生双击事件，用短时间内同一行两次点击模拟
)

// focusPanel 表示当前聚焦的面板类型
type focusPanel int

const (
	focusShelf  focusPanel = iota // 书架面板
	focusLog                      // 日志面板
	focusQRCode                   // QRCode 面板
)

// appScreen 表示当前 TUI 主界面。
type appScreen int

const (
	screenShelf  appScreen = iota // 书架/日志/QRCode 总览界面
	screenReader                  // 终端阅读界面
)

// qrAction 表示 QRCode 面板按钮动作。
type qrAction int

const (
	qrActionScrollMode     qrAction = iota // 切换到卷轴阅读 URL
	qrActionFlipMode                       // 切换到翻页阅读 URL
	qrActionCopyURL                        // 复制当前 URL
	qrActionTerminalReader                 // 在 TUI 内打开终端阅读
	qrActionOpenBrowser                    // 调用系统浏览器打开当前 URL
)

// qrButtonHitbox 记录 QRCode 面板按钮的点击范围。
type qrButtonHitbox struct {
	action qrAction
	row    int
	start  int
	end    int
}

// shelfItemKind 书架列表条目类型
type shelfItemKind int

const (
	shelfItemHeader shelfItemKind = iota // 标题行（书库名称等，不可选中）
	shelfItemBack                        // 返回上一级
	shelfItemGroup                       // 子书架（书籍分组）
	shelfItemBook                        // 可阅读的书籍
)

// Bubble Tea 消息类型
type (
	tickMsg           time.Time // 定时刷新消息
	backendStartedMsg struct{}  // 后台服务启动成功
	backendErrorMsg   struct {  // 后台服务启动失败
		err error
	}
)

// openURLResultMsg 浏览器打开 URL 的结果回调
type openURLResultMsg struct {
	url string
	err error
}

// shelfItem 书架列表中的单个条目
type shelfItem struct {
	Kind       shelfItemKind // 条目类型
	Title      string        // 显示标题
	Subtitle   string        // 副标题（页数、子项数等）
	BookID     string        // 书籍/书架 ID
	TargetURL  string        // 点击后打开的 URL
	Selectable bool          // 是否可被选中
}

// shelfLevel 书架导航栈的一层，用于记录进入子书架的路径
type shelfLevel struct {
	BookID string // 当前层对应的书架 ID
	Title  string // 显示名称
}

// systemSnapshot 保存当前仍会渲染到 TUI 的状态文本。
type systemSnapshot struct {
	StatusText string
}

// panelRect 表示一个面板在终端中的矩形区域（含边框）
type panelRect struct {
	x int // 左上角列
	y int // 左上角行
	w int // 总宽度（含左右边框各 1 列）
	h int // 总高度（含上下边框各 1 行）
}

// innerWidth 返回面板内容区宽度（去掉左右边框）
func (r panelRect) innerWidth() int {
	if r.w <= 2 {
		return 0
	}
	return r.w - 2
}

// innerHeight 返回面板内容区高度（去掉上下边框）
func (r panelRect) innerHeight() int {
	if r.h <= 2 {
		return 0
	}
	return r.h - 2
}

// layoutState 四个面板的布局位置
type layoutState struct {
	shelf panelRect // 书架面板
	cover panelRect // 封面预览面板
	log   panelRect // 日志面板
	qr    panelRect // QRCode 面板
}

// appModel 是 Bubble Tea 的核心模型，持有全部 TUI 状态。
type appModel struct {
	logBuffer *LogBuffer // 日志缓冲区（由 logger 镜像写入）

	width  int        // 终端当前宽度
	height int        // 终端当前高度
	focus  focusPanel // 当前聚焦的面板
	screen appScreen  // 当前主界面

	logs           []string // 当前快照的日志行
	logOffset      int      // 日志滚动偏移（首行索引）
	autoFollowLogs bool     // 是否自动跟随最新日志

	backendReady bool       // 后台服务是否已启动完成
	backendError string     // 后台启动错误信息（空表示无错误）
	actionMsg    string     // 最近一次操作的状态提示
	modal        modalState // 全局提示弹窗，优先接管键盘和鼠标事件

	stack               []shelfLevel // 书架导航栈（从根到当前层级）
	items               []shelfItem  // 当前层级的书架条目列表
	selected            int          // 当前选中的条目索引
	shelfOffset         int          // 书架列表滚动偏移（首个可见的 item 索引）
	shelfRowToID        map[int]int  // 面板内容行号 → items 索引的映射（用于鼠标点击定位）
	lastShelfClickIndex int          // 上一次书架左键点击的条目索引，用于判断双击
	lastShelfClickAt    time.Time    // 上一次书架左键点击时间，用于判断双击

	currentShelfURL string   // 当前书架层级对应的 Web URL
	qrLines         []string // QR 码的 Unicode 字符行
	qrButtonFocus   qrAction // QR 面板当前聚焦按钮
	qrButtonHitbox  []qrButtonHitbox
	readMode        int            // 阅读模式：0=scroll（卷轴阅读）, 1=flip（翻页阅读）
	status          systemSnapshot // 最新的系统状态快照

	coverProtocol  tuiImageProtocol             // 终端图片协议，自动检测后可用 COMIGO_TUI_IMAGE 覆盖
	coverPreview   coverPreviewState            // 当前选中书籍的封面预览状态
	coverCache     map[string]coverPreviewState // 已渲染封面缓存，避免频繁滚动时重复解码
	coverRequestID int                          // 封面异步加载序号，用于丢弃过期结果
	coverSetupKey  string                       // 已发送到 Kitty 的封面 setup key，避免每帧重复传图

	readerProtocol            tuiImageProtocol               // 终端阅读使用的图片协议，现代 Kitty 协议终端会尝试原生图像
	terminalReader            terminalReaderState            // 终端阅读状态
	terminalReaderCache       map[string]terminalReaderState // 终端阅读页缓存
	readerRequestID           int                            // 终端阅读异步加载序号
	readerSetupKey            string                         // 已发送到 Kitty 的阅读页 setup key，避免每帧重复传图
	readerAutoFlip            bool                           // 是否正在自动翻页
	readerAutoInterval        int                            // 自动翻页间隔秒数
	readerNextAutoAt          time.Time                      // 下一次自动翻页时间
	terminalReaderFullscreen  bool                           // 终端阅读是否隐藏顶部和底部状态栏
	readerPendingPage         bool                           // 是否有后台加载中的目标页，加载期间继续显示当前页
	readerPendingPageIndex    int                            // 后台加载中的目标页索引
	readerPendingRequestKey   string                         // 后台加载请求 key，避免同一目标页重复发起渲染
	clearKittyImagesNextFrame bool                           // 切换界面时下一帧清理 Kitty 图像层，避免旧 overlay 残留
}

// LogBuffer 用来缓存 TUI 需要展示的实时日志。
type LogBuffer struct {
	lines []string
	limit int
	mu    sync.RWMutex
}

// NewLogBuffer 创建一个新的日志缓冲区实例。
func NewLogBuffer() *LogBuffer {
	return &LogBuffer{limit: logBufferLimit}
}

// Write 实现 io.Writer 接口，将日志按行追加到缓冲区，并在超出上限时自动淘汰旧行。
func (lb *LogBuffer) Write(p []byte) (int, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	logText := string(p)
	parts := strings.Split(logText, "\n")
	for _, line := range parts {
		if line == "" {
			continue
		}
		lb.lines = append(lb.lines, line)
	}
	if overflow := len(lb.lines) - lb.limit; overflow > 0 {
		lb.lines = append([]string(nil), lb.lines[overflow:]...)
	}
	return len(p), nil
}

// GetLines 返回当前所有日志副本，避免并发读写冲突。
func (lb *LogBuffer) GetLines() []string {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	result := make([]string, len(lb.lines))
	copy(result, lb.lines)
	return result
}

// InitialModel 构建 TUI 初始模型，设定默认焦点和状态。
func InitialModel(lb *LogBuffer) *appModel {
	m := &appModel{
		logBuffer:           lb,
		focus:               focusShelf,
		screen:              screenShelf,
		autoFollowLogs:      true,
		shelfRowToID:        make(map[int]int),
		currentShelfURL:     "",
		qrButtonFocus:       qrActionTerminalReader,
		coverProtocol:       detectTUIImageProtocol(),
		coverCache:          make(map[string]coverPreviewState),
		readerProtocol:      detectTUIReaderImageProtocol(),
		terminalReaderCache: make(map[string]terminalReaderState),
		readerAutoInterval:  defaultReaderAutoInterval,
	}
	m.setActionMsg(locale.GetString("tui_starting_service"))
	m.refreshData()
	return m
}

// Run 启动 TUI 模式；如果没有终端，则退回到普通服务模式。
func Run() error {
	// 现代终端（iTerm2、Terminal.app 等）将 East Asian Ambiguous 字符（含 Box Drawing）渲染为宽度 1，
	// 但 go-runewidth 在 zh_CN 等 CJK locale 下默认将其视为宽度 2，导致面板宽度计算偏差。
	runewidth.DefaultCondition.EastAsianWidth = false

	for _, arg := range os.Args {
		if arg == "-v" || arg == "--version" || arg == "-h" || arg == "--help" ||
			arg == "-u" || arg == "--upgrade" {
			cmd.Execute()
			return nil
		}
	}

	if shouldBypassTUI(os.Args) {
		return runWithoutTUI()
	}

	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return runWithoutTUI()
	}

	logBuffer := NewLogBuffer()
	stdLogOutput := stdlog.Writer()
	stdLogFlags := stdlog.Flags()
	stdLogPrefix := stdlog.Prefix()
	stdlog.SetOutput(logBuffer)
	logger.SetSuppressStdout(true)
	logger.SetOutput(io.Discard)
	logger.SetMirrorOutput(logBuffer)
	defer func() {
		stdlog.SetOutput(stdLogOutput)
		stdlog.SetFlags(stdLogFlags)
		stdlog.SetPrefix(stdLogPrefix)
		logger.SetMirrorOutput(nil)
		logger.SetSuppressStdout(false)
		logger.SetOutput(os.Stderr)
	}()

	program := tea.NewProgram(
		InitialModel(logBuffer),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := program.Run(); err != nil {
		return err
	}
	if err := routers.StopWebServer(); err != nil {
		return err
	}
	return nil
}

// shouldBypassTUI 在 Cobra 正式解析前识别 --no-tui/-n。
// TUI 会先于 cmd.Execute() 创建，因此这里必须提前判断一次启动入口。
func shouldBypassTUI(args []string) bool {
	bypass := false
	for _, arg := range args[1:] {
		if arg == "--" {
			break
		}
		if arg == "--no-tui" || arg == "-n" {
			bypass = true
			continue
		}
		if strings.HasPrefix(arg, "--no-tui=") {
			value, err := strconv.ParseBool(strings.TrimPrefix(arg, "--no-tui="))
			if err == nil {
				bypass = value
			}
			continue
		}
		if strings.HasPrefix(arg, "-n=") {
			value, err := strconv.ParseBool(strings.TrimPrefix(arg, "-n="))
			if err == nil {
				bypass = value
			}
		}
	}
	return bypass
}

// runWithoutTUI 在非终端环境下（如管道、重定向）退回普通服务模式启动。
func runWithoutTUI() error {
	cmd.Execute()
	if err := routers.StartWebServer(); err != nil {
		return err
	}
	routers.StartTailscale()
	cmd.LoadUserPlugins()
	cmd.AddStoreUrls(cmd.Args)
	cmd.SetCwdAsScanPathIfNeed()
	cmd.LoadMetadata()
	cmd.ScanStore()
	cmd.SaveMetadata()
	config.StartOrStopAutoRescan()
	// 非 TUI 模式没有右侧 QRCode 面板，需要在命令行最后打印阅读链接和二维码。
	cmd.ShowQRCode()
	cmd.OpenReaderBrowserIfNeeded()
	cmd.SetShutdownHandler()
	return nil
}

// tickCmd 返回一个定时 Cmd，每 tickInterval 触发一次数据刷新。
func tickCmd() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// startBackendCmd 在后台协程中依次启动 Web 服务、扫描书籍等全部后端流程。
func startBackendCmd() tea.Cmd {
	return func() (msg tea.Msg) {
		defer func() {
			if recovered := recover(); recovered != nil {
				msg = backendErrorMsg{err: fmt.Errorf("tui 后台启动异常: %v", recovered)}
			}
		}()
		// Comigo 后台服务启动
		cmd.Execute()
		if err := routers.StartWebServer(); err != nil {
			return backendErrorMsg{err: err}
		}
		routers.StartTailscale()
		cmd.LoadUserPlugins()
		cmd.AddStoreUrls(cmd.Args)
		cmd.SetCwdAsScanPathIfNeed()
		cmd.LoadMetadata()
		cmd.ScanStore()
		cmd.SaveMetadata()
		config.StartOrStopAutoRescan()
		go cmd.SetShutdownHandler()
		return backendStartedMsg{}
	}
}

// openURLCmd 调用系统默认浏览器打开指定 URL（异步执行，结果通过 openURLResultMsg 回传）。
func openURLCmd(target string) tea.Cmd {
	return func() tea.Msg {
		err := tools.OpenURL(target)
		return openURLResultMsg{url: target, err: err}
	}
}

// Init 实现 tea.Model 接口，启动定时器和后台服务。
func (m *appModel) Init() tea.Cmd {
	return tea.Batch(tickCmd(), startBackendCmd())
}

// Update 实现 tea.Model 接口，处理所有消息分发：窗口变化、定时刷新、后台事件、键盘/鼠标输入。
func (m *appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.refreshData()
		return m, m.syncActiveImageCmd()
	case tickMsg:
		if m.screen == screenReader {
			return m, tea.Batch(tickCmd(), m.handleReaderAutoFlip(time.Time(msg)))
		}
		m.refreshData()
		return m, tea.Batch(tickCmd(), m.handleReaderAutoFlip(time.Time(msg)), m.syncActiveImageCmd())
	case backendStartedMsg:
		m.backendReady = true
		m.backendError = ""
		m.setActionMsg(locale.GetString("tui_service_started"))
		m.refreshData()
		return m, m.syncActiveImageCmd()
	case backendErrorMsg:
		m.backendError = msg.err.Error()
		m.setActionMsg(locale.GetString("tui_backend_failed"))
		m.refreshData()
		return m, m.syncActiveImageCmd()
	case openURLResultMsg:
		if msg.err != nil {
			text := fmt.Sprintf(locale.GetString("tui_open_browser_failed"), msg.err.Error())
			m.setActionMsg(shortenText(text, maxActionMessage))
			m.showModal(locale.GetString("tui_modal_title_notice"), text)
			logger.Infof(locale.GetString("tui_open_browser_failed"), msg.err.Error())
		} else {
			m.setActionMsg(shortenText(fmt.Sprintf(locale.GetString("tui_opened_url"), msg.url), maxActionMessage))
			logger.Infof(locale.GetString("log_opening_browser"), msg.url)
		}
		m.refreshData()
		return m, m.syncActiveImageCmd()
	case coverPreviewMsg:
		m.applyCoverPreviewMsg(msg)
		return m, nil
	case terminalReaderPageMsg:
		m.applyTerminalReaderPageMsg(msg)
		return m, nil
	case tea.MouseMsg:
		return m.handleMouse(msg)
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

// handleKey 处理键盘输入，包含全局快捷键（退出、Tab 切换面板）和各面板专属按键。
func (m *appModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	}
	if m.modal.Visible {
		return m.handleModalKey(msg)
	}

	switch msg.String() {
	case "q":
		if m.screen == screenReader {
			return m, m.exitTerminalReader()
		}
		return m, tea.Quit
	case "c":
		return m, m.toggleTUIImageMode()
	case "tab":
		if m.screen == screenReader {
			return m, nil
		}
		m.moveFocus(1)
		return m, nil
	case "shift+tab":
		if m.screen == screenReader {
			return m, nil
		}
		m.moveFocus(-1)
		return m, nil
	}

	if m.screen == screenReader {
		return m.handleTerminalReaderKey(msg)
	}

	switch m.focus {
	case focusShelf:
		switch msg.String() {
		case "up", "k":
			m.moveSelection(-1)
		case "down", "j":
			m.moveSelection(1)
		case "pgup":
			m.moveSelection(-5)
		case "pgdown":
			m.moveSelection(5)
		case "home":
			m.selectFirst()
		case "end":
			m.selectLast()
		case "left", "h", "backspace", "esc":
			m.goBack()
		case "enter", " ":
			return m, m.activateSelectedItem()
		case "r":
			m.refreshData()
		}
	case focusLog:
		switch msg.String() {
		case "up", "k":
			m.scrollLogs(-1)
		case "down", "j":
			m.scrollLogs(1)
		case "pgup", "b":
			m.scrollLogs(-max(1, m.layout().log.innerHeight()/2))
		case "pgdown", "f":
			m.scrollLogs(max(1, m.layout().log.innerHeight()/2))
		case "home", "g":
			m.logOffset = 0
			m.autoFollowLogs = false
		case "end", "G":
			m.autoFollowLogs = true
			m.syncLogOffset()
		}
	case focusQRCode:
		switch msg.String() {
		case "up", "k":
			m.moveQRButtonVertical(-1)
		case "down", "j":
			m.moveQRButtonVertical(1)
		case "left", "h":
			m.moveQRButtonHorizontal(-1)
		case "right", "l":
			m.moveQRButtonHorizontal(1)
		case "tab":
			m.moveQRButtonTab()
		case "enter", " ":
			return m, m.executeQRButton()
		}
	}

	m.refreshStatus()
	return m, m.syncActiveImageCmd()
}

// handleMouse 处理鼠标事件：点击切换面板焦点、书架选中项、QR 面板按钮及滚轮操作。
func (m *appModel) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.modal.Visible {
		return m.handleModalMouse(msg)
	}
	if m.screen == screenReader {
		return m.handleTerminalReaderMouse(msg)
	}

	layout := m.layout()
	if msg.Action != tea.MouseActionPress {
		return m, nil
	}

	switch {
	case contains(layout.shelf, msg.X, msg.Y):
		m.focus = focusShelf
		if msg.Button == tea.MouseButtonWheelUp {
			m.moveSelection(-1)
			return m, m.syncActiveImageCmd()
		}
		if msg.Button == tea.MouseButtonWheelDown {
			m.moveSelection(1)
			return m, m.syncActiveImageCmd()
		}
		if msg.Button == tea.MouseButtonLeft {
			row := msg.Y - layout.shelf.y - 1
			if idx, ok := m.shelfRowToID[row]; ok {
				m.selected = idx
				m.syncShelfOffset()
				m.refreshQRCode()
				m.refreshStatus()
				if m.isShelfDoubleClick(idx, time.Now()) {
					return m, m.activateShelfDoubleClick()
				}
				return m, m.syncActiveImageCmd()
			}
		}
	case contains(layout.log, msg.X, msg.Y):
		m.focus = focusLog
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			m.scrollLogs(-2)
		case tea.MouseButtonWheelDown:
			m.scrollLogs(2)
		}
	case contains(layout.qr, msg.X, msg.Y):
		m.focus = focusQRCode
		if msg.Button == tea.MouseButtonLeft {
			clickRow := msg.Y - layout.qr.y - 1
			col := msg.X - layout.qr.x - 1
			for _, hitbox := range m.qrButtonHitbox {
				if clickRow == hitbox.row && col >= hitbox.start && col < hitbox.end {
					m.qrButtonFocus = hitbox.action
					return m, m.executeQRButton()
				}
			}
		}
	}
	m.refreshStatus()
	return m, m.syncActiveImageCmd()
}

// isShelfDoubleClick 判断本次书架点击是否构成双击；单击仍会立即更新选中项。
func (m *appModel) isShelfDoubleClick(index int, now time.Time) bool {
	elapsed := now.Sub(m.lastShelfClickAt)
	doubleClick := index == m.lastShelfClickIndex && !m.lastShelfClickAt.IsZero() && elapsed >= 0 && elapsed <= doubleClickGap
	m.lastShelfClickIndex = index
	m.lastShelfClickAt = now
	return doubleClick
}

// activateShelfDoubleClick 双击书籍时按当前打开方式执行；双击分组/返回项仍沿用原激活逻辑。
func (m *appModel) activateShelfDoubleClick() tea.Cmd {
	item := m.currentItem()
	if item == nil {
		return nil
	}
	if item.Kind != shelfItemBook {
		return m.activateSelectedItem()
	}
	if m.qrButtonFocus == qrActionOpenBrowser {
		return m.executeQRButton()
	}
	return m.startTerminalReader(item)
}

// View 实现 tea.Model 接口，根据当前终端尺寸渲染全部面板并拼接为最终输出字符串。
func (m *appModel) View() string {
	if m.modal.Visible {
		return m.renderModalView()
	}
	if m.screen == screenReader {
		return m.renderTerminalReaderView()
	}

	if m.width < minTUIWidth || m.height < minTUIHeight {
		return fmt.Sprintf(locale.GetString("tui_terminal_too_small"), minTUIWidth, minTUIHeight) + "\n" +
			fmt.Sprintf(locale.GetString("tui_current_size"), m.width, m.height) + "\n"
	}

	layout := m.layout()

	if m.isNarrow() {
		shelfLines := m.makePanel(locale.GetString("tui_panel_shelf"), m.renderShelfContent(layout.shelf), layout.shelf, m.focus == focusShelf)
		qrLines := m.makePanel("QRCode", m.renderQRCodeContent(layout.qr), layout.qr, m.focus == focusQRCode)
		logLines := m.makePanel(locale.GetString("tui_panel_log"), m.renderLogContent(layout.log), layout.log, m.focus == focusLog)
		var all []string
		all = append(all, shelfLines...)
		all = append(all, qrLines...)
		all = append(all, logLines...)
		return strings.Join(all, "\n")
	}

	shelfLines := m.makePanel(locale.GetString("tui_panel_shelf"), m.renderShelfContent(layout.shelf), layout.shelf, m.focus == focusShelf)
	qrLines := m.makePanel("QRCode", m.renderQRCodeContent(layout.qr), layout.qr, m.focus == focusQRCode)
	logLines := m.makePanel(locale.GetString("tui_panel_log"), m.renderLogContent(layout.log), layout.log, m.focus == focusLog)
	coverLines := m.makeStyledPanel(locale.GetString("tui_panel_preview"), m.renderCoverPreviewContent(layout.cover), layout.cover, false)

	rows := mergeWideLayoutRows(layout, shelfLines, logLines, qrLines, coverLines)
	return m.renderCoverClearPrefix(layout.cover) + m.renderCoverSetupPrefix() + strings.Join(rows, "\n") + m.renderCoverOverlay(layout.cover)
}

// refreshData 一次性刷新日志、书架、系统状态和二维码等全部数据。
func (m *appModel) refreshData() {
	m.logs = m.logBuffer.GetLines()
	m.syncLogOffset()
	m.refreshShelf()
	m.refreshStatus()
	m.refreshQRCode()
}

// refreshShelf 重新构建当前层级的书架条目列表并更新选中项和书架 URL。
func (m *appModel) refreshShelf() {
	m.items = m.buildCurrentShelfItems()
	m.ensureSelectedItem()
	m.currentShelfURL = m.buildCurrentShelfURL()
}

// refreshStatus 更新当前界面底部需要展示的服务状态。
func (m *appModel) refreshStatus() {
	statusText := locale.GetString("tui_status_starting")
	if m.backendError != "" {
		statusText = locale.GetString("tui_status_failed")
	} else if m.backendReady {
		statusText = locale.GetString("tui_status_running")
	}
	if m.actionMsg != "" {
		statusText = m.actionMsg
	}

	m.status = systemSnapshot{
		StatusText: statusText,
	}
}

// setActionMsg 统一设置操作提示消息，用于 QR 面板和信息面板展示。
func (m *appModel) setActionMsg(msg string) {
	m.actionMsg = msg
}

// selectedURL 返回当前选中项的 URL；若无可选项则回退到当前书架 URL。
func (m *appModel) selectedURL() string {
	if item := m.currentItem(); item != nil && item.TargetURL != "" {
		return item.TargetURL
	}
	return m.currentShelfURL
}

// refreshQRCode 根据当前选中项的 URL 重新生成 QR 码字符行。
func (m *appModel) refreshQRCode() {
	lines, err := renderQRCodeLines(m.selectedURL())
	if err != nil {
		m.qrLines = []string{locale.GetString("tui_qr_gen_failed"), err.Error()}
		return
	}
	m.qrLines = lines
}

// buildCurrentShelfItems 根据导航栈层级构建当前应显示的书架条目列表。
func (m *appModel) buildCurrentShelfItems() []shelfItem {
	if modelpkg.IStore == nil {
		return []shelfItem{{
			Kind:       shelfItemHeader,
			Title:      locale.GetString("tui_shelf_waiting_init"),
			Subtitle:   locale.GetString("tui_shelf_waiting_hint"),
			Selectable: false,
		}}
	}

	if len(m.stack) == 0 {
		return buildTopShelfItems(m.readMode)
	}
	return buildChildShelfItems(m.stack[len(m.stack)-1], m.readMode)
}

// buildTopShelfItems 构建顶层书架列表（按 StoreUrl 分组，depth=0 的书籍）。
func buildTopShelfItems(readMode int) []shelfItem {
	if modelpkg.IStore == nil {
		return []shelfItem{{
			Kind:       shelfItemHeader,
			Title:      locale.GetString("tui_shelf_empty"),
			Subtitle:   locale.GetString("tui_shelf_empty_hint"),
			Selectable: false,
		}}
	}

	allBooks, err := modelpkg.IStore.ListBooks()
	if err != nil || len(allBooks) == 0 {
		return []shelfItem{{
			Kind:       shelfItemHeader,
			Title:      locale.GetString("tui_shelf_empty"),
			Subtitle:   locale.GetString("tui_shelf_empty_hint"),
			Selectable: false,
		}}
	}

	var items []shelfItem
	for _, rawStoreURL := range config.GetCfg().StoreUrls {
		storeID := normalizeStoreID(rawStoreURL)
		var topBooks modelpkg.BookInfos
		childBookNum := 0
		for _, book := range allBooks {
			if book.StoreUrl != storeID {
				continue
			}
			if book.Depth == 0 {
				topBooks = append(topBooks, book.BookInfo)
			}
			if book.Type != modelpkg.TypeBooksGroup {
				childBookNum++
			}
		}
		topBooks.SortBooks("default")
		items = append(items, shelfItem{
			Kind:       shelfItemHeader,
			Title:      displayStoreName(rawStoreURL),
			Subtitle:   fmt.Sprintf(locale.GetString("tui_readable_books_count"), childBookNum),
			Selectable: false,
		})
		for _, book := range topBooks {
			items = append(items, convertBookInfo(book, readMode))
		}
	}
	return items
}

// buildChildShelfItems 构建子书架内容列表，包含返回按钮和该书架下的所有子项。
func buildChildShelfItems(level shelfLevel, readMode int) []shelfItem {
	if modelpkg.IStore == nil {
		return []shelfItem{{Kind: shelfItemHeader, Title: locale.GetString("tui_shelf_not_initialized"), Selectable: false}}
	}

	parentBook, err := modelpkg.IStore.GetBook(level.BookID)
	if err != nil {
		return []shelfItem{
			{Kind: shelfItemBack, Title: locale.GetString("tui_go_back"), Selectable: true},
			{Kind: shelfItemHeader, Title: locale.GetString("tui_sub_shelf_not_found"), Selectable: false},
		}
	}

	var childBooks modelpkg.BookInfos
	for _, childID := range parentBook.ChildBooksID {
		book, childErr := modelpkg.IStore.GetBook(childID)
		if childErr != nil {
			continue
		}
		childBooks = append(childBooks, book.BookInfo)
	}
	childBooks.SortBooks("default")
	if len(childBooks) == 0 {
		return []shelfItem{
			{
				Kind:       shelfItemBack,
				Title:      locale.GetString("tui_go_back"),
				Subtitle:   "",
				Selectable: true,
			},
			{
				Kind:       shelfItemHeader,
				Title:      locale.GetString("tui_sub_shelf_no_content"),
				Subtitle:   "",
				Selectable: false,
			},
		}
	}

	items := []shelfItem{{
		Kind:       shelfItemBack,
		Title:      locale.GetString("tui_go_back"),
		Subtitle:   level.Title,
		Selectable: true,
	}}
	for _, book := range childBooks {
		items = append(items, convertBookInfo(book, readMode))
	}
	return items
}

// isShelfEffectivelyEmpty 判断当前层级是否无可展示书籍（忽略 header/back，仅 group/book 视为有内容）。
func isShelfEffectivelyEmpty(items []shelfItem) bool {
	for _, item := range items {
		if item.Kind == shelfItemBook || item.Kind == shelfItemGroup {
			return false
		}
	}
	return true
}

// convertBookInfo 将 BookInfo 转换为 shelfItem，自动判断类型（书籍/子书架）并生成目标 URL。
func convertBookInfo(book modelpkg.BookInfo, readMode int) shelfItem {
	kind := shelfItemBook
	subtitle := buildBookSubtitle(book)
	if book.Type == modelpkg.TypeBooksGroup {
		kind = shelfItemGroup
	}
	return shelfItem{
		Kind:       kind,
		Title:      book.ShortName(),
		Subtitle:   subtitle,
		BookID:     book.BookID,
		TargetURL:  buildBookTargetURL(book, readMode),
		Selectable: true,
	}
}

// buildBookSubtitle 根据书籍类型生成副标题文字（页数、子项数或媒体类型）。
func buildBookSubtitle(book modelpkg.BookInfo) string {
	switch book.Type {
	case modelpkg.TypeBooksGroup:
		return fmt.Sprintf(locale.GetString("tui_sub_shelf_items"), book.ChildBooksNum)
	case modelpkg.TypeVideo:
		return locale.GetString("tui_type_video")
	case modelpkg.TypeAudio:
		return locale.GetString("tui_type_audio")
	case modelpkg.TypeHTML:
		return locale.GetString("tui_type_html")
	case modelpkg.TypeUnknownFile:
		return locale.GetString("tui_type_raw")
	default:
		pageText := ""
		if book.PageCount > 0 {
			pageText = fmt.Sprintf(locale.GetString("tui_page_count"), book.PageCount)
		}
		if pageText == "" {
			return string(book.Type)
		}
		return fmt.Sprintf("%s | %s", string(book.Type), pageText)
	}
}

// buildBookTargetURL 根据书籍类型和阅读模式构建打开书籍的完整 URL。
// 对普通图片类书籍，readMode=0 使用 /scroll/ 路径，readMode=1 使用 /flip/ 路径。
func buildBookTargetURL(book modelpkg.BookInfo, readMode int) string {
	baseURL := buildBaseURL()
	base := strings.TrimRight(baseURL, "/")
	switch book.Type {
	case modelpkg.TypeBooksGroup:
		return base + config.PrefixPath("/shelf/"+book.BookID)
	case modelpkg.TypeVideo, modelpkg.TypeAudio:
		return base + config.PrefixPath("/player/"+book.BookID)
	case modelpkg.TypeHTML, modelpkg.TypeUnknownFile:
		return base + config.PrefixPath("/api/raw/"+book.BookID+"/"+url.QueryEscape(book.Title))
	default:
		prefix := "/scroll/"
		if readMode == 1 {
			prefix = "/flip/"
		}
		target := base + config.PrefixPath(prefix+book.BookID)
		if readMode == 0 && modelpkg.IStore != nil {
			if marks, err := modelpkg.IStore.GetBookMarks(book.BookID); err == nil && marks != nil {
				if start := marks.GetLastReadPage(); start > 1 {
					target += "?start=" + strconv.Itoa(start)
				}
			}
		}
		return target
	}
}

// buildBaseURL 根据当前配置（Host、Port、TLS 等）构建 Web 服务的基础 URL。
func buildBaseURL() string {
	cfg := config.GetCfg()
	protocol := "http://"
	if (cfg.CertFile != "" && cfg.KeyFile != "") || cfg.AutoTLSCertificate {
		protocol = "https://"
	}

	host := cfg.Host
	if host == "" {
		if cfg.DisableLAN {
			host = "127.0.0.1"
		} else {
			// TUI 只展示一个可访问地址：未指定 Host 时使用系统默认路由选出的出站 IP，避免 ZeroTier 等虚拟网卡因枚举顺序排在前面。
			if outboundIP, err := tools.LookupOutboundIP(); err == nil {
				host = outboundIP.String()
			} else {
				host = "127.0.0.1"
			}
		}
	} else {
		if tools.IsLoopbackHost(host) {
			host = tools.GetOutboundIP().String()
		}
	}
	if cfg.AutoTLSCertificate {
		return protocol + host
	}
	return fmt.Sprintf("%s%s:%d", protocol, host, cfg.Port)
}

// buildCurrentShelfURL 构建当前书架层级对应的 Web URL（顶层为根 URL，子层为 /shelf/{id}）。
func (m *appModel) buildCurrentShelfURL() string {
	baseURL := strings.TrimRight(buildBaseURL(), "/")
	if len(m.stack) == 0 {
		return baseURL + config.PrefixPath("/")
	}
	return baseURL + config.PrefixPath("/shelf/"+m.stack[len(m.stack)-1].BookID)
}

// displayStoreName 从 Store URL/路径中提取适合显示的简短名称。
func displayStoreName(storeURL string) string {
	if strings.Contains(storeURL, "://") {
		return storeURL
	}
	base := filepath.Base(storeURL)
	if base == "." || base == string(filepath.Separator) {
		return storeURL
	}
	return base
}

// normalizeStoreID 将本地路径转为绝对路径作为 Store 的统一标识；远程 URL 保持原样。
func normalizeStoreID(storeURL string) string {
	if strings.Contains(storeURL, "://") {
		return storeURL
	}
	storePathAbs, err := filepath.Abs(storeURL)
	if err != nil {
		return storeURL
	}
	return storePathAbs
}

// currentItem 返回当前选中的可选条目指针；若索引越界或不可选则返回 nil。
func (m *appModel) currentItem() *shelfItem {
	if len(m.items) == 0 || m.selected < 0 || m.selected >= len(m.items) {
		return nil
	}
	if !m.items[m.selected].Selectable {
		return nil
	}
	return &m.items[m.selected]
}

// ensureSelectedItem 保证 selected 指向一个有效的可选条目；若当前不可选则搜索第一个可选项。
func (m *appModel) ensureSelectedItem() {
	if len(m.items) == 0 {
		m.selected = 0
		m.shelfOffset = 0
		return
	}
	if m.selected < 0 {
		m.selected = 0
	}
	if m.selected >= len(m.items) {
		m.selected = len(m.items) - 1
	}
	if m.items[m.selected].Selectable {
		m.syncShelfOffset()
		return
	}
	for i := range m.items {
		if m.items[i].Selectable {
			m.selected = i
			m.syncShelfOffset()
			return
		}
	}
	m.selected = 0
	m.shelfOffset = 0
}

// selectFirst 选中列表中第一个可选条目。
func (m *appModel) selectFirst() {
	for i := range m.items {
		if m.items[i].Selectable {
			m.selected = i
			m.syncShelfOffset()
			return
		}
	}
}

// selectLast 选中列表中最后一个可选条目。
func (m *appModel) selectLast() {
	for i := len(m.items) - 1; i >= 0; i-- {
		if m.items[i].Selectable {
			m.selected = i
			m.syncShelfOffset()
			return
		}
	}
}

// syncShelfOffset 调整滚动偏移，使当前选中项始终在可视区域内。
func (m *appModel) syncShelfOffset() {
	const headerLines = 2 // 面包屑 + 操作提示
	visible := m.layout().shelf.innerHeight() - headerLines
	if visible <= 0 {
		m.shelfOffset = 0
		return
	}
	if m.selected < m.shelfOffset {
		m.shelfOffset = m.selected
	}
	if m.selected >= m.shelfOffset+visible {
		m.shelfOffset = m.selected - visible + 1
	}
	maxOffset := max(0, len(m.items)-visible)
	if m.shelfOffset > maxOffset {
		m.shelfOffset = maxOffset
	}
	if m.shelfOffset < 0 {
		m.shelfOffset = 0
	}
}

// moveSelection 将选中项向上（delta<0）或向下（delta>0）移动指定步数，自动跳过不可选条目。
func (m *appModel) moveSelection(delta int) {
	if len(m.items) == 0 {
		return
	}
	index := m.selected
	step := 1
	if delta < 0 {
		step = -1
		delta = -delta
	}
	for i := 0; i < delta; i++ {
		next := index
		for {
			next += step
			if next < 0 || next >= len(m.items) {
				break
			}
			if m.items[next].Selectable {
				index = next
				break
			}
		}
	}
	m.selected = index
	m.syncShelfOffset()
}

// goBack 弹出导航栈顶层，返回上一级书架。
func (m *appModel) goBack() {
	if len(m.stack) == 0 {
		return
	}
	m.stack = m.stack[:len(m.stack)-1]
	m.selected = 0
	m.shelfOffset = 0
	m.refreshData()
}

// activateSelectedItem 激活当前选中项：返回上级、进入子书架或打开书籍 URL。
func (m *appModel) activateSelectedItem() tea.Cmd {
	item := m.currentItem()
	if item == nil {
		return nil
	}

	switch item.Kind {
	case shelfItemBack:
		m.goBack()
		return m.syncCoverPreviewCmd()
	case shelfItemGroup:
		m.stack = append(m.stack, shelfLevel{
			BookID: item.BookID,
			Title:  item.Title,
		})
		m.selected = 0
		m.shelfOffset = 0
		m.setActionMsg(fmt.Sprintf(locale.GetString("tui_entered_sub_shelf"), item.Title))
		m.refreshData()
		return m.syncCoverPreviewCmd()
	case shelfItemBook:
		return m.startTerminalReader(item)
	default:
		return nil
	}
}

// executeQRButton 执行 QR 面板当前聚焦按钮的动作。
func (m *appModel) executeQRButton() tea.Cmd {
	switch m.qrButtonFocus {
	case qrActionScrollMode:
		m.setReadMode(0)
		return nil
	case qrActionFlipMode:
		m.setReadMode(1)
		return nil
	case qrActionTerminalReader:
		if item := m.currentItem(); item != nil && item.Kind == shelfItemBook {
			return m.startTerminalReader(item)
		}
		m.setActionMsg(locale.GetString("tui_terminal_reader_no_book"))
		m.refreshStatus()
		return nil
	}

	target := m.selectedURL()
	if target == "" {
		m.setActionMsg(locale.GetString("tui_no_url_available"))
		m.refreshStatus()
		return nil
	}
	switch m.qrButtonFocus {
	case qrActionOpenBrowser:
		m.setActionMsg(fmt.Sprintf(locale.GetString("tui_opening_url"), shortenText(target, maxActionMessage-6)))
		m.refreshStatus()
		return openURLCmd(target)
	case qrActionCopyURL:
		if err := clipboard.WriteAll(target); err != nil {
			m.setActionMsg(fmt.Sprintf(locale.GetString("tui_copy_failed"), err.Error()))
			logger.Infof(locale.GetString("tui_copy_failed"), err.Error())
		} else {
			m.setActionMsg(fmt.Sprintf(locale.GetString("tui_url_copied"), shortenText(target, maxActionMessage-8)))
			logger.Infof(locale.GetString("log_copied_url_to_clipboard"), target)
		}
		m.refreshStatus()
		return nil
	default:
		return nil
	}
}

// setReadMode 设置 Web 阅读 URL 使用的模式，并刷新书架和二维码。
func (m *appModel) setReadMode(readMode int) {
	if readMode != 0 && readMode != 1 {
		return
	}
	m.readMode = readMode
	m.refreshShelf()
	m.refreshQRCode()
	m.refreshStatus()
}

// scrollLogs 滚动日志面板偏移量，正数向下、负数向上，自动限定边界。
func (m *appModel) scrollLogs(delta int) {
	visible := max(1, m.layout().log.innerHeight())
	maxOffset := max(0, len(m.logs)-visible)
	m.logOffset += delta
	if m.logOffset < 0 {
		m.logOffset = 0
	}
	if m.logOffset > maxOffset {
		m.logOffset = maxOffset
	}
	m.autoFollowLogs = m.logOffset >= maxOffset
}

// syncLogOffset 同步日志滚动偏移，若开启自动跟随则始终定位到最新日志。
func (m *appModel) syncLogOffset() {
	visible := max(1, m.layout().log.innerHeight())
	maxOffset := max(0, len(m.logs)-visible)
	if m.autoFollowLogs {
		m.logOffset = maxOffset
		return
	}
	if m.logOffset > maxOffset {
		m.logOffset = maxOffset
	}
}

// isNarrow 判断当前终端宽度是否小于窄屏阈值。
func (m *appModel) isNarrow() bool {
	return m.width < narrowThreshold
}

// moveFocus 在当前可交互面板之间切换。
func (m *appModel) moveFocus(delta int) {
	order := []focusPanel{focusShelf, focusLog, focusQRCode}
	current := 0
	for i, panel := range order {
		if panel == m.focus {
			current = i
			break
		}
	}
	next := (current + delta) % len(order)
	if next < 0 {
		next += len(order)
	}
	m.focus = order[next]
}

// renderWidth 返回 TUI 实际写入的行宽。
// 终端最后一列写满时容易触发自动换行，下一帧局部刷新会留下上一帧的文字残影。
func (m *appModel) renderWidth() int {
	return max(0, m.width-1)
}

// layout 计算 TUI 面板矩形布局。宽屏使用 2×2 网格，窄屏使用垂直单栏（隐藏封面预览）。
func (m *appModel) layout() layoutState {
	width := m.renderWidth()
	height := max(0, m.height)

	if m.isNarrow() {
		// 窄屏单栏：书架(3) / QR(5) / 日志(2)，封面预览隐藏。
		logH := height * 2 / 10
		shelfH := height * 3 / 10
		qrH := height - shelfH - logH
		qrY := shelfH
		logY := shelfH + qrH
		return layoutState{
			shelf: panelRect{x: 0, y: 0, w: width, h: shelfH},
			cover: panelRect{},
			qr:    panelRect{x: 0, y: qrY, w: width, h: qrH},
			log:   panelRect{x: 0, y: logY, w: width, h: logH},
		}
	}

	leftWidth := (width - layoutGap) * 2 / 3
	rightWidth := width - layoutGap - leftWidth
	split := height - layoutGap
	leftTopHeight := (split * 2) / 3
	leftBottomHeight := split - leftTopHeight
	// 右侧 QRCode 与预览区按 1:1 切分；奇数高度时把多出的 1 行给 QRCode。
	rightTopHeight := (split + 1) / 2
	rightBottomHeight := split - rightTopHeight

	return layoutState{
		shelf: panelRect{x: 0, y: 0, w: leftWidth, h: leftTopHeight},
		log:   panelRect{x: 0, y: leftTopHeight + layoutGap, w: leftWidth, h: leftBottomHeight},
		qr:    panelRect{x: leftWidth + layoutGap, y: 0, w: rightWidth, h: rightTopHeight},
		cover: panelRect{x: leftWidth + layoutGap, y: rightTopHeight + layoutGap, w: rightWidth, h: rightBottomHeight},
	}
}

// mergeWideLayoutRows 按绝对 y 坐标合并左右两列，允许左右列使用不同的上下高度比例。
func mergeWideLayoutRows(layout layoutState, shelfLines []string, logLines []string, qrLines []string, coverLines []string) []string {
	height := max(layout.shelf.y+layout.shelf.h, layout.log.y+layout.log.h)
	height = max(height, layout.qr.y+layout.qr.h)
	height = max(height, layout.cover.y+layout.cover.h)
	widthLeft := layout.shelf.w
	widthRight := layout.qr.w
	rows := make([]string, 0, height)
	for y := 0; y < height; y++ {
		left := panelLineAt(y, layout.shelf, shelfLines)
		if left == "" {
			left = panelLineAt(y, layout.log, logLines)
		}
		right := panelLineAt(y, layout.qr, qrLines)
		if right == "" {
			right = panelLineAt(y, layout.cover, coverLines)
		}
		rows = append(rows, clipAndPadStyled(left, widthLeft)+strings.Repeat(" ", layoutGap)+clipAndPadStyled(right, widthRight))
	}
	return rows
}

// panelLineAt 取指定绝对行对应的面板内容；不在面板范围内时返回空字符串。
func panelLineAt(row int, rect panelRect, lines []string) string {
	if rect.w <= 0 || rect.h <= 0 || row < rect.y || row >= rect.y+rect.h {
		return ""
	}
	index := row - rect.y
	if index < 0 || index >= len(lines) {
		return ""
	}
	return lines[index]
}

// renderShelfContent 渲染书架面板内容：面包屑、操作提示和可滚动条目列表。
func (m *appModel) renderShelfContent(rect panelRect) []string {
	inner := rect.innerHeight()
	if inner <= 0 {
		return nil
	}
	w := rect.innerWidth()
	lines := make([]string, 0, inner)
	m.shelfRowToID = make(map[int]int)

	rootDir := locale.GetString("tui_root_dir")
	breadcrumb := rootDir
	if len(m.stack) > 0 {
		parts := []string{rootDir}
		for _, level := range m.stack {
			parts = append(parts, level.Title)
		}
		breadcrumb = strings.Join(parts, " / ")
	}
	lines = append(lines, locale.GetString("tui_path_prefix")+breadcrumb)
	lines = append(lines, locale.GetString("tui_controls_hint"))

	visibleItems := max(0, inner-len(lines))
	if isShelfEffectivelyEmpty(m.items) {
		if visibleItems > 0 {
			lines = append(lines, locale.GetString("tui_no_shelf_content"))
		}
	} else {
		endIdx := min(len(m.items), m.shelfOffset+visibleItems)
		for idx := m.shelfOffset; idx < endIdx; idx++ {
			item := m.items[idx]
			line := m.formatShelfLine(item, idx == m.selected)
			m.shelfRowToID[len(lines)] = idx
			if !item.Selectable {
				delete(m.shelfRowToID, len(lines))
			}
			lines = append(lines, line)
		}
	}

	return fitLines(lines, w, inner)
}

// formatShelfLine 格式化单个书架条目为显示行，选中项以 "> " 前缀标记。
func (m *appModel) formatShelfLine(item shelfItem, selected bool) string {
	prefix := "  "
	if selected {
		prefix = "> "
	}

	switch item.Kind {
	case shelfItemHeader:
		return locale.GetString("tui_tag_store") + " " + item.Title + " | " + item.Subtitle
	case shelfItemBack:
		return prefix + locale.GetString("tui_tag_back") + " " + item.Title
	case shelfItemGroup:
		return prefix + locale.GetString("tui_tag_group") + " " + item.Title + " | " + item.Subtitle
	default:
		return prefix + locale.GetString("tui_tag_book") + " " + item.Title + " | " + item.Subtitle
	}
}

// renderLogContent 渲染日志面板内容，支持滚动和自动跟随最新日志；底部固定显示当前状态提示。
func (m *appModel) renderLogContent(rect panelRect) []string {
	height := rect.innerHeight()
	if height <= 0 {
		return nil
	}
	w := rect.innerWidth()
	if height == 1 {
		return []string{clipAndPad(shortenText(m.status.StatusText, w), w)}
	}

	const footerLines = 2 // 分隔线 + 底行服务状态
	bodyHeight := max(0, height-footerLines)
	lines := make([]string, 0, height)
	if len(m.logs) == 0 {
		if bodyHeight > 0 {
			lines = append(lines, locale.GetString("tui_no_logs"))
		}
	} else {
		start := min(m.logOffset, len(m.logs))
		end := min(len(m.logs), start+bodyHeight)
		lines = append(lines, m.logs[start:end]...)
		if !m.autoFollowLogs && len(lines) < bodyHeight {
			lines = append(lines, fmt.Sprintf(locale.GetString("tui_log_scrolling"), end, len(m.logs)))
		}
	}

	for len(lines) < bodyHeight {
		lines = append(lines, "")
	}
	lines = append(lines, padRightWith("", w, "─"))
	lines = append(lines, clipAndPad(shortenText(m.status.StatusText, w), w))
	return fitLines(lines, w, height)
}

// renderQRCodeContent 渲染 QRCode 面板内容：选中项 URL、二维码、阅读模式和操作按钮。
func (m *appModel) renderQRCodeContent(rect panelRect) []string {
	w := rect.innerWidth()
	m.qrButtonHitbox = nil
	selURL := m.selectedURL()
	label := locale.GetString("tui_qr_shelf_url")
	if item := m.currentItem(); item != nil && item.TargetURL != "" {
		label = fmt.Sprintf(locale.GetString("tui_qr_selected"), shortenText(item.Title, max(10, w-6)))
	}
	lines := []string{label, selURL}
	if len(m.qrLines) == 0 {
		lines = append(lines, centerText(locale.GetString("tui_qr_unavailable"), w))
	} else {
		for _, qrLine := range m.qrLines {
			lines = append(lines, centerText(qrLine, w))
		}
	}

	lines = m.appendQRCodeButtonRow(lines, w, []qrAction{qrActionScrollMode, qrActionFlipMode, qrActionCopyURL})
	lines = m.appendQRCodeButtonRow(lines, w, []qrAction{qrActionTerminalReader, qrActionOpenBrowser})

	// 操作结果提示行（浏览器打开/URL复制等；加标记强调，不自动隐藏）
	if m.actionMsg != "" {
		lines = append(lines, "")
		lines = append(lines, formatTUIActionHint(m.actionMsg, w))
	}

	return fitLines(lines, w, rect.innerHeight())
}

// appendQRCodeButtonRow 追加一行 QRCode 操作按钮，并记录鼠标点击范围。
func (m *appModel) appendQRCodeButtonRow(lines []string, width int, actions []qrAction) []string {
	if width <= 0 || len(actions) == 0 {
		return lines
	}
	parts := make([]string, len(actions))
	totalWidth := 0
	for i, action := range actions {
		parts[i] = m.formatQRButton(action)
		totalWidth += runewidth.StringWidth(parts[i])
		if i > 0 {
			totalWidth += 2
		}
	}

	row := len(lines)
	offset := 0
	if totalWidth < width {
		offset = (width - totalWidth) / 2
	}
	col := offset
	var builder strings.Builder
	builder.WriteString(strings.Repeat(" ", offset))
	for i, action := range actions {
		if i > 0 {
			builder.WriteString("  ")
			col += 2
		}
		partWidth := runewidth.StringWidth(parts[i])
		m.qrButtonHitbox = append(m.qrButtonHitbox, qrButtonHitbox{
			action: action,
			row:    row,
			start:  col,
			end:    col + partWidth,
		})
		builder.WriteString(parts[i])
		col += partWidth
	}
	return append(lines, builder.String())
}

// formatQRButton 根据当前阅读模式和按钮焦点渲染按钮文本。
func (m *appModel) formatQRButton(action qrAction) string {
	label := qrActionLabel(action)
	activeMode := (action == qrActionScrollMode && m.readMode == 0) || (action == qrActionFlipMode && m.readMode == 1)
	if activeMode || m.qrButtonFocus == action {
		return "> " + label + " <"
	}
	return "[ " + label + " ]"
}

func qrActionLabel(action qrAction) string {
	switch action {
	case qrActionScrollMode:
		return locale.GetString("tui_mode_scroll")
	case qrActionFlipMode:
		return locale.GetString("tui_mode_flip")
	case qrActionCopyURL:
		return locale.GetString("tui_btn_copy_url")
	case qrActionTerminalReader:
		return locale.GetString("tui_btn_terminal_reader")
	case qrActionOpenBrowser:
		return locale.GetString("tui_btn_open_browser")
	default:
		return ""
	}
}

func qrButtonRows() [][]qrAction {
	return [][]qrAction{
		{qrActionScrollMode, qrActionFlipMode, qrActionCopyURL},
		{qrActionTerminalReader, qrActionOpenBrowser},
	}
}

func qrButtonPosition(action qrAction) (row int, col int, ok bool) {
	for rowIndex, actions := range qrButtonRows() {
		for colIndex, candidate := range actions {
			if candidate == action {
				return rowIndex, colIndex, true
			}
		}
	}
	return 0, 0, false
}

// moveQRButtonHorizontal 在当前按钮行内移动焦点。
func (m *appModel) moveQRButtonHorizontal(delta int) {
	row, col, ok := qrButtonPosition(m.qrButtonFocus)
	if !ok {
		m.qrButtonFocus = qrActionTerminalReader
		return
	}
	actions := qrButtonRows()[row]
	next := col + delta
	if next < 0 {
		next = 0
	}
	if next >= len(actions) {
		next = len(actions) - 1
	}
	m.qrButtonFocus = actions[next]
}

// moveQRButtonVertical 在上下两行按钮之间移动焦点，并尽量保持列位置。
func (m *appModel) moveQRButtonVertical(delta int) {
	row, col, ok := qrButtonPosition(m.qrButtonFocus)
	if !ok {
		m.qrButtonFocus = qrActionTerminalReader
		return
	}
	rows := qrButtonRows()
	nextRow := row + delta
	if nextRow < 0 {
		nextRow = 0
	}
	if nextRow >= len(rows) {
		nextRow = len(rows) - 1
	}
	if col >= len(rows[nextRow]) {
		col = len(rows[nextRow]) - 1
	}
	m.qrButtonFocus = rows[nextRow][col]
}

// moveQRButtonTab 按视觉顺序循环 QRCode 按钮焦点。
func (m *appModel) moveQRButtonTab() {
	order := []qrAction{qrActionScrollMode, qrActionFlipMode, qrActionCopyURL, qrActionTerminalReader, qrActionOpenBrowser}
	for i, action := range order {
		if action == m.qrButtonFocus {
			m.qrButtonFocus = order[(i+1)%len(order)]
			return
		}
	}
	m.qrButtonFocus = qrActionTerminalReader
}

// makePanel 将内容行包装进带边框的面板，聚焦面板使用双线边框，非聚焦使用单线边框。
func (m *appModel) makePanel(title string, content []string, rect panelRect, focused bool) []string {
	if rect.w <= 1 || rect.h <= 1 {
		return []string{""}
	}
	border := singleBorder()
	if focused {
		border = doubleBorder()
	}

	lines := make([]string, 0, rect.h)
	lines = append(lines, border.top(title, rect.w))
	content = fitLines(content, rect.innerWidth(), rect.innerHeight())
	for _, line := range content {
		lines = append(lines, border.middle(line, rect.w))
	}
	for len(lines) < rect.h-1 {
		lines = append(lines, border.middle("", rect.w))
	}
	lines = append(lines, border.bottom(rect.w))
	return lines
}

// boxBorder 面板边框字符集，支持单线和双线两种样式。
type boxBorder struct {
	leftTop     string
	rightTop    string
	leftBottom  string
	rightBottom string
	horizontal  string
	vertical    string
}

// singleBorder 返回单线边框字符（┌─┐│└─┘）。
func singleBorder() boxBorder {
	return boxBorder{
		leftTop:     "┌",
		rightTop:    "┐",
		leftBottom:  "└",
		rightBottom: "┘",
		horizontal:  "─",
		vertical:    "│",
	}
}

// doubleBorder 返回双线边框字符（╔═╗║╚═╝），用于聚焦面板。
func doubleBorder() boxBorder {
	return boxBorder{
		leftTop:     "╔",
		rightTop:    "╗",
		leftBottom:  "╚",
		rightBottom: "╝",
		horizontal:  "═",
		vertical:    "║",
	}
}

// top 渲染面板顶部边框行，标题嵌入在水平线中。
func (b boxBorder) top(title string, width int) string {
	inner := max(0, width-2)
	text := " " + title + " "
	return b.leftTop + padRightWith(text, inner, b.horizontal) + b.rightTop
}

// middle 渲染面板中间内容行，左右各一个竖线边框。
func (b boxBorder) middle(line string, width int) string {
	inner := max(0, width-2)
	return b.vertical + clipAndPad(line, inner) + b.vertical
}

// bottom 渲染面板底部边框行。
func (b boxBorder) bottom(width int) string {
	inner := max(0, width-2)
	return b.leftBottom + strings.Repeat(b.horizontal, inner) + b.rightBottom
}

// fitLines 将内容行裁剪/填充到指定宽高，多余行截断，不足行补空行。
func fitLines(lines []string, width int, height int) []string {
	if height <= 0 {
		return nil
	}
	result := make([]string, 0, height)
	for _, line := range lines {
		if len(result) >= height {
			break
		}
		result = append(result, clipAndPad(line, width))
	}
	for len(result) < height {
		result = append(result, clipAndPad("", width))
	}
	return result
}

// clipAndPad 将文本裁剪到指定显示宽度，不足部分用空格填充，确保每行等宽。
func clipAndPad(text string, width int) string {
	if width <= 0 {
		return ""
	}
	text = strings.ReplaceAll(text, "\t", "    ")
	if runewidth.StringWidth(text) <= width {
		return text + strings.Repeat(" ", width-runewidth.StringWidth(text))
	}

	var builder strings.Builder
	current := 0
	for _, r := range text {
		rw := runewidth.RuneWidth(r)
		if current+rw > width {
			break
		}
		builder.WriteRune(r)
		current += rw
	}
	return builder.String() + strings.Repeat(" ", width-current)
}

// padRightWith 将文本用指定填充字符（如边框线 "─"）补齐到目标宽度。
func padRightWith(text string, width int, pad string) string {
	if width <= 0 {
		return ""
	}
	tw := runewidth.StringWidth(text)
	if tw >= width {
		return clipAndPad(text, width)
	}
	remaining := width - tw
	pw := runewidth.StringWidth(pad)
	if pw <= 0 {
		pw = 1
	}
	count := remaining / pw
	extra := remaining % pw
	return text + strings.Repeat(pad, count) + strings.Repeat(" ", extra)
}

// renderQRCodeLines 将文本编码为 QR 码，并转换为 Unicode 半高块字符行（▀▄█ ），适合终端显示。
func renderQRCodeLines(text string) ([]string, error) {
	return tools.RenderQRCodeLinesTerminal(text)
}

// contains 判断坐标 (x, y) 是否在面板矩形区域内。
func contains(rect panelRect, x int, y int) bool {
	return x >= rect.x && x < rect.x+rect.w && y >= rect.y && y < rect.y+rect.h
}

// centerText 在指定宽度内将文本水平居中（左侧补空格）。
func centerText(text string, width int) string {
	tw := runewidth.StringWidth(text)
	if tw >= width {
		return text
	}
	pad := (width - tw) / 2
	return strings.Repeat(" ", pad) + text
}

// formatTUIActionHint 为操作反馈加两侧标记并居中，便于在 QR 面板中一眼识别（不用 ANSI，避免与 clipAndPad 宽度不一致）。
func formatTUIActionHint(msg string, width int) string {
	if width <= 0 || msg == "" {
		return ""
	}
	const deco = "▶ " // 与右侧 " ◀" 合计 4 列（均为半角符号）
	const decoR = " ◀"
	maxCore := width - runewidth.StringWidth(deco+decoR)
	if maxCore < 1 {
		maxCore = 1
	}
	core := shortenText(msg, maxCore)
	line := deco + core + decoR
	return centerText(line, width)
}

// shortenText 将文本截断到 maxWidth 以内，超长部分以 "…" 替代。
func shortenText(text string, maxWidth int) string {
	if maxWidth <= 0 || runewidth.StringWidth(text) <= maxWidth {
		return text
	}
	if maxWidth <= 1 {
		return "…"
	}
	return strings.TrimRight(clipWidth(text, maxWidth-1), " ") + "…"
}

// clipWidth 按显示宽度（考虑 CJK 全角字符）截取文本，不补空格。
func clipWidth(text string, width int) string {
	if width <= 0 {
		return ""
	}
	var builder strings.Builder
	current := 0
	for _, r := range text {
		rw := runewidth.RuneWidth(r)
		if current+rw > width {
			break
		}
		builder.WriteRune(r)
		current += rw
	}
	return builder.String()
}
