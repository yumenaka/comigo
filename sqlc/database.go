package sqlc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

// StoreBook 向数据库中插入一本书
func (db *StoreDatabase) StoreBook(book *model.Book) error {
	if book == nil {
		return fmt.Errorf("book is nil")
	}
	if book.BookID == "" {
		return fmt.Errorf("book ID is empty: %s", book.BookPath)
	}
	if err := DbStore.CheckDBQueries(); err != nil {
		return fmt.Errorf("StoreBook: %v", err)
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
		logger.Infof("Updated existing book: %s %s", book.BookID, book.BookPath)
	}
	// 保存书籍的页面信息（媒体文件）
	if len(book.PageInfos) > 0 {
		err = db.SaveBookPageInfos(ctx, book.BookID, book.PageInfos)
		if err != nil {
			return fmt.Errorf("book media files error: %v", err)
		}
	}
	if err := db.SaveBookBookmarks(ctx, book.BookID, book.BookMarks); err != nil {
		return fmt.Errorf("book bookmarks error: %v", err)
	}
	return nil
}

func (db *StoreDatabase) StoreBookMark(mark *model.BookMark) error {
	// 获取书籍
	b, err := db.GetBook(mark.BookID)
	if err != nil {
		return fmt.Errorf(locale.GetString("err_storebookmark_cannot_find"), mark.BookID)
	}
	switch mark.Type {
	case model.UserMark:
		// 用户书签的处理逻辑（用户书签可以有多个，但同一页只能有一个用户书签）
		exits := false
		for i, existingMark := range b.BookMarks {
			if existingMark.Type == mark.Type && existingMark.PageIndex == mark.PageIndex {
				// 更新现有书签
				b.BookMarks[i] = *mark
				exits = true
			}
		}
		if !exits {
			b.BookMarks = append(b.BookMarks, *mark)
		}
	case model.AutoMark:
		// 自动书签的处理逻辑（每本书只有一个自动书签）
		exits := false
		for i, existingMark := range b.BookMarks {
			if existingMark.Type == mark.Type {
				// 更新现有书签
				b.BookMarks[i] = *mark
				exits = true
			}
		}
		if !exits {
			b.BookMarks = append(b.BookMarks, *mark)
		}
	default:
		// 目前没有其他类型书签
		return errors.New(locale.GetString("err_storebookmark_unknown_type"))
	}
	err = db.StoreBook(b)
	if err != nil {
		return err
	}
	return nil
}

func (db *StoreDatabase) GetBookMarks(bookID string) (*model.BookMarks, error) {
	// 获取书籍
	b, err := db.GetBook(bookID)
	if err != nil {
		return nil, fmt.Errorf(locale.GetString("err_getbookmark_cannot_find"), bookID)
	}
	return &b.BookMarks, nil
}

// SaveBookPageInfos  保存书籍的媒体文件信息
func (db *StoreDatabase) SaveBookPageInfos(ctx context.Context, bookID string, pageInfos []model.PageInfo) error {
	// 先删除旧的媒体文件记录
	err := db.queries.DeletePageInfosByBookID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("delete old media files error: %v", err)
	}

	// 插入新的媒体文件记录
	for _, pageInfo := range pageInfos {
		// 设置页码
		createParams := ToSQLCCreatePageInfoParams(pageInfo, bookID)
		_, err := db.queries.CreatePageInfo(ctx, createParams)
		if err != nil {
			return fmt.Errorf("create media file %s error: %v", pageInfo.Name, err)
		}
	}
	if config.GetCfg().Debug {
		logger.Infof("Saved %d media files for book %s", len(pageInfos), bookID)
	}
	return nil
}

// SaveBookBookmarks 保存书籍的书签信息
func (db *StoreDatabase) SaveBookBookmarks(ctx context.Context, bookID string, bookmarks model.BookMarks) error {
	if err := db.queries.DeleteBookmarksByBookID(ctx, bookID); err != nil {
		return fmt.Errorf("delete old bookmarks error: %v", err)
	}
	for _, bookmark := range bookmarks {
		params := ToSQLCCreateBookmarkParams(bookID, bookmark)
		if _, err := db.queries.CreateBookmark(ctx, params); err != nil {
			return fmt.Errorf("create bookmark %s error: %v", string(bookmark.Type), err)
		}
	}
	if config.GetCfg().Debug {
		logger.Infof("Saved %d bookmarks for book %s", len(bookmarks), bookID)
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
	pagesMap := make(map[string][]model.PageInfo)
	bookmarksMap := make(map[string]model.BookMarks)
	for _, sqlcBook := range sqlcBooks {
		// 过滤掉已删除的书籍
		if sqlcBook.Deleted.Valid && sqlcBook.Deleted.Bool {
			continue
		}
		sqlcPageInfos, err := db.queries.GetPageInfosByBookID(ctx, sqlcBook.BookID)
		if err != nil {
			logger.Infof("Get media files for book %s error: %s", sqlcBook.BookID, err.Error())
			pagesMap[sqlcBook.BookID] = []model.PageInfo{}
		} else {
			pagesMap[sqlcBook.BookID] = FromSQLCPageInfos(sqlcPageInfos)
		}
		sqlcBookmarks, err := db.queries.ListBookmarksByBookID(ctx, sqlcBook.BookID)
		if err != nil {
			logger.Infof("Get bookmarks for book %s error: %s", sqlcBook.BookID, err.Error())
			bookmarksMap[sqlcBook.BookID] = nil
		} else {
			bookmarksMap[sqlcBook.BookID] = FromSQLCBookmarks(sqlcBookmarks)
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
	books := FromSQLCBooks(validBooks, pagesMap, bookmarksMap)
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
	imagesSQL, err := db.queries.GetPageInfosByBookID(ctx, sqlcBook.BookID)
	if err == nil {
		book.PageInfos = FromSQLCPageInfos(imagesSQL)
		book.SortPages("default") // 对页面进行排序
	}
	bookmarksSQL, err := db.queries.ListBookmarksByBookID(ctx, sqlcBook.BookID)
	if err == nil {
		book.BookMarks = FromSQLCBookmarks(bookmarksSQL)
	} else if config.GetCfg().Debug {
		logger.Infof("Get bookmarks for book %s error: %s", sqlcBook.BookID, err.Error())
	}
	return book, nil
}

// GenerateBookGroup 分析所有子书库，并并生成书籍组
func (db *StoreDatabase) GenerateBookGroup() (e error) {
	if err := DbStore.CheckDBQueries(); err != nil {
		return fmt.Errorf("GetAllBook: %v", err)
	}
	ctx := context.Background()
	// 遍历所有子书库
	storeUrls, err := DbStore.queries.ListAllBookStoreURLs(ctx)
	if err != nil {
		return fmt.Errorf("ListAllBookStoreURLs error: %v", err)
	}
	for _, storeUrl := range storeUrls {
		sqlcBook, err := DbStore.queries.ListBooksByStorePath(ctx, storeUrl)
		if err != nil {
			return fmt.Errorf("ListBooksByStorePath error: %v", err)
		}
		storeBooks := FromSQLCBooks(sqlcBook, nil, nil)
		// 遍历 BookMap ，删除所有 BooksGroup 类型的书籍
		for _, b := range storeBooks {
			if b.Type == model.TypeBooksGroup {
				err := db.DeleteBook(b.BookID)
				if err != nil {
					return err
				}
			}
		}
		// 然后再重新生成 BooksGroup
		depthBooksMap := make(map[int][]*model.Book) // key是Depth的临时map
		// 计算最大深度
		maxDepth := 0
		for _, b := range storeBooks {
			depthBooksMap[b.Depth] = append(depthBooksMap[b.Depth], b)
			if b.Depth > maxDepth {
				maxDepth = b.Depth
			}
		}
		// 从深往浅遍历
		// 如果有几本书同时有同一个父文件夹，那么应该【新建】一本书(组)，并加入到depth-1层里面
		for depth := maxDepth; depth >= 0; depth-- {
			// 用父文件夹做key的parentMap，后面遍历用
			parentTempMap := make(map[string][]*model.Book)
			// //遍历depth等于i的所有book
			for _, b := range depthBooksMap[depth] {
				parentTempMap[b.ParentFolder] = append(parentTempMap[b.ParentFolder], b)
			}
			// 循环parentMap，把有相同parent的书创建为一个书组
			for parent, sameParentBookList := range parentTempMap {
				// 新建一本书,类型是书籍组
				// 获取文件夹信息
				pathInfo, err := os.Stat(sameParentBookList[0].BookPath)
				if err != nil {
					return err
				}
				// 获取修改时间
				modTime := pathInfo.ModTime()
				tempBook, err := model.NewBook(filepath.Dir(sameParentBookList[0].BookPath), modTime, 0, storeUrl, depth-1, model.TypeBooksGroup)
				if err != nil {
					if config.GetCfg().Debug {
						logger.Infof("Error creating new book group: %s", err)
					}
					continue
				}
				newBookGroup := tempBook
				// 书名应该设置成parent
				if newBookGroup.Title != parent {
					newBookGroup.Title = parent
				}
				// 初始化ChildBook
				// 然后把同一parent的书，都加进某个书籍组
				for _, bookInList := range sameParentBookList {
					newBookGroup.ChildBooksID = append(newBookGroup.ChildBooksID, bookInList.BookID)
				}
				newBookGroup.ChildBooksNum = len(sameParentBookList)
				// 如果书籍组的子书籍数量等于0，那么不需要添加
				if newBookGroup.ChildBooksNum == 0 {
					continue
				}
				// 检测是否已经生成并添加过
				Added := false
				sqlAllBook, err := db.queries.ListBooks(ctx)
				if err != nil {
					return fmt.Errorf("ListBooks error: %v", err)
				}
				allBooks := FromSQLCBooks(sqlAllBook, nil, nil)
				for _, bookGroup := range allBooks {
					if bookGroup.Type != model.TypeBooksGroup {
						continue // 只关心书籍组类型
					}
					if bookGroup.BookPath == newBookGroup.BookPath {
						Added = true
					}
				}
				// 添加过的不需要添加
				if Added {
					continue
				}
				if (depth - 1) < 0 {
					continue
				}
				depthBooksMap[depth-1] = append(depthBooksMap[depth-1], newBookGroup)
				// 将这本书加到Store的 BookMap 表里面去
				err = db.StoreBook(newBookGroup)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	return e
}

//// UpdateBook 更新书籍信息
//func (db *StoreDatabase) UpdateBook(book *model.Book) error {
//	ctx := context.Background()
//	params := ToSQLCUpdateBookParams(book)
//	return db.queries.UpdateBook(ctx, params)
//}

// DeleteBook 删除书籍信息
func (db *StoreDatabase) DeleteBook(bookID string) error {
	if err := DbStore.CheckDBQueries(); err != nil {
		return err
	}
	ctx := context.Background()
	// 清理书籍信息
	err := db.queries.DeleteBook(ctx, bookID)
	if err != nil {
		return fmt.Errorf("DeleteBook error: %v", err)
	}
	// 清理书籍相关的媒体文件记录
	err = db.queries.DeletePageInfosByBookID(ctx, bookID)
	if err != nil {
		return fmt.Errorf("DeleteBook media files error: %v", err)
	}
	if err := db.queries.DeleteBookmarksByBookID(ctx, bookID); err != nil {
		return fmt.Errorf("DeleteBook bookmarks error: %v", err)
	}
	return nil
}
