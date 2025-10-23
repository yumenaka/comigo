package model

type StoreInterface interface {
	// 01-12 书库相关基本操作，可能需要优化
	AddBook(b *Book, minPageNum int) error
	AddBooks(list []*Book, minPageNum int) error
	GetBook(id string) (*Book, error)
	UpdateBook(b *Book) error
	ListBooks() []*Book
	ListBookSkipBookGroup() []*Book
	GetAllBooksNumber() int
	GetParentBook(childID string) (*Book, error)
	GetChildBooksInfo(BookID string, sortBy string) (*BookInfoList, error)
	ClearAllBook()
	CheckBookFileExist(filePath string, bookType SupportFileType) bool // CheckRawFileExist 查看书库中是否已经有了这本书，有了就返回 true，让调用者跳过
	GetBookInfoListByParentFolder(parentFolder string, sortBy string) (*BookInfoList, error)

	// 13. TODO：重写
	TopOfShelfInfo(sortBy string) (*BookInfoList, error) // 返回按照书库分组的数据

	// 14-17. 研究函数的作用，移走或优化掉
	GenerateAllBookGroup() (e error)
	ClearBookNotExist()
	GetShortBookID(fullID string, minLength int) string
	ClearTempFilesALL(debug bool, cachePath string)
}

var IStore StoreInterface
