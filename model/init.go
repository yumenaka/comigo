package model

type StoreInterface interface {
	ClearTempFilesALL(debug bool, cachePath string)
	AddBook(storeURL string, b *Book, minPageNum int) error
	GetBooksNumber() int
	ListBooks() []*Book
	CheckRawFileExist(filePath string, bookType SupportFileType) bool
	GetBookByID(id string, sortBy string) (*Book, error)
	GetShortBookID(fullID string, minLength int) string
	CheckAllNotExistBooks()
	GetParentBook(childID string) (*Book, error)
	TopOfShelfInfo(sortBy string) (*BookInfoList, error)
	GetBookInfoListByParentFolder(parentFolder string, sortBy string) (*BookInfoList, error)
	GetChildBooksInfo(BookID string, sortBy string) (*BookInfoList, error)
	AddBooks(storeURL string, list []*Book, minPageNum int) error
	GenerateAllBookGroup() (e error)
	ClearAllBookData()
	GetAllBookSkipBookGroup() []*Book
}

var IStore StoreInterface
