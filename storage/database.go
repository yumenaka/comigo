package storage

import (
	"fmt"
	"path"
	"time"

	"github.com/asdine/storm/v3"
	bolt "go.etcd.io/bbolt"

	"github.com/yumenaka/comi/book"
)

//嵌入式数据库storm的使用:
//https://ystyle.top/2019/12/22/how-to-use-storm-db/

//web查看工具：go install github.com/evnix/boltdbweb@latest
//命令行：go install github.com/br0xen/boltbrowser@latest

//go install github.com/yumenaka/comi@latest

//func InitMapBooksByDatabase() {
//	tempMap, err := GetAllBookFromDatabase()
//	if err != nil {
//		mapBooks = tempMap
//	}
//}
//
////根据扫描完成的书籍数据，覆盖本地数据库
//func CleanMapBooksByLocalData() {
//	err := SaveAllBookToDatabase(mapBooks)
//	if err != nil {
//		fmt.Println(err)
//	}
//}

// GetBookFromDatabase 根据文件路径，从数据库里面取出一本书。用来减少重复扫描压缩包。
func GetBookFromDatabase(filepath string, databasePath string) (*book.Book, error) {
	//Open函数默认创建的db文件权限是0600(属主读写),超时选项,可以自行传值修改
	db, err := storm.Open(path.Join(databasePath, "comigo.db"), storm.BoltOptions(0600, &bolt.Options{Timeout: 50 * time.Millisecond}))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer db.Close()
	var b book.Book
	fmt.Println("从数据库中获取前：", b)
	err = db.One("filepath", filepath, &b)
	fmt.Println("从数据库中获取后：", b)
	return &b, err
}

// GetAllBookFromDatabase 从本地数据库里面取出全部书籍信息，并以Map形式返回
func GetAllBookFromDatabase(databasePath string) (map[string]*book.Book, error) {
	//Open函数默认创建的db文件权限是0600(属主读写),1秒超时的选项,可以自行传值修改
	db, err := storm.Open(path.Join(databasePath, "comigo.db"), storm.BoltOptions(0600, &bolt.Options{Timeout: 100 * time.Millisecond}))
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	var allBook []book.Book
	err = db.All(&allBook)
	if err == nil {
		var temp map[string]*book.Book
		temp = make(map[string]*book.Book)
		for _, b := range allBook {
			temp[b.BookID] = &b
		}
		fmt.Println("成功读取数据库")
		return temp, nil
	}
	return nil, err
}

// SaveAllBookToDatabase 将Map里面的书籍信息，全部保存到本地数据库中
func SaveAllBookToDatabase(databaseFilePath string, m map[string]*book.Book) error {
	//Open函数默认创建的db文件权限是0600(属主读写),1秒超时的选项,可以自行传值修改
	db, err := storm.Open(path.Join(databaseFilePath, "comigo.db"), storm.BoltOptions(0600, &bolt.Options{Timeout: 100 * time.Millisecond}))
	defer db.Close()
	for _, b := range m {
		var book = *b
		db.Save(&book)
	}
	return err
}

// CleanAndSaveAllBookToDatabase  同时清空Map里面不存在的书。然后将Map里面的书籍信息，全部保存到本地数据库中。
func CleanAndSaveAllBookToDatabase(databaseFilePath string, m map[string]*book.Book) error {
	//Open函数默认创建的db文件权限是0600(属主读写),1秒超时的选项,可以自行传值修改
	db, err := storm.Open(path.Join(databaseFilePath, "comigo.db"), storm.BoltOptions(0600, &bolt.Options{Timeout: 100 * time.Millisecond}))
	if err != nil {
		return err
	}
	defer db.Close()
	////Drop方法 删除表
	//err = db.Drop(&Book{})
	//if err != nil {
	//	return err
	//}
	////Init方法 初始化表
	//err = db.Init(&Book{})

	//再次添加
	for _, b := range m {
		var book = *b
		db.Save(&book)
	}
	return err
}
