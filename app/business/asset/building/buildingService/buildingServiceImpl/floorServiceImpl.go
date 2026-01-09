package buildingServiceImpl

import (
	"nova-factory-server/app/business/asset/building/buildingDao"
	"nova-factory-server/app/business/asset/building/buildingModels"
	"nova-factory-server/app/business/asset/building/buildingService"

	"github.com/gin-gonic/gin"
)

type FloorServiceImpl struct {
	dao buildingDao.FloorDao
}

func NewFloorServiceImpl(dao buildingDao.FloorDao) buildingService.FloorService {
	return &FloorServiceImpl{
		dao: dao,
	}
}

func (b *FloorServiceImpl) Set(c *gin.Context, data *buildingModels.SetSysFloor) (*buildingModels.SysFloor, error) {
	return b.dao.Set(c, data)
}

func (b *FloorServiceImpl) List(c *gin.Context, req *buildingModels.SetSysFloorListReq) (*buildingModels.SetSysFloorList, error) {
	return b.dao.List(c, req)
}
func (b *FloorServiceImpl) Remove(c *gin.Context, ids []string) error {
	return b.dao.Remove(c, ids)
}

func (b *FloorServiceImpl) CheckUniqueFloor(c *gin.Context, id int64, buildingId int64, level int8) (int64, error) {
	return b.dao.CheckUniqueFloor(c, id, buildingId, level)
}
