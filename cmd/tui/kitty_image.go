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
	xdraw "golang.org/x/image/draw"
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

	placementImage := rasterizeKittyPlacementImage(img, cols, rows)
	var buf bytes.Buffer
	if err := png.Encode(&buf, placementImage); err != nil {
		return "", nil, fmt.Errorf("encode kitty png: %w", err)
	}

	imageID := nextTUIKittyImageID()
	setup := renderKittyUnicodeSetup(buf.Bytes(), imageID, cols, rows)
	lines := renderKittyPlaceholderLines(imageID, cols, rows)
	return setup, lines, nil
}

// rasterizeKittyPlacementImage 把图片预先缩放到 virtual placement 的等效像素矩形。
// Kitty 会在 c/r 矩形内保持源图比例 fit；这里让 PNG 自身比例先匹配字符格矩形，避免终端二次居中造成横向或纵向偏移。
func rasterizeKittyPlacementImage(img image.Image, cols int, rows int) image.Image {
	cellW, cellH := protocolCellPixels(termimg.Kitty)
	targetW := max(1, cols*cellW)
	targetH := max(1, rows*cellH)
	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), xdraw.Src, nil)
	return dst
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

// kittySetupKey 提取首段中的图片 ID 用于去重，兼容 tmux 包装且避免额外持有整份图片数据。
func kittySetupKey(setup string) string {
	_, id, _ := strings.Cut(setup, ",i=")
	id, _, _ = strings.Cut(id, ",")
	return id
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

// kittyPlaceholderCell 使用官方行列表；go-termimg v0.1.26 从索引 62 起编码不符。
func kittyPlaceholderCell(row uint16, column uint16, idExtra byte) string {
	cell := []rune{0x10eeee, kittyDiacritics[row], kittyDiacritics[column]}
	if idExtra != 0 {
		cell = append(cell, kittyDiacritics[idExtra])
	}
	return string(cell)
}

// Kitty 官方 Unicode 占位符行列表，顺序是协议的一部分，不能替换为通用组合字符表。
// https://github.com/kovidgoyal/kitty/blob/master/gen/rowcolumn-diacritics.txt
var kittyDiacritics = [...]rune{
	0x0305, 0x030d, 0x030e, 0x0310, 0x0312, 0x033d, 0x033e, 0x033f, 0x0346, 0x034a,
	0x034b, 0x034c, 0x0350, 0x0351, 0x0352, 0x0357, 0x035b, 0x0363, 0x0364, 0x0365,
	0x0366, 0x0367, 0x0368, 0x0369, 0x036a, 0x036b, 0x036c, 0x036d, 0x036e, 0x036f,
	0x0483, 0x0484, 0x0485, 0x0486, 0x0487, 0x0592, 0x0593, 0x0594, 0x0595, 0x0597,
	0x0598, 0x0599, 0x059c, 0x059d, 0x059e, 0x059f, 0x05a0, 0x05a1, 0x05a8, 0x05a9,
	0x05ab, 0x05ac, 0x05af, 0x05c4, 0x0610, 0x0611, 0x0612, 0x0613, 0x0614, 0x0615,
	0x0616, 0x0617, 0x0657, 0x0658, 0x0659, 0x065a, 0x065b, 0x065d, 0x065e, 0x06d6,
	0x06d7, 0x06d8, 0x06d9, 0x06da, 0x06db, 0x06dc, 0x06df, 0x06e0, 0x06e1, 0x06e2,
	0x06e4, 0x06e7, 0x06e8, 0x06eb, 0x06ec, 0x0730, 0x0732, 0x0733, 0x0735, 0x0736,
	0x073a, 0x073d, 0x073f, 0x0740, 0x0741, 0x0743, 0x0745, 0x0747, 0x0749, 0x074a,
	0x07eb, 0x07ec, 0x07ed, 0x07ee, 0x07ef, 0x07f0, 0x07f1, 0x07f3, 0x0816, 0x0817,
	0x0818, 0x0819, 0x081b, 0x081c, 0x081d, 0x081e, 0x081f, 0x0820, 0x0821, 0x0822,
	0x0823, 0x0825, 0x0826, 0x0827, 0x0829, 0x082a, 0x082b, 0x082c, 0x082d, 0x0951,
	0x0953, 0x0954, 0x0f82, 0x0f83, 0x0f86, 0x0f87, 0x135d, 0x135e, 0x135f, 0x17dd,
	0x193a, 0x1a17, 0x1a75, 0x1a76, 0x1a77, 0x1a78, 0x1a79, 0x1a7a, 0x1a7b, 0x1a7c,
	0x1b6b, 0x1b6d, 0x1b6e, 0x1b6f, 0x1b70, 0x1b71, 0x1b72, 0x1b73, 0x1cd0, 0x1cd1,
	0x1cd2, 0x1cda, 0x1cdb, 0x1ce0, 0x1dc0, 0x1dc1, 0x1dc3, 0x1dc4, 0x1dc5, 0x1dc6,
	0x1dc7, 0x1dc8, 0x1dc9, 0x1dcb, 0x1dcc, 0x1dd1, 0x1dd2, 0x1dd3, 0x1dd4, 0x1dd5,
	0x1dd6, 0x1dd7, 0x1dd8, 0x1dd9, 0x1dda, 0x1ddb, 0x1ddc, 0x1ddd, 0x1dde, 0x1ddf,
	0x1de0, 0x1de1, 0x1de2, 0x1de3, 0x1de4, 0x1de5, 0x1de6, 0x1dfe, 0x20d0, 0x20d1,
	0x20d4, 0x20d5, 0x20d6, 0x20d7, 0x20db, 0x20dc, 0x20e1, 0x20e7, 0x20e9, 0x20f0,
	0x2cef, 0x2cf0, 0x2cf1, 0x2de0, 0x2de1, 0x2de2, 0x2de3, 0x2de4, 0x2de5, 0x2de6,
	0x2de7, 0x2de8, 0x2de9, 0x2dea, 0x2deb, 0x2dec, 0x2ded, 0x2dee, 0x2def, 0x2df0,
	0x2df1, 0x2df2, 0x2df3, 0x2df4, 0x2df5, 0x2df6, 0x2df7, 0x2df8, 0x2df9, 0x2dfa,
	0x2dfb, 0x2dfc, 0x2dfd, 0x2dfe, 0x2dff, 0xa66f, 0xa67c, 0xa67d, 0xa6f0, 0xa6f1,
	0xa8e0, 0xa8e1, 0xa8e2, 0xa8e3, 0xa8e4, 0xa8e5, 0xa8e6, 0xa8e7, 0xa8e8, 0xa8e9,
	0xa8ea, 0xa8eb, 0xa8ec, 0xa8ed, 0xa8ee, 0xa8ef, 0xa8f0, 0xa8f1, 0xaab0, 0xaab2,
	0xaab3, 0xaab7, 0xaab8, 0xaabe, 0xaabf, 0xaac1, 0xfe20, 0xfe21, 0xfe22, 0xfe23,
	0xfe24, 0xfe25, 0xfe26, 0x10a0f, 0x10a38, 0x1d185, 0x1d186, 0x1d187, 0x1d188, 0x1d189,
	0x1d1aa, 0x1d1ab, 0x1d1ac, 0x1d1ad, 0x1d242, 0x1d243, 0x1d244,
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
