package financeserviceimpl

import (
	"errors"
	"strings"

	"nova-factory-server/app/business/erp/erpbiz"
	"nova-factory-server/app/business/erp/finance/financedao"
	"nova-factory-server/app/business/erp/finance/financemodels"
	"nova-factory-server/app/business/erp/finance/financeservice"

	"github.com/gin-gonic/gin"
)

// FinancePaymentItemServiceImpl 提供业务实现。
type FinancePaymentItemServiceImpl struct {
	dao          financedao.IFinancePaymentItemDao
	uniqueFields []erpbiz.UniqueField
}

// NewFinancePaymentItemService 创建服务。
func NewFinancePaymentItemService(dao financedao.IFinancePaymentItemDao) financeservice.IFinancePaymentItemService {
	return &FinancePaymentItemServiceImpl{
		dao:          dao,
		uniqueFields: []erpbiz.UniqueField{},
	}
}

func (s *FinancePaymentItemServiceImpl) create(c *gin.Context, req *financemodels.FinancePaymentItemUpsert) (*financemodels.FinancePaymentItem, error) {
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

func (s *FinancePaymentItemServiceImpl) update(c *gin.Context, req *financemodels.FinancePaymentItemUpsert) (*financemodels.FinancePaymentItem, error) {
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

func (s *FinancePaymentItemServiceImpl) Create(c *gin.Context, req *financemodels.FinancePaymentItemUpsert) (*financemodels.FinancePaymentItem, error) {
	return s.create(c, req)
}

func (s *FinancePaymentItemServiceImpl) Update(c *gin.Context, req *financemodels.FinancePaymentItemUpsert) (*financemodels.FinancePaymentItem, error) {
	return s.update(c, req)
}

func (s *FinancePaymentItemServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *FinancePaymentItemServiceImpl) GetByID(c *gin.Context, id int64) (*financemodels.FinancePaymentItem, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *FinancePaymentItemServiceImpl) ListPage(c *gin.Context, req *financemodels.FinancePaymentItemQuery) (*erpbiz.PageResult[financemodels.FinancePaymentItem], error) {
	if req != nil {
		erpbiz.TrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *FinancePaymentItemServiceImpl) validateUniqueFields(c *gin.Context, req *financemodels.FinancePaymentItemUpsert, currentID int64) error {
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

func (s *FinancePaymentItemServiceImpl) List(c *gin.Context, req *financemodels.FinancePaymentItemQuery) (*financemodels.FinancePaymentItemListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinancePaymentItemListData{Rows: result.Rows, Total: result.Total}, nil
}
