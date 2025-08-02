package repository

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

// getShortBookID 生成短的 BookID，避免冲突
func getShortBookID(fullID string, minLength int) string {
	if len(fullID) <= minLength {
		logger.Infof("Cannot shorten ID: %s", fullID)
		return fullID
	}
	shortID := fullID[:minLength]
	add := 0
	for {
		conflict := false
		for _, value := range MainStore.mapBooks.Range {
			b := value.(*model.Book)
			if b.BookID == shortID {
				conflict = true
				break
			}
		}
		for _, value := range MainStore.mapBookGroup.Range {
			group := value.(*model.BookInfo)
			if group.BookID == shortID {
				conflict = true
				break
			}
		}
		if !conflict {
			break
		}
		add++
		if minLength+add > len(fullID) {
			break
		}
		shortID = fullID[:minLength+add]
	}
	return shortID
}

// AddBook 添加一本书
func AddBook(b *model.Book, basePath string, minPageNum int) error {
	if b.BookID == "" {
		return errors.New("add book error: empty BookID")
	}
	if b.GetPageCount() < minPageNum {
		return errors.New("add book error: minPageNum = " + strconv.Itoa(b.GetPageCount()))
	}
	if _, ok := MainStore.SubStores.Load(basePath); !ok {
		if err := MainStore.AddSubStore(basePath); err != nil {
			logger.Infof("Error adding subfolder: %s", err)
		}
	}
	MainStore.mapBooks.Store(b.BookID, b)
	return MainStore.AddBookToSubStore(basePath, &b.BookInfo)
}

// AddBooks 添加一组书
func AddBooks(list []*model.Book, basePath string, minPageNum int) error {
	for _, b := range list {
		if b.GetPageCount() < minPageNum {
			continue
		}
		if err := AddBook(b, basePath, minPageNum); err != nil {
			return err
		}
	}
	// 生成虚拟书籍组
	if err := MainStore.AnalyzeStore(); err != nil {
		logger.Infof("%s", err)
	}
	return nil
}

// DeleteBookByID 删除一本书
func DeleteBookByID(bookID string) {
	MainStore.mapBooks.Delete(bookID)
}

// GetBookGroupByBookID 通过数组 ID 获取书组信息
func GetBookGroupByBookID(id string) (*model.BookInfo, error) {
	if value, ok := MainStore.mapBookGroup.Load(id); ok {
		return value.(*model.BookInfo), nil
	}
	return nil, errors.New("GetBookGroupByBookID：cannot find book group, id=" + id)
}

// GetBookGroupIDByBookID 通过子书籍 ID 获取所属书组 ID
func GetBookGroupIDByBookID(id string) (string, error) {
	for _, value := range MainStore.mapBookGroup.Range {
		group := value.(*model.BookInfo)
		for _, childID := range group.ChildBooksID {
			if childID == id {
				return group.BookID, nil
			}
		}
	}
	return "", errors.New("cannot find group, id=" + id)
}

// GetBooksNumber 获取书籍总数，不包括 BookInfo
func GetBooksNumber() int {
	// 用于计数的变量
	var count int
	// 遍历 map 并递增计数器
	for range MainStore.mapBooks.Range {
		count++
	}
	return count
}

// GetAllBookList 获取所有书籍列表
func GetAllBookList() []*model.Book {
	var list []*model.Book
	for _, value := range MainStore.mapBooks.Range {
		b := value.(*model.Book)
		list = append(list, b)
	}
	return list
}

// CheckBookExist 查看内存中是否已经有了这本书，有了就返回 true，让调用者跳过
func CheckBookExist(filePath string, bookType model.SupportFileType) bool {
	if bookType == model.TypeDir || bookType == model.TypeBooksGroup {
		return false
	}
	for _, value := range MainStore.mapBooks.Range {
		realBook := value.(*model.Book)
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			logger.Infof("Error getting absolute path: %v", err)
			continue
		}
		if realBook.FilePath == absFilePath && realBook.Type == bookType {
			return true
		}
	}
	return false
}

// CheckAllBookFileExist 检查内存中的书的源文件是否存在，不存在就删掉
func CheckAllBookFileExist() {
	logger.Infof("Checking book files exist...")
	var deletedBooks []string
	// 遍历所有书籍
	for _, value := range MainStore.mapBooks.Range {
		book := value.(*model.Book)
		if _, err := os.Stat(book.FilePath); os.IsNotExist(err) {
			deletedBooks = append(deletedBooks, book.FilePath)
			DeleteBookByID(book.BookID)
		}
	}
	// 删除不存在的书组
	if len(deletedBooks) > 0 {
		ResetBookGroupData()
	}
}

// GetArchiveBooks 获取所有压缩包格式的书籍
func GetArchiveBooks() []*model.Book {
	var list []*model.Book
	for _, value := range MainStore.mapBooks.Range {
		b := value.(*model.Book)
		if b.Type == model.TypeZip || b.Type == model.TypeRar || b.Type == model.TypeCbz || b.Type == model.TypeCbr || b.Type == model.TypeTar || b.Type == model.TypeEpub {
			list = append(list, b)
		}
	}
	return list
}

// GetBookByID 根据 BookID 获取书籍
func GetBookByID(id string, sortBy string) (*model.Book, error) {
	if value, ok := MainStore.mapBooks.Load(id); ok {
		b := value.(*model.Book)
		b.SortPages(sortBy)
		return b, nil
	}
	if value, ok := MainStore.mapBookGroup.Load(id); ok {
		group := value.(*model.BookInfo)
		return &model.Book{
			BookInfo: *group,
		}, nil
	}
	return nil, errors.New("GetBookByID：cannot find book, id=" + id)
}

// GetRandomBook 随机获取一本书
func GetRandomBook() (*model.Book, error) {
	for _, value := range GetAllBookList() {
		return value, nil // 这里可以改为随机选择
	}
	return nil, errors.New("cannot find any book")
}

// GetBookGroupInfoByChildBookID 通过子书籍 ID 获取所属书组信息
func GetBookGroupInfoByChildBookID(id string) (*model.BookInfo, error) {
	for _, value := range MainStore.mapBookGroup.Range {
		group := value.(*model.BookInfo)
		for _, childID := range group.ChildBooksID {
			if childID == id {
				return group, nil
			}
		}
	}
	return nil, errors.New("cannot find group, id=" + id)
}

// GetBookByAuthor 获取同一作者的书籍
func GetBookByAuthor(author string, sortBy string) ([]*model.Book, error) {
	var bookList []*model.Book
	for _, value := range MainStore.mapBooks.Range {
		b := value.(*model.Book)
		if b.Author == author {
			b.SortPages(sortBy)
			bookList = append(bookList, b)
		}
	}

	if len(bookList) > 0 {
		return bookList, nil
	}
	return nil, errors.New("cannot find book, author=" + author)
}

// ClearAllBookData  清空所有书籍与虚拟书组数据
func ClearAllBookData() {
	ClearBookData()
	ClearBookGroupData()
}

// ClearBookData 清空书籍数据
func ClearBookData() {
	MainStore.mapBooks.Clear()
}

// ClearBookGroupData  清空书组相关数据
func ClearBookGroupData() {
	// Clear 会删除所有条目，从而产生一个空的 Map。
	MainStore.mapBookGroup.Clear()
	MainStore.SubStores.Clear()
}

// ResetBookGroupData 重置虚拟书库
func ResetBookGroupData() {
	ClearBookGroupData()
	if err := MainStore.AnalyzeStore(); err != nil {
		logger.Infof("Error initializing main folder: %s", err)
	}
}

// AddBookToSubStore 将某一本书，放到BookMap里面去
func (folder *memoryStore) AddBookToSubStore(searchPath string, b *model.BookInfo) error {
	if f, ok := folder.SubStores.Load(searchPath); !ok {
		// 创建一个新子书库，并添加一本书
		newStore := subMemoryStore{
			Path: searchPath,
		}
		newStore.BookMap.Store(b.BookID, b)
		folder.SubStores.Store(searchPath, &newStore)
		return errors.New("add Bookstore Error： The key not found [" + searchPath + "]")
	} else {
		// 给已有子书库添加一本书
		temp := f.(*subMemoryStore)
		temp.BookMap.Store(b.BookID, b)
		return nil
	}
}

// AddSubStore 创建一个新文件夹
func (folder *memoryStore) AddSubStore(basePath string) error {
	if _, ok := folder.SubStores.Load(basePath); ok {
		// 已经有这个key了
		return errors.New("add Bookstore Error： The key already exists [" + basePath + "]")
	}
	s := subMemoryStore{
		Path: basePath,
	}
	folder.SubStores.Store(basePath, &s)
	return nil
}

// AnalyzeStore 分析并生成有书籍组（设置深度与层级关系）
func (folder *memoryStore) AnalyzeStore() (e error) {
	// 遍历所有子书库
	for _, value := range folder.SubStores.Range {
		s := value.(*subMemoryStore)
		err := s.AnalyzeFolder()
		if err != nil {
			e = err
		}
	}
	return e
}

func (s *subMemoryStore) AnalyzeFolder() error {
	count := 0
	for range s.BookMap.Range {
		count++
	}
	if count == 0 {
		return errors.New("empty Bookstore")
	}
	depthBooksMap := make(map[int][]model.BookInfo) // key是Depth的临时map
	// 定义一个最大深度
	maxDepth := 0
	for _, value := range s.BookMap.Range {
		b := value.(*model.BookInfo)
		depthBooksMap[b.Depth] = append(depthBooksMap[b.Depth], *b)
		// 找到最大深度
		if b.Depth > maxDepth {
			maxDepth = b.Depth
		}
	}

	// 从深往浅遍历
	// 如果有几本书同时有同一个父文件夹，那么应该【新建]一本书(组)，并加入到depth-1层里面
	for depth := maxDepth; depth >= 0; depth-- {
		// 用父文件夹做key的parentMap，后面遍历用
		parentTempMap := make(map[string][]model.BookInfo)
		// //遍历depth等于i的所有book
		for _, b := range depthBooksMap[depth] {
			parentTempMap[b.ParentFolder] = append(parentTempMap[b.ParentFolder], b)
		}
		// 循环parentMap，把有相同parent的书创建为一个书组
		for parent, sameParentBookList := range parentTempMap {
			// 新建一本书,类型是书籍组
			// 获取文件夹信息
			pathInfo, err := os.Stat(sameParentBookList[0].FilePath)
			if err != nil {
				return err
			}
			// 获取修改时间
			modTime := pathInfo.ModTime()
			newBookGroup, err := model.NewbookinfoBookgroup(filepath.Dir(sameParentBookList[0].FilePath), modTime, 0, s.Path, depth-1, model.TypeBooksGroup)
			if err != nil {
				logger.Infof("%s", err)
				continue
			}
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
			for _, value := range MainStore.mapBookGroup.Range {
				group := value.(*model.BookInfo)
				if group.FilePath == newBookGroup.FilePath {
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
			depthBooksMap[depth-1] = append(depthBooksMap[depth-1], *newBookGroup)
			newBookGroup.SetAuthor()
			// 将这本书加到子书库的BookGroup表（Images.BookGroupMap）里面去
			s.BookGroupMap.Store(newBookGroup.BookID, &newBookGroup)
			// 将这本书加到BookGroup总表（mapBookGroup）里面去
			MainStore.mapBookGroup.Store(newBookGroup.BookID, newBookGroup)
		}
	}
	return nil
}
