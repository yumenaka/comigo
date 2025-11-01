//go:build !(windows && 386) && !js

package scan

import (
	"strconv"

	"github.com/yumenaka/comigo/sqlc"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/tools/logger"
)

// SaveBooksToDatabase 4，保存扫描结果到数据库
func SaveBooksToDatabase(cfg ConfigInterface) error {
	InitConfig(cfg)
	allBooks, err := store.RamStore.ListBooks()
	if err != nil {
		logger.Infof("Error listing books: %s", err)
	}
	for _, b := range allBooks {
		saveErr := sqlc.DbStore.AddBook(b)
		if saveErr != nil {
			logger.Info(saveErr)
			return saveErr
		}
	}
	logger.Infof("SaveBooksToDatabase: Books saved to database successfully: " + strconv.Itoa(len(allBooks)))
	return nil
}
