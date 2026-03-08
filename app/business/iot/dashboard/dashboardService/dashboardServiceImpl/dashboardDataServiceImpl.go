package dashboardServiceImpl

import (
	"nova-factory-server/app/business/iot/dashboard/dashboardDao"
	"nova-factory-server/app/business/iot/dashboard/dashboardModels"
	"nova-factory-server/app/business/iot/dashboard/dashboardService"

	"github.com/gin-gonic/gin"
)

type DashboardDataServiceImpl struct {
	dao dashboardDao.DashboardDataDao
}

func NewDashboardDataServiceImpl(dao dashboardDao.DashboardDataDao) dashboardService.DashboardDataService {
	return &DashboardDataServiceImpl{
		dao: dao,
	}
}

func (d *DashboardDataServiceImpl) Set(c *gin.Context, data *dashboardModels.SetSysDashboardData) (*dashboardModels.SysDashboardData, error) {
	return d.dao.Set(c, data)
}
func (d *DashboardDataServiceImpl) Remove(c *gin.Context, ids []string) error {
	return d.dao.Remove(c, ids)
}
func (d *DashboardDataServiceImpl) Info(c *gin.Context, dashboardId int64) (*dashboardModels.SysDashboardData, error) {
	return d.dao.Info(c, dashboardId)
}
