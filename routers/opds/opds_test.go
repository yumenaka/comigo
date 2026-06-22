package opds

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
)

type fakeStore struct {
	books []*model.Book
}

func (s *fakeStore) StoreBook(b *model.Book) error { return nil }

func (s *fakeStore) GetBook(id string) (*model.Book, error) {
	for _, book := range s.books {
		if book.BookID == id {
			return book, nil
		}
	}
	return nil, errors.New("book not found")
}

func (s *fakeStore) DeleteBook(id string) error { return nil }

func (s *fakeStore) ListBooks() ([]*model.Book, error) { return s.books, nil }

func (s *fakeStore) GenerateBookGroup() error { return nil }

func (s *fakeStore) StoreBookMark(mark *model.BookMark) error { return nil }

func (s *fakeStore) GetBookMarks(bookID string) (*model.BookMarks, error) {
	bookMarks := model.BookMarks{}
	return &bookMarks, nil
}

func (s *fakeStore) DeleteBookMark(bookID string, markType model.MarkType, pageIndex int) error {
	return nil
}

func withOPDSTestStore(t *testing.T, store *fakeStore) {
	t.Helper()
	oldStore := model.IStore
	model.IStore = store
	t.Cleanup(func() {
		model.IStore = oldStore
	})
}

// 验证 OPDS 根入口返回导航订阅源。
func TestRootHandlerReturnsNavigationFeed(t *testing.T) {
	withOPDSTestStore(t, &fakeStore{books: testBooks()})
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "http://example.test/opds", nil)
	rec := httptest.NewRecorder()

	if err := RootHandler(e.NewContext(req, rec)); err != nil {
		t.Fatalf("RootHandler error: %v", err)
	}

	body := rec.Body.String()
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, body)
	}
	if contentType := rec.Header().Get(echo.HeaderContentType); !strings.Contains(contentType, "kind=navigation") {
		t.Fatalf("Content-Type = %q, want OPDS navigation", contentType)
	}
	for _, want := range []string{
		`<title>ComiGo OPDS</title>`,
		`href="http://example.test/opds/books"`,
		`href="http://example.test/opds/books/group1"`,
		`type="application/atom+xml;profile=opds-catalog;kind=acquisition"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("navigation feed missing %q\n%s", want, body)
		}
	}
}

// 验证 OPDS 书籍接口返回可获取书籍的订阅源。
func TestBooksHandlerReturnsAcquisitionFeed(t *testing.T) {
	withOPDSTestStore(t, &fakeStore{books: testBooks()})
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "http://example.test/opds/books/group1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("group1")

	if err := BooksHandler(c); err != nil {
		t.Fatalf("BooksHandler error: %v", err)
	}

	body := rec.Body.String()
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, body)
	}
	if contentType := rec.Header().Get(echo.HeaderContentType); !strings.Contains(contentType, "kind=acquisition") {
		t.Fatalf("Content-Type = %q, want OPDS acquisition", contentType)
	}
	for _, want := range []string{
		`rel="http://opds-spec.org/acquisition/open-access"`,
		`href="http://example.test/api/raw/book1/Book%201.cbz"`,
		`href="http://example.test/api/download-zip?id=book2"`,
		`rel="http://opds-spec.org/image/thumbnail"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("acquisition feed missing %q\n%s", want, body)
		}
	}
}

// 验证 OPDS 链接会带上反向代理基础路径。
func TestOPDSLinksRespectBasePath(t *testing.T) {
	cfg := config.GetCfg()
	oldBasePath := cfg.BasePath
	cfg.BasePath = "/proxy"
	t.Cleanup(func() { cfg.BasePath = oldBasePath })

	withOPDSTestStore(t, &fakeStore{books: testBooks()})
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "http://example.test/proxy/opds", nil)
	rec := httptest.NewRecorder()

	if err := RootHandler(e.NewContext(req, rec)); err != nil {
		t.Fatalf("RootHandler error: %v", err)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `href="http://example.test/proxy/opds/books"`) {
		t.Fatalf("base path was not applied to OPDS links:\n%s", body)
	}
}

func testBooks() []*model.Book {
	modified := time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)
	return []*model.Book{
		{
			BookInfo: model.BookInfo{
				BookID:       "group1",
				Title:        "Series",
				Type:         model.TypeBooksGroup,
				Depth:        0,
				Modified:     modified,
				ChildBooksID: []string{"book1", "book2"},
			},
		},
		{
			BookInfo: model.BookInfo{
				BookID:       "book1",
				Title:        "Book 1",
				Author:       "Author",
				Type:         model.TypeCbz,
				BookPath:     "/library/Book 1.cbz",
				ParentFolder: "library",
				PageCount:    10,
				Modified:     modified.Add(time.Hour),
			},
		},
		{
			BookInfo: model.BookInfo{
				BookID:       "book2",
				Title:        "Book 2",
				Type:         model.TypeDir,
				BookPath:     "/library/Book 2",
				ParentFolder: "library",
				PageCount:    5,
				Modified:     modified.Add(2 * time.Hour),
			},
		},
	}
}
