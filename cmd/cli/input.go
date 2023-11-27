package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Input interface {
	Blink() tea.Msg                  //初始化光标闪烁的命令
	Blur() tea.Msg                   //失去焦点
	Focus() tea.Cmd                  //获取焦点
	SetValue(string)                 //设置值
	Value() string                   //获取值
	Update(tea.Msg) (Input, tea.Cmd) //更新
	View() string                    //渲染
}

// ShortDescriptionField 短描述输入框 只能输入一行
type ShortDescriptionField struct {
	textinput textinput.Model
}

// NewShortDescriptionField 创建一个新的短描述输入框
func NewShortDescriptionField() *ShortDescriptionField {
	d := ShortDescriptionField{}
	model := textinput.New()
	model.Placeholder = "please input your description"
	// 获取焦点
	model.Focus()
	d.textinput = model
	return &d
}

/* 实现 Input 接口 */

// Blink 是用于初始化光标闪烁的命令。注意，这里是textinput.Blink()，而不是d.textinput.Blink()！
func (d *ShortDescriptionField) Blink() tea.Msg {
	return textinput.Blink() //Blink 是用于初始化光标闪烁的命令。注意，这里是textinput.Blink()，而不是d.textinput.Blink()！
}

func (d *ShortDescriptionField) Init() tea.Cmd {
	return nil
}

func (sd *ShortDescriptionField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	sd.textinput, cmd = sd.textinput.Update(msg)
	return sd, cmd
}

func (sd *ShortDescriptionField) View() string {
	return sd.textinput.View()
}

func (sd *ShortDescriptionField) Focus() tea.Cmd {
	return sd.textinput.Focus()
}

func (sd *ShortDescriptionField) SetValue(s string) {
	sd.textinput.SetValue(s)
}
func (sd *ShortDescriptionField) Blur() tea.Msg {
	return sd.textinput.Blur
}
func (sd *ShortDescriptionField) Value() string {
	return sd.textinput.Value()
}

// LongDescriptionField 长描述输入框 可以输入多行
type LongDescriptionField struct {
	textarea textarea.Model
}

// NewLongDescriptionField 创建一个新的长描述输入框
func NewLongDescriptionField() *LongDescriptionField {
	d := LongDescriptionField{}
	ta := textarea.New()
	ta.Placeholder = "please input your description"
	// 获取焦点
	ta.Focus()
	d.textarea = ta
	return &d
}

/* 实现 Input 接口 */

// Blink 是用于初始化光标闪烁的命令。注意，这里是textinput.Blink()，而不是d.textinput.Blink()！
func (ld *LongDescriptionField) Blink() tea.Msg {
	return textarea.Blink() //Blink 是用于初始化光标闪烁的命令。注意，这里是textinput.Blink()，而不是d.textinput.Blink()！
}
func (ld *LongDescriptionField) Init() tea.Cmd {
	return nil
}
func (ld *LongDescriptionField) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	ld.textarea, cmd = ld.textarea.Update(msg)
	return ld, cmd
}
func (ld *LongDescriptionField) View() string {
	return ld.textarea.View()
}
func (ld *LongDescriptionField) Focus() tea.Cmd {
	return ld.textarea.Focus()
}
func (ld *LongDescriptionField) SetValue(s string) {
	ld.textarea.SetValue(s)
}
func (ld *LongDescriptionField) Blur() tea.Msg {
	return ld.textarea.Blur
}
func (ld *LongDescriptionField) Value() string {
	return ld.textarea.Value()
}
