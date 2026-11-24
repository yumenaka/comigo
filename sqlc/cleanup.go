//go:build !js

package sqlc

import (
	"context"
	"fmt"

	"github.com/yumenaka/comigo/model"
)

// ListCleanupBooks 包含软删除记录；先关闭游标，再加载关联页面，适用于 SQLite 单连接。
func (db *StoreDatabase) ListCleanupBooks() ([]*model.Book, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	rows, err := db.connection.Query("SELECT book_id FROM books ORDER BY book_id")
	if err != nil {
		return nil, err
	}
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}
	books := make([]*model.Book, 0, len(ids))
	for _, id := range ids {
		book, err := db.getBook(id)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// ListOrphanRows 查找历史遗留的无主页面和书签，不把其他业务表纳入清理。
func (db *StoreDatabase) ListOrphanRows() (map[string][]int64, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	result := map[string][]int64{}
	for _, table := range []string{"page_infos", "bookmarks"} {
		rows, err := db.connection.Query("SELECT id FROM " + table + " WHERE NOT EXISTS (SELECT 1 FROM books WHERE books.book_id = " + table + ".book_id) ORDER BY id")
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var id int64
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				return nil, err
			}
			result[table] = append(result[table], id)
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// DeleteOrphanRows 只删除已预览且仍无主的记录；表名固定，不接受外部 SQL 标识符。
func (db *StoreDatabase) DeleteOrphanRows(records map[string][]int64) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	tx, err := db.connection.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, table := range []string{"page_infos", "bookmarks"} {
		for _, id := range records[table] {
			if _, err := tx.Exec("DELETE FROM "+table+" WHERE id = $1 AND NOT EXISTS (SELECT 1 FROM books WHERE books.book_id = "+table+".book_id)", id); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

// ReclaimCleanupSpace 仅在手动清理后回收 SQLite 空闲页，不在日常写入时压缩或转换旧库。
func (db *StoreDatabase) ReclaimCleanupSpace() error {
	if db.postgres {
		return nil
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	var freePages int
	if err := db.connection.QueryRow("PRAGMA freelist_count").Scan(&freePages); err != nil {
		return err
	}
	// 以开始时的空闲页数限定工作量；每批最多回收 1000 页，无需完整重建数据库。
	for remaining := freePages; remaining > 0; remaining -= 1000 {
		rows, err := db.connection.Query("PRAGMA incremental_vacuum(1000)")
		if err != nil {
			return err
		}
		// 此 PRAGMA 会逐页返回结果，必须读完；Exec 可能只推进第一步。
		for rows.Next() {
		}
		err = rows.Err()
		closeErr := rows.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
	}
	var busy, logFrames, checkpointed int
	if err := db.connection.QueryRow("PRAGMA wal_checkpoint(TRUNCATE)").Scan(&busy, &logFrames, &checkpointed); err != nil {
		return err
	}
	if busy != 0 {
		return fmt.Errorf("SQLite WAL checkpoint busy: %d/%d pages checkpointed", checkpointed, logFrames)
	}
	return nil
}
