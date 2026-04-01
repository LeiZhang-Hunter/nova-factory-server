package dashboardserviceimpl

import (
	"nova-factory-server/app/business/iot/dashboard/dashboarddao"
	"nova-factory-server/app/business/iot/dashboard/dashboardmodels"
	"nova-factory-server/app/business/iot/dashboard/dashboardservice"
	"nova-factory-server/app/business/iot/metric/device/metricdao"
	"nova-factory-server/app/business/iot/metric/device/metricmodels"

	"github.com/gin-gonic/gin"
)

type DashboardServiceImpl struct {
	dao        dashboarddao.DashboardDao
	metricCDao metricdao.IMetricDao
}

func NewDashboardServiceImpl(dao dashboarddao.DashboardDao, metricCDao metricdao.IMetricDao) dashboardservice.DashboardService {
	return &DashboardServiceImpl{
		dao:        dao,
		metricCDao: metricCDao,
	}
}

func (d *DashboardServiceImpl) List(c *gin.Context, req *dashboardmodels.SysDashboardReq) (*dashboardmodels.SysDashboardList, error) {
	return d.dao.List(c, req)
}
func (d *DashboardServiceImpl) Set(c *gin.Context, data *dashboardmodels.SetSysDashboard) (*dashboardmodels.SysDashboard, error) {
	return d.dao.Set(c, data)
}
func (d *DashboardServiceImpl) Remove(c *gin.Context, ids []string) error {
	return d.dao.Remove(c, ids)
}

func (d *DashboardServiceImpl) Query(c *gin.Context, req *metricmodels.MetricDataQueryReq) (*metricmodels.MetricQueryData, error) {
	return d.metricCDao.Query(c, req)
}
