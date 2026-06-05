package store

import (
	"context"
	"net/http"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"
	"sync"

	"github.com/gin-gonic/gin"
)

type IShopCategoryStore interface {
	Set(rows []*shopmodels.CategoryInfo)
	Get() ([]*shopmodels.CategoryInfo, bool)
	GetCategoryIDs(id int64) []int64
	Clear()
}

type shopCategoryStore struct {
	service shopservice.IShopCategoryService
	mu      sync.RWMutex
	rows    []*shopmodels.CategoryInfo
}

func NewShopCategoryStore(service shopservice.IShopCategoryService) IShopCategoryStore {
	store := &shopCategoryStore{service: service}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	if rows, err := service.All(&gin.Context{Request: req}); err == nil {
		store.Set(rows)
	}
	return store
}

func (s *shopCategoryStore) Set(rows []*shopmodels.CategoryInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rows = copyCategoryInfos(rows)
}

func (s *shopCategoryStore) Get() ([]*shopmodels.CategoryInfo, bool) {
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

func copyCategoryInfos(rows []*shopmodels.CategoryInfo) []*shopmodels.CategoryInfo {
	if rows == nil {
		return nil
	}
	copied := make([]*shopmodels.CategoryInfo, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			copied = append(copied, nil)
			continue
		}
		item := *row
		item.Children = copyCategoryInfos(row.Children)
		copied = append(copied, &item)
	}
	return copied
}

func findCategoryInfo(row *shopmodels.CategoryInfo, id int64) *shopmodels.CategoryInfo {
	if row == nil {
		return nil
	}
	if row.ID == id {
		return row
	}
	for _, child := range row.Children {
		if target := findCategoryInfo(child, id); target != nil {
			return target
		}
	}
	return nil
}

func collectCategoryIDs(row *shopmodels.CategoryInfo, ids *[]int64) {
	if row == nil {
		return
	}
	*ids = append(*ids, row.ID)
	for _, child := range row.Children {
		collectCategoryIDs(child, ids)
	}
}
