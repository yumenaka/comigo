package sqlc

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/yumenaka/comigo/model"
)

// TestCleanupDatabase 验证软删除记录可见，并且确认前恢复关联的书签不会被误删。
func TestCleanupDatabase(t *testing.T) {
	conn, db := newTestStoreDatabase(t)
	b := &model.Book{BookInfo: model.BookInfo{BookID: "Deleted", BookPath: "/library/book.cbz", StoreUrl: "/library", Type: model.TypeCbz, Deleted: true}}
	if err := db.StoreBook(b); err != nil {
		t.Fatal(err)
	}
	books, err := db.ListCleanupBooks()
	if err != nil || len(books) != 1 || !books[0].Deleted {
		t.Fatalf("soft deleted books: %+v %v", books, err)
	}
	if _, err := conn.Exec("INSERT INTO bookmarks (book_id,type,page_index) VALUES ('Orphan1',0,1),('Orphan2',0,2)"); err != nil {
		t.Fatal(err)
	}
	// 模拟旧版本未启用外键时留下的无主页面。
	if _, err := conn.Exec("PRAGMA foreign_keys = OFF"); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec("INSERT INTO page_infos (book_id,name) VALUES ('Orphan2','page.jpg')"); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatal(err)
	}
	rows, err := db.ListOrphanRows()
	if err != nil || len(rows["bookmarks"]) != 2 || len(rows["page_infos"]) != 1 {
		t.Fatalf("orphans: %+v %v", rows, err)
	}
	b.BookID = "Orphan1"
	b.BookPath = "/library/recovered.cbz"
	b.Deleted = false
	if err := db.StoreBook(b); err != nil {
		t.Fatal(err)
	}
	// StoreBook 会重写自己的书签，这里恢复同一行，验证删除端仍复核关联。
	if _, err := conn.Exec("INSERT INTO bookmarks (id,book_id,type,page_index) VALUES ($1,'Orphan1',0,1)", rows["bookmarks"][0]); err != nil {
		t.Fatal(err)
	}
	if err := db.DeleteOrphanRows(rows); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := conn.QueryRow("SELECT COUNT(*) FROM bookmarks WHERE book_id='Orphan1'").Scan(&count); err != nil || count != 1 {
		t.Fatalf("restored bookmark lost: %d %v", count, err)
	}
	rows, err = db.ListOrphanRows()
	if err != nil || len(rows["bookmarks"]) != 0 || len(rows["page_infos"]) != 0 {
		t.Fatalf("orphans remain: %+v %v", rows, err)
	}
}

// TestReclaimCleanupSpace 验证超过一批的空闲页实际归还磁盘，并且不会损坏保留的数据。
func TestReclaimCleanupSpace(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cleanup.sqlite")
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if err := configureSQLitePragmas(context.Background(), conn); err != nil {
		t.Fatal(err)
	}
	// 禁止自动 checkpoint，让测试能单独观察主动回收后的文件状态。
	if _, err := conn.Exec(`PRAGMA wal_autocheckpoint=0;
 CREATE TABLE payload (id INTEGER PRIMARY KEY, data BLOB);
 WITH RECURSIVE n(x) AS (VALUES(1) UNION ALL SELECT x+1 FROM n WHERE x<1500)
 INSERT INTO payload SELECT x,zeroblob(8192) FROM n;
 PRAGMA wal_checkpoint(TRUNCATE);
 DELETE FROM payload WHERE id>1;`); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	var freePages int
	if err := conn.QueryRow("PRAGMA freelist_count").Scan(&freePages); err != nil || freePages <= 1000 {
		t.Fatalf("need multiple batches: pages=%d err=%v", freePages, err)
	}
	db := NewDBStore(conn)
	if err := db.ReclaimCleanupSpace(); err != nil {
		t.Fatal(err)
	}
	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if after.Size() >= before.Size() {
		t.Fatalf("database did not shrink: before=%d after=%d", before.Size(), after.Size())
	}
	if err := conn.QueryRow("PRAGMA freelist_count").Scan(&freePages); err != nil || freePages != 0 {
		t.Fatalf("free pages remain: %d %v", freePages, err)
	}
	wal, err := os.Stat(path + "-wal")
	if err != nil || wal.Size() != 0 {
		t.Fatalf("WAL not truncated: %v %v", wal, err)
	}
	var count, size int
	if err := conn.QueryRow("SELECT count(*),sum(length(data)) FROM payload").Scan(&count, &size); err != nil || count != 1 || size != 8192 {
		t.Fatalf("retained data damaged: %d %d %v", count, size, err)
	}
	var integrity string
	if err := conn.QueryRow("PRAGMA integrity_check").Scan(&integrity); err != nil || integrity != "ok" {
		t.Fatalf("integrity=%q err=%v", integrity, err)
	}
	// 外部读事务占用旧 WAL 快照时必须报告 busy，结束读事务后可以再次回收。
	reader, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	tx, err := reader.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	if err := tx.QueryRow("SELECT count(*) FROM payload").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec("PRAGMA busy_timeout=1; DELETE FROM payload"); err != nil {
		t.Fatal(err)
	}
	if err := db.ReclaimCleanupSpace(); err == nil {
		t.Fatal("busy WAL checkpoint was reported as successful")
	}
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}
	if err := db.ReclaimCleanupSpace(); err != nil {
		t.Fatal(err)
	}
	if err := conn.QueryRow("SELECT count(*) FROM payload").Scan(&count); err != nil || count != 0 {
		t.Fatalf("deletion lost: %d %v", count, err)
	}
}

// TestReclaimCleanupSpaceSkipsPostgres 确保 PostgreSQL 分支不执行 SQLite 专属语句。
func TestReclaimCleanupSpaceSkipsPostgres(t *testing.T) {
	db := &StoreDatabase{postgres: true}
	if err := db.ReclaimCleanupSpace(); err != nil {
		t.Fatal(err)
	}
}
