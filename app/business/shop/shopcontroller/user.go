package shopcontroller

import (
	"nova-factory-server/app/business/shop/shopmodels"
	"nova-factory-server/app/business/shop/shopservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

type User struct {
	service shopservice.IShopUserService
}

func NewUser(service shopservice.IShopUserService) *User {
	return &User{service: service}
}

func (s *User) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/user")
	group.GET("/list", middlewares.HasPermission("shop:user:list"), s.List)
	group.GET("/:id", middlewares.HasPermission("shop:user:query"), s.GetByID)
	group.POST("", middlewares.HasPermission("shop:user:add"), s.Create)
	group.PUT("", middlewares.HasPermission("shop:user:edit"), s.Update)
	group.DELETE("/:ids", middlewares.HasPermission("shop:user:remove"), s.Delete)
}

func (s *User) List(c *gin.Context) {
	req := new(shopmodels.UserQuery)
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

func (s *User) GetByID(c *gin.Context) {
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

func (s *User) Create(c *gin.Context) {
	req := new(shopmodels.UserUpsert)
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

func (s *User) Update(c *gin.Context) {
	req := new(shopmodels.UserUpsert)
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

func (s *User) Delete(c *gin.Context) {
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
