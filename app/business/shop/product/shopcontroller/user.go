package shopcontroller

import (
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"
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

// PrivateRoutes 注册商城用户路由
func (s *User) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/user")
	group.GET("/list", middlewares.HasPermission("shop:user:list"), s.List)
	group.GET("/:id", middlewares.HasPermission("shop:user:query"), s.GetByID)
	group.POST("", middlewares.HasPermission("shop:user:add"), s.Create)
	group.PUT("", middlewares.HasPermission("shop:user:edit"), s.Update)
	group.DELETE("/:ids", middlewares.HasPermission("shop:user:remove"), s.Delete)
}

// List 获取商城用户列表
// @Summary 获取商城用户列表
// @Description 获取商城用户列表
// @Tags 商城/用户管理
// @Param object query shopmodels.UserQuery true "商城用户查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/list [get]
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

// GetByID 获取商城用户详情
// @Summary 获取商城用户详情
// @Description 根据ID获取商城用户详情
// @Tags 商城/用户管理
// @Param id path int true "商城用户ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/{id} [get]
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

// Create 新增商城用户
// @Summary 新增商城用户
// @Description 新增商城用户
// @Tags 商城/用户管理
// @Param object body shopmodels.UserUpsert true "商城用户新增参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "新增成功"
// @Router /shop/user [post]
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

// Update 修改商城用户
// @Summary 修改商城用户
// @Description 修改商城用户
// @Tags 商城/用户管理
// @Param object body shopmodels.UserUpsert true "商城用户修改参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "修改成功"
// @Router /shop/user [put]
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

// Delete 删除商城用户
// @Summary 删除商城用户
// @Description 根据ID删除商城用户
// @Tags 商城/用户管理
// @Param ids path string true "商城用户ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/user/{ids} [delete]
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
