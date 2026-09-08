package manual

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yumenaka/comigo/assets"
)

// 验证各语言的章节路径及未知路径处理。
func TestParsePath(t *testing.T) {
	tests := []struct {
		path, language, page string
		ok                   bool
	}{
		{"/manual/", "zh", "index", true},
		{"/manual/install", "zh", "install", true},
		{"/manual/ja-JP/reading", "ja-JP", "reading", true},
		{"/manual/en-US/faq", "en-US", "faq", true},
		{"/manual/comigo-omarchy", "zh", "comigo-omarchy", true},
		{"/manual/en-US/comigo-omarchy", "en-US", "comigo-omarchy", true},
		{"/manual/ja-JP/comigo-omarchy", "ja-JP", "comigo-omarchy", true},
		{"/manual/unknown", "", "", false},
		{"/manual/en-US/faq/extra", "", "", false},
	}
	for _, tt := range tests {
		language, page, ok := parsePath(tt.path)
		if language != tt.language || page != tt.page || ok != tt.ok {
			t.Fatalf("parsePath(%q) = %q, %q, %v", tt.path, language, page, ok)
		}
	}
}

// 验证内嵌的三种语言均包含完整章节。
func TestEmbeddedContent(t *testing.T) {
	for _, file := range []string{"zh_CN.json", "ja_JP.json", "en_US.json"} {
		data := assets.GetData("static/manual-content/" + file)
		var content struct {
			Pages map[string]struct {
				Title string `json:"title"`
				HTML  string `json:"html"`
			} `json:"pages"`
		}
		if err := json.Unmarshal(data, &content); err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		if len(content.Pages) != len(validPages) {
			t.Fatalf("%s: got %d pages", file, len(content.Pages))
		}
		for page, entry := range content.Pages {
			if entry.Title == "" || entry.HTML == "" {
				t.Fatalf("%s/%s: empty title or content", file, page)
			}
		}
		if strings.Contains(string(data), "data:image/") {
			t.Fatalf("%s embeds an image", file)
		}
	}
}
