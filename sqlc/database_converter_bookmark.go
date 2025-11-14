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

	description := ""
	if sqlcBookmark.Description.Valid {
		description = sqlcBookmark.Description.String
	}

	return model.BookMark{
		Type:        model.MarkType(sqlcBookmark.Type),
		BookID:      sqlcBookmark.BookID,
		PageIndex:   int(sqlcBookmark.PageIndex),
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

// FromSQLCBookmarks 批量转换 sqlc.Bookmark 为 model.BookMarks
func FromSQLCBookmarks(sqlcBookmarks []Bookmark) model.BookMarks {
	if len(sqlcBookmarks) == 0 {
		return nil
	}
	bookmarks := make(model.BookMarks, len(sqlcBookmarks))
	for i, b := range sqlcBookmarks {
		bookmarks[i] = FromSQLCBookmark(b)
	}
	return bookmarks
}

// ToSQLCCreateBookmarkParams 将 model.BookMark 转换为 sqlc.CreateBookmarkParams
func ToSQLCCreateBookmarkParams(bookID string, bookmark model.BookMark) CreateBookmarkParams {
	resolvedBookID := bookmark.BookID
	if resolvedBookID == "" {
		resolvedBookID = bookID
	}
	markType := bookmark.Type
	if markType == "" {
		markType = model.UserMark
	}
	return CreateBookmarkParams{
		Type:        string(markType),
		BookID:      resolvedBookID,
		PageIndex:   int64(bookmark.PageIndex),
		Description: sql.NullString{String: bookmark.Description, Valid: bookmark.Description != ""},
	}
}

// ToSQLCUpdateBookmarkParams 将 model.BookMark 转换为 sqlc.UpdateBookmarkParams
func ToSQLCUpdateBookmarkParams(bookID string, bookmark model.BookMark) UpdateBookmarkParams {
	resolvedBookID := bookmark.BookID
	if resolvedBookID == "" {
		resolvedBookID = bookID
	}
	markType := bookmark.Type
	if markType == "" {
		markType = model.UserMark
	}
	description := sql.NullString{String: bookmark.Description, Valid: bookmark.Description != ""}
	return UpdateBookmarkParams{
		Description:   description,
		PageIndex:     int64(bookmark.PageIndex),
		Description_2: description,
		BookID:        resolvedBookID,
		Type:          string(markType),
	}
}
