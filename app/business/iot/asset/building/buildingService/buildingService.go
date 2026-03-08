package buildingService

import (
	"nova-factory-server/app/business/iot/asset/building/buildingModels"

	"github.com/gin-gonic/gin"
)

type BuildingService interface {
	Set(c *gin.Context, data *buildingModels.SetSysBuilding) (*buildingModels.SysBuilding, error)
	List(c *gin.Context, req *buildingModels.SetSysBuildingListReq) (*buildingModels.SetSysBuildingList, error)
	Remove(c *gin.Context, ids []string) error
	AllDetail(c *gin.Context) ([]*buildingModels.SysBuilding, error)
}
