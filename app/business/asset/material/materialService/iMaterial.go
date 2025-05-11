package materialService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/asset/material/materialModels"
)

type IMaterialService interface {
	InsertMaterial(c *gin.Context, job *materialModels.MaterialInfo) (*materialModels.MaterialVO, error)
	UpdateMaterial(c *gin.Context, job *materialModels.MaterialInfo) (*materialModels.MaterialVO, error)
	SelectMaterialList(c *gin.Context, req *materialModels.MaterialListReq) (*materialModels.MaterialInfoListValue, error)
	DeleteByMaterialIds(c *gin.Context, ids []int64) error
	GetByMaterialId(c *gin.Context, id int64) (*materialModels.MaterialVO, error)
	Inbound(c *gin.Context, info *materialModels.InboundInfo) (*materialModels.InboundVO, error)
	Outbound(c *gin.Context, info *materialModels.OutboundInfo) (*materialModels.OutboundVO, error)
	InboundList(c *gin.Context, info *materialModels.InboundListReq) (*materialModels.InboundListData, error)
	OutboundList(c *gin.Context, info *materialModels.OutboundListReq) (*materialModels.OutboundListData, error)
}
