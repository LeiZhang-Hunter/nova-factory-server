package craftrouteservice

import (
	"nova-factory-server/app/business/iot/craft/craftroutemodels"

	"github.com/gin-gonic/gin"
)

type ICraftProcessContextService interface {
	Add(c *gin.Context, processContext *craftroutemodels.SysProSetProcessContent) (*craftroutemodels.SysProProcessContent, error)

	List(c *gin.Context, processContext *craftroutemodels.SysProProcessContextListReq) (*craftroutemodels.SysProProcessContextListData, error)

	Update(c *gin.Context, processContext *craftroutemodels.SysProSetProcessContent) (*craftroutemodels.SysProProcessContent, error)

	Remove(c *gin.Context, ids []string) error
}
