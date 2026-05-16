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

// AccountServiceImpl 提供业务实现。
type AccountServiceImpl struct {
	dao          masterdao.IAccountDao
	uniqueFields []erpbiz.UniqueField
}

// NewAccountService 创建服务。
func NewAccountService(dao masterdao.IAccountDao) masterservice.IAccountService {
	return &AccountServiceImpl{
		dao:          dao,
		uniqueFields: []erpbiz.UniqueField{{Field: "No", Column: "no", Label: "账户编码，对接 erp_order_account.finance_code"}},
	}
}

func (s *AccountServiceImpl) create(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error) {
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

func (s *AccountServiceImpl) update(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error) {
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

func (s *AccountServiceImpl) Create(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error) {
	return s.create(c, req)
}

func (s *AccountServiceImpl) Update(c *gin.Context, req *mastermodels.AccountUpsert) (*mastermodels.Account, error) {
	return s.update(c, req)
}

func (s *AccountServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *AccountServiceImpl) GetByID(c *gin.Context, id int64) (*mastermodels.Account, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *AccountServiceImpl) ListPage(c *gin.Context, req *mastermodels.AccountQuery) (*erpbiz.PageResult[mastermodels.Account], error) {
	if req != nil {
		erpbiz.TrimStringFields(req)
	}
	return s.dao.ListPage(c, req)
}

func (s *AccountServiceImpl) validateUniqueFields(c *gin.Context, req *mastermodels.AccountUpsert, currentID int64) error {
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

func (s *AccountServiceImpl) List(c *gin.Context, req *mastermodels.AccountQuery) (*mastermodels.AccountListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.AccountListData{Rows: result.Rows, Total: result.Total}, nil
}
