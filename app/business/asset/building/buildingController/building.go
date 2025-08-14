package buildingController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/building/buildingService"
)

type Building struct {
	service buildingService.BuildingService
}

func NewBuilding(service buildingService.BuildingService) *Building {
	return &Building{
		service: service,
	}
}

func (b *Building) PrivateRoutes(router *gin.RouterGroup) {

}

func (b *Building) Set(C *gin.Context) {

}

func (b *Building) Remove(C *gin.Context) {

}

func (b *Building) List(C *gin.Context) {

}
