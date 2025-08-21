package dashboardService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/dashboard/dashboardModels"
	"nova-factory-server/app/business/metric/device/metricModels"
)

type DashboardService interface {
	List(c *gin.Context, req *dashboardModels.SysDashboardReq) (*dashboardModels.SysDashboardList, error)
	Set(c *gin.Context, data *dashboardModels.SetSysDashboard) (*dashboardModels.SysDashboard, error)
	Remove(c *gin.Context, ids []string) error
	Query(c *gin.Context, req *metricModels.MetricDataQueryReq) (*metricModels.MetricQueryData, error)
}
