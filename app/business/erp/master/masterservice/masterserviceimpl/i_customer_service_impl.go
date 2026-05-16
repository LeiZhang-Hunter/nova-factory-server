package masterserviceimpl

import (
	"errors"
	"strings"

	"nova-factory-server/app/business/erp/erpbiz"
	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"

	"github.com/gin-gonic/gin"
)

// CustomerServiceImpl 提供业务实现。
type CustomerServiceImpl struct {
	dao          masterdao.ICustomerDao
	uniqueFields []erpbiz.UniqueField
}

// NewCustomerService 创建服务。
func NewCustomerService(dao masterdao.ICustomerDao) masterservice.ICustomerService {
	return &CustomerServiceImpl{
		dao:          dao,
		uniqueFields: []erpbiz.UniqueField{{Field: "Code", Column: "code", Label: "客户编码，对接 erp_order.b_type_code"}},
	}
}

func (s *CustomerServiceImpl) create(c *gin.Context, req *mastermodels.CustomerUpsert) (*mastermodels.Customer, error) {
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

func (s *CustomerServiceImpl) update(c *gin.Context, req *mastermodels.CustomerUpsert) (*mastermodels.Customer, error) {
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

func (s *CustomerServiceImpl) Create(c *gin.Context, req *mastermodels.CustomerUpsert) (*mastermodels.Customer, error) {
	return s.create(c, req)
}

func (s *CustomerServiceImpl) Update(c *gin.Context, req *mastermodels.CustomerUpsert) (*mastermodels.Customer, error) {
	return s.update(c, req)
}

func (s *CustomerServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *CustomerServiceImpl) GetByID(c *gin.Context, id int64) (*mastermodels.Customer, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *CustomerServiceImpl) ListPage(c *gin.Context, req *mastermodels.CustomerQuery) (*erpbiz.PageResult[mastermodels.Customer], error) {
	if req != nil {
		erpbiz.TrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *CustomerServiceImpl) validateUniqueFields(c *gin.Context, req *mastermodels.CustomerUpsert, currentID int64) error {
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

func (s *CustomerServiceImpl) List(c *gin.Context, req *mastermodels.CustomerQuery) (*mastermodels.CustomerListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.CustomerListData{Rows: result.Rows, Total: result.Total}, nil
}
