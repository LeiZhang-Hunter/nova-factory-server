package dashboardServiceImpl

import (
	"nova-factory-server/app/business/iot/dashboard/dashboardDao"
	"nova-factory-server/app/business/iot/dashboard/dashboardModels"
	"nova-factory-server/app/business/iot/dashboard/dashboardService"
	"nova-factory-server/app/business/iot/metric/device/metricDao"
	"nova-factory-server/app/business/iot/metric/device/metricModels"

	"github.com/gin-gonic/gin"
)

type DashboardServiceImpl struct {
	dao        dashboardDao.DashboardDao
	metricCDao metricDao.IMetricDao
}

func NewDashboardServiceImpl(dao dashboardDao.DashboardDao, metricCDao metricDao.IMetricDao) dashboardService.DashboardService {
	return &DashboardServiceImpl{
		dao:        dao,
		metricCDao: metricCDao,
	}
}

func (d *DashboardServiceImpl) List(c *gin.Context, req *dashboardModels.SysDashboardReq) (*dashboardModels.SysDashboardList, error) {
	return d.dao.List(c, req)
}
func (d *DashboardServiceImpl) Set(c *gin.Context, data *dashboardModels.SetSysDashboard) (*dashboardModels.SysDashboard, error) {
	return d.dao.Set(c, data)
}
func (d *DashboardServiceImpl) Remove(c *gin.Context, ids []string) error {
	return d.dao.Remove(c, ids)
}

func (d *DashboardServiceImpl) Query(c *gin.Context, req *metricModels.MetricDataQueryReq) (*metricModels.MetricQueryData, error) {
	return d.metricCDao.Query(c, req)
}
