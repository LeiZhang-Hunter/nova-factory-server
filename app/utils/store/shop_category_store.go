package store

import (
	"go.uber.org/zap"
	"sync"
)

type IShopCategoryStore interface {
	Set(rows []ShopCategoryData)
	Get() ([]ShopCategoryData, bool)
	GetCategoryIDs(id int64) []int64
	Clear()
}

type shopCategoryStore struct {
	mu   sync.RWMutex
	rows []ShopCategoryData
}

func NewShopCategoryStore() IShopCategoryStore {
	store := &shopCategoryStore{}
	return store
}

func (s *shopCategoryStore) Set(rows []ShopCategoryData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rows = copyCategoryInfos(rows)
}

func (s *shopCategoryStore) Get() ([]ShopCategoryData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.rows == nil {
		return nil, false
	}
	return copyCategoryInfos(s.rows), true
}

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

func (s *shopCategoryStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rows = nil
}

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

func collectCategoryIDs(row ShopCategoryData, ids *[]int64) {
	if row == nil {
		return
	}
	*ids = append(*ids, row.CategoryID())
	for _, child := range row.ChildrenData() {
		collectCategoryIDs(child, ids)
	}
}
