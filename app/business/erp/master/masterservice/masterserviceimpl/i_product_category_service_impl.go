package masterserviceimpl

import (
	"errors"
	"strings"

	"nova-factory-server/app/business/erp/erpbiz"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"
	"nova-factory-server/app/utils/category"

	"github.com/gin-gonic/gin"
)

// ProductCategoryServiceImpl 提供业务实现。
type ProductCategoryServiceImpl struct {
	dao          masterdao.IProductCategoryDao
	uniqueFields []erpbiz.UniqueField
}

// NewProductCategoryService 创建服务。
func NewProductCategoryService(dao masterdao.IProductCategoryDao) masterservice.IProductCategoryService {
	return &ProductCategoryServiceImpl{
		dao:          dao,
		uniqueFields: []erpbiz.UniqueField{},
	}
}

func (s *ProductCategoryServiceImpl) create(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	erpbiz.TrimStringFields(req)
	if err := erpbiz.ValidateRequiredFields(req); err != nil {
		return nil, err
	}
	if err := s.validateUniqueFields(c, req, 0); err != nil {
		return nil, err
	}
	return s.dao.Create(c, req)
}

func (s *ProductCategoryServiceImpl) update(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	erpbiz.TrimStringFields(req)
	if err := erpbiz.ValidateRequiredFields(req); err != nil {
		return nil, err
	}
	if err := s.validateUniqueFields(c, req, id); err != nil {
		return nil, err
	}
	return s.dao.Update(c, req)
}

// Create 新增 ERP 产品分类。
func (s *ProductCategoryServiceImpl) Create(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error) {
	s.ensureCode(req)
	return s.create(c, req)
}

// Update 修改 ERP 产品分类。
func (s *ProductCategoryServiceImpl) Update(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error) {
	s.ensureCode(req)
	return s.update(c, req)
}

func (s *ProductCategoryServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ProductCategoryServiceImpl) GetByID(c *gin.Context, id int64) (*mastermodels.ProductCategory, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *ProductCategoryServiceImpl) ListPage(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*erpbiz.PageResult[mastermodels.ProductCategory], error) {
	if req != nil {
		erpbiz.TrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *ProductCategoryServiceImpl) validateUniqueFields(c *gin.Context, req *mastermodels.ProductCategoryUpsert, currentID int64) error {
	if len(s.uniqueFields) == 0 {
		return nil
	}
	for _, field := range s.uniqueFields {
		value, ok := erpbiz.GetFieldValue(req, field.Field)
		if !ok {
			continue
		}
		normalized, empty := erpbiz.NormalizeValue(value)
		if empty {
			continue
		}
		exists, err := s.dao.GetByColumn(c, field.Column, normalized)
		if err != nil {
			return err
		}
		if exists == nil {
			continue
		}
		if erpbiz.GetIntField(exists, "ID") != currentID {
			label := strings.TrimSpace(field.Label)
			if label == "" {
				label = field.Column
			}
			return errors.New(label + "已存在")
		}
	}
	return nil
}

// List 查询 ERP 产品分类列表。
func (s *ProductCategoryServiceImpl) List(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*mastermodels.ProductCategoryListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.ProductCategoryListData{Rows: result.Rows, Total: result.Total}, nil
}

func (s *ProductCategoryServiceImpl) ensureCode(req *mastermodels.ProductCategoryUpsert) {
	if req == nil {
		return
	}
	req.Code = strings.TrimSpace(req.Code)
	if req.Code == "" {
		req.Code = category.GenerateProductCategoryCode()
	}
}
