package models

import "nova-factory-server/app/baize"

// Order 订单主表
type Order struct {
	ID                   int64   `json:"id,string" db:"id"`                                      // 主键ID
	OrderNo              string  `json:"orderNo" db:"order_no"`                                 // 订单号
	UserID               int64   `json:"userId" db:"user_id"`                                   // 用户ID
	TotalAmount          float64 `json:"totalAmount" db:"total_amount"`                         // 订单总金额
	PayAmount            float64 `json:"payAmount" db:"pay_amount"`                             // 实付金额
	FreightAmount        float64 `json:"freightAmount" db:"freight_amount"`                      // 运费金额
	DiscountAmount       float64 `json:"discountAmount" db:"discount_amount"`                    // 优惠金额
	Status               int32   `json:"status" db:"status"`                                    // 订单状态：0待支付，1已支付，2已发货，3已完成，4已取消
	PayTime              string  `json:"payTime" db:"pay_time"`                                 // 支付时间
	ShipTime             string  `json:"shipTime" db:"ship_time"`                               // 发货时间
	CompleteTime         string  `json:"completeTime" db:"complete_time"`                        // 完成时间
	CancelTime           string  `json:"cancelTime" db:"cancel_time"`                           // 取消时间
	CancelReason         string  `json:"cancelReason" db:"cancel_reason"`                        // 取消原因
	ReceiverName         string  `json:"receiverName" db:"receiver_name"`                        // 收货人姓名
	ReceiverPhone        string  `json:"receiverPhone" db:"receiver_phone"`                      // 收货人手机号
	ReceiverProvince     string  `json:"receiverProvince" db:"receiver_province"`                // 省份
	ReceiverCity         string  `json:"receiverCity" db:"receiver_city"`                        // 城市
	ReceiverDistrict     string  `json:"receiverDistrict" db:"receiver_district"`                // 区县
	ReceiverDetailAddress string `json:"receiverDetailAddress" db:"receiver_detail_address"`      // 详细地址
	Remark               string  `json:"remark" db:"remark"`                                    // 订单备注
	Version              int32   `json:"version" db:"version"`                                  // 乐观锁版本号
	DeptID               int64   `json:"deptId" db:"dept_id"`                                   // 部门ID
	baize.BaseEntity
	State int32 `json:"state" db:"state"` // 操作状态
}

// OrderItem 订单商品明细
type OrderItem struct {
	ID         int64   `json:"id,string" db:"id"`                                  // 主键ID
	OrderID    int64   `json:"orderId" db:"order_id"`                             // 订单ID
	OrderNo    string  `json:"orderNo" db:"order_no"`                             // 订单号
	GoodsID    string  `json:"goodsId" db:"goods_id"`                            // 商品ID
	SkuID      string  `json:"skuId" db:"sku_id"`                                 // SKU ID
	GoodsName  string  `json:"goodsName" db:"goods_name"`                        // 商品名称快照
	SkuName    string  `json:"skuName" db:"sku_name"`                             // SKU名称快照
	ImageURL   string  `json:"imageUrl" db:"image_url"`                           // 商品图片快照
	Price      float64 `json:"price" db:"price"`                                 // 商品单价快照
	Quantity   int32   `json:"quantity" db:"quantity"`                           // 购买数量
	TotalAmount float64 `json:"totalAmount" db:"total_amount"`                    // 商品小计金额
	DeptID     int64   `json:"deptId" db:"dept_id"`                               // 部门ID
	baize.BaseEntity
	State int32 `json:"state" db:"state"` // 操作状态
}

// OrderSetReq 订单创建/更新请求
type OrderSetReq struct {
	ID                    int64   `json:"id,string"`                                // 主键ID（更新时传入）
	UserID                int64   `json:"-"`                                        // 用户ID（从上下文获取）
	ReceiverName          string  `json:"receiverName" binding:"required"`           // 收货人姓名
	ReceiverPhone         string  `json:"receiverPhone" binding:"required"`         // 收货人手机号
	ReceiverProvince      string  `json:"receiverProvince"`                          // 省份
	ReceiverCity          string  `json:"receiverCity"`                              // 城市
	ReceiverDistrict      string  `json:"receiverDistrict"`                          // 区县
	ReceiverDetailAddress string  `json:"receiverDetailAddress" binding:"required"` // 详细地址
	Remark                string  `json:"remark"`                                   // 订单备注
	Items                 []OrderItemReq `json:"items" binding:"required"`          // 订单商品明细
}

// OrderItemReq 订单商品明细请求
type OrderItemReq struct {
	GoodsID    string  `json:"goodsId" binding:"required"`  // 商品ID
	SkuID      string  `json:"skuId" binding:"required"`    // SKU ID
	GoodsName  string  `json:"goodsName" binding:"required"` // 商品名称
	SkuName    string  `json:"skuName"`                     // SKU名称
	ImageURL   string  `json:"imageUrl"`                    // 图片
	Price      float64 `json:"price" binding:"required"`     // 单价
	Quantity   int32   `json:"quantity" binding:"required"` // 数量
}

// OrderQuery 订单查询参数
type OrderQuery struct {
	UserID  int64  `form:"userId"`  // 用户ID
	Status  *int32 `form:"status"`  // 订单状态
	OrderNo string `form:"orderNo"`  // 订单号
	Page    int64  `form:"page"`    // 页码
	Size    int64  `form:"size"`    // 每页数量
}

// OrderListData 订单列表结果
type OrderListData struct {
	Rows  []*OrderVO `json:"rows"`  // 数据列表
	Total int64      `json:"total"` // 总数
}

// OrderVO 订单视图对象（包含商品明细）
type OrderVO struct {
	Order
	Items []*OrderItem `json:"items"` // 订单商品明细
}

// OrderStatusReq 订单状态更新请求
type OrderStatusReq struct {
	ID     int64  `json:"id,string" binding:"required"` // 订单ID
	Status int32  `json:"status" binding:"required"`   // 目标状态
	Reason string `json:"reason"`                      // 变更原因（如取消原因）
}

// 订单状态常量
const (
	OrderStatusPending   int32 = 0 // 待支付
	OrderStatusPaid      int32 = 1 // 已支付
	OrderStatusShipped   int32 = 2 // 已发货
	OrderStatusCompleted int32 = 3 // 已完成
	OrderStatusCancelled int32 = 4 // 已取消
)

// GetStatusText 获取状态文本
func GetStatusText(status int32) string {
	switch status {
	case OrderStatusPending:
		return "待支付"
	case OrderStatusPaid:
		return "已支付"
	case OrderStatusShipped:
		return "已发货"
	case OrderStatusCompleted:
		return "已完成"
	case OrderStatusCancelled:
		return "已取消"
	default:
		return "未知"
	}
}

// OrderStatistics 订单统计
type OrderStatistics struct {
	PendingPay     int64 `json:"pendingPay"`      // 待付款
	PendingSend    int64 `json:"pendingSend"`     // 待发货
	PendingReceive int64 `json:"pendingReceive"` // 待收货
	Completed      int64 `json:"completed"`       // 已完成
	Cancelled      int64 `json:"cancelled"`       // 已取消
}
