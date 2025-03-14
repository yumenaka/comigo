package scan

import (
	"io/fs"
	"net/url"
	"path"
	"strings"

	"github.com/yumenaka/comigo/config/stores"
	"github.com/yumenaka/comigo/internal/database"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/file"
	"github.com/yumenaka/comigo/util/locale"
	"github.com/yumenaka/comigo/util/logger"
)

type Config interface {
	GetLocalStores() []string
	GetStores() []stores.Store
	GetMaxScanDepth() int
	GetMinImageNum() int
	GetTimeoutLimitForScan() int
	GetExcludePath() []string
	GetSupportMediaType() []string
	GetSupportFileType() []string
	GetSupportTemplateFile() []string
	GetZipFileTextEncoding() string
	GetEnableDatabase() bool
	GetClearDatabaseWhenExit() bool
	GetDebug() bool
}
type Option struct {
	ReScanFile bool // 是否重新扫描文件
	Cfg        Config
}

func NewOption(
	reScanFile bool,
	scanConfig Config,
) Option {
	return Option{
		ReScanFile: reScanFile,
		Cfg:        scanConfig,
	}
}

// IsSupportTemplate 判断压缩包内的文件是否是支持的模板文件
func (o *Option) IsSupportTemplate(checkPath string) bool {
	for _, ex := range o.Cfg.GetSupportTemplateFile() {
		suffix := strings.ToLower(path.Ext(checkPath)) //strings.ToLower():某些文件会用大写文件名
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSupportMedia 判断文件是否需要展示
func (o *Option) IsSupportMedia(checkPath string) bool {
	for _, ex := range o.Cfg.GetSupportMediaType() {
		suffix := strings.ToLower(path.Ext(checkPath)) //strings.ToLower():某些文件会用大写文件名
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSupportFile 判断压缩包内的文件是否需要展示
func (o *Option) IsSupportFile(checkPath string) bool {
	for _, ex := range o.Cfg.GetSupportFileType() {
		suffix := strings.ToLower(path.Ext(checkPath)) //strings.ToLower():某些文件会用大写文件名
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSupportArchiver 是否是支持的压缩文件
func (o *Option) IsSupportArchiver(checkPath string) bool {
	for _, ex := range o.Cfg.GetSupportFileType() {
		suffix := path.Ext(checkPath)
		if ex == suffix {
			return true
		}
	}
	return false
}

// IsSkipDir  检查路径是否应该跳过（排除文件，文件夹列表）。
func (o *Option) IsSkipDir(path string) bool {
	for _, substr := range o.Cfg.GetExcludePath() {
		if strings.HasSuffix(path, substr) {
			return true
		}
	}
	return false
}

// AllStore 3、扫描路径，取得路径里的书籍
func AllStore(option Option) error {
	// 重置所有书籍与书组信息
	model.ClearAllBookData()
	for _, localPath := range option.Cfg.GetLocalStores() {
		addList, err := ScanDirectory(localPath, option)
		if err != nil {
			logger.Infof(locale.GetString("scan_error")+" path:%s %s", localPath, err)
			continue
		}
		AddBooksToStore(addList, localPath, option.Cfg.GetMinImageNum())
	}
	//for _, server := range scanStores {
	//	addList, err := Smb(option)
	//	if err != nil {
	//		logger.Infof("smb scan_error"+" path:%s %s", server.ShareName, err)
	//		continue
	//	}
	//	AddBooksToStore(addList, server.ShareName, scanMinImageNum)
	//}
	return nil
}

// SaveResultsToDatabase 4，保存扫描结果到数据库，并清理不存在的书籍
func SaveResultsToDatabase(ConfigPath string, ClearDatabaseWhenExit bool) error {
	err := database.InitDatabase(ConfigPath)
	if err != nil {
		return err
	}
	saveErr := database.SaveBookListToDatabase(model.GetArchiveBooks())
	if saveErr != nil {
		logger.Info(saveErr)
		return saveErr
	}
	return nil
}

func ClearDatabaseWhenExit(ConfigPath string) {
	AllBook := model.GetAllBookList()
	for _, b := range AllBook {
		database.ClearBookData(b)
	}
}

// AddBooksToStore 添加一组书到书库
func AddBooksToStore(bookList []*model.Book, basePath string, MinImageNum int) {
	err := model.AddBooks(bookList, basePath, MinImageNum)
	if err != nil {
		logger.Infof(locale.GetString("AddBook_error")+"%s", basePath)
	}
	// 生成虚拟书籍组
	if err := model.MainStore.AnalyzeStore(); err != nil {
		logger.Infof("%s", err)
	}
}

func scanNonUTF8ZipFile(filePath string, b *model.Book, option Option) error {
	b.NonUTF8Zip = true
	reader, err := file.ScanNonUTF8Zip(filePath, option.Cfg.GetZipFileTextEncoding())
	if err != nil {
		return err
	}
	for _, f := range reader.File {
		if option.IsSupportMedia(f.Name) {
			// 如果是压缩文件
			// 替换特殊字符的时候，额外将“+替换成"%2b"，因为gin会将+解析为空格。
			TempURL := "/api/get_file?id=" + b.BookID + "&filename=" + url.QueryEscape(f.Name)
			b.Pages.Images = append(b.Pages.Images, model.MediaFileInfo{Path: "", Size: f.FileInfo().Size(), ModTime: f.FileInfo().ModTime(), Name: f.Name, Url: TempURL})
		} else {
			if option.Cfg.GetDebug() {
				logger.Infof(locale.GetString("unsupported_file_type")+" %s", f.Name)
			}
		}
	}
	b.SortPages("default")
	return err
}

// 手动写的递归查找，功能与fs.WalkDir()相同。发现一个Archiver/V4的BUG：zip文件的虚拟文件系统，找不到正确的多级文件夹？
// https://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter06/06.3.html
func walkUTF8ZipFs(fsys fs.FS, parent, base string, b *model.Book, option Option) error {
	// 一般zip文件的处理流程
	// logger.Infof("parent:" + parent + " base:" + base)
	dirName := path.Join(parent, base)
	dirEntries, err := fs.ReadDir(fsys, dirName)
	for _, dirEntry := range dirEntries {
		name := dirEntry.Name()
		f, errInfo := dirEntry.Info()
		if errInfo != nil {
			continue
		}
		if dirEntry.IsDir() == true {
			switch name {
			case ".comigo":
				return fs.SkipDir
			case "flutter_ui":
				return fs.SkipDir
			case "node_modules":
				return fs.SkipDir
			default:
			}
			joinPath := path.Join(parent, name)
			err = walkUTF8ZipFs(fsys, joinPath, base, b, option)
		} else if option.IsSupportMedia(name) {
			inArchiveName := path.Join(parent, f.Name())
			TempURL := "/api/get_file?id=" + b.BookID + "&filename=" + url.QueryEscape(inArchiveName)
			// 替换特殊字符的时候,不要用url.PathEscape()，PathEscape不会把“+“替换成"%2b"，会导致BUG，让gin会将+解析为空格。
			b.Pages.Images = append(b.Pages.Images, model.MediaFileInfo{Path: "", Size: f.Size(), ModTime: f.ModTime(), Name: inArchiveName, Url: TempURL})
		} else {
			if option.Cfg.GetDebug() {
				logger.Infof(locale.GetString("unsupported_file_type")+" %s", name)
			}
		}
	}
	b.SortPages("default")
	return err
}
