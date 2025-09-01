//go:build !(windows && 386) || !js

package sqlc // Package sqlc 编译条件的注释和 package 语句之间一定要隔一行，不然无法识别编译条件。go:build 是1.18以后"条件编译"的推荐语法。

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

// ClearBookData   清空指定书籍的所有页面信息
func (repo *Repository) ClearBookData(book *model.Book) {
	if book == nil || book.BookID == "" {
		logger.Infof("ClearBookData: book or bookID is empty")
		return
	}

	if err := Repo.CheckDBQueries(); err != nil {
		logger.Infof("ClearBookData: %v", err)
		return
	}

	ctx := context.Background()

	// 清空该书籍的所有媒体文件记录
	err := repo.queries.DeleteMediaFilesByBookID(ctx, book.BookID)
	if err != nil {
		logger.Infof("Clear book media files error: %s", err.Error())
		return
	}

	// 清空书籍的 Images 切片
	book.Images = []model.MediaFileInfo{}

	logger.Infof("Clear book %s media files completed", book.BookID)
}

// DeleteAllBookInDatabase  清空数据库所有 book 与 media_file_info
func (repo *Repository) DeleteAllBookInDatabase(debug bool) {
	if err := Repo.CheckDBQueries(); err != nil {
		logger.Infof("DeleteAllBookInDatabase: %v", err)
		return
	}

	ctx := context.Background()

	// 获取书籍总数
	totalBooks, err := repo.queries.CountBooks(ctx)
	if err != nil {
		logger.Infof("Count books error: %s", err.Error())
		return
	}

	if debug {
		logger.Infof("Starting to delete all books and media files. Total books: %d", totalBooks)
	}

	// 获取所有书籍来删除对应的媒体文件
	books, err := repo.queries.ListBooks(ctx)
	if err != nil {
		logger.Infof("List books error: %s", err.Error())
		return
	}

	// 删除所有媒体文件
	deletedMediaFiles := 0
	for _, book := range books {
		err := repo.queries.DeleteMediaFilesByBookID(ctx, book.BookID)
		if err != nil {
			logger.Infof("Delete media files for book %s error: %s", book.BookID, err.Error())
		} else {
			// 统计该书籍删除的媒体文件数量
			count, _ := repo.queries.CountMediaFilesByBookID(ctx, book.BookID)
			deletedMediaFiles += int(count)
		}
	}

	// 删除所有书籍记录
	// 注意：由于没有 DeleteAllBooks 查询，我们需要逐个删除
	deletedBooks := 0
	for _, book := range books {
		err := repo.queries.DeleteBook(ctx, book.BookID)
		if err != nil {
			logger.Infof("Delete book %s error: %s", book.BookID, err.Error())
		} else {
			deletedBooks++
		}
	}

	if debug {
		logger.Infof("Delete completed. Books deleted: %d, Media files deleted: %d", deletedBooks, deletedMediaFiles)
	} else {
		logger.Infof("Delete all books and media files completed")
	}
}

// SaveAllBookToDatabase 将Map里面的书籍信息，全部保存到本地数据库中
func (repo *Repository) SaveAllBookToDatabase(m map[string]*model.Book) {
	for _, b := range m {
		c := *b
		err := repo.SaveBookToDatabase(&c)
		if err != nil {
			logger.Infof("SaveAllBookToDatabase error :%s", err.Error())
		}
	}
}

// SaveBookListToDatabase  向数据库中插入一组书
func (repo *Repository) SaveBookListToDatabase(bookList []*model.Book) error {
	for _, b := range bookList {
		err := repo.SaveBookToDatabase(b)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveBookToDatabase 向数据库中插入一本书
func (repo *Repository) SaveBookToDatabase(save *model.Book) error {
	if save == nil {
		return fmt.Errorf("book is nil")
	}

	if save.BookID == "" {
		return fmt.Errorf("book ID is empty")
	}

	if err := Repo.CheckDBQueries(); err != nil {
		return fmt.Errorf("SaveBookToDatabase: %v", err)
	}

	ctx := context.Background()

	// 检查书籍是否已存在
	_, err := repo.queries.GetBookByID(ctx, save.BookID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("check existing book error: %v", err)
	}

	// 根据是否存在决定创建或更新
	if errors.Is(err, sql.ErrNoRows) {
		// 书籍不存在，创建新记录
		createParams := ToSQLCCreateBookParams(save)
		_, err = repo.queries.CreateBook(ctx, createParams)
		if err != nil {
			return fmt.Errorf("create book error: %v", err)
		}
		logger.Infof("Created new book: %s", save.BookID)
	} else {
		// 书籍已存在，更新记录
		updateParams := ToSQLCUpdateBookParams(save)
		err = repo.queries.UpdateBook(ctx, updateParams)
		if err != nil {
			return fmt.Errorf("update book error: %v", err)
		}
		logger.Infof("Updated existing book: %s", save.BookID)
	}

	// 保存书籍的页面信息（媒体文件）
	if len(save.Images) > 0 {
		err = repo.SaveBookMediaFiles(ctx, save.BookID, save.Images)
		if err != nil {
			return fmt.Errorf("save media files error: %v", err)
		}
	}

	return nil
}

// saveBookMediaFiles 保存书籍的媒体文件信息
func (repo *Repository) SaveBookMediaFiles(ctx context.Context, bookID string, mediaFiles []model.MediaFileInfo) error {
	// 先删除旧的媒体文件记录
	err := repo.queries.DeleteMediaFilesByBookID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("delete old media files error: %v", err)
	}

	// 插入新的媒体文件记录
	for i, mediaFile := range mediaFiles {
		// 设置页码
		mediaFile.PageNum = i + 1

		createParams := ToSQLCCreateMediaFileParams(mediaFile, bookID)
		_, err := repo.queries.CreateMediaFile(ctx, createParams)
		if err != nil {
			return fmt.Errorf("create media file %s error: %v", mediaFile.Name, err)
		}
	}

	logger.Infof("Saved %d media files for book %s", len(mediaFiles), bookID)
	return nil
}

// GetBookFromDatabase 根据文件路径，从数据库查询一本书的详细信息,避免重复扫描压缩包
func (repo *Repository) GetBookFromDatabase(filepath string) (*model.Book, error) {
	if filepath == "" {
		return nil, fmt.Errorf("filepath is empty")
	}

	if err := Repo.CheckDBQueries(); err != nil {
		return nil, fmt.Errorf("GetBookFromDatabase: %v", err)
	}

	ctx := context.Background()

	// 根据文件路径查询书籍
	sqlcBook, err := repo.queries.GetBookByFilePath(ctx, filepath)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 没有找到记录，返回 nil 而不是错误
		}
		return nil, fmt.Errorf("get book by filepath error: %v", err)
	}

	// 转换为 model.Book
	book := FromSQLCBook(sqlcBook)

	// 查询书籍的媒体文件信息
	sqlcMediaFiles, err := repo.queries.GetMediaFilesByBookID(ctx, book.BookID)
	if err != nil {
		logger.Infof("Get media files for book %s error: %s", book.BookID, err.Error())
		// 即使媒体文件查询失败，我们仍然返回书籍信息
		book.Images = []model.MediaFileInfo{}
	} else {
		// 转换媒体文件信息
		book.Images = FromSQLCMediaFiles(sqlcMediaFiles)
	}

	logger.Infof("Found book %s with %d pages from database", book.BookID, len(book.Images))
	return book, nil
}

// GetBooksFromDatabase  从数据库查询所有书的详细信息,避免重复扫描压缩包。忽略文件夹型的书籍
func (repo *Repository) GetBooksFromDatabase() (list []*model.Book, err error) {
	if err := Repo.CheckDBQueries(); err != nil {
		return nil, fmt.Errorf("GetBooksFromDatabase: %v", err)
	}

	ctx := context.Background()

	// 查询所有书籍（排除已删除的）
	sqlcBooks, err := repo.queries.ListBooks(ctx)
	if err != nil {
		return nil, fmt.Errorf("list books error: %v", err)
	}

	if len(sqlcBooks) == 0 {
		return []*model.Book{}, nil
	}

	// 为每本书查询媒体文件信息
	pagesMap := make(map[string][]model.MediaFileInfo)
	for _, sqlcBook := range sqlcBooks {
		// 过滤掉已删除的书籍
		if sqlcBook.Deleted.Valid && sqlcBook.Deleted.Bool {
			continue
		}

		sqlcMediaFiles, err := repo.queries.GetMediaFilesByBookID(ctx, sqlcBook.BookID)
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

	// 过滤掉文件夹型的书籍（根据注释要求）
	var result []*model.Book
	for _, book := range books {
		// 根据 SupportFileType 定义，文件夹型书籍的类型为 TypeDir
		if book.Type != model.TypeDir {
			result = append(result, book)
		}
	}

	logger.Infof("Found %d books from database (excluding directories)", len(result))
	return result, nil
}
