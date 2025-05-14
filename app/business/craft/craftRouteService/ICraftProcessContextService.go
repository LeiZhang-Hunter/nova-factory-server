package craftRouteService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type ICraftProcessContextService interface {
	Add(c *gin.Context, processContext *craftRouteModels.SysProSetProcessContent) (*craftRouteModels.SysProProcessContent, error)

	List(c *gin.Context, processContext *craftRouteModels.SysProProcessContextListReq) (*craftRouteModels.SysProProcessContextListData, error)

	Update(c *gin.Context, processContext *craftRouteModels.SysProSetProcessContent) (*craftRouteModels.SysProProcessContent, error)

	Remove(c *gin.Context, ids []string) error
}
