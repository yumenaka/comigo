package sqlc

import (
	"database/sql"
	"time"

	"github.com/yumenaka/comigo/model"
)

// ==================== Bookmark 相关转换 ====================

// FromSQLCBookmark 将 sqlc.Bookmark 转换为 model.BookMark
func FromSQLCBookmark(sqlcBookmark Bookmark) model.BookMark {
	var createdAt time.Time
	if sqlcBookmark.CreatedAt.Valid {
		createdAt = sqlcBookmark.CreatedAt.Time
	}
	var updatedAt time.Time
	if sqlcBookmark.UpdatedAt.Valid {
		updatedAt = sqlcBookmark.UpdatedAt.Time
	}

	return model.BookMark{
		BookID:      sqlcBookmark.BookID,
		PageIndex:   int(sqlcBookmark.PageIndex),
		Description: sqlcBookmark.Description.String,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// FromSQLCBookmarks 批量转换 sqlc.Bookmark 为 model.BookMarks
func FromSQLCBookmarks(sqlcBookmarks []Bookmark) model.BookMarks {
	bookmarks := make(model.BookMarks, len(sqlcBookmarks))
	for i, b := range sqlcBookmarks {
		bookmarks[i] = FromSQLCBookmark(b)
	}
	return bookmarks
}

// ToSQLCCreateBookmarkParams 将 model.BookMark 转换为 sqlc.CreateBookmarkParams
func ToSQLCCreateBookmarkParams(bookmark model.BookMark) CreateBookmarkParams {
	return CreateBookmarkParams{
		BookID:      bookmark.BookID,
		PageIndex:   int64(bookmark.PageIndex),
		Description: sql.NullString{String: bookmark.Description, Valid: bookmark.Description != ""},
	}
}

// ToSQLCUpdateBookmarkParams 将 model.BookMark 转换为 sqlc.UpdateBookmarkParams
func ToSQLCUpdateBookmarkParams(bookmark model.BookMark) UpdateBookmarkParams {
	return UpdateBookmarkParams{
		Type:        int64(bookmark.PageIndex),
		BookID:      bookmark.BookID,
		PageIndex:   int64(bookmark.PageIndex),
		Description: sql.NullString{String: bookmark.Description, Valid: bookmark.Description != ""},
	}
}
