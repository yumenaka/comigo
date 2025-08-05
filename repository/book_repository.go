package repository

import (
	"context"
	"database/sql"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/sqlc"
)

// BookRepository 书籍数据访问层
type BookRepository struct {
	queries *sqlc.Queries
}

// NewBookRepository 创建新的BookRepository实例
func NewBookRepository(db sqlc.DBTX) *BookRepository {
	return &BookRepository{
		queries: sqlc.New(db),
	}
}

// GetBookByID 根据ID获取书籍
func (r *BookRepository) GetBookByID(ctx context.Context, bookID string) (*model.Book, error) {
	sqlcBook, err := r.queries.GetBookByID(ctx, bookID)
	if err != nil {
		return nil, err
	}
	// 补充页面信息
	book := model.FromSQLCBook(sqlcBook)
	imagesSQL, err := r.queries.GetMediaFilesByBookID(ctx, sqlcBook.BookID)
	if err == nil {
		book.Pages.Images = model.FromSQLCMediaFiles(imagesSQL)
		book.SortPages("filename") // 对页面进行排序
	}
	return book, nil
}

// GetBookByFilePath 根据文件路径获取书籍
func (r *BookRepository) GetBookByFilePath(ctx context.Context, filePath string) (*model.Book, error) {
	sqlcBook, err := r.queries.GetBookByFilePath(ctx, filePath)
	if err != nil {
		return nil, err
	}
	return model.FromSQLCBook(sqlcBook), nil
}

// GetImagesMap 获取书籍的媒体文件信息的工具函数
func (r *BookRepository) GetImagesMap(ctx context.Context, sqlcBooks []sqlc.Book) map[string][]model.MediaFileInfo {
	imagesMap := make(map[string][]model.MediaFileInfo)
	for _, book := range sqlcBooks {
		imagesSQL, err := r.queries.GetMediaFilesByBookID(ctx, book.BookID)
		if err != nil {
			continue // 如果获取媒体文件失败，跳过该书籍
		}
		imagesMap[book.BookID] = model.FromSQLCMediaFiles(imagesSQL)
	}
	return imagesMap
}

// ListBooks 获取所有书籍列表
func (r *BookRepository) ListBooks(ctx context.Context) ([]*model.Book, error) {
	sqlcBooks, err := r.queries.ListBooks(ctx)
	if err != nil {
		return nil, err
	}
	imagesMap := r.GetImagesMap(ctx, sqlcBooks)
	return model.FromSQLCBooks(sqlcBooks, imagesMap), nil
}

// ListBooksByType 根据类型获取书籍列表
func (r *BookRepository) ListBooksByType(ctx context.Context, bookType string) ([]*model.Book, error) {
	sqlcBooks, err := r.queries.ListBooksByType(ctx, bookType)
	if err != nil {
		return nil, err
	}
	imagesMap := r.GetImagesMap(ctx, sqlcBooks)
	return model.FromSQLCBooks(sqlcBooks, imagesMap), nil
}

// ListBooksByStorePath 根据书库路径获取书籍列表
func (r *BookRepository) ListBooksByStorePath(ctx context.Context, storePath string) ([]*model.Book, error) {
	sqlcBooks, err := r.queries.ListBooksByStorePath(ctx, storePath)
	if err != nil {
		return nil, err
	}
	imagesMap := r.GetImagesMap(ctx, sqlcBooks)
	return model.FromSQLCBooks(sqlcBooks, imagesMap), nil
}

// SearchBooksByTitle 根据标题搜索书籍
func (r *BookRepository) SearchBooksByTitle(ctx context.Context, title string) ([]*model.Book, error) {
	titleParam := sql.NullString{String: title, Valid: title != ""}
	sqlcBooks, err := r.queries.SearchBooksByTitle(ctx, titleParam)
	if err != nil {
		return nil, err
	}
	imagesMap := r.GetImagesMap(ctx, sqlcBooks)
	return model.FromSQLCBooks(sqlcBooks, imagesMap), nil
}

// Create 创建新书籍
func (r *BookRepository) Create(ctx context.Context, book *model.Book) error {
	params := model.ToSQLCCreateBookParams(book)
	_, err := r.queries.CreateBook(ctx, params)
	return err
}

// Update 更新书籍信息
func (r *BookRepository) Update(ctx context.Context, book *model.Book) error {
	params := model.ToSQLCUpdateBookParams(book)
	return r.queries.UpdateBook(ctx, params)
}

// UpdateReadPercent 更新阅读进度
func (r *BookRepository) UpdateReadPercent(ctx context.Context, bookID string, readPercent float64) error {
	params := sqlc.UpdateReadPercentParams{
		ReadPercent: sql.NullFloat64{Float64: readPercent, Valid: true},
		BookID:      bookID,
	}
	return r.queries.UpdateReadPercent(ctx, params)
}

// MarkAsDeleted 标记书籍为已删除
func (r *BookRepository) MarkAsDeleted(ctx context.Context, bookID string) error {
	return r.queries.MarkBookAsDeleted(ctx, bookID)
}

// Delete 删除书籍
func (r *BookRepository) Delete(ctx context.Context, bookID string) error {
	return r.queries.DeleteBook(ctx, bookID)
}

// GetMediaFiles 获取书籍的媒体文件列表
func (r *BookRepository) GetMediaFiles(ctx context.Context, bookID string) ([]model.MediaFileInfo, error) {
	sqlcMediaFiles, err := r.queries.GetMediaFilesByBookID(ctx, bookID)
	if err != nil {
		return nil, err
	}
	return model.FromSQLCMediaFiles(sqlcMediaFiles), nil
}

// GetMediaFile 获取特定页面的媒体文件
func (r *BookRepository) GetMediaFile(ctx context.Context, bookID string, pageNum int) (*model.MediaFileInfo, error) {
	pageNumParam := sql.NullInt64{Int64: int64(pageNum), Valid: true}
	params := sqlc.GetMediaFileByBookIDAndPageParams{
		BookID:  bookID,
		PageNum: pageNumParam,
	}
	sqlcMediaFile, err := r.queries.GetMediaFileByBookIDAndPage(ctx, params)
	if err != nil {
		return nil, err
	}
	mediaFile := model.FromSQLCMediaFile(sqlcMediaFile)
	return &mediaFile, nil
}

// GetCover 获取书籍封面
func (r *BookRepository) GetCover(ctx context.Context, bookID string) (*model.MediaFileInfo, error) {
	sqlcMediaFile, err := r.queries.GetBookCover(ctx, bookID)
	if err != nil {
		return nil, err
	}
	mediaFile := model.FromSQLCMediaFile(sqlcMediaFile)
	return &mediaFile, nil
}

// CreateMediaFile 创建媒体文件记录
func (r *BookRepository) CreateMediaFile(ctx context.Context, mediaFile model.MediaFileInfo, bookID string) error {
	params := model.ToSQLCCreateMediaFileParams(mediaFile, bookID)
	_, err := r.queries.CreateMediaFile(ctx, params)
	return err
}

// UpdateMediaFile 更新媒体文件信息
func (r *BookRepository) UpdateMediaFile(ctx context.Context, mediaFile model.MediaFileInfo, bookID string) error {
	params := model.ToSQLCUpdateMediaFileParams(mediaFile, bookID)
	return r.queries.UpdateMediaFile(ctx, params)
}

// DeleteMediaFiles 删除书籍的所有媒体文件
func (r *BookRepository) DeleteMediaFiles(ctx context.Context, bookID string) error {
	return r.queries.DeleteMediaFilesByBookID(ctx, bookID)
}

// Count 统计书籍总数
func (r *BookRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountBooks(ctx)
}

// CountByType 根据类型统计书籍数量
func (r *BookRepository) CountByType(ctx context.Context, bookType string) (int64, error) {
	return r.queries.CountBooksByType(ctx, bookType)
}

// CountMediaFiles 统计书籍的媒体文件数量
func (r *BookRepository) CountMediaFiles(ctx context.Context, bookID string) (int64, error) {
	return r.queries.CountMediaFilesByBookID(ctx, bookID)
}

// GetTotalFileSize 获取总文件大小
func (r *BookRepository) GetTotalFileSize(ctx context.Context) (float64, error) {
	result, err := r.queries.GetTotalFileSize(ctx)
	if err != nil {
		return 0, err
	}
	if result.Valid {
		return result.Float64, nil
	}
	return 0, nil
}
