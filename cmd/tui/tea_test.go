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
		status:         systemSnapshot{StatusText: "running"},
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
}

func TestMoveFocusCyclesInteractivePanels(t *testing.T) {
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

func TestShouldBypassTUI(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{name: "long flag", args: []string{"comigo", "--no-tui"}, want: true},
		{name: "short flag", args: []string{"comigo", "-n"}, want: true},
		{name: "long true", args: []string{"comigo", "--no-tui=true"}, want: true},
		{name: "long false", args: []string{"comigo", "--no-tui=false"}, want: false},
		{name: "short false", args: []string{"comigo", "-n=false"}, want: false},
		{name: "last value wins", args: []string{"comigo", "--no-tui=false", "-n"}, want: true},
		{name: "after separator", args: []string{"comigo", "--", "-n"}, want: false},
		{name: "absent", args: []string{"comigo", "--open-browser"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldBypassTUI(tt.args); got != tt.want {
				t.Fatalf("shouldBypassTUI(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
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

	t.Setenv("COMIGO_TUI_IMAGE", "kitty")
	if got := detectTUIImageProtocol(); got != termimg.Kitty {
		t.Fatalf("detectTUIImageProtocol() = %v, want Kitty", got)
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
	if got := detectNativeTUIImageProtocol(); got != termimg.Kitty {
		t.Fatalf("detectNativeTUIImageProtocol() = %v, want Kitty for manual Ghostty image mode", got)
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

func TestHalfblocksRenderDoesNotEmitTerminalQueries(t *testing.T) {
	rendered, err := renderTUIImageWithoutQuery(image.NewRGBA(image.Rect(0, 0, 4, 4)), termimg.Halfblocks, 4, 4)
	if err != nil {
		t.Fatalf("renderTUIImageWithoutQuery() error = %v", err)
	}
	if rendered == "" {
		t.Fatal("rendered halfblocks image should not be empty")
	}
	for _, query := range []string{
		"\x1b]1337;ReportCellSize",
		"\x1b_Gi=42",
		"\x1b[?1;1;0S",
	} {
		if strings.Contains(rendered, query) {
			t.Fatalf("halfblocks output should not contain terminal query %q", query)
		}
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

func TestCenterStyledLineKeepsKittyPlaceholderEncoding(t *testing.T) {
	placeholderRow := "\x1b[38;2;0;0;1m" + termimg.CreatePlaceholder(0, 0, 0) + strings.Repeat(termimg.PLACEHOLDER_CHAR, 2) + "\x1b[39m"
	centered := centerStyledLine(placeholderRow, 7)
	if !strings.Contains(centered, termimg.CreatePlaceholder(0, 0, 0)+strings.Repeat(termimg.PLACEHOLDER_CHAR, 2)) {
		t.Fatalf("centered placeholder should keep row/column encoding, got %q", centered)
	}
	if got := xansi.StringWidth(centered); got != 7 {
		t.Fatalf("centered placeholder width = %d, want 7", got)
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

func TestRenderKittyUnicodeImageSplitsSetupAndPlaceholders(t *testing.T) {
	setup, lines, err := renderKittyUnicodeImage(image.NewRGBA(image.Rect(0, 0, 4, 4)), 3, 2)
	if err != nil {
		t.Fatalf("renderKittyUnicodeImage() error = %v", err)
	}
	if !strings.Contains(setup, "a=T,f=100,t=d,i=") || !strings.Contains(setup, "U=1,c=3,r=2") {
		t.Fatalf("kitty setup should transmit PNG and create virtual placement in one command, got %q", setup)
	}
	if strings.Contains(setup, termimg.PLACEHOLDER_CHAR) {
		t.Fatalf("kitty setup should not include visible placeholders, got %q", setup)
	}
	if len(lines) != 2 {
		t.Fatalf("placeholder line count = %d, want 2", len(lines))
	}
	for _, line := range lines {
		if strings.Count(line, termimg.PLACEHOLDER_CHAR) != 3 {
			t.Fatalf("placeholder line should contain 3 cells, got %q", line)
		}
		if got := xansi.StringWidth(line); got != 3 {
			t.Fatalf("placeholder line width = %d, want 3", got)
		}
	}
}

func TestRasterizeKittyPlacementImageUsesCellRectangle(t *testing.T) {
	cellW, cellH := protocolCellPixels(termimg.Kitty)
	got := rasterizeKittyPlacementImage(image.NewRGBA(image.Rect(0, 0, 4, 8)), 3, 2)
	if got.Bounds().Dx() != 3*cellW || got.Bounds().Dy() != 2*cellH {
		t.Fatalf("kitty placement image size = %dx%d, want %dx%d", got.Bounds().Dx(), got.Bounds().Dy(), 3*cellW, 2*cellH)
	}
}

func TestKittyPlaceholderCellOmitsUnusedIDExtraDiacritic(t *testing.T) {
	if got := len([]rune(kittyPlaceholderCell(0, 0, 0))); got != 3 {
		t.Fatalf("24-bit placeholder rune count = %d, want placeholder + row + column", got)
	}
	if got := len([]rune(kittyPlaceholderCell(0, 0, 1))); got != 4 {
		t.Fatalf("32-bit placeholder rune count = %d, want placeholder + row + column + id extra", got)
	}
}

func TestMeasureTUIPlaceholderLines(t *testing.T) {
	_, lines, err := renderKittyUnicodeImage(image.NewRGBA(image.Rect(0, 0, 4, 4)), 3, 2)
	if err != nil {
		t.Fatalf("renderKittyUnicodeImage() error = %v", err)
	}
	metrics := measureTUIPlaceholderLines(lines)
	if metrics.FirstWidth != 3 || metrics.LastWidth != 3 || metrics.MinWidth != 3 || metrics.MaxWidth != 3 {
		t.Fatalf("placeholder widths = %+v, want all 3", metrics)
	}
	if metrics.FirstPlaceholders != 3 || metrics.LastPlaceholders != 3 || metrics.MinPlaceholders != 3 || metrics.MaxPlaceholders != 3 {
		t.Fatalf("placeholder counts = %+v, want all 3", metrics)
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
		status:        systemSnapshot{StatusText: "running"},
	}
	lines := model.renderCoverPreviewContent(panelRect{w: 60, h: 28})
	content := strings.Join(lines, "\n")
	if strings.Contains(content, protocolName(termimg.Halfblocks)) {
		t.Fatalf("preview content should hide protocol details, got:\n%s", content)
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
	if got := model.renderCoverSetupPrefix(); got != "" {
		t.Fatalf("cover setup prefix should only be sent once, got %q", got)
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

func TestShelfSingleClickOnlySelectsItem(t *testing.T) {
	model := &appModel{
		width:         120,
		height:        40,
		focus:         focusShelf,
		coverProtocol: termimg.Halfblocks,
		items: []shelfItem{
			{Kind: shelfItemBook, Title: "Book 1", BookID: "book1", Selectable: true},
			{Kind: shelfItemBook, Title: "Book 2", BookID: "book2", Selectable: true},
		},
		selected:      0,
		shelfRowToID:  make(map[int]int),
		qrButtonFocus: qrActionTerminalReader,
	}
	layout := model.layout()
	_ = model.renderShelfContent(layout.shelf)

	_, _ = model.handleMouse(tea.MouseMsg{
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonLeft,
		X:      layout.shelf.x + 2,
		Y:      layout.shelf.y + 1 + 3,
	})
	if model.selected != 1 {
		t.Fatalf("single click selected = %d, want 1", model.selected)
	}
	if model.screen != screenShelf {
		t.Fatalf("single click screen = %v, want shelf", model.screen)
	}
}

func TestShelfDoubleClickStartsTerminalReaderByDefault(t *testing.T) {
	restoreModelStore(t, &tuiTestStore{books: map[string]*modelpkg.Book{
		"book2": {
			BookInfo:  modelpkg.BookInfo{BookID: "book2", Title: "Book 2", Type: modelpkg.TypeDir},
			PageInfos: modelpkg.PageInfos{{Name: "1.jpg"}},
		},
	}})
	model := &appModel{
		width:               120,
		height:              40,
		focus:               focusShelf,
		readerProtocol:      termimg.Halfblocks,
		terminalReaderCache: make(map[string]terminalReaderState),
		items: []shelfItem{
			{Kind: shelfItemBook, Title: "Book 1", BookID: "book1", Selectable: true},
			{Kind: shelfItemBook, Title: "Book 2", BookID: "book2", Selectable: true},
		},
		selected:            0,
		shelfRowToID:        make(map[int]int),
		qrButtonFocus:       qrActionTerminalReader,
		lastShelfClickIndex: 1,
		lastShelfClickAt:    time.Now(),
	}
	layout := model.layout()
	_ = model.renderShelfContent(layout.shelf)

	_, cmd := model.handleMouse(tea.MouseMsg{
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonLeft,
		X:      layout.shelf.x + 2,
		Y:      layout.shelf.y + 1 + 3,
	})
	if model.screen != screenReader {
		t.Fatalf("double click screen = %v, want terminal reader", model.screen)
	}
	if model.terminalReader.BookID != "book2" {
		t.Fatalf("terminalReader book = %q, want book2", model.terminalReader.BookID)
	}
	if cmd == nil {
		t.Fatal("double click terminal reader should request page render")
	}
}

func TestShelfDoubleClickUsesCurrentBrowserAction(t *testing.T) {
	model := &appModel{
		width:  120,
		height: 40,
		focus:  focusShelf,
		items: []shelfItem{{
			Kind:       shelfItemBook,
			Title:      "Book 1",
			BookID:     "book1",
			TargetURL:  "http://127.0.0.1:1234/reader/book1",
			Selectable: true,
		}},
		shelfRowToID:        make(map[int]int),
		qrButtonFocus:       qrActionOpenBrowser,
		lastShelfClickIndex: 0,
		lastShelfClickAt:    time.Now(),
	}
	layout := model.layout()
	_ = model.renderShelfContent(layout.shelf)

	_, cmd := model.handleMouse(tea.MouseMsg{
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonLeft,
		X:      layout.shelf.x + 2,
		Y:      layout.shelf.y + 1 + 2,
	})
	if model.screen != screenShelf {
		t.Fatalf("browser double click should stay on shelf before command runs, screen=%v", model.screen)
	}
	if cmd == nil {
		t.Fatal("browser double click should request open browser command")
	}
	if !strings.Contains(model.actionMsg, "http://127.0.0.1:1234/reader/book1") {
		t.Fatalf("action message = %q, want opened target URL", model.actionMsg)
	}
}

func TestStartTerminalReaderClearsGhosttyCoverOverlay(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "ghostty")
	t.Setenv("TERM", "xterm-ghostty")
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
		coverProtocol:       termimg.Kitty,
		readerProtocol:      termimg.Halfblocks,
		terminalReaderCache: make(map[string]terminalReaderState),
		items: []shelfItem{{
			Kind:       shelfItemBook,
			Title:      "Book 1",
			BookID:     "book1",
			Selectable: true,
		}},
	}

	_, _ = model.handleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if !model.clearKittyImagesNextFrame {
		t.Fatal("entering reader from Ghostty cover overlay should request Kitty image clear")
	}
	prefix := model.renderTerminalReaderClearPrefix()
	if !strings.Contains(prefix, "a=d,d=A") {
		t.Fatalf("clear prefix should delete Kitty images, got %q", prefix)
	}
	if model.clearKittyImagesNextFrame {
		t.Fatal("clear flag should be consumed after rendering prefix")
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
	if defaultReaderAutoInterval != 10 {
		t.Fatalf("default auto interval = %d, want 10", defaultReaderAutoInterval)
	}
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
	if model.terminalReader.PageIndex != 1 {
		t.Fatalf("right key should keep visible page = %d, want 1", model.terminalReader.PageIndex)
	}
	if !model.readerPendingPage || model.readerPendingPageIndex != 2 {
		t.Fatalf("right key pending = %v/%d, want true/2", model.readerPendingPage, model.readerPendingPageIndex)
	}
	_, _ = model.handleTerminalReaderKey(tea.KeyMsg{Type: tea.KeyLeft})
	if model.readerPendingPage || model.terminalReader.PageIndex != 1 {
		t.Fatalf("left key should cancel pending and keep page 1, pending=%v page=%d", model.readerPendingPage, model.terminalReader.PageIndex)
	}
	_, _ = model.handleTerminalReaderKey(tea.KeyMsg{Type: tea.KeyDelete})
	if model.screen != screenShelf {
		t.Fatalf("delete key screen = %v, want shelf", model.screen)
	}
	model.screen = screenReader
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

func TestTerminalReaderMouseClickPagesAndTitleBack(t *testing.T) {
	model := &appModel{
		screen:                    screenReader,
		width:                     80,
		height:                    24,
		readerProtocol:            termimg.Halfblocks,
		terminalReaderCache:       make(map[string]terminalReaderState),
		clearKittyImagesNextFrame: false,
		terminalReader: terminalReaderState{
			BookID:    "book1",
			Title:     "Book 1",
			PageIndex: 1,
			PageCount: 3,
		},
	}

	_, _ = model.handleMouse(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonLeft, X: 60, Y: 12})
	if model.terminalReader.PageIndex != 1 {
		t.Fatalf("right half click should keep visible page = %d, want 1", model.terminalReader.PageIndex)
	}
	if !model.readerPendingPage || model.readerPendingPageIndex != 2 {
		t.Fatalf("right half click pending = %v/%d, want true/2", model.readerPendingPage, model.readerPendingPageIndex)
	}
	_, _ = model.handleMouse(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonLeft, X: 10, Y: 12})
	if model.readerPendingPage || model.terminalReader.PageIndex != 1 {
		t.Fatalf("left half click should cancel pending and keep page 1, pending=%v page=%d", model.readerPendingPage, model.terminalReader.PageIndex)
	}
	_, _ = model.handleMouse(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonLeft, X: 0, Y: 0})
	if model.screen != screenShelf {
		t.Fatalf("title click screen = %v, want shelf", model.screen)
	}
}

func TestTerminalReaderMouseIgnoresBarsOutsideTitle(t *testing.T) {
	model := &appModel{
		screen: screenReader,
		width:  80,
		height: 24,
		terminalReader: terminalReaderState{
			BookID:    "book1",
			Title:     "Book 1",
			PageIndex: 1,
			PageCount: 3,
		},
	}

	_, _ = model.handleMouse(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonLeft, X: 79, Y: 0})
	if model.screen != screenReader || model.terminalReader.PageIndex != 1 {
		t.Fatalf("top bar outside title should be ignored, screen=%v page=%d", model.screen, model.terminalReader.PageIndex)
	}
	_, _ = model.handleMouse(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonLeft, X: 60, Y: 23})
	if model.terminalReader.PageIndex != 1 {
		t.Fatalf("footer click should not turn page, page=%d", model.terminalReader.PageIndex)
	}
}

func TestGlobalCTogglesImageAndANSIMode(t *testing.T) {
	t.Setenv("COMIGO_TUI_IMAGE", "auto")
	t.Setenv("TERM_PROGRAM", "iTerm.app")
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("LC_TERMINAL", "")
	t.Setenv("ITERM_SESSION_ID", "")
	t.Setenv("GHOSTTY_RESOURCES_DIR", "")
	t.Setenv("WEZTERM_EXECUTABLE", "")
	t.Setenv("WEZTERM_PANE", "")
	t.Setenv("KITTY_WINDOW_ID", "")
	model := &appModel{
		coverProtocol:       termimg.ITerm2,
		readerProtocol:      termimg.ITerm2,
		terminalReaderCache: make(map[string]terminalReaderState),
	}

	_, _ = model.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	if model.coverProtocol != termimg.Halfblocks || model.readerProtocol != termimg.Halfblocks {
		t.Fatalf("c should switch to ANSI mode, cover=%v reader=%v", model.coverProtocol, model.readerProtocol)
	}

	_, _ = model.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	if model.coverProtocol != termimg.ITerm2 || model.readerProtocol != termimg.ITerm2 {
		t.Fatalf("c should switch back to native image mode, cover=%v reader=%v", model.coverProtocol, model.readerProtocol)
	}
	if model.modal.Visible {
		t.Fatal("supported terminal should not show incompatible modal")
	}
}

func TestGlobalCShowsModalWhenImageModeUnsupported(t *testing.T) {
	t.Setenv("COMIGO_TUI_IMAGE", "off")
	model := &appModel{
		width:          80,
		height:         24,
		coverProtocol:  termimg.Halfblocks,
		readerProtocol: termimg.Halfblocks,
	}

	_, _ = model.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	if !model.modal.Visible {
		t.Fatal("unsupported image mode should show modal")
	}
	if !strings.Contains(model.modal.Message, "当前终端不兼容") {
		t.Fatalf("modal message = %q, want incompatible terminal hint", model.modal.Message)
	}
}

func TestModalClosesOnEnterAndMouseOK(t *testing.T) {
	model := &appModel{width: 80, height: 24}
	model.showModal("提示", "消息")
	if !strings.Contains(model.View(), "OK") || model.modal.OKRect.w == 0 {
		t.Fatalf("modal view should render OK button, rect=%+v", model.modal.OKRect)
	}
	_, _ = model.handleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if model.modal.Visible {
		t.Fatal("enter should close modal")
	}

	model.showModal("提示", "消息")
	_ = model.View()
	_, _ = model.handleMouse(tea.MouseMsg{
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonLeft,
		X:      model.modal.OKRect.x,
		Y:      model.modal.OKRect.y,
	})
	if model.modal.Visible {
		t.Fatal("clicking OK should close modal")
	}
}

func TestOpenBrowserFailureShowsModal(t *testing.T) {
	model := &appModel{logBuffer: NewLogBuffer(), shelfRowToID: make(map[int]int)}

	updated, _ := model.Update(openURLResultMsg{url: "http://127.0.0.1:1234", err: errors.New("no browser")})
	got := updated.(*appModel)
	if !got.modal.Visible {
		t.Fatal("browser failure should show modal")
	}
	if !strings.Contains(got.modal.Message, "打开浏览器失败") {
		t.Fatalf("modal message = %q, want browser failure", got.modal.Message)
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
	right := "Comigo 1.2.3"
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

func TestTerminalReaderVersionLineOmitsClock(t *testing.T) {
	line := terminalReaderVersionLine()
	if want := "Comigo " + config.GetVersion(); line != want {
		t.Fatalf("terminal reader version line = %q, want %q", line, want)
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
	if !strings.HasPrefix(title, readerTitleBackPrefix+"Book") {
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

func TestTerminalReaderPageTurnDoesNotClearBeforeNewKittyPage(t *testing.T) {
	model := &appModel{
		readerProtocol:            termimg.Kitty,
		clearKittyImagesNextFrame: false,
		terminalReader: terminalReaderState{
			BookID:    "book1",
			Protocol:  termimg.Kitty,
			PageIndex: 0,
			PageCount: 3,
			Lines:     []string{"old"},
		},
	}

	if !model.moveTerminalReaderPage(1) {
		t.Fatal("page turn should succeed")
	}
	if model.clearKittyImagesNextFrame {
		t.Fatal("page turn should not clear Kitty images until the new page is ready")
	}
	if model.terminalReader.PageIndex != 0 || len(model.terminalReader.Lines) != 1 || model.terminalReader.Lines[0] != "old" {
		t.Fatalf("page turn should keep old visible page, state=%+v", model.terminalReader)
	}
}

func TestTerminalReaderPendingPageReplacesWhenReady(t *testing.T) {
	model := &appModel{
		readerRequestID:        7,
		readerPendingPage:      true,
		readerPendingPageIndex: 2,
		terminalReaderCache:    make(map[string]terminalReaderState),
		terminalReader: terminalReaderState{
			BookID:    "book1",
			Title:     "Book 1",
			PageIndex: 1,
			PageCount: 4,
			Width:     80,
			Height:    20,
			Protocol:  termimg.ITerm2,
			Overlay:   "OLD",
		},
	}

	next := terminalReaderState{
		BookID:    "book1",
		Title:     "Book 1",
		PageIndex: 2,
		PageCount: 4,
		Width:     80,
		Height:    20,
		Protocol:  termimg.ITerm2,
		Overlay:   "NEW",
	}
	model.applyTerminalReaderPageMsg(terminalReaderPageMsg{requestID: 7, state: next})
	if model.readerPendingPage || model.terminalReader.PageIndex != 2 || model.terminalReader.Overlay != "NEW" {
		t.Fatalf("ready page should replace visible state, pending=%v state=%+v", model.readerPendingPage, model.terminalReader)
	}
	key := terminalReaderCacheKey("book1", 2, 80, 20, termimg.ITerm2)
	if _, ok := model.terminalReaderCache[key]; !ok {
		t.Fatalf("ready page should be cached with key %q", key)
	}
}

func TestITerm2PendingPageKeepsOldImageLayer(t *testing.T) {
	model := &appModel{
		readerPendingPage: true,
		terminalReader: terminalReaderState{
			Protocol: termimg.ITerm2,
			Overlay:  "OLD",
			ImageW:   10,
			ImageH:   8,
		},
	}
	if prefix := model.renderTerminalReaderClearPrefix(); strings.Contains(prefix, "\x1b[2J") {
		t.Fatalf("pending iTerm2 page should not clear the full screen, got %q", prefix)
	}
	if overlay := model.renderTerminalReaderOverlay(terminalReaderImageArea{w: 20, h: 12}); overlay != "" {
		t.Fatalf("pending iTerm2 page should keep old overlay without re-rendering it, got %q", overlay)
	}
}

func TestTerminalReaderSetupPrefixOnlyForKitty(t *testing.T) {
	model := &appModel{terminalReader: terminalReaderState{Protocol: termimg.Kitty, Setup: "SETUP"}}
	if got := model.renderTerminalReaderSetupPrefix(); got != "SETUP" {
		t.Fatalf("reader setup prefix = %q, want SETUP", got)
	}
	if got := model.renderTerminalReaderSetupPrefix(); got != "" {
		t.Fatalf("reader setup prefix should only be sent once, got %q", got)
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

func TestReaderProtocolDefaultsToKittyForGhosttyAndPreview(t *testing.T) {
	t.Setenv("COMIGO_TUI_IMAGE", "auto")
	t.Setenv("TERM_PROGRAM", "ghostty")
	t.Setenv("TERM", "xterm-ghostty")
	t.Setenv("LC_TERMINAL", "")
	t.Setenv("ITERM_SESSION_ID", "")
	t.Setenv("WEZTERM_EXECUTABLE", "")
	t.Setenv("WEZTERM_PANE", "")

	if got := detectTUIReaderImageProtocol(); got != termimg.Kitty {
		t.Fatalf("reader protocol = %v, want Kitty for Ghostty auto mode", got)
	}
	if got := detectTUIImageProtocol(); got != termimg.Kitty {
		t.Fatalf("preview protocol = %v, want Kitty for Ghostty auto mode", got)
	}
	if got := detectNativeTUIReaderImageProtocol(); got != termimg.Kitty {
		t.Fatalf("native reader protocol = %v, want Kitty for manual Ghostty image mode", got)
	}
	if got := detectNativeTUIImageProtocol(); got != termimg.Kitty {
		t.Fatalf("native preview protocol = %v, want Kitty for manual Ghostty image mode", got)
	}
}

func TestReaderProtocolDefaultsToKittyForKittyTerminal(t *testing.T) {
	t.Setenv("COMIGO_TUI_IMAGE", "auto")
	t.Setenv("TERM", "xterm-kitty")
	t.Setenv("TERM_PROGRAM", "")
	t.Setenv("KITTY_WINDOW_ID", "1")
	t.Setenv("LC_TERMINAL", "")
	t.Setenv("ITERM_SESSION_ID", "")
	t.Setenv("GHOSTTY_RESOURCES_DIR", "")
	t.Setenv("WEZTERM_PANE", "")

	if got := detectTUIReaderImageProtocol(); got != termimg.Kitty {
		t.Fatalf("reader protocol = %v, want Kitty for Kitty auto mode", got)
	}
	if got := detectTUIImageProtocol(); got != termimg.Kitty {
		t.Fatalf("preview protocol = %v, want Kitty for Kitty auto mode", got)
	}
	if got := detectNativeTUIReaderImageProtocol(); got != termimg.Kitty {
		t.Fatalf("native reader protocol = %v, want Kitty for manual Kitty image mode", got)
	}
	if got := detectNativeTUIImageProtocol(); got != termimg.Kitty {
		t.Fatalf("native preview protocol = %v, want Kitty for manual Kitty image mode", got)
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
