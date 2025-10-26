package model

// x6 书库相关基本操作接口
type StoreInterface interface {
	AddBook(b *Book) error
	GetBook(id string) (*Book, error)
	UpdateBook(b *Book) error
	DeleteBook(id string) error
	ListBooks() ([]*Book, error)
	GenerateBookGroup() error
}

var IStore StoreInterface
