package dashboardservice

import (
	"nova-factory-server/app/business/iot/dashboard/dashboardmodels"
	"nova-factory-server/app/business/iot/metric/device/metricmodels"

	"github.com/gin-gonic/gin"
)

type DashboardService interface {
	List(c *gin.Context, req *dashboardmodels.SysDashboardReq) (*dashboardmodels.SysDashboardList, error)
	Set(c *gin.Context, data *dashboardmodels.SetSysDashboard) (*dashboardmodels.SysDashboard, error)
	Remove(c *gin.Context, ids []string) error
	Query(c *gin.Context, req *metricmodels.MetricDataQueryReq) (*metricmodels.MetricQueryData, error)
}
