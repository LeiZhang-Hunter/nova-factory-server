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

// FinanceReceiptServiceImpl 提供业务实现。
type FinanceReceiptServiceImpl struct {
	dao          financedao.IFinanceReceiptDao
	uniqueFields []erpbiz.UniqueField
}

// NewFinanceReceiptService 创建服务。
func NewFinanceReceiptService(dao financedao.IFinanceReceiptDao) financeservice.IFinanceReceiptService {
	return &FinanceReceiptServiceImpl{
		dao:          dao,
		uniqueFields: []erpbiz.UniqueField{{Field: "No", Column: "no", Label: "收款单号"}},
	}
}

func (s *FinanceReceiptServiceImpl) create(c *gin.Context, req *financemodels.FinanceReceiptUpsert) (*financemodels.FinanceReceipt, error) {
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

func (s *FinanceReceiptServiceImpl) update(c *gin.Context, req *financemodels.FinanceReceiptUpsert) (*financemodels.FinanceReceipt, error) {
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

func (s *FinanceReceiptServiceImpl) Create(c *gin.Context, req *financemodels.FinanceReceiptUpsert) (*financemodels.FinanceReceipt, error) {
	return s.create(c, req)
}

func (s *FinanceReceiptServiceImpl) Update(c *gin.Context, req *financemodels.FinanceReceiptUpsert) (*financemodels.FinanceReceipt, error) {
	return s.update(c, req)
}

func (s *FinanceReceiptServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *FinanceReceiptServiceImpl) GetByID(c *gin.Context, id int64) (*financemodels.FinanceReceipt, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *FinanceReceiptServiceImpl) ListPage(c *gin.Context, req *financemodels.FinanceReceiptQuery) (*erpbiz.PageResult[financemodels.FinanceReceipt], error) {
	if req != nil {
		erpbiz.TrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *FinanceReceiptServiceImpl) validateUniqueFields(c *gin.Context, req *financemodels.FinanceReceiptUpsert, currentID int64) error {
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

func (s *FinanceReceiptServiceImpl) List(c *gin.Context, req *financemodels.FinanceReceiptQuery) (*financemodels.FinanceReceiptListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &financemodels.FinanceReceiptListData{Rows: result.Rows, Total: result.Total}, nil
}
