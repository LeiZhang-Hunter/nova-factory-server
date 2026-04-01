package buildingservice

import (
	"nova-factory-server/app/business/iot/asset/building/buildingmodels"

	"github.com/gin-gonic/gin"
)

type FloorService interface {
	Set(c *gin.Context, data *buildingmodels.SetSysFloor) (*buildingmodels.SysFloor, error)
	List(c *gin.Context, req *buildingmodels.SetSysFloorListReq) (*buildingmodels.SetSysFloorList, error)
	Remove(c *gin.Context, ids []string) error
	CheckUniqueFloor(c *gin.Context, id int64, buildingId int64, level int8) (int64, error)
	SaveLayout(c *gin.Context, id int64, layout *buildingmodels.FloorLayout) error
	Info(c *gin.Context, id int64) (*buildingmodels.SysFloor, error)
}
