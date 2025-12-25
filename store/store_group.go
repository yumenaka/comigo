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
	"github.com/yumenaka/comigo/assets/locale"
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
		return fmt.Errorf(locale.GetString("err_add_bookstore_key_exists"), storeURL)
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

// SaveBooksToJson  保存书籍信息到本地硬盘（JSON 文件）
func (ramStore *StoreInRam) SaveBooksToJson() error {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}
	savePath := filepath.Join(configDir, "books")
	logger.Infof(locale.GetString("log_saving_books_meta_data_to"), savePath)
	allBooks, err := ramStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
		return err
	}
	// 遍历并保存每本书
	for _, book := range allBooks {
		err := SaveBookJson(book)
		if err != nil {
			logger.Infof(locale.GetString("log_error_saving_book"), book.BookID, err)
		}
	}
	logger.Infof(locale.GetString("log_successfully_saved_books"), len(allBooks), savePath)
	return nil
}

// SaveBookJson 将单本书籍信息保存为 JSON 文件
func SaveBookJson(book *model.Book) error {
	configDir, err := config.GetConfigDir()
	savePath := filepath.Join(configDir, "books")
	if err != nil {
		return err
	}
	// 序列化书籍为 JSON 格式
	jsonData, err := json.MarshalIndent(book, "", "  ")
	if err != nil {
		return err
	}
	// StoreID是书库URL路径的 base62 编码
	cacheDir := filepath.Join(savePath, book.GetStoreID())
	// 创建保存目录
	err = os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return err
	}
	// 构造文件路径
	fileName := filepath.Join(cacheDir, book.BookID+".json")
	// 写入文件
	err = os.WriteFile(fileName, jsonData, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func DeleteBookJson(book *model.Book) error {
	configDir, err := config.GetConfigDir()
	savePath := filepath.Join(configDir, "books")
	if err != nil {
		return err
	}
	// StoreID是书库URL路径的 base62 编码
	cacheDir := filepath.Join(savePath, book.GetStoreID())
	// 创建保存目录
	err = os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return err
	}
	// 构造文件路径
	fileName := filepath.Join(cacheDir, book.BookID+".json")
	// 写入文件
	err = os.Remove(fileName)
	if err != nil {
		return err
	}
	return nil
}

// LoadBooks 从本地路径加载书籍信息
func (ramStore *StoreInRam) LoadBooks() error {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return err
	}
	savePath := filepath.Join(configDir, "books")
	logger.Infof(locale.GetString("log_loading_books_from"), savePath)
	logger.Infof(locale.GetString("log_configured_store_urls"), config.GetCfg().StoreUrls)
	// 遍历所有 storeUrl 对应的目录
	for _, storeUrl := range config.GetCfg().StoreUrls {
		// 计算 storeUrl 的绝对路径
		storePathAbs, err := filepath.Abs(storeUrl)
		if err != nil {
			logger.Infof(locale.GetString("log_error_getting_absolute_path"), err)
			storePathAbs = storeUrl
		}
		// 计算缓存目录路径
		cacheDir := filepath.Join(savePath, base62.EncodeToString([]byte(storePathAbs)))
		// 检查目录是否存在
		_, err = os.Stat(cacheDir)
		// 目录不存在是正常的（首次运行），不返回错误
		if os.IsNotExist(err) {
			logger.Infof(locale.GetString("log_book_data_directory_not_exist"), cacheDir)
			continue
		}
		// 其他错误
		if err != nil {
			logger.Infof(locale.GetString("log_error_accessing_book_data_directory"), err)
			continue
		}
		// 读取目录中的所有文件
		entries, err := os.ReadDir(cacheDir)
		if err != nil {
			logger.Infof(locale.GetString("log_error_reading_book_data_directory"), err)
			return err
		}
		// 统计加载成功的书籍数量
		loadedCount := 0
		// 遍历所有文件
		for _, entry := range entries {
			// 跳过目录，只处理文件
			if entry.IsDir() {
				logger.Infof(locale.GetString("log_skipping_directory"), entry.Name())
				continue
			}
			// 只处理 .json 文件
			fileName := entry.Name()
			if !strings.HasSuffix(fileName, ".json") {
				logger.Infof(locale.GetString("log_skipping_non_json_file"), fileName)
				continue
			}
			// 读取文件内容
			filePath := filepath.Join(cacheDir, fileName)
			jsonData, err := os.ReadFile(filePath)
			if err != nil {
				logger.Infof(locale.GetString("log_error_reading_file"), fileName, err)
				continue // 跳过这个文件，继续处理其他文件
			}
			// 反序列化为 Book 对象
			var book model.Book
			err = json.Unmarshal(jsonData, &book)
			if err != nil {
				logger.Infof(locale.GetString("log_warning_corrupted_json_file"), fileName, err)
				// 尝试删除损坏的json文件
				errDel := os.Remove(filePath)
				if errDel != nil {
					logger.Infof(locale.GetString("log_error_deleting_corrupted_file"), fileName, errDel)
				}
				continue // 继续处理其他文件
			}
			// 添加书籍到内存书库
			err = ramStore.StoreBook(&book)
			if err != nil {
				logger.Infof(locale.GetString("log_error_adding_book_to_store"), book.BookID, err)
				continue // 跳过这本书，继续处理其他书籍
			}
			loadedCount++
			if loadedCount%50 == 0 {
				logger.Infof(locale.GetString("log_loaded_books_so_far"), loadedCount, cacheDir)
			}
		}
		logger.Infof(locale.GetString("log_successfully_loaded_books"), loadedCount, cacheDir)
	}

	return nil
}

// AddBook 添加一本书
func (ramStore *StoreInRam) StoreBook(b *model.Book) error {
	if b.BookID == "" {
		return errors.New(locale.GetString("err_add_book_empty_bookid"))
	}
	// Load 返回存储在映射中的键值，如果不存在值，则返回 nil。ok 结果指示是否在映射中找到了值。
	if _, ok := ramStore.ChildStores.Load(b.StoreUrl); !ok {
		if err := ramStore.AddStore(b.StoreUrl); err != nil {
			logger.Infof(locale.GetString("log_error_adding_subfolder"), err)
		}
	}
	err := SaveBookJson(b)
	if err != nil {
		logger.Infof(locale.GetString("log_error_saving_book_to_json"), b.BookID, err)
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
		return fmt.Errorf(locale.GetString("err_add_bookstore_key_not_found"), storeURL)
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
		if err := ramStore.StoreBook(b); err != nil {
			logger.Infof(locale.GetString("log_error_adding_book"), b.BookID, err)
		}
		err := SaveBookJson(b)
		if err != nil {
			logger.Infof(locale.GetString("log_error_saving_book_to_json"), b.BookID, err)
		}
	}
	return nil
}

// GetParentBookID 通过子书籍 ID 获取所属书组 ID
func (ramStore *StoreInRam) GetParentBookID(childID string) (string, error) {
	allBooks, err := ramStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
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
	return "", fmt.Errorf(locale.GetString("err_cannot_find_group"), childID)
}

// GetArchiveBooks 获取所有压缩包格式的书籍
func (ramStore *StoreInRam) GetArchiveBooks() []*model.Book {
	var list []*model.Book
	allBooks, err := ramStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
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
	// 遍历 ChildStores ，获取指定 ID 的书籍
	for _, value := range ramStore.ChildStores.Range {
		childStore := value.(*Store)
		if value, ok := childStore.BookMap.Load(id); ok {
			b := value.(*model.Book)
			return b, nil
		}
	}
	return nil, fmt.Errorf(locale.GetString("err_getbook_cannot_find"), id)
}

func (ramStore *StoreInRam) StoreBookMark(mark *model.BookMark) error {
	// 获取书籍
	b, err := ramStore.GetBook(mark.BookID)
	if err != nil {
		return fmt.Errorf(locale.GetString("err_storebookmark_cannot_find"), mark.BookID)
	}
	switch mark.Type {
	case model.UserMark:
		// 用户书签的处理逻辑(用户书签可以有多个，但同一页只能有一个用户书签）
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
	err = ramStore.StoreBook(b)
	if err != nil {
		return err
	}
	return nil
}

func (ramStore *StoreInRam) GetBookMarks(bookID string) (*model.BookMarks, error) {
	// 获取书籍
	b, err := ramStore.GetBook(bookID)
	if err != nil {
		return nil, fmt.Errorf(locale.GetString("err_getbookmark_cannot_find"), bookID)
	}
	return &b.BookMarks, nil
}

func (ramStore *StoreInRam) DeleteBook(id string) error {
	for _, value := range ramStore.ChildStores.Range {
		childStore := value.(*Store)
		if b, ok := childStore.BookMap.Load(id); ok {
			childStore.BookMap.Delete(id)
			// 删除本地缓存的 JSON 文件
			err := DeleteBookJson(b.(*model.Book))
			if err != nil {
				logger.Infof(locale.GetString("log_error_deleting_book_json_file"), id, err)
			}
			return nil
		}
	}
	return fmt.Errorf(locale.GetString("err_deletebook_cannot_find"), id)
}

// TopOfShelfInfo 获取顶层书架信息
func TopOfShelfInfo(sortBy string) ([]model.StoreBookInfo, error) {
	model.ClearBookWhenStoreUrlNotExist(config.GetCfg().StoreUrls)
	model.ClearBookNotExist()
	// 显示顶层书库的书籍
	var topBookList model.BookInfos
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
	}
	for _, b := range allBooks {
		if b.Depth == 0 {
			topBookList = append(topBookList, b.BookInfo)
		}
	}
	var storeBookInfoList []model.StoreBookInfo
	storeUrls := config.GetCfg().StoreUrls
	for _, storeUrl := range storeUrls {
		storePathAbs, err := filepath.Abs(storeUrl)
		if err != nil {
			logger.Infof(locale.GetString("log_error_getting_absolute_path"), err)
			storePathAbs = storeUrl
		}
		newStoreBookInfo := model.StoreBookInfo{
			StoreUrl: storePathAbs,
		}
		for _, topBook := range topBookList {
			if topBook.StoreUrl == storePathAbs {
				newStoreBookInfo.BookInfos = append(newStoreBookInfo.BookInfos, topBook)
			}
		}
		newStoreBookInfo.BookInfos.SortBooks(sortBy)
		childBookNum := 0
		for _, b := range allBooks {
			if b.StoreUrl == storePathAbs && b.Type != model.TypeBooksGroup {
				childBookNum++
				// logger.Infof("[%v]Counting book %s in store %s, BookID=%s", childBookNum, b.Title, storePathAbs, b.BookID)
			}
		}
		newStoreBookInfo.ChildBookNum = childBookNum
		storeBookInfoList = append(storeBookInfoList, newStoreBookInfo)
	}
	if len(storeBookInfoList) > 0 {
		return storeBookInfoList, nil
	}
	// 没找到任何书
	return nil, errors.New(locale.GetString("err_cannot_find_book_topofshelf"))
}

// GetChildBooksInfo 根据 ID 获取书籍列表
func GetChildBooksInfo(BookID string) (*model.BookInfos, error) {
	model.ClearBookWhenStoreUrlNotExist(config.GetCfg().StoreUrls)
	model.ClearBookNotExist()
	var infoList model.BookInfos
	parentBook, err := model.IStore.GetBook(BookID)
	if err != nil {
		return nil, fmt.Errorf(locale.GetString("err_cannot_find_child_books"), BookID)
	}
	for _, childID := range parentBook.ChildBooksID {
		b, err := model.IStore.GetBook(childID)
		if err != nil {
			return nil, fmt.Errorf(locale.GetString("err_getparentbook_cannot_find"), childID)
		}
		infoList = append(infoList, b.BookInfo)
	}
	if len(infoList) > 0 {
		return &infoList, nil
	} else {
		return nil, fmt.Errorf(locale.GetString("err_cannot_find_child_books"), BookID)
	}
}

// GetBookInfoListByParentFolder 根据父文件夹获取书籍列表
func GetBookInfoListByParentFolder(parentFolder string) (*model.BookInfos, error) {
	var infoList model.BookInfos
	allBooks, err := model.IStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
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
	return nil, fmt.Errorf(locale.GetString("err_cannot_find_book_parentfolder"), parentFolder)
}
