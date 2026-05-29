package masterserviceimpl

import (
	"errors"
	"strings"

	"nova-factory-server/app/business/erp/master/masterdao"
	"nova-factory-server/app/business/erp/master/mastermodels"
	"nova-factory-server/app/business/erp/master/masterservice"

	"github.com/gin-gonic/gin"
)

type WarehouseServiceImpl struct {
	dao masterdao.IWarehouseDao
}

func NewWarehouseService(dao masterdao.IWarehouseDao) masterservice.IWarehouseService {
	return &WarehouseServiceImpl{dao: dao}
}

func (s *WarehouseServiceImpl) Create(c *gin.Context, req *mastermodels.WarehouseUpsert) (*mastermodels.Warehouse, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	trimWarehouseUpsert(req)
	if err := validateWarehouseRequired(req); err != nil {
		return nil, err
	}
	if err := s.validateUnique(c, req, 0); err != nil {
		return nil, err
	}
	return s.dao.Create(c, req)
}

func (s *WarehouseServiceImpl) Update(c *gin.Context, req *mastermodels.WarehouseUpsert) (*mastermodels.Warehouse, error) {
	if req == nil {
		return nil, errors.New("参数不能为空")
	}
	if req.ID <= 0 {
		return nil, errors.New("id不能为空")
	}
	trimWarehouseUpsert(req)
	if err := validateWarehouseRequired(req); err != nil {
		return nil, err
	}
	if err := s.validateUnique(c, req, req.ID); err != nil {
		return nil, err
	}
	return s.dao.Update(c, req)
}

func (s *WarehouseServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.dao.DeleteByIDs(c, ids)
}

func (s *WarehouseServiceImpl) GetByID(c *gin.Context, id int64) (*mastermodels.Warehouse, error) {
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	return s.dao.GetByID(c, id)
}

func (s *WarehouseServiceImpl) GetByIDs(c *gin.Context, ids []int64) ([]*mastermodels.Warehouse, error) {
	if len(ids) == 0 {
		return []*mastermodels.Warehouse{}, nil
	}
	return s.dao.GetByIDs(c, ids)
}

func (s *WarehouseServiceImpl) ListPage(c *gin.Context, req *mastermodels.WarehouseQuery) (*mastermodels.WarehouseListData, error) {
	if req != nil {
		req.Name = strings.TrimSpace(req.Name)
	}
	return s.dao.ListPage(c, req)
}

func (s *WarehouseServiceImpl) List(c *gin.Context, req *mastermodels.WarehouseQuery) (*mastermodels.WarehouseListData, error) {
	result, err := s.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &mastermodels.WarehouseListData{Rows: result.Rows, Total: result.Total}, nil
}

func (s *WarehouseServiceImpl) validateUnique(c *gin.Context, req *mastermodels.WarehouseUpsert, currentID int64) error {
	code := strings.TrimSpace(req.Code)
	if code != "" {
		exists, err := s.dao.GetByColumn(c, "code", code)
		if err != nil {
			return err
		}
		if exists != nil && exists.ID != currentID {
			return errors.New("仓库编号已存在")
		}
	}
	name := strings.TrimSpace(req.Name)
	if name != "" {
		exists, err := s.dao.GetByColumn(c, "name", name)
		if err != nil {
			return err
		}
		if exists != nil && exists.ID != currentID {
			return errors.New("仓库名称已存在")
		}
	}
	return nil
}

func trimWarehouseUpsert(req *mastermodels.WarehouseUpsert) {
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	req.Address = strings.TrimSpace(req.Address)
	req.Remark = strings.TrimSpace(req.Remark)
	req.Principal = strings.TrimSpace(req.Principal)
}

func validateWarehouseRequired(req *mastermodels.WarehouseUpsert) error {
	if strings.TrimSpace(req.Code) == "" {
		return errors.New("仓库编号不能为空")
	}
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("仓库名称不能为空")
	}
	return nil
}
