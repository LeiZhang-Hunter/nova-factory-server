package dashboardServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/dashboard/dashboardDao"
	"nova-factory-server/app/business/dashboard/dashboardModels"
	"nova-factory-server/app/business/dashboard/dashboardService"
)

type DashboardServiceImpl struct {
	dao dashboardDao.DashboardDao
}

func NewDashboardServiceImpl(dao dashboardDao.DashboardDao) dashboardService.DashboardService {
	return &DashboardServiceImpl{
		dao: dao,
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
