package buildingdao

import (
	"nova-factory-server/app/business/iot/asset/building/buildingmodels"

	"github.com/gin-gonic/gin"
)

type BuildingDao interface {
	Set(c *gin.Context, data *buildingmodels.SetSysBuilding) (*buildingmodels.SysBuilding, error)
	List(c *gin.Context, req *buildingmodels.SetSysBuildingListReq) (*buildingmodels.SetSysBuildingList, error)
	Remove(c *gin.Context, ids []string) error
	GetByIds(c *gin.Context, ids []uint64) ([]*buildingmodels.SysBuilding, error)
	All(c *gin.Context) ([]*buildingmodels.SysBuilding, error)
	Count(c *gin.Context) (int64, error)
}
