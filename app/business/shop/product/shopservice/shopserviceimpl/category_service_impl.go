package shopserviceimpl

import (
	"fmt"
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/utils/fileUtils"
	"nova-factory-server/app/utils/snowflake"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ShopCategoryServiceImpl struct {
	dao shopdao.IShopCategoryDao
}

func NewShopCategoryService(dao shopdao.IShopCategoryDao) shopservice.IShopCategoryService {
	return &ShopCategoryServiceImpl{dao: dao}
}

func (s *ShopCategoryServiceImpl) Create(c *gin.Context, req *shopmodels.CategoryUpsert) (*shopmodels.Category, error) {
	req.CategoryCode = normalizeCategoryCode(req.CategoryCode)
	req.ImageURL = strings.TrimSpace(req.ImageURL)
	req.Description = strings.TrimSpace(req.Description)
	if req.CategoryCode == "" {
		req.CategoryCode = generateCategoryCode()
	}
	return s.dao.Create(c, req)
}

func (s *ShopCategoryServiceImpl) Update(c *gin.Context, req *shopmodels.CategoryUpsert) (*shopmodels.Category, error) {
	req.CategoryCode = normalizeCategoryCode(req.CategoryCode)
	req.ImageURL = strings.TrimSpace(req.ImageURL)
	req.Description = strings.TrimSpace(req.Description)
	if req.CategoryCode == "" {
		req.CategoryCode = generateCategoryCode()
	}
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
	for k, v := range rows {
		rows[k].ImageURL = fileUtils.BuildAbsoluteURL(c, v.ImageURL)
	}
	return buildCategoryTree(rows), nil
}

func (s *ShopCategoryServiceImpl) List(c *gin.Context, req *shopmodels.CategoryQuery) (*shopmodels.CategoryListData, error) {
	data, err := s.dao.List(c, req)
	if err != nil {
		return nil, err
	}
	for _, row := range data.Rows {
		if row == nil {
			continue
		}
		row.ImageURL = fileUtils.BuildAbsoluteURL(c, row.ImageURL)
	}
	return data, nil
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
		ImageURL:     row.ImageURL,
		Description:  row.Description,
		Sort:         row.Sort,
		Status:       row.Status,
		CreateTime:   row.CreateTime,
		UpdateTime:   row.UpdateTime,
	}
}

func normalizeCategoryCode(code string) string {
	return strings.TrimSpace(code)
}

func generateCategoryCode() string {
	now := time.Now()
	uniqueNumber := snowflake.GenID()
	return fmt.Sprintf("CAT%s%d", now.Format("20060102150405"), uniqueNumber)
}
