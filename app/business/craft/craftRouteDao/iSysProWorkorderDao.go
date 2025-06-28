package craftRouteDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type ISysProWorkorderDao interface {
	Add(c *gin.Context, info *craftRouteModels.SysSetProWorkorder) (*craftRouteModels.SysProWorkorder, error)
	Update(c *gin.Context, info *craftRouteModels.SysSetProWorkorder) (*craftRouteModels.SysProWorkorder, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *craftRouteModels.SysProWorkorderReq) (*craftRouteModels.SysProWorkorderList, error)
}
