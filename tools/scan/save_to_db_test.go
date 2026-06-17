//go:build !js

package scan

import (
	"strings"
	"testing"

	"github.com/yumenaka/comigo/sqlc"
)

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
