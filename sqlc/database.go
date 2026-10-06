package sqlc

import (
	"context"
	"errors"
	"fmt"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/tools/logger"
)

// StoreBook 向数据库中插入一本书
func (db *StoreDatabase) StoreBook(book *model.Book) error {
	return db.write(func(tx *StoreDatabase) error { return tx.storeBook(book) })
}

func (db *StoreDatabase) storeBook(book *model.Book) error {
	if book == nil {
		return fmt.Errorf("book is nil")
	}
	if book.BookID == "" {
		return fmt.Errorf("book ID is empty: %s", book.BookPath)
	}
	if err := db.CheckDBQueries(); err != nil {
		return fmt.Errorf("StoreBook: %v", err)
	}
	ctx := context.Background()
	if err := db.queries.UpsertBook(ctx, ToSQLCUpsertBookParams(book)); err != nil {
		return fmt.Errorf("store book: %w", err)
	}
	// 保存书籍的页面信息。即使列表为空也要清理旧记录，保持和 JSON 元数据一致。
	err := db.saveBookPageInfos(ctx, book.BookID, book.PageInfos)
	if err != nil {
		return fmt.Errorf("book media files error: %v", err)
	}
	if err := db.saveBookBookmarks(ctx, book.BookID, book.BookMarks); err != nil {
		return fmt.Errorf("book bookmarks error: %v", err)
	}
	return nil
}

func (db *StoreDatabase) StoreBookMark(mark *model.BookMark) error {
	return db.write(func(tx *StoreDatabase) error { return tx.storeBookMark(mark) })
}

func (db *StoreDatabase) storeBookMark(mark *model.BookMark) error {
	if mark == nil {
		return errors.New("bookmark is nil")
	}
	// 获取书籍
	b, err := db.getBook(mark.BookID)
	if err != nil {
		return fmt.Errorf(locale.GetString("err_storebookmark_cannot_find"), mark.BookID)
	}
	switch mark.Type {
	case model.UserMark:
		// 用户书签的处理逻辑（用户书签可以有多个，但同一页只能有一个用户书签）
		exists := false
		for i, existingMark := range b.BookMarks {
			if existingMark.Type == mark.Type && existingMark.PageIndex == mark.PageIndex {
				// 更新现有书签
				b.BookMarks[i] = *mark
				exists = true
			}
		}
		if !exists {
			b.BookMarks = append(b.BookMarks, *mark)
		}
	case model.AutoMark:
		// 自动书签的处理逻辑（每本书只有一个自动书签）
		exists := false
		for i, existingMark := range b.BookMarks {
			if existingMark.Type == mark.Type {
				// 更新现有书签
				b.BookMarks[i] = *mark
				exists = true
			}
		}
		if !exists {
			b.BookMarks = append(b.BookMarks, *mark)
		}
	default:
		// 目前没有其他类型书签
		return errors.New(locale.GetString("err_storebookmark_unknown_type"))
	}
	err = db.saveBookBookmarks(context.Background(), b.BookID, b.BookMarks)
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

// DeleteBookMark 删除指定书籍的特定书签
// 根据 bookID + markType + pageIndex 唯一确定一个书签
func (db *StoreDatabase) DeleteBookMark(bookID string, markType model.MarkType, pageIndex int) error {
	return db.write(func(tx *StoreDatabase) error { return tx.deleteBookMark(bookID, markType, pageIndex) })
}

func (db *StoreDatabase) deleteBookMark(bookID string, markType model.MarkType, pageIndex int) error {
	// 获取书籍
	b, err := db.getBook(bookID)
	if err != nil {
		return fmt.Errorf(locale.GetString("err_getbookmark_cannot_find"), bookID)
	}
	// 查找并删除匹配的书签
	found := false
	newBookMarks := make(model.BookMarks, 0, len(b.BookMarks))
	for _, mark := range b.BookMarks {
		if mark.Type == markType && mark.PageIndex == pageIndex {
			found = true
			continue // 跳过要删除的书签
		}
		newBookMarks = append(newBookMarks, mark)
	}
	if !found {
		return fmt.Errorf("bookmark not found: bookID=%s, type=%s, pageIndex=%d", bookID, markType, pageIndex)
	}
	b.BookMarks = newBookMarks
	// 持久化更新
	return db.saveBookBookmarks(context.Background(), b.BookID, b.BookMarks)
}

// saveBookPageInfos  保存书籍的媒体文件信息
func (db *StoreDatabase) saveBookPageInfos(ctx context.Context, bookID string, pageInfos []model.PageInfo) error {
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
		logger.Infof(locale.GetString("log_saved_media_files_for_book"), len(pageInfos), bookID)
	}
	return nil
}

// saveBookBookmarks 保存书籍的书签信息
func (db *StoreDatabase) saveBookBookmarks(ctx context.Context, bookID string, bookmarks model.BookMarks) error {
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
		logger.Infof(locale.GetString("log_saved_bookmarks_for_book"), len(bookmarks), bookID)
	}
	return nil
}

// ListBooks  从数据库查询所有书籍的详细信息,避免重复扫描压缩包。忽略已删除书籍
func (db *StoreDatabase) ListBooks() (list []*model.Book, err error) {
	if err := db.CheckDBQueries(); err != nil {
		return nil, err
	}
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.listBooks()
}

func (db *StoreDatabase) listBooks() (list []*model.Book, err error) {
	if err := db.CheckDBQueries(); err != nil {
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
		sqlcPageInfos, err := db.queries.GetPageInfosByBookID(ctx, sqlcBook.BookID)
		if err != nil {
			return nil, err
		} else {
			pagesMap[sqlcBook.BookID] = FromSQLCPageInfos(sqlcPageInfos)
		}
		sqlcBookmarks, err := db.queries.ListBookmarksByBookID(ctx, sqlcBook.BookID)
		if err != nil {
			return nil, err
		} else {
			bookmarksMap[sqlcBook.BookID] = FromSQLCBookmarks(sqlcBookmarks)
		}
	}

	// 批量转换
	books := FromSQLCBooks(sqlcBooks, pagesMap, bookmarksMap)
	return books, nil
}

// GetBook 根据ID获取书籍信息
func (db *StoreDatabase) GetBook(bookID string) (*model.Book, error) {
	if err := db.CheckDBQueries(); err != nil {
		return nil, err
	}
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.getBook(bookID)
}

func (db *StoreDatabase) getBook(bookID string) (*model.Book, error) {
	ctx := context.Background()
	// 查询书籍基本信息
	sqlcBook, err := db.queries.GetBookByID(ctx, bookID)
	if err != nil {
		return nil, err
	}
	// 补充页面信息
	book := FromSQLCBook(sqlcBook)
	imagesSQL, err := db.queries.GetPageInfosByBookID(ctx, sqlcBook.BookID)
	if err != nil {
		return nil, err
	}
	book.PageInfos = FromSQLCPageInfos(imagesSQL)
	bookmarksSQL, err := db.queries.ListBookmarksByBookID(ctx, sqlcBook.BookID)
	if err != nil {
		return nil, err
	}
	book.BookMarks = FromSQLCBookmarks(bookmarksSQL)
	return book, nil
}

// GenerateBookGroup 复用内存书库的分组规则；临时 Store 仅计算拓扑，不写 JSON。
func (db *StoreDatabase) GenerateBookGroup() error {
	if err := db.CheckDBQueries(); err != nil {
		return err
	}
	db.groupsMu.Lock()
	defer db.groupsMu.Unlock()
	books, err := db.ListBooks()
	if err != nil {
		return err
	}
	stores := map[string]*store.Store{}
	for _, book := range books {
		s := stores[book.StoreUrl]
		if s == nil {
			s = &store.Store{StoreInfo: store.StoreInfo{BackendURL: book.StoreUrl}}
			stores[book.StoreUrl] = s
		}
		s.BookMap.Store(book.BookID, book)
	}
	var groups []*model.Book
	for _, s := range stores {
		if err := s.GenerateBookGroup(); err != nil {
			return err
		}
		for _, value := range s.BookMap.Range {
			book := value.(*model.Book)
			if book.Type == model.TypeBooksGroup && book.RemoteBookID == "" {
				groups = append(groups, book)
			}
		}
	}
	// 计算或写入失败时保留旧拓扑，不向读者暴露半份书组。
	return db.write(func(tx *StoreDatabase) error {
		for _, book := range books {
			if book.Type == model.TypeBooksGroup && book.RemoteBookID == "" {
				if err := tx.deleteBook(book.BookID); err != nil {
					return err
				}
			}
		}
		for _, group := range groups {
			if err := tx.storeBook(group); err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteBook 删除书籍信息
func (db *StoreDatabase) DeleteBook(bookID string) error {
	return db.write(func(tx *StoreDatabase) error { return tx.deleteBook(bookID) })
}

func (db *StoreDatabase) deleteBook(bookID string) error {
	if err := db.CheckDBQueries(); err != nil {
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
