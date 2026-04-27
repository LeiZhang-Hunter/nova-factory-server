package shopcontroller

import (
	"fmt"
	"net/url"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type Category struct {
	service shopservice.IShopCategoryService
}

var categoryUnsafeHTMLPattern = regexp.MustCompile(`(?i)<[^>]+>`)

func NewCategory(service shopservice.IShopCategoryService) *Category {
	return &Category{service: service}
}

// PrivateRoutes 注册商品分类路由
func (s *Category) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/shop/category")
	group.GET("/list", middlewares.HasPermission("shop:category:list"), s.List)
	group.GET("/info/:id", middlewares.HasPermission("shop:category:info"), s.GetByID)
	group.POST("/add", middlewares.HasPermission("shop:category:add"), s.Create)
	group.PUT("/update", middlewares.HasPermission("shop:category:update"), s.Update)
	group.DELETE("/:ids", middlewares.HasPermission("shop:category:remove"), s.Delete)
}

// List 获取商品分类列表
// @Summary 获取商品分类列表
// @Description 获取商品分类列表
// @Tags 商城/商品分类
// @Param object query shopmodels.CategoryQuery true "商品分类查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/category/list [get]
func (s *Category) List(c *gin.Context) {
	req := new(shopmodels.CategoryQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	req.CategoryName = strings.TrimSpace(req.CategoryName)
	data, err := s.service.List(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetByID 获取商品分类详情
// @Summary 获取商品分类详情
// @Description 根据ID获取商品分类详情
// @Tags 商城/商品分类
// @Param id path int true "商品分类ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/category/{id} [get]
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

// Create 新增商品分类
// @Summary 新增商品分类
// @Description 新增商品分类
// @Tags 商城/商品分类
// @Param object body shopmodels.CategoryUpsert true "商品分类新增参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "新增成功"
// @Router /shop/category [post]
func (s *Category) Create(c *gin.Context) {
	req := new(shopmodels.CategoryUpsert)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if err := validateCategoryUpsertSecurity(req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	data, err := s.service.Create(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Update 修改商品分类
// @Summary 修改商品分类
// @Description 修改商品分类
// @Tags 商城/商品分类
// @Param object body shopmodels.CategoryUpsert true "商品分类修改参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "修改成功"
// @Router /shop/category [put]
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
	if err := validateCategoryUpsertSecurity(req); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	data, err := s.service.Update(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Delete 删除商品分类
// @Summary 删除商品分类
// @Description 根据ID删除商品分类
// @Tags 商城/商品分类
// @Param ids path string true "商品分类ID，多个以逗号分隔"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "删除成功"
// @Router /shop/category/{ids} [delete]
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

func validateCategoryUpsertSecurity(req *shopmodels.CategoryUpsert) error {
	req.CategoryName = strings.TrimSpace(req.CategoryName)
	req.CategoryCode = strings.TrimSpace(req.CategoryCode)
	req.ImageURL = strings.TrimSpace(req.ImageURL)
	req.Description = strings.TrimSpace(req.Description)

	if req.CategoryName == "" {
		return fmt.Errorf("分类名称不能为空")
	}
	if hasUnsafeCategoryText(req.CategoryName) {
		return fmt.Errorf("分类名称包含不安全内容")
	}
	if req.CategoryCode != "" && hasUnsafeCategoryText(req.CategoryCode) {
		return fmt.Errorf("分类编号包含不安全内容")
	}
	if req.ImageURL == "" {
		return fmt.Errorf("分类图片不能为空")
	}
	if req.Description == "" {
		return fmt.Errorf("分类描述不能为空")
	}
	if req.Description != "" && hasUnsafeCategoryText(req.Description) {
		return fmt.Errorf("分类描述包含不安全内容")
	}
	if err := validateCategoryImageURL(req.ImageURL); err != nil {
		return err
	}

	return nil
}

func hasUnsafeCategoryText(value string) bool {
	lowerValue := strings.ToLower(strings.TrimSpace(value))
	if lowerValue == "" {
		return false
	}
	if categoryUnsafeHTMLPattern.MatchString(value) {
		return true
	}

	dangerousTokens := []string{
		"javascript:",
		"vbscript:",
		"onerror=",
		"onload=",
		"<script",
		"</script",
		"<iframe",
		"<svg",
		"data:text/html",
	}
	for _, token := range dangerousTokens {
		if strings.Contains(lowerValue, token) {
			return true
		}
	}

	return false
}

func validateCategoryImageURL(imageURL string) error {
	if imageURL == "" {
		return nil
	}
	if hasUnsafeCategoryText(imageURL) {
		return fmt.Errorf("分类图片地址包含不安全内容")
	}
	if strings.HasPrefix(imageURL, "/") {
		return nil
	}

	parsedURL, err := url.ParseRequestURI(imageURL)
	if err != nil {
		return fmt.Errorf("分类图片地址格式不正确")
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("分类图片地址仅支持http或https")
	}

	return nil
}
