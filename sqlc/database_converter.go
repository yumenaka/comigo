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
			StoreUrl:        sqlcBook.StoreUrl,
			ChildBooksNum:   int(sqlcBook.ChildBooksNum.Int64),
			ChildBooksID:    ConvertCommaSeparatedString(sqlcBook.ChildBooksID),
			Deleted:         sqlcBook.Deleted.Bool,
			Depth:           int(sqlcBook.Depth.Int64),
			ExtractPath:     sqlcBook.ExtractPath.String,
			ExtractNum:      int(sqlcBook.ExtractNum.Int64),
			BookPath:        sqlcBook.BookPath,
			FileSize:        sqlcBook.FileSize.Int64,
			ISBN:            sqlcBook.Isbn.String,
			InitComplete:    sqlcBook.InitComplete.Bool,
			Modified:        sqlcBook.ModifiedTime.Time,
			NonUTF8Zip:      sqlcBook.NonUtf8zip.Bool,
			PageCount:       int(sqlcBook.PageCount.Int64),
			ParentFolder:    sqlcBook.ParentFolder.String,
			Press:           sqlcBook.Press.String,
			PublishedAt:     sqlcBook.PublishedAt.String,
			Type:            model.SupportFileType(sqlcBook.Type),
			ZipTextEncoding: sqlcBook.ZipTextEncoding.String,
		},
	}
}

// ToSQLCCreateBookParams 将model.Book转换为sqlc.CreateBookParams //"Valid"必须是验证条件或true
func ToSQLCCreateBookParams(book *model.Book) CreateBookParams {
	return CreateBookParams{
		Title:           book.Title,
		BookID:          book.BookID,
		Owner:           sql.NullString{String: "admin", Valid: true},
		BookPath:        book.BookPath,
		StoreUrl:        book.StoreUrl,
		Type:            string(book.Type),
		ChildBooksNum:   sql.NullInt64{Int64: int64(book.ChildBooksNum), Valid: true},
		ChildBooksID:    sql.NullString{String: strings.Join(book.ChildBooksID, ", "), Valid: len(book.ChildBooksID) > 0},
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
		NonUtf8zip:      sql.NullBool{Bool: book.NonUTF8Zip, Valid: true},
		ZipTextEncoding: sql.NullString{String: book.ZipTextEncoding, Valid: book.ZipTextEncoding != ""},
	}
}

// ToSQLCUpdateBookParams 将model.Book转换为sqlc.UpdateBookParams //"Valid"必须是验证条件或true
func ToSQLCUpdateBookParams(book *model.Book) UpdateBookParams {
	return UpdateBookParams{
		Title:           book.Title,
		Owner:           sql.NullString{String: "admin", Valid: true},
		BookPath:        book.BookPath,
		StoreUrl:        book.StoreUrl,
		Type:            string(book.Type),
		ChildBooksNum:   sql.NullInt64{Int64: int64(book.ChildBooksNum), Valid: true},
		ChildBooksID:    sql.NullString{String: strings.Join(book.ChildBooksID, ", "), Valid: len(book.ChildBooksID) > 0},
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
		NonUtf8zip:      sql.NullBool{Bool: book.NonUTF8Zip, Valid: true},
		ZipTextEncoding: sql.NullString{String: book.ZipTextEncoding, Valid: book.ZipTextEncoding != ""},
		BookID:          book.BookID,
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
