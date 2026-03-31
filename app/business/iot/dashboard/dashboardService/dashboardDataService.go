package dashboardService

import (
	"nova-factory-server/app/business/iot/dashboard/dashboardModels"

	"github.com/gin-gonic/gin"
)

type DashboardDataService interface {
	Set(c *gin.Context, data *dashboardModels.SetSysDashboardData) (*dashboardModels.SysDashboardData, error)
	Remove(c *gin.Context, ids []string) error
	Info(c *gin.Context, dashboardId int64) (*dashboardModels.SysDashboardData, error)
}
