package tui

import (
	"strings"

	termimg "github.com/blacktop/go-termimg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"

	"github.com/yumenaka/comigo/assets/locale"
)

const (
	modalMinWidth             = 32 // 弹窗常规最小宽度，避免短提示看起来过窄。
	modalTitleMessagePadding  = 6  // 标题/正文左右合计预留 6 列，让文字不贴边框。
	modalButtonPadding        = 8  // OK 按钮左右额外预留 8 列，让按钮居中后仍有呼吸空间。
	modalScreenHorizontalGaps = 4  // 弹窗左右至少离屏幕各 2 列，避免贴住终端边缘。
	modalAbsoluteMinWidth     = 12 // 极窄终端下仍保证边框、正文和 OK 按钮可显示。
	modalBorderWidth          = 2  // 单线边框左右各 1 列。
	modalOKPanelRow           = 4  // 当前弹窗布局中 OK 按钮固定在面板第 5 行。
)

// modalState 保存 TUI 内部的轻量提示弹窗；Bubble Tea 本身没有通用 modal 组件，因此在 View 中直接渲染。
type modalState struct {
	Visible    bool
	Title      string
	Message    string
	OKRect     panelRect // OK 按钮的点击区域，View 渲染后写入，鼠标点击时复用。
	NeedsClear bool      // 弹窗首次出现时清屏一次，避免底层终端图片残留到弹窗上。
}

// showModal 打开全局提示弹窗；弹窗会在下一帧接管按键和鼠标事件。
func (m *appModel) showModal(title string, message string) {
	m.modal = modalState{Visible: true, Title: title, Message: message, NeedsClear: true}
}

// closeModal 关闭弹窗，并同步当前界面的图片状态，恢复被弹窗覆盖的阅读页或预览区。
func (m *appModel) closeModal() tea.Cmd {
	m.modal = modalState{}
	return m.syncActiveImageCmd()
}

// handleModalKey 让弹窗优先接管确认/取消键，避免用户按 q 时误退出整个 TUI。
func (m *appModel) handleModalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ", "esc", "q":
		return m, m.closeModal()
	default:
		return m, nil
	}
}

// handleModalMouse 只响应 OK 按钮范围内的左键点击，其它鼠标事件不透传到底层界面。
func (m *appModel) handleModalMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft && contains(m.modal.OKRect, msg.X, msg.Y) {
		return m, m.closeModal()
	}
	return m, nil
}

// renderModalView 用整屏提示替代原画面，并先清理可能残留的终端图片层。
func (m *appModel) renderModalView() string {
	width := max(1, m.renderWidth())
	height := max(1, m.height)
	lines := fitLines(nil, width, height)
	m.modal.OKRect = panelRect{}
	if width <= 0 || height <= 0 {
		return ""
	}

	title := m.modal.Title
	if title == "" {
		title = locale.GetString("tui_modal_title_notice")
	}
	okText := "[ " + locale.GetString("tui_modal_ok") + " ]"
	message := m.modal.Message
	if message == "" {
		message = title
	}

	desiredWidth := max(modalMinWidth, runewidth.StringWidth(title)+modalTitleMessagePadding)
	desiredWidth = max(desiredWidth, runewidth.StringWidth(message)+modalTitleMessagePadding)
	desiredWidth = max(desiredWidth, runewidth.StringWidth(okText)+modalButtonPadding)
	if width > modalScreenHorizontalGaps {
		desiredWidth = min(desiredWidth, width-modalScreenHorizontalGaps)
	}
	modalWidth := max(modalAbsoluteMinWidth, min(width, desiredWidth))
	inner := max(0, modalWidth-modalBorderWidth)
	border := doubleBorder()
	title = shortenText(title, max(0, inner-modalBorderWidth))
	message = shortenText(message, inner)
	okText = shortenText(okText, inner)
	panel := []string{
		border.top(title, modalWidth),
		border.middle("", modalWidth),
		border.middle(centerText(message, inner), modalWidth),
		border.middle("", modalWidth),
		border.middle(centerText(okText, inner), modalWidth),
		border.bottom(modalWidth),
	}
	panel = fitLines(panel, modalWidth, min(len(panel), height))

	startY := max(0, (height-len(panel))/2)
	startX := max(0, (width-modalWidth)/2)
	for i, line := range panel {
		y := startY + i
		if y >= len(lines) {
			break
		}
		lines[y] = clipAndPad(strings.Repeat(" ", startX)+line, width)
	}

	okLine := startY + modalOKPanelRow
	okWidth := runewidth.StringWidth(okText)
	okStart := max(0, (inner-okWidth)/2)
	if okLine < height {
		m.modal.OKRect = panelRect{x: startX + 1 + okStart, y: okLine, w: okWidth, h: 1}
	}
	prefix := ""
	if m.modal.NeedsClear {
		prefix = m.renderModalClearPrefix()
		m.modal.NeedsClear = false
	}
	return prefix + strings.Join(lines, "\n")
}

// renderModalClearPrefix 清掉独立图像层，避免弹窗上方残留上一帧封面或阅读页。
func (m *appModel) renderModalClearPrefix() string {
	var builder strings.Builder
	if m.coverProtocol == termimg.Kitty || m.readerProtocol == termimg.Kitty ||
		m.coverPreview.Protocol == termimg.Kitty || m.terminalReader.Protocol == termimg.Kitty {
		builder.WriteString(termimg.ClearAllString())
	}
	builder.WriteString("\x1b[2J\x1b[H")
	return builder.String()
}
