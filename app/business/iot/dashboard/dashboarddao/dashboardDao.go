package dashboarddao

import (
	"nova-factory-server/app/business/iot/dashboard/dashboardmodels"

	"github.com/gin-gonic/gin"
)

type DashboardDao interface {
	List(c *gin.Context, req *dashboardmodels.SysDashboardReq) (*dashboardmodels.SysDashboardList, error)
	Set(c *gin.Context, data *dashboardmodels.SetSysDashboard) (*dashboardmodels.SysDashboard, error)
	Remove(c *gin.Context, ids []string) error
}
