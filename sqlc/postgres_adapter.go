package sqlc

import (
	"context"

	"github.com/yumenaka/comigo/sqlc/postgres"
)

type postgresAdapter struct {
	queries *postgres.Queries
}

// newPostgresAdapter 把 PostgreSQL 生成代码适配到顶层 sqlc 规范类型。
func newPostgresAdapter(queries *postgres.Queries) *postgresAdapter {
	return &postgresAdapter{queries: queries}
}

func (a *postgresAdapter) GetBookByID(ctx context.Context, bookID string) (Book, error) {
	book, err := a.queries.GetBookByID(ctx, bookID)
	return Book(book), err
}

func (a *postgresAdapter) CreateBook(ctx context.Context, arg CreateBookParams) (Book, error) {
	book, err := a.queries.CreateBook(ctx, postgres.CreateBookParams(arg))
	return Book(book), err
}

func (a *postgresAdapter) UpdateBook(ctx context.Context, arg UpdateBookParams) error {
	return a.queries.UpdateBook(ctx, postgres.UpdateBookParams(arg))
}

func (a *postgresAdapter) ListBooks(ctx context.Context) ([]Book, error) {
	books, err := a.queries.ListBooks(ctx)
	return fromPostgresBooks(books), err
}

func (a *postgresAdapter) ListAllBookStoreURLs(ctx context.Context) ([]string, error) {
	return a.queries.ListAllBookStoreURLs(ctx)
}

func (a *postgresAdapter) ListBooksByStorePath(ctx context.Context, storeUrl string) ([]Book, error) {
	books, err := a.queries.ListBooksByStorePath(ctx, storeUrl)
	return fromPostgresBooks(books), err
}

func (a *postgresAdapter) DeleteBook(ctx context.Context, bookID string) error {
	return a.queries.DeleteBook(ctx, bookID)
}

func (a *postgresAdapter) GetPageInfosByBookID(ctx context.Context, bookID string) ([]PageInfo, error) {
	pageInfos, err := a.queries.GetPageInfosByBookID(ctx, bookID)
	return fromPostgresPageInfos(pageInfos), err
}

func (a *postgresAdapter) DeletePageInfosByBookID(ctx context.Context, bookID string) error {
	return a.queries.DeletePageInfosByBookID(ctx, bookID)
}

func (a *postgresAdapter) CreatePageInfo(ctx context.Context, arg CreatePageInfoParams) (PageInfo, error) {
	pageInfo, err := a.queries.CreatePageInfo(ctx, postgres.CreatePageInfoParams(arg))
	return PageInfo(pageInfo), err
}

func (a *postgresAdapter) ListBookmarksByBookID(ctx context.Context, bookID string) ([]Bookmark, error) {
	bookmarks, err := a.queries.ListBookmarksByBookID(ctx, bookID)
	return fromPostgresBookmarks(bookmarks), err
}

func (a *postgresAdapter) DeleteBookmarksByBookID(ctx context.Context, bookID string) error {
	return a.queries.DeleteBookmarksByBookID(ctx, bookID)
}

func (a *postgresAdapter) CreateBookmark(ctx context.Context, arg CreateBookmarkParams) (Bookmark, error) {
	bookmark, err := a.queries.CreateBookmark(ctx, postgres.CreateBookmarkParams(arg))
	return Bookmark(bookmark), err
}

func fromPostgresBooks(books []postgres.Book) []Book {
	result := make([]Book, len(books))
	for i, book := range books {
		result[i] = Book(book)
	}
	return result
}

func fromPostgresPageInfos(pageInfos []postgres.PageInfo) []PageInfo {
	result := make([]PageInfo, len(pageInfos))
	for i, pageInfo := range pageInfos {
		result[i] = PageInfo(pageInfo)
	}
	return result
}

func fromPostgresBookmarks(bookmarks []postgres.Bookmark) []Bookmark {
	result := make([]Bookmark, len(bookmarks))
	for i, bookmark := range bookmarks {
		result[i] = Bookmark(bookmark)
	}
	return result
}
