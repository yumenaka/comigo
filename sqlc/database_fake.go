//go:build (windows && 386) || (js && wasm)

package sqlc

import (
	"errors"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools/logger"
)

func InitDatabase(configFilePath string) error {
	logger.Infof("%s", "Not Support DateBase")
	return nil
}

func CloseDatabase() {
	logger.Infof("%s", "Not Support DateBase")
}

func (repo *Repository) GetBooksFromDatabase() (list []*model.Book, err error) {
	return nil, errors.New("Not Support DateBase")
}

// DeleteAllBookInDatabase  清空数据库的Book与SinglePageInfo表
// 后台并发执行，所以不能保证结果如预期，不用这个函数。
func ClearBookData(clearBook *model.Book) {
	logger.Infof("%s", "Not Support DateBase")
}

func GetBooksFromDatabase() (list []*model.Book, err error) {
	return nil, errors.New("Not Support DateBase")
}

// DeleteAllBookInDatabase  清空数据库的Book与SinglePageInfo表
// 后台并发执行，不能保证结果如预期，不用这个函数。
func DeleteAllBookInDatabase(debug bool) {
	logger.Infof("%s", "Not Support DateBase")
}

// SaveAllBookToDatabase 将Map里面的书籍信息，全部保存到本地数据库中
func SaveAllBookToDatabase(databaseFilePath string, m map[string]*model.Book) {
	logger.Infof("%s", "Not Support DateBase")
}

// SaveBookListToDatabase  向数据库中插入一组书
func SaveBookListToDatabase(bookList []*model.Book) error {
	return errors.New("Not Support DateBase")
}

// SaveBookToDatabase 向数据库中插入一本书
func SaveBookToDatabase(save *model.Book) error {
	return errors.New("Not Support DateBase")
}

// GetBookFromDatabase 根据文件路径，从数据库查询一本书的详细信息,避免重复扫描压缩包。
func GetBookFromDatabase(filepath string) (*model.Book, error) {
	return nil, errors.New("Not Support DateBase")
}

// GetBooksFromDatabase  根据文件路径，从数据库查询书的详细信息,避免重复扫描压缩包。//忽略文件夹型的书籍
func GetArchiveBookFromDatabase() (list []*model.Book, err error) {
	return nil, errors.New("Not Support DateBase")
}

// CleanDatabaseByLocalData 根据扫描完成的书籍数据，清理本地数据库当中不存在的书籍
func CleanDatabaseByLocalData() {
	logger.Infof("%s", "Not Support DateBase")
}
