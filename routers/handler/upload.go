package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/tools"
	"log"
	"net/http"
	"os"
	"path"
)

var LocalRescanBroadcast *chan string

// UploadHandler 下载服务器配置
// 除了设置头像以外，也可以做上传文件并阅读功能
// Set a lower memory limit for multipart forms (default is 32 MiB)
// https://github.com/gin-gonic/examples/blob/master/upload-file/single/main.go
// 也能上传多个文件，示例：
//https://github.com/gin-gonic/examples/blob/master/upload-file/multiple/main.go
//engine.MaxMultipartMemory = 60 << 20  // 60 MiB  只限制程序在上传文件时可以使用多少内存，而是不限制上传文件的大小。(default is 32 MiB)
func UploadHandler(c *gin.Context) {
	// single file
	file, err := c.FormFile("file")
	if err != nil { //没有传文件会报错，处理这个错误。
		fmt.Println(err)
	}
	UploadDir := "ComigoUpload"
	if !tools.ChickExists(UploadDir) {
		// 创建文件夹
		err := os.MkdirAll(UploadDir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir success!\n")
		}
	}
	// Upload the file to specific dst.
	err = c.SaveUploadedFile(file, path.Join(UploadDir, file.Filename))
	if err != nil {
		fmt.Println(err)
	}
	/*
	   也可以直接使用io操作，拷贝文件数据。
	   out, err := os.Create(filename)
	   defer out.Close()
	   _, err = io.Copy(out, file)
	*/
	log.Println(file.Filename)
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	*LocalRescanBroadcast <- UploadDir
}
