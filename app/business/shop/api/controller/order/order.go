package order

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	apiDao "nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Order 订单控制器
type Order struct {
	service            service.IApiShopOrderService
	orderRefundService service.IApiShopOrderRefundService
	configDao          apiDao.IApiShopSysConfigDao
}

// NewOrder 创建订单控制器
func NewOrder(service service.IApiShopOrderService, orderRefundService service.IApiShopOrderRefundService, configDao apiDao.IApiShopSysConfigDao) *Order {
	return &Order{service: service, orderRefundService: orderRefundService, configDao: configDao}
}

// PrivateRoutes 注册商城订单路由（商城模块只检查登录，不检查权限）
func (s *Order) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/order")
	group.GET("/list", s.List)
	group.GET("/info/:id", s.GetByID)
	group.GET("/count", s.GetStatistics)
	group.POST("/confirm", s.Confirm)
	group.POST("/create", s.Create)
	group.POST("/pay/:id", s.Pay)
	group.POST("/cancel/:id", s.Cancel)
	group.POST("/confirm/:id", s.ConfirmReceive)
	// 售后/退款申请
	group.POST("/refund/apply", s.Apply)
	// 退货物流提交
	group.POST("/refund/submit-logistics", s.SubmitReturnLogistics)
	// 支付凭证提交
	group.POST("/submit-payment-voucher", s.SubmitPaymentVoucher)
	// 企业账户收款信息
	router.GET("/api/v1/app/shop/config/enterprise-account", s.GetEnterpriseAccountConfig)
	// 文件上传（支付凭证）
	router.POST("/api/v1/app/shop/upload", s.UploadFile)
}

// List 获取订单列表
// @Summary 获取订单列表
// @Description 获取当前用户的订单列表
// @Tags 商城/用户订单
// @Param object query models.OrderQuery true "订单查询参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/order/list [get]
func (s *Order) List(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	req := new(models.OrderQuery)
	if err := c.ShouldBindQuery(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.List(c, userID, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// GetByID 获取订单详情
// @Summary 获取订单详情
// @Description 根据ID获取订单详情
// @Tags 商城/用户订单
// @Param id path int true "订单ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/order/info/{id} [get]
func (s *Order) GetByID(c *gin.Context) {
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

// Confirm 获取确认单数据
// @Summary 获取确认单数据
// @Description 根据当前地址、配送方式等信息实时试算确认单
// @Tags app接口/商城/App订单
// @Param object body models.OrderConfirmReq true "确认单参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /api/v1/app/shop/order/confirm [post]
func (s *Order) Confirm(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	req := new(models.OrderConfirmReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.Confirm(c, userID, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Create 创建订单
// @Summary 创建订单
// @Description 根据预订单 orderKey 正式落订单
// @Tags app接口/商城/App订单
// @Param object body models.OrderCreateReq true "订单创建参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "创建成功"
// @Router /api/v1/app/shop/order/create [post]
func (s *Order) Create(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	req := new(models.OrderCreateReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	data, err := s.service.Create(c, userID, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Pay 支付订单
// @Summary 支付订单
// @Description 支付待付款订单，可通过payChannel指定支付通道（0=使用订单默认通道 1=微信 2=支付宝）
// @Tags 商城/用户订单
// @Param id path int true "订单ID"
// @Param object body models.OrderPayReq false "支付通道参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "支付成功"
// @Router /api/v1/app/shop/order/pay/{id} [post]
func (s *Order) Pay(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	var req models.OrderPayReq
	_ = c.ShouldBindJSON(&req) // 空 body/格式错误 → PayChannel=0 → 走订单默认通道
	data, err := s.service.Pay(c, userID, id, req.PayChannel)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Cancel 取消订单
// @Summary 取消订单
// @Description 取消待支付订单
// @Tags 商城/用户订单
// @Param id path int true "订单ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "取消成功"
// @Router /shop/user/order/cancel/{id} [post]
func (s *Order) Cancel(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	reason := c.Query("reason")
	if err := s.service.Cancel(c, userID, id, reason); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// ConfirmReceive 确认收货
// @Summary 确认收货
// @Description 确认已发货订单收货
// @Tags 商城/用户订单
// @Param id path int true "订单ID"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "确认成功"
// @Router /shop/user/order/confirm/{id} [post]
func (s *Order) ConfirmReceive(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	id := baizeContext.ParamInt64(c, "id")
	if id == 0 {
		baizeContext.ParameterError(c)
		return
	}
	if err := s.service.ConfirmReceive(c, userID, id); err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// GetStatistics 获取订单统计
// @Summary 获取订单统计
// @Description 获取用户各状态订单数量统计
// @Tags 商城/用户订单
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "获取成功"
// @Router /shop/user/order/count [get]
func (s *Order) GetStatistics(c *gin.Context) {
	userID := baizeContext.GetUserId(c)
	data, err := s.service.GetStatistics(c, userID)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, data)
}

// Apply 申请售后
// @Summary 申请售后
// @Description 申请售后
// @Tags 商城/售后
// @Param object body models.RefundApplyReq true "售后申请参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "申请成功"
// @Router /shop/user/refund/apply [post]
func (s *Order) Apply(c *gin.Context) {
	var req models.RefundApplyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	userID := baizeContext.GetUserId(c)
	if userID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "请先登录"})
		return
	}
	resp, err := s.orderRefundService.Apply(c, userID, &req)
	if err != nil {
		zap.L().Error("申请售后失败", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": resp})
}

// SubmitReturnLogistics 提交退货物流信息
// @Summary 提交退货物流信息
// @Description 已审核的退货退款类型售后单，填写退货物流单号后再次同步ERP
// @Tags 商城/售后
// @Param object body models.SubmitReturnLogisticsReq true "退货物流参数"
// @Security BearerAuth
// @Produce application/json
// @Success 200 {object} response.ResponseData "提交成功"
// @Router /api/v1/app/shop/order/refund/submit-logistics [post]
func (s *Order) SubmitReturnLogistics(c *gin.Context) {
	var req models.SubmitReturnLogisticsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	userID := baizeContext.GetUserId(c)
	if userID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "请先登录"})
		return
	}
	resp, err := s.orderRefundService.SubmitReturnLogistics(c, userID, &req)
	if err != nil {
		zap.L().Error("提交退货物流失败", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": resp})
}

// GetEnterpriseAccountConfig 获取企业账户收款信息
func (s *Order) GetEnterpriseAccountConfig(c *gin.Context) {
	data, err := s.configDao.GetEnterpriseAccountConfig(c)
	if err != nil {
		zap.L().Error("获取企业账户配置失败", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "获取失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": data})
}

// SubmitPaymentVoucher 提交线下打款支付凭证
func (s *Order) SubmitPaymentVoucher(c *gin.Context) {
	var req models.PaymentVoucherReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "参数错误"})
		return
	}
	userID := baizeContext.GetUserId(c)
	if userID == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "请先登录"})
		return
	}
	if err := s.service.SubmitPaymentVoucher(c, userID, &req); err != nil {
		zap.L().Error("提交支付凭证失败", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success"})
}

// UploadFile 上传文件（支付凭证图片）
func (s *Order) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "msg": "请选择文件"})
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := fmt.Sprintf("voucher_%d%s", time.Now().UnixNano(), ext)
	uploadDir := "uploads/payment"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "创建上传目录失败"})
		return
	}

	savePath := filepath.Join(uploadDir, filename)
	out, err := os.Create(savePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "文件保存失败"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "msg": "文件写入失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "data": map[string]string{
		"url": "/" + filepath.ToSlash(savePath),
	}})
}
