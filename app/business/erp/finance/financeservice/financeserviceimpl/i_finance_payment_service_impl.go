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

// FinancePaymentServiceImpl 提供业务实现。
type FinancePaymentServiceImpl struct {
	dao          financedao.IFinancePaymentDao
	uniqueFields []erpbiz.UniqueField
}

// NewFinancePaymentService 创建服务。
func NewFinancePaymentService(dao financedao.IFinancePaymentDao) financeservice.IFinancePaymentService {
	return &FinancePaymentServiceImpl{
		dao:          dao,
		uniqueFields: []erpbiz.UniqueField{{Field: "No", Column: "no", Label: "付款单号"}},
	}
}

func (s *FinancePaymentServiceImpl) create(c *gin.Context, req *financemodels.FinancePaymentUpsert) (*financemodels.FinancePayment, error) {
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

func (s *FinancePaymentServiceImpl) update(c *gin.Context, req *financemodels.FinancePaymentUpsert) (*financemodels.FinancePayment, error) {
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

func (s *FinancePaymentServiceImpl) Create(c *gin.Context, req *financemodels.FinancePaymentUpsert) (*financemodels.FinancePayment, error) {
	return s.create(c, req)
}

func (s *FinancePaymentServiceImpl) Update(c *gin.Context, req *financemodels.FinancePaymentUpsert) (*financemodels.FinancePayment, error) {
	return s.update(c, req)
}

func (s *FinancePaymentServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *FinancePaymentServiceImpl) GetByID(c *gin.Context, id int64) (*financemodels.FinancePayment, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *FinancePaymentServiceImpl) ListPage(c *gin.Context, req *financemodels.FinancePaymentQuery) (*erpbiz.PageResult[financemodels.FinancePayment], error) {
	if req != nil {
		erpbiz.TrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *FinancePaymentServiceImpl) validateUniqueFields(c *gin.Context, req *financemodels.FinancePaymentUpsert, currentID int64) error {
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

func (s *FinancePaymentServiceImpl) List(c *gin.Context, req *financemodels.FinancePaymentQuery) (*financemodels.FinancePaymentListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinancePaymentListData{Rows: result.Rows, Total: result.Total}, nil
}
