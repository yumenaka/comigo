package model

// x5 书库相关基本操作接口
type StoreInterface interface {
	AddBook(b *Book) error
	GetBook(id string) (*Book, error)
	UpdateBook(b *Book) error
	ListBooks() []*Book
	DeleteBook(id string) error
}

var IStore StoreInterface
