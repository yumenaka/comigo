//go:build windows && 386

package database

import (
	"errors"

	"github.com/yumenaka/comi/logger"
	mainTypes "github.com/yumenaka/comi/types"
)

func InitDatabase(configFilePath string) error {
	logger.Info("Not Support DateBase")
	return nil
}

func CloseDatabase() {
	logger.Info("Not Support DateBase")
}

// DeleteAllBookInDatabase  清空数据库的Book与SinglePageInfo表
// 后台并发执行，所以不能保证结果如预期，不用这个函数。
func ClearBookData(clearBook *mainTypes.Book) {
	logger.Info("Not Support DateBase")
}

func GetBooksFromDatabase() (list []*mainTypes.Book, err error) {
	return nil, errors.New("Not Support DateBase")
}

// DeleteAllBookInDatabase  清空数据库的Book与SinglePageInfo表
// 后台并发执行，不能保证结果如预期，不用这个函数。
func DeleteAllBookInDatabase(debug bool) {
	logger.Info("Not Support DateBase")
}

// SaveAllBookToDatabase 将Map里面的书籍信息，全部保存到本地数据库中
func SaveAllBookToDatabase(databaseFilePath string, m map[string]*mainTypes.Book) {
	logger.Info("Not Support DateBase")
}

// SaveBookListToDatabase  向数据库中插入一组书
func SaveBookListToDatabase(bookList []*mainTypes.Book) error {
	return errors.New("Not Support DateBase")
}

// SaveBookToDatabase 向数据库中插入一本书
func SaveBookToDatabase(save *mainTypes.Book) error {
	return errors.New("Not Support DateBase")
}

// GetBookFromDatabase 根据文件路径，从数据库查询一本书的详细信息,避免重复扫描压缩包。
func GetBookFromDatabase(filepath string) (*mainTypes.Book, error) {
	return nil, errors.New("Not Support DateBase")
}

// GetBooksFromDatabase  根据文件路径，从数据库查询书的详细信息,避免重复扫描压缩包。//忽略文件夹型的书籍
func GetArchiveBookFromDatabase() (list []*mainTypes.Book, err error) {
	return nil, errors.New("Not Support DateBase")
}

// CleanDatabaseByLocalData 根据扫描完成的书籍数据，清理本地数据库当中不存在的书籍
func CleanDatabaseByLocalData() {
	logger.Info("Not Support DateBase")
}

//func InitMapBooksByDatabase() {
//	tempMap, err := GetBooksFromDatabase()
//	if err != nil {
//		mapBooks = tempMap
//	}
//}
//

//// CleanAndSaveAllBookToDatabase  同时清空Map里面不存在的书。然后将Map里面的书籍信息，全部保存到本地数据库中。
//func CleanAndSaveAllBookToDatabase(databaseFilePath string, m map[string]*mainTypes.Book) error {
//	//Open函数默认创建的db文件权限是0600(属主读写),1秒超时的选项,可以自行传值修改
//	//db, err := storm.Open(path.Join(databaseFilePath, "comigo.db"), storm.BoltOptions(0600, &bolt.Options{Timeout: 100 * time.Millisecond}))
//	//if err != nil {
//	//	return err
//	//}
//	//defer db.Close()
//	////Drop方法 删除表
//	//err = db.Drop(&Book{})
//	//if err != nil {
//	//	return err
//	//}
//	////Init方法 初始化表
//	//err = db.Init(&Book{})
//
//	////再次添加
//	//for _, b := range m {
//	//	var book = *b
//	//	db.Save(&book)
//	//}
//	return err
//}
