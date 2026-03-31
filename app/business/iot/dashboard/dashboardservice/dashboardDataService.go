package dashboardservice

import (
	"nova-factory-server/app/business/iot/dashboard/dashboardmodels"

	"github.com/gin-gonic/gin"
)

type DashboardDataService interface {
	Set(c *gin.Context, data *dashboardmodels.SetSysDashboardData) (*dashboardmodels.SysDashboardData, error)
	Remove(c *gin.Context, ids []string) error
	Info(c *gin.Context, dashboardId int64) (*dashboardmodels.SysDashboardData, error)
}
