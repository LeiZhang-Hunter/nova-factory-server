package settingdao

import (
	"nova-factory-server/app/business/erp/setting/settingmodels"

	"github.com/gin-gonic/gin"
)

// ILogisticsCompanyDao ERP物流公司数据访问接口
type ILogisticsCompanyDao interface {
	Create(c *gin.Context, req *settingmodels.LogisticsCompanyUpsert) (*settingmodels.LogisticsCompany, error)
	Update(c *gin.Context, req *settingmodels.LogisticsCompanyUpsert) (*settingmodels.LogisticsCompany, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*settingmodels.LogisticsCompany, error)
	List(c *gin.Context, req *settingmodels.LogisticsCompanyQuery) (*settingmodels.LogisticsCompanyListData, error)
}
