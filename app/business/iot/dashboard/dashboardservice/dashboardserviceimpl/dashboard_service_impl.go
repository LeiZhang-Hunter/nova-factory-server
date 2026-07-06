package dashboardserviceimpl

import (
	"nova-factory-server/app/business/iot/dashboard/dashboarddao"
	"nova-factory-server/app/business/iot/dashboard/dashboardmodels"
	"nova-factory-server/app/business/iot/dashboard/dashboardservice"
	"nova-factory-server/app/business/iot/devicemonitor/devicemonitordao"
	"nova-factory-server/app/business/iot/metric/device/metricdao"
	"nova-factory-server/app/business/iot/metric/device/metricmodels"

	"github.com/gin-gonic/gin"
)

type DashboardServiceImpl struct {
	dao        dashboarddao.DashboardDao
	devMaoDao  devicemonitordao.IDeviceDataReportDao
	metricCDao metricdao.IMetricDao
}

func NewDashboardServiceImpl(dao dashboarddao.DashboardDao,
	metricCDao metricdao.IMetricDao, devMaoDao devicemonitordao.IDeviceDataReportDao) dashboardservice.DashboardService {
	return &DashboardServiceImpl{
		dao:        dao,
		metricCDao: metricCDao,
		devMaoDao:  devMaoDao,
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
	dev, _ := d.devMaoDao.GetByDev(c, req.Name)
	if dev != nil {
		if req.QueryMetric == nil {
			req.QueryMetric = make([]*metricmodels.MetricQueryCondition, 0)
		}
		req.QueryMetric = append(req.QueryMetric, &metricmodels.MetricQueryCondition{
			DeviceId:   dev.DeviceID,
			TemplateId: dev.TemplateID,
			DataId:     dev.DataID,
		})
	}
	return d.metricCDao.Query(c, req)
}
