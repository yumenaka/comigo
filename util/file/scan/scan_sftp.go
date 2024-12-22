package scan

import (
	"fmt"
	"github.com/hirochachacha/go-smb2"
	"github.com/yumenaka/comigo/model"
	iofs "io/fs"
	"net"
	"strconv"
)

// https://github.com/moov-io/go-sftp  ?
// https://github.com/spf13/afero  SftpFs

// TODO:SFTP扫描书籍
func SFTP(scanOption Option) (newBookList []*model.Book, err error) {
	conn, err := net.Dial("tcp", scanOption.RemoteStores[0].Smb.Host+":"+strconv.Itoa(scanOption.RemoteStores[0].Smb.Port))
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
			User:     scanOption.RemoteStores[0].Smb.Username,
			Password: scanOption.RemoteStores[0].Smb.Password,
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

	fs, err := s.Mount(scanOption.RemoteStores[0].Smb.ShareName)
	if err != nil {
		panic(err)
	}
	defer func(fs *smb2.Share) {
		err := fs.Umount()
		if err != nil {
			fmt.Println(err)
		}
	}(fs)
	////fs.DirFS(".") 创建一个表示当前目录(".")的文件系统。"." 表示当前工作目录
	//// iofs.Glob(fs.DirFS("."), "*") 调用Glob函数在当前目录下搜索匹配给定模式的文件名。"*" 模式意味着匹配所有文件和目录。
	//matches, err := iofs.Glob(fs.DirFS("test"), "*")
	//if err != nil {
	//	panic(err)
	//}
	//for _, match := range matches {
	//	fmt.Println(match)
	//}

	// iofs.WalkDir(fs.DirFS("."), ".", func...) 调用WalkDir函数遍历当前目录（以及其下的所有子目录）中的所有文件和目录。
	err = iofs.WalkDir(
		//fs.DirFS(".") 指定特定目录作为遍历的起点。"." 表示当前工作目录。
		fs.DirFS("test"),
		".",
		//对于目录中的每一个项（无论是文件还是目录），指定的函数都会被调用。
		func(path string, d iofs.DirEntry, err error) error {
			//这个函数接收三个参数：path（项的路径），d（一个DirEntry对象，表示文件或目录的信息），和err（如果在访问该项时出现错误）
			fmt.Println(path, d, err)
			return nil
		})
	if err != nil {
		panic(err)
	}

	//// 路径不存在的Error，不过目前并不会打印出来
	//if !util.PathExists(storePath) {
	//	return nil, errors.New(locale.GetString("path_not_exist"))
	//}
	//storePathAbs, err := filepath.Abs(storePath)
	//if err != nil {
	//	storePathAbs = storePath
	//	logger.Infof("%s", err)
	//}
	//logger.Infof(locale.GetString("scan_start_hint")+"%s", storePathAbs)
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
	//				logger.Infof(locale.GetString("found_in_bookstore")+"%s", walkPath)
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
	//		logger.Infof(locale.GetString("exceeds_maximum_depth")+" %d，base：%s scan: %s:", scanOption.MaxScanDepth, storePathAbs, walkPath)
	//		return filepath.SkipDir // 当WalkFunc的返回值是filepath.SkipDir时，Walk将会跳过这个目录，照常执行下一个文件。
	//	}
	//	if scanOption.IsSkipDir(walkPath) {
	//		logger.Infof(locale.GetString("skip_path")+"%s", walkPath)
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
	//	logger.Infof(locale.GetString("found_in_path"), len(newBookList), storePathAbs)
	//}
	//return newBookList, err
	return newBookList, err
}
