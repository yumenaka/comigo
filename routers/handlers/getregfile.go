package handlers

import (
	"github.com/yumenaka/comi/logger"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yumenaka/comi/locale"
)

// HandlerGetRegFile 下载服务器配置
func HandlerGetRegFile(c *gin.Context) {
	if runtime.GOOS != "windows" {
		logger.Info("Now system not windows,can't generate reg file.")
		return
	}
	// Windows特定文件添加右键菜单
	// 参考资料：https://blog.csdn.net/yang382197207/article/details/80079052

	// 带后缀的执行文件名 comi.exe  sketch.exe
	exeFilePath := os.Args[0]
	// 在创建字符串类型的键值时，如果该字符串中包含路径分隔符，这个路径分隔符要用双斜杠“\\"表示。
	newStr := strings.Replace(exeFilePath, `\`, `\\`, -1)
	logger.Info("exe_file_path:", exeFilePath)
	logger.Info("newStr:", newStr)
	//HKEY_CLASSES_ROOT\*：系统所有文件，右键系统任一文件都会添加右键菜单
	//1% 用来传递文件名，一定要加引号，不然当文件名含有空格时，只能得到空格前的部分。
	//最后要以一个换行结束，为了保证汉字正常，最好用ANSI编码
	var regText = `Windows Registry Editor Version 5.00

[HKEY_CLASSES_ROOT\AllFilesystemObjects\shell\ComiGo]
@="ComiGo"
"Icon"="C:\\Users\\%USERNAME%\\Desktop\\comi.exe,0"

[HKEY_CLASSES_ROOT\AllFilesystemObjects\shell\ComiGo\command]
@="\"C:\\Users\\%USERNAME%\\Desktop\\comi.exe\"  \"%1\""

;HKEY_CLASSES_ROOT\Directory\Background：文件夹空白处右键的菜单 
[HKEY_CLASSES_ROOT\Directory\Background\shell\ComiGo]
"Icon"="C:\\Users\\%USERNAME%\\Desktop\\comi.exe,0"
@="ComiGo Here"

[HKEY_CLASSES_ROOT\Directory\Background\shell\ComiGo\command]
@="\"ComigoExePath\"  \"%V\""

`

	//替换Icon那一行
	regText = strings.Replace(regText, `C:\\Users\\%USERNAME%\\Desktop\\comi.exe`, newStr, -1)
	//替换 ComigoExePath
	regText = strings.Replace(regText, "ComigoExePath", newStr, -1)
	//命令行打印 如果值中有中文，则需要将.reg文件以ascii编码保存，否则会出现乱码。
	logger.Info(regText)
	regFileName := strings.Replace("comigo(XXX).reg", "XXX", locale.GetString("REG_FILE_HINT"), 1)
	//用gin实现下载文件的功能，只需要在接口返回时设置Response-Header中的Content-Type为文件类型，并设置Content-Disposition指定默认的文件名，然后将文件数据返回浏览器即可
	fileContentDisposition := "attachment;filename=\"" + regFileName + "\""
	c.Header("Content-Type", "application/octet-stream") // 这里是压缩文件类型 .zip
	c.Header("Content-Disposition", fileContentDisposition)
	c.Data(http.StatusOK, "application/octet-stream", []byte(regText))
	//在程序执行目录创建一个REG文件，并写入内容
	//targetPath, _ := os.Getwd()
	//filePath := path.Join(targetPath, regFileName)
	//file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	//if err != nil {
	//	logger.Infof("open file error=%v\n", err)
	//	return
	//}
	//defer func(file *os.File) {
	//	err := file.Close()
	//	if err != nil {
	//		logger.Infof("Close file error=%v\n", err)
	//	}
	//}(file)
	//writer := bufio.NewWriter(file)
	//_, err = writer.WriteString(regText)
	//if err != nil {
	//	return
	//}
	//err := writer.Flush()
	//if err != nil {
	//	return
	//}
}
