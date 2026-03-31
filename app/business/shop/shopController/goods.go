package shopController

import (
	"nova-factory-server/app/business/shop/shopModels"
	"nova-factory-server/app/business/shop/shopService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Goods struct {
	service shopService.IShopGoodsService
}

func NewGoods(service shopService.IShopGoodsService) *Goods {
	return &Goods{service: service}
}

func (s *Goods) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/goods")
	group.GET("/list", middlewares.HasPermission("shop:goods:list"), s.List)
	group.GET("/:id", middlewares.HasPermission("shop:goods:query"), s.GetByID)
	group.POST("", middlewares.HasPermission("shop:goods:add"), s.Create)
	group.PUT("", middlewares.HasPermission("shop:goods:edit"), s.Update)
	group.DELETE("/:ids", middlewares.HasPermission("shop:goods:remove"), s.Delete)
}

func (s *Goods) List(c *gin.Context) {
	req := new(shopModels.GoodsQuery)
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

func (s *Goods) GetByID(c *gin.Context) {
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

func (s *Goods) Create(c *gin.Context) {
	req := new(shopModels.GoodsUpsert)
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

func (s *Goods) Update(c *gin.Context) {
	req := new(shopModels.GoodsUpsert)
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

func (s *Goods) Delete(c *gin.Context) {
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
