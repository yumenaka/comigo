package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"strings"
	"sync"
	"time"
)

func main() {
	// 1. 初始化自定义的日志缓冲区
	logBuffer := NewLogBuffer()
	// 将标准日志的输出重定向到 logBuffer
	log.SetOutput(logBuffer)

	// 2. 创建 Bubble Tea 程序
	m := initialModel(logBuffer)
	p := tea.NewProgram(m)

	// 3. 运行 TUI 程序
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}

type model struct {
	logBuffer *LogBuffer
	logs      []string
}

func initialModel(lb *LogBuffer) model {
	return model{
		logBuffer: lb,
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

func (m model) Init() tea.Cmd {
	go func() {
		// 模拟一些日志输出
		for i := 0; i < 1000; i++ {
			log.Printf("日志行 %d", i)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 启动的时候立即请求一次 tick
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//switch msg.(type) {
	//下面是类型断言的写法，会把 msg 转换为 tickMsg 类型 或者 tea.KeyMsg 类型
	switch msg := msg.(type) {
	case tickMsg:
		// 获取最新日志并更新 model.logs
		m.logs = m.logBuffer.GetLines()
		// 返回下一个 tickCmd()，以便持续刷新
		return m, tickCmd()
	case tea.KeyMsg:
		// 这里也可以添加各种按键处理，比如退出等
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if len(m.logs) == 0 {
		return "暂无日志\n"
	}

	view := "=== 日志输出 ===\n"
	for _, line := range m.logs {
		view += line + "\n"
	}
	view += "\n(按 Ctrl+C 退出)\n"
	return view
}
