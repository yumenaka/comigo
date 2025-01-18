package settings_page

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/angelofallars/htmx-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/htmx/state"
	"github.com/yumenaka/comigo/util"
	"github.com/yumenaka/comigo/util/file/scan"
	"github.com/yumenaka/comigo/util/logger"
)

// -------------------------
// 使用templ模板响应htmx请求
// -------------------------

func Tab1(c *gin.Context) {
	template := tab1(&state.Global) // define body content
	// 用模板渲染 html 元素
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, template); renderErr != nil {
		// 如果出错，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func Tab2(c *gin.Context) {
	template := tab2(&state.Global) // define body content
	// 用模板渲染 html 元素
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, template); renderErr != nil {
		// 如果出错，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func Tab3(c *gin.Context) {
	template := tab3(&state.Global) // define body content
	// 用模板渲染 html 元素
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, template); renderErr != nil {
		// 如果出错，返回 HTTP 500 错误。
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

// -------------------------
// 抽取更新配置的公共逻辑
// -------------------------

// parseSingleHTMXFormPair 提取并返回表单中的“第一对” key/value。
// 如果不是 HTMX 请求或解析失败等情况，返回对应的错误。
func parseSingleHTMXFormPair(c *gin.Context) (string, string, error) {
	if !htmx.IsHTMX(c.Request) {
		return "", "", errors.New("non-htmx request")
	}
	if err := c.Request.ParseForm(); err != nil {
		return "", "", fmt.Errorf("parseForm error: %v", err)
	}
	formData := c.Request.PostForm
	if len(formData) == 0 {
		return "", "", errors.New("no form data")
	}
	var name, newValue string
	for key, values := range formData {
		name = key
		if len(values) > 0 {
			newValue = values[0]
		}
		break
	}
	return name, newValue, nil
}

// updateConfigGeneric 抽取了更新配置的大部分公共逻辑，返回 name 和 newValue 等。
func updateConfigGeneric(c *gin.Context) (string, string, error) {
	name, newValue, err := parseSingleHTMXFormPair(c)
	if err != nil {
		return "", "", err
	}

	logger.Infof("Update config: %s = %s", name, newValue)

	// 旧配置做个备份（有需要对比）
	oldConfig := config.CopyCfg()

	// 写入配置文件
	if writeErr := config.WriteConfigFile(); writeErr != nil {
		logger.Infof("Failed to update local config: %v", writeErr)
	}

	// 更新配置
	if setErr := state.ServerConfig.SetConfigValue(name, newValue); setErr != nil {
		logger.Errorf("Failed to set config value: %v", setErr)
		return "", "", setErr
	}

	// 根据配置的变化，做相应操作。比如打开浏览器,重新扫描等
	beforeConfigUpdate(&oldConfig, config.GetCfg())

	return name, newValue, nil
}

// -------------------------
// 各类配置的更新 Handler
// -------------------------

// UpdateStringConfigHandler 处理 String 类型
func UpdateStringConfigHandler(c *gin.Context) {
	name, newValue, err := updateConfigGeneric(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// 渲染对应的模板
	updatedHTML := StringConfig(name, newValue, name+"_Description")
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, updatedHTML); renderErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

// UpdateBoolConfigHandler 处理 Bool 类型
func UpdateBoolConfigHandler(c *gin.Context) {
	name, newValue, err := updateConfigGeneric(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 将字符串解析为 bool
	boolVal, parseErr := strconv.ParseBool(newValue)
	if parseErr != nil {
		logger.Errorf("无法将 '%s' 解析为 bool: %v", newValue, parseErr)
		c.String(http.StatusBadRequest, "parse bool error")
		return
	}
	// 渲染对应的模板
	updatedHTML := BoolConfig(name, boolVal, name+"_Description")
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, updatedHTML); renderErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

// UpdateNumberConfigHandler 处理 Number 类型
func UpdateNumberConfigHandler(c *gin.Context) {
	name, newValue, err := updateConfigGeneric(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// 将字符串解析为 int
	intVal, parseErr := strconv.ParseInt(newValue, 10, 64)
	if parseErr != nil {
		logger.Errorf("无法将 '%s' 解析为 int: %v", newValue, parseErr)
		c.String(http.StatusBadRequest, "parse int error")
		return
	}
	// 渲染对应的模板
	updatedHTML := NumberConfig(name, int(intVal), name+"_Description", 0, 65535)
	if renderErr := htmx.NewResponse().RenderTempl(c.Request.Context(), c.Writer, updatedHTML); renderErr != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

// -------------------------
// 其他辅助逻辑保持不变
// -------------------------

// beforeConfigUpdate 根据配置的变化，判断是否需要打开浏览器重新扫描等
func beforeConfigUpdate(oldConfig *config.Config, newConfig *config.Config) {
	openBrowserIfNeeded(oldConfig, newConfig)
	needScan, reScanFile := checkReScanStatus(oldConfig, newConfig)
	if needScan {
		startReScan(reScanFile)
	} else {
		if newConfig.Debug {
			logger.Info("No changes in cfg, skipped scan store path\n")
		}
	}
}

func openBrowserIfNeeded(oldConfig *config.Config, newConfig *config.Config) {
	if !oldConfig.OpenBrowser && newConfig.OpenBrowser {
		protocol := "http://"
		if newConfig.EnableTLS {
			protocol = "https://"
		}
		util.OpenBrowser(protocol + "localhost:" + strconv.Itoa(newConfig.Port))
	}
}

// checkReScanStatus 检查旧的和新的配置是否需要更新，并返回需要重新扫描和重新扫描文件的布尔值
func checkReScanStatus(oldConfig *config.Config, newConfig *config.Config) (reScanDir bool, reScanFile bool) {
	if !reflect.DeepEqual(oldConfig.LocalStores, newConfig.LocalStores) {
		reScanDir = true
	}
	if oldConfig.MaxScanDepth != newConfig.MaxScanDepth {
		reScanDir = true
	}
	if !reflect.DeepEqual(oldConfig.SupportMediaType, newConfig.SupportMediaType) {
		reScanDir = true
		reScanFile = true
	}
	if !reflect.DeepEqual(oldConfig.SupportFileType, newConfig.SupportFileType) {
		reScanDir = true
	}
	if oldConfig.MinImageNum != newConfig.MinImageNum {
		reScanDir = true
		reScanFile = true
	}
	if !reflect.DeepEqual(oldConfig.ExcludePath, newConfig.ExcludePath) {
		reScanDir = true
	}
	if oldConfig.EnableDatabase != newConfig.EnableDatabase {
		reScanDir = true
	}
	return
}

// startReScan 扫描并相应地更新数据库
func startReScan(reScanFile bool) {
	option := scan.NewScanOption(
		reScanFile,
		config.GetLocalStoresList(),
		config.GetStores(),
		config.GetMaxScanDepth(),
		config.GetMinImageNum(),
		config.GetTimeoutLimitForScan(),
		config.GetExcludePath(),
		config.GetSupportMediaType(),
		config.GetSupportFileType(),
		config.GetSupportTemplateFile(),
		config.GetZipFileTextEncoding(),
		config.GetEnableDatabase(),
		config.GetClearDatabaseWhenExit(),
		config.GetDebug(),
	)
	if err := scan.AllStore(option); err != nil {
		logger.Infof("Failed to scan store path: %v", err)
	}
	if config.GetEnableDatabase() {
		saveResultsToDatabase(viper.ConfigFileUsed(), config.GetClearDatabaseWhenExit())
	}
}

func saveResultsToDatabase(configPath string, clearDatabaseWhenExit bool) {
	if err := scan.SaveResultsToDatabase(configPath, clearDatabaseWhenExit); err != nil {
		logger.Infof("Failed to save results to database: %v", err)
	}
}
