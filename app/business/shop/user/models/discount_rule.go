package models

import (
	"time"

	"nova-factory-server/app/baize"
)

// DateTimeYYYYMMDDHHMMSS 自定义时间格式：2006-01-02 15:04:05
type DateTimeYYYYMMDDHHMMSS struct {
	time.Time
}

func (t *DateTimeYYYYMMDDHHMMSS) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" || str == `""` {
		return nil
	}
	// 去掉首尾引号
	str = str[1 : len(str)-1]
	if str == "" {
		return nil
	}
	parsed, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

// UserDiscountRule 用户折扣规则
type UserDiscountRule struct {
	ID           int64      `json:"id,string" gorm:"primaryKey;autoIncrement"`      // 主键
	UserID       int64      `json:"userId,string" gorm:"index;not null"`            // 用户ID
	TargetType   string     `json:"targetType" gorm:"type:varchar(20);not null"`    // sku=SKU折扣, category=分类折扣
	TargetID     string     `json:"targetId" gorm:"index;not null"`                 // SKU ID 或分类 ID
	DiscountRate float64    `json:"discountRate" gorm:"type:decimal(5,2);not null"` // 折扣率 (0.85 = 85折)
	ValidFrom    *time.Time `json:"validFrom" gorm:"type:datetime"`                 // 有效期开始
	ValidTo      *time.Time `json:"validTo" gorm:"type:datetime"`                   // 有效期结束
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"` // 状态
}

// TableName 表名
func (UserDiscountRule) TableName() string {
	return "shop_user_discount_rules"
}

// UserDiscountRuleUpsert 用户折扣规则新增修改参数
type UserDiscountRuleUpsert struct {
	ID           int64  `json:"id,string"`     // 主键ID
	UserID       int64  `json:"userId,string"` // 用户ID
	TargetType   string `json:"targetType"`    // goods=商品, category=分类
	TargetID     string `json:"targetId"`      // 商品ID或分类ID
	DiscountRate int64  `json:"discountRate"`  // 折扣率
	ValidFrom    string `json:"validFrom"`     // 有效期开始 (格式: 2006-01-02 15:04:05)
	ValidTo      string `json:"validTo"`       // 有效期结束 (格式: 2006-01-02 15:04:05)
}

// UserDiscountRuleQuery 用户折扣规则查询参数
type UserDiscountRuleQuery struct {
	UserID     int64  `form:"userId"`     // 用户ID
	TargetType string `form:"targetType"` // 类型: goods, category
	TargetID   string `form:"targetId"`   // 商品ID或分类ID
	Page       int64  `form:"page"`       // 页码
	Size       int64  `form:"size"`       // 每页数量
}

// UserDiscountRuleListData 用户折扣规则列表结果
type UserDiscountRuleListData struct {
	Rows  []*UserDiscountRule `json:"rows"`  // 数据列表
	Total int64               `json:"total"` // 总数
}

// BatchDiscountRuleCreate 批量创建折扣规则请求
type BatchDiscountRuleCreate struct {
	UserIDs      []int64    `json:"userIds"`      // 用户ID列表
	TargetType   string     `json:"targetType"`   // goods=商品, category=分类
	TargetIDs    []string   `json:"targetIds"`    // 商品ID或分类ID列表
	DiscountRate float64    `json:"discountRate"` // 折扣率
	ValidFrom    *time.Time `json:"validFrom"`    // 有效期开始
	ValidTo      *time.Time `json:"validTo"`      // 有效期结束
}
