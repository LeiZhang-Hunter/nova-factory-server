package materialDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/material/materialModels"
)

type IMaterialDao interface {
	InsertMaterial(c *gin.Context, material *materialModels.MaterialInfo) (*materialModels.MaterialVO, error)
	UpdateMaterial(c *gin.Context, material *materialModels.MaterialInfo) (*materialModels.MaterialVO, error)
	GetMaterialGroupByName(c *gin.Context, name string) (*materialModels.MaterialVO, error)
	GetNoExitIdMaterialGroupByName(c *gin.Context, name string, id uint64) (*materialModels.MaterialVO, error)
	SelectMaterialList(c *gin.Context, req *materialModels.MaterialListReq) (*materialModels.MaterialInfoListData, error)
	DeleteByMaterialIds(c *gin.Context, ids []int64) error
	GetByMaterialId(c *gin.Context, id int64) (*materialModels.MaterialVO, error)
	// Inbound 入库
	Inbound(c *gin.Context, info *materialModels.InboundInfo) (*materialModels.InboundVO, error)
	// Outbound 出库
	Outbound(c *gin.Context, info *materialModels.OutboundInfo) (*materialModels.OutboundVO, error)
	// InboundList 入库列表
	InboundList(c *gin.Context, req *materialModels.InboundListReq) (*materialModels.InboundListData, error)
	// OutboundList 出库列表
	OutboundList(c *gin.Context, req *materialModels.OutboundListReq) (*materialModels.OutboundListData, error)
}
