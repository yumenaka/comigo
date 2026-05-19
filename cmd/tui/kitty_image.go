package tui

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"
	"sync/atomic"
	"time"

	termimg "github.com/blacktop/go-termimg"
)

const (
	kittyBase64ChunkSize = 3072
	kittyDiacriticLimit  = 297
	kittyImageIDMask     = 0x00ffffff
)

var tuiKittyImageID = (uint32(os.Getpid()<<8) ^ uint32(time.Now().UnixNano())) & kittyImageIDMask

// renderKittyUnicodeImage 按 Kitty 官方 Unicode placeholder 流程生成图片：
// 先安静传输 PNG 数据并创建 virtual placement，再把可见占位符行交给 Bubble Tea 布局。
func renderKittyUnicodeImage(img image.Image, cols int, rows int) (string, []string, error) {
	if img == nil {
		return "", nil, fmt.Errorf("empty image")
	}
	cols = max(1, cols)
	rows = min(max(1, rows), kittyDiacriticLimit)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", nil, fmt.Errorf("encode kitty png: %w", err)
	}

	imageID := nextTUIKittyImageID()
	setup := renderKittyUnicodeSetup(buf.Bytes(), imageID, cols, rows)
	lines := renderKittyPlaceholderLines(imageID, cols, rows)
	return setup, lines, nil
}

func nextTUIKittyImageID() uint32 {
	id := atomic.AddUint32(&tuiKittyImageID, 1) & kittyImageIDMask
	if id == 0 {
		id = atomic.AddUint32(&tuiKittyImageID, 1) & kittyImageIDMask
	}
	return id
}

func renderKittyUnicodeSetup(data []byte, imageID uint32, cols int, rows int) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	var builder strings.Builder
	first := true
	for len(encoded) > 0 {
		chunk := encoded
		if len(chunk) > kittyBase64ChunkSize {
			chunk = encoded[:kittyBase64ChunkSize]
		}
		encoded = encoded[len(chunk):]
		more := 0
		if len(encoded) > 0 {
			more = 1
		}

		var seq string
		if first {
			// a=T 加 U=1 会在完整接收图片后只创建 virtual placement，不把图片直接画到光标处。
			seq = fmt.Sprintf("\x1b_Ga=T,f=100,t=d,i=%d,U=1,c=%d,r=%d,q=2,m=%d;%s\x1b\\", imageID, cols, rows, more, chunk)
			first = false
		} else {
			seq = fmt.Sprintf("\x1b_Gm=%d;%s\x1b\\", more, chunk)
		}
		builder.WriteString(wrapKittyTmuxPassthrough(seq))
	}
	return builder.String()
}

func renderKittyPlaceholderLines(imageID uint32, cols int, rows int) []string {
	lines := make([]string, 0, rows)
	colorStart := kittyPlaceholderColorStart(imageID)
	idExtra := byte(imageID >> 24)
	for row := 0; row < rows; row++ {
		var line strings.Builder
		line.WriteString(colorStart)
		for col := 0; col < cols; col++ {
			if col < kittyDiacriticLimit {
				// 可编码范围内显式写入 row/column/id，避免普通宽度下依赖左侧 cell 继承导致错位。
				line.WriteString(kittyPlaceholderCell(uint16(row), uint16(col), idExtra))
				continue
			}
			// 官方 diacritic 表只有有限长度；超宽终端的剩余列按协议从左侧 placeholder 继承递增。
			line.WriteString(termimg.PLACEHOLDER_CHAR)
		}
		line.WriteString("\x1b[39m")
		lines = append(lines, line.String())
	}
	return lines
}

func kittyPlaceholderCell(row uint16, column uint16, idExtra byte) string {
	cell := termimg.CreatePlaceholder(row, column, idExtra)
	if idExtra != 0 {
		return cell
	}
	runes := []rune(cell)
	if len(runes) <= 3 {
		return cell
	}
	return string(runes[:len(runes)-1])
}

func kittyPlaceholderColorStart(imageID uint32) string {
	r := (imageID >> 16) & 0xFF
	g := (imageID >> 8) & 0xFF
	b := imageID & 0xFF
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
}

func wrapKittyTmuxPassthrough(seq string) string {
	if os.Getenv("TMUX") == "" && os.Getenv("TERM_PROGRAM") != "tmux" {
		return seq
	}
	return "\x1bPtmux;\x1b" + strings.ReplaceAll(seq, "\x1b", "\x1b\x1b") + "\x1b\\"
}
