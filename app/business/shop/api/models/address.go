package models

import (
	"nova-factory-server/app/baize"
)

// // SnowflakeID 自定义类型，用于处理雪花算法的 64 位 ID 在 JSON 序列化时的精度丢失问题
// // JSON 输出时序列化为字符串，解析时支持字符串和数字两种格式
// type SnowflakeID int64
//
//	func (f SnowflakeID) MarshalJSON() ([]byte, error) {
//		return json.Marshal(strconv.FormatInt(int64(f), 10))
//	}
//
//	func (f *SnowflakeID) UnmarshalJSON(data []byte) error {
//		// 尝试解析为字符串
//		var s string
//		if err := json.Unmarshal(data, &s); err == nil {
//			v, err := strconv.ParseInt(s, 10, 64)
//			if err != nil {
//				return err
//			}
//			*f = SnowflakeID(v)
//			return nil
//		}
//		// 尝试解析为数字
//		var n int64
//		if err := json.Unmarshal(data, &n); err != nil {
//			return err
//		}
//		*f = SnowflakeID(n)
//		return nil
//	}
//
// ShopUserAddressApp 移动端用户地址
type ShopUserAddressApp struct {
	ID             int64  `json:"id,string" db:"id"`
	UserID         int64  `json:"userId,string" db:"user_id"`
	ReceiverName   string `json:"receiverName" db:"receiver_name"`
	ReceiverMobile string `json:"receiverMobile" db:"receiver_mobile"`
	ProvinceCode   string `json:"provinceCode" db:"province_code"`
	ProvinceName   string `json:"provinceName" db:"province_name"`
	CityCode       string `json:"cityCode" db:"city_code"`
	CityName       string `json:"cityName" db:"city_name"`
	DistrictCode   string `json:"districtCode" db:"district_code"`
	DistrictName   string `json:"districtName" db:"district_name"`
	DetailAddress  string `json:"detailAddress" db:"detail_address"`
	AddressLabel   string `json:"addressLabel" db:"address_label"`
	IsDefault      int    `json:"isDefault" db:"is_default"` // 0-否 1-是
	baize.BaseEntity
}

// AddressSetReq 地址设置请求
type AddressSetReq struct {
	ID             int64  `json:"id,string"`
	UserID         int64  `json:"-"`
	ReceiverName   string `json:"receiverName" binding:"required"`
	ReceiverMobile string `json:"receiverMobile" binding:"required"`
	ProvinceCode   string `json:"provinceCode" binding:"required"`
	ProvinceName   string `json:"provinceName" binding:"required"`
	CityCode       string `json:"cityCode" binding:"required"`
	CityName       string `json:"cityName" binding:"required"`
	DistrictCode   string `json:"districtCode"`
	DistrictName   string `json:"districtName"`
	DetailAddress  string `json:"detailAddress" binding:"required"`
	AddressLabel   string `json:"addressLabel"`
	IsDefault      *bool  `json:"isDefault"`
}

// AddressQuery 地址查询参数
type AddressQuery struct {
	UserId int64 `form:"userId"`
	Page   int64 `form:"page"`
	Size   int64 `form:"size"`
}

// AddressListData 地址列表结果
type AddressListData struct {
	Rows  []*ShopUserAddressApp `json:"rows"`
	Total int64                 `json:"total"`
}

// AddressRegionItem 行政区节点
type AddressRegionItem struct {
	Code       string `json:"code"`
	Name       string `json:"name"`
	Level      string `json:"level"`
	ParentCode string `json:"parentCode"`
}
