package main

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"entgo.io/ent/dialect"
	"github.com/disintegration/imaging"
	"github.com/yumenaka/archiver/v4"
	"github.com/yumenaka/comi/ent"
	"github.com/yumenaka/comi/logger"
	"modernc.org/sqlite"
)

// 参考：
// Go製CGOフリーなSQLiteドライバーでentを使う
// https://zenn.dev/nobonobo/articles/e9f17d183c19f6

// 初始化数据库为sqlite3
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

func init() {
	sql.Register("sqlite3", sqliteDriver{Driver: &sqlite.Driver{}})
}

// func main() {
// 	testDatabase()
// }

func testDatabase() {
	// 链接或创建数据库
	var entOptions []ent.Option
	entOptions = append(entOptions, ent.Debug())
	// 连接器
	client, err := ent.Open(dialect.SQLite, "file:comigo.sqlite?cache=shared", entOptions...)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer func(client *ent.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalf("failed closing connection to sqlite: %v", err)
		}
	}(client)

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	// 弄个随机数、为那些不能重复的字段加个随机后缀
	r := rand.New(rand.NewSource(time.Now().Unix()))
	randString := strconv.Itoa(r.Intn(100000))

	// 如何增删查改： https://entgo.io/zh/docs/crud

	// 插入一个User
	ctx := context.Background()
	u, err := client.User.
		Create().
		SetAge(32).
		SetUsername("test Username" + randString).
		SetPassword("12345").
		SetName("Test Admin" + randString).
		Save(ctx)
	if err != nil {
		log.Fatalf("failed creating user: %v", err)
	}
	log.Println("user was created: ", u)

	// 插入一本书
	b, err := client.Book.
		Create().
		SetTitle("Test Book title" + randString).
		SetBookID("BookID" + randString).
		SetFilePath("path" + randString).
		SetBookStorePath("path2").
		SetChildBookNum(0).
		SetType("zip").
		SetDepth(1).
		SetParentFolder("ParentPath").
		SetPageCount(99).
		SetFileSize(66).
		SetAuthors("unknown").
		SetISBN("").
		SetPress("").
		SetPublishedAt("").
		SetExtractPath("").
		SetInitComplete(true).
		SetReadPercent(0.01).
		SetNonUTF8Zip(false).
		SetZipTextEncoding("").
		SetExtractNum(0).
		Save(ctx) // 创建并返回 //还有一个SaveX(ctx)，和 Save() 不一样， SaveX 在出错时 panic。
	if err != nil {
		log.Fatalf("failed creating book: %v", err)
	}
	log.Println("book was created: ", b)
}

func testPDF() {
	//pageCount, err := CountPagesOfPDFFile("01.pdf")
	//if err != nil {
	//	logger.Infof("%s", err)
	//}
	//for i := 0; i < pageCount; i++ {
	//	ExportImageFromPDF("01.pdf", i+1)
	//}
	//ExportAllImageFromPDF("01.pdf")
}

func ImageResize() {
	// 读取本地文件，本地文件尺寸300*400
	imgData, _ := os.ReadFile("d:/1.jpg")
	buf := bytes.NewBuffer(imgData)
	image, err := imaging.Decode(buf)
	if err != nil {
		logger.Infof("%s", err)
		return
	}
	// 生成缩略图，尺寸150*200，并保持到为文件2.jpg
	image = imaging.Resize(image, 150, 200, imaging.Lanczos)
	err = imaging.Save(image, "d:/2.jpg")
	if err != nil {
		logger.Infof("%s", err)
	}
}

// UnArchiveZip 一次性解压zip文件
func UnArchiveZip(filePath string, extractPath string, textEncoding string) error {
	extractPath = getAbsPath(extractPath)
	// 如果解压路径不存在，创建路径
	err := os.MkdirAll(extractPath, os.ModePerm)
	if err != nil {
		logger.Infof("%s", err)
	}
	// 打开文件，只读模式
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0o400) // Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
	if err != nil {
		logger.Infof("%s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Infof("%s", err)
		}
	}(file)
	// 是否是压缩包
	format, _, err := archiver.Identify(filePath, file)
	if err != nil {
		return err
	}
	// 如果是zip
	if ex, ok := format.(archiver.Zip); ok {
		ex.TextEncoding = textEncoding // “”  "shiftjis" "gbk"
		ctx := context.Background()
		// WithValue返回parent的一个副本，该副本保存了传入的key/value，而调用Context接口的Value(key)方法就可以得到val。注意在同一个context中设置key/value，若key相同，值会被覆盖。
		ctx = context.WithValue(ctx, "extractPath", extractPath)
		_, err := ex.LsAllFile(ctx, file, extractFileHandler)
		if err != nil {
			return err
		}
		logger.Infof("zip文件解压完成：%s 解压到：%s", getAbsPath(filePath), getAbsPath(extractPath))
	}
	return nil
}

//// UnArchiveFle 一次性解压rar文件
//func UnArchiveRar(filePath string, extractPath string) error {
//	extractPath = getAbsPath(extractPath)
//	//如果解压路径不存在，创建路径
//	err := os.MkdirAll(extractPath, os.ModePerm)
//	if err != nil {
//		logger.Infof("%s", err)
//	}
//	//打开文件，只读模式
//	file, err := os.OpenFile(filePath, os.O_RDONLY, 0400) //Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
//	if err != nil {
//		logger.Infof("%s", err)
//	}
//	defer file.Close()
//	//是否是压缩包
//	format, _, err := archiver.Identify(filePath, file)
//	if err != nil {
//		return err
//	}
//	//如果是rar
//	if ex, ok := format.(archiver.Rar); ok {
//		ctx := context.Background()
//		//WithValue返回parent的一个副本，该副本保存了传入的key/value，而调用Context接口的Value(key)方法就可以得到val。注意在同一个context中设置key/value，若key相同，值会被覆盖。
//		ctx = context.WithValue(ctx, "extractPath", extractPath)
//		err := ex.LsAllFile(ctx, file, extractFileHandler)
//		if err != nil {
//			return err
//		}
//		logger.Infof("rar文件解压完成：" + getAbsPath(filePath) + " 解压到：" + getAbsPath(extractPath))
//	}
//	return nil
//}

// 解压文件的函数
func extractFileHandler(ctx context.Context, f archiver.File) error {
	extractPath := ""
	if e, ok := ctx.Value("extractPath").(string); ok {
		extractPath = e
	}
	// 取得压缩文件
	file, err := f.Open()
	if err != nil {
		logger.Infof("%s", err)
	}
	defer func(file io.ReadCloser) {
		err := file.Close()
		if err != nil {
			logger.Infof("%s", err)
		}
	}(file)
	// 如果是文件夹，直接创建文件夹
	if f.IsDir() {
		err = os.MkdirAll(filepath.Join(extractPath, f.NameInArchive), os.ModePerm)
		if err != nil {
			logger.Infof("%s", err)
		}
		return err
	}
	// 如果是一般文件，将文件写入磁盘
	writeFilePath := filepath.Join(extractPath, f.NameInArchive)
	// 写文件前，如果对应文件夹不存在，就创建对应文件夹
	checkDir := filepath.Dir(writeFilePath)
	if !isExist(checkDir) {
		err = os.MkdirAll(checkDir, os.ModePerm)
		if err != nil {
			logger.Infof("%s", err)
		}
		return err
	}
	// 具体内容
	content, err := io.ReadAll(file)
	if err != nil {
		logger.Infof("%s", err)
	}
	// 写入文件
	err = os.WriteFile(writeFilePath, content, 0o644)
	if err != nil {
		logger.Infof("%s", err)
	}
	return err
}

// 判断文件夹或文件是否存在
func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		logger.Infof("%s", err)
		return false
	}
	return true
}

// 获取绝对路径
func getAbsPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		abs = path
	}
	return abs
}
