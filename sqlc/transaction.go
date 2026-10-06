package sqlc

import (
	"context"

	"github.com/yumenaka/comigo/sqlc/postgres"
)

// write 保证书籍、页面及书签一起提交；读改写共享锁，避免并发书签互相覆盖。
// ponytail: 每个书库实例串行写入；写吞吐成为瓶颈时再按书籍拆分锁。
func (db *StoreDatabase) write(fn func(*StoreDatabase) error) error {
	if err := db.CheckDBQueries(); err != nil {
		return err
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	tx, err := db.connection.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	queries := bookQueries(New(tx))
	if db.postgres {
		queries = newPostgresAdapter(postgres.New(tx))
	}
	if err := fn(&StoreDatabase{queries: queries}); err != nil {
		return err
	}
	return tx.Commit()
}
