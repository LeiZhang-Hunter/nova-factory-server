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

// SaleOutServiceImpl 提供业务实现。
type SaleOutServiceImpl struct {
	dao          saledao.ISaleOutDao
	uniqueFields []erpbiz.UniqueField
}

// NewSaleOutService 创建服务。
func NewSaleOutService(dao saledao.ISaleOutDao) saleservice.ISaleOutService {
	return &SaleOutServiceImpl{
		dao:          dao,
		uniqueFields: []erpbiz.UniqueField{{Field: "No", Column: "no", Label: "销售出库单号"}},
	}
}

func (s *SaleOutServiceImpl) create(c *gin.Context, req *salemodels.SaleOutUpsert) (*salemodels.SaleOut, error) {
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

func (s *SaleOutServiceImpl) update(c *gin.Context, req *salemodels.SaleOutUpsert) (*salemodels.SaleOut, error) {
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

func (s *SaleOutServiceImpl) Create(c *gin.Context, req *salemodels.SaleOutUpsert) (*salemodels.SaleOut, error) {
	return s.create(c, req)
}

func (s *SaleOutServiceImpl) Update(c *gin.Context, req *salemodels.SaleOutUpsert) (*salemodels.SaleOut, error) {
	return s.update(c, req)
}

func (s *SaleOutServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *SaleOutServiceImpl) GetByID(c *gin.Context, id int64) (*salemodels.SaleOut, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *SaleOutServiceImpl) ListPage(c *gin.Context, req *salemodels.SaleOutQuery) (*erpbiz.PageResult[salemodels.SaleOut], error) {
	if req != nil {
		erpbiz.TrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *SaleOutServiceImpl) validateUniqueFields(c *gin.Context, req *salemodels.SaleOutUpsert, currentID int64) error {
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

func (s *SaleOutServiceImpl) List(c *gin.Context, req *salemodels.SaleOutQuery) (*salemodels.SaleOutListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &salemodels.SaleOutListData{Rows: result.Rows, Total: result.Total}, nil
}
