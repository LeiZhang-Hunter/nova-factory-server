package models

// Pink 拼团记录
type Pink struct {
	ID               int64   `json:"id,string" gorm:"column:id"`
	UID              int64   `json:"uid,string" gorm:"column:uid"`
	Nickname         string  `json:"nickname" gorm:"column:nickname"`
	Avatar           string  `json:"avatar" gorm:"column:avatar"`
	OrderID          string  `json:"orderId" gorm:"column:order_id"`
	OrderIDKey       int64   `json:"orderIdKey,string" gorm:"column:order_id_key"`
	TotalNum         int64   `json:"totalNum" gorm:"column:total_num"`
	TotalPrice       float64 `json:"totalPrice" gorm:"column:total_price"`
	CID              int64   `json:"cid,string" gorm:"column:cid"`
	PID              int64   `json:"pid,string" gorm:"column:pid"`
	People           int64   `json:"people" gorm:"column:people"`
	Price            float64 `json:"price" gorm:"column:price"`
	AddTime          string  `json:"addTime" gorm:"column:add_time"`
	StopTime         string  `json:"stopTime" gorm:"column:stop_time"`
	KID              int64   `json:"kId,string" gorm:"column:k_id"`
	IsTpl            int32   `json:"isTpl" gorm:"column:is_tpl"`
	IsRefund         int32   `json:"isRefund" gorm:"column:is_refund"`
	Status           int32   `json:"status" gorm:"column:status"`
	IsVirtual        int32   `json:"isVirtual" gorm:"column:is_virtual"`
	CombinationTitle string  `json:"combinationTitle" gorm:"column:combination_title"`
	CombinationImage string  `json:"combinationImage" gorm:"column:combination_image"`
}

// PinkQuery 拼团记录查询参数
type PinkQuery struct {
	OrderID  string `form:"orderId"`
	CID      int64  `form:"cid"`
	UID      int64  `form:"uid"`
	Status   *int32 `form:"status"`
	IsRefund *int32 `form:"isRefund"`
	Page     int64  `form:"page"`
	Size     int64  `form:"size"`
}

// PinkListData 拼团记录列表数据
type PinkListData struct {
	Rows  []*Pink `json:"rows"`
	Total int64   `json:"total"`
}
