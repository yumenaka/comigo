package tui

import (
	"fmt"
	"io"
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
	"github.com/skip2/go-qrcode"
	"golang.org/x/term"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/cmd"
	"github.com/yumenaka/comigo/config"
	modelpkg "github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers"
	wsrouter "github.com/yumenaka/comigo/routers/websocket"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

const (
	logBufferLimit   = 5000
	tickInterval     = 500 * time.Millisecond
	minTUIWidth      = 72
	minTUIHeight     = 20
	layoutGap        = 1
	maxActionMessage = 120
)

type focusPanel int

const (
	focusShelf focusPanel = iota
	focusLog
	focusInfo
	focusQRCode
)

type shelfItemKind int

const (
	shelfItemHeader shelfItemKind = iota
	shelfItemBack
	shelfItemGroup
	shelfItemBook
)

type (
	tickMsg           time.Time
	backendStartedMsg struct{}
	backendErrorMsg   struct {
		err error
	}
)

type openURLResultMsg struct {
	url string
	err error
}

type shelfItem struct {
	Kind       shelfItemKind
	Title      string
	Subtitle   string
	BookID     string
	TargetURL  string
	Selectable bool
}

type shelfLevel struct {
	BookID string
	Title  string
}

type systemSnapshot struct {
	CPUPercent   float64
	RAMPercent   float64
	OnlineUsers  int
	Books        int
	ServerPort   int
	ShelfURL     string
	SelectedText string
	TargetURL    string
	StatusText   string
}

type panelRect struct {
	x int
	y int
	w int
	h int
}

func (r panelRect) innerWidth() int {
	if r.w <= 2 {
		return 0
	}
	return r.w - 2
}

func (r panelRect) innerHeight() int {
	if r.h <= 2 {
		return 0
	}
	return r.h - 2
}

type layoutState struct {
	shelf panelRect
	log   panelRect
	info  panelRect
	qr    panelRect
}

type appModel struct {
	logBuffer *LogBuffer

	width  int
	height int
	focus  focusPanel

	logs           []string
	logOffset      int
	autoFollowLogs bool

	backendReady bool
	backendError string
	actionMsg    string

	stack        []shelfLevel
	items        []shelfItem
	selected     int
	shelfOffset  int // 书架列表滚动偏移（首个可见的 item 索引）
	shelfRowToID map[int]int

	currentShelfURL string
	qrLines         []string
	qrButtonFocus   int // 0 = 浏览器打开, 1 = 复制URL, 2 = 模式切换
	qrButtonRow     int // 操作按钮在 QR 面板内容区的行号
	qrModeRow       int // 模式切换行在 QR 面板内容区的行号
	readMode        int // 0 = scroll（卷轴阅读）, 1 = flip（翻页阅读）
	status          systemSnapshot
}

// LogBuffer 用来缓存 TUI 需要展示的实时日志。
type LogBuffer struct {
	lines []string
	limit int
	mu    sync.RWMutex
}

func NewLogBuffer() *LogBuffer {
	return &LogBuffer{limit: logBufferLimit}
}

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

func InitialModel(lb *LogBuffer) *appModel {
	m := &appModel{
		logBuffer:       lb,
		focus:           focusShelf,
		autoFollowLogs:  true,
		shelfRowToID:    make(map[int]int),
		actionMsg:       locale.GetString("tui_starting_service"),
		currentShelfURL: "",
	}
	m.refreshData()
	return m
}

// Run 启动 TUI 模式；如果没有终端，则退回到普通服务模式。
func Run() error {
	// 现代终端（iTerm2、Terminal.app 等）将 East Asian Ambiguous 字符（含 Box Drawing）渲染为宽度 1，
	// 但 go-runewidth 在 zh_CN 等 CJK locale 下默认将其视为宽度 2，导致面板宽度计算偏差。
	runewidth.DefaultCondition.EastAsianWidth = false

	for _, arg := range os.Args {
		if arg == "-v" || arg == "--version" || arg == "-h" || arg == "--help" {
			cmd.Execute()
			return nil
		}
	}

	if !term.IsTerminal(int(os.Stdout.Fd())) {
		runWithoutTUI()
		return nil
	}

	logBuffer := NewLogBuffer()
	logger.SetSuppressStdout(true)
	logger.SetOutput(io.Discard)
	logger.SetMirrorOutput(logBuffer)
	defer func() {
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

func runWithoutTUI() {
	cmd.Execute()
	routers.StartWebServer()
	routers.StartTailscale()
	cmd.LoadUserPlugins()
	cmd.AddStoreUrls(cmd.Args)
	cmd.SetCwdAsScanPathIfNeed()
	cmd.LoadMetadata()
	cmd.ScanStore()
	cmd.SaveMetadata()
	config.StartOrStopAutoRescan()
	cmd.SetShutdownHandler()
}

func tickCmd() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func startBackendCmd() tea.Cmd {
	return func() (msg tea.Msg) {
		defer func() {
			if recovered := recover(); recovered != nil {
				msg = backendErrorMsg{err: fmt.Errorf("tui 后台启动异常: %v", recovered)}
			}
		}()

		cmd.Execute()
		routers.StartWebServer()
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

func openURLCmd(target string) tea.Cmd {
	return func() tea.Msg {
		err := tools.OpenURL(target)
		return openURLResultMsg{url: target, err: err}
	}
}

func (m *appModel) Init() tea.Cmd {
	return tea.Batch(tickCmd(), startBackendCmd())
}

func (m *appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.refreshData()
		return m, nil
	case tickMsg:
		m.refreshData()
		return m, tickCmd()
	case backendStartedMsg:
		m.backendReady = true
		m.backendError = ""
		m.actionMsg = locale.GetString("tui_service_started")
		m.refreshData()
		return m, nil
	case backendErrorMsg:
		m.backendError = msg.err.Error()
		m.actionMsg = locale.GetString("tui_backend_failed")
		m.refreshData()
		return m, nil
	case openURLResultMsg:
		if msg.err != nil {
			m.actionMsg = shortenText(fmt.Sprintf(locale.GetString("tui_open_browser_failed"), msg.err.Error()), maxActionMessage)
		} else {
			m.actionMsg = shortenText(fmt.Sprintf(locale.GetString("tui_opened_url"), msg.url), maxActionMessage)
		}
		m.refreshStatus()
		return m, nil
	case tea.MouseMsg:
		return m.handleMouse(msg)
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m *appModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "tab":
		m.focus = (m.focus + 1) % 4
		return m, nil
	case "shift+tab":
		if m.focus == focusShelf {
			m.focus = focusQRCode
		} else {
			m.focus--
		}
		return m, nil
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
			if m.qrButtonFocus < 2 {
				m.qrButtonFocus = 2
			}
		case "down", "j":
			if m.qrButtonFocus == 2 {
				m.qrButtonFocus = 0
			}
		case "left", "h":
			if m.qrButtonFocus < 2 {
				m.qrButtonFocus = 0
			}
		case "right", "l":
			if m.qrButtonFocus < 2 {
				m.qrButtonFocus = 1
			}
		case "tab":
			m.qrButtonFocus = (m.qrButtonFocus + 1) % 3
		case "enter", " ":
			if m.qrButtonFocus == 2 {
				m.toggleReadMode()
			} else {
				return m, m.executeQRButton()
			}
		}
	}

	m.refreshStatus()
	return m, nil
}

func (m *appModel) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	layout := m.layout()
	if msg.Action != tea.MouseActionPress {
		return m, nil
	}

	switch {
	case contains(layout.shelf, msg.X, msg.Y):
		m.focus = focusShelf
		if msg.Button == tea.MouseButtonWheelUp {
			m.moveSelection(-1)
			return m, nil
		}
		if msg.Button == tea.MouseButtonWheelDown {
			m.moveSelection(1)
			return m, nil
		}
		if msg.Button == tea.MouseButtonLeft {
			row := msg.Y - layout.shelf.y - 1
			if idx, ok := m.shelfRowToID[row]; ok {
				m.selected = idx
				m.syncShelfOffset()
				m.refreshQRCode()
				m.refreshStatus()
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
	case contains(layout.info, msg.X, msg.Y):
		m.focus = focusInfo
	case contains(layout.qr, msg.X, msg.Y):
		m.focus = focusQRCode
		if msg.Button == tea.MouseButtonLeft {
			clickRow := msg.Y - layout.qr.y - 1
			innerW := layout.qr.innerWidth()
			col := msg.X - layout.qr.x - 1
			if clickRow == m.qrModeRow {
				m.qrButtonFocus = 2
				scrollText := locale.GetString("tui_mode_scroll")
				btn0Str := "[ " + scrollText + " ]"
				if m.readMode == 0 {
					btn0Str = "> " + scrollText + " <"
				}
				modeFull := btn0Str + "  "
				modeFullW := runewidth.StringWidth(modeFull)
				allText := modeFull
				flipText := locale.GetString("tui_mode_flip")
				if m.readMode == 1 {
					allText += "> " + flipText + " <"
				} else {
					allText += "[ " + flipText + " ]"
				}
				allW := runewidth.StringWidth(allText)
				offset := 0
				if allW < innerW {
					offset = (innerW - allW) / 2
				}
				if col < offset+modeFullW {
					if m.readMode != 0 {
						m.toggleReadMode()
					}
				} else {
					if m.readMode != 1 {
						m.toggleReadMode()
					}
				}
				return m, nil
			}
			if clickRow == m.qrButtonRow {
				btn0Text := "[ " + locale.GetString("tui_btn_open_browser") + " ]"
				btn1Text := "[ " + locale.GetString("tui_btn_copy_url") + " ]"
				btnLine := btn0Text + "  " + btn1Text
				btnW := runewidth.StringWidth(btnLine)
				offset := 0
				if btnW < innerW {
					offset = (innerW - btnW) / 2
				}
				btn0End := offset + runewidth.StringWidth(btn0Text)
				if col >= offset && col < btn0End {
					m.qrButtonFocus = 0
				} else {
					m.qrButtonFocus = 1
				}
				return m, m.executeQRButton()
			}
		}
	}
	m.refreshStatus()
	return m, nil
}

func (m *appModel) View() string {
	if m.width < minTUIWidth || m.height < minTUIHeight {
		return fmt.Sprintf(locale.GetString("tui_terminal_too_small"), minTUIWidth, minTUIHeight) + "\n" +
			fmt.Sprintf(locale.GetString("tui_current_size"), m.width, m.height) + "\n"
	}

	layout := m.layout()
	shelfLines := m.makePanel(locale.GetString("tui_panel_shelf"), m.renderShelfContent(layout.shelf), layout.shelf, m.focus == focusShelf)
	logLines := m.makePanel(locale.GetString("tui_panel_log"), m.renderLogContent(layout.log), layout.log, m.focus == focusLog)
	infoLines := m.makePanel(locale.GetString("tui_panel_info"), m.renderInfoContent(layout.info), layout.info, m.focus == focusInfo)
	qrLines := m.makePanel("QRCode", m.renderQRCodeContent(layout.qr), layout.qr, m.focus == focusQRCode)

	topRow := mergeRows(shelfLines, layout.shelf.w, qrLines)
	bottomRow := mergeRows(logLines, layout.log.w, infoLines)
	return strings.Join(append(topRow, bottomRow...), "\n")
}

func (m *appModel) refreshData() {
	m.logs = m.logBuffer.GetLines()
	m.syncLogOffset()
	m.refreshShelf()
	m.refreshStatus()
	m.refreshQRCode()
}

func (m *appModel) refreshShelf() {
	m.items = m.buildCurrentShelfItems()
	m.ensureSelectedItem()
	m.currentShelfURL = m.buildCurrentShelfURL()
}

func (m *appModel) refreshStatus() {
	selectedText := locale.GetString("tui_info_none")
	targetURL := ""
	if item := m.currentItem(); item != nil {
		selectedText = item.Title
		targetURL = item.TargetURL
	}

	statusText := locale.GetString("tui_status_starting")
	if m.backendError != "" {
		statusText = locale.GetString("tui_status_failed")
	} else if m.backendReady {
		statusText = locale.GetString("tui_status_running")
	}
	if m.actionMsg != "" {
		statusText = m.actionMsg
	}

	bookCount := 0
	if modelpkg.IStore != nil {
		bookCount = modelpkg.GetAllBooksNumber()
	}

	cfg := config.GetCfg()
	m.status = systemSnapshot{
		CPUPercent:   0,
		RAMPercent:   0,
		OnlineUsers:  wsrouter.ClientCount(),
		Books:        bookCount,
		ServerPort:   cfg.Port,
		ShelfURL:     m.currentShelfURL,
		SelectedText: selectedText,
		TargetURL:    targetURL,
		StatusText:   statusText,
	}
	sys := tools.GetSystemStatus()
	m.status.CPUPercent = sys.CPUUsedPercent
	m.status.RAMPercent = sys.MemoryUsedPercent
}

// selectedURL 返回当前选中项的 URL；若无可选项则回退到当前书架 URL。
func (m *appModel) selectedURL() string {
	if item := m.currentItem(); item != nil && item.TargetURL != "" {
		return item.TargetURL
	}
	return m.currentShelfURL
}

func (m *appModel) refreshQRCode() {
	lines, err := renderQRCodeLines(m.selectedURL())
	if err != nil {
		m.qrLines = []string{locale.GetString("tui_qr_gen_failed"), err.Error()}
		return
	}
	m.qrLines = lines
}

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

func buildBookTargetURL(book modelpkg.BookInfo, readMode int) string {
	baseURL := buildBaseURL()
	base := strings.TrimRight(baseURL, "/")
	switch book.Type {
	case modelpkg.TypeBooksGroup:
		return base + "/shelf/" + book.BookID
	case modelpkg.TypeVideo, modelpkg.TypeAudio:
		return base + "/player/" + book.BookID
	case modelpkg.TypeHTML, modelpkg.TypeUnknownFile:
		return base + "/api/raw/" + book.BookID + "/" + url.QueryEscape(book.Title)
	default:
		prefix := "/scroll/"
		if readMode == 1 {
			prefix = "/flip/"
		}
		target := base + prefix + book.BookID
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

func buildBaseURL() string {
	cfg := config.GetCfg()
	protocol := "http://"
	if (cfg.CertFile != "" && cfg.KeyFile != "") || cfg.AutoTLSCertificate {
		protocol = "https://"
	}

	host := cfg.Host
	if host == "" {
		if ipList, err := tools.GetIPList(); err == nil && len(ipList) > 0 {
			host = ipList[0]
		} else {
			host = "127.0.0.1"
		}
	}
	if cfg.AutoTLSCertificate {
		return protocol + host
	}
	return fmt.Sprintf("%s%s:%d", protocol, host, cfg.Port)
}

func (m *appModel) buildCurrentShelfURL() string {
	baseURL := strings.TrimRight(buildBaseURL(), "/")
	if len(m.stack) == 0 {
		return baseURL
	}
	return baseURL + "/shelf/" + m.stack[len(m.stack)-1].BookID
}

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

func (m *appModel) currentItem() *shelfItem {
	if len(m.items) == 0 || m.selected < 0 || m.selected >= len(m.items) {
		return nil
	}
	if !m.items[m.selected].Selectable {
		return nil
	}
	return &m.items[m.selected]
}

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

func (m *appModel) selectFirst() {
	for i := range m.items {
		if m.items[i].Selectable {
			m.selected = i
			m.syncShelfOffset()
			return
		}
	}
}

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

func (m *appModel) goBack() {
	if len(m.stack) == 0 {
		return
	}
	m.stack = m.stack[:len(m.stack)-1]
	m.selected = 0
	m.shelfOffset = 0
	m.refreshData()
}

func (m *appModel) activateSelectedItem() tea.Cmd {
	item := m.currentItem()
	if item == nil {
		return nil
	}

	switch item.Kind {
	case shelfItemBack:
		m.goBack()
		return nil
	case shelfItemGroup:
		m.stack = append(m.stack, shelfLevel{
			BookID: item.BookID,
			Title:  item.Title,
		})
		m.selected = 0
		m.shelfOffset = 0
		m.actionMsg = fmt.Sprintf(locale.GetString("tui_entered_sub_shelf"), item.Title)
		m.refreshData()
		return nil
	case shelfItemBook:
		m.actionMsg = fmt.Sprintf(locale.GetString("tui_opening_url"), item.TargetURL)
		m.refreshStatus()
		return openURLCmd(item.TargetURL)
	default:
		return nil
	}
}

// executeQRButton 执行 QR 面板当前聚焦按钮的动作。
func (m *appModel) executeQRButton() tea.Cmd {
	target := m.selectedURL()
	if target == "" {
		m.actionMsg = locale.GetString("tui_no_url_available")
		m.refreshStatus()
		return nil
	}
	switch m.qrButtonFocus {
	case 0:
		m.actionMsg = fmt.Sprintf(locale.GetString("tui_opening_url"), shortenText(target, maxActionMessage-6))
		m.refreshStatus()
		return openURLCmd(target)
	case 1:
		if err := clipboard.WriteAll(target); err != nil {
			m.actionMsg = fmt.Sprintf(locale.GetString("tui_copy_failed"), err.Error())
		} else {
			m.actionMsg = fmt.Sprintf(locale.GetString("tui_url_copied"), shortenText(target, maxActionMessage-8))
		}
		m.refreshStatus()
		return nil
	default:
		return nil
	}
}

// toggleReadMode 在卷轴阅读和翻页阅读之间切换，并刷新所有受影响的数据。
func (m *appModel) toggleReadMode() {
	if m.readMode == 0 {
		m.readMode = 1
	} else {
		m.readMode = 0
	}
	m.refreshShelf()
	m.refreshQRCode()
	m.refreshStatus()
}

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

func (m *appModel) layout() layoutState {
	width := max(0, m.width)
	height := max(0, m.height)
	leftWidth := (width - layoutGap) * 2 / 3
	rightWidth := width - layoutGap - leftWidth
	topHeight := (height - layoutGap) * 2 / 3
	bottomHeight := height - layoutGap - topHeight

	return layoutState{
		shelf: panelRect{x: 0, y: 0, w: leftWidth, h: topHeight},
		qr:    panelRect{x: leftWidth + layoutGap, y: 0, w: rightWidth, h: topHeight},
		log:   panelRect{x: 0, y: topHeight + layoutGap, w: leftWidth, h: bottomHeight},
		info:  panelRect{x: leftWidth + layoutGap, y: topHeight + layoutGap, w: rightWidth, h: bottomHeight},
	}
}

func (m *appModel) renderShelfContent(rect panelRect) []string {
	lines := make([]string, 0, rect.innerHeight())
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

	visibleItems := rect.innerHeight() - len(lines)
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

	if len(lines) == 0 {
		return []string{locale.GetString("tui_no_shelf_content")}
	}
	return fitLines(lines, rect.innerWidth(), rect.innerHeight())
}

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

func (m *appModel) renderLogContent(rect panelRect) []string {
	height := rect.innerHeight()
	if height <= 0 {
		return nil
	}
	if len(m.logs) == 0 {
		return fitLines([]string{locale.GetString("tui_no_logs")}, rect.innerWidth(), height)
	}

	start := min(m.logOffset, len(m.logs))
	end := min(len(m.logs), start+height)
	lines := append([]string(nil), m.logs[start:end]...)
	if !m.autoFollowLogs {
		lines = append(lines, fmt.Sprintf(locale.GetString("tui_log_scrolling"), end, len(m.logs)))
	}
	return fitLines(lines, rect.innerWidth(), height)
}

func (m *appModel) renderInfoContent(rect panelRect) []string {
	w := rect.innerWidth()
	h := rect.innerHeight()
	lines := []string{
		fmt.Sprintf(locale.GetString("tui_info_service_status"), m.status.StatusText),
		fmt.Sprintf(locale.GetString("tui_info_cpu"), m.status.CPUPercent),
		fmt.Sprintf(locale.GetString("tui_info_ram"), m.status.RAMPercent),
		fmt.Sprintf(locale.GetString("tui_info_online_users"), m.status.OnlineUsers),
		fmt.Sprintf(locale.GetString("tui_info_books_total"), m.status.Books),
		fmt.Sprintf(locale.GetString("tui_info_server_port"), m.status.ServerPort),
		locale.GetString("tui_info_shelf_url"),
		m.status.ShelfURL,
		fmt.Sprintf(locale.GetString("tui_info_selected"), m.status.SelectedText),
	}
	if m.status.TargetURL != "" {
		lines = append(lines, locale.GetString("tui_info_target_url"), m.status.TargetURL)
	}

	// 底部右对齐：时间 + 版本
	versionLine := time.Now().Format("2006-01-02 15:04:05") + "  Comigo " + config.GetVersion()
	if h > len(lines)+1 {
		for len(lines) < h-1 {
			lines = append(lines, "")
		}
		vw := runewidth.StringWidth(versionLine)
		if vw < w {
			versionLine = strings.Repeat(" ", w-vw) + versionLine
		}
		lines = append(lines, versionLine)
	}
	return fitLines(lines, w, h)
}

func (m *appModel) renderQRCodeContent(rect panelRect) []string {
	w := rect.innerWidth()
	selURL := m.selectedURL()
	label := locale.GetString("tui_qr_shelf_url")
	if item := m.currentItem(); item != nil && item.TargetURL != "" {
		label = fmt.Sprintf(locale.GetString("tui_qr_selected"), shortenText(item.Title, max(10, w-6)))
	}
	lines := []string{label, selURL}
	if len(m.qrLines) == 0 {
		lines = append(lines, centerText(locale.GetString("tui_qr_unavailable"), w))
	} else {
		lines = append(lines, "")
		for _, qrLine := range m.qrLines {
			lines = append(lines, centerText(qrLine, w))
		}
	}

	// 模式切换行
	lines = append(lines, "")
	scrollText := locale.GetString("tui_mode_scroll")
	flipText := locale.GetString("tui_mode_flip")
	modeBtn0 := "[ " + scrollText + " ]"
	modeBtn1 := "[ " + flipText + " ]"
	if m.readMode == 0 {
		modeBtn0 = "> " + scrollText + " <"
	} else {
		modeBtn1 = "> " + flipText + " <"
	}
	modeLine := modeBtn0 + "  " + modeBtn1
	m.qrModeRow = len(lines)
	lines = append(lines, centerText(modeLine, w))

	// 操作按钮行
	lines = append(lines, "")
	btnOpenText := locale.GetString("tui_btn_open_browser")
	btnCopyText := locale.GetString("tui_btn_copy_url")
	btn0 := "[ " + btnOpenText + " ]"
	btn1 := "[ " + btnCopyText + " ]"
	btnLine := btn0 + "  " + btn1
	m.qrButtonRow = len(lines)
	lines = append(lines, centerText(btnLine, w))

	return fitLines(lines, w, rect.innerHeight())
}

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

type boxBorder struct {
	leftTop     string
	rightTop    string
	leftBottom  string
	rightBottom string
	horizontal  string
	vertical    string
}

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

func (b boxBorder) top(title string, width int) string {
	inner := max(0, width-2)
	text := " " + title + " "
	return b.leftTop + padRightWith(text, inner, b.horizontal) + b.rightTop
}

func (b boxBorder) middle(line string, width int) string {
	inner := max(0, width-2)
	return b.vertical + clipAndPad(line, inner) + b.vertical
}

func (b boxBorder) bottom(width int) string {
	inner := max(0, width-2)
	return b.leftBottom + strings.Repeat(b.horizontal, inner) + b.rightBottom
}

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

func mergeRows(left []string, leftWidth int, right []string) []string {
	count := max(len(left), len(right))
	rows := make([]string, 0, count)
	for i := 0; i < count; i++ {
		l := ""
		r := ""
		if i < len(left) {
			l = left[i]
		}
		if i < len(right) {
			r = right[i]
		}
		l = clipAndPad(l, leftWidth)
		rows = append(rows, l+strings.Repeat(" ", layoutGap)+r)
	}
	return rows
}

func renderQRCodeLines(text string) ([]string, error) {
	if strings.TrimSpace(text) == "" {
		return []string{locale.GetString("tui_qr_unavailable")}, nil
	}
	qr, err := qrcode.New(text, qrcode.Low)
	if err != nil {
		return nil, err
	}
	bitmap := qr.Bitmap()
	if len(bitmap)%2 != 0 {
		padding := make([]bool, len(bitmap[0]))
		bitmap = append(bitmap, padding)
	}

	lines := make([]string, 0, len(bitmap)/2)
	for row := 0; row < len(bitmap); row += 2 {
		var builder strings.Builder
		for col := 0; col < len(bitmap[row]); col++ {
			top := bitmap[row][col]
			bottom := bitmap[row+1][col]
			builder.WriteRune(qrBlock(top, bottom))
		}
		lines = append(lines, builder.String())
	}
	return lines, nil
}

func qrBlock(top bool, bottom bool) rune {
	switch {
	case top && bottom:
		return '█'
	case top && !bottom:
		return '▀'
	case !top && bottom:
		return '▄'
	default:
		return ' '
	}
}

func contains(rect panelRect, x int, y int) bool {
	return x >= rect.x && x < rect.x+rect.w && y >= rect.y && y < rect.y+rect.h
}

func centerText(text string, width int) string {
	tw := runewidth.StringWidth(text)
	if tw >= width {
		return text
	}
	pad := (width - tw) / 2
	return strings.Repeat(" ", pad) + text
}

func shortenText(text string, maxWidth int) string {
	if maxWidth <= 0 || runewidth.StringWidth(text) <= maxWidth {
		return text
	}
	if maxWidth <= 1 {
		return "…"
	}
	return strings.TrimRight(clipWidth(text, maxWidth-1), " ") + "…"
}

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
