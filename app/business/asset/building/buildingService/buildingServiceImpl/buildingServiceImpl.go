package buildingServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/building/buildingDao"
	"nova-factory-server/app/business/asset/building/buildingModels"
	"nova-factory-server/app/business/asset/building/buildingService"
)

type BuildingServiceImpl struct {
	dao buildingDao.BuildingDao
}

func NewBuildingServiceImpl(dao buildingDao.BuildingDao) buildingService.BuildingService {
	return &BuildingServiceImpl{
		dao: dao,
	}
}

func (b *BuildingServiceImpl) Set(c *gin.Context, data *buildingModels.SetSysBuilding) (*buildingModels.SysBuilding, error) {
	return b.dao.Set(c, data)
}

func (b *BuildingServiceImpl) List(c *gin.Context, req *buildingModels.SetSysBuildingListReq) (*buildingModels.SetSysBuildingList, error) {
	return b.dao.List(c, req)
}
func (b *BuildingServiceImpl) Remove(c *gin.Context, ids []string) error {
	return b.dao.Remove(c, ids)
}
