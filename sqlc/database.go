package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

// AddBook 向数据库中插入一本书
func (db *StoreDatabase) AddBook(book *model.Book) error {
	if book == nil {
		return fmt.Errorf("book is nil")
	}
	if book.BookID == "" {
		return fmt.Errorf("book ID is empty" + book.FilePath)
	}
	if err := DbStore.CheckDBQueries(); err != nil {
		return fmt.Errorf("AddBook: %v", err)
	}
	ctx := context.Background()
	// 检查书籍是否已存在
	_, err := db.queries.GetBookByID(ctx, book.BookID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("check existing book error: %v", err)
	}
	// 根据是否存在决定创建或更新
	if errors.Is(err, sql.ErrNoRows) {
		// 书籍不存在，创建新记录
		createParams := ToSQLCCreateBookParams(book)
		_, err = db.queries.CreateBook(ctx, createParams)
		if err != nil {
			return fmt.Errorf("create book error: %v", err)
		}
		logger.Infof("Created new book: %s", book.BookID)
	} else {
		// 书籍已存在，更新记录
		updateParams := ToSQLCUpdateBookParams(book)
		err = db.queries.UpdateBook(ctx, updateParams)
		if err != nil {
			return fmt.Errorf("update book error: %v", err)
		}
		logger.Infof("Updated existing book: %s %s", book.BookID, book.FilePath)
	}
	// 保存书籍的页面信息（媒体文件）
	if len(book.Images) > 0 {
		err = db.SaveBookMediaFiles(ctx, book.BookID, book.Images)
		if err != nil {
			return fmt.Errorf("book media files error: %v", err)
		}
	}
	return nil
}

// SaveBookMediaFiles  保存书籍的媒体文件信息
func (db *StoreDatabase) SaveBookMediaFiles(ctx context.Context, bookID string, mediaFiles []model.MediaFileInfo) error {
	// 先删除旧的媒体文件记录
	err := db.queries.DeleteMediaFilesByBookID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("delete old media files error: %v", err)
	}

	// 插入新的媒体文件记录
	for i, mediaFile := range mediaFiles {
		// 设置页码
		mediaFile.PageNum = i + 1

		createParams := ToSQLCCreateMediaFileParams(mediaFile, bookID)
		_, err := db.queries.CreateMediaFile(ctx, createParams)
		if err != nil {
			return fmt.Errorf("create media file %s error: %v", mediaFile.Name, err)
		}
	}
	if config.GetCfg().Debug {
		logger.Infof("Saved %d media files for book %s", len(mediaFiles), bookID)
	}
	return nil
}

// ListBooks  从数据库查询所有书籍的详细信息,避免重复扫描压缩包。忽略已删除书籍
func (db *StoreDatabase) ListBooks() (list []*model.Book, err error) {
	if err := DbStore.CheckDBQueries(); err != nil {
		return nil, fmt.Errorf("GetAllBook: %v", err)
	}
	ctx := context.Background()

	// 查询所有书籍（排除已删除的书籍）
	sqlcBooks, err := db.queries.ListBooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("list books error: %v", err)
	}
	// 为每本书查询媒体文件信息
	pagesMap := make(map[string][]model.MediaFileInfo)
	for _, sqlcBook := range sqlcBooks {
		// 过滤掉已删除的书籍
		if sqlcBook.Deleted.Valid && sqlcBook.Deleted.Bool {
			continue
		}
		sqlcMediaFiles, err := db.queries.GetMediaFilesByBookID(ctx, sqlcBook.BookID)
		if err != nil {
			logger.Infof("Get media files for book %s error: %s", sqlcBook.BookID, err.Error())
			pagesMap[sqlcBook.BookID] = []model.MediaFileInfo{}
		} else {
			pagesMap[sqlcBook.BookID] = FromSQLCMediaFiles(sqlcMediaFiles)
		}
	}

	// 过滤未删除的书籍
	var validBooks []Book
	for _, book := range sqlcBooks {
		if !book.Deleted.Valid || !book.Deleted.Bool {
			validBooks = append(validBooks, book)
		}
	}
	// 批量转换
	books := FromSQLCBooks(validBooks, pagesMap)
	return books, nil
}

// GetBook 根据ID获取书籍信息
func (db *StoreDatabase) GetBook(bookID string) (*model.Book, error) {
	ctx := context.Background()
	// 查询书籍基本信息
	sqlcBook, err := db.queries.GetBookByID(ctx, bookID)
	if err != nil {
		return nil, err
	}
	// 补充页面信息
	book := FromSQLCBook(sqlcBook)
	imagesSQL, err := db.queries.GetMediaFilesByBookID(ctx, sqlcBook.BookID)
	if err == nil {
		book.Images = FromSQLCMediaFiles(imagesSQL)
		book.SortPages("default") // 对页面进行排序
	}
	return book, nil
}

// UpdateBook 更新书籍信息
func (db *StoreDatabase) UpdateBook(book *model.Book) error {
	ctx := context.Background()
	params := ToSQLCUpdateBookParams(book)
	return db.queries.UpdateBook(ctx, params)
}

// DeleteBook 删除书籍信息
func (db *StoreDatabase) DeleteBook(bookID string) error {
	ctx := context.Background()
	return db.queries.DeleteBook(ctx, bookID)
}

// ClearBookData   清空指定书籍的所有页面信息
func (db *StoreDatabase) ClearBookData(book *model.Book) {
	if book == nil || book.BookID == "" {
		logger.Infof("ClearBookData: book or bookID is empty")
		return
	}
	if err := DbStore.CheckDBQueries(); err != nil {
		logger.Infof("ClearBookData: %v", err)
		return
	}
	ctx := context.Background()
	// 清空该书籍的所有媒体文件记录
	err := db.queries.DeleteMediaFilesByBookID(ctx, book.BookID)
	if err != nil {
		logger.Infof("ClearAllBook book media files error: %s", err.Error())
		return
	}
	// 清空书籍的 Images 切片
	book.Images = []model.MediaFileInfo{}
	logger.Infof("ClearAllBook book %s media files completed", book.BookID)
}
