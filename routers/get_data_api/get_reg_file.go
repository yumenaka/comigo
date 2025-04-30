package get_data_api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/util/logger"
)

// TODO：GetRegFile 设置注册表文件
func GetRegFile(c echo.Context) error {
	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		logger.Infof("%s", err)
		return c.String(http.StatusInternalServerError, "Error getting executable path")
	}

	// 获取当前可执行文件的目录
	exeDir := filepath.Dir(exePath)
	// 获取当前可执行文件的名称（不含扩展名）
	exeName := strings.TrimSuffix(filepath.Base(exePath), filepath.Ext(exePath))

	// 构建注册表文件内容
	regContent := fmt.Sprintf(`Windows Registry Editor Version 5.00

[HKEY_CLASSES_ROOT\comigo]
@="URL:comigo Protocol"
"URL Protocol"=""

[HKEY_CLASSES_ROOT\comigo\DefaultIcon]
@="\"%s\""

[HKEY_CLASSES_ROOT\comigo\shell]

[HKEY_CLASSES_ROOT\comigo\shell\open]

[HKEY_CLASSES_ROOT\comigo\shell\open\command]
@="\"%s\" \"%%1\""`, exePath, exePath)

	// 构建注册表文件路径
	regFilePath := filepath.Join(exeDir, exeName+".reg")

	// 写入注册表文件
	err = os.WriteFile(regFilePath, []byte(regContent), 0o644)
	if err != nil {
		logger.Infof("%s", err)
		return c.String(http.StatusInternalServerError, "Error writing reg file")
	}

	// 设置响应头
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.reg", exeName))
	c.Response().Header().Set("Content-Type", "application/x-windows-registry-script")

	// 返回文件
	return c.File(regFilePath)
}
