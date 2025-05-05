package scan

import (
	"context"
	"errors"
	"github.com/klauspost/compress/zip"
	"github.com/yumenaka/archives"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/file"
	"github.com/yumenaka/comigo/util/logger"
	"io/fs"
	"net/url"
	"path"
	"time"
)

// 处理 ZIP 和 EPUB 文件
func handleZipAndEpubFiles(filePath string, newBook *model.Book, option Option) error {
	fsys, err := zip.OpenReader(filePath)
	if err != nil {
		return errors.New(locale.GetString("not_a_valid_zip_file") + filePath)
	}
	defer fsys.Close()

	err = walkUTF8ZipFs(fsys, "", ".", newBook, option)
	if err != nil {
		var pathError *fs.PathError
		if errors.As(err, &pathError) {
			if option.Cfg.GetDebug() {
				logger.Infof("NonUTF-8 ZIP: %s, Error: %s", filePath, err.Error())
			}
			err = scanNonUTF8ZipFile(filePath, newBook, option)
		} else {
			return err
		}
	}

	if newBook.Type == model.TypeEpub {
		imageList, err := file.GetImageListFromEpubFile(newBook.FilePath)
		if err == nil {
			newBook.SortPagesByImageList(imageList)
		} else {
			logger.Infof("Failed to get image list from EPUB: %s, error: %v", newBook.FilePath, err)
		}

		metaData, err := file.GetEpubMetadata(newBook.FilePath)
		if err == nil {
			newBook.Author = metaData.Creator
			newBook.Press = metaData.Publisher
		} else {
			logger.Infof("Failed to get metadata from EPUB: %s, error: %v", newBook.FilePath, err)
		}
	}
	return nil
}

// 处理其他类型的压缩文件
func handleOtherArchiveFiles(filePath string, newBook *model.Book, option Option) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fsys, err := archives.FileSystem(ctx, filePath, nil)
	if err != nil {
		return err
	}

	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			logger.Infof("Failed to access path %s in archive: %v", path, err)
			return err
		}

		if option.IsSkipDir(path) {
			logger.Infof("Skip Scan: %s", path)
			return fs.SkipDir
		}

		f, err := d.Info()
		if err != nil {
			logger.Infof("Failed to get file info in archive: %v", err)
			return fs.SkipDir
		}

		if option.IsSupportMedia(path) {
			archivedFile, ok := f.(archives.FileInfo)
			var tempURL string
			if ok {
				tempURL = "/api/get_file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(archivedFile.NameInArchive)
				newBook.Pages.Images = append(newBook.Pages.Images, model.MediaFileInfo{
					Name: archivedFile.NameInArchive,
					Url:  tempURL,
				})
			} else {
				tempURL = "/api/get_file?id=" + newBook.BookID + "&filename=" + url.QueryEscape(path)
				newBook.Pages.Images = append(newBook.Pages.Images, model.MediaFileInfo{
					Url: tempURL,
				})
			}
		} else {
			if option.Cfg.GetDebug() {
				logger.Infof(locale.GetString("unsupported_file_type")+" %s", path)
			}
		}
		return nil
	})
	return err
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
