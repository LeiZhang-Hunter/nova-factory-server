package purchaseserviceimpl

import (
	"errors"
	"strings"

	"nova-factory-server/app/business/erp/erpbiz"
	"nova-factory-server/app/business/erp/purchase/purchasedao"
	"nova-factory-server/app/business/erp/purchase/purchasemodels"
	"nova-factory-server/app/business/erp/purchase/purchaseservice"

	"github.com/gin-gonic/gin"
)

// PurchaseReturnServiceImpl 提供业务实现。
type PurchaseReturnServiceImpl struct {
	dao          purchasedao.IPurchaseReturnDao
	uniqueFields []erpbiz.UniqueField
}

// NewPurchaseReturnService 创建服务。
func NewPurchaseReturnService(dao purchasedao.IPurchaseReturnDao) purchaseservice.IPurchaseReturnService {
	return &PurchaseReturnServiceImpl{
		dao:          dao,
		uniqueFields: []erpbiz.UniqueField{{Field: "No", Column: "no", Label: "采购退货单号"}},
	}
}

func (s *PurchaseReturnServiceImpl) create(c *gin.Context, req *purchasemodels.PurchaseReturnUpsert) (*purchasemodels.PurchaseReturn, error) {
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

func (s *PurchaseReturnServiceImpl) update(c *gin.Context, req *purchasemodels.PurchaseReturnUpsert) (*purchasemodels.PurchaseReturn, error) {
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

func (s *PurchaseReturnServiceImpl) Create(c *gin.Context, req *purchasemodels.PurchaseReturnUpsert) (*purchasemodels.PurchaseReturn, error) {
	return s.create(c, req)
}

func (s *PurchaseReturnServiceImpl) Update(c *gin.Context, req *purchasemodels.PurchaseReturnUpsert) (*purchasemodels.PurchaseReturn, error) {
	return s.update(c, req)
}

func (s *PurchaseReturnServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *PurchaseReturnServiceImpl) GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseReturn, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *PurchaseReturnServiceImpl) ListPage(c *gin.Context, req *purchasemodels.PurchaseReturnQuery) (*erpbiz.PageResult[purchasemodels.PurchaseReturn], error) {
	if req != nil {
		erpbiz.TrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *PurchaseReturnServiceImpl) validateUniqueFields(c *gin.Context, req *purchasemodels.PurchaseReturnUpsert, currentID int64) error {
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

func (s *PurchaseReturnServiceImpl) List(c *gin.Context, req *purchasemodels.PurchaseReturnQuery) (*purchasemodels.PurchaseReturnListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseReturnListData{Rows: result.Rows, Total: result.Total}, nil
}
