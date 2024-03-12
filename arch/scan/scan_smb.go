package scan

import (
	"fmt"
	"github.com/hirochachacha/go-smb2"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/types"
	iofs "io/fs"
	"net"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// Smb 扫描smb书籍
func Smb(scanOption Option) (newBookList []*types.Book, err error) {
	// connection
	connection, err := net.Dial("tcp", scanOption.RemoteStores[0].Host+":"+strconv.Itoa(scanOption.RemoteStores[0].Port))
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
	}(connection)

	dialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     scanOption.RemoteStores[0].Username,
			Password: scanOption.RemoteStores[0].Password,
		},
	}

	session, err := dialer.Dial(connection)
	if err != nil {
		panic(err)
	}
	defer func(s *smb2.Session) {
		err := s.Logoff()
		if err != nil {
			fmt.Println(err)
		}
	}(session)

	fs, err := session.Mount(scanOption.RemoteStores[0].ShareName)
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
			//这个函数接收三个参数：path（项的路径），dialer（一个DirEntry对象，表示文件或目录的信息），和err（如果在访问该项时出现错误）
			fmt.Println("smb://"+scanOption.RemoteStores[0].Host+"/"+scanOption.RemoteStores[0].ShareName+"/test/"+path, d.IsDir(), err)
			return nil
		})
	if err != nil {
		panic(err)
	}
	err = iofs.WalkDir(
		//fs.DirFS(".") 指定特定目录作为遍历的起点。"." 表示当前工作目录。
		fs.DirFS("test"),
		".",
		//对于目录中的每一个项（无论是文件还是目录），指定的函数都会被调用。
		func(walkPath string, fileInfo iofs.DirEntry, err error) error {
			smbFilePath := "smb://" + scanOption.RemoteStores[0].Host + "/" + scanOption.RemoteStores[0].ShareName + "/" + walkPath
			if !scanOption.ReScanFile {
				for _, p := range types.GetArchiveBooks() {
					if smbFilePath == p.FilePath {
						//跳过已经在数据库里面的文件
						logger.Infof(locale.GetString("FoundInDatabase")+"%path", walkPath)
						return nil
					}
				}
			}
			// TODO：SMB路径深度。这里的深度是指相对于扫描的根目录的深度。
			depth := strings.Count(walkPath, "/")
			if runtime.GOOS == "windows" {
				depth = strings.Count(walkPath, "\\")
			}
			if depth > scanOption.MaxScanDepth {
				logger.Infof(locale.GetString("ExceedsMaximumDepth")+" %dialer，base：%session scan: %session:", scanOption.MaxScanDepth, scanOption.RemoteStores[0].ShareName, walkPath)
				return filepath.SkipDir // 当WalkFunc的返回值是filepath.SkipDir时，Walk将会跳过这个目录，照常执行下一个文件。
			}
			if scanOption.IsSkipDir(walkPath) {
				logger.Infof(locale.GetString("SkipPath")+"%p", walkPath)
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
				getBook, err := scanFileGetBook(walkPath, scanOption.RemoteStores[0].ShareName, depth, scanOption)
				if err != nil {
					logger.Infof("%session", err)
					return nil
				}
				newBookList = append(newBookList, getBook)
			}
			// 如果是文件夹
			if fileInfo.IsDir() {
				// 得到书籍文件数据
				getBook, err := scanDirGetBook(walkPath, scanOption.RemoteStores[0].ShareName, depth, scanOption)
				if err != nil {
					logger.Infof("%e", err)
					return nil
				}
				newBookList = append(newBookList, getBook)
			}
			return nil
		})
	// 所有可用书籍，包括压缩包与文件夹
	if len(newBookList) > 0 {
		logger.Infof(locale.GetString("FOUND_IN_PATH"), len(newBookList), scanOption.RemoteStores[0].ShareName)
	}
	return newBookList, err
}
