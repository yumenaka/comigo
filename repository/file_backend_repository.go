package repository

import (
	"context"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/sqlc"
)

// FileBackendRepository 文件后端数据访问层
type FileBackendRepository struct {
	queries *sqlc.Queries
}

// NewFileBackendRepository 创建新的FileBackendRepository实例
func NewFileBackendRepository(db sqlc.DBTX) *FileBackendRepository {
	return &FileBackendRepository{
		queries: sqlc.New(db),
	}
}

// GetByID 根据ID获取文件后端
func (r *FileBackendRepository) GetByID(ctx context.Context, url string) (*model.FileBackend, error) {
	sqlcFileBackend, err := r.queries.GetFileBackendByID(ctx, url)
	if err != nil {
		return nil, err
	}
	return model.FromSQLCFileBackend(sqlcFileBackend), nil
}

// List 获取所有文件后端列表
func (r *FileBackendRepository) List(ctx context.Context) ([]*model.FileBackend, error) {
	sqlcFileBackends, err := r.queries.ListFileBackends(ctx)
	if err != nil {
		return nil, err
	}
	return model.FromSQLCFileBackends(sqlcFileBackends), nil
}

// ListByType 根据类型获取文件后端列表
func (r *FileBackendRepository) ListByType(ctx context.Context, backendType model.FileBackendType) ([]*model.FileBackend, error) {
	sqlcFileBackends, err := r.queries.ListFileBackendsByType(ctx, int64(backendType))
	if err != nil {
		return nil, err
	}
	return model.FromSQLCFileBackends(sqlcFileBackends), nil
}

// Create 创建新文件后端
func (r *FileBackendRepository) Create(ctx context.Context, fileBackend *model.FileBackend) (*model.FileBackend, error) {
	params := model.ToSQLCCreateFileBackendParams(fileBackend)
	sqlcFileBackend, err := r.queries.CreateFileBackend(ctx, params)
	if err != nil {
		return nil, err
	}
	return model.FromSQLCFileBackend(sqlcFileBackend), nil
}

// Update 更新文件后端信息
func (r *FileBackendRepository) Update(ctx context.Context, fileBackend *model.FileBackend) error {
	params := model.ToSQLCUpdateFileBackendParams(fileBackend)
	return r.queries.UpdateFileBackend(ctx, params)
}

// Delete 删除文件后端
func (r *FileBackendRepository) Delete(ctx context.Context, url string) error {
	return r.queries.DeleteFileBackend(ctx, url)
}

// CountByType 根据类型统计文件后端数量
func (r *FileBackendRepository) CountByType(ctx context.Context, backendType model.FileBackendType) (int64, error) {
	return r.queries.CountFileBackendsByType(ctx, int64(backendType))
}
