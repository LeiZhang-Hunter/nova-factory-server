package dashboardDao

import (
	"nova-factory-server/app/business/iot/dashboard/dashboardModels"

	"github.com/gin-gonic/gin"
)

type DashboardDao interface {
	List(c *gin.Context, req *dashboardModels.SysDashboardReq) (*dashboardModels.SysDashboardList, error)
	Set(c *gin.Context, data *dashboardModels.SetSysDashboard) (*dashboardModels.SysDashboard, error)
	Remove(c *gin.Context, ids []string) error
}
