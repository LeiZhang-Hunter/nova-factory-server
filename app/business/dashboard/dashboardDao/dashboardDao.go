package dashboardDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/dashboard/dashboardModels"
)

type DashboardDao interface {
	List(c *gin.Context, req *dashboardModels.SysDashboardReq) (*dashboardModels.SysDashboardList, error)
	Set(c *gin.Context, data *dashboardModels.SetSysDashboard) (*dashboardModels.SysDashboard, error)
	Remove(c *gin.Context, ids []string) error
}
