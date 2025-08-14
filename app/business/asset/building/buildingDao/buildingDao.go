package buildingDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/building/buildingModels"
)

type BuildingDao interface {
	Set(c *gin.Context, data *buildingModels.SetSysBuilding) (*buildingModels.SysBuilding, error)
	List(c *gin.Context, req *buildingModels.SetSysBuildingListReq) (*buildingModels.SetSysBuildingList, error)
	Remove(c *gin.Context, ids []string) error
}
