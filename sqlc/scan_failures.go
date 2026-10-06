package sqlc

import (
	"context"

	"github.com/yumenaka/comigo/model"
)

// LoadScanFailures 读取失败文件指纹，重启后仍可跳过未变化的损坏文件。
func (db *StoreDatabase) LoadScanFailures() (map[string]model.ScanFailureRecord, error) {
	if err := db.CheckDBQueries(); err != nil {
		return nil, err
	}
	db.mu.RLock()
	defer db.mu.RUnlock()
	rows, err := db.connection.QueryContext(context.Background(), "SELECT key, store_url, file_path, file_size, modified_unix_nano, created_by_version, failed_at, error, is_remote FROM scan_failures")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	records := map[string]model.ScanFailureRecord{}
	for rows.Next() {
		var key string
		var record model.ScanFailureRecord
		if err := rows.Scan(&key, &record.StoreURL, &record.FilePath, &record.FileSize, &record.ModifiedUnixNano, &record.CreatedByVersion, &record.FailedAt, &record.Error, &record.IsRemote); err != nil {
			return nil, err
		}
		records[key] = record
	}
	return records, rows.Err()
}

// SaveScanFailures 用事务替换失败缓存；失败时保留上一次完整记录。
func (db *StoreDatabase) SaveScanFailures(records map[string]model.ScanFailureRecord) error {
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
	if _, err := tx.Exec("DELETE FROM scan_failures"); err != nil {
		return err
	}
	for key, r := range records {
		if _, err := tx.Exec("INSERT INTO scan_failures (key, store_url, file_path, file_size, modified_unix_nano, created_by_version, failed_at, error, is_remote) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)", key, r.StoreURL, r.FilePath, r.FileSize, r.ModifiedUnixNano, r.CreatedByVersion, r.FailedAt, r.Error, r.IsRemote); err != nil {
			return err
		}
	}
	return tx.Commit()
}
