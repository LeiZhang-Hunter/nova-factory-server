package shopcontroller

import (
	"nova-factory-server/app/business/shop/shopmodels"
	"nova-factory-server/app/business/shop/shopservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type Category struct {
	service shopservice.IShopCategoryService
}

func NewCategory(service shopservice.IShopCategoryService) *Category {
	return &Category{service: service}
}

func (s *Category) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/category")
	group.GET("/list", middlewares.HasPermission("shop:category:list"), s.List)
	group.GET("/:id", middlewares.HasPermission("shop:category:query"), s.GetByID)
	group.POST("", middlewares.HasPermission("shop:category:add"), s.Create)
	group.PUT("", middlewares.HasPermission("shop:category:edit"), s.Update)
	group.DELETE("/:ids", middlewares.HasPermission("shop:category:remove"), s.Delete)
}

func (s *Category) List(c *gin.Context) {
	req := new(shopmodels.CategoryQuery)
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

func (s *Category) GetByID(c *gin.Context) {
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

func (s *Category) Create(c *gin.Context) {
	req := new(shopmodels.CategoryUpsert)
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

func (s *Category) Update(c *gin.Context) {
	req := new(shopmodels.CategoryUpsert)
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

func (s *Category) Delete(c *gin.Context) {
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
