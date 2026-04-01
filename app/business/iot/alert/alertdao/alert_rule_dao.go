package alertdao

import (
	"nova-factory-server/app/business/iot/alert/alertmodels"

	"github.com/gin-gonic/gin"
)

type AlertRuleDao interface {
	Create(c *gin.Context, data *alertmodels.SetSysAlert) (*alertmodels.SysAlert, error)
	Update(c *gin.Context, data *alertmodels.SetSysAlert) (*alertmodels.SysAlert, error)
	List(c *gin.Context, req *alertmodels.SysAlertListReq) (*alertmodels.SysAlertList, error)
	Remove(c *gin.Context, ids []string) error
	GetByGatewayId(c *gin.Context, gatewayId uint64) (*alertmodels.SysAlert, error)
	Change(c *gin.Context, data *alertmodels.ChangeSysAlert) error
	GetOnlineByGatewayId(c *gin.Context, gatewayId uint64) (*alertmodels.SysAlert, error)
	GetById(c *gin.Context, id uint64) (*alertmodels.SysAlert, error)
	FindOpen(c *gin.Context, gatewayId int64) (*alertmodels.SysAlert, error)
	// Count 调度策略统计
	Count(c *gin.Context) (int64, error)
	// GetByIds 调度策略统计
	GetByIds(c *gin.Context, ids []uint64) ([]*alertmodels.SysAlert, error)
}
