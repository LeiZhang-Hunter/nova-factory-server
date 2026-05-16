package stockserviceimpl

import (
	"errors"
	"strings"

	"nova-factory-server/app/business/erp/erpbiz"
	"nova-factory-server/app/business/erp/stock/stockdao"
	"nova-factory-server/app/business/erp/stock/stockmodels"
	"nova-factory-server/app/business/erp/stock/stockservice"

	"github.com/gin-gonic/gin"
)

// StockOutServiceImpl 提供业务实现。
type StockOutServiceImpl struct {
	dao          stockdao.IStockOutDao
	uniqueFields []erpbiz.UniqueField
}

// NewStockOutService 创建服务。
func NewStockOutService(dao stockdao.IStockOutDao) stockservice.IStockOutService {
	return &StockOutServiceImpl{
		dao:          dao,
		uniqueFields: []erpbiz.UniqueField{{Field: "No", Column: "no", Label: "出库单号"}},
	}
}

func (s *StockOutServiceImpl) create(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error) {
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

func (s *StockOutServiceImpl) update(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error) {
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

func (s *StockOutServiceImpl) Create(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error) {
	return s.create(c, req)
}

func (s *StockOutServiceImpl) Update(c *gin.Context, req *stockmodels.StockOutUpsert) (*stockmodels.StockOut, error) {
	return s.update(c, req)
}

func (s *StockOutServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *StockOutServiceImpl) GetByID(c *gin.Context, id int64) (*stockmodels.StockOut, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *StockOutServiceImpl) ListPage(c *gin.Context, req *stockmodels.StockOutQuery) (*erpbiz.PageResult[stockmodels.StockOut], error) {
	if req != nil {
		erpbiz.TrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *StockOutServiceImpl) validateUniqueFields(c *gin.Context, req *stockmodels.StockOutUpsert, currentID int64) error {
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

func (s *StockOutServiceImpl) List(c *gin.Context, req *stockmodels.StockOutQuery) (*stockmodels.StockOutListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockOutListData{Rows: result.Rows, Total: result.Total}, nil
}
