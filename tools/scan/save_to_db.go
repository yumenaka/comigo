//go:build !(windows && 386) && !js

package scan

import (
	"fmt"
	"strconv"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/sqlc"
	"github.com/yumenaka/comigo/tools/logger"
)

// SaveResultsToDatabase 4，保存扫描结果到数据库，并清理不存在的书籍
func SaveResultsToDatabase(ConfigPath string, ClearDatabaseWhenExit bool) error {
	allBooks := model.MainStoreGroup.GetAllBookList()
	saveErr := sqlc.Repo.SaveBookListToDatabase(allBooks)
	if saveErr != nil {
		logger.Info(saveErr)
		return saveErr
	}
	fmt.Println("SaveResultsToDatabase: Books saved to database successfully.!!!!!!!!!!!!" + strconv.Itoa(len(allBooks)))
	return nil
}
