package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/config/dao"
	"nova-factory-server/app/business/shop/config/models"
	"nova-factory-server/app/business/shop/config/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// LogisticsCompanyServiceImpl 提供 ERP 物流公司业务实现。
type LogisticsCompanyServiceImpl struct {
	dao dao.ILogisticsCompanyDao
}

// NewLogisticsCompanyService 创建 ERP 物流公司服务。
func NewLogisticsCompanyService(dao dao.ILogisticsCompanyDao) service.ILogisticsCompanyService {
	return &LogisticsCompanyServiceImpl{dao: dao}
}

// Create 新增 ERP 物流公司。
func (l *LogisticsCompanyServiceImpl) Create(c *gin.Context, req *models.LogisticsCompanyUpsert) (*models.LogisticsCompany, error) {
	if err := l.validateUniqueFields(c, req); err != nil {
		return nil, err
	}
	return l.dao.Create(c, req)
}

// Update 修改 ERP 物流公司。
func (l *LogisticsCompanyServiceImpl) Update(c *gin.Context, req *models.LogisticsCompanyUpsert) (*models.LogisticsCompany, error) {
	if err := l.validateUniqueFields(c, req); err != nil {
		return nil, err
	}
	return l.dao.Update(c, req)
}

// DeleteByIDs 删除 ERP 物流公司。
func (l *LogisticsCompanyServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return l.dao.DeleteByIDs(c, ids)
}

// GetByID 查询 ERP 物流公司详情。
func (l *LogisticsCompanyServiceImpl) GetByID(c *gin.Context, id int64) (*models.LogisticsCompany, error) {
	return l.dao.GetByID(c, id)
}

// GetByCode 根据编码查询 ERP 物流公司。
func (l *LogisticsCompanyServiceImpl) GetByCode(c *gin.Context, code string) (*models.LogisticsCompany, error) {
	return l.dao.GetByCode(c, code)
}

// GetByName 根据名称查询 ERP 物流公司。
func (l *LogisticsCompanyServiceImpl) GetByName(c *gin.Context, name string) (*models.LogisticsCompany, error) {
	return l.dao.GetByName(c, name)
}

// List 分页查询 ERP 物流公司。
func (l *LogisticsCompanyServiceImpl) List(c *gin.Context, req *models.LogisticsCompanyQuery) (*models.LogisticsCompanyListData, error) {
	return l.dao.List(c, req)
}

func (l *LogisticsCompanyServiceImpl) validateUniqueFields(c *gin.Context, req *models.LogisticsCompanyUpsert) error {
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	if req.Code == "" {
		return errors.New("物流公司编码不能为空")
	}
	if req.Name == "" {
		return errors.New("物流公司名称不能为空")
	}

	codeInfo, err := l.dao.GetByCode(c, req.Code)
	if err != nil {
		return err
	}
	if codeInfo != nil && codeInfo.ID != req.ID {
		return errors.New("物流公司编码已存在")
	}

	nameInfo, err := l.dao.GetByName(c, req.Name)
	if err != nil {
		return err
	}
	if nameInfo != nil && nameInfo.ID != req.ID {
		return errors.New("物流公司名称已存在")
	}
	return nil
}
