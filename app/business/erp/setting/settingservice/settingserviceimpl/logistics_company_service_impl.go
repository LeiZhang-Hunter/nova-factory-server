package settingserviceimpl

import (
	"nova-factory-server/app/business/erp/setting/settingdao"
	"nova-factory-server/app/business/erp/setting/settingmodels"
	"nova-factory-server/app/business/erp/setting/settingservice"

	"github.com/gin-gonic/gin"
)

// LogisticsCompanyServiceImpl 提供 ERP 物流公司业务实现。
type LogisticsCompanyServiceImpl struct {
	dao settingdao.ILogisticsCompanyDao
}

// NewLogisticsCompanyService 创建 ERP 物流公司服务。
func NewLogisticsCompanyService(dao settingdao.ILogisticsCompanyDao) settingservice.ILogisticsCompanyService {
	return &LogisticsCompanyServiceImpl{dao: dao}
}

// Create 新增 ERP 物流公司。
func (l *LogisticsCompanyServiceImpl) Create(c *gin.Context, req *settingmodels.LogisticsCompanyUpsert) (*settingmodels.LogisticsCompany, error) {
	return l.dao.Create(c, req)
}

// Update 修改 ERP 物流公司。
func (l *LogisticsCompanyServiceImpl) Update(c *gin.Context, req *settingmodels.LogisticsCompanyUpsert) (*settingmodels.LogisticsCompany, error) {
	return l.dao.Update(c, req)
}

// DeleteByIDs 删除 ERP 物流公司。
func (l *LogisticsCompanyServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return l.dao.DeleteByIDs(c, ids)
}

// GetByID 查询 ERP 物流公司详情。
func (l *LogisticsCompanyServiceImpl) GetByID(c *gin.Context, id int64) (*settingmodels.LogisticsCompany, error) {
	return l.dao.GetByID(c, id)
}

// List 分页查询 ERP 物流公司。
func (l *LogisticsCompanyServiceImpl) List(c *gin.Context, req *settingmodels.LogisticsCompanyQuery) (*settingmodels.LogisticsCompanyListData, error) {
	return l.dao.List(c, req)
}
