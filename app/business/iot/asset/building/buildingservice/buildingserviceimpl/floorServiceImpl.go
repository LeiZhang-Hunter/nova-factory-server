package buildingserviceimpl

import (
	"nova-factory-server/app/business/iot/asset/building/buildingdao"
	"nova-factory-server/app/business/iot/asset/building/buildingmodels"
	"nova-factory-server/app/business/iot/asset/building/buildingservice"

	"github.com/gin-gonic/gin"
)

type FloorServiceImpl struct {
	dao buildingdao.FloorDao
}

func NewFloorServiceImpl(dao buildingdao.FloorDao) buildingservice.FloorService {
	return &FloorServiceImpl{
		dao: dao,
	}
}

func (b *FloorServiceImpl) Set(c *gin.Context, data *buildingmodels.SetSysFloor) (*buildingmodels.SysFloor, error) {
	return b.dao.Set(c, data)
}

func (b *FloorServiceImpl) List(c *gin.Context, req *buildingmodels.SetSysFloorListReq) (*buildingmodels.SetSysFloorList, error) {
	return b.dao.List(c, req)
}
func (b *FloorServiceImpl) Remove(c *gin.Context, ids []string) error {
	return b.dao.Remove(c, ids)
}

func (b *FloorServiceImpl) CheckUniqueFloor(c *gin.Context, id int64, buildingId int64, level int8) (int64, error) {
	return b.dao.CheckUniqueFloor(c, id, buildingId, level)
}
func (b *FloorServiceImpl) SaveLayout(c *gin.Context, id int64, layout *buildingmodels.FloorLayout) error {
	return b.dao.SaveLayout(c, id, layout)
}

func (b *FloorServiceImpl) Info(c *gin.Context, id int64) (*buildingmodels.SysFloor, error) {
	return b.dao.Info(c, id)
}
