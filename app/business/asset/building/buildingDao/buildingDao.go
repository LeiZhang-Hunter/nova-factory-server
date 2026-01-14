package buildingDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/building/buildingModels"
)

type BuildingDao interface {
	Set(c *gin.Context, data *buildingModels.SetSysBuilding) (*buildingModels.SysBuilding, error)
	List(c *gin.Context, req *buildingModels.SetSysBuildingListReq) (*buildingModels.SetSysBuildingList, error)
	Remove(c *gin.Context, ids []string) error
	GetByIds(c *gin.Context, ids []uint64) ([]*buildingModels.SysBuilding, error)
	All(c *gin.Context) ([]*buildingModels.SysBuilding, error)
	Count(c *gin.Context) (int64, error)
}
