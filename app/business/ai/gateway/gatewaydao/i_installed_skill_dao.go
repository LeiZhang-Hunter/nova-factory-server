package gatewaydao

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

type IInstalledSkillDao interface {
	Create(c *gin.Context, req *gatewaymodels.InstalledSkillUpsert) (*gatewaymodels.InstalledSkill, error)
	Update(c *gin.Context, req *gatewaymodels.InstalledSkillUpsert) (*gatewaymodels.InstalledSkill, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*gatewaymodels.InstalledSkill, error)
	List(c *gin.Context, req *gatewaymodels.InstalledSkillQuery) (*gatewaymodels.InstalledSkillListData, error)
}
