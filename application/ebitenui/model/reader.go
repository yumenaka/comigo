package model

import "github.com/hajimehoshi/ebiten/v2"

type ReaderMode int

const (
	// ScrollMode 卷轴模式
	ScrollMode ReaderMode = iota
	// FlipMode 翻页模式
	FlipMode
)

// ReaderConfig 阅读器设定
type ReaderConfig struct {
	Title string
	//阅读器默认阅读模式
	ReaderMode ReaderMode
	Router     func(to string)
	//是否全屏
	WindowFullScreen bool
	//是否有边框和标题栏
	WindowDecorated bool
	//是否允许调整窗口大小
	WindowResizingModeEnabled ebiten.WindowResizingModeType
	//窗口宽度
	Width int
	//窗口高度
	Height int
	//ebiten运行选项
	RunOptions ebiten.RunGameOptions
}

// SetTitle 设置阅读器标题
func (rc *ReaderConfig) SetTitle(title string) *ReaderConfig {
	rc.Title = title
	return rc
}

// SetReaderMode 设置阅读器模式
func (rc *ReaderConfig) SetReaderMode(readerMode ReaderMode) *ReaderConfig {
	rc.ReaderMode = readerMode
	return rc
}

// SetWindowFullScreen 设置是否全屏
func (rc *ReaderConfig) SetWindowFullScreen(windowFullScreen bool) *ReaderConfig {
	rc.WindowFullScreen = windowFullScreen
	return rc
}

// SetWindowDecorated 设置是否有边框和标题栏
func (rc *ReaderConfig) SetWindowDecorated(windowDecorated bool) *ReaderConfig {
	rc.WindowDecorated = windowDecorated
	return rc
}

// SetWindowResizingModeEnabled 设置是否允许调整窗口大小
func (rc *ReaderConfig) SetWindowResizingModeEnabled(windowResizingModeEnabled ebiten.WindowResizingModeType) *ReaderConfig {
	rc.WindowResizingModeEnabled = windowResizingModeEnabled
	return rc
}

// SetWidth 设置窗口宽度
func (rc *ReaderConfig) SetWidth(width int) *ReaderConfig {
	rc.Width = width
	return rc
}

// SetHeight 设置窗口高度
func (rc *ReaderConfig) SetHeight(height int) *ReaderConfig {
	rc.Height = height
	return rc
}

func (rc *ReaderConfig) SetWindowSize(width, height int) *ReaderConfig {
	rc.Width = width
	rc.Height = height
	return rc
}

// SetRunOptions 设置ebiten运行选项
func (rc *ReaderConfig) SetRunOptions(runOptions ebiten.RunGameOptions) *ReaderConfig {
	rc.RunOptions = runOptions
	return rc
}

// NewReaderConfig 创建一个新的阅读器设定
func NewReaderConfig() *ReaderConfig {
	return &ReaderConfig{
		ReaderMode:                ScrollMode,
		WindowFullScreen:          false,
		WindowDecorated:           true,
		WindowResizingModeEnabled: ebiten.WindowResizingModeEnabled,
		Width:                     900,
		Height:                    800,
		RunOptions:                ebiten.RunGameOptions{},
	}
}
