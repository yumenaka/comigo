//go:build !js

package scan

import (
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/sqlc"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/tools/logger"
)

// SaveBooksToDatabase 4，保存扫描结果到数据库
func SaveBooksToDatabase(cfg ConfigInterface) error {
	InitConfig(cfg)
	allBooks, err := store.RamStore.ListBooks()
	if err != nil {
		logger.Infof(locale.GetString("log_error_listing_books"), err)
	}
	for _, b := range allBooks {
		saveErr := sqlc.DbStore.StoreBook(b)
		if saveErr != nil {
			logger.Info(saveErr)
			return saveErr
		}
	}
	logger.Infof(locale.GetString("log_books_saved_to_database_successfully"), len(allBooks))
	return nil
}
