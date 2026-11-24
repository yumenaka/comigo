package model

import (
	"testing"
	"time"
)

// TestBookIDBatchCollision 使用真实的七位摘要碰撞，验证未入库的新书也会占用 ID。
func TestBookIDBatchCollision(t *testing.T) {
	ids, err := NewBookIDIndex(nil)
	if err != nil {
		t.Fatal(err)
	}
	a, err := ids.NewBook("/library/book-513274.cbz", time.Time{}, 10, "/library", 0, TypeCbz)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ids.NewBook("/library/book-2231337.cbz", time.Time{}, 10, "/library", 0, TypeCbz)
	if err != nil {
		t.Fatal(err)
	}
	if a.BookID != "YSTdWZm" || len(b.BookID) != 8 || b.BookID[:7] != a.BookID {
		t.Fatalf("collision was not expanded: %q, %q", a.BookID, b.BookID)
	}
	// 第一册移除后，重新加载仍保留第二册的加长 ID。
	reloaded, err := NewBookIDIndex([]*Book{b})
	if err != nil {
		t.Fatal(err)
	}
	updated, err := reloaded.NewBook(b.BookPath, time.Now(), 10, b.StoreUrl, 0, b.Type)
	if err != nil || updated.BookID != b.BookID {
		t.Fatalf("rescan changed ID: book=%+v err=%v", updated, err)
	}
	changed, err := reloaded.NewBook(b.BookPath, time.Now(), 20, b.StoreUrl, 0, b.Type)
	if err != nil || changed.BookID == b.BookID {
		t.Fatalf("size change must invalidate cache ID: book=%+v err=%v", changed, err)
	}
}

// TestBookIDIndexRejectsConflictingSnapshot 不允许已有不同书籍的重复 ID 被静默合并。
func TestBookIDIndexRejectsConflictingSnapshot(t *testing.T) {
	a := &Book{BookInfo: BookInfo{BookID: "conflict", BookPath: "a"}}
	b := &Book{BookInfo: BookInfo{BookID: "conflict", BookPath: "b"}}
	if _, err := NewBookIDIndex([]*Book{a, b}); err == nil {
		t.Fatal("expected conflicting ID error")
	}
}

// TestBookIDExhaustion 完整摘要也被占用时必须报错，不能返回已存在的 ID。
func TestBookIDExhaustion(t *testing.T) {
	book := BookInfo{StoreUrl: "library", BookPath: "book.cbz", Type: TypeCbz}
	fullID := book.fullBookID()
	ids, err := NewBookIDIndex(nil)
	if err != nil {
		t.Fatal(err)
	}
	for length := 7; length <= len(fullID); length++ {
		ids.byID[fullID[:length]] = [3]string{"another", "book", ".cbz"}
	}
	if err := ids.Assign(&book); err == nil || book.BookID != "" {
		t.Fatalf("expected allocation failure, ID=%q err=%v", book.BookID, err)
	}
}
