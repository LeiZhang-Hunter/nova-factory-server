package buildingservice

import (
	"nova-factory-server/app/business/iot/asset/building/buildingmodels"

	"github.com/gin-gonic/gin"
)

type BuildingService interface {
	Set(c *gin.Context, data *buildingmodels.SetSysBuilding) (*buildingmodels.SysBuilding, error)
	List(c *gin.Context, req *buildingmodels.SetSysBuildingListReq) (*buildingmodels.SetSysBuildingList, error)
	Remove(c *gin.Context, ids []string) error
	AllDetail(c *gin.Context) ([]*buildingmodels.SysBuilding, error)
}
