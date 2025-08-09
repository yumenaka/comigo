package model

import (
	"database/sql"
	"strings"

	"github.com/yumenaka/comigo/sqlc"
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
func FromSQLCBook(sqlcBook sqlc.Book) *Book {
	return &Book{
		BookInfo: BookInfo{
			Title:           sqlcBook.Title,
			BookID:          sqlcBook.BookID,
			BookStorePath:   sqlcBook.BookStorePath,
			ChildBooksNum:   int(sqlcBook.ChildBooksNum.Int64),
			ChildBooksID:    ConvertCommaSeparatedString(sqlcBook.ChildBooksID),
			Deleted:         sqlcBook.Deleted.Bool,
			Depth:           int(sqlcBook.Depth.Int64),
			ExtractPath:     sqlcBook.ExtractPath.String,
			ExtractNum:      int(sqlcBook.ExtractNum.Int64),
			FilePath:        sqlcBook.FilePath,
			FileSize:        sqlcBook.FileSize.Int64,
			ISBN:            sqlcBook.Isbn.String,
			InitComplete:    sqlcBook.InitComplete.Bool,
			Modified:        sqlcBook.ModifiedTime.Time,
			NonUTF8Zip:      sqlcBook.NonUtf8zip.Bool,
			PageCount:       int(sqlcBook.PageCount.Int64),
			ParentFolder:    sqlcBook.ParentFolder.String,
			Press:           sqlcBook.Press.String,
			PublishedAt:     sqlcBook.PublishedAt.String,
			ReadPercent:     sqlcBook.ReadPercent.Float64,
			Type:            SupportFileType(sqlcBook.Type),
			ZipTextEncoding: sqlcBook.ZipTextEncoding.String,
		},
	}
}

// ToSQLCCreateBookParams 将model.Book转换为sqlc.CreateBookParams
func ToSQLCCreateBookParams(book *Book) sqlc.CreateBookParams {
	return sqlc.CreateBookParams{
		Title:           book.Title,
		BookID:          book.BookID,
		Owner:           sql.NullString{String: "admin", Valid: true},
		FilePath:        book.FilePath,
		BookStorePath:   book.BookStorePath,
		Type:            string(book.Type),
		ChildBooksNum:   sql.NullInt64{Int64: int64(book.ChildBooksNum), Valid: true},
		ChildBooksID:    sql.NullString{String: strings.Join(book.ChildBooksID, ", "), Valid: book.ChildBooksID != nil && len(book.ChildBooksID) > 0},
		Depth:           sql.NullInt64{Int64: int64(book.Depth), Valid: true},
		ParentFolder:    sql.NullString{String: book.ParentFolder, Valid: book.ParentFolder != ""},
		PageCount:       sql.NullInt64{Int64: int64(book.PageCount), Valid: true},
		FileSize:        sql.NullInt64{Int64: book.FileSize, Valid: true},
		Author:          sql.NullString{String: book.Author, Valid: book.Author != ""},
		Isbn:            sql.NullString{String: book.ISBN, Valid: book.ISBN != ""},
		Press:           sql.NullString{String: book.Press, Valid: book.Press != ""},
		PublishedAt:     sql.NullString{String: book.PublishedAt, Valid: book.PublishedAt != ""},
		ExtractPath:     sql.NullString{String: book.ExtractPath, Valid: book.ExtractPath != ""},
		ExtractNum:      sql.NullInt64{Int64: int64(book.ExtractNum), Valid: true},
		InitComplete:    sql.NullBool{Bool: book.InitComplete, Valid: true},
		ReadPercent:     sql.NullFloat64{Float64: book.ReadPercent, Valid: true},
		NonUtf8zip:      sql.NullBool{Bool: book.NonUTF8Zip, Valid: true},
		ZipTextEncoding: sql.NullString{String: book.ZipTextEncoding, Valid: book.ZipTextEncoding != ""},
	}
}

// ToSQLCUpdateBookParams 将model.Book转换为sqlc.UpdateBookParams
func ToSQLCUpdateBookParams(book *Book) sqlc.UpdateBookParams {
	return sqlc.UpdateBookParams{
		Title:           book.Title,
		Owner:           sql.NullString{String: "admin", Valid: true},
		FilePath:        book.FilePath,
		BookStorePath:   book.BookStorePath,
		Type:            string(book.Type),
		ChildBooksNum:   sql.NullInt64{Int64: int64(book.ChildBooksNum), Valid: true},
		Depth:           sql.NullInt64{Int64: int64(book.Depth), Valid: true},
		ParentFolder:    sql.NullString{String: book.ParentFolder, Valid: book.ParentFolder != ""},
		PageCount:       sql.NullInt64{Int64: int64(book.PageCount), Valid: true},
		FileSize:        sql.NullInt64{Int64: book.FileSize, Valid: true},
		Author:          sql.NullString{String: book.Author, Valid: book.Author != ""},
		Isbn:            sql.NullString{String: book.ISBN, Valid: book.ISBN != ""},
		Press:           sql.NullString{String: book.Press, Valid: book.Press != ""},
		PublishedAt:     sql.NullString{String: book.PublishedAt, Valid: book.PublishedAt != ""},
		ExtractPath:     sql.NullString{String: book.ExtractPath, Valid: book.ExtractPath != ""},
		ExtractNum:      sql.NullInt64{Int64: int64(book.ExtractNum), Valid: true},
		InitComplete:    sql.NullBool{Bool: book.InitComplete, Valid: true},
		ReadPercent:     sql.NullFloat64{Float64: book.ReadPercent, Valid: true},
		NonUtf8zip:      sql.NullBool{Bool: book.NonUTF8Zip, Valid: true},
		ZipTextEncoding: sql.NullString{String: book.ZipTextEncoding, Valid: book.ZipTextEncoding != ""},
		BookID:          book.BookID,
	}
}

// ==================== MediaFileInfo 相关转换 ====================

// FromSQLCMediaFile 将sqlc.MediaFile转换为model.MediaFileInfo
func FromSQLCMediaFile(sqlcMediaFile sqlc.MediaFile) MediaFileInfo {
	return MediaFileInfo{
		Name:       sqlcMediaFile.Name,
		Path:       sqlcMediaFile.Path.String,
		Size:       sqlcMediaFile.Size.Int64,
		ModTime:    sqlcMediaFile.ModTime.Time,
		Url:        sqlcMediaFile.Url.String,
		PageNum:    int(sqlcMediaFile.PageNum.Int64),
		Blurhash:   sqlcMediaFile.Blurhash.String,
		Height:     int(sqlcMediaFile.Height.Int64),
		Width:      int(sqlcMediaFile.Width.Int64),
		ImgType:    sqlcMediaFile.ImgType.String,
		InsertHtml: sqlcMediaFile.InsertHtml.String,
	}
}

// ToSQLCCreateMediaFileParams 将model.MediaFileInfo转换为sqlc.CreateMediaFileParams
func ToSQLCCreateMediaFileParams(mediaFile MediaFileInfo, bookID string) sqlc.CreateMediaFileParams {
	return sqlc.CreateMediaFileParams{
		BookID:     bookID,
		Name:       mediaFile.Name,
		Path:       sql.NullString{String: mediaFile.Path, Valid: mediaFile.Path != ""},
		Size:       sql.NullInt64{Int64: mediaFile.Size, Valid: true},
		ModTime:    sql.NullTime{Time: mediaFile.ModTime, Valid: !mediaFile.ModTime.IsZero()},
		Url:        sql.NullString{String: mediaFile.Url, Valid: mediaFile.Url != ""},
		PageNum:    sql.NullInt64{Int64: int64(mediaFile.PageNum), Valid: true},
		Blurhash:   sql.NullString{String: mediaFile.Blurhash, Valid: mediaFile.Blurhash != ""},
		Height:     sql.NullInt64{Int64: int64(mediaFile.Height), Valid: true},
		Width:      sql.NullInt64{Int64: int64(mediaFile.Width), Valid: true},
		ImgType:    sql.NullString{String: mediaFile.ImgType, Valid: mediaFile.ImgType != ""},
		InsertHtml: sql.NullString{String: mediaFile.InsertHtml, Valid: mediaFile.InsertHtml != ""},
	}
}

// ToSQLCUpdateMediaFileParams 将model.MediaFileInfo转换为sqlc.UpdateMediaFileParams
func ToSQLCUpdateMediaFileParams(mediaFile MediaFileInfo, bookID string) sqlc.UpdateMediaFileParams {
	return sqlc.UpdateMediaFileParams{
		Name:       mediaFile.Name,
		Path:       sql.NullString{String: mediaFile.Path, Valid: mediaFile.Path != ""},
		Size:       sql.NullInt64{Int64: mediaFile.Size, Valid: true},
		ModTime:    sql.NullTime{Time: mediaFile.ModTime, Valid: !mediaFile.ModTime.IsZero()},
		Url:        sql.NullString{String: mediaFile.Url, Valid: mediaFile.Url != ""},
		Blurhash:   sql.NullString{String: mediaFile.Blurhash, Valid: mediaFile.Blurhash != ""},
		Height:     sql.NullInt64{Int64: int64(mediaFile.Height), Valid: true},
		Width:      sql.NullInt64{Int64: int64(mediaFile.Width), Valid: true},
		ImgType:    sql.NullString{String: mediaFile.ImgType, Valid: mediaFile.ImgType != ""},
		InsertHtml: sql.NullString{String: mediaFile.InsertHtml, Valid: mediaFile.InsertHtml != ""},
		BookID:     bookID,
		PageNum:    sql.NullInt64{Int64: int64(mediaFile.PageNum), Valid: true},
	}
}

// ==================== FileBackend 相关转换 ====================

// FromSQLCFileBackend 将sqlc.FileBackend转换为model.FileBackend
func FromSQLCFileBackend(sqlcFileBackend sqlc.FileBackend) *FileBackend {
	return &FileBackend{
		Type:         FileBackendType(sqlcFileBackend.Type),
		URL:          sqlcFileBackend.Url,
		ServerHost:   sqlcFileBackend.ServerHost.String,
		ServerPort:   int(sqlcFileBackend.ServerPort.Int64),
		NeedAuth:     sqlcFileBackend.NeedAuth.Bool,
		AuthUsername: sqlcFileBackend.AuthUsername.String,
		AuthPassword: sqlcFileBackend.AuthPassword.String,
		SMBShareName: sqlcFileBackend.SmbShareName.String,
		SMBPath:      sqlcFileBackend.SmbPath.String,
	}
}

// ToSQLCCreateFileBackendParams 将model.FileBackend转换为sqlc.CreateFileBackendParams
func ToSQLCCreateFileBackendParams(fileBackend *FileBackend) sqlc.CreateFileBackendParams {
	return sqlc.CreateFileBackendParams{
		Type:         int64(fileBackend.Type),
		Url:          fileBackend.URL,
		ServerHost:   sql.NullString{String: fileBackend.ServerHost, Valid: fileBackend.ServerHost != ""},
		ServerPort:   sql.NullInt64{Int64: int64(fileBackend.ServerPort), Valid: fileBackend.ServerPort > 0},
		NeedAuth:     sql.NullBool{Bool: fileBackend.NeedAuth, Valid: true},
		AuthUsername: sql.NullString{String: fileBackend.AuthUsername, Valid: fileBackend.AuthUsername != ""},
		AuthPassword: sql.NullString{String: fileBackend.AuthPassword, Valid: fileBackend.AuthPassword != ""},
		SmbShareName: sql.NullString{String: fileBackend.SMBShareName, Valid: fileBackend.SMBShareName != ""},
		SmbPath:      sql.NullString{String: fileBackend.SMBPath, Valid: fileBackend.SMBPath != ""},
	}
}

// ToSQLCUpdateFileBackendParams 将model.FileBackend转换为sqlc.UpdateFileBackendParams
func ToSQLCUpdateFileBackendParams(fileBackend *FileBackend) sqlc.UpdateFileBackendParams {
	return sqlc.UpdateFileBackendParams{
		Type:         int64(fileBackend.Type),
		Url:          fileBackend.URL,
		ServerHost:   sql.NullString{String: fileBackend.ServerHost, Valid: fileBackend.ServerHost != ""},
		ServerPort:   sql.NullInt64{Int64: int64(fileBackend.ServerPort), Valid: fileBackend.ServerPort > 0},
		NeedAuth:     sql.NullBool{Bool: fileBackend.NeedAuth, Valid: true},
		AuthUsername: sql.NullString{String: fileBackend.AuthUsername, Valid: fileBackend.AuthUsername != ""},
		AuthPassword: sql.NullString{String: fileBackend.AuthPassword, Valid: fileBackend.AuthPassword != ""},
		SmbShareName: sql.NullString{String: fileBackend.SMBShareName, Valid: fileBackend.SMBShareName != ""},
		SmbPath:      sql.NullString{String: fileBackend.SMBPath, Valid: fileBackend.SMBPath != ""},
	}
}

// ==================== StoreInfo 相关转换 ====================

// FromSQLCStore 将sqlc.Store转换为model.StoreInfo
func FromSQLCStore(sqlcStore sqlc.Store) *StoreInfo {
	return nil // TODO:这里需要根据实际情况实现转换逻辑
	// return &StoreInfo{
	// 	ID:          int(sqlcStore.ID),
	// 	Name:        sqlcStore.Name,
	// 	Description: sqlcStore.Description.String,
	// 	CreatedAt:   sqlcStore.CreatedAt.Time,
	// 	UpdatedAt:   sqlcStore.UpdatedAt.Time,
	// 	// Backender字段需要通过关联查询获取
	// }
}

// ToSQLCCreateStoreParams 将model.Store转换为sqlc.CreateStoreParams
func ToSQLCCreateStoreParams(store *StoreInfo, url string) sqlc.CreateStoreParams {
	// TODO:这里需要根据实际情况实现转换逻辑
	// return sqlc.CreateStoreParams{
	// }
	return sqlc.CreateStoreParams{
		// Name:          store.Name,
		// Description:   sql.NullString{String: store.Description, Valid: store.Description != ""},
		Url: url,
	}
}

// ToSQLCUpdateStoreParams 将model.Store转换为sqlc.UpdateStoreParams
func ToSQLCUpdateStoreParams(store *StoreInfo, url string) sqlc.UpdateStoreParams {
	return sqlc.UpdateStoreParams{
		// Name:          store.Name,// TODO:这里需要根据实际情况实现转换逻辑
		// Description:   sql.NullString{String: store.Description, Valid: store.Description != ""},
		Url: url,
	}
}

// ==================== 批量转换函数 ====================

// FromSQLCBooks 批量转换sqlc.Book为model.Book
func FromSQLCBooks(sqlcBooks []sqlc.Book, pagesMap map[string][]MediaFileInfo) []*Book {
	books := make([]*Book, len(sqlcBooks))
	for i, sqlcBook := range sqlcBooks {
		books[i] = FromSQLCBook(sqlcBook)
		if pages, exists := pagesMap[sqlcBook.BookID]; exists {
			books[i].Pages.Images = pages
		} else {
			books[i].Pages.Images = []MediaFileInfo{} // 确保即使没有页面也不会出错
		}
	}
	return books
}

// FromSQLCMediaFiles 批量转换sqlc.MediaFile为model.MediaFileInfo
func FromSQLCMediaFiles(sqlcMediaFiles []sqlc.MediaFile) []MediaFileInfo {
	mediaFiles := make([]MediaFileInfo, len(sqlcMediaFiles))
	for i, sqlcMediaFile := range sqlcMediaFiles {
		mediaFiles[i] = FromSQLCMediaFile(sqlcMediaFile)
	}
	return mediaFiles
}

// FromSQLCFileBackends 批量转换sqlc.FileBackend为model.FileBackend
func FromSQLCFileBackends(sqlcFileBackends []sqlc.FileBackend) []*FileBackend {
	fileBackends := make([]*FileBackend, len(sqlcFileBackends))
	for i, sqlcFileBackend := range sqlcFileBackends {
		fileBackends[i] = FromSQLCFileBackend(sqlcFileBackend)
	}
	return fileBackends
}

// FromSQLCStores 批量转换sqlc.Store为model.StoreInfo
func FromSQLCStores(sqlcStores []sqlc.Store) []*StoreInfo {
	stores := make([]*StoreInfo, len(sqlcStores))
	for i, sqlcStore := range sqlcStores {
		stores[i] = FromSQLCStore(sqlcStore)
	}
	return stores
}

// ==================== 关联查询结果转换 ====================

// FromSQLCStoreWithBackendRow 将sqlc.GetStoreWithBackendRow转换为model.StoreInfo
func FromSQLCStoreWithBackendRow(row sqlc.GetStoreWithBackendRow) *StoreInfo {
	store := &StoreInfo{
		// ID:          int(row.ID),// TODO:这里需要根据实际情况实现转换逻辑
		// Name:        row.Name,// TODO:这里需要根据实际情况实现转换逻辑
		// Description: row.Description.String,
		// CreatedAt:   row.CreatedAt.Time,
		// UpdatedAt:   row.UpdatedAt.Time,
		// Backender: &FileBackend{
		// 	Type:         FileBackendType(row.Type),
		// 	URL:          row.Url,
		// 	ServerHost:   row.ServerHost.String,
		// 	ServerPort:   int(row.ServerPort.Int64),
		// 	NeedAuth:     row.NeedAuth.Bool,
		// 	AuthUsername: row.AuthUsername.String,
		// 	AuthPassword: row.AuthPassword.String,
		// 	SMBShareName: row.SmbShareName.String,
		// 	SMBPath:      row.SmbPath.String,
		// },
	}
	return store
}

// FromSQLCListStoresWithBackendRow 批量转换sqlc.ListStoresWithBackendRow为model.StoreInfo
func FromSQLCListStoresWithBackendRow(rows []sqlc.ListStoresWithBackendRow) []*StoreInfo {
	stores := make([]*StoreInfo, len(rows))
	for i, row := range rows {
		stores[i] = FromSQLCStoreWithBackendRow(sqlc.GetStoreWithBackendRow(row))
	}
	return stores
}
