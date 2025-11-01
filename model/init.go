package model

// StoreInterface 书库相关基本操作接口 x6
type StoreInterface interface {
	AddBook(b *Book) error
	GetBook(id string) (*Book, error)
	UpdateBook(b *Book) error
	DeleteBook(id string) error
	ListBooks() ([]*Book, error)
	GenerateBookGroup() error
}
