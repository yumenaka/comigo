package tools

import (
	"bytes"
	"image"
	"image/png"
	"testing"
	"unicode/utf8"
)

// TestGetAuthorUnicode 覆盖全角、半角及已有的混用括号支持。
func TestGetAuthorUnicode(t *testing.T) {
	for _, input := range []string{"[作者]漫画", "【作者】漫画", "［作者］漫画", "［作者]漫画", "[作者］漫画", "（作者）漫画"} {
		if got := GetAuthor(input); got != "作者" || !utf8.ValidString(got) {
			t.Errorf("%q -> %q", input, got)
		}
	}
}

// TestBlurHashImageFailureKeepsOriginal 解码失败和不支持的 components 都应保留原图。
func TestBlurHashImageFailureKeepsOriginal(t *testing.T) {
	var pngData bytes.Buffer
	if err := png.Encode(&pngData, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		data       []byte
		components int
	}{{[]byte("bad image"), 2}, {pngData.Bytes(), 10}} {
		if got := GetImageDataBlurHashImage(tc.data, tc.components); !bytes.Equal(got, tc.data) {
			t.Fatal("处理失败没有保留原图")
		}
	}
}
