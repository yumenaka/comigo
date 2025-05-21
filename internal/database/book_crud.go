//go:build !((windows && 386) || js)

package database // Package database 编译条件的注释和 package 语句之间一定要隔一行，不然无法识别编译条件。go:build 是1.18以后“条件编译”的推荐语法。

import (
	"context"
	"errors"

	"github.com/yumenaka/comigo/internal/ent"
	entbook "github.com/yumenaka/comigo/internal/ent/book"
	"github.com/yumenaka/comigo/internal/ent/singlepageinfo"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

// ClearBookData   清空数据库的Book与SinglePageInfo表  // 后台并发执行，所以不能保证结果如预期，不用这个函数???
func ClearBookData(clearBook *model.Book) {
	// 如何增删查改： https://entgo.io/zh/docs/crud
	ctx := context.Background()
	_, err := client.Book.
		Delete().
		Where(entbook.BookIDEQ(clearBook.BookID)).
		Exec(ctx)
	if err != nil {
		logger.Infof("%s", "ClearBookData Book:"+err.Error())
	}
	logger.Infof("%s", "Clear Book ："+clearBook.Title)
	deletePageInfoNum, err := client.SinglePageInfo.
		Delete().
		Where(singlepageinfo.BookIDEQ(clearBook.BookID)).
		Exec(ctx)
	if err != nil {
		logger.Infof("ClearBookData SinglePageInfo:%s", err.Error())
	}
	logger.Infof("Clear SinglePageInfo Num：%d", deletePageInfoNum)
}

// DeleteAllBookInDatabase  清空数据库的Book与SinglePageInfo表
// 后台并发执行，不能保证结果如预期，不用这个函数。
func DeleteAllBookInDatabase(debug bool) {
	// 如何增删查改： https://entgo.io/zh/docs/crud
	ctx := context.Background()
	deleteBookNum, err := client.Book.
		Delete().
		Where(entbook.PageCountNEQ(-99999)).
		Exec(ctx)
	if err != nil {
		logger.Infof("%s", err)
	}
	if debug {
		logger.Infof("Delete Book Num：%d", deleteBookNum)
	}
	deletePageInfoNum, err := client.SinglePageInfo.
		Delete().
		Where(singlepageinfo.WidthNEQ(-99999)).
		Exec(ctx)
	if err != nil {
		logger.Infof("%s", err)
	}
	if debug {
		logger.Infof("Delete SinglePageInfo Num：%d", deletePageInfoNum)
	}
}

// SaveAllBookToDatabase 将Map里面的书籍信息，全部保存到本地数据库中
func SaveAllBookToDatabase(m map[string]*model.Book) {
	for _, b := range m {
		c := *b
		err := SaveBookToDatabase(&c)
		if err != nil {
			logger.Infof("SaveAllBookToDatabase error :%s", err.Error())
		}
	}
}

// SaveBookListToDatabase  向数据库中插入一组书
func SaveBookListToDatabase(bookList []*model.Book) error {
	for _, b := range bookList {
		err := SaveBookToDatabase(b)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveBookToDatabase 向数据库中插入一本书
func SaveBookToDatabase(save *model.Book) error {
	// 如何增删查改： https://entgo.io/zh/docs/crud
	ctx := context.Background()
	b, err := client.Book.
		Create().
		SetTitle(save.BookInfo.Title).
		SetBookID(save.BookInfo.BookID).
		SetOwner("").
		SetFilePath(save.BookInfo.FilePath).
		SetBookStorePath(save.BookInfo.BookStorePath).
		SetChildBookNum(save.BookInfo.ChildBookNum).
		SetType(string(save.BookInfo.Type)).
		SetDepth(save.BookInfo.Depth).
		SetParentFolder(save.BookInfo.ParentFolder).
		SetPageCount(save.BookInfo.PageCount).
		SetSize(save.BookInfo.FileSize).
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
		// log.Fatalf("failed creating book: %v", err)
		return err
	}

	// 保存封面与页面信息
	bulk := make([]*ent.SinglePageInfoCreate, len(save.Pages.Images))
	for i, p := range save.Pages.Images {
		bulk[i] = client.SinglePageInfo.
			Create().
			SetBookID(save.BookID).
			SetPageNum(p.PageNum).
			SetName(p.Name).
			SetURL(p.Url).
			SetBlurHash(p.Blurhash).
			SetHeight(p.Height).
			SetWidth(p.Width).
			SetModTime(p.ModTime).
			SetSize(p.Size)
	}
	pages, err := client.SinglePageInfo.CreateBulk(bulk...).Save(ctx)
	if b != nil && pages != nil {
		// log.Println("book was created: ", b)
		// log.Println("book pages info was created: ", pages)
	}
	return err
}

// GetBookFromDatabase 根据文件路径，从数据库查询一本书的详细信息,避免重复扫描压缩包
func GetBookFromDatabase(filepath string) (*model.Book, error) {
	ctx := context.Background()
	books, err := client.Book. // UserClient.
					Query(). // 用户查询生成器。
					Where(entbook.FilePath(filepath)).
					All(ctx) // query and return.
	if err != nil {
		logger.Infof("%s", err)
	}
	if len(books) == 0 {
		return nil, errors.New("not found in database,filepath:" + filepath)
	}
	temp := books[0]
	b := model.Book{
		BookInfo: model.BookInfo{
			Title:           temp.Title,
			BookID:          temp.BookID,
			FilePath:        temp.FilePath,
			BookStorePath:   temp.BookStorePath,
			Type:            model.SupportFileType(temp.Type),
			ChildBookNum:    temp.ChildBookNum,
			Depth:           temp.Depth,
			ParentFolder:    temp.ParentFolder,
			PageCount:       temp.PageCount,
			FileSize:        temp.Size,
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

	// 查询数据库里的封面与页面信息
	// https://entgo.io/zh/docs/crud
	pages, err := client.SinglePageInfo. // UserClient.
						Query(). // 用户查询生成器。
						Where(singlepageinfo.BookID(temp.BookID)).
						All(ctx) // query and return.
	for _, v := range pages {
		b.Pages.Images = append(b.Pages.Images, model.MediaFileInfo{
			PageNum:  v.PageNum,
			Name:     v.Name,
			Url:      v.URL,
			Blurhash: v.BlurHash,
			Height:   v.Height,
			Width:    v.Width,
			ModTime:  v.ModTime,
			Size:     v.Size,
			Path:     v.Path,
			ImgType:  v.ImgType,
		})
	}
	if err != nil {
		logger.Infof("%s", err)
	}
	return &b, err
}

// GetBooksFromDatabase  根据文件路径，从数据库查询书的详细信息,避免重复扫描压缩包。//忽略文件夹型的书籍
func GetBooksFromDatabase() (list []*model.Book, err error) {
	ctx := context.Background()
	books, err := client.Book. // UserClient.
					Query(). // 用户查询生成器。
		// Where(ent_book.Not(ent_book.Type("dir"))). //忽略文件夹型的书籍
		All(ctx) // query and return.
	if err != nil {
		logger.Infof("%s", err)
	}
	if len(books) == 0 {
		return nil, errors.New("not found in database")
	}
	for _, temp := range books {
		b := model.Book{
			BookInfo: model.BookInfo{
				Title:           temp.Title,
				BookID:          temp.BookID,
				FilePath:        temp.FilePath,
				BookStorePath:   temp.BookStorePath,
				Type:            model.SupportFileType(temp.Type),
				ChildBookNum:    temp.ChildBookNum,
				Depth:           temp.Depth,
				ParentFolder:    temp.ParentFolder,
				PageCount:       temp.PageCount,
				FileSize:        temp.Size,
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
		// 查询数据库里的封面与页面信息
		// https://entgo.io/zh/docs/crud
		pages, err := client.SinglePageInfo. // UserClient.
							Query(). // 用户查询生成器。
							Where(singlepageinfo.BookID(temp.BookID)).
							All(ctx) // query and return.
		if err != nil {
			logger.Infof("%s", err)
		}
		for _, v := range pages {
			b.Pages.Images = append(b.Pages.Images, model.MediaFileInfo{
				PageNum:  v.PageNum,
				Name:     v.Name,
				Url:      v.URL,
				Blurhash: v.BlurHash,
				Height:   v.Height,
				Width:    v.Width,
				ModTime:  v.ModTime,
				Size:     v.Size,
				Path:     v.Path,
				ImgType:  v.ImgType,
			})
		}
		list = append(list, &b)
	}
	return list, err
}

// todo： 根据扫描完成的书籍数据，清理本地数据库当中不存在的书籍
