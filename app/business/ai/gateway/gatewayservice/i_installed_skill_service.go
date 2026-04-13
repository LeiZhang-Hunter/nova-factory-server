package gatewayservice

import (
	"nova-factory-server/app/business/ai/gateway/gatewaymodels"

	"github.com/gin-gonic/gin"
)

// IInstalledSkillService 已安装技能服务接口
type IInstalledSkillService interface {
	Create(c *gin.Context, req *gatewaymodels.InstalledSkillUpsert) (*gatewaymodels.InstalledSkill, error)
	Update(c *gin.Context, req *gatewaymodels.InstalledSkillUpsert) (*gatewaymodels.InstalledSkill, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*gatewaymodels.InstalledSkill, error)
	List(c *gin.Context, req *gatewaymodels.InstalledSkillQuery) (*gatewaymodels.InstalledSkillListData, error)
}
