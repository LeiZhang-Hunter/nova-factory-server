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

// StockMoveServiceImpl 提供业务实现。
type StockMoveServiceImpl struct {
	dao          stockdao.IStockMoveDao
	uniqueFields []erpbiz.UniqueField
}

// NewStockMoveService 创建服务。
func NewStockMoveService(dao stockdao.IStockMoveDao) stockservice.IStockMoveService {
	return &StockMoveServiceImpl{
		dao:          dao,
		uniqueFields: []erpbiz.UniqueField{{Field: "No", Column: "no", Label: "调拨单号"}},
	}
}

func (s *StockMoveServiceImpl) create(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error) {
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

func (s *StockMoveServiceImpl) update(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error) {
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

func (s *StockMoveServiceImpl) Create(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error) {
	return s.create(c, req)
}

func (s *StockMoveServiceImpl) Update(c *gin.Context, req *stockmodels.StockMoveUpsert) (*stockmodels.StockMove, error) {
	return s.update(c, req)
}

func (s *StockMoveServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *StockMoveServiceImpl) GetByID(c *gin.Context, id int64) (*stockmodels.StockMove, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *StockMoveServiceImpl) ListPage(c *gin.Context, req *stockmodels.StockMoveQuery) (*erpbiz.PageResult[stockmodels.StockMove], error) {
	if req != nil {
		erpbiz.TrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *StockMoveServiceImpl) validateUniqueFields(c *gin.Context, req *stockmodels.StockMoveUpsert, currentID int64) error {
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

func (s *StockMoveServiceImpl) List(c *gin.Context, req *stockmodels.StockMoveQuery) (*stockmodels.StockMoveListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &stockmodels.StockMoveListData{Rows: result.Rows, Total: result.Total}, nil
}
