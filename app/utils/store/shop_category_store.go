package store

import (
	"go.uber.org/zap"
	"sync"
)

// IShopCategoryStore 定义商城分类树的内存缓存能力。
// 实现方需要保证读写并发安全，并在返回数据时避免外部修改内部缓存。
type IShopCategoryStore interface {
	Set(rows []ShopCategoryData)
	Get() ([]ShopCategoryData, bool)
	GetCategoryIDs(id int64) []int64
	Clear()
}

// shopCategoryStore 使用读写锁保存一份分类树快照。
type shopCategoryStore struct {
	mu   sync.RWMutex
	rows []ShopCategoryData
}

// NewShopCategoryStore 创建默认的商城分类缓存实例。
func NewShopCategoryStore() IShopCategoryStore {
	store := &shopCategoryStore{}
	return store
}

// Set 覆盖当前缓存，并拷贝传入分类树，避免调用方后续修改影响缓存内容。
func (s *shopCategoryStore) Set(rows []ShopCategoryData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rows = copyCategoryInfos(rows)
}

// Get 返回分类树快照；第二个返回值表示缓存是否已初始化。
func (s *shopCategoryStore) Get() ([]ShopCategoryData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.rows == nil {
		return nil, false
	}
	return copyCategoryInfos(s.rows), true
}

// GetCategoryIDs 查找指定分类，并返回该分类及所有子分类的 ID。
func (s *shopCategoryStore) GetCategoryIDs(id int64) []int64 {
	if id == 0 {
		return []int64{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.rows == nil {
		return []int64{}
	}
	for _, row := range s.rows {
		if target := findCategoryInfo(row, id); target != nil {
			ids := make([]int64, 0)
			collectCategoryIDs(target, &ids)
			return ids
		}
	}
	return []int64{}
}

// Clear 清空当前缓存。
func (s *shopCategoryStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rows = nil
}

// copyCategoryInfos 递归复制分类树切片，避免暴露缓存中的可变 children。
func copyCategoryInfos(rows []ShopCategoryData) []ShopCategoryData {
	if rows == nil {
		return nil
	}
	copied := make([]ShopCategoryData, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			copied = append(copied, nil)
			continue
		}
		item := row
		err := item.SetChildren(copyCategoryInfos(row.ChildrenData()))
		if err != nil {
			zap.L().Error("set children failed", zap.Error(err))
			continue
		}
		copied = append(copied, item)
	}
	return copied
}

// findCategoryInfo 在分类树中深度优先查找指定分类。
func findCategoryInfo(row ShopCategoryData, id int64) ShopCategoryData {
	if row == nil {
		return nil
	}
	if row.CategoryID() == id {
		return row
	}
	for _, child := range row.ChildrenData() {
		if target := findCategoryInfo(child, id); target != nil {
			return target
		}
	}
	return nil
}

// collectCategoryIDs 以先序遍历收集当前分类及所有后代分类 ID。
func collectCategoryIDs(row ShopCategoryData, ids *[]int64) {
	if row == nil {
		return
	}
	*ids = append(*ids, row.CategoryID())
	for _, child := range row.ChildrenData() {
		collectCategoryIDs(child, ids)
	}
}
