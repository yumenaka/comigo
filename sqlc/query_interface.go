package sqlc

import "context"

// bookQueries 是 StoreDatabase 实际使用的最小查询集合。
// SQLite 的 *Queries 直接满足，PostgreSQL 通过 adapter 转换生成类型后满足。
type bookQueries interface {
	GetBookByID(ctx context.Context, bookID string) (Book, error)
	CreateBook(ctx context.Context, arg CreateBookParams) (Book, error)
	UpdateBook(ctx context.Context, arg UpdateBookParams) error
	ListBooks(ctx context.Context) ([]Book, error)
	ListAllBookStoreURLs(ctx context.Context) ([]string, error)
	ListBooksByStorePath(ctx context.Context, storeUrl string) ([]Book, error)
	DeleteBook(ctx context.Context, bookID string) error
	GetPageInfosByBookID(ctx context.Context, bookID string) ([]PageInfo, error)
	DeletePageInfosByBookID(ctx context.Context, bookID string) error
	CreatePageInfo(ctx context.Context, arg CreatePageInfoParams) (PageInfo, error)
	ListBookmarksByBookID(ctx context.Context, bookID string) ([]Bookmark, error)
	DeleteBookmarksByBookID(ctx context.Context, bookID string) error
	CreateBookmark(ctx context.Context, arg CreateBookmarkParams) (Bookmark, error)
}
