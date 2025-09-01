package sqlc

import (
	"context"

	"github.com/yumenaka/comigo/model"
)

// StoreRepository 书库数据访问层
type StoreRepository struct {
	queries *Queries
}

// NewStoreRepository 创建新的StoreRepository实例
func NewStoreRepository(db DBTX) *StoreRepository {
	return &StoreRepository{
		queries: New(db),
	}
}

// GetByBackendUrl  根据ID获取书库
func (r *StoreRepository) GetByBackendUrl(ctx context.Context, backendUrl string) (*model.StoreInfo, error) {
	sqlcStore, err := r.queries.GetStoreByBackendURL(ctx, backendUrl)
	if err != nil {
		return nil, err
	}
	return FromSQLCStore(sqlcStore), nil
}

// GetByName 根据名称获取书库
func (r *StoreRepository) GetByName(ctx context.Context, name string) (*model.StoreInfo, error) {
	sqlcStore, err := r.queries.GetStoreByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return FromSQLCStore(sqlcStore), nil
}

// List 获取所有书库列表
func (r *StoreRepository) List(ctx context.Context) ([]*model.StoreInfo, error) {
	sqlcStores, err := r.queries.ListStores(ctx)
	if err != nil {
		return nil, err
	}
	return FromSQLCStores(sqlcStores), nil
}

// GetWithBackend 获取书库及其后端信息
func (r *StoreRepository) GetWithBackend(ctx context.Context, fileBackendUrl string) (*model.StoreInfo, error) {
	sqlcStoreWithBackend, err := r.queries.GetStoreWithBackend(ctx, fileBackendUrl)
	if err != nil {
		return nil, err
	}
	return FromSQLCStoreWithBackendRow(sqlcStoreWithBackend), nil
}

// ListWithBackend 获取所有书库及其后端信息
func (r *StoreRepository) ListWithBackend(ctx context.Context) ([]*model.StoreInfo, error) {
	sqlcStoresWithBackend, err := r.queries.ListStoresWithBackend(ctx)
	if err != nil {
		return nil, err
	}
	return FromSQLCListStoresWithBackendRow(sqlcStoresWithBackend), nil
}

// Create 创建新书库
func (r *StoreRepository) Create(ctx context.Context, store *model.StoreInfo) error {
	params := ToSQLCCreateStoreParams(store)
	_, err := r.queries.CreateStore(ctx, params)
	return err
}

// Update 更新书库信息
func (r *StoreRepository) Update(ctx context.Context, store *model.StoreInfo) error {
	params := ToSQLCUpdateStoreParams(store)
	return r.queries.UpdateStore(ctx, params)
}

// Delete 删除书库
func (r *StoreRepository) Delete(ctx context.Context, fileBackendUrl string) error {
	return r.queries.DeleteStore(ctx, fileBackendUrl)
}

// Count 统计书库总数
func (r *StoreRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountStores(ctx)
}

// CreateStoreWithBackend 创建书库及其文件后端（事务操作）
// 注意：这里需要在实际使用时实现事务逻辑
// 示例代码，实际使用时需要根据具体的数据库驱动实现事务
func (r *StoreRepository) CreateStoreWithBackend(ctx context.Context, store *model.StoreInfo, fileBackend *model.Backend) error {
	// 1. 创建文件后端
	params := ToSQLCCreateFileBackendParams(fileBackend)
	_, err := r.queries.CreateFileBackend(ctx, params)
	if err != nil {
		return err
	}

	// 2. 创建书库
	return r.Create(ctx, store)
}
