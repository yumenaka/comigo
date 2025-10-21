package model

type StoreInterface interface {
	// 1. 书库相关基本操作，可能需要优化
	AddBook(storeURL string, b *Book, minPageNum int) error
	AddBooks(storeURL string, list []*Book, minPageNum int) error
	GetBook(id string) (*Book, error)
	GetBookAndSort(id string, sortBy string) (*Book, error)
	GetBooksNumber() int
	ListBooks() []*Book
	ClearAll()
	GetAllBookSkipBookGroup() []*Book
	GetParentBook(childID string) (*Book, error)
	GetChildBooksInfo(BookID string, sortBy string) (*BookInfoList, error)

	// 2. 需要重写细节
	TopOfShelfInfo(sortBy string) (*BookInfoList, error) // 返回按照书库分组的数据

	// 3. 需要研究函数的作用，看看能不能优化掉
	CheckRawFileExist(filePath string, bookType SupportFileType) bool // CheckRawFileExist 查看内存中是否已经有了这本书，有了就返回 true，让调用者跳过
	GetBookInfoListByParentFolder(parentFolder string, sortBy string) (*BookInfoList, error)
	GenerateAllBookGroup() (e error)
	ClearBookNotExist()
	GetShortBookID(fullID string, minLength int) string

	// 4. 非泛用函数，想办法干掉
	ClearTempFilesALL(debug bool, cachePath string)
}

var IStore StoreInterface
