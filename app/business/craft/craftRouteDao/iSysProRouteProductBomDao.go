package craftRouteDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type ISysProRouteProductBomDao interface {
	Add(c *gin.Context, info *craftRouteModels.SysSetProRouteProductBom) (*craftRouteModels.SysProRouteProductBom, error)
	Update(c *gin.Context, info *craftRouteModels.SysSetProRouteProductBom) (*craftRouteModels.SysProRouteProductBom, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *craftRouteModels.SysProRouteProductBomReq) (*craftRouteModels.SysProRouteProductBomList, error)
}
