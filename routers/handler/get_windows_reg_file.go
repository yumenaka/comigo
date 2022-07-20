package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yumenaka/comi/locale"
)

// GetRegFIleHandler 下载服务器配置
func GetRegFIleHandler(c *gin.Context) {
	// 带后缀的执行文件名 comi.exe  sketch.exe
	exeFilePath := os.Args[0]
	newStr := strings.Replace(exeFilePath, `\`, `\\`, -1)
	fmt.Println("exe_file_path:", exeFilePath)
	fmt.Println("newStr:", newStr)
	var regText = `Windows Registry Editor Version 5.00

[HKEY_CLASSES_ROOT\Directory\Background\shell\ComiGo]
"Icon"="C:\\Users\\%USERNAME%\\Desktop\\comi.exe,0"
@="ComiGo Here"

[HKEY_CLASSES_ROOT\Directory\Background\shell\ComiGo\command]
@="ComigoExePath  \"%V\""`

	//替换Icon那一行
	regText = strings.Replace(regText, `C:\\Users\\%USERNAME%\\Desktop\\comi.exe`, newStr, 1)
	//替换 ComigoExePath
	regText = strings.Replace(regText, "ComigoExePath", newStr, 1)
	//命令行打印
	fmt.Println(regText)
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
	//	fmt.Printf("open file error=%v\n", err)
	//	return
	//}
	//defer func(file *os.File) {
	//	err := file.Close()
	//	if err != nil {
	//		fmt.Printf("Close file error=%v\n", err)
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
