package settingservice

import (
	"nova-factory-server/app/business/erp/setting/settingmodels"

	"github.com/gin-gonic/gin"
)

// ILogisticsCompanyService ERP物流公司服务接口
type ILogisticsCompanyService interface {
	Create(c *gin.Context, req *settingmodels.LogisticsCompanyUpsert) (*settingmodels.LogisticsCompany, error)
	Update(c *gin.Context, req *settingmodels.LogisticsCompanyUpsert) (*settingmodels.LogisticsCompany, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*settingmodels.LogisticsCompany, error)
	List(c *gin.Context, req *settingmodels.LogisticsCompanyQuery) (*settingmodels.LogisticsCompanyListData, error)
}
