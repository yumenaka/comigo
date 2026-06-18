package shelf

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/yumenaka/comigo/model"
)

func TestBookCardRendersWailsContextMenuHook(t *testing.T) {
	var html bytes.Buffer
	marks := model.BookMarks{}
	book := model.BookInfo{
		BookID:    "book-id",
		Title:     "book.zip",
		Type:      model.TypeZip,
		PageCount: 1,
	}

	if err := BookCard(nil, book, marks).Render(context.Background(), &html); err != nil {
		t.Fatalf("render BookCard: %v", err)
	}
	if !strings.Contains(html.String(), "data-wails-book-card") {
		t.Fatalf("BookCard missing Wails context menu hook in %s", html.String())
	}
}
