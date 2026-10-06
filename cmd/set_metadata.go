package cmd

import (
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/sqlc"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/tools/logger"
)

// LoadMetadata 加载书籍元数据
func LoadMetadata() error {
	// 从数据库加载书籍信息
	if config.GetCfg().EnableDatabase {
		// 数据库是唯一持久化来源；打开失败直接返回，不回退到 JSON。
		configDir, err := config.GetConfigDir()
		if err != nil {
			return err
		}
		if err := sqlc.OpenDatabase(sqlc.DBOptions{
			Type:      config.GetCfg().DBType,
			DSN:       config.GetCfg().DBDSN,
			ConfigDir: configDir,
		}); err != nil {
			return err
		}
		model.IStore = sqlc.DbStore
		model.ClearBookWhenStoreUrlNotExist(config.GetCfg().StoreUrls)
		model.ClearBookNotExist()
	}
	// 从本地文件加载书籍信息
	if !config.GetCfg().EnableDatabase {
		model.IStore = store.RamStore
		err := store.RamStore.LoadBooks()
		if err != nil {
			logger.Infof(locale.GetString("log_loadbooks_error"), err)
		}
		model.ClearBookWhenStoreUrlNotExist(config.GetCfg().StoreUrls)
		model.ClearBookNotExist()
		// 旧数据需要重建时立即扫描；嵌入式宿主关闭启动扫描也应自动完成版本迁移。
		for range store.RamStore.PendingBooks.Range {
			ScanStore()
			break
		}
		model.GenerateBookGroup()
	}
	return nil
}

// SaveMetadata 保存书籍元数据
func SaveMetadata() {
	// 未启用数据库的时候，保存书籍元数据到本地Json文件
	if !config.GetCfg().EnableDatabase {
		err := store.RamStore.SaveAllBooksMetaJson()
		if err != nil {
			logger.Infof(locale.GetString("log_savebooks_error"), err)
		}
	}
}
