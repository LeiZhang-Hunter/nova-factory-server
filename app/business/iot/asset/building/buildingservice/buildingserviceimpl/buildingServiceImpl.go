package buildingserviceimpl

import (
	buildingDao2 "nova-factory-server/app/business/iot/asset/building/buildingdao"
	buildingModels2 "nova-factory-server/app/business/iot/asset/building/buildingmodels"
	"nova-factory-server/app/business/iot/asset/building/buildingservice"

	"github.com/gin-gonic/gin"
)

type BuildingServiceImpl struct {
	dao      buildingDao2.BuildingDao
	floorDao buildingDao2.FloorDao
}

func NewBuildingServiceImpl(dao buildingDao2.BuildingDao, floorDao buildingDao2.FloorDao) buildingservice.BuildingService {
	return &BuildingServiceImpl{
		dao:      dao,
		floorDao: floorDao,
	}
}

func (b *BuildingServiceImpl) Set(c *gin.Context, data *buildingModels2.SetSysBuilding) (*buildingModels2.SysBuilding, error) {
	return b.dao.Set(c, data)
}

func (b *BuildingServiceImpl) List(c *gin.Context, req *buildingModels2.SetSysBuildingListReq) (*buildingModels2.SetSysBuildingList, error) {
	return b.dao.List(c, req)
}
func (b *BuildingServiceImpl) Remove(c *gin.Context, ids []string) error {
	return b.dao.Remove(c, ids)
}

func (b *BuildingServiceImpl) AllDetail(c *gin.Context) ([]*buildingModels2.SysBuilding, error) {
	allBuild, err := b.dao.All(c)
	if err != nil {
		return nil, err
	}

	allFloor, err := b.floorDao.All(c)
	if err != nil {
		return nil, err
	}

	var floorMap = make(map[int64][]*buildingModels2.SysFloor)
	for _, floor := range allFloor {
		_, ok := floorMap[floor.BuildingID]
		if !ok {
			floorMap[floor.BuildingID] = make([]*buildingModels2.SysFloor, 0)
		}
		floorMap[floor.BuildingID] = append(floorMap[floor.BuildingID], floor)
	}

	for k, build := range allBuild {
		floorList, ok := floorMap[build.ID]
		if !ok {
			continue
		}
		allBuild[k].Floor = floorList
	}

	return allBuild, nil
}
