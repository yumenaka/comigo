package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jxskiss/base62"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

// StoreInRam 内存书库，扫描后生成。可以有多个子书库。
type StoreInRam struct {
	StoreInfo
	ChildStores sync.Map // key为路径 存储 *Store
}

// AddStore 创建一个新书库
func (ramStore *StoreInRam) AddStore(storeURL string) error {
	if _, ok := ramStore.ChildStores.Load(storeURL); ok {
		// 已经有这个key了
		return errors.New("add Bookstore Error： The key already exists [" + storeURL + "]")
	}
	s := Store{
		StoreInfo: StoreInfo{
			BackendURL:  storeURL,
			Name:        filepath.Base(storeURL), // 使用路径的最后一部分作为名称
			Description: "Comigo StoreInfo for " + storeURL,
		},
	}
	ramStore.ChildStores.Store(storeURL, &s)
	return nil
}

// GenerateBookGroup 分析所有子书库，并并生成书籍组
func (ramStore *StoreInRam) GenerateBookGroup() (e error) {
	// 遍历所有子书库
	for _, value := range ramStore.ChildStores.Range {
		s := value.(*Store)
		err := s.GenerateBookGroup()
		if err != nil {
			e = err
		}
	}
	return e
}

func (ramStore *StoreInRam) ListBooks() ([]*model.Book, error) {
	var books []*model.Book
	// 遍历 ChildStores 中的所有书籍
	for _, value := range ramStore.ChildStores.Range {
		childStore := value.(*Store)
		for _, value := range childStore.BookMap.Range {
			book := value.(*model.Book)
			books = append(books, book)
		}
	}
	return books, nil
}

// SaveBooks  保存书籍信息到本地硬盘
func (ramStore *StoreInRam) SaveBooks() error {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}
	savePath := filepath.Join(configDir, "books")

	allBooks, err := ramStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
		return err
	}

	// 遍历并保存每本书
	for _, book := range allBooks {
		// 序列化书籍为 JSON 格式
		jsonData, err := json.MarshalIndent(book, "", "  ")
		if err != nil {
			logger.Infof("Error marshaling book %s: %s", book.BookID, err)
			continue // 跳过这本书，继续处理其他书籍
		}
		storePathAbs, err := filepath.Abs(book.StoreUrl)
		if err != nil {
			logger.Infof("Failed to get absolute path: %s", err)
			storePathAbs = book.StoreUrl
		}
		cacheDir := filepath.Join(savePath, base62.EncodeToString([]byte(storePathAbs)))
		// 创建保存目录
		err = os.MkdirAll(cacheDir, os.ModePerm)
		if err != nil {
			logger.Infof("Error creating book_data directory: %s", err)
			return err
		}
		// 构造文件路径
		fileName := filepath.Join(cacheDir, book.BookID+".json")

		// 写入文件
		err = os.WriteFile(fileName, jsonData, 0o644)
		if err != nil {
			logger.Infof("Error saving book %s to file: %s", book.BookID, err)
			continue // 跳过这本书，继续处理其他书籍
		}
	}

	logger.Infof("Successfully saved %d books to %s", len(allBooks), savePath)
	return nil
}

// LoadBooks 从本地路径加载书籍信息
func (ramStore *StoreInRam) LoadBooks() error {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}
	savePath := filepath.Join(configDir, "books")
	// 遍历所有 storeUrl 对应的目录
	for _, storeUrl := range config.GetCfg().StoreUrls {
		storePathAbs, err := filepath.Abs(storeUrl)
		if err != nil {
			logger.Infof("Failed to get absolute path: %s", err)
			storePathAbs = storeUrl
		}
		cacheDir := filepath.Join(savePath, base62.EncodeToString([]byte(storePathAbs)))
		// 检查目录是否存在
		_, err = os.Stat(cacheDir)
		// 目录不存在是正常的（首次运行），不返回错误
		if os.IsNotExist(err) {
			logger.Infof("Book data directory does not exist yet: %s", cacheDir)
			continue
		}
		// 其他错误
		if err != nil {
			logger.Infof("Error accessing book_data directory: %s", err)
			continue
		}
		// 读取目录中的所有文件
		entries, err := os.ReadDir(cacheDir)
		if err != nil {
			logger.Infof("Error reading book_data directory: %s", err)
			return err
		}
		// 统计加载成功的书籍数量
		loadedCount := 0
		// 遍历所有文件
		for _, entry := range entries {
			// 跳过目录，只处理文件
			if entry.IsDir() {
				continue
			}
			// 只处理 .json 文件
			fileName := entry.Name()
			if !strings.HasSuffix(fileName, ".json") {
				continue
			}
			// 读取文件内容
			filePath := filepath.Join(cacheDir, fileName)
			jsonData, err := os.ReadFile(filePath)
			if err != nil {
				logger.Infof("Error reading file %s: %s", fileName, err)
				continue // 跳过这个文件，继续处理其他文件
			}
			// 反序列化为 Book 对象
			var book model.Book
			err = json.Unmarshal(jsonData, &book)
			if err != nil {
				logger.Infof("Warning: corrupted JSON file %s, skipping: %s", fileName, err)
				continue // 跳过损坏的文件，继续处理其他文件
			}
			// 添加书籍到内存
			err = ramStore.AddBook(&book)
			if err != nil {
				logger.Infof("Error adding book %s to store: %s", book.BookID, err)
				continue // 跳过这本书，继续处理其他书籍
			}
			loadedCount++
		}
		logger.Infof("Successfully loaded %d books from %s", loadedCount, cacheDir)
	}

	return nil
}

// AddBook 添加一本书
func (ramStore *StoreInRam) AddBook(b *model.Book) error {
	if b.BookID == "" {
		return errors.New("add book error: empty BookID")
	}
	if _, ok := ramStore.ChildStores.Load(b.StoreUrl); !ok {
		if err := ramStore.AddStore(b.StoreUrl); err != nil {
			logger.Infof("Error adding subfolder: %s", err)
		}
	}
	return ramStore.AddBookToSubStore(b.StoreUrl, b)
}

// AddBookToSubStore 将某一本书，放到BookMap里面去
func (ramStore *StoreInRam) AddBookToSubStore(storeURL string, b *model.Book) error {
	if f, ok := ramStore.ChildStores.Load(storeURL); !ok {
		// 创建一个新子书库，并添加一本书
		newSubStore := Store{
			StoreInfo: StoreInfo{
				BackendURL:  storeURL,
				Name:        filepath.Base(storeURL), // 使用路径的最后一部分作为名称
				Description: "Comigo StoreInfo for " + storeURL,
			},
		}
		newSubStore.BookMap.Store(b.BookID, b)
		ramStore.ChildStores.Store(storeURL, &newSubStore)
		return errors.New("add Bookstore Error： The key not found [" + storeURL + "]")
	} else {
		// 给已有子书库添加一本书
		temp := f.(*Store)
		temp.BookMap.Store(b.BookID, b)
		return nil
	}
}

// AddBooks 添加多本书
func (ramStore *StoreInRam) AddBooks(books []*model.Book) error {
	for _, b := range books {
		if err := ramStore.AddBook(b); err != nil {
			logger.Infof("Error adding book %s: %s", b.BookID, err)
		}
	}
	return nil
}

// GetParentBookID 通过子书籍 ID 获取所属书组 ID
func (ramStore *StoreInRam) GetParentBookID(childID string) (string, error) {
	allBooks, err := ramStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, bookGroup := range allBooks {
		if bookGroup.Type != model.TypeBooksGroup {
			continue // 只处理书组类型
		}
		for _, id := range bookGroup.ChildBooksID {
			if id == childID {
				fmt.Println("Found group for book ID:", childID, "Group ID:", bookGroup.BookID)
				return bookGroup.BookID, nil
			}
		}
	}
	return "", errors.New("cannot find group, id=" + childID)
}

// GetArchiveBooks 获取所有压缩包格式的书籍
func (ramStore *StoreInRam) GetArchiveBooks() []*model.Book {
	var list []*model.Book
	allBooks, err := ramStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, b := range allBooks {
		if b.Type == model.TypeZip || b.Type == model.TypeRar || b.Type == model.TypeCbz || b.Type == model.TypeCbr || b.Type == model.TypeTar || b.Type == model.TypeEpub {
			list = append(list, b)
		}
	}
	return list
}

// GetBook 根据 BookID 获取书籍
// GetBookByID 根据 BookID 获取书籍
func (ramStore *StoreInRam) GetBook(id string) (*model.Book, error) {
	// 遍历 ChildStores ，删除指定 ID 的书籍
	for _, value := range ramStore.ChildStores.Range {
		childStore := value.(*Store)
		if value, ok := childStore.BookMap.Load(id); ok {
			b := value.(*model.Book)
			return b, nil
		}
	}
	return nil, errors.New("GetBook：cannot find book, id=" + id)
}

func (ramStore *StoreInRam) DeleteBook(id string) error {
	for _, value := range ramStore.ChildStores.Range {
		childStore := value.(*Store)
		if _, ok := childStore.BookMap.Load(id); ok {
			childStore.BookMap.Delete(id)
			return nil
		}
	}
	return errors.New("DeleteBook：cannot find book, id=" + id)
}

func (ramStore *StoreInRam) UpdateBook(b *model.Book) error {
	for _, value := range ramStore.ChildStores.Range {
		childStore := value.(*Store)
		if _, ok := childStore.BookMap.Load(b.BookID); ok {
			childStore.BookMap.Store(b.BookID, b)
			return nil
		}
	}
	return errors.New("UpdateBook：cannot find book, id=" + b.BookID)
}

// GetBookByAuthor 获取同一作者的书籍
func (ramStore *StoreInRam) GetBookByAuthor(author string, sortBy string) ([]*model.Book, error) {
	var bookList []*model.Book
	allBooks, err := ramStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, b := range allBooks {
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

// TopOfShelfInfo 获取顶层书架信息
func TopOfShelfInfo(sortBy string) (*model.BookInfos, error) {
	// 显示顶层书库的书籍
	var infoList model.BookInfos
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, b := range allBooks {
		if b.Depth == 0 {
			infoList = append(infoList, b.BookInfo)
		}
	}
	if len(infoList) > 0 {
		infoList.SortBooks(sortBy)
		return &infoList, nil
	}
	// 没找到任何书
	return nil, errors.New("error: cannot find book in TopOfShelfInfo")
}

// GetChildBooksInfo 根据 ID 获取书籍列表
func GetChildBooksInfo(BookID string) (*model.BookInfos, error) {
	var infoList model.BookInfos
	parentBook, err := model.IStore.GetBook(BookID)
	if err != nil {
		return nil, errors.New("cannot find child books info，BookID：" + BookID)
	}
	for _, childID := range parentBook.ChildBooksID {
		b, err := model.IStore.GetBook(childID)
		if err != nil {
			return nil, errors.New("GetParentBook: cannot find book by childID=" + childID)
		}
		infoList = append(infoList, b.BookInfo)
	}
	if len(infoList) > 0 {
		return &infoList, nil
	} else {
		return nil, errors.New("cannot find child books info，BookID：" + BookID)
	}
}

// GetBookInfoListByParentFolder 根据父文件夹获取书籍列表
func GetBookInfoListByParentFolder(parentFolder string) (*model.BookInfos, error) {
	var infoList model.BookInfos
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, b := range allBooks {
		if b.ParentFolder == parentFolder {
			infoList = append(infoList, b.BookInfo)
		}
	}
	if len(infoList) > 0 {
		infoList.SortBooks("filename")
		return &infoList, nil
	}
	return nil, errors.New("cannot find book, parentFolder=" + parentFolder)
}
