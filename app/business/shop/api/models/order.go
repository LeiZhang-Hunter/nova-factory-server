package models

import (
	"fmt"
	"nova-factory-server/app/baize"
	shopordermodels "nova-factory-server/app/business/shop/order/models"
	shopusermodels "nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/utils/stringUtils"
)

// OrderItem 订单商品明细
type OrderItem struct {
	ID          int64   `json:"id,string" gorm:"column:id"`             // 主键ID
	OrderID     int64   `json:"orderId,string" gorm:"column:order_id"`  // 订单ID
	OrderNo     string  `json:"orderNo" gorm:"column:order_no"`         // 订单号
	GoodsID     string  `json:"goodsId" gorm:"column:goods_id"`         // 商品ID
	SkuID       string  `json:"skuId" gorm:"column:sku_id"`             // SKU ID
	GoodsName   string  `json:"goodsName" gorm:"column:goods_name"`     // 商品名称快照
	SkuName     string  `json:"skuName" gorm:"column:sku_name"`         // SKU名称快照
	ImageURL    string  `json:"imageUrl" gorm:"column:image_url"`       // 商品图片快照
	Price       float64 `json:"price" gorm:"column:price"`              // 商品单价快照
	Quantity    int32   `json:"quantity" gorm:"column:quantity"`        // 购买数量
	TotalAmount float64 `json:"totalAmount" gorm:"column:total_amount"` // 商品小计金额
	DeptID      int64   `json:"deptId" gorm:"column:dept_id"`           // 部门ID
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"` // 操作状态
}

// OrderSetReq 订单创建/更新请求
type OrderSetReq struct {
	ID                    int64          `json:"id,string"`                                // 主键ID（更新时传入）
	UserID                int64          `json:"-"`                                        // 用户ID（从上下文获取）
	ReceiverName          string         `json:"receiverName" binding:"required"`          // 收货人姓名
	ReceiverPhone         string         `json:"receiverPhone" binding:"required"`         // 收货人手机号
	ReceiverProvince      string         `json:"receiverProvince"`                         // 省份
	ReceiverCity          string         `json:"receiverCity"`                             // 城市
	ReceiverDistrict      string         `json:"receiverDistrict"`                         // 区县
	ReceiverDetailAddress string         `json:"receiverDetailAddress" binding:"required"` // 详细地址
	Remark                string         `json:"remark"`                                   // 订单备注
	Items                 []OrderItemReq `json:"items" binding:"required"`                 // 订单商品明细
}

// OrderItemReq 订单商品明细请求
type OrderItemReq struct {
	GoodsID   string  `json:"goodsId" binding:"required"`   // 商品ID
	SkuID     string  `json:"skuId" binding:"required"`     // SKU ID
	GoodsName string  `json:"goodsName" binding:"required"` // 商品名称
	SkuName   string  `json:"skuName"`                      // SKU名称
	ImageURL  string  `json:"imageUrl"`                     // 图片
	Price     float64 `json:"price" binding:"required"`     // 单价
	Quantity  int32   `json:"quantity" binding:"required"`  // 数量
}

// OrderQuery 订单查询参数
type OrderQuery struct {
	UserID  int64  `form:"userId"`  // 用户ID
	Status  *int32 `form:"status"`  // 订单状态
	OrderNo string `form:"orderNo"` // 订单号
	Page    int64  `form:"page"`    // 页码
	Size    int64  `form:"size"`    // 每页数量
}

// OrderListData 订单列表结果
type OrderListData struct {
	Rows  []*ApiOrderVO `json:"rows"`  // 数据列表
	Total int64         `json:"total"` // 总数
}

// OrderVO 订单视图对象（包含商品明细）
type OrderVO struct {
	shopordermodels.Order
	Items []*OrderItem `json:"items"` // 订单商品明细
}

// OrderStatusReq 订单状态更新请求
type OrderStatusReq struct {
	ID     int64  `json:"id,string" binding:"required"` // 订单ID
	Status int32  `json:"status" binding:"required"`    // 目标状态
	Reason string `json:"reason"`                       // 变更原因（如取消原因）
}

// OrderStatistics 订单统计
type OrderStatistics struct {
	PendingPay     int64 `json:"pendingPay"`     // 待付款
	PendingSend    int64 `json:"pendingSend"`    // 待发货
	PendingReceive int64 `json:"pendingReceive"` // 待收货
	Completed      int64 `json:"completed"`      // 已完成
	Cancelled      int64 `json:"cancelled"`      // 已取消
}

// OrderPayResp 支付响应（微信小程序调起支付参数）
type OrderPayResp struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

func ToShopOrderVO(order *shopordermodels.Order) *OrderVO {
	if order == nil {
		return nil
	}
	items := make([]*OrderItem, 0, len(order.Details))
	for _, detail := range order.Details {
		if detail == nil {
			continue
		}
		items = append(items, toShopOrderItem(detail))
	}
	return &OrderVO{
		Order: *order,
		Items: items,
	}
}

func toShopOrderItem(detail *shopordermodels.OrderDetail) *OrderItem {
	if detail == nil {
		return nil
	}
	price := detail.Payment
	if detail.Num > 0 {
		price = detail.Payment / detail.Num
	}
	return &OrderItem{
		ID:          int64(detail.ID),
		OrderID:     int64(detail.OrderID),
		OrderNo:     detail.Tid,
		GoodsID:     stringUtils.FirstNonEmpty(detail.EShopGoodsID, fmt.Sprintf("%d", detail.NumIID)),
		SkuID:       stringUtils.FirstNonEmpty(detail.EShopSkuID, fmt.Sprintf("%d", detail.SkuID)),
		GoodsName:   detail.EShopGoodsName,
		SkuName:     detail.EShopSkuName,
		ImageURL:    detail.PicPath,
		Price:       price,
		Quantity:    int32(detail.Num),
		TotalAmount: detail.Payment,
		DeptID:      detail.DeptID,
		State:       detail.State,
		BaseEntity:  detail.BaseEntity,
	}
}

func BuildOrderBuyerNick(shopUser *shopusermodels.User) string {
	if shopUser == nil {
		return ""
	}
	return fmt.Sprintf("shop-user-%s", shopUser.UserID)
}
