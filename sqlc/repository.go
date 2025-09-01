package sqlc

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"path"
	"path/filepath"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

// 参考：
// https://docs.sqlc.dev/en/stable/tutorials/getting-started-sqlite.html#setting-up

//go:embed schema.sql
var ddl string
var client *sql.DB
var Repo *Repository

// Repository 书籍数据访问层
type Repository struct {
	queries *Queries
}

// NewBookRepository 创建新的BookRepository实例
func NewBookRepository(db DBTX) *Repository {
	return &Repository{
		queries: New(db),
	}
}

func OpenDatabase(configFilePath string) error {
	// 使用内存数据库
	// dataSourceName := ":memory:"
	// 如果没有配置文件的话，默认在当前目录下创建一个数据库文件
	dataSourceName := "file:comigo.sqlite?cache=shared"
	// 如果有配置文件的话，把数据库文件在同一文件夹内: dataSourceName = "file:comigo.sqlite?cache=shared"
	if configFilePath != "" {
		// 如果有配置文件的话，数据库文件在同一文件夹内: dataSourceName = "file:comigo.sqlite?cache=shared"
		configDir := filepath.Dir(configFilePath) // 不能用path.Dir()，因为windows返回 "."
		dataSourceName = "file:" + path.Join(configDir, "comigo.sqlite") + "?cache=shared"
		logger.Infof(locale.GetString("init_database")+"%s", dataSourceName)
	}
	ctx := context.Background()
	var err error
	client, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		logger.Infof("Failed to open database: %v", err)
		return err
	}

	// Test database connection
	if err = client.PingContext(ctx); err != nil {
		logger.Infof("Failed to ping database: %v", err)
		return err
	}

	// create tables - 现在使用 IF NOT EXISTS，所以即使表已存在也不会报错
	if _, err := client.ExecContext(ctx, ddl); err != nil {
		logger.Infof("Failed to create tables: %v", err)
		// 即使创建表失败，我们也要尝试创建 DBQueries，因为表可能已经存在
		// 只要数据库连接正常，就应该能正常工作
	}

	// 创建 Repository 实例
	Repo = NewBookRepository(client)
	logger.Infof("Database initialized successfully")
	return nil
}

func CloseDatabase() {
	err := client.Close()
	if err != nil {
		logger.Infof("%s", err)
	}
}

// CheckDBQueries 检查 queries 是否已初始化
func (repo *Repository) CheckDBQueries() error {
	if repo.queries == nil {
		return fmt.Errorf("database not initialized, DBQueries is nil")
	}
	return nil
}

// GetBookByID 根据ID获取书籍
func (repo *Repository) GetBookByID(ctx context.Context, bookID string) (*model.Book, error) {
	sqlcBook, err := repo.queries.GetBookByID(ctx, bookID)
	if err != nil {
		return nil, err
	}
	// 补充页面信息
	book := FromSQLCBook(sqlcBook)
	imagesSQL, err := repo.queries.GetMediaFilesByBookID(ctx, sqlcBook.BookID)
	if err == nil {
		book.Images = FromSQLCMediaFiles(imagesSQL)
		book.SortPages("default") // 对页面进行排序
	}
	return book, nil
}

// GetBookByFilePath 根据文件路径获取书籍
func (repo *Repository) GetBookByFilePath(ctx context.Context, filePath string) (*model.Book, error) {
	sqlcBook, err := repo.queries.GetBookByFilePath(ctx, filePath)
	if err != nil {
		return nil, err
	}
	return FromSQLCBook(sqlcBook), nil
}

// GetImagesMap 获取书籍的媒体文件信息的工具函数
func (repo *Repository) GetImagesMap(ctx context.Context, sqlcBooks []Book) map[string][]model.MediaFileInfo {
	imagesMap := make(map[string][]model.MediaFileInfo)
	for _, book := range sqlcBooks {
		imagesSQL, err := repo.queries.GetMediaFilesByBookID(ctx, book.BookID)
		if err != nil {
			continue // 如果获取媒体文件失败，跳过该书籍
		}
		imagesMap[book.BookID] = FromSQLCMediaFiles(imagesSQL)
	}
	return imagesMap
}

// ListBooks 获取所有书籍列表
func (repo *Repository) ListBooks(ctx context.Context) ([]*model.Book, error) {
	sqlcBooks, err := repo.queries.ListBooks(ctx)
	if err != nil {
		return nil, err
	}
	imagesMap := repo.GetImagesMap(ctx, sqlcBooks)
	return FromSQLCBooks(sqlcBooks, imagesMap), nil
}

// ListBooksByType 根据类型获取书籍列表
func (repo *Repository) ListBooksByType(ctx context.Context, bookType string) ([]*model.Book, error) {
	sqlcBooks, err := repo.queries.ListBooksByType(ctx, bookType)
	if err != nil {
		return nil, err
	}
	imagesMap := repo.GetImagesMap(ctx, sqlcBooks)
	return FromSQLCBooks(sqlcBooks, imagesMap), nil
}

// ListBooksByStorePath 根据书库路径获取书籍列表
func (repo *Repository) ListBooksByStorePath(ctx context.Context, storePath string) ([]*model.Book, error) {
	sqlcBooks, err := repo.queries.ListBooksByStorePath(ctx, storePath)
	if err != nil {
		return nil, err
	}
	imagesMap := repo.GetImagesMap(ctx, sqlcBooks)
	return FromSQLCBooks(sqlcBooks, imagesMap), nil
}

// SearchBooksByTitle 根据标题搜索书籍
func (repo *Repository) SearchBooksByTitle(ctx context.Context, title string) ([]*model.Book, error) {
	titleParam := sql.NullString{String: title, Valid: title != ""}
	sqlcBooks, err := repo.queries.SearchBooksByTitle(ctx, titleParam)
	if err != nil {
		return nil, err
	}
	imagesMap := repo.GetImagesMap(ctx, sqlcBooks)
	return FromSQLCBooks(sqlcBooks, imagesMap), nil
}

// Create 创建新书籍
func (repo *Repository) Create(ctx context.Context, book *model.Book) error {
	params := ToSQLCCreateBookParams(book)
	_, err := repo.queries.CreateBook(ctx, params)
	return err
}

// Update 更新书籍信息
func (repo *Repository) Update(ctx context.Context, book *model.Book) error {
	params := ToSQLCUpdateBookParams(book)
	return repo.queries.UpdateBook(ctx, params)
}

// UpdateReadPercent 更新阅读进度
func (repo *Repository) UpdateReadPercent(ctx context.Context, bookID string, readPercent float64) error {
	params := UpdateReadPercentParams{
		ReadPercent: sql.NullFloat64{Float64: readPercent, Valid: true},
		BookID:      bookID,
	}
	return repo.queries.UpdateReadPercent(ctx, params)
}

// MarkAsDeleted 标记书籍为已删除
func (repo *Repository) MarkAsDeleted(ctx context.Context, bookID string) error {
	return repo.queries.MarkBookAsDeleted(ctx, bookID)
}

// Delete 删除书籍
func (repo *Repository) Delete(ctx context.Context, bookID string) error {
	return repo.queries.DeleteBook(ctx, bookID)
}

// GetMediaFiles 获取书籍的媒体文件列表
func (repo *Repository) GetMediaFiles(ctx context.Context, bookID string) ([]model.MediaFileInfo, error) {
	sqlcMediaFiles, err := repo.queries.GetMediaFilesByBookID(ctx, bookID)
	if err != nil {
		return nil, err
	}
	return FromSQLCMediaFiles(sqlcMediaFiles), nil
}

// GetMediaFile 获取特定页面的媒体文件
func (repo *Repository) GetMediaFile(ctx context.Context, bookID string, pageNum int) (*model.MediaFileInfo, error) {
	pageNumParam := sql.NullInt64{Int64: int64(pageNum), Valid: true}
	params := GetMediaFileByBookIDAndPageParams{
		BookID:  bookID,
		PageNum: pageNumParam,
	}
	sqlcMediaFile, err := repo.queries.GetMediaFileByBookIDAndPage(ctx, params)
	if err != nil {
		return nil, err
	}
	mediaFile := FromSQLCMediaFile(sqlcMediaFile)
	return &mediaFile, nil
}

// GetCover 获取书籍封面
func (repo *Repository) GetCover(ctx context.Context, bookID string) (*model.MediaFileInfo, error) {
	sqlcMediaFile, err := repo.queries.GetBookCover(ctx, bookID)
	if err != nil {
		return nil, err
	}
	mediaFile := FromSQLCMediaFile(sqlcMediaFile)
	return &mediaFile, nil
}

// CreateMediaFile 创建媒体文件记录
func (repo *Repository) CreateMediaFile(ctx context.Context, mediaFile model.MediaFileInfo, bookID string) error {
	params := ToSQLCCreateMediaFileParams(mediaFile, bookID)
	_, err := repo.queries.CreateMediaFile(ctx, params)
	return err
}

// UpdateMediaFile 更新媒体文件信息
func (repo *Repository) UpdateMediaFile(ctx context.Context, mediaFile model.MediaFileInfo, bookID string) error {
	params := ToSQLCUpdateMediaFileParams(mediaFile, bookID)
	return repo.queries.UpdateMediaFile(ctx, params)
}

// DeleteMediaFiles 删除书籍的所有媒体文件
func (repo *Repository) DeleteMediaFiles(ctx context.Context, bookID string) error {
	return repo.queries.DeleteMediaFilesByBookID(ctx, bookID)
}

// Count 统计书籍总数
func (repo *Repository) Count(ctx context.Context) (int64, error) {
	return repo.queries.CountBooks(ctx)
}

// CountByType 根据类型统计书籍数量
func (repo *Repository) CountByType(ctx context.Context, bookType string) (int64, error) {
	return repo.queries.CountBooksByType(ctx, bookType)
}

// CountMediaFiles 统计书籍的媒体文件数量
func (repo *Repository) CountMediaFiles(ctx context.Context, bookID string) (int64, error) {
	return repo.queries.CountMediaFilesByBookID(ctx, bookID)
}

// GetTotalFileSize 获取总文件大小
func (repo *Repository) GetTotalFileSize(ctx context.Context) (float64, error) {
	result, err := repo.queries.GetTotalFileSize(ctx)
	if err != nil {
		return 0, err
	}
	if result.Valid {
		return result.Float64, nil
	}
	return 0, nil
}
