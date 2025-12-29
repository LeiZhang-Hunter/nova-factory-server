package productModels

import "nova-factory-server/app/baize"

type SysProductLaboratoryDQL struct {
	Material  string `form:"material" db:"material" jsonschema:"description=物料编码"`
	Number    string `form:"number" db:"number" jsonschema:"description=物料化验单样号"`
	Contact   string `form:"contact" db:"contact" jsonschema:"description=采样人"`
	Type      int    `form:"type" db:"type" jsonschema:"description=化验类型，0是物料化验，1是配料化验"`
	State     string `form:"state" db:"state"`
	BeginTime string `form:"beginTime" jsonschema:"description=开始时间，时间格式是:2006-01-02 15:04:05"`
	EndTime   string `form:"endTime" jsonschema:"description=结束时间，时间格式是:2006-01-02 15:04:05"`
	baize.BaseEntityDQL
}

type SysProductLaboratory struct {
	Id             int64   `json:"id,string" db:"id"`
	Material       string  `json:"material" db:"material" binding:"required"`
	DryHigh        float64 `gorm:"column:dry_high;not null;default:0.00;comment:干燥基高位" json:"dry_high"`         // 干燥基高位
	ReceivedLow    float64 `gorm:"column:received_low;not null;default:0.00;comment:收到基低位" json:"received_low"` // 收到基低位
	Contact        string  `json:"contact" db:"contact" binding:"required"`
	Date           string  `json:"date" db:"string" `
	Address        string  `json:"address" db:"address" binding:"required"`
	Heat           float64 `json:"heat" db:"heat"`
	Sulphur        float64 `json:"sulphur" db:"sulphur"`
	Volatility     float64 `json:"volatility" db:"volatility"`
	Water          float64 `json:"water" db:"water"`
	Weight         float64 `json:"weight" db:"weight"`
	Number         string  `json:"number" db:"number" binding:"required"`
	Img            string  `json:"img" db:"img"`
	Type           int     `form:"type" db:"type" comment:"化验类型，0是物料化验，1是配料化验"`
	DeptId         int64   `json:"deptId" db:"dept_id"`
	State          bool    `json:"state" db:"state"`
	CreateUserName string  `json:"createUserName" gorm:"-"`
	UpdateUserName string  `json:"updateUserName" gorm:"-"`
	baize.BaseEntity
}

func ToSysProductLaboratory(vo *SysProductLaboratoryVo) *SysProductLaboratory {
	return &SysProductLaboratory{
		Id:          vo.Id,
		Material:    vo.Material,
		Date:        vo.Date,
		DryHigh:     vo.DryHigh,
		ReceivedLow: vo.ReceivedLow,
		Contact:     vo.Contact,
		Address:     vo.Address,
		Heat:        vo.Heat,
		Sulphur:     vo.Sulphur,
		Volatility:  vo.Volatility,
		Water:       vo.Water,
		Weight:      vo.Weight,
		Number:      vo.Number,
		Img:         vo.Img,
		Type:        vo.Type,
	}
}

type SysProductLaboratoryVo struct {
	Id          int64   `json:"id,string" db:"id"`
	Type        int     `form:"type" db:"type" comment:"化验类型，0是物料化验，1是配料化验"`
	Material    string  `json:"material" db:"material" binding:"required"`
	DryHigh     float64 `gorm:"column:dry_high;not null;default:0.00;comment:干燥基高位" json:"dry_high"`         // 干燥基高位
	ReceivedLow float64 `gorm:"column:received_low;not null;default:0.00;comment:收到基低位" json:"received_low"` // 收到基低位
	Contact     string  `json:"contact" db:"contact" `
	Date        string  `json:"date" db:"string" `
	Address     string  `json:"address" db:"address"`
	Heat        float64 `json:"heat" db:"heat"`
	Sulphur     float64 `json:"sulphur" db:"sulphur"`
	Volatility  float64 `json:"volatility" db:"volatility"`
	Water       float64 `json:"water" db:"water"`
	Weight      float64 `json:"weight" db:"weight"`
	Number      string  `json:"number" db:"number" binding:"required"`
	Img         string  `json:"img" db:"img"`
}

type SysProductLaboratoryList struct {
	Rows  []*SysProductLaboratory `json:"rows"`
	Total int64                   `json:"total"`
}

type SysProductLaboratoryInfoDQL struct {
	Type int `form:"type" db:"type" jsonschema:"description=化验类型，0是物料化验，1是配料化验"`
}
