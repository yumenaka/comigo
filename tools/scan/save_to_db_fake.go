//go:build (windows && 386) || js

package scan

// SaveResultsToDatabase 4，保存扫描结果到数据库，并清理不存在的书籍
func SaveResultsToDatabase(ConfigPath string, ClearDatabaseWhenExit bool) error {
	return nil
}
