package shopserviceimpl

import (
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"

	"github.com/gin-gonic/gin"
)

type ShopCategoryServiceImpl struct {
	dao shopdao.IShopCategoryDao
}

func NewShopCategoryService(dao shopdao.IShopCategoryDao) shopservice.IShopCategoryService {
	return &ShopCategoryServiceImpl{dao: dao}
}

func (s *ShopCategoryServiceImpl) Create(c *gin.Context, req *shopmodels.CategoryUpsert) (*shopmodels.Category, error) {
	return s.dao.Create(c, req)
}

func (s *ShopCategoryServiceImpl) Update(c *gin.Context, req *shopmodels.CategoryUpsert) (*shopmodels.Category, error) {
	return s.dao.Update(c, req)
}

func (s *ShopCategoryServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ShopCategoryServiceImpl) GetByID(c *gin.Context, id int64) (*shopmodels.Category, error) {
	return s.dao.GetByID(c, id)
}

func (s *ShopCategoryServiceImpl) All(c *gin.Context) ([]*shopmodels.CategoryInfo, error) {
	rows, err := s.dao.All(c)
	if err != nil {
		return nil, err
	}
	return buildCategoryTree(rows), nil
}

func (s *ShopCategoryServiceImpl) List(c *gin.Context, req *shopmodels.CategoryQuery) (*shopmodels.CategoryListData, error) {
	return s.dao.List(c, req)
}

func buildCategoryTree(rows []*shopmodels.Category) []*shopmodels.CategoryInfo {
	if len(rows) == 0 {
		return []*shopmodels.CategoryInfo{}
	}

	categoryMap := make(map[int64]*shopmodels.CategoryInfo, len(rows))
	roots := make([]*shopmodels.CategoryInfo, 0)

	for _, row := range rows {
		if row == nil {
			continue
		}
		info := toCategoryInfo(row)
		info.Children = make([]*shopmodels.CategoryInfo, 0)
		categoryMap[row.ID] = info
	}

	for _, row := range rows {
		if row == nil {
			continue
		}
		info := categoryMap[row.ID]
		if info == nil {
			continue
		}

		if info.ParentID > 0 {
			if parent, ok := categoryMap[info.ParentID]; ok {
				parent.Children = append(parent.Children, info)
				continue
			}
		}

		roots = append(roots, info)
	}

	return roots
}

func toCategoryInfo(row *shopmodels.Category) *shopmodels.CategoryInfo {
	if row == nil {
		return nil
	}

	return &shopmodels.CategoryInfo{
		ID:           row.ID,
		ParentID:     row.ParentID,
		AncestorPath: row.AncestorPath,
		Depth:        row.Depth,
		CategoryName: row.CategoryName,
		CategoryCode: row.CategoryCode,
		Sort:         row.Sort,
		Status:       row.Status,
		CreateTime:   row.CreateTime,
		UpdateTime:   row.UpdateTime,
	}
}
