package scan

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/types"
	"github.com/yumenaka/comi/util"
)

// Local 扫描路径，取得路径里的书籍
func Local(storePath string, scanOption Option) (newBookList []*types.Book, err error) {
	// 路径不存在的Error，不过目前并不会打印出来
	if !util.PathExists(storePath) {
		return nil, errors.New(locale.GetString("PATH_NOT_EXIST"))
	}
	storePathAbs, err := filepath.Abs(storePath)
	if err != nil {
		storePathAbs = storePath
		logger.Infof("%s", err)
	}
	logger.Infof(locale.GetString("SCAN_START_HINT")+"%s", storePathAbs)
	err = filepath.Walk(storePathAbs, func(walkPath string, fileInfo os.FileInfo, err error) error {
		if !scanOption.ReScanFile {
			for _, p := range types.GetArchiveBooks() {
				AbsW, err := filepath.Abs(walkPath) // 取得绝对路径
				if err != nil {
					// 无法取得的情况下，用相对路径
					AbsW = walkPath
					logger.Info(err, AbsW)
				}
				if walkPath == p.FilePath || AbsW == p.FilePath {
					//跳过已经在数据库里面的文件
					logger.Infof(locale.GetString("FoundInDatabase")+"%s", walkPath)
					return nil
				}
			}
		}
		// 路径深度
		depth := strings.Count(walkPath, "/") - strings.Count(storePathAbs, "/")
		if runtime.GOOS == "windows" {
			depth = strings.Count(walkPath, "\\") - strings.Count(storePathAbs, "\\")
		}
		if depth > scanOption.MaxScanDepth {
			logger.Infof(locale.GetString("ExceedsMaximumDepth")+" %d，base：%s scan: %s:", scanOption.MaxScanDepth, storePathAbs, walkPath)
			return filepath.SkipDir // 当WalkFunc的返回值是filepath.SkipDir时，Walk将会跳过这个目录，照常执行下一个文件。
		}
		if scanOption.IsSkipDir(walkPath) {
			logger.Infof(locale.GetString("SkipPath")+"%s", walkPath)
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
			// 得到书籍文件数据
			getBook, err := scanFileGetBook(walkPath, storePathAbs, depth, scanOption)
			if err != nil {
				logger.Infof("%s", err)
				return nil
			}
			newBookList = append(newBookList, getBook)
		}
		// 如果是文件夹
		if fileInfo.IsDir() {
			// 得到书籍文件数据
			getBook, err := scanDirGetBook(walkPath, storePathAbs, depth, scanOption)
			if err != nil {
				logger.Infof("%s", err)
				return nil
			}
			newBookList = append(newBookList, getBook)
		}
		return nil
	})
	// 所有可用书籍，包括压缩包与文件夹
	if len(newBookList) > 0 {
		logger.Infof(locale.GetString("FOUND_IN_PATH"), len(newBookList), storePathAbs)
	}
	return newBookList, err
}
