package cmd

import (
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
			logger.Errorf("Failed to get config dir: %v", err)
			configDir = ""
		}
		if err := sqlc.OpenDatabase(configDir); err != nil {
			logger.Infof("OpenDatabase Error: %s", err)
			model.IStore = store.RamStore
		} else {
			model.IStore = sqlc.DbStore
		}
	}
	// 从本地文件加载书籍信息
	if !config.GetCfg().EnableDatabase {
		err := store.RamStore.LoadBooks()
		if err != nil {
			logger.Infof("LoadBooks_error %s", err)
		}
		model.ClearBookNotExist()
		// 生成虚拟书籍组
		if err := model.IStore.GenerateBookGroup(); err != nil {
			logger.Infof("%s", err)
		}
	}
}

// SaveMetadata 保存书籍元数据
func SaveMetadata() {
	// 未启用数据库的时候，保存书籍元数据到本地Json文件
	if !config.GetCfg().EnableDatabase {
		err := store.RamStore.SaveBooksToJson()
		if err != nil {
			logger.Infof("SaveBooks_error %s", err)
		}
	}
	// 启用数据库的时候，同步书籍元数据到RamStore
	if config.GetCfg().EnableDatabase {
		allBooks, err := sqlc.DbStore.ListBooks()
		if err != nil {
			logger.Infof("Error listing books from database: %s", err)
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
			logger.Infof("Failed SaveBooksToDatabase: %v", err)
			return
		}
	}
	// 重新生成书组
	if err := model.IStore.GenerateBookGroup(); err != nil {
		logger.Infof("%s", err)
	}
}
