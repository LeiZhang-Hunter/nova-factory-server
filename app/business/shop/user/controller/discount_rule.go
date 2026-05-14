package shopcontroller

import (
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/business/shop/user/service"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
)

// DiscountRule 折扣规则管理
type DiscountRule struct {
	service service.IShopUserDiscountRuleService
}

func NewDiscountRule(svc service.IShopUserDiscountRuleService) *DiscountRule {
	return &DiscountRule{service: svc}
}

// PrivateRoutes 注册折扣规则路由
func (s *DiscountRule) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/user/discount")
	group.GET("/list", middlewares.HasPermission("shop:user:discount:list"), s.List)
	group.GET("/info/:id", middlewares.HasPermission("shop:user:discount:query"), s.GetByID)
	group.GET("/ids/:userId", middlewares.HasPermission("shop:user:discount:list"), s.GetUsedTargetIds)
	group.POST("/add", middlewares.HasPermission("shop:user:discount:add"), s.Create)
	group.POST("/batch", middlewares.HasPermission("shop:user:discount:add"), s.BatchCreate)
	group.PUT("/update", middlewares.HasPermission("shop:user:discount:update"), s.Update)
	group.DELETE("/remove/:ids", middlewares.HasPermission("shop:user:discount:remove"), s.Delete)
}

// List 获取折扣规则列表
// @Summary 获取折扣规则列表
// @Description 获取折扣规则列表
// @Tags 商城/用户管理/折扣规则
// @Param object query models.UserDiscountRuleQuery true "折扣规则查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/discount/list [get]
func (s *DiscountRule) List(c *gin.Context) {
	req := new(models.UserDiscountRuleQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	var data interface{}
	var err error

	if req.UserID > 0 && req.TargetType != "" {
		data, err = s.service.ListByUserIDAndType(c, req.UserID, req.TargetType, req.Page, req.Size)
	} else if req.UserID > 0 {
		data, err = s.service.ListByUserID(c, req.UserID, req.Page, req.Size)
	} else {
		baizeContext.ParameterError(c)
		return
	}

	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetByID 获取折扣规则详情
// @Summary 获取折扣规则详情
// @Description 根据ID获取折扣规则详情
// @Tags 商城/用户管理/折扣规则
// @Param id path int true "折扣规则ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/discount/info/{id} [get]
func (s *DiscountRule) GetByID(c *gin.Context) {
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

// GetUsedTargetIds 获取用户所有已使用的目标ID
// @Summary 获取用户已使用目标ID
// @Description 获取用户所有折扣规则关联的商品ID和分类ID列表
// @Tags 商城/用户管理/折扣规则
// @Param userId path int true "用户ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData{data=map[string][]int64} "获取成功"
// @Router /shop/user/discount/ids/{userId} [get]
func (s *DiscountRule) GetUsedTargetIds(c *gin.Context) {
	userId := baizeContext.ParamInt64(c, "userId")
	if userId == 0 {
		baizeContext.ParameterError(c)
		return
	}
	goodsIds, categoryIds, err := s.service.GetUsedTargetIds(c, userId)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, gin.H{
		"goodsIds":    goodsIds,
		"categoryIds": categoryIds,
	})
}

// Create 新增折扣规则
// @Summary 新增折扣规则
// @Description 新增折扣规则
// @Tags 商城/用户管理/折扣规则
// @Param object body models.UserDiscountRuleUpsert true "折扣规则新增参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "新增成功"
// @Router /shop/user/discount/add [post]
func (s *DiscountRule) Create(c *gin.Context) {
	req := new(models.UserDiscountRuleUpsert)
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

// BatchCreate 批量创建折扣规则
// @Summary 批量创建折扣规则
// @Description 批量创建折扣规则
// @Tags 商城/用户管理/折扣规则
// @Param object body models.BatchDiscountRuleCreate true "批量创建参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "创建成功"
// @Router /shop/user/discount/batch [post]
func (s *DiscountRule) BatchCreate(c *gin.Context) {
	req := new(models.BatchDiscountRuleCreate)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if len(req.UserIDs) == 0 || len(req.TargetIDs) == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if req.TargetType != "goods" && req.TargetType != "category" {
		baizeContext.ParameterError(c)
		return
	}
	if req.DiscountRate <= 0 || req.DiscountRate > 1 {
		baizeContext.ParameterError(c)
		return
	}

	count, err := s.service.BatchCreate(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, map[string]int64{"count": count})
}

// Update 修改折扣规则
// @Summary 修改折扣规则
// @Description 修改折扣规则
// @Tags 商城/用户管理/折扣规则
// @Param object body models.UserDiscountRuleUpsert true "折扣规则修改参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "修改成功"
// @Router /shop/user/discount/update [put]
func (s *DiscountRule) Update(c *gin.Context) {
	req := new(models.UserDiscountRuleUpsert)
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

// Delete 删除折扣规则
// @Summary 删除折扣规则
// @Description 根据ID删除折扣规则
// @Tags 商城/用户管理/折扣规则
// @Param ids path string true "折扣规则ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/user/discount/remove/{ids} [delete]
func (s *DiscountRule) Delete(c *gin.Context) {
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
