package craftRouteService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/craft/craftRouteModels"
)

type ISysProRouteProductService interface {
	Add(c *gin.Context, data *craftRouteModels.SysProRouteSetProduct) (*craftRouteModels.SysProRouteProduct, error)
	Update(c *gin.Context, data *craftRouteModels.SysProRouteSetProduct) (*craftRouteModels.SysProRouteProduct, error)
	List(c *gin.Context, req *craftRouteModels.SysProRouteProductReq) (*craftRouteModels.SysProRouteProductList, error)
	Remove(c *gin.Context, ids []string) error
}
