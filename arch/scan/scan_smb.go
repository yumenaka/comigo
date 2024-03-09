package scan

import (
	"fmt"
	"github.com/hirochachacha/go-smb2"
	"github.com/yumenaka/comi/types"
	iofs "io/fs"
	"net"
	"strconv"
)

// Smb 扫描smb书籍
func Smb(scanOption Option) (newBookList []*types.Book, err error) {
	conn, err := net.Dial("tcp", scanOption.RemoteStores[0].Host+":"+strconv.Itoa(scanOption.RemoteStores[0].Port))
	if err != nil {
		fmt.Println(err)
		panic(err)
		//return nil, err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     scanOption.RemoteStores[0].Username,
			Password: scanOption.RemoteStores[0].Password,
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		panic(err)
	}
	defer func(s *smb2.Session) {
		err := s.Logoff()
		if err != nil {
			fmt.Println(err)
		}
	}(s)

	fs, err := s.Mount(scanOption.RemoteStores[0].ShareName)
	if err != nil {
		panic(err)
	}
	defer func(fs *smb2.Share) {
		err := fs.Umount()
		if err != nil {
			fmt.Println(err)
		}
	}(fs)

	matches, err := iofs.Glob(fs.DirFS("./cvgo"), "*")
	if err != nil {
		panic(err)
	}
	for _, match := range matches {
		fmt.Println(match)
	}

	err = iofs.WalkDir(fs.DirFS("."), ".", func(path string, d iofs.DirEntry, err error) error {
		fmt.Println(path, d, err)
		return nil
	})
	if err != nil {
		panic(err)
	}

	//// 路径不存在的Error，不过目前并不会打印出来
	//if !util.PathExists(storePath) {
	//	return nil, errors.New(locale.GetString("PATH_NOT_EXIST"))
	//}
	//storePathAbs, err := filepath.Abs(storePath)
	//if err != nil {
	//	storePathAbs = storePath
	//	logger.Infof("%s", err)
	//}
	//logger.Infof(locale.GetString("SCAN_START_HINT")+"%s", storePathAbs)
	//err = filepath.Walk(storePathAbs, func(walkPath string, fileInfo os.FileInfo, err error) error {
	//	if !scanOption.ReScanFile {
	//		for _, p := range types.GetArchiveBooks() {
	//			AbsW, err := filepath.Abs(walkPath) // 取得绝对路径
	//			if err != nil {
	//				// 无法取得的情况下，用相对路径
	//				AbsW = walkPath
	//				logger.Info(err, AbsW)
	//			}
	//			if walkPath == p.FilePath || AbsW == p.FilePath {
	//				//跳过已经在数据库里面的文件
	//				logger.Infof(locale.GetString("FoundInDatabase")+"%s", walkPath)
	//				return nil
	//			}
	//		}
	//	}
	//	// 路径深度
	//	depth := strings.Count(walkPath, "/") - strings.Count(storePathAbs, "/")
	//	if runtime.GOOS == "windows" {
	//		depth = strings.Count(walkPath, "\\") - strings.Count(storePathAbs, "\\")
	//	}
	//	if depth > scanOption.MaxScanDepth {
	//		logger.Infof(locale.GetString("ExceedsMaximumDepth")+" %d，base：%s scan: %s:", scanOption.MaxScanDepth, storePathAbs, walkPath)
	//		return filepath.SkipDir // 当WalkFunc的返回值是filepath.SkipDir时，Walk将会跳过这个目录，照常执行下一个文件。
	//	}
	//	if scanOption.IsSkipDir(walkPath) {
	//		logger.Infof(locale.GetString("SkipPath")+"%s", walkPath)
	//		return filepath.SkipDir
	//	}
	//	if fileInfo == nil {
	//		return err
	//	}
	//	// 如果不是文件夹
	//	if !fileInfo.IsDir() {
	//		if !scanOption.IsSupportArchiver(walkPath) {
	//			return nil
	//		}
	//		// 得到书籍文件数据
	//		getBook, err := scanFileGetBook(walkPath, storePathAbs, depth, scanOption)
	//		if err != nil {
	//			logger.Infof("%s", err)
	//			return nil
	//		}
	//		newBookList = append(newBookList, getBook)
	//	}
	//	// 如果是文件夹
	//	if fileInfo.IsDir() {
	//		// 得到书籍文件数据
	//		getBook, err := scanDirGetBook(walkPath, storePathAbs, depth, scanOption)
	//		if err != nil {
	//			logger.Infof("%s", err)
	//			return nil
	//		}
	//		newBookList = append(newBookList, getBook)
	//	}
	//	return nil
	//})
	//// 所有可用书籍，包括压缩包与文件夹
	//if len(newBookList) > 0 {
	//	logger.Infof(locale.GetString("FOUND_IN_PATH"), len(newBookList), storePathAbs)
	//}
	//return newBookList, err
	return newBookList, err
}
