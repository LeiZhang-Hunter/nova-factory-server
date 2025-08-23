package dashboardDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/dashboard/dashboardModels"
)

type DashboardDataDao interface {
	Set(c *gin.Context, data *dashboardModels.SetSysDashboardData) (*dashboardModels.SysDashboardData, error)
	Remove(c *gin.Context, ids []string) error
}
