package dashboardServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/dashboard/dashboardDao"
	"nova-factory-server/app/business/dashboard/dashboardModels"
	"nova-factory-server/app/business/dashboard/dashboardService"
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
