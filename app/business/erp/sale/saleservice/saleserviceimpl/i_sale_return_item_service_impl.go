package saleserviceimpl

import (
	"errors"
	"strings"

	"nova-factory-server/app/business/erp/erpbiz"
	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/business/erp/sale/saleservice"

	"github.com/gin-gonic/gin"
)

// SaleReturnItemServiceImpl 提供业务实现。
type SaleReturnItemServiceImpl struct {
	dao          saledao.ISaleReturnItemDao
	uniqueFields []erpbiz.UniqueField
}

// NewSaleReturnItemService 创建服务。
func NewSaleReturnItemService(dao saledao.ISaleReturnItemDao) saleservice.ISaleReturnItemService {
	return &SaleReturnItemServiceImpl{
		dao:          dao,
		uniqueFields: []erpbiz.UniqueField{},
	}
}

func (s *SaleReturnItemServiceImpl) create(c *gin.Context, req *salemodels.SaleReturnItemUpsert) (*salemodels.SaleReturnItem, error) {
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

func (s *SaleReturnItemServiceImpl) update(c *gin.Context, req *salemodels.SaleReturnItemUpsert) (*salemodels.SaleReturnItem, error) {
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

func (s *SaleReturnItemServiceImpl) Create(c *gin.Context, req *salemodels.SaleReturnItemUpsert) (*salemodels.SaleReturnItem, error) {
	return s.create(c, req)
}

func (s *SaleReturnItemServiceImpl) Update(c *gin.Context, req *salemodels.SaleReturnItemUpsert) (*salemodels.SaleReturnItem, error) {
	return s.update(c, req)
}

func (s *SaleReturnItemServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *SaleReturnItemServiceImpl) GetByID(c *gin.Context, id int64) (*salemodels.SaleReturnItem, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *SaleReturnItemServiceImpl) ListPage(c *gin.Context, req *salemodels.SaleReturnItemQuery) (*erpbiz.PageResult[salemodels.SaleReturnItem], error) {
	if req != nil {
		erpbiz.TrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *SaleReturnItemServiceImpl) validateUniqueFields(c *gin.Context, req *salemodels.SaleReturnItemUpsert, currentID int64) error {
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

func (s *SaleReturnItemServiceImpl) List(c *gin.Context, req *salemodels.SaleReturnItemQuery) (*salemodels.SaleReturnItemListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &salemodels.SaleReturnItemListData{Rows: result.Rows, Total: result.Total}, nil
}
