package scan

import (
	"github.com/yumenaka/comigo/internal/database"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/util/logger"
)

// SaveResultsToDatabase 4，保存扫描结果到数据库，并清理不存在的书籍
func SaveResultsToDatabase(ConfigPath string, ClearDatabaseWhenExit bool) error {
	err := database.InitDatabase(ConfigPath)
	if err != nil {
		return err
	}
	saveErr := database.SaveBookListToDatabase(model.MainStores.GetArchiveBooks())
	if saveErr != nil {
		logger.Info(saveErr)
		return saveErr
	}
	return nil
}

func ClearDatabaseWhenExit(ConfigPath string) {
	AllBook := model.MainStores.GetAllBookList()
	for _, b := range AllBook {
		database.ClearBookData(b)
	}
}
