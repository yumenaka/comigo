//go:build !(windows && 386)

// Package storage 编译条件的注释和 package 语句之间一定要隔一行，不然无法识别编译条件。go:build 是1.18以后“条件编译”的推荐语法。
package storage

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strconv"

	"entgo.io/ent/dialect"
	"modernc.org/sqlite"

	comigoBook "github.com/yumenaka/comi/book"
	"github.com/yumenaka/comi/ent"
	"github.com/yumenaka/comi/ent/book"
	"github.com/yumenaka/comi/ent/singlepageinfo"
	"github.com/yumenaka/comi/locale"
)

// 参考：
// Go製CGOフリーなSQLiteドライバーでentを使う
// https://zenn.dev/nobonobo/articles/e9f17d183c19f6

// 数据库为sqlite3
// 查看工具：SQLiteStudio https://github.com/pawelsalawa/sqlitestudio/releases
// 查看工具： DB Browser for SQLite  https://sqlitebrowser.org/dl/
type sqliteDriver struct {
	*sqlite.Driver
}

func (d sqliteDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.Driver.Open(name)
	if err != nil {
		return conn, err
	}
	c := conn.(interface {
		Exec(stmt string, args []driver.Value) (driver.Result, error)
	})
	if _, err := c.Exec("PRAGMA foreign_keys = on;", nil); err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to enable enable foreign keys: %w", err)
	}
	return conn, nil
}

// 注册 sqlite
func init() {
	sql.Register("sqlite3", sqliteDriver{Driver: &sqlite.Driver{}})
}

var client *ent.Client

func InitDatabase(configFilePath string) error {

	if client != nil {
		fmt.Println("database already initialized")
		return nil
	}
	//链接或创建数据库
	var entOptions []ent.Option
	//是否打印log
	//entOptions = append(entOptions, ent.Debug())
	//连接器
	var err error
	dataSourceName := "file:comigo.sqlite?cache=shared"
	//如果有配置文件的话，数据库文件，就在同一文件夹内
	if configFilePath != "" {
		configDir := filepath.Dir(configFilePath) //不能用path.Dir()，因为windows返回 "."
		dataSourceName = "file:" + path.Join(configDir, "comigo.sqlite") + "?cache=shared"
	}
	fmt.Println(locale.GetString("InitDatabase") + dataSourceName)
	client, err = ent.Open(dialect.SQLite, dataSourceName, entOptions...)
	if err != nil {
		return fmt.Errorf("failed opening connection to sqlite: %v", err)
		//time.Sleep(3 * time.Second)
		//log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	//defer client.Close()
	if err := client.Schema.Create(context.Background()); err != nil {
		return fmt.Errorf("failed creating schema resources: %v", err)
		//time.Sleep(3 * time.Second)
		//log.Fatalf("failed creating schema resources: %v", err)
	}
	return nil
}

func CloseDatabase() {
	err := client.Close()
	if err != nil {
		fmt.Println(err)
	}
}

// ClearBookData   清空数据库的Book与SinglePageInfo表
// 后台并发执行，所以不能保证结果如预期，不用这个函数。
func ClearBookData(clearBook *comigoBook.Book, debug bool) {
	//如何增删查改： https://entgo.io/zh/docs/crud
	ctx := context.Background()
	_, err := client.Book.
		Delete().
		Where(book.BookIDEQ(clearBook.BookID)).
		Exec(ctx)
	if err != nil {
		fmt.Println("ClearBookData Book:" + err.Error())
	}
	if debug {
		fmt.Println("Clear Book ：" + clearBook.Name)
	}
	deletePageInfoNum, err := client.SinglePageInfo.
		Delete().
		Where(singlepageinfo.BookIDEQ(clearBook.BookID)).
		Exec(ctx)
	if err != nil {
		fmt.Println("ClearBookData SinglePageInfo:" + err.Error())
	}
	if debug {
		fmt.Println("Clear SinglePageInfo Num：" + strconv.Itoa(deletePageInfoNum))
	}
}

// DeleteAllBookInDatabase  清空数据库的Book与SinglePageInfo表
// 后台并发执行，不能保证结果如预期，不用这个函数。
func DeleteAllBookInDatabase(debug bool) {
	//如何增删查改： https://entgo.io/zh/docs/crud
	ctx := context.Background()
	deleteBookNum, err := client.Book.
		Delete().
		Where(book.AllPageNumNEQ(-99999)).
		Exec(ctx)
	if err != nil {
		fmt.Println(err)
	}
	if debug {
		fmt.Println("Delete Book Num：" + strconv.Itoa(deleteBookNum))
	}
	deletePageInfoNum, err := client.SinglePageInfo.
		Delete().
		Where(singlepageinfo.WidthNEQ(-99999)).
		Exec(ctx)
	if err != nil {
		fmt.Println(err)
	}
	if debug {
		fmt.Println("Delete SinglePageInfo Num：" + strconv.Itoa(deletePageInfoNum))
	}
}

// SaveAllBookToDatabase 将Map里面的书籍信息，全部保存到本地数据库中
func SaveAllBookToDatabase(databaseFilePath string, m map[string]*comigoBook.Book) {
	for _, b := range m {
		var c = *b
		err := SaveBookToDatabase(&c)
		if err != nil {
			fmt.Println("SaveAllBookToDatabase error :" + err.Error())
		}
	}
}

// SaveBookListToDatabase  向数据库中插入一组书
func SaveBookListToDatabase(bookList []*comigoBook.Book) error {
	for _, b := range bookList {
		err := SaveBookToDatabase(b)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveBookToDatabase 向数据库中插入一本书
func SaveBookToDatabase(save *comigoBook.Book) error {
	//如何增删查改： https://entgo.io/zh/docs/crud
	ctx := context.Background()
	b, err := client.Book.
		Create().
		SetName(save.BookInfo.Name).
		SetBookID(save.BookInfo.BookID).
		SetOwner("").
		SetFilePath(save.BookInfo.FilePath).
		SetBookStorePath(save.BookInfo.BookStorePath).
		SetChildBookNum(save.BookInfo.ChildBookNum).
		SetType(string(save.BookInfo.Type)).
		SetDepth(save.BookInfo.Depth).
		SetParentFolder(save.BookInfo.ParentFolder).
		SetAllPageNum(save.BookInfo.AllPageNum).
		SetFileSize(save.BookInfo.FileSize).
		SetAuthors(save.GetAuthor()).
		SetISBN(save.BookInfo.ISBN).
		SetPress(save.BookInfo.Press).
		SetPublishedAt(save.BookInfo.PublishedAt).
		SetExtractPath(save.BookInfo.ExtractPath).
		SetInitComplete(save.BookInfo.InitComplete).
		SetReadPercent(save.BookInfo.ReadPercent).
		SetNonUTF8Zip(save.BookInfo.NonUTF8Zip).
		SetZipTextEncoding(save.BookInfo.ZipTextEncoding).
		SetExtractNum(save.BookInfo.ExtractNum).
		Save(ctx) // 创建并返回 //还有一个SaveX(ctx)，和 Save() 不一样， SaveX 在出错时 panic。
	if err != nil {
		//log.Fatalf("failed creating book: %v", err)
		return err
	}

	//保存封面与页面信息
	bulk := make([]*ent.SinglePageInfoCreate, len(save.Pages.Images))
	for i, p := range save.Pages.Images {
		bulk[i] = client.SinglePageInfo.
			Create().
			SetBookID(save.BookID).
			SetPageNum(p.PageNum).
			SetNameInArchive(p.NameInArchive).
			SetURL(p.Url).
			SetBlurHash(p.Blurhash).
			SetHeight(p.Height).
			SetWidth(p.Width).
			SetModeTime(p.ModeTime).
			SetFileSize(p.FileSize).
			SetRealImageFilePATH(p.RealImageFilePATH).
			SetImgType(p.ImgType)
	}
	pages, err := client.SinglePageInfo.CreateBulk(bulk...).Save(ctx)
	if b != nil && pages != nil {
		//log.Println("book was created: ", b)
		//log.Println("book pages info was created: ", pages)
	}
	return err
}

// GetBookFromDatabase 根据文件路径，从数据库查询一本书的详细信息,避免重复扫描压缩包。
func GetBookFromDatabase(filepath string) (*comigoBook.Book, error) {
	ctx := context.Background()
	books, err := client.Book. // UserClient.
					Query(). // 用户查询生成器。
					Where(book.FilePath(filepath)).
					All(ctx) // query and return.
	if err != nil {
		fmt.Println(err)
	}
	if len(books) == 0 {
		return nil, errors.New("not found in database,filepath:" + filepath)
	}
	temp := books[0]
	b := comigoBook.Book{
		BookInfo: comigoBook.BookInfo{
			Name:            temp.Name,
			BookID:          temp.BookID,
			FilePath:        temp.FilePath,
			BookStorePath:   temp.BookStorePath,
			ChildBookNum:    temp.ChildBookNum,
			Depth:           temp.Depth,
			ParentFolder:    temp.ParentFolder,
			AllPageNum:      temp.AllPageNum,
			FileSize:        temp.FileSize,
			ISBN:            temp.ISBN,
			Press:           temp.Press,
			PublishedAt:     temp.PublishedAt,
			ExtractPath:     temp.ExtractPath,
			Modified:        temp.Modified,
			ExtractNum:      temp.ExtractNum,
			InitComplete:    temp.InitComplete,
			ReadPercent:     temp.ReadPercent,
			NonUTF8Zip:      temp.NonUTF8Zip,
			ZipTextEncoding: temp.ZipTextEncoding,
		},
	}
	b.Type = comigoBook.GetBookTypeByFilename(temp.FilePath)
	//查询数据库里的封面与页面信息
	//https://entgo.io/zh/docs/crud
	pages, err := client.SinglePageInfo. // UserClient.
						Query(). // 用户查询生成器。
						Where(singlepageinfo.BookID(temp.BookID)).
						All(ctx) // query and return.
	for _, v := range pages {
		b.Pages.Images = append(b.Pages.Images, comigoBook.ImageInfo{
			PageNum:           v.PageNum,
			NameInArchive:     v.NameInArchive,
			Url:               v.URL,
			Blurhash:          v.BlurHash,
			Height:            v.Height,
			Width:             v.Width,
			ModeTime:          v.ModeTime,
			FileSize:          v.FileSize,
			RealImageFilePATH: v.RealImageFilePATH,
			ImgType:           v.ImgType,
		})
	}
	//设置封面
	if len(b.Pages.Images) > 0 {
		b.Cover = b.Pages.Images[0]
	}
	if err != nil {
		fmt.Println(err)
	}
	return &b, err
}

// GetArchiveBookFromDatabase  根据文件路径，从数据库查询书的详细信息,避免重复扫描压缩包。//忽略文件夹型的书籍
func GetArchiveBookFromDatabase() (list []*comigoBook.Book, err error) {
	ctx := context.Background()
	books, err := client.Book. // UserClient.
					Query(). // 用户查询生成器。
		//Where(book.Not(book.Type("dir"))). //忽略文件夹型的书籍
		All(ctx) // query and return.
	if err != nil {
		fmt.Println(err)
	}
	if len(books) == 0 {
		return nil, errors.New("not found in database")
	}
	for _, temp := range books {
		b := comigoBook.Book{
			BookInfo: comigoBook.BookInfo{
				Name:            temp.Name,
				BookID:          temp.BookID,
				FilePath:        temp.FilePath,
				BookStorePath:   temp.BookStorePath,
				ChildBookNum:    temp.ChildBookNum,
				Depth:           temp.Depth,
				ParentFolder:    temp.ParentFolder,
				AllPageNum:      temp.AllPageNum,
				FileSize:        temp.FileSize,
				ISBN:            temp.ISBN,
				Press:           temp.Press,
				PublishedAt:     temp.PublishedAt,
				ExtractPath:     temp.ExtractPath,
				Modified:        temp.Modified,
				ExtractNum:      temp.ExtractNum,
				InitComplete:    temp.InitComplete,
				ReadPercent:     temp.ReadPercent,
				NonUTF8Zip:      temp.NonUTF8Zip,
				ZipTextEncoding: temp.ZipTextEncoding,
			},
		}
		b.Type = comigoBook.GetBookTypeByFilename(temp.FilePath)
		if b.ChildBookNum > 0 {
			b.Type = comigoBook.TypeBooksGroup
		}
		//查询数据库里的封面与页面信息
		//https://entgo.io/zh/docs/crud
		pages, err := client.SinglePageInfo. // UserClient.
							Query(). // 用户查询生成器。
							Where(singlepageinfo.BookID(temp.BookID)).
							All(ctx) // query and return.
		if err != nil {
			fmt.Println(err)
		}
		for _, v := range pages {
			b.Pages.Images = append(b.Pages.Images, comigoBook.ImageInfo{
				PageNum:           v.PageNum,
				NameInArchive:     v.NameInArchive,
				Url:               v.URL,
				Blurhash:          v.BlurHash,
				Height:            v.Height,
				Width:             v.Width,
				ModeTime:          v.ModeTime,
				FileSize:          v.FileSize,
				RealImageFilePATH: v.RealImageFilePATH,
				ImgType:           v.ImgType,
			})
		}
		//设置封面
		if len(b.Pages.Images) > 0 {
			b.Cover = b.Pages.Images[0]
		}
		//硬写一个封面
		switch b.Type {
		case comigoBook.TypePDF:
			b.Cover = comigoBook.ImageInfo{NameInArchive: "pdf.png", Url: "/images/pdf.png"}
		case comigoBook.TypeVideo:
			b.Cover = comigoBook.ImageInfo{NameInArchive: "video.png", Url: "/images/video.png"}
		case comigoBook.TypeAudio:
			b.Cover = comigoBook.ImageInfo{NameInArchive: "audio.png", Url: "/images/audio.png"}
		case comigoBook.TypeUnknownFile:
			b.Cover = comigoBook.ImageInfo{NameInArchive: "unknown.png", Url: "/images/unknown.png"}
		}
		list = append(list, &b)
	}
	return list, err
}

// CleanDatabaseByLocalData 根据扫描完成的书籍数据，清理本地数据库当中不存在的书籍
func CleanDatabaseByLocalData() {

}

//func InitMapBooksByDatabase() {
//	tempMap, err := GetArchiveBookFromDatabase()
//	if err != nil {
//		mapBooks = tempMap
//	}
//}
//

//// CleanAndSaveAllBookToDatabase  同时清空Map里面不存在的书。然后将Map里面的书籍信息，全部保存到本地数据库中。
//func CleanAndSaveAllBookToDatabase(databaseFilePath string, m map[string]*comigoBook.Book) error {
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
