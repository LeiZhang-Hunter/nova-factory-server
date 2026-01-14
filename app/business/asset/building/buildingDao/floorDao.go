package buildingDao

import (
	"nova-factory-server/app/business/asset/building/buildingModels"

	"github.com/gin-gonic/gin"
)

type FloorDao interface {
	Set(c *gin.Context, data *buildingModels.SetSysFloor) (*buildingModels.SysFloor, error)
	List(c *gin.Context, req *buildingModels.SetSysFloorListReq) (*buildingModels.SetSysFloorList, error)
	Remove(c *gin.Context, ids []string) error
	GetByIds(c *gin.Context, ids []uint64) ([]*buildingModels.SysFloor, error)
	CheckUniqueFloor(c *gin.Context, id int64, buildingId int64, level int8) (int64, error)
	SaveLayout(c *gin.Context, id int64, layout *buildingModels.FloorLayout) error
	Info(c *gin.Context, id int64) (*buildingModels.SysFloor, error)
	All(c *gin.Context) ([]*buildingModels.SysFloor, error)
}
