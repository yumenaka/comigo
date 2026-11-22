package data_api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/model"
)

// bookmarkTestStore 只记录请求产生的持久化副作用。
type bookmarkTestStore struct {
	model.StoreInterface
	book            *model.Book
	writes, deletes int
	listErr         error
}

// GetBook 模拟书籍不存在的错误。
func (s *bookmarkTestStore) GetBook(id string) (*model.Book, error) {
	if id != s.book.BookID {
		return nil, errors.New("book not found")
	}
	return s.book, nil
}

// StoreBookMark 记录写入次数。
func (s *bookmarkTestStore) StoreBookMark(*model.BookMark) error { s.writes++; return nil }

// DeleteBookMark 记录删除次数。
func (s *bookmarkTestStore) DeleteBookMark(string, model.MarkType, int) error {
	s.deletes++
	return nil
}

// ListBooks 返回共享对象，检查输出封面时是否修改原对象。
func (s *bookmarkTestStore) ListBooks() ([]*model.Book, error) {
	return []*model.Book{s.book}, s.listErr
}

// TestBookmarkValidationStopsMutation 保证失败响应只有一个 JSON，且没有存储副作用。
func TestBookmarkValidationStopsMutation(t *testing.T) {
	old := model.IStore
	defer func() { model.IStore = old }()
	for _, tc := range []struct{ name, method, query, body, code string }{
		{"invalid JSON", "POST", "", `{`, "invalid_json"},
		{"missing book", "POST", "", `{"type":"user","page_index":1}`, "missing_book_id"},
		{"unknown book", "POST", "", `{"type":"user","book_id":"missing","page_index":1}`, "book_not_found"},
		{"unknown type", "POST", "", `{"type":"invalid","book_id":"book","page_index":1}`, "invalid_mark_type"},
		{"missing type", "POST", "", `{"book_id":"book","page_index":1}`, "missing_mark_type"},
		{"zero page", "POST", "", `{"type":"user","book_id":"book","page_index":0}`, "missing_page_index"},
		{"past last page", "POST", "", `{"type":"user","book_id":"book","page_index":2}`, "page_index_out_of_range"},
		{"delete missing page", "DELETE", "?book_id=book&mark_type=user", "", "missing_page_index"},
		{"delete invalid page", "DELETE", "?book_id=book&mark_type=user&page_index=x", "", "invalid_page_index"},
		{"delete invalid type", "DELETE", "?book_id=book&mark_type=invalid&page_index=1", "", "invalid_mark_type"},
		{"delete valid", "DELETE", "?book_id=book&mark_type=user&page_index=1", "", ""},
		{"store valid", "POST", "", `{"type":"user","book_id":"book","page_index":1}`, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := &bookmarkTestStore{book: &model.Book{BookInfo: model.BookInfo{BookID: "book", PageCount: 1}}}
			model.IStore = s
			req := httptest.NewRequest(tc.method, "/api/bookmarks"+tc.query, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			handler := StoreBookmark
			if tc.method == "DELETE" {
				handler = DeleteBookmark
			}
			if err := handler(c); err != nil {
				t.Fatal(err)
			}
			if tc.code != "" {
				if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), `"code":"`+tc.code+`"`) || s.writes+s.deletes != 0 || strings.Count(rec.Body.String(), `"code":`) != 1 {
					t.Fatalf("拒绝请求仍有副作用或响应错误: %d %s writes=%d deletes=%d", rec.Code, rec.Body, s.writes, s.deletes)
				}
			} else if rec.Code != http.StatusOK || s.writes+s.deletes != 1 {
				t.Fatalf("有效请求失败: %d %s", rec.Code, rec.Body)
			}
		})
	}
}

// TestBookmarkViewsDoNotMutateBook 保留共享书籍的封面，同时在响应副本里生成封面。
func TestBookmarkViewsDoNotMutateBook(t *testing.T) {
	old := model.IStore
	defer func() { model.IStore = old }()
	book := &model.Book{BookInfo: model.BookInfo{BookID: "book", Type: model.TypePDF}, BookMarks: []model.BookMark{{PageIndex: 1, Type: model.UserMark}}}
	model.IStore = &bookmarkTestStore{book: book}
	for _, handler := range []echo.HandlerFunc{GetAllBookmarks, GetReadingHistory} {
		rec := httptest.NewRecorder()
		if err := handler(echo.New().NewContext(httptest.NewRequest("GET", "/", nil), rec)); err != nil {
			t.Fatal(err)
		}
		if book.Cover.Name != "" || !strings.Contains(rec.Body.String(), "1.jpg") {
			t.Fatalf("封面副本错误: book=%+v response=%s", book.Cover, rec.Body)
		}
	}
}

// TestBookmarkListFailure 查询失败应能与正常的空书库区分。
func TestBookmarkListFailure(t *testing.T) {
	old := model.IStore
	defer func() { model.IStore = old }()
	model.IStore = &bookmarkTestStore{listErr: errors.New("failed")}
	for _, handler := range []echo.HandlerFunc{GetAllBookmarks, GetReadingHistory} {
		r := httptest.NewRecorder()
		err := handler(echo.New().NewContext(httptest.NewRequest("GET", "/", nil), r))
		if err != nil || r.Code != 500 || !strings.Contains(r.Body.String(), "list_books_failed") {
			t.Fatalf("查询失败响应错误: %v %d %s", err, r.Code, r.Body)
		}
	}
}
