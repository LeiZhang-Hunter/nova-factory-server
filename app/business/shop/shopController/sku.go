package shopController

import (
	"nova-factory-server/app/business/shop/shopModels"
	"nova-factory-server/app/business/shop/shopService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Sku struct {
	service shopService.IShopSkuService
}

func NewSku(service shopService.IShopSkuService) *Sku {
	return &Sku{service: service}
}

func (s *Sku) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/sku")
	group.GET("/list", middlewares.HasPermission("shop:sku:list"), s.List)
	group.GET("/:id", middlewares.HasPermission("shop:sku:query"), s.GetByID)
	group.POST("", middlewares.HasPermission("shop:sku:add"), s.Create)
	group.PUT("", middlewares.HasPermission("shop:sku:edit"), s.Update)
	group.DELETE("/:ids", middlewares.HasPermission("shop:sku:remove"), s.Delete)
}

func (s *Sku) List(c *gin.Context) {
	req := new(shopModels.GoodsSkuQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

func (s *Sku) GetByID(c *gin.Context) {
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.GetByID(c, id)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

func (s *Sku) Create(c *gin.Context) {
	req := new(shopModels.GoodsSkuUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.Create(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

func (s *Sku) Update(c *gin.Context) {
	req := new(shopModels.GoodsSkuUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.ID == 0 {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.Update(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

func (s *Sku) Delete(c *gin.Context) {
	ids := baizeContext.ParamInt64Array(c, "ids")
	if len(ids) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := s.service.DeleteByIDs(c, ids); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}
