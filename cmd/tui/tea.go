package tui

import (
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yumenaka/comigo/config"
)

// Examples
// https://github.com/charmbracelet/bubbletea/tree/main/examples

// 存储程序的状态。可以是任何类型，但结构通常是合理的选择
type model struct {
	logBuffer *LogBuffer
	logs      []string
	config    *config.Config // Comigo 配置
}

// 返回我们的初始模型的函数
func InitialModel(lb *LogBuffer) model {
	return model{
		logBuffer: lb,
		config:    config.GetCfg(),
	}
}

type LogBuffer struct {
	lines []string
	mu    sync.RWMutex
}

func NewLogBuffer() *LogBuffer {
	return &LogBuffer{}
}

func (lb *LogBuffer) Write(p []byte) (int, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// 这里简单按行分割，或者也可直接保存为一个完整字符串
	logText := string(p)
	parts := strings.Split(logText, "\n")
	for _, line := range parts {
		if line == "" {
			continue
		}
		lb.lines = append(lb.lines, line)
	}
	return len(p), nil
}

// GetLines 返回当前所有的日志
func (lb *LogBuffer) GetLines() []string {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	// 复制一份，避免并发读写冲突
	result := make([]string, len(lb.lines))
	copy(result, lb.lines)
	return result
}

// 定义一个消息类型，用于周期性刷新日志
type tickMsg time.Time

// 一个命令，用于周期性发送 tickMsg
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/10, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init 用于返回应用程序运行的初始命令的函数。
// 注意，我们并没有调用该函数；Bubble Tea 运行时将在合适的时候进行调用。
// Init 程序启动时如果需要做一些操作，通常就会在Init()方法中返回一个tea.Cmd。tea后台会执行这个函数，最终将返回的tea.Msg传给模型的Update()方法。
func (m model) Init() tea.Cmd {
	// 启动的时候立即请求一次 tick
	return tickCmd()
}

// Update 处理传入事件并相应更新模型的函数
// 事件可以是用户输入、定时器、网络请求等。Bubble Tea 运行时会将这些事件传递给 Update 函数。
// Update 在 “事件发生 " 时被调用。它的工作是查看已经发生的事情并返回一个更新的模型（Model）作为回应。还可以返回一个 Cmd 来使更多的事件发生。

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch msg.(type) {
	// 下面是类型断言的写法，会把 msg 转换为 tickMsg 类型 或者 tea.KeyMsg 类型
	switch msg := msg.(type) {
	case tickMsg:
		// 获取最新日志并更新 model.logs
		m.logs = m.logBuffer.GetLines()
		// 返回下一个 tickCmd()，以便持续刷新
		return m, tickCmd()
	case tea.KeyMsg:
		// 这里也可以添加各种按键处理，比如退出等
		switch msg.String() {
		//  ctrl+c 和 q 返回一个带有模型 (Model) 的 tea.Quit 命令。那是一个特殊的命令，它指示 Bubble Tea 运行时退出，退出程序。
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+j", "j":
			return m, tea.Quit
		}
	}
	// 将更新后的模型（`Model`）返回给Bubble Tea运行时进行处理。
	// 请注意，这里并没有返回一个命令。
	return m, nil
}

// View 根据模型中的数据渲染 UI 的函数
// 所有的方法中，视图是最简单的： 看一下模型的当前状态，然后用它来返回一个字符串。这个字符串就是我们的 UI
// Bubble Tea 的姐妹库 bubbles 提供很多常用组件 https://github.com/charmbracelet/bubbles
// bubbles 的 Viewport 可以用来显示log，并配合 reflow 来实现自动换行 https://github.com/muesli/reflow
// Bubble Tea 可以用 lipgloss 库给文本添加各种颜色 https://github.com/charmbracelet/lipgloss
// Glamour 库可以用来渲染 markdown 文本 https://github.com/charmbracelet/glamour
func (m model) View() string {
	if len(m.logs) == 0 {
		return "暂无日志\n"
	}
	view := "logs：\n"
	// view := "=== 日志输出 ===\n"
	for i, line := range m.logs {
		// 只显示最后 10 行日志
		if i < len(m.logs)-10 {
			continue
		}
		view += line + "\n"
	}
	// view += "\n(按 Ctrl+C 退出)\n"
	return view
}
