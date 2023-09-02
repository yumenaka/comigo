package cmd

import (
	"fmt"
	"github.com/yumenaka/comi/common"
	"os"
	"path"
)

// AddPathToStore 添加默认扫描路径
func AddPathToStore(args []string) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get working directory:", err)
		return
	}
	fmt.Println("Working directory:", wd)
	//没有指定路径或文件的情况下
	if len(args) == 0 {
		common.Config.StoresPath = append(common.Config.StoresPath, wd)
	} else {
		//指定了多个参数的话，都扫描一遍
		for _, arg := range args {
			common.Config.StoresPath = append(common.Config.StoresPath, arg)
		}
	}

	//启用上传，则添加upload目录
	if common.Config.EnableUpload {
		if common.Config.UploadPath != "" {
			common.Config.StoresPath = append(common.Config.StoresPath, common.Config.UploadPath)
		} else {
			common.Config.StoresPath = append(common.Config.StoresPath, path.Join(wd, "upload"))
		}
	}
}
