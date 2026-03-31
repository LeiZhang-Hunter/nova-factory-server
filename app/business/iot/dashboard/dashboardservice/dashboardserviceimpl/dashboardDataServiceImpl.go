package dashboardserviceimpl

import (
	"nova-factory-server/app/business/iot/dashboard/dashboarddao"
	"nova-factory-server/app/business/iot/dashboard/dashboardmodels"
	"nova-factory-server/app/business/iot/dashboard/dashboardservice"

	"github.com/gin-gonic/gin"
)

type DashboardDataServiceImpl struct {
	dao dashboarddao.DashboardDataDao
}

func NewDashboardDataServiceImpl(dao dashboarddao.DashboardDataDao) dashboardservice.DashboardDataService {
	return &DashboardDataServiceImpl{
		dao: dao,
	}
}

func (d *DashboardDataServiceImpl) Set(c *gin.Context, data *dashboardmodels.SetSysDashboardData) (*dashboardmodels.SysDashboardData, error) {
	return d.dao.Set(c, data)
}
func (d *DashboardDataServiceImpl) Remove(c *gin.Context, ids []string) error {
	return d.dao.Remove(c, ids)
}
func (d *DashboardDataServiceImpl) Info(c *gin.Context, dashboardId int64) (*dashboardmodels.SysDashboardData, error) {
	return d.dao.Info(c, dashboardId)
}
