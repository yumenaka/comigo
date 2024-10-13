package state

type ReaderMode int

const (
	// ScrollMode 卷轴模式
	ScrollMode ReaderMode = iota
	// FlipMode 翻页模式
	FlipMode
)

type UserState struct {
	IsLogin    bool
	IsAdmin    bool
	ReaderMode ReaderMode //阅读器默认阅读模式
}

var User UserState
