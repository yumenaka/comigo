package main

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
	title string
	//阅读器默认阅读模式
	readerMode ReaderMode
	//是否全屏
	windowFullScreen bool
	//是否有边框和标题栏
	windowDecorated bool
	//是否允许调整窗口大小
	windowResizingModeEnabled ebiten.WindowResizingModeType
	//窗口宽度
	width int
	//窗口高度
	height int
	//ebiten运行选项
	runOptions ebiten.RunGameOptions
}

//// 获取各种值的方法

// Title 返回阅读器标题
func (rc *ReaderConfig) Title() string {
	return rc.title
}

// ReaderMode 返回阅读器模式
func (rc *ReaderConfig) ReaderMode() ReaderMode {
	return rc.readerMode
}

// WindowFullScreen 返回是否全屏
func (rc *ReaderConfig) WindowFullScreen() bool {
	return rc.windowFullScreen
}

// WindowDecorated 返回是否有边框和标题栏
func (rc *ReaderConfig) WindowDecorated() bool {
	return rc.windowDecorated
}

// WindowResizingModeEnabled 返回是否允许调整窗口大小
func (rc *ReaderConfig) WindowResizingModeEnabled() ebiten.WindowResizingModeType {
	return rc.windowResizingModeEnabled
}

// Width 返回窗口宽度
func (rc *ReaderConfig) Width() int {
	return rc.width
}

// Height 返回窗口高度
func (rc *ReaderConfig) Height() int {
	return rc.height
}

// RunOptions 返回ebiten运行选项
func (rc *ReaderConfig) RunOptions() *ebiten.RunGameOptions {
	return &rc.runOptions
}

//// 设置各种值的方法

// SetTitle 设置阅读器标题
func (rc *ReaderConfig) SetTitle(title string) *ReaderConfig {
	rc.title = title
	return rc
}

// SetReaderMode 设置阅读器模式
func (rc *ReaderConfig) SetReaderMode(readerMode ReaderMode) *ReaderConfig {
	rc.readerMode = readerMode
	return rc
}

// SetWindowFullScreen 设置是否全屏
func (rc *ReaderConfig) SetWindowFullScreen(windowFullScreen bool) *ReaderConfig {
	rc.windowFullScreen = windowFullScreen
	return rc
}

// SetWindowDecorated 设置是否有边框和标题栏
func (rc *ReaderConfig) SetWindowDecorated(windowDecorated bool) *ReaderConfig {
	rc.windowDecorated = windowDecorated
	return rc
}

// SetWindowResizingModeEnabled 设置是否允许调整窗口大小
func (rc *ReaderConfig) SetWindowResizingModeEnabled(windowResizingModeEnabled ebiten.WindowResizingModeType) *ReaderConfig {
	rc.windowResizingModeEnabled = windowResizingModeEnabled
	return rc
}

// SetWidth 设置窗口宽度
func (rc *ReaderConfig) SetWidth(width int) *ReaderConfig {
	rc.width = width
	return rc
}

// SetHeight 设置窗口高度
func (rc *ReaderConfig) SetHeight(height int) *ReaderConfig {
	rc.height = height
	return rc
}

func (rc *ReaderConfig) SetWindowSize(width, height int) *ReaderConfig {
	rc.width = width
	rc.height = height
	return rc
}

// SetRunOptions 设置ebiten运行选项
func (rc *ReaderConfig) SetRunOptions(runOptions ebiten.RunGameOptions) *ReaderConfig {
	rc.runOptions = runOptions
	return rc
}

// NewReaderConfig 创建一个新的阅读器设定
func NewReaderConfig() *ReaderConfig {
	return &ReaderConfig{
		readerMode:                ScrollMode,
		windowFullScreen:          false,
		windowDecorated:           true,
		windowResizingModeEnabled: ebiten.WindowResizingModeEnabled,
		width:                     900,
		height:                    800,
		runOptions:                ebiten.RunGameOptions{},
	}
}
