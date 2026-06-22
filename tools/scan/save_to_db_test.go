//go:build !js

package scan

import (
	"strings"
	"testing"

	"github.com/yumenaka/comigo/sqlc"
)

// 验证保存书籍到数据库前必须先初始化存储。
func TestSaveBooksToDatabaseRequiresInitializedStore(t *testing.T) {
	oldStore := sqlc.DbStore
	sqlc.DbStore = nil
	t.Cleanup(func() {
		sqlc.DbStore = oldStore
	})

	err := SaveBooksToDatabase(&testCfgScan{})
	if err == nil || !strings.Contains(err.Error(), "database store is not initialized") {
		t.Fatalf("SaveBooksToDatabase error = %v, want uninitialized store error", err)
	}
}
