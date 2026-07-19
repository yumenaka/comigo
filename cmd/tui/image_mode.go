package tui

import (
	"fmt"

	termimg "github.com/blacktop/go-termimg"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/yumenaka/comigo/assets/locale"
)

// toggleTUIImageMode 在低兼容的 ANSI/Halfblocks 和终端原生图片协议之间切换。
func (m *appModel) toggleTUIImageMode() tea.Cmd {
	wasKitty := m.coverProtocol == termimg.Kitty || m.readerProtocol == termimg.Kitty
	wasOverlay := m.activeImageUsesOverlay()
	if m.isANSIImageMode() {
		coverProtocol := detectNativeTUIImageProtocol()
		readerProtocol := detectNativeTUIReaderImageProtocol()
		if !isNativeTUIImageProtocol(coverProtocol) || !isNativeTUIImageProtocol(readerProtocol) {
			m.showModal(locale.GetString("tui_modal_title_notice"), locale.GetString("tui_image_mode_incompatible"))
			return clearTUIScreenCmd(nil)
		}
		m.coverProtocol = coverProtocol
		m.readerProtocol = readerProtocol
		m.setActionMsg(fmt.Sprintf(locale.GetString("tui_image_mode_native_enabled"), protocolName(readerProtocol)))
	} else {
		m.coverProtocol = termimg.Halfblocks
		m.readerProtocol = termimg.Halfblocks
		m.setActionMsg(locale.GetString("tui_image_mode_ansi_enabled"))
	}
	if wasKitty || m.coverProtocol == termimg.Kitty || m.readerProtocol == termimg.Kitty {
		m.clearKittyImagesNextFrame = true
	}
	m.invalidateTUIImageState()
	m.refreshStatus()
	syncCmd := m.syncActiveImageCmd()
	if wasOverlay {
		return clearTUIScreenCmd(syncCmd)
	}
	return syncCmd
}

func (m *appModel) isANSIImageMode() bool {
	return isTextTUIImageProtocol(m.coverProtocol) && isTextTUIImageProtocol(m.readerProtocol)
}

// isNativeTUIImageProtocol 判断协议是否是真正的终端图片协议；ANSI/Halfblocks 只是文本兜底。
func isNativeTUIImageProtocol(protocol tuiImageProtocol) bool {
	return protocol != termimg.Unsupported && protocol != termimg.Halfblocks
}

// isTextTUIImageProtocol 判断当前是否处于纯文本图片模式，用于 c 键在图片模式和 ANSI 模式之间切换。
func isTextTUIImageProtocol(protocol tuiImageProtocol) bool {
	return protocol == termimg.Halfblocks || protocol == termimg.Unsupported
}

// invalidateTUIImageState 作废当前帧尺寸和协议状态，切换协议后让封面/阅读页重新按新模式渲染。
func (m *appModel) invalidateTUIImageState() {
	m.coverRequestID++
	m.coverPreview = coverPreviewState{}
	m.markKittyImagesCleared()
	if m.terminalReader.BookID == "" {
		return
	}
	m.readerRequestID++
	m.readerPendingPage = false
	m.readerPendingRequestKey = ""
	m.terminalReader.Width = 0
	m.terminalReader.Height = 0
	m.terminalReader.Protocol = m.readerProtocol
	m.terminalReader.Loading = true
	m.terminalReader.ErrText = ""
	m.terminalReader.Setup = ""
	m.terminalReader.Lines = nil
	m.terminalReader.Overlay = ""
}

// markKittyImagesCleared 标记终端侧 Kitty 图像缓存已经不可用，下一帧需要重新发送 setup。
func (m *appModel) markKittyImagesCleared() {
	m.coverSetupKey = ""
	m.readerSetupKey = ""
}

// renderPendingKittyClearPrefix 消费待清理标记，并作废终端侧 Kitty 图片状态。
func (m *appModel) renderPendingKittyClearPrefix() string {
	if !m.clearKittyImagesNextFrame {
		return ""
	}
	m.clearKittyImagesNextFrame = false
	m.markKittyImagesCleared()
	return termimg.ClearAllString()
}
