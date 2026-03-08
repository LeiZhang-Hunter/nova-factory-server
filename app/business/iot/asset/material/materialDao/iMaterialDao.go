package materialDao

import (
	materialModels2 "nova-factory-server/app/business/iot/asset/material/materialModels"

	"github.com/gin-gonic/gin"
)

type IMaterialDao interface {
	InsertMaterial(c *gin.Context, material *materialModels2.MaterialInfo) (*materialModels2.MaterialVO, error)
	UpdateMaterial(c *gin.Context, material *materialModels2.MaterialInfo) (*materialModels2.MaterialVO, error)
	GetMaterialGroupByName(c *gin.Context, name string) (*materialModels2.MaterialVO, error)
	GetNoExitIdMaterialGroupByName(c *gin.Context, name string, id uint64) (*materialModels2.MaterialVO, error)
	SelectMaterialList(c *gin.Context, req *materialModels2.MaterialListReq) (*materialModels2.MaterialInfoListData, error)
	DeleteByMaterialIds(c *gin.Context, ids []int64) error
	GetByMaterialId(c *gin.Context, id int64) (*materialModels2.MaterialVO, error)
	// Inbound 入库
	Inbound(c *gin.Context, info *materialModels2.InboundInfo) (*materialModels2.InboundVO, error)
	// Outbound 出库
	Outbound(c *gin.Context, info *materialModels2.OutboundInfo) (*materialModels2.OutboundVO, error)
	// InboundList 入库列表
	InboundList(c *gin.Context, req *materialModels2.InboundListReq) (*materialModels2.InboundListData, error)
	// OutboundList 出库列表
	OutboundList(c *gin.Context, req *materialModels2.OutboundListReq) (*materialModels2.OutboundListData, error)
}
