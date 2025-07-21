package scan

import (
	"context"
	"errors"
	iofs "io/fs"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/cloudsoda/go-smb2"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

// Smb 扫描smb书籍  github.com/hirochachacha/go-smb2
// 换用一个持续更新，rclone用的库：github.com/cloudsoda/go-smb2
// https://github.com/rclone/rclone/blob/master/go.mod
func Smb(scanOption Option) (newBookList []*model.Book, err error) {
	// connection
	dialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     scanOption.Cfg.GetStores()[0].Smb.Username,
			Password: scanOption.Cfg.GetStores()[0].Smb.Password,
		},
	}

	session, err := dialer.Dial(context.Background(), scanOption.Cfg.GetStores()[0].Smb.Host+":"+strconv.Itoa(scanOption.Cfg.GetStores()[0].Smb.Port))
	if err != nil {
		panic(err)
	}
	defer session.Logoff()

	fs, err := session.Mount(scanOption.Cfg.GetStores()[0].Smb.ShareName)
	if err != nil {
		panic(err)
	}
	defer fs.Umount()

	// // iofs.WalkDir(fs.DirFS("."), ".", func...) 调用WalkDir函数遍历当前目录（以及其下的所有子目录）中的所有文件和目录。
	// err = iofs.WalkDir(
	//	//fs.DirFS(".") 指定特定目录作为遍历的起点。"." 表示当前工作目录。
	//	fs.DirFS("test"),
	//	".",
	//	//对于目录中的每一个项（无论是文件还是目录），指定的函数都会被调用。
	//	func(path string, d iofs.DirEntry, err error) error {
	//		//这个函数接收三个参数：path（项的路径），dialer（一个DirEntry对象，表示文件或目录的信息），和err（如果在访问该项时出现错误）
	//		logger.Info("smb hint： smb://"+scanOption.Stores[0].Host+"/"+scanOption.Stores[0].ShareName+"/test/"+path, d.IsDir(), err)
	//		return nil
	//	})
	// if err != nil {
	//	panic(err)
	// }

	err = iofs.WalkDir(
		// fs.DirFS(".") 指定特定目录作为遍历的起点。"." 表示当前工作目录。
		fs.DirFS("test"),
		".",
		// 对于目录中的每一个项（无论是文件还是目录），指定的函数都会被调用。
		func(walkPath string, fileInfo iofs.DirEntry, err error) error {
			smbFilePath := "smb://" + scanOption.Cfg.GetStores()[0].Smb.Host + "/" + scanOption.Cfg.GetStores()[0].Smb.ShareName + "/" + walkPath

			for _, p := range model.GetArchiveBooks() {
				if smbFilePath == p.FilePath {
					// 跳过已扫描文件
					logger.Infof(locale.GetString("found_in_bookstore")+"%path", walkPath)
					return nil
				}
			}

			// SMB路径深度。这里的深度是指相对于扫描的根目录的深度。
			depth := strings.Count(walkPath, "/")
			if runtime.GOOS == "windows" {
				depth = strings.Count(walkPath, "\\")
			}
			if depth > scanOption.Cfg.GetMaxScanDepth() {
				logger.Infof(
					locale.GetString("exceeds_maximum_depth")+" %dialer，base：%session scan: %session:",
					scanOption.Cfg.GetMaxScanDepth(),
					scanOption.Cfg.GetStores()[0].Smb.ShareName, walkPath)
				return filepath.SkipDir // 当WalkFunc的返回值是filepath.SkipDir时，Walk将会跳过这个目录，照常执行下一个文件。
			}
			if scanOption.IsSkipDir(walkPath) {
				logger.Infof(locale.GetString("skip_path")+"%p", walkPath)
				return filepath.SkipDir
			}
			if fileInfo == nil {
				return err
			}
			// 如果不是文件夹
			if !fileInfo.IsDir() {
				if !scanOption.IsSupportArchiver(walkPath) {
					return nil
				}
				// 打开文件
				file, err := fs.OpenFile("test/"+walkPath, os.O_RDONLY, 0o400) // Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
				if err != nil {
					logger.Infof("Smb OpenFile: %s", err.Error())
					return nil
				}

				// 得到书籍文件数据
				getBook, err := smbScanFile(walkPath, file, scanOption.Cfg.GetStores()[0].Smb.ShareName, depth, scanOption)
				if err != nil {
					logger.Infof("%session", err)
					return nil
				}
				newBookList = append(newBookList, getBook)
			}
			// // 如果是文件夹
			// if fileInfo.IsDir() {
			//	// 得到书籍文件数据
			//	getBook, err := smbScanDir(walkPath, scanOption.Stores[0].ShareName, depth, scanOption)
			//	if err != nil {
			//		logger.Infof("%e", err)
			//		return nil
			//	}
			//	newBookList = append(newBookList, getBook)
			// }
			return nil
		})
	// 所有可用书籍，包括压缩包与文件夹
	if len(newBookList) > 0 {
		logger.Infof(locale.GetString("found_in_path"), scanOption.Cfg.GetStores()[0].Smb.ShareName, len(newBookList))
		return newBookList, err
	}
	return nil, errors.New("NO_BOOKS_FOUND in SMB:" + scanOption.Cfg.GetStores()[0].Smb.ShareName)
}

// 扫描本地路径，并返回对应书籍
func smbScanFile(filePath string, file *smb2.File, storePath string, depth int, scanOption Option) (newBook *model.Book, err error) {
	// 设置了一个defer函数来捕获可能的panic
	// defer func() {
	//	if err := recover(); err != nil {
	//		logger.Info("Recovered from panic:", err)
	//		// 可以在这里执行一些处理逻辑，比如记录日志、返回错误信息等
	//	}
	// }()
	defer func(file *smb2.File) {
		err := file.Close()
		if err != nil {
			logger.Infof("%s", err)
		}
	}(file)
	logger.Info(file.Name())
	// MediaFileInfo, err := file.Stat()

	return nil, err
	// if err != nil {
	//	logger.Infof("%s", err.Error())
	//	return nil, err
	// }
	// // 初始化一本书，设置文件路径等等
	// newBook, err := types.NewBook(filePath, MediaFileInfo.ModTime(), MediaFileInfo.Size(), storePath, depth, types.GetBookTypeByFilename(filePath))
	// if err != nil {
	//	return nil, err
	// }
	// // 根据文件类型，走不同的初始化流程
	// switch newBook.Type {
	// // 为解决archiver/v4的BUG “zip文件无法读取2级目录” 单独处理zip文件
	// case types.TypeZip, types.TypeCbz, types.TypeEpub:
	//	// 使用Archiver的虚拟文件系统，无法处理非UTF-8编码
	//	fsys, zipErr := zip.OpenReader(filePath)
	//	if zipErr != nil {
	//		// logger.Infof(zipErr)
	//		return nil, errors.New(locale.GetString("not_a_valid_zip_file") + filePath)
	//	}
	//	err = walkUTF8ZipFs(fsys, "", ".", newBook, scanOption)
	//	// 如果扫描ZIP文件的时候遇到了 fs.PathError ，则扫描到NonUTF-8 ZIP文件，需要特殊处理
	//	var pathError *iofs.PathError
	//	if errors.As(err, &pathError) {
	//		if scanOption.Debug {
	//			logger.Infof("NonUTF-8 ZIP:%s  Error:%s", filePath, err.Error())
	//		}
	//		// 忽略 fs.PathError 并换个方式扫描
	//		err = scanNonUTF8ZipFile(filePath, newBook, scanOption)
	//	}
	//	// epub文件，需要根据 META-INF/container.xml 里面定义的rootfile （.opf文件）来重新排序
	//	if newBook.Type == types.TypeEpub {
	//		imageList, err := arch.GetImageListFromEpubFile(newBook.FilePath)
	//		if err != nil {
	//			logger.Infof("%s", err)
	//		} else {
	//			newBook.SortPagesByImageList(imageList)
	//		}
	//		// 根据metadata，改写书籍信息
	//		metaData, err := arch.GetEpubMetadata(newBook.FilePath)
	//		if err != nil {
	//			logger.Infof("%s", err)
	//		} else {
	//			newBook.Author = metaData.Creator
	//			newBook.Press = metaData.Publisher
	//		}
	//	}
	// case types.TypePDF:
	//	pageCount, pdfErr := arch.CountPagesOfPDF(filePath)
	//	if pdfErr != nil {
	//		return nil, pdfErr
	//	}
	//	if pageCount < 1 {
	//		return nil, errors.New(locale.GetString("no_pages_in_pdf") + filePath)
	//	}
	//	logger.Infof(locale.GetString("scan_pdf")+"%s: %d", filePath, pageCount)
	//	newBook.PageCount = pageCount
	//	newBook.InitComplete = true
	//	newBook.Cover = types.MediaFileInfo{RealImageFilePATH: "", Size: MediaFileInfo.Size(), ModTime: MediaFileInfo.ModTime(), Name: "", Url: "/images/pdf.png"}
	//	for i := 1; i <= pageCount; i++ {
	//		TempURL := "/api/get_file?id=" + newBook.BookID + "&filename=" + strconv.Itoa(i) + ".jpg"
	//		newBook.Pages.Images = append(newBook.Pages.Images, types.MediaFileInfo{RealImageFilePATH: "", Size: MediaFileInfo.Size(), ModTime: MediaFileInfo.ModTime(), Name: strconv.Itoa(i), Url: TempURL})
	//	}
	// case types.TypeVideo:
	//	newBook.PageCount = 1
	//	newBook.InitComplete = true
	//	newBook.Cover = types.MediaFileInfo{Name: "video.png", Url: "/images/video.png"}
	// case types.TypeAudio:
	//	newBook.PageCount = 1
	//	newBook.InitComplete = true
	//	newBook.Cover = types.MediaFileInfo{Name: "audio.png", Url: "/images/audio.png"}
	// case types.TypeUnknownFile:
	//	newBook.PageCount = 1
	//	newBook.InitComplete = true
	//	newBook.Cover = types.MediaFileInfo{Name: "unknown.png", Url: "/images/unknown.png"}
	// // 其他类型的压缩文件或文件夹
	// default:
	//	// archiver.FileSystem可以配合ctx了，加个默认超时时间
	//	const shortDuration = 10 * 1000 * time.Millisecond // 超时时间，10秒
	//	ctx, cancel := context.WithTimeout(context.Background(), shortDuration)
	//	defer cancel()
	//	fsys, err := archiver.FileSystem(ctx, filePath)
	//	if err != nil {
	//		return nil, err
	//	}
	//	err = iofs.WalkDir(fsys, ".", func(path string, d iofs.DirEntry, err error) error {
	//		if scanOption.IsSkipDir(path) {
	//			logger.Infof("Skip Scan:", path)
	//			return iofs.SkipDir
	//		}
	//		f, errInfo := d.Info()
	//		if errInfo != nil {
	//			logger.Info(errInfo)
	//			return iofs.SkipDir
	//		}
	//		if !scanOption.IsSupportMedia(path) {
	//			if scanOption.Debug {
	//				logger.Infof(locale.GetString("unsupported_file_type")+"%s", path)
	//			}
	//		} else {
	//			u, ok := f.(archiver.File) // f.Name不包含路径信息.需要转换一下
	//			if !ok {
	//				// 如果是文件夹+图片
	//				newBook.Type = types.TypeDir
	//				////用Archiver的虚拟文件系统提供图片文件，理论上现在不应该用到
	//				//newBook.Pages = append(newBook.Pages, MediaFileInfo{RealImageFilePATH: "", Size: f.Size(), ModTime: f.ModTime(), Name: "", Url: "/cache/" + newBook.BookID + "/" + url.QueryEscape(path)})
	//				//实验：用get_file接口提供文件服务
	//				TempURL := "/api/get_file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(path)
	//				newBook.Pages.Images = append(newBook.Pages.Images, types.MediaFileInfo{RealImageFilePATH: "", Size: f.Size(), ModTime: f.ModTime(), Name: "", Url: TempURL})
	//				// logger.Infof(locale.GetString("unsupported_extract")+" %s", f)
	//			} else {
	//				// 替换特殊字符的时候，额外将“+替换成"%2b"，因为gin会将+解析为空格。
	//				TempURL := "/api/get_file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(u.Name)
	//				// 不替换特殊字符
	//				// TempURL := "/api/get_file?id=" + newBook.BookID + "&filename=" + u.Name
	//				newBook.Pages.Images = append(newBook.Pages.Images, types.MediaFileInfo{RealImageFilePATH: "", Size: f.Size(), ModTime: f.ModTime(), Name: u.Name, Url: TempURL})
	//			}
	//		}
	//		return nil
	//	})
	//	if err != nil {
	//		return nil, err
	//	}
	// }
	// // 不管页数，直接返回：在添加到书库时判断页数
	// newBook.SortPages("default")
	// return newBook, err
}

func smbScanDir(dirPath string, storePath string, depth int, scanOption Option) (*model.Book, error) {
	// 初始化，生成UUID
	newBook, err := model.NewBook(dirPath, time.Now(), 0, storePath, depth, model.TypeDir)
	if err != nil {
		return nil, err
	}
	// // 获取目录中的文件和子目录的详细信息
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	infos := make([]iofs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	for _, file := range infos {
		// 跳过子目录, 只搜寻目录中的文件
		if file.IsDir() {
			continue
		}
		// 输出绝对路径
		strAbsPath, errPath := filepath.Abs(dirPath + "/" + file.Name())
		if errPath != nil {
			logger.Info(errPath)
		}
		if scanOption.IsSupportMedia(file.Name()) {
			TempURL := "/api/get_file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(file.Name())
			newBook.Pages.Images = append(newBook.Pages.Images, model.MediaFileInfo{Path: strAbsPath, Size: file.Size(), ModTime: file.ModTime(), Name: file.Name(), Url: TempURL})
		}
	}
	newBook.SortPages("default")
	// 在添加到书库时判断页数
	return newBook, err
}
