package models

import (
	"fmt"
	"nova-factory-server/app/baize"
	shopordermodels "nova-factory-server/app/business/shop/order/models"
	orderConstant "nova-factory-server/app/constant/order"
	"nova-factory-server/app/utils/stringUtils"
	"time"
)

// ApiOrderVO 订单视图对象（包含商品明细）
type ApiOrderVO struct {
	Order
	Items []*OrderItem `json:"items"` // 订单商品明细
}

// Order 订单主表
type Order struct {
	ID                    int64      `json:"id,string" gorm:"column:id"`                                  // 主键ID
	OrderNo               string     `json:"orderNo" gorm:"column:order_no"`                              // 订单号
	UserID                int64      `json:"userId,string" gorm:"column:user_id"`                         // 用户ID
	TotalAmount           float64    `json:"totalAmount" gorm:"column:total_amount"`                      // 订单总金额
	PayAmount             float64    `json:"payAmount" gorm:"column:pay_amount"`                          // 实付金额
	FreightAmount         float64    `json:"freightAmount" gorm:"column:freight_amount"`                  // 运费金额
	DiscountAmount        float64    `json:"discountAmount" gorm:"column:discount_amount"`                // 优惠金额
	Status                int32      `json:"status" gorm:"column:status"`                                 // 订单状态：0待支付，1已支付，2已发货，3已完成，4已取消
	PayTime               *time.Time `json:"payTime" gorm:"column:pay_time"`                              // 支付时间
	ShipTime              *time.Time `json:"shipTime" gorm:"column:ship_time"`                            // 发货时间
	CompleteTime          *time.Time `json:"completeTime" gorm:"column:complete_time"`                    // 完成时间
	CancelTime            *time.Time `json:"cancelTime" gorm:"column:cancel_time"`                        // 取消时间
	CancelReason          string     `json:"cancelReason" gorm:"column:cancel_reason"`                    // 取消原因
	ReceiverName          string     `json:"receiverName" gorm:"column:receiver_name"`                    // 收货人姓名
	ReceiverPhone         string     `json:"receiverPhone" gorm:"column:receiver_phone"`                  // 收货人手机号
	ReceiverProvince      string     `json:"receiverProvince" gorm:"column:receiver_province"`            // 省份
	ReceiverCity          string     `json:"receiverCity" gorm:"column:receiver_city"`                    // 城市
	ReceiverDistrict      string     `json:"receiverDistrict" gorm:"column:receiver_district"`            // 区县
	ReceiverDetailAddress string     `json:"receiverDetailAddress" gorm:"column:receiver_detail_address"` // 详细地址
	Remark                string     `json:"remark" gorm:"column:remark"`                                 // 订单备注
	Version               int32      `json:"version" gorm:"column:version"`                               // 乐观锁版本号
	DeptID                int64      `json:"deptId" gorm:"column:dept_id"`                                // 部门ID
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"` // 操作状态
}

// ToApiShopOrderVO 订单转换为 api订单,为了兼容小程序数据结构返回，历史原因
func ToApiShopOrderVO(order *shopordermodels.Order) *ApiOrderVO {
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
	return &ApiOrderVO{
		Order: *ToShopOrder(order),
		Items: items,
	}
}

// ToShopOrder 转换为 api shop订单
func ToShopOrder(order *shopordermodels.Order) *Order {
	if order == nil {
		return nil
	}
	return &Order{
		ID:                    int64(order.ID),
		OrderNo:               order.Tid,
		TotalAmount:           order.Total,
		PayAmount:             order.Total - order.Privilege + order.PostFee,
		FreightAmount:         order.PostFee,
		DiscountAmount:        order.Privilege,
		Status:                orderConstant.OrderStatusToShopStatus(order.Status),
		PayTime:               order.PayTime,
		ReceiverName:          order.ReceiverName,
		ReceiverPhone:         stringUtils.FirstNonEmpty(order.ReceiverMobile, order.ReceiverPhone),
		ReceiverProvince:      stringUtils.FirstNonEmpty(order.ReceiverProvinceName, order.ReceiverProvince),
		ReceiverCity:          stringUtils.FirstNonEmpty(order.ReceiverCityName, order.ReceiverCity),
		ReceiverDistrict:      stringUtils.FirstNonEmpty(order.ReceiverDistrictName, order.ReceiverDistrict),
		ReceiverDetailAddress: order.ReceiverAddress,
		Remark:                stringUtils.FirstNonEmpty(order.SellerMemo, order.BuyerMessage),
	}
}

func ToShopOrderItem(detail *shopordermodels.OrderDetail) *OrderItem {
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
