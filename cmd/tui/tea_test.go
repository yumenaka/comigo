package tui

import (
	"errors"
	"image"
	"strconv"
	"strings"
	"testing"
	"time"

	termimg "github.com/blacktop/go-termimg"
	tea "github.com/charmbracelet/bubbletea"
	xansi "github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"
	"github.com/yumenaka/comigo/config"
	modelpkg "github.com/yumenaka/comigo/model"
)

func TestBuildBaseURLUsesConfiguredHost(t *testing.T) {
	restoreConfig(t)
	cfg := config.GetCfg()
	cfg.Host = "example.com"
	cfg.Port = 1234
	cfg.CertFile = ""
	cfg.KeyFile = ""
	cfg.AutoTLSCertificate = false

	if got, want := buildBaseURL(), "http://example.com:1234"; got != want {
		t.Fatalf("buildBaseURL() = %q, want %q", got, want)
	}
}

func TestBuildBaseURLUsesLocalhostWhenLANDisabled(t *testing.T) {
	restoreConfig(t)
	cfg := config.GetCfg()
	cfg.Host = ""
	cfg.Port = 1234
	cfg.DisableLAN = true
	cfg.CertFile = ""
	cfg.KeyFile = ""
	cfg.AutoTLSCertificate = false

	if got, want := buildBaseURL(), "http://127.0.0.1:1234"; got != want {
		t.Fatalf("buildBaseURL() = %q, want %q", got, want)
	}
}

func TestViewLeavesRightmostColumnUnused(t *testing.T) {
	width := 120
	model := &appModel{
		width:          width,
		height:         30,
		logBuffer:      NewLogBuffer(),
		focus:          focusShelf,
		shelfRowToID:   make(map[int]int),
		autoFollowLogs: true,
		status: systemSnapshot{
			CPUPercent: 10,
			RAMPercent: 20,
			StatusText: "running",
		},
	}

	for lineNumber, line := range strings.Split(model.View(), "\n") {
		if got, maxWidth := runewidth.StringWidth(line), width-1; got > maxWidth {
			t.Fatalf("View() line %d width = %d, want <= %d", lineNumber+1, got, maxWidth)
		}
	}
}

func TestWideLayoutPlacesQRCodeTopRightAndPreviewBottomRight(t *testing.T) {
	model := &appModel{width: 120, height: 30}

	layout := model.layout()
	if layout.qr.x <= layout.shelf.x || layout.qr.y != layout.shelf.y {
		t.Fatalf("QR panel should be placed to the right of shelf: shelf=%+v qr=%+v", layout.shelf, layout.qr)
	}
	if layout.cover.x != layout.qr.x || layout.cover.y <= layout.qr.y {
		t.Fatalf("preview panel should be placed below QR on the right: qr=%+v cover=%+v", layout.qr, layout.cover)
	}
	if layout.shelf.w <= layout.cover.w {
		t.Fatalf("right preview column should be narrower than shelf: shelf=%+v cover=%+v", layout.shelf, layout.cover)
	}
	if diff := layout.qr.h - layout.cover.h; diff < 0 || diff > 1 {
		t.Fatalf("right column should split QR and preview 1:1: qr=%+v cover=%+v", layout.qr, layout.cover)
	}
	if layout.qr.h+layout.cover.h+layoutGap != model.height {
		t.Fatalf("right column should fill full height: qr=%+v cover=%+v height=%d", layout.qr, layout.cover, model.height)
	}
	if layout.info.w != 0 || layout.info.h != 0 {
		t.Fatalf("info panel should stay hidden: %+v", layout.info)
	}
}

func TestMoveFocusSkipsHiddenInfoPanel(t *testing.T) {
	model := &appModel{focus: focusShelf}

	model.moveFocus(1)
	if model.focus != focusLog {
		t.Fatalf("first focus move = %v, want focusLog", model.focus)
	}
	model.moveFocus(1)
	if model.focus != focusQRCode {
		t.Fatalf("second focus move = %v, want focusQRCode", model.focus)
	}
	model.moveFocus(1)
	if model.focus != focusShelf {
		t.Fatalf("third focus move = %v, want focusShelf", model.focus)
	}
	model.moveFocus(-1)
	if model.focus != focusQRCode {
		t.Fatalf("reverse focus move = %v, want focusQRCode", model.focus)
	}
}

func TestShelfContentDoesNotRenderBottomStatus(t *testing.T) {
	model := &appModel{
		status:         systemSnapshot{StatusText: "RUNNING_SENTINEL"},
		shelfRowToID:   make(map[int]int),
		autoFollowLogs: true,
	}

	lines := model.renderShelfContent(panelRect{w: 60, h: 10})
	if strings.Contains(strings.Join(lines, "\n"), "RUNNING_SENTINEL") {
		t.Fatalf("shelf content should not render bottom status, got:\n%s", strings.Join(lines, "\n"))
	}
}

func TestLogContentRendersBottomStatus(t *testing.T) {
	model := &appModel{
		logs:           []string{"log line"},
		status:         systemSnapshot{StatusText: "RUNNING_SENTINEL"},
		autoFollowLogs: true,
	}

	lines := model.renderLogContent(panelRect{w: 60, h: 8})
	if len(lines) == 0 || !strings.Contains(lines[len(lines)-1], "RUNNING_SENTINEL") {
		t.Fatalf("log content should render status on bottom line, got:\n%s", strings.Join(lines, "\n"))
	}
}

func TestQRCodeContentDoesNotInsertBlankAroundQRCode(t *testing.T) {
	model := &appModel{
		qrLines:         []string{"QR-A", "QR-B"},
		currentShelfURL: "http://127.0.0.1:1234",
		readMode:        0,
	}

	lines := model.renderQRCodeContent(panelRect{w: 60, h: 20})
	if strings.TrimSpace(lines[2]) != "QR-A" || strings.TrimSpace(lines[3]) != "QR-B" {
		t.Fatalf("QR lines should start immediately after URL, got:\n%s", strings.Join(lines, "\n"))
	}
	if strings.TrimSpace(lines[4]) == "" {
		t.Fatalf("mode row should follow QR lines without an empty spacer, got:\n%s", strings.Join(lines, "\n"))
	}
	if len(model.qrButtonHitbox) != 5 {
		t.Fatalf("QR buttons = %d, want 5", len(model.qrButtonHitbox))
	}
	if !strings.Contains(strings.Join(lines, "\n"), "终端阅读") {
		t.Fatalf("QR content should include terminal reader button, got:\n%s", strings.Join(lines, "\n"))
	}
}

func TestClipAndPadStyledIgnoresANSIWidth(t *testing.T) {
	line := clipAndPadStyled("\x1b[31mHi\x1b[0m", 5)
	if got := xansi.StringWidth(line); got != 5 {
		t.Fatalf("StringWidth(%q) = %d, want 5", line, got)
	}
	if !strings.Contains(line, "\x1b[31m") {
		t.Fatalf("styled line should preserve ANSI escape sequence: %q", line)
	}
}

func TestDetectTUIImageProtocolAllowsEnvOverride(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "iTerm.app")
	t.Setenv("COMIGO_TUI_IMAGE", "auto")
	if got := detectTUIImageProtocol(); got != termimg.ITerm2 {
		t.Fatalf("detectTUIImageProtocol() = %v, want ITerm2", got)
	}

	t.Setenv("COMIGO_TUI_IMAGE", "ansi")
	if got := detectTUIImageProtocol(); got != termimg.Halfblocks {
		t.Fatalf("detectTUIImageProtocol() = %v, want Halfblocks", got)
	}

	t.Setenv("COMIGO_TUI_IMAGE", "off")
	if got := detectTUIImageProtocol(); got != termimg.Unsupported {
		t.Fatalf("detectTUIImageProtocol() = %v, want Unsupported", got)
	}
}

func TestDetectTUIImageProtocolUsesTerminalSpecificProtocols(t *testing.T) {
	t.Setenv("COMIGO_TUI_IMAGE", "auto")
	t.Setenv("TERM", "xterm-ghostty")
	t.Setenv("TERM_PROGRAM", "")
	t.Setenv("LC_TERMINAL", "")
	t.Setenv("ITERM_SESSION_ID", "")
	t.Setenv("GHOSTTY_RESOURCES_DIR", "")
	t.Setenv("WEZTERM_EXECUTABLE", "")
	t.Setenv("WEZTERM_PANE", "")
	if got := detectTUIImageProtocol(); got != termimg.Kitty {
		t.Fatalf("detectTUIImageProtocol() = %v, want Kitty for Ghostty auto mode", got)
	}

	t.Setenv("TERM", "xterm-256color")
	t.Setenv("WEZTERM_PANE", "1")
	if got := detectTUIImageProtocol(); got != termimg.ITerm2 {
		t.Fatalf("detectTUIImageProtocol() = %v, want ITerm2 for WezTerm preview", got)
	}
}

func TestPreviewImageFrameUsesAllMiddleRows(t *testing.T) {
	frame, ok := previewImageFrameFor(panelRect{w: 90, h: 40})
	if !ok {
		t.Fatal("previewImageFrameFor() should return a frame")
	}
	if frame.w != 88 || frame.h != 36 {
		t.Fatalf("frame size = %dx%d, want full preview content width and middle rows 88x36", frame.w, frame.h)
	}
	if frame.x != 0 || frame.y != 1 {
		t.Fatalf("frame position = %d,%d, want 0,1", frame.x, frame.y)
	}
	if frame.innerW != 86 || frame.innerH != 34 {
		t.Fatalf("inner frame size = %dx%d, want 86x34", frame.innerW, frame.innerH)
	}
}

func TestFitImageCellsKeepsPortraitCenteredSize(t *testing.T) {
	w, h := fitImageCells(image.Rect(0, 0, 500, 1000), 50, 12)
	if w != 12 || h != 12 {
		t.Fatalf("portrait fit = %dx%d, want 12x12", w, h)
	}
}

func TestFitImageCellsCanFillWidth(t *testing.T) {
	w, h := fitImageCells(image.Rect(0, 0, 1000, 300), 50, 12)
	if w != 50 || h != 8 {
		t.Fatalf("wide fit = %dx%d, want 50x8", w, h)
	}
}

func TestProtocolFitUsesVisibleDisplayArea(t *testing.T) {
	w, h := fitImageCellsForProtocol(image.Rect(0, 0, 500, 1000), 50, 12, termimg.Halfblocks)
	if w != 12 || h != 12 {
		t.Fatalf("halfblocks portrait fit = %dx%d, want 12x12", w, h)
	}

	w, h = fitImageCellsForProtocol(image.Rect(0, 0, 500, 1000), 50, 12, termimg.Kitty)
	if w != 12 || h != 12 {
		t.Fatalf("kitty portrait fit = %dx%d, want 12x12", w, h)
	}
}

func TestFitImageCellsWithCellPixelsUsesNativeCellRatio(t *testing.T) {
	w, h := fitImageCellsWithCellPixels(image.Rect(0, 0, 100, 100), 10, 10, 10, 20)
	if w != 10 || h != 5 {
		t.Fatalf("native cell fit = %dx%d, want 10x5", w, h)
	}
}

func TestProtocolCellPixelsUsesWezTermGeometryBeforeITerm2Protocol(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "WezTerm")
	t.Setenv("WEZTERM_PANE", "1")

	w, h := protocolCellPixels(termimg.ITerm2)
	if w != termimg.WezTermWidth || h != termimg.WezTermHeight {
		t.Fatalf("WezTerm iTerm2 geometry = %dx%d, want %dx%d", w, h, termimg.WezTermWidth, termimg.WezTermHeight)
	}
}

func TestHalfblocksRenderSizeCompensatesMosaicCells(t *testing.T) {
	w, h := termImageRenderSizeForProtocol(50, 12, termimg.Halfblocks)
	if w != 100 || h != 24 {
		t.Fatalf("halfblocks render size = %dx%d, want 100x24", w, h)
	}

	w, h = termImageRenderSizeForProtocol(50, 12, termimg.Kitty)
	if w != 50 || h != 12 {
		t.Fatalf("kitty render size = %dx%d, want 50x12", w, h)
	}
}

func TestSplitRenderedImageLinesSeparatesKittySetup(t *testing.T) {
	rendered := "\x1b_Ga=T,i=1\x1b\\\x1b_Ga=p,U=1,i=1,c=2,r=2\x1b\\" +
		"\x1b[38;2;0;0;1m" + termimg.PLACEHOLDER_CHAR + termimg.PLACEHOLDER_CHAR + "\x1b[39m\n" +
		"\x1b[38;2;0;0;1m" + termimg.PLACEHOLDER_CHAR + termimg.PLACEHOLDER_CHAR + "\x1b[39m"

	setup, lines := splitRenderedImageLines(rendered, termimg.Kitty)
	if !strings.Contains(setup, "\x1b_Ga=T") || !strings.Contains(setup, "U=1") {
		t.Fatalf("setup should keep Kitty control sequences, got %q", setup)
	}
	if len(lines) != 2 {
		t.Fatalf("line count = %d, want 2", len(lines))
	}
	if strings.Contains(strings.Join(lines, "\n"), "\x1b_G") {
		t.Fatalf("visible lines should not contain Kitty control sequences: %#v", lines)
	}
	if !strings.Contains(lines[0], "\x1b[38;2;0;0;1m") {
		t.Fatalf("visible placeholder line should keep Kitty image-id color, got %q", lines[0])
	}
}

func TestSplitRenderedImageLinesKeepsAnsiInline(t *testing.T) {
	setup, lines := splitRenderedImageLines("a\nb", termimg.Halfblocks)
	if setup != "" {
		t.Fatalf("ANSI setup = %q, want empty", setup)
	}
	if got := strings.Join(lines, "|"); got != "a|b" {
		t.Fatalf("ANSI lines = %q, want a|b", got)
	}
}

func TestCoverResizeHeightUsesHighResolutionForTUI(t *testing.T) {
	if got := coverResizeHeight(10, termimg.Halfblocks); got < coverPreviewMinResizeHeight {
		t.Fatalf("halfblocks resize height = %d, want at least %d", got, coverPreviewMinResizeHeight)
	}
	if got := coverResizeHeight(10, termimg.Kitty); got < coverPreviewMinResizeHeight {
		t.Fatalf("kitty resize height = %d, want at least %d", got, coverPreviewMinResizeHeight)
	}
}

func TestCenterPreviewImageLinesCentersHorizontallyAndVertically(t *testing.T) {
	lines := centerPreviewImageLines([]string{"xx"}, 6, 3)
	if len(lines) != 3 {
		t.Fatalf("line count = %d, want 3", len(lines))
	}
	if lines[1] != "  xx  " {
		t.Fatalf("centered line = %q, want %q", lines[1], "  xx  ")
	}
}

func TestCenterPreviewImageLinesCropsFromCenter(t *testing.T) {
	lines := centerPreviewImageLines([]string{"1", "2", "3", "4", "5"}, 3, 3)
	got := strings.Join(lines, "|")
	if got != " 2 | 3 | 4 " {
		t.Fatalf("center cropped lines = %q", got)
	}
}

func TestPreviewContentShowsOnlyVersionAtBottom(t *testing.T) {
	model := &appModel{
		items: []shelfItem{{
			Kind:       shelfItemBook,
			Title:      "Book",
			BookID:     "book1",
			Selectable: true,
		}},
		selected:      0,
		coverProtocol: termimg.Halfblocks,
		status: systemSnapshot{
			CPUPercent: 10,
			RAMPercent: 20,
			StatusText: "running",
			TargetURL:  "http://127.0.0.1:1234/scroll/book1",
		},
	}
	lines := model.renderCoverPreviewContent(panelRect{w: 60, h: 28})
	content := strings.Join(lines, "\n")
	if strings.Contains(content, "10.0") || strings.Contains(content, protocolName(termimg.Halfblocks)) || strings.Contains(content, model.status.TargetURL) {
		t.Fatalf("preview content should hide status details, protocol and target URL, got:\n%s", content)
	}
	if got := strings.TrimSpace(lines[len(lines)-1]); !strings.Contains(got, "Comigo ") {
		t.Fatalf("bottom line = %q, want Comigo version line", got)
	}
}

func TestRenderCoverOverlayClearsForcedOverlayImageArea(t *testing.T) {
	model := &appModel{
		coverProtocol: termimg.ITerm2,
		items: []shelfItem{{
			Kind:       shelfItemBook,
			Title:      "Book",
			BookID:     "book1",
			Selectable: true,
		}},
		selected: 0,
		coverPreview: coverPreviewState{
			BookID:   "book1",
			Protocol: termimg.ITerm2,
			Overlay:  "IMAGE",
			ImageW:   10,
			ImageH:   5,
		},
	}
	frame, _ := previewImageFrameFor(panelRect{w: 60, h: 24})
	overlay := model.renderCoverOverlay(panelRect{w: 60, h: 24})
	if !strings.Contains(overlay, strings.Repeat(" ", frame.innerW)) {
		t.Fatalf("overlay should actively clear the image area")
	}
	if !strings.Contains(overlay, "\x1b["+strconv.Itoa(frame.innerW)+"X") {
		t.Fatalf("iTerm2 overlay should clear cells with ECH, got %q", overlay)
	}
	if !strings.Contains(overlay, "IMAGE") {
		t.Fatalf("overlay should include the rendered image sequence")
	}
}

func TestRenderCoverOverlaySupportsGhosttyKittyPreview(t *testing.T) {
	t.Setenv("TERM", "xterm-ghostty")
	t.Setenv("TERM_PROGRAM", "")
	t.Setenv("GHOSTTY_RESOURCES_DIR", "")
	model := &appModel{
		coverProtocol: termimg.Kitty,
		items: []shelfItem{{
			Kind:       shelfItemBook,
			Title:      "Book",
			BookID:     "book1",
			Selectable: true,
		}},
		selected: 0,
		coverPreview: coverPreviewState{
			BookID:   "book1",
			Protocol: termimg.Kitty,
			Overlay:  "IMAGE",
			ImageW:   10,
			ImageH:   5,
		},
	}
	overlay := model.renderCoverOverlay(panelRect{w: 60, h: 24})
	if !strings.Contains(overlay, "a=d,d=A") || !strings.Contains(overlay, "IMAGE") {
		t.Fatalf("Ghostty Kitty preview overlay should clear and draw image, got %q", overlay)
	}
}

func TestRenderITerm2InlineImageUsesCellUnits(t *testing.T) {
	rendered := renderITerm2InlineImage([]byte("abc"), 10, 5, 20, 5)
	if !strings.Contains(rendered, "width=auto;height=5") {
		t.Fatalf("height-limited iTerm2 image should leave width auto, got %q", rendered)
	}
	rendered = renderITerm2InlineImage([]byte("abc"), 10, 5, 10, 20)
	if !strings.Contains(rendered, "width=10;height=auto") {
		t.Fatalf("width-limited iTerm2 image should leave height auto, got %q", rendered)
	}
	rendered = renderITerm2InlineImage([]byte("abc"), 10, 5, 10, 5)
	if !strings.Contains(rendered, "width=10;height=5") {
		t.Fatalf("exact-fit iTerm2 image should use both cell dimensions, got %q", rendered)
	}
	if strings.Contains(rendered, "px") {
		t.Fatalf("iTerm2 image should not use px units, got %q", rendered)
	}
	if !strings.Contains(rendered, "preserveAspectRatio=1") || !strings.Contains(rendered, "doNotMoveCursor=1") {
		t.Fatalf("iTerm2 image should preserve aspect ratio and not move cursor, got %q", rendered)
	}
}

func TestRenderCoverClearPrefixForForcedOverlayProtocols(t *testing.T) {
	model := &appModel{coverProtocol: termimg.ITerm2}
	if prefix := model.renderCoverClearPrefix(panelRect{w: 60, h: 24}); prefix == "" {
		t.Fatal("iTerm2 protocol should clear the screen before redraw")
	}

	model.coverProtocol = termimg.Sixel
	if prefix := model.renderCoverClearPrefix(panelRect{w: 60, h: 24}); prefix != "" {
		t.Fatalf("Sixel should not use iTerm2 clear prefix, got %q", prefix)
	}

	model.coverProtocol = termimg.Kitty
	if prefix := model.renderCoverClearPrefix(panelRect{w: 60, h: 24}); prefix != "" {
		t.Fatalf("Kitty placeholder path should not use overlay clear prefix, got %q", prefix)
	}
}

func TestRenderCoverSetupPrefixOnlyForCurrentKittyImage(t *testing.T) {
	model := &appModel{
		coverPreview: coverPreviewState{BookID: "book1", Protocol: termimg.Kitty, Setup: "SETUP"},
		items: []shelfItem{{
			Kind:       shelfItemBook,
			BookID:     "book1",
			Selectable: true,
		}},
	}
	if got := model.renderCoverSetupPrefix(); got != "SETUP" {
		t.Fatalf("cover setup prefix = %q, want SETUP", got)
	}
	model.coverPreview.BookID = "book2"
	if got := model.renderCoverSetupPrefix(); got != "" {
		t.Fatalf("stale cover setup prefix = %q, want empty", got)
	}
}

func TestSyncCoverPreviewCacheInvalidatesStaleLoadingRequest(t *testing.T) {
	model := &appModel{
		width:         120,
		height:        30,
		coverProtocol: termimg.Halfblocks,
		coverCache:    make(map[string]coverPreviewState),
		items: []shelfItem{{
			Kind:       shelfItemBook,
			Title:      "Book 2",
			BookID:     "book2",
			Selectable: true,
		}},
		selected:       0,
		coverRequestID: 7,
		coverPreview:   coverPreviewState{BookID: "book1", Loading: true, Protocol: termimg.Halfblocks},
		shelfRowToID:   make(map[int]int),
		autoFollowLogs: true,
	}
	frame, ok := previewImageFrameFor(model.layout().cover)
	if !ok {
		t.Fatal("preview frame should exist")
	}
	cached := coverPreviewState{BookID: "book2", Width: frame.innerW, Height: frame.innerH, Protocol: termimg.Halfblocks}
	model.coverCache[coverPreviewCacheKey("book2", frame.innerW, frame.innerH, termimg.Halfblocks)] = cached

	if cmd := model.syncCoverPreviewCmd(); cmd != nil {
		t.Fatal("cache hit should not start a new command")
	}
	if model.coverRequestID != 8 {
		t.Fatalf("coverRequestID = %d, want 8", model.coverRequestID)
	}
	if model.coverPreview.BookID != "book2" || model.coverPreview.Loading {
		t.Fatalf("coverPreview = %+v, want cached book2", model.coverPreview)
	}
}

func TestShelfEnterStartsTerminalReader(t *testing.T) {
	restoreModelStore(t, &tuiTestStore{books: map[string]*modelpkg.Book{
		"book1": {
			BookInfo:  modelpkg.BookInfo{BookID: "book1", Title: "Book 1", Type: modelpkg.TypeDir},
			PageInfos: modelpkg.PageInfos{{Name: "1.jpg"}},
		},
	}})
	model := &appModel{
		width:               100,
		height:              40,
		focus:               focusShelf,
		readerProtocol:      termimg.Halfblocks,
		terminalReaderCache: make(map[string]terminalReaderState),
		items: []shelfItem{{
			Kind:       shelfItemBook,
			Title:      "Book 1",
			BookID:     "book1",
			Selectable: true,
		}},
	}

	_, cmd := model.handleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if model.screen != screenReader {
		t.Fatalf("screen = %v, want terminal reader", model.screen)
	}
	if cmd == nil {
		t.Fatal("starting terminal reader should request first page render")
	}
	if model.terminalReader.BookID != "book1" || model.terminalReader.PageIndex != 0 || model.terminalReader.PageCount != 1 {
		t.Fatalf("terminalReader = %+v, want first page of book1", model.terminalReader)
	}
}

func TestQRTerminalReaderRequiresBookItem(t *testing.T) {
	model := &appModel{
		qrButtonFocus: qrActionTerminalReader,
		items: []shelfItem{{
			Kind:       shelfItemGroup,
			Title:      "Group",
			BookID:     "group1",
			Selectable: true,
		}},
		selected: 0,
	}

	if cmd := model.executeQRButton(); cmd != nil {
		t.Fatalf("group selection should not start terminal reader, got %v", cmd)
	}
	if model.screen == screenReader {
		t.Fatal("group selection should stay on shelf screen")
	}
}

func TestTerminalReaderQReturnsToShelf(t *testing.T) {
	model := &appModel{
		screen:         screenReader,
		coverProtocol:  termimg.Halfblocks,
		terminalReader: terminalReaderState{Protocol: termimg.Kitty, Overlay: "IMAGE"},
	}

	_, cmd := model.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if model.screen != screenShelf {
		t.Fatalf("screen = %v, want shelf", model.screen)
	}
	if model.readerAutoFlip {
		t.Fatal("reader auto flip should stop when leaving terminal reader")
	}
	if cmd != nil {
		t.Fatalf("empty shelf return should not need a command, got %v", cmd)
	}
}

func TestTerminalReaderKeyBindings(t *testing.T) {
	model := &appModel{
		screen:             screenReader,
		readerAutoInterval: defaultReaderAutoInterval,
		terminalReader: terminalReaderState{
			BookID:    "book1",
			Title:     "Book 1",
			PageIndex: 1,
			PageCount: 3,
		},
	}

	_, _ = model.handleTerminalReaderKey(tea.KeyMsg{Type: tea.KeyRight})
	if model.terminalReader.PageIndex != 2 {
		t.Fatalf("right key page = %d, want 2", model.terminalReader.PageIndex)
	}
	_, _ = model.handleTerminalReaderKey(tea.KeyMsg{Type: tea.KeyLeft})
	if model.terminalReader.PageIndex != 1 {
		t.Fatalf("left key page = %d, want 1", model.terminalReader.PageIndex)
	}
	_, _ = model.handleTerminalReaderKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	if !model.terminalReaderFullscreen {
		t.Fatal("f should enable fullscreen")
	}
	_, _ = model.handleTerminalReaderKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	if !model.readerAutoFlip || model.readerNextAutoAt.IsZero() {
		t.Fatal("a should start auto flip and schedule next page")
	}
	_, _ = model.handleTerminalReaderKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}})
	if model.readerAutoInterval != defaultReaderAutoInterval+1 {
		t.Fatalf("auto interval = %d, want %d", model.readerAutoInterval, defaultReaderAutoInterval+1)
	}
}

func TestReaderAutoFlipStopsAtLastPage(t *testing.T) {
	model := &appModel{
		screen:             screenReader,
		readerAutoFlip:     true,
		readerAutoInterval: defaultReaderAutoInterval,
		readerNextAutoAt:   time.Now().Add(-time.Second),
		terminalReader: terminalReaderState{
			BookID:    "book1",
			PageIndex: 0,
			PageCount: 1,
		},
	}

	if cmd := model.handleReaderAutoFlip(time.Now()); cmd != nil {
		t.Fatalf("last page should not request render command, got %v", cmd)
	}
	if model.readerAutoFlip {
		t.Fatal("auto flip should stop at last page")
	}
}

func TestTerminalReaderFooterKeepsRightStatusVisible(t *testing.T) {
	right := "2026-05-17 15:43:09  Comigo 1.2.3"
	line := formatThreePartStatusLine("very long shortcut hint that should be shortened first", "1/33", right, 80)
	if !strings.Contains(line, right) {
		t.Fatalf("footer should keep right status visible, got %q", line)
	}
	if !strings.Contains(line, "1/33") {
		t.Fatalf("footer should include centered page status, got %q", line)
	}
	if got := runewidth.StringWidth(line); got != 80 {
		t.Fatalf("footer width = %d, want 80", got)
	}
}

func TestTerminalReaderVersionLineUsesMinutePrecision(t *testing.T) {
	line := terminalReaderVersionLineAt(time.Date(2026, 5, 17, 15, 43, 9, 0, time.Local))
	if !strings.Contains(line, "2026-05-17 15:43") {
		t.Fatalf("terminal reader version line should include minute precision time, got %q", line)
	}
	if strings.Contains(line, "15:43:09") {
		t.Fatalf("terminal reader version line should not include seconds, got %q", line)
	}
}

func TestTerminalReaderLoadingAppearsInTitleBar(t *testing.T) {
	model := &appModel{
		terminalReader: terminalReaderState{
			Title:   "Book",
			Loading: true,
		},
	}
	title := model.renderTerminalReaderTitle(24)
	if !strings.HasPrefix(title, "Book") {
		t.Fatalf("title line should keep book title on the left, got %q", title)
	}
	if !strings.HasSuffix(title, "正在加载页面...") {
		t.Fatalf("title line should show loading text on the right, got %q", title)
	}

	imageLines := model.renderTerminalReaderImageLines(terminalReaderImageArea{w: 24, h: 3})
	if strings.TrimSpace(strings.Join(imageLines, "")) != "" {
		t.Fatalf("normal loading state should not render text inside image area, got %#v", imageLines)
	}

	model.terminalReaderFullscreen = true
	imageLines = model.renderTerminalReaderImageLines(terminalReaderImageArea{w: 24, h: 3})
	if strings.TrimSpace(strings.Join(imageLines, "")) != "" {
		t.Fatalf("fullscreen loading state should not render loading text, got %#v", imageLines)
	}
}

func TestTerminalReaderSetupPrefixOnlyForKitty(t *testing.T) {
	model := &appModel{terminalReader: terminalReaderState{Protocol: termimg.Kitty, Setup: "SETUP"}}
	if got := model.renderTerminalReaderSetupPrefix(); got != "SETUP" {
		t.Fatalf("reader setup prefix = %q, want SETUP", got)
	}
	model.terminalReader.Protocol = termimg.Halfblocks
	if got := model.renderTerminalReaderSetupPrefix(); got != "" {
		t.Fatalf("halfblocks setup prefix = %q, want empty", got)
	}
}

func TestTerminalReaderUsesPlaceholderForKitty(t *testing.T) {
	if isTerminalReaderOverlayProtocol(termimg.Kitty) {
		t.Fatal("Kitty terminal reader should use placeholder text rendering")
	}
	if isTerminalReaderOverlayProtocol(termimg.Halfblocks) {
		t.Fatal("Halfblocks terminal reader should stay in text mode")
	}
}

func TestTopRightPreviewTextPlacesLoadingAtTopRight(t *testing.T) {
	lines := topRightPreviewText("loading", 10, 3)
	if len(lines) != 3 {
		t.Fatalf("line count = %d, want 3", len(lines))
	}
	if lines[0] != "   loading" {
		t.Fatalf("first line = %q, want right aligned loading text", lines[0])
	}
	if strings.TrimSpace(lines[1]) != "" || strings.TrimSpace(lines[2]) != "" {
		t.Fatalf("only first line should contain loading text, got %#v", lines)
	}
}

func TestReaderProtocolUsesKittyForGhosttyAndPreview(t *testing.T) {
	t.Setenv("COMIGO_TUI_IMAGE", "auto")
	t.Setenv("TERM_PROGRAM", "ghostty")
	t.Setenv("TERM", "xterm-ghostty")
	t.Setenv("LC_TERMINAL", "")
	t.Setenv("ITERM_SESSION_ID", "")
	t.Setenv("WEZTERM_EXECUTABLE", "")
	t.Setenv("WEZTERM_PANE", "")

	if got := detectTUIReaderImageProtocol(); got != termimg.Kitty {
		t.Fatalf("reader protocol = %v, want Kitty for Ghostty", got)
	}
	if got := detectTUIImageProtocol(); got != termimg.Kitty {
		t.Fatalf("preview protocol = %v, want Kitty for Ghostty preview", got)
	}
}

func TestReaderProtocolUsesKittyForKittyTerminal(t *testing.T) {
	t.Setenv("COMIGO_TUI_IMAGE", "auto")
	t.Setenv("TERM", "xterm-kitty")
	t.Setenv("TERM_PROGRAM", "")
	t.Setenv("KITTY_WINDOW_ID", "1")
	t.Setenv("LC_TERMINAL", "")
	t.Setenv("ITERM_SESSION_ID", "")
	t.Setenv("GHOSTTY_RESOURCES_DIR", "")
	t.Setenv("WEZTERM_PANE", "")

	if got := detectTUIReaderImageProtocol(); got != termimg.Kitty {
		t.Fatalf("reader protocol = %v, want Kitty for Kitty terminal", got)
	}
}

func TestReaderProtocolUsesITerm2ForWezTerm(t *testing.T) {
	t.Setenv("COMIGO_TUI_IMAGE", "auto")
	t.Setenv("TERM_PROGRAM", "WezTerm")
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("LC_TERMINAL", "")
	t.Setenv("ITERM_SESSION_ID", "")
	t.Setenv("GHOSTTY_RESOURCES_DIR", "")
	t.Setenv("WEZTERM_PANE", "1")

	if got := detectTUIReaderImageProtocol(); got != termimg.ITerm2 {
		t.Fatalf("reader protocol = %v, want ITerm2 for WezTerm", got)
	}
	if got := detectTUIImageProtocol(); got != termimg.ITerm2 {
		t.Fatalf("preview protocol = %v, want ITerm2 for WezTerm preview", got)
	}
}

func restoreConfig(t *testing.T) {
	t.Helper()
	original := config.CopyCfg()
	t.Cleanup(func() {
		*config.GetCfg() = original
	})
}

func restoreModelStore(t *testing.T, store modelpkg.StoreInterface) {
	t.Helper()
	original := modelpkg.IStore
	modelpkg.IStore = store
	t.Cleanup(func() {
		modelpkg.IStore = original
	})
}

type tuiTestStore struct {
	books map[string]*modelpkg.Book
}

func (s *tuiTestStore) StoreBook(b *modelpkg.Book) error {
	if s.books == nil {
		s.books = make(map[string]*modelpkg.Book)
	}
	s.books[b.BookID] = b
	return nil
}

func (s *tuiTestStore) GetBook(id string) (*modelpkg.Book, error) {
	if book, ok := s.books[id]; ok {
		return book, nil
	}
	return nil, errors.New("book not found")
}

func (s *tuiTestStore) DeleteBook(id string) error {
	delete(s.books, id)
	return nil
}

func (s *tuiTestStore) ListBooks() ([]*modelpkg.Book, error) {
	books := make([]*modelpkg.Book, 0, len(s.books))
	for _, book := range s.books {
		books = append(books, book)
	}
	return books, nil
}

func (s *tuiTestStore) GenerateBookGroup() error {
	return nil
}

func (s *tuiTestStore) StoreBookMark(_ *modelpkg.BookMark) error {
	return nil
}

func (s *tuiTestStore) GetBookMarks(_ string) (*modelpkg.BookMarks, error) {
	return nil, nil
}

func (s *tuiTestStore) DeleteBookMark(_ string, _ modelpkg.MarkType, _ int) error {
	return nil
}
