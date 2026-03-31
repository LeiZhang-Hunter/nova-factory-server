package materialService

import (
	materialModels2 "nova-factory-server/app/business/iot/asset/material/materialModels"

	"github.com/gin-gonic/gin"
)

type IMaterialService interface {
	InsertMaterial(c *gin.Context, job *materialModels2.MaterialInfo) (*materialModels2.MaterialVO, error)
	UpdateMaterial(c *gin.Context, job *materialModels2.MaterialInfo) (*materialModels2.MaterialVO, error)
	SelectMaterialList(c *gin.Context, req *materialModels2.MaterialListReq) (*materialModels2.MaterialInfoListValue, error)
	DeleteByMaterialIds(c *gin.Context, ids []int64) error
	GetByMaterialId(c *gin.Context, id int64) (*materialModels2.MaterialVO, error)
	Inbound(c *gin.Context, info *materialModels2.InboundInfo) (*materialModels2.InboundVO, error)
	Outbound(c *gin.Context, info *materialModels2.OutboundInfo) (*materialModels2.OutboundVO, error)
	InboundList(c *gin.Context, info *materialModels2.InboundListReq) (*materialModels2.InboundListData, error)
	OutboundList(c *gin.Context, info *materialModels2.OutboundListReq) (*materialModels2.OutboundListData, error)
}
