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
func SaveResultsToDatabase(cfg ConfigInterface) error {
	InitConfig(cfg)
	books := model.IStore.ListBookSkipBookGroup()
	saveErr := sqlc.Repo.SaveBookListToDatabase(books)
	if saveErr != nil {
		logger.Info(saveErr)
		return saveErr
	}
	fmt.Println("SaveResultsToDatabase: Books saved to database successfully: " + strconv.Itoa(len(books)))
	return nil
}
