package dashboardService

import (
	"nova-factory-server/app/business/iot/dashboard/dashboardModels"
	"nova-factory-server/app/business/iot/metric/device/metricModels"

	"github.com/gin-gonic/gin"
)

type DashboardService interface {
	List(c *gin.Context, req *dashboardModels.SysDashboardReq) (*dashboardModels.SysDashboardList, error)
	Set(c *gin.Context, data *dashboardModels.SetSysDashboard) (*dashboardModels.SysDashboard, error)
	Remove(c *gin.Context, ids []string) error
	Query(c *gin.Context, req *metricModels.MetricDataQueryReq) (*metricModels.MetricQueryData, error)
}
