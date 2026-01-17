package cmd

import (
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/sqlc"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/tools/logger"
	"github.com/yumenaka/comigo/tools/scan"
)

// LoadMetadata 加载书籍元数据
func LoadMetadata() {
	// 从数据库加载书籍信息
	if config.GetCfg().EnableDatabase {
		// 从数据库中读取书籍信息并持久化
		configDir, err := config.GetConfigDir()
		if err != nil {
			logger.Infof(locale.GetString("err_failed_to_get_config_dir"), err)
			model.IStore = store.RamStore
			return
		}
		if err := sqlc.OpenDatabase(configDir); err != nil {
			logger.Infof(locale.GetString("log_open_database_error"), err)
			model.IStore = store.RamStore
		} else {
			model.IStore = sqlc.DbStore
		}
	}
	// 从本地文件加载书籍信息
	if !config.GetCfg().EnableDatabase {
		err := store.RamStore.LoadBooks()
		if err != nil {
			logger.Infof(locale.GetString("log_loadbooks_error"), err)
		}
		model.ClearBookWhenStoreUrlNotExist(config.GetCfg().StoreUrls)
		model.ClearBookNotExist()
		model.GenerateBookGroup()
	}
}

// SaveMetadata 保存书籍元数据
func SaveMetadata() {
	// 未启用数据库的时候，保存书籍元数据到本地Json文件
	if !config.GetCfg().EnableDatabase {
		err := store.RamStore.SaveBooksToJson()
		if err != nil {
			logger.Infof(locale.GetString("log_savebooks_error"), err)
		}
	}
	// 启用数据库的时候，同步书籍元数据到RamStore
	if config.GetCfg().EnableDatabase {
		allBooks, err := sqlc.DbStore.ListBooks()
		if err != nil {
			logger.Infof(locale.GetString("log_error_listing_books_from_database"), err)
		} else {
			// 兜底：万一数据库无效，至把书加回RamStore
			err = store.RamStore.AddBooks(allBooks)
			if err != nil {
				return
			}
		}
	}
	// 启用数据库的时候，保存书籍元数据到到数据库
	if config.GetCfg().EnableDatabase {
		err := scan.SaveBooksToDatabase(config.GetCfg())
		if err != nil {
			logger.Infof(locale.GetString("log_failed_savebookstodatabase"), err)
			return
		}
	}
}
