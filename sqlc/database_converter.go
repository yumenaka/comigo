package sqlc

import (
	"database/sql"
	"strings"
	"time"

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
			Type:            model.SupportFileType(sqlcBook.Type),
			ZipTextEncoding: sqlcBook.ZipTextEncoding.String,
		},
	}
}

// ToSQLCCreateBookParams 将model.Book转换为sqlc.CreateBookParams
func ToSQLCCreateBookParams(book *model.Book) CreateBookParams {
	return CreateBookParams{
		Title:           book.Title,
		BookID:          book.BookID,
		Owner:           sql.NullString{String: "admin", Valid: true},
		FilePath:        book.FilePath,
		BookStorePath:   book.BookStorePath,
		Type:            string(book.Type),
		ChildBooksNum:   sql.NullInt64{Int64: int64(book.ChildBooksNum), Valid: true},
		ChildBooksID:    sql.NullString{String: strings.Join(book.ChildBooksID, ", "), Valid: len(book.ChildBooksID) > 0},
		Depth:           sql.NullInt64{Int64: int64(book.Depth), Valid: true},
		ParentFolder:    sql.NullString{String: book.ParentFolder, Valid: true},
		PageCount:       sql.NullInt64{Int64: int64(book.PageCount), Valid: true},
		FileSize:        sql.NullInt64{Int64: book.FileSize, Valid: true},
		Author:          sql.NullString{String: book.Author},
		Isbn:            sql.NullString{String: book.ISBN},
		Press:           sql.NullString{String: book.Press},
		PublishedAt:     sql.NullString{String: book.PublishedAt},
		ExtractPath:     sql.NullString{String: book.ExtractPath},
		ExtractNum:      sql.NullInt64{Int64: int64(book.ExtractNum), Valid: true},
		InitComplete:    sql.NullBool{Bool: book.InitComplete, Valid: true},
		ReadPercent:     sql.NullFloat64{Float64: book.ReadPercent, Valid: true},
		NonUtf8zip:      sql.NullBool{Bool: book.NonUTF8Zip},
		ZipTextEncoding: sql.NullString{String: book.ZipTextEncoding},
	}
}

// ToSQLCUpdateBookParams 将model.Book转换为sqlc.UpdateBookParams
func ToSQLCUpdateBookParams(book *model.Book) UpdateBookParams {
	return UpdateBookParams{
		Title:           book.Title,
		Owner:           sql.NullString{String: "admin", Valid: true},
		FilePath:        book.FilePath,
		BookStorePath:   book.BookStorePath,
		Type:            string(book.Type),
		ChildBooksNum:   sql.NullInt64{Int64: int64(book.ChildBooksNum), Valid: true},
		Depth:           sql.NullInt64{Int64: int64(book.Depth), Valid: true},
		ParentFolder:    sql.NullString{String: book.ParentFolder},
		PageCount:       sql.NullInt64{Int64: int64(book.PageCount), Valid: true},
		FileSize:        sql.NullInt64{Int64: book.FileSize, Valid: true},
		Author:          sql.NullString{String: book.Author},
		Isbn:            sql.NullString{String: book.ISBN},
		Press:           sql.NullString{String: book.Press},
		PublishedAt:     sql.NullString{String: book.PublishedAt},
		ExtractPath:     sql.NullString{String: book.ExtractPath},
		ExtractNum:      sql.NullInt64{Int64: int64(book.ExtractNum), Valid: true},
		InitComplete:    sql.NullBool{Bool: book.InitComplete},
		ReadPercent:     sql.NullFloat64{Float64: book.ReadPercent},
		NonUtf8zip:      sql.NullBool{Bool: book.NonUTF8Zip},
		ZipTextEncoding: sql.NullString{String: book.ZipTextEncoding},
		BookID:          book.BookID,
	}
}

// ==================== MediaFileInfo 相关转换 ====================

// FromSQLCMediaFile 将sqlc.MediaFile转换为model.MediaFileInfo
func FromSQLCMediaFile(sqlcMediaFile MediaFile) model.MediaFileInfo {
	return model.MediaFileInfo{
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
func ToSQLCCreateMediaFileParams(mediaFile model.MediaFileInfo, bookID string) CreateMediaFileParams {
	return CreateMediaFileParams{
		BookID:     bookID,
		Name:       mediaFile.Name,
		Path:       sql.NullString{String: mediaFile.Path},
		Size:       sql.NullInt64{Int64: mediaFile.Size, Valid: true},
		ModTime:    sql.NullTime{Time: mediaFile.ModTime},
		Url:        sql.NullString{String: mediaFile.Url},
		PageNum:    sql.NullInt64{Int64: int64(mediaFile.PageNum), Valid: true},
		Blurhash:   sql.NullString{String: mediaFile.Blurhash},
		Height:     sql.NullInt64{Int64: int64(mediaFile.Height), Valid: true},
		Width:      sql.NullInt64{Int64: int64(mediaFile.Width), Valid: true},
		ImgType:    sql.NullString{String: mediaFile.ImgType},
		InsertHtml: sql.NullString{String: mediaFile.InsertHtml},
	}
}

// ToSQLCUpdateMediaFileParams 将model.MediaFileInfo转换为sqlc.UpdateMediaFileParams
func ToSQLCUpdateMediaFileParams(mediaFile model.MediaFileInfo, bookID string) UpdateMediaFileParams {
	return UpdateMediaFileParams{
		Name:       mediaFile.Name,
		Path:       sql.NullString{String: mediaFile.Path},
		Size:       sql.NullInt64{Int64: mediaFile.Size, Valid: true},
		ModTime:    sql.NullTime{Time: mediaFile.ModTime},
		Url:        sql.NullString{String: mediaFile.Url},
		Blurhash:   sql.NullString{String: mediaFile.Blurhash},
		Height:     sql.NullInt64{Int64: int64(mediaFile.Height), Valid: true},
		Width:      sql.NullInt64{Int64: int64(mediaFile.Width), Valid: true},
		ImgType:    sql.NullString{String: mediaFile.ImgType},
		InsertHtml: sql.NullString{String: mediaFile.InsertHtml},
		BookID:     bookID,
	}
}

// ==================== Backend 相关转换 ====================

// FromSQLCFileBackend 将sqlc.FileBackend转换为model.Backend
func FromSQLCFileBackend(sqlcFileBackend FileBackend) *model.Backend {
	return &model.Backend{
		Type:         model.BackendType(sqlcFileBackend.Type),
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
func ToSQLCCreateFileBackendParams(fileBackend *model.Backend) CreateFileBackendParams {
	return CreateFileBackendParams{
		Type:         int64(fileBackend.Type),
		Url:          fileBackend.URL,
		ServerHost:   sql.NullString{String: fileBackend.ServerHost},
		ServerPort:   sql.NullInt64{Int64: int64(fileBackend.ServerPort)},
		NeedAuth:     sql.NullBool{Bool: fileBackend.NeedAuth, Valid: true},
		AuthUsername: sql.NullString{String: fileBackend.AuthUsername},
		AuthPassword: sql.NullString{String: fileBackend.AuthPassword},
		SmbShareName: sql.NullString{String: fileBackend.SMBShareName},
		SmbPath:      sql.NullString{String: fileBackend.SMBPath},
	}
}

// ToSQLCUpdateFileBackendParams 将model.Backend转换为sqlc.UpdateFileBackendParams
func ToSQLCUpdateFileBackendParams(fileBackend *model.Backend) UpdateFileBackendParams {
	return UpdateFileBackendParams{
		Url:          fileBackend.URL,
		Type:         int64(fileBackend.Type),
		ServerHost:   sql.NullString{String: fileBackend.ServerHost},
		ServerPort:   sql.NullInt64{Int64: int64(fileBackend.ServerPort)},
		NeedAuth:     sql.NullBool{Bool: fileBackend.NeedAuth, Valid: true},
		AuthUsername: sql.NullString{String: fileBackend.AuthUsername},
		AuthPassword: sql.NullString{String: fileBackend.AuthPassword},
		SmbShareName: sql.NullString{String: fileBackend.SMBShareName},
		SmbPath:      sql.NullString{String: fileBackend.SMBPath},
		Url_2:        fileBackend.URL, // WHERE条件中的URL参数
	}
}

// ==================== StoreInfo 相关转换 ====================

// FromSQLCStore 将sqlc.Store转换为model.StoreInfo
func FromSQLCStore(sqlcStore Store) *model.StoreInfo {
	backend := model.Backend{
		URL: sqlcStore.BackendUrl,
	}
	err := backend.ParseStoreURL(sqlcStore.BackendUrl)
	if err != nil {
		return nil
	}
	return &model.StoreInfo{
		BackendURL:  sqlcStore.BackendUrl,
		Name:        sqlcStore.Name,
		Description: sqlcStore.Description.String,
		Backend:     backend,
	}
}

// ToSQLCCreateStoreParams 将model.StoreInfo转换为sqlc.CreateStoreParams
func ToSQLCCreateStoreParams(store *model.StoreInfo) CreateStoreParams {
	return CreateStoreParams{
		BackendUrl:  store.BackendURL,
		Name:        store.Name,
		Description: sql.NullString{String: store.Description, Valid: store.Description != ""},
	}
}

// ToSQLCUpdateStoreParams 将model.StoreInfo转换为sqlc.UpdateStoreParams
func ToSQLCUpdateStoreParams(store *model.StoreInfo) UpdateStoreParams {
	return UpdateStoreParams{
		Name:        store.Name,
		Description: sql.NullString{String: store.Description, Valid: store.Description != ""},
		BackendUrl:  store.BackendURL,
	}
}

// ==================== 批量转换函数 ====================

// FromSQLCBooks 批量转换sqlc.Book为model.Book
func FromSQLCBooks(sqlcBooks []Book, pagesMap map[string][]model.MediaFileInfo) []*model.Book {
	books := make([]*model.Book, len(sqlcBooks))
	for i, sqlcBook := range sqlcBooks {
		books[i] = FromSQLCBook(sqlcBook)
		if pages, exists := pagesMap[sqlcBook.BookID]; exists {
			books[i].Images = pages
		} else {
			books[i].Images = []model.MediaFileInfo{} // 确保即使没有页面也不会出错
		}
	}
	return books
}

// FromSQLCMediaFiles 批量转换sqlc.MediaFile为model.MediaFileInfo
func FromSQLCMediaFiles(sqlcMediaFiles []MediaFile) []model.MediaFileInfo {
	mediaFiles := make([]model.MediaFileInfo, len(sqlcMediaFiles))
	for i, sqlcMediaFile := range sqlcMediaFiles {
		mediaFiles[i] = FromSQLCMediaFile(sqlcMediaFile)
	}
	return mediaFiles
}

// FromSQLCFileBackends 批量转换sqlc.FileBackend为model.Backend
func FromSQLCFileBackends(sqlcFileBackends []FileBackend) []*model.Backend {
	fileBackends := make([]*model.Backend, len(sqlcFileBackends))
	for i, sqlcFileBackend := range sqlcFileBackends {
		fileBackends[i] = FromSQLCFileBackend(sqlcFileBackend)
	}
	return fileBackends
}

// FromSQLCStores 批量转换sqlc.Store为model.StoreInfo
func FromSQLCStores(sqlcStores []Store) []*model.StoreInfo {
	stores := make([]*model.StoreInfo, len(sqlcStores))
	for i, sqlcStore := range sqlcStores {
		stores[i] = FromSQLCStore(sqlcStore)
	}
	return stores
}

// ==================== 关联查询结果转换 ====================

// FromSQLCStoreWithBackendRow 将sqlc.GetStoreWithBackendRow转换为model.StoreInfo
func FromSQLCStoreWithBackendRow(row GetStoreWithBackendRow) *model.StoreInfo {
	store := &model.StoreInfo{
		BackendURL:  row.BackendUrl,
		Name:        row.Name,
		Description: row.Description.String,
		Backend: model.Backend{
			Type:         model.BackendType(row.Type),
			URL:          row.Url,
			ServerHost:   row.ServerHost.String,
			ServerPort:   int(row.ServerPort.Int64),
			NeedAuth:     row.NeedAuth.Bool,
			AuthUsername: row.AuthUsername.String,
			AuthPassword: row.AuthPassword.String,
			SMBShareName: row.SmbShareName.String,
			SMBPath:      row.SmbPath.String,
		},
	}
	return store
}

// FromSQLCListStoresWithBackendRow 批量转换sqlc.ListStoresWithBackendRow为model.StoreInfo
func FromSQLCListStoresWithBackendRow(rows []ListStoresWithBackendRow) []*model.StoreInfo {
	stores := make([]*model.StoreInfo, len(rows))
	for i, row := range rows {
		// 将ListStoresWithBackendRow转换为GetStoreWithBackendRow格式
		convertedRow := GetStoreWithBackendRow{
			BackendUrl:   row.BackendUrl,
			Name:         row.Name,
			Description:  row.Description,
			CreatedAt:    row.CreatedAt,
			UpdatedAt:    row.UpdatedAt,
			Type:         row.Type,
			Url:          row.Url,
			ServerHost:   row.ServerHost,
			ServerPort:   row.ServerPort,
			NeedAuth:     row.NeedAuth,
			AuthUsername: row.AuthUsername,
			AuthPassword: row.AuthPassword,
			SmbShareName: row.SmbShareName,
			SmbPath:      row.SmbPath,
		}
		stores[i] = FromSQLCStoreWithBackendRow(convertedRow)
	}
	return stores
}

// ==================== User 相关转换 ====================

// parseExpireAtString 解析过期时间字符串为time.Time
func parseExpireAtString(expireAt string) time.Time {
	if expireAt == "" {
		return time.Time{}
	}
	// 尝试解析多种时间格式
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02",
	}
	for _, format := range formats {
		if t, err := time.Parse(format, expireAt); err == nil {
			return t
		}
	}
	return time.Time{}
}

// formatTimeToString 将time.Time格式化为字符串
func formatTimeToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// FromSQLCUser 将sqlc.User转换为model.User
func FromSQLCUser(sqlcUser User) *model.User {
	return &model.User{
		ID:       int(sqlcUser.ID),
		Username: sqlcUser.Username,
		Password: sqlcUser.Password,
		Role:     sqlcUser.Role.String,
		Email:    sqlcUser.Email.String,
		Key:      sqlcUser.Key.String,
		ExpireAt: formatTimeToString(sqlcUser.ExpiresAt.Time),
	}
}

// ToSQLCCreateUserParams 将model.User转换为sqlc.CreateUserParams
func ToSQLCCreateUserParams(user *model.User) CreateUserParams {
	expiresAt := parseExpireAtString(user.ExpireAt)
	return CreateUserParams{
		Username:  user.Username,
		Password:  user.Password,
		Role:      sql.NullString{String: user.Role, Valid: user.Role != ""},
		Email:     sql.NullString{String: user.Email, Valid: user.Email != ""},
		Key:       sql.NullString{String: user.Key, Valid: user.Key != ""},
		ExpiresAt: sql.NullTime{Time: expiresAt, Valid: !expiresAt.IsZero()},
	}
}

// ToSQLCUpdateUserParams 将model.User转换为sqlc.UpdateUserParams
func ToSQLCUpdateUserParams(user *model.User) UpdateUserParams {
	expiresAt := parseExpireAtString(user.ExpireAt)
	return UpdateUserParams{
		Username:  user.Username,
		Password:  user.Password,
		Role:      sql.NullString{String: user.Role, Valid: user.Role != ""},
		Email:     sql.NullString{String: user.Email, Valid: user.Email != ""},
		Key:       sql.NullString{String: user.Key, Valid: user.Key != ""},
		ExpiresAt: sql.NullTime{Time: expiresAt, Valid: !expiresAt.IsZero()},
		ID:        int64(user.ID),
	}
}

// ToSQLCUpdateUserPasswordParams 将用户ID和新密码转换为sqlc.UpdateUserPasswordParams
func ToSQLCUpdateUserPasswordParams(userID int, newPassword string) UpdateUserPasswordParams {
	return UpdateUserPasswordParams{
		Password: newPassword,
		ID:       int64(userID),
	}
}

// ToSQLCUpdateUserKeyParams 将用户ID、key和过期时间转换为sqlc.UpdateUserKeyParams
func ToSQLCUpdateUserKeyParams(userID int, key string, expiresAt time.Time) UpdateUserKeyParams {
	return UpdateUserKeyParams{
		Key:       sql.NullString{String: key, Valid: key != ""},
		ExpiresAt: sql.NullTime{Time: expiresAt, Valid: !expiresAt.IsZero()},
		ID:        int64(userID),
	}
}

// FromSQLCUsers 批量转换sqlc.User为model.User
func FromSQLCUsers(sqlcUsers []User) []*model.User {
	users := make([]*model.User, len(sqlcUsers))
	for i, sqlcUser := range sqlcUsers {
		users[i] = FromSQLCUser(sqlcUser)
	}
	return users
}
