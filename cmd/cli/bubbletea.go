package main

// 参考视频：
// https://youtu.be/Gl31diSVP8M
// go run bubbletea.go input.go
import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// 样式
type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	// 将边框设置为正常边框，而不是粗体边框。同时设置边框颜色。间距设置为1，宽度设置为80。
	s.InputField = lipgloss.NewStyle().
		Bold(true).
		Foreground(s.BorderColor).
		Background(lipgloss.Color("#7D56F4")).
		//当前版本的边框似乎有BUG，设置为隐藏边框。
		BorderStyle(lipgloss.HiddenBorder()).
		Padding(1).
		Width(40)
	//s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

// Main 模型 用来存储 tui 的状态
type Main struct {
	styles *Styles
	index  int
	books  []Book
	width  int // 窗口宽度
	height int // 窗口高度
	done   bool
}
type Book struct {
	title       string
	description string
	input       Input
}

func NewBook(title string) Book {
	return Book{title: title}
}

func NewBookWithShortDescription(title string) Book {
	book := NewBook(title)
	model := NewShortDescriptionField()
	book.input = model
	return book
}

func NewBookWithLongDescription(title string) Book {
	book := NewBook(title)
	model := NewLongDescriptionField()
	book.input = model
	return book
}

// New 创建一个新的模型
func New(books []Book) *Main {
	styles := DefaultStyles()
	return &Main{styles: styles, books: books}
}

func (m Main) Init() tea.Cmd {
	return m.books[m.index].input.Blink
}

// Update 响应操作，更新模型的地方
func (m Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	current := &m.books[m.index]
	var cmd tea.Cmd
	switch msg := msg.(type) {
	// 当窗口大小改变时，更新模型的宽度和高度。一般只在初始化的时候赋值一次。
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		// 按下ctrl+c或q，退出程序
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.index == len(m.books)-1 {
				m.done = true
			}
			// 将输入框的值赋值给当前的书籍的描述
			// Value() 返回输入框的值，需要在input包装器当中实现
			current.description = current.input.Value()
			// 打印当前书籍的标题和描述
			log.Printf("title: %s, description: %s", current.title, current.description)
			// 下一本书
			m.Next()
			// Blur()删除焦点
			return m, current.input.Blur
		}
	}
	current.input, cmd = current.input.Update(msg)
	return m, cmd
}

// View 渲染模型，显示在终端上
func (m Main) View() string {
	current := m.books[m.index]
	if m.done {
		var output string
		for _, d := range m.books {
			output += fmt.Sprintf("title: %s, description: %s\n", d.title, d.description)
		}
		return output
	}
	// 如果窗口大小为0，说明还没有初始化完成
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}
	// 如果窗口太小，显示提示信息
	if m.width < 20 || m.height < 10 {
		return "Window too small!"
	}
	// lipgloss 是 bubbletea 的布局工具
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.books[m.index].title,
			//应用 model 自身定义的样式
			m.styles.InputField.Render(current.input.View()),
		),
	)
}
func (m *Main) Next() {
	if m.index == len(m.books)-1 {
		m.index = 0
	} else {
		m.index = m.index + 1
	}
}
func main() {
	books := []Book{
		NewBookWithShortDescription("无所事事的世界中心小姐"),
		NewBookWithShortDescription("颠扑不破"),
		NewBookWithLongDescription("人间失格"),
		//NewBook("逗你玩"),
		//NewBook("解忧杂货店"),

	}
	main := New(books)
	//LogToFile 设置默认日志记录以记录到文件。//这很有必要，我们无法打印到终端，因为我们的 TUI 正在占用它。//如果该文件不存在，则会创建该文件。使用完毕后不要忘记关闭该文件。
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("could not create log file: %v", err)
	}
	defer f.Close()
	//WithAltScreen 在启用备用屏幕缓冲区的情况下启动程序（即程序以全窗口模式启动）。请注意，当程序退出时，Alt屏幕将自动退出。
	//在终端环境中，"Alt Screen"（或称为"Alternate Screen"）是一个特殊的终端模式，它允许应用程序使用两个独立的屏幕缓冲区：一个是主屏幕（Main Screen），另一个是备用屏幕（Alternate Screen）。
	//这个功能用于那些需要临时清除屏幕以显示信息，但完成后又希望恢复原始屏幕内容的情况(vim就是这样的)。 这个功能依赖于终端仿真器的支持，并不是所有的终端都支持。
	p := tea.NewProgram(*main, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("could not start program: %v", err)
	}
}
