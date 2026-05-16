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

// PurchaseOrderItemServiceImpl 提供业务实现。
type PurchaseOrderItemServiceImpl struct {
	dao          purchasedao.IPurchaseOrderItemDao
	uniqueFields []erpbiz.UniqueField
}

// NewPurchaseOrderItemService 创建服务。
func NewPurchaseOrderItemService(dao purchasedao.IPurchaseOrderItemDao) purchaseservice.IPurchaseOrderItemService {
	return &PurchaseOrderItemServiceImpl{
		dao:          dao,
		uniqueFields: []erpbiz.UniqueField{},
	}
}

func (s *PurchaseOrderItemServiceImpl) create(c *gin.Context, req *purchasemodels.PurchaseOrderItemUpsert) (*purchasemodels.PurchaseOrderItem, error) {
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

func (s *PurchaseOrderItemServiceImpl) update(c *gin.Context, req *purchasemodels.PurchaseOrderItemUpsert) (*purchasemodels.PurchaseOrderItem, error) {
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

func (s *PurchaseOrderItemServiceImpl) Create(c *gin.Context, req *purchasemodels.PurchaseOrderItemUpsert) (*purchasemodels.PurchaseOrderItem, error) {
	return s.create(c, req)
}

func (s *PurchaseOrderItemServiceImpl) Update(c *gin.Context, req *purchasemodels.PurchaseOrderItemUpsert) (*purchasemodels.PurchaseOrderItem, error) {
	return s.update(c, req)
}

func (s *PurchaseOrderItemServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *PurchaseOrderItemServiceImpl) GetByID(c *gin.Context, id int64) (*purchasemodels.PurchaseOrderItem, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *PurchaseOrderItemServiceImpl) ListPage(c *gin.Context, req *purchasemodels.PurchaseOrderItemQuery) (*erpbiz.PageResult[purchasemodels.PurchaseOrderItem], error) {
	if req != nil {
		erpbiz.TrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *PurchaseOrderItemServiceImpl) validateUniqueFields(c *gin.Context, req *purchasemodels.PurchaseOrderItemUpsert, currentID int64) error {
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

func (s *PurchaseOrderItemServiceImpl) List(c *gin.Context, req *purchasemodels.PurchaseOrderItemQuery) (*purchasemodels.PurchaseOrderItemListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &purchasemodels.PurchaseOrderItemListData{Rows: result.Rows, Total: result.Total}, nil
}
