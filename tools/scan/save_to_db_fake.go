//go:build (windows && 386) || js

package scan

// SaveBooksToDatabase 4，保存扫描结果到数据库，并清理不存在的书籍
func SaveResultsToDatabase(ConfigDir string, ClearDatabaseWhenExit bool) error {
	return nil
}
