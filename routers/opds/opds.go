package opds

import (
	"encoding/xml"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/routers/apiresp"
	"github.com/yumenaka/comigo/tools"
)

const (
	opdsNavigationContentType  = `application/atom+xml;profile=opds-catalog;kind=navigation`
	opdsAcquisitionContentType = `application/atom+xml;profile=opds-catalog;kind=acquisition`
	opdsAcquisitionRel         = "http://opds-spec.org/acquisition/open-access"
	opdsThumbnailRel           = "http://opds-spec.org/image/thumbnail"
)

type feed struct {
	XMLName xml.Name `xml:"feed"`
	Xmlns   string   `xml:"xmlns,attr"`
	XmlnsDC string   `xml:"xmlns:dc,attr"`
	ID      string   `xml:"id"`
	Title   string   `xml:"title"`
	Updated string   `xml:"updated"`
	Links   []link   `xml:"link,omitempty"`
	Entries []entry  `xml:"entry,omitempty"`
}

type entry struct {
	ID           string  `xml:"id"`
	Title        string  `xml:"title"`
	Updated      string  `xml:"updated"`
	Author       *author `xml:"author,omitempty"`
	Summary      string  `xml:"summary,omitempty"`
	DCIdentifier string  `xml:"dc:identifier,omitempty"`
	Links        []link  `xml:"link,omitempty"`
}

type author struct {
	Name string `xml:"name"`
}

type link struct {
	Rel   string `xml:"rel,attr,omitempty"`
	Href  string `xml:"href,attr"`
	Type  string `xml:"type,attr,omitempty"`
	Title string `xml:"title,attr,omitempty"`
}

// RootHandler 输出 OPDS 1.2 导航 feed，入口只链接到 acquisition feed，避免混合 feed 类型。
func RootHandler(c echo.Context) error {
	books, err := model.IStore.ListBooks()
	if err != nil {
		return apiresp.Error(c, http.StatusInternalServerError, "opds_list_books_failed", err.Error(), nil)
	}
	bookGroups := topLevelBookGroups(books)
	updated := newestModified(books)
	if updated.IsZero() {
		updated = time.Now()
	}

	f := newFeed(c, "urn:comigo:opds:root", "ComiGo OPDS", updated, opdsNavigationContentType)
	f.Entries = append(f.Entries, entry{
		ID:      "urn:comigo:opds:all-books",
		Title:   "All Books",
		Updated: formatOPDSTime(updated),
		Links: []link{{
			Rel:   "subsection",
			Href:  absolutePath(c, "/opds/books"),
			Type:  opdsAcquisitionContentType,
			Title: "All Books",
		}},
	})
	for _, group := range bookGroups {
		f.Entries = append(f.Entries, navigationEntry(c, group))
	}
	return renderFeed(c, f, opdsNavigationContentType)
}

// BooksHandler 输出 acquisition feed。未指定书组时返回所有非书组书籍，指定 id 时递归展开该书组。
func BooksHandler(c echo.Context) error {
	parentID := c.Param("id")
	var books []*model.Book
	var title string
	switch parentID {
	case "":
		allBooks, err := model.IStore.ListBooks()
		if err != nil {
			return apiresp.Error(c, http.StatusInternalServerError, "opds_list_books_failed", err.Error(), nil)
		}
		books = nonGroupBooks(allBooks)
		title = "All Books"
	default:
		parent, err := model.IStore.GetBook(parentID)
		if err != nil || parent.Type != model.TypeBooksGroup {
			return apiresp.Error(c, http.StatusNotFound, "opds_group_not_found", "OPDS group not found", map[string]string{"id": parentID})
		}
		books = collectGroupBooks(parent, map[string]bool{})
		title = parent.Title
	}
	sortBooksByTitle(books)

	updated := newestModified(books)
	if updated.IsZero() {
		updated = time.Now()
	}
	f := newFeed(c, "urn:comigo:opds:books:"+parentID, title, updated, opdsAcquisitionContentType)
	f.Links = append(f.Links, link{
		Rel:   "start",
		Href:  absolutePath(c, "/opds"),
		Type:  opdsNavigationContentType,
		Title: "ComiGo OPDS",
	})
	for _, book := range books {
		if bookEntry, ok := acquisitionEntry(c, book); ok {
			f.Entries = append(f.Entries, bookEntry)
		}
	}
	return renderFeed(c, f, opdsAcquisitionContentType)
}

func newFeed(c echo.Context, id, title string, updated time.Time, contentType string) feed {
	return feed{
		Xmlns:   "http://www.w3.org/2005/Atom",
		XmlnsDC: "http://purl.org/dc/terms/",
		ID:      id,
		Title:   title,
		Updated: formatOPDSTime(updated),
		Links: []link{{
			Rel:  "self",
			Href: absolutePath(c, config.StripBasePath(c.Request().URL.RequestURI())),
			Type: contentType,
		}},
	}
}

func navigationEntry(c echo.Context, book *model.Book) entry {
	updated := book.Modified
	if updated.IsZero() {
		updated = time.Now()
	}
	return entry{
		ID:      "urn:comigo:opds:group:" + book.BookID,
		Title:   book.Title,
		Updated: formatOPDSTime(updated),
		Links: []link{{
			Rel:   "subsection",
			Href:  absolutePath(c, "/opds/books/"+url.PathEscape(book.BookID)),
			Type:  opdsAcquisitionContentType,
			Title: book.Title,
		}},
	}
}

func acquisitionEntry(c echo.Context, book *model.Book) (entry, bool) {
	acquisitionHref, contentType, ok := acquisitionLink(c, book)
	if !ok {
		return entry{}, false
	}
	updated := book.Modified
	if updated.IsZero() {
		updated = time.Now()
	}
	bookAuthor := strings.TrimSpace(book.Author)
	if bookAuthor == "" {
		bookAuthor = "ComiGo"
	}
	return entry{
		ID:           "urn:comigo:book:" + book.BookID,
		Title:        book.Title,
		Updated:      formatOPDSTime(updated),
		Author:       &author{Name: bookAuthor},
		Summary:      bookSummary(book),
		DCIdentifier: book.BookID,
		Links: []link{
			{
				Rel:  opdsAcquisitionRel,
				Href: acquisitionHref,
				Type: contentType,
			},
			{
				Rel:  opdsThumbnailRel,
				Href: absolutePath(c, "/api/get-cover?id="+url.QueryEscape(book.BookID)+"&resize_height=352"),
				Type: "image/jpeg",
			},
		},
	}, true
}

func acquisitionLink(c echo.Context, book *model.Book) (href string, contentType string, ok bool) {
	switch book.Type {
	case model.TypeDir:
		return absolutePath(c, "/api/download-zip?id="+url.QueryEscape(book.BookID)), "application/zip", true
	case model.TypeZip:
		contentType = "application/zip"
	case model.TypeCbz:
		contentType = "application/vnd.comicbook+zip"
	case model.TypeRar:
		contentType = "application/vnd.rar"
	case model.TypeCbr:
		contentType = "application/vnd.comicbook-rar"
	case model.TypeTar:
		contentType = "application/x-tar"
	case model.TypeEpub:
		contentType = "application/epub+zip"
	case model.TypePDF:
		contentType = "application/pdf"
	case model.TypeHTML:
		contentType = "text/html"
	default:
		return "", "", false
	}
	fileName := filepath.Base(book.BookPath)
	return absolutePath(c, "/api/raw/"+url.PathEscape(book.BookID)+"/"+url.PathEscape(fileName)), contentType, true
}

func bookSummary(book *model.Book) string {
	parts := []string{string(book.Type)}
	if book.PageCount > 0 {
		parts = append(parts, strconv.Itoa(book.PageCount)+" pages")
	}
	if book.ParentFolder != "" {
		parts = append(parts, book.ParentFolder)
	}
	return strings.Join(parts, " / ")
}

func topLevelBookGroups(books []*model.Book) []*model.Book {
	groups := make([]*model.Book, 0)
	for _, book := range books {
		if book.Type == model.TypeBooksGroup && book.Depth == 0 {
			groups = append(groups, book)
		}
	}
	sortBooksByTitle(groups)
	return groups
}

func nonGroupBooks(books []*model.Book) []*model.Book {
	result := make([]*model.Book, 0, len(books))
	for _, book := range books {
		if book.Type != model.TypeBooksGroup {
			result = append(result, book)
		}
	}
	return result
}

func collectGroupBooks(group *model.Book, seen map[string]bool) []*model.Book {
	if group == nil || seen[group.BookID] {
		return nil
	}
	seen[group.BookID] = true
	result := make([]*model.Book, 0, len(group.ChildBooksID))
	for _, id := range group.ChildBooksID {
		child, err := model.IStore.GetBook(id)
		if err != nil || child == nil {
			continue
		}
		if child.Type == model.TypeBooksGroup {
			result = append(result, collectGroupBooks(child, seen)...)
			continue
		}
		result = append(result, child)
	}
	return result
}

func sortBooksByTitle(books []*model.Book) {
	sort.Slice(books, func(i, j int) bool {
		return tools.Compare(books[i].Title, books[j].Title)
	})
}

func newestModified(books []*model.Book) time.Time {
	var newest time.Time
	for _, book := range books {
		if book.Modified.After(newest) {
			newest = book.Modified
		}
	}
	return newest
}

func absolutePath(c echo.Context, relativePath string) string {
	if strings.HasPrefix(relativePath, "http://") || strings.HasPrefix(relativePath, "https://") {
		return relativePath
	}
	scheme := c.Scheme()
	if scheme == "" {
		scheme = "http"
	}
	if !strings.HasPrefix(relativePath, "/") {
		relativePath = "/" + relativePath
	}
	return scheme + "://" + c.Request().Host + config.PrefixPath(relativePath)
}

func formatOPDSTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func renderFeed(c echo.Context, f feed, contentType string) error {
	payload, err := xml.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	payload = append([]byte(xml.Header), payload...)
	return c.Blob(http.StatusOK, contentType, payload)
}
