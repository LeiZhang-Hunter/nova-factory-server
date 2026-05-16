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

// SaleReturnServiceImpl 提供业务实现。
type SaleReturnServiceImpl struct {
	dao          saledao.ISaleReturnDao
	uniqueFields []erpbiz.UniqueField
}

// NewSaleReturnService 创建服务。
func NewSaleReturnService(dao saledao.ISaleReturnDao) saleservice.ISaleReturnService {
	return &SaleReturnServiceImpl{
		dao:          dao,
		uniqueFields: []erpbiz.UniqueField{{Field: "No", Column: "no", Label: "销售退货单号"}},
	}
}

func (s *SaleReturnServiceImpl) create(c *gin.Context, req *salemodels.SaleReturnUpsert) (*salemodels.SaleReturn, error) {
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

func (s *SaleReturnServiceImpl) update(c *gin.Context, req *salemodels.SaleReturnUpsert) (*salemodels.SaleReturn, error) {
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

func (s *SaleReturnServiceImpl) Create(c *gin.Context, req *salemodels.SaleReturnUpsert) (*salemodels.SaleReturn, error) {
	return s.create(c, req)
}

func (s *SaleReturnServiceImpl) Update(c *gin.Context, req *salemodels.SaleReturnUpsert) (*salemodels.SaleReturn, error) {
	return s.update(c, req)
}

func (s *SaleReturnServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *SaleReturnServiceImpl) GetByID(c *gin.Context, id int64) (*salemodels.SaleReturn, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *SaleReturnServiceImpl) ListPage(c *gin.Context, req *salemodels.SaleReturnQuery) (*erpbiz.PageResult[salemodels.SaleReturn], error) {
	if req != nil {
		erpbiz.TrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *SaleReturnServiceImpl) validateUniqueFields(c *gin.Context, req *salemodels.SaleReturnUpsert, currentID int64) error {
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

func (s *SaleReturnServiceImpl) List(c *gin.Context, req *salemodels.SaleReturnQuery) (*salemodels.SaleReturnListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &salemodels.SaleReturnListData{Rows: result.Rows, Total: result.Total}, nil
}
