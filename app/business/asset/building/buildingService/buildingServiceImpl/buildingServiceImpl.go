package buildingServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/building/buildingDao"
	"nova-factory-server/app/business/asset/building/buildingModels"
	"nova-factory-server/app/business/asset/building/buildingService"
)

type BuildingServiceImpl struct {
	dao      buildingDao.BuildingDao
	floorDao buildingDao.FloorDao
}

func NewBuildingServiceImpl(dao buildingDao.BuildingDao, floorDao buildingDao.FloorDao) buildingService.BuildingService {
	return &BuildingServiceImpl{
		dao:      dao,
		floorDao: floorDao,
	}
}

func (b *BuildingServiceImpl) Set(c *gin.Context, data *buildingModels.SetSysBuilding) (*buildingModels.SysBuilding, error) {
	return b.dao.Set(c, data)
}

func (b *BuildingServiceImpl) List(c *gin.Context, req *buildingModels.SetSysBuildingListReq) (*buildingModels.SetSysBuildingList, error) {
	return b.dao.List(c, req)
}
func (b *BuildingServiceImpl) Remove(c *gin.Context, ids []string) error {
	return b.dao.Remove(c, ids)
}

func (b *BuildingServiceImpl) AllDetail(c *gin.Context) ([]*buildingModels.SysBuilding, error) {
	allBuild, err := b.dao.All(c)
	if err != nil {
		return nil, err
	}

	allFloor, err := b.floorDao.All(c)
	if err != nil {
		return nil, err
	}

	var floorMap = make(map[int64][]*buildingModels.SysFloor)
	for _, floor := range allFloor {
		_, ok := floorMap[floor.BuildingID]
		if !ok {
			floorMap[floor.BuildingID] = make([]*buildingModels.SysFloor, 0)
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
