package model

// StoreInterface 书库相关基本操作接口
type StoreInterface interface {
	StoreBook(b *Book) error
	GetBook(id string) (*Book, error)
	DeleteBook(id string) error
	ListBooks() ([]*Book, error)
	GenerateBookGroup() error
	StoreBookMark(mark *BookMark) error
	GetBookMarks(bookID string) (*BookMarks, error)
	DeleteBookMark(bookID string, markType MarkType, pageIndex int) error
}

var IStore StoreInterface
