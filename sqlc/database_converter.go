package sqlc

import (
	"database/sql"
	"strings"

	"github.com/yumenaka/comigo/model"
)

// ==================== Book 相关转换 ====================

// ConvertCommaSeparatedString 将英文逗号分隔的字符串转换为 []string
func ConvertCommaSeparatedString(nullString sql.NullString) []string {
	if !nullString.Valid {
		return []string{} // 如果字符串无效，返回空切片
	}
	// 使用 strings.Split 分割
	parts := strings.Split(nullString.String, ",")
	// 去除每个元素前后的空白字符
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// FromSQLCBook 将sqlc.Book转换为model.Book
func FromSQLCBook(sqlcBook Book) *model.Book {
	return &model.Book{
		BookInfo: model.BookInfo{
			Title:            sqlcBook.Title,
			BookID:           sqlcBook.BookID,
			Author:           sqlcBook.Author.String,
			StoreUrl:         sqlcBook.StoreUrl,
			ChildBooksNum:    int(sqlcBook.ChildBooksNum.Int64),
			ChildBooksID:     ConvertCommaSeparatedString(sqlcBook.ChildBooksID),
			Deleted:          sqlcBook.Deleted.Bool,
			Depth:            int(sqlcBook.Depth.Int64),
			ExtractPath:      sqlcBook.ExtractPath.String,
			ExtractNum:       int(sqlcBook.ExtractNum.Int64),
			BookPath:         sqlcBook.BookPath,
			FileSize:         sqlcBook.FileSize.Int64,
			ISBN:             sqlcBook.Isbn.String,
			BookComplete:     sqlcBook.BookComplete.Bool,
			InitComplete:     sqlcBook.InitComplete.Bool,
			Modified:         sqlcBook.ModifiedTime.Time,
			NonUTF8Zip:       sqlcBook.NonUtf8zip.Bool,
			PageCount:        int(sqlcBook.PageCount.Int64),
			ParentFolder:     sqlcBook.ParentFolder.String,
			Press:            sqlcBook.Press.String,
			PublishedAt:      sqlcBook.PublishedAt.String,
			Type:             model.SupportFileType(sqlcBook.Type),
			ZipTextEncoding:  sqlcBook.ZipTextEncoding.String,
			CreatedByVersion: sqlcBook.CreatedByVersion.String,
			IsRemote:         sqlcBook.IsRemote.Bool,
			RemoteURL:        sqlcBook.RemoteUrl.String,
		},
	}
}

// ToSQLCCreateBookParams 将model.Book转换为sqlc.CreateBookParams //"Valid"必须是验证条件或true
func ToSQLCCreateBookParams(book *model.Book) CreateBookParams {
	return CreateBookParams{
		Title:            book.Title,
		BookID:           book.BookID,
		Owner:            sql.NullString{String: "admin", Valid: true},
		BookPath:         book.BookPath,
		StoreUrl:         book.StoreUrl,
		Type:             string(book.Type),
		ChildBooksNum:    sql.NullInt64{Int64: int64(book.ChildBooksNum), Valid: true},
		ChildBooksID:     sql.NullString{String: strings.Join(book.ChildBooksID, ", "), Valid: len(book.ChildBooksID) > 0},
		Depth:            sql.NullInt64{Int64: int64(book.Depth), Valid: true},
		ParentFolder:     sql.NullString{String: book.ParentFolder, Valid: book.ParentFolder != ""},
		PageCount:        sql.NullInt64{Int64: int64(book.PageCount), Valid: true},
		FileSize:         sql.NullInt64{Int64: book.FileSize, Valid: true},
		Author:           sql.NullString{String: book.Author, Valid: book.Author != ""},
		Isbn:             sql.NullString{String: book.ISBN, Valid: book.ISBN != ""},
		Press:            sql.NullString{String: book.Press, Valid: book.Press != ""},
		PublishedAt:      sql.NullString{String: book.PublishedAt, Valid: book.PublishedAt != ""},
		ExtractPath:      sql.NullString{String: book.ExtractPath, Valid: book.ExtractPath != ""},
		ExtractNum:       sql.NullInt64{Int64: int64(book.ExtractNum), Valid: true},
		BookComplete:     sql.NullBool{Bool: book.BookComplete, Valid: true},
		InitComplete:     sql.NullBool{Bool: book.InitComplete, Valid: true},
		NonUtf8zip:       sql.NullBool{Bool: book.NonUTF8Zip, Valid: true},
		ZipTextEncoding:  sql.NullString{String: book.ZipTextEncoding, Valid: book.ZipTextEncoding != ""},
		CreatedByVersion: sql.NullString{String: book.CreatedByVersion, Valid: book.CreatedByVersion != ""},
		IsRemote:         sql.NullBool{Bool: book.IsRemote, Valid: true},
		RemoteUrl:        sql.NullString{String: book.RemoteURL, Valid: book.RemoteURL != ""},
	}
}

// ToSQLCUpdateBookParams 将model.Book转换为sqlc.UpdateBookParams //"Valid"必须是验证条件或true
func ToSQLCUpdateBookParams(book *model.Book) UpdateBookParams {
	return UpdateBookParams{
		Title:            book.Title,
		Owner:            sql.NullString{String: "admin", Valid: true},
		BookPath:         book.BookPath,
		StoreUrl:         book.StoreUrl,
		Type:             string(book.Type),
		ChildBooksNum:    sql.NullInt64{Int64: int64(book.ChildBooksNum), Valid: true},
		ChildBooksID:     sql.NullString{String: strings.Join(book.ChildBooksID, ", "), Valid: len(book.ChildBooksID) > 0},
		Depth:            sql.NullInt64{Int64: int64(book.Depth), Valid: true},
		ParentFolder:     sql.NullString{String: book.ParentFolder, Valid: book.ParentFolder != ""},
		PageCount:        sql.NullInt64{Int64: int64(book.PageCount), Valid: true},
		FileSize:         sql.NullInt64{Int64: book.FileSize, Valid: true},
		Author:           sql.NullString{String: book.Author, Valid: book.Author != ""},
		Isbn:             sql.NullString{String: book.ISBN, Valid: book.ISBN != ""},
		Press:            sql.NullString{String: book.Press, Valid: book.Press != ""},
		PublishedAt:      sql.NullString{String: book.PublishedAt, Valid: book.PublishedAt != ""},
		ExtractPath:      sql.NullString{String: book.ExtractPath, Valid: book.ExtractPath != ""},
		ExtractNum:       sql.NullInt64{Int64: int64(book.ExtractNum), Valid: true},
		BookComplete:     sql.NullBool{Bool: book.BookComplete, Valid: true},
		InitComplete:     sql.NullBool{Bool: book.InitComplete, Valid: true},
		NonUtf8zip:       sql.NullBool{Bool: book.NonUTF8Zip, Valid: true},
		ZipTextEncoding:  sql.NullString{String: book.ZipTextEncoding, Valid: book.ZipTextEncoding != ""},
		CreatedByVersion: sql.NullString{String: book.CreatedByVersion, Valid: book.CreatedByVersion != ""},
		IsRemote:         sql.NullBool{Bool: book.IsRemote, Valid: true},
		RemoteUrl:        sql.NullString{String: book.RemoteURL, Valid: book.RemoteURL != ""},
		BookID:           book.BookID,
	}
}

// ==================== PageInfo 相关转换 ====================

// FromSQLCPageInfo 将sqlc.PageInfo转换为model.PageInfo
func FromSQLCPageInfo(sqlcPageInfo PageInfo) model.PageInfo {
	return model.PageInfo{
		Name:       sqlcPageInfo.Name,
		Path:       sqlcPageInfo.Path.String,
		Size:       sqlcPageInfo.Size.Int64,
		ModTime:    sqlcPageInfo.ModTime.Time,
		Url:        sqlcPageInfo.Url.String,
		PageNum:    int(sqlcPageInfo.PageNum.Int64),
		Blurhash:   sqlcPageInfo.Blurhash.String,
		Height:     int(sqlcPageInfo.Height.Int64),
		Width:      int(sqlcPageInfo.Width.Int64),
		ImgType:    sqlcPageInfo.ImgType.String,
		InsertHtml: sqlcPageInfo.InsertHtml.String,
	}
}

// ToSQLCCreatePageInfoParams 将model.PageInfo 转换为sqlc.CreatePageInfoParams //"Valid"必须是验证条件或true
func ToSQLCCreatePageInfoParams(pageInfo model.PageInfo, bookID string) CreatePageInfoParams {
	return CreatePageInfoParams{
		BookID:     bookID,
		Name:       pageInfo.Name,
		Path:       sql.NullString{String: pageInfo.Path, Valid: pageInfo.Path != ""},
		Size:       sql.NullInt64{Int64: pageInfo.Size, Valid: true},
		ModTime:    sql.NullTime{Time: pageInfo.ModTime, Valid: !pageInfo.ModTime.IsZero()},
		Url:        sql.NullString{String: pageInfo.Url, Valid: pageInfo.Url != ""},
		PageNum:    sql.NullInt64{Int64: int64(pageInfo.PageNum), Valid: true},
		Blurhash:   sql.NullString{String: pageInfo.Blurhash, Valid: pageInfo.Blurhash != ""},
		Height:     sql.NullInt64{Int64: int64(pageInfo.Height), Valid: true},
		Width:      sql.NullInt64{Int64: int64(pageInfo.Width), Valid: true},
		ImgType:    sql.NullString{String: pageInfo.ImgType, Valid: pageInfo.ImgType != ""},
		InsertHtml: sql.NullString{String: pageInfo.InsertHtml, Valid: pageInfo.InsertHtml != ""},
	}
}

// ToSQLCUpdatePageInfoParams 将model.PageInfo转换为sqlc.UpdatePageInfoParams //"Valid"必须是验证条件或true
func ToSQLCUpdatePageInfoParams(pageInfo model.PageInfo, bookID string) UpdatePageInfoParams {
	return UpdatePageInfoParams{
		Name:       pageInfo.Name,
		Path:       sql.NullString{String: pageInfo.Path, Valid: pageInfo.Path != ""},
		Size:       sql.NullInt64{Int64: pageInfo.Size, Valid: true},
		ModTime:    sql.NullTime{Time: pageInfo.ModTime, Valid: !pageInfo.ModTime.IsZero()},
		Url:        sql.NullString{String: pageInfo.Url, Valid: pageInfo.Url != ""},
		Blurhash:   sql.NullString{String: pageInfo.Blurhash, Valid: pageInfo.Blurhash != ""},
		Height:     sql.NullInt64{Int64: int64(pageInfo.Height), Valid: true},
		Width:      sql.NullInt64{Int64: int64(pageInfo.Width), Valid: true},
		ImgType:    sql.NullString{String: pageInfo.ImgType, Valid: pageInfo.ImgType != ""},
		InsertHtml: sql.NullString{String: pageInfo.InsertHtml, Valid: pageInfo.InsertHtml != ""},
		BookID:     bookID,
	}
}

// ==================== 批量转换函数 ====================

// FromSQLCBooks 批量转换sqlc.Book为model.Book
func FromSQLCBooks(sqlcBooks []Book, pagesMap map[string][]model.PageInfo, bookmarksMap map[string]model.BookMarks) []*model.Book {
	books := make([]*model.Book, len(sqlcBooks))
	for i, sqlcBook := range sqlcBooks {
		books[i] = FromSQLCBook(sqlcBook)
		if pagesMap != nil {
			if pages, exists := pagesMap[sqlcBook.BookID]; exists {
				books[i].PageInfos = pages
			}
		}
		if bookmarksMap != nil {
			if marks, exists := bookmarksMap[sqlcBook.BookID]; exists {
				books[i].BookMarks = marks
			}
		}
	}
	return books
}

// FromSQLCPageInfos 批量转换sqlc.PageInfo为model.PageInfo
func FromSQLCPageInfos(sqlcPageInfos []PageInfo) []model.PageInfo {
	pageInfos := make([]model.PageInfo, len(sqlcPageInfos))
	for i, sqlcPageInfo := range sqlcPageInfos {
		pageInfos[i] = FromSQLCPageInfo(sqlcPageInfo)
	}
	return pageInfos
}
