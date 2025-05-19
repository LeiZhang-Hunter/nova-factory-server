package craftRouteDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

// ISysProRouteProductDao 产品制程
type ISysProRouteProductDao interface {
	Add(c *gin.Context, data *craftRouteModels.SysProRouteSetProduct) (*craftRouteModels.SysProRouteProduct, error)
	Update(c *gin.Context, data *craftRouteModels.SysProRouteSetProduct) (*craftRouteModels.SysProRouteProduct, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *craftRouteModels.SysProRouteProductReq) (*craftRouteModels.SysProRouteProductList, error)
}
