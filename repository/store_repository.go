package repository

import (
	"context"

	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/sqlc"
)

// StoreRepository 书库数据访问层
type StoreRepository struct {
	queries *sqlc.Queries
}

// NewStoreRepository 创建新的StoreRepository实例
func NewStoreRepository(db sqlc.DBTX) *StoreRepository {
	return &StoreRepository{
		queries: sqlc.New(db),
	}
}

// GetByID 根据ID获取书库
func (r *StoreRepository) GetByID(ctx context.Context, id int) (*model.Store, error) {
	sqlcStore, err := r.queries.GetStoreByID(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	return model.FromSQLCStore(sqlcStore), nil
}

// GetByName 根据名称获取书库
func (r *StoreRepository) GetByName(ctx context.Context, name string) (*model.Store, error) {
	sqlcStore, err := r.queries.GetStoreByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return model.FromSQLCStore(sqlcStore), nil
}

// List 获取所有书库列表
func (r *StoreRepository) List(ctx context.Context) ([]*model.Store, error) {
	sqlcStores, err := r.queries.ListStores(ctx)
	if err != nil {
		return nil, err
	}
	return model.FromSQLCStores(sqlcStores), nil
}

// GetWithBackend 获取书库及其后端信息
func (r *StoreRepository) GetWithBackend(ctx context.Context, id int) (*model.Store, error) {
	sqlcStoreWithBackend, err := r.queries.GetStoreWithBackend(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	return model.FromSQLCStoreWithBackendRow(sqlcStoreWithBackend), nil
}

// ListWithBackend 获取所有书库及其后端信息
func (r *StoreRepository) ListWithBackend(ctx context.Context) ([]*model.Store, error) {
	sqlcStoresWithBackend, err := r.queries.ListStoresWithBackend(ctx)
	if err != nil {
		return nil, err
	}
	return model.FromSQLCListStoresWithBackendRow(sqlcStoresWithBackend), nil
}

// Create 创建新书库
func (r *StoreRepository) Create(ctx context.Context, store *model.Store, fileBackendID int64) error {
	params := model.ToSQLCCreateStoreParams(store, fileBackendID)
	_, err := r.queries.CreateStore(ctx, params)
	return err
}

// Update 更新书库信息
func (r *StoreRepository) Update(ctx context.Context, store *model.Store, fileBackendID int64) error {
	params := model.ToSQLCUpdateStoreParams(store, fileBackendID)
	return r.queries.UpdateStore(ctx, params)
}

// Delete 删除书库
func (r *StoreRepository) Delete(ctx context.Context, id int) error {
	return r.queries.DeleteStore(ctx, int64(id))
}

// Count 统计书库总数
func (r *StoreRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountStores(ctx)
}

// CreateStoreWithBackend 创建书库及其文件后端（事务操作）
// 注意：这里需要在实际使用时实现事务逻辑
// 示例代码，实际使用时需要根据具体的数据库驱动实现事务
func (r *StoreRepository) CreateStoreWithBackend(ctx context.Context, store *model.Store, fileBackend *model.FileBackend) error {
	// 1. 创建文件后端
	params := model.ToSQLCCreateFileBackendParams(fileBackend)
	sqlcFileBackend, err := r.queries.CreateFileBackend(ctx, params)
	if err != nil {
		return err
	}

	// 2. 创建书库
	return r.Create(ctx, store, sqlcFileBackend.ID)
}
