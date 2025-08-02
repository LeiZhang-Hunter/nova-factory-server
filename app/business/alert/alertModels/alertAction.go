package alertModels

import (
	"encoding/json"
	"go.uber.org/zap"
	"nova-factory-server/app/baize"
)

type Receiver struct {
	UserId      int64  `json:"userId,string" db:"user_id"swaggerignore:"true"` //用户ID
	UserName    string `json:"userName" db:"user_name" binding:"required"`     //用户名
	NickName    string `json:"nickName" db:"nick_name" binding:"required"`     //用户昵称
	Email       string `json:"email" db:"email"`                               //邮箱
	Phonenumber string `json:"phonenumber" db:"phonenumber"`                   //手机号
}

type UserNotify struct {
	Receiver  []Receiver `json:"receiver"`
	Period    []string   `json:"period"`
	TimeStart uint64     `gorm:"column:time_start;comment:通知开始时间" json:"time_start"` // 通知开始时间
	TimeEnd   uint64     `gorm:"column:time_end;comment:通知结束时间" json:"time_end"`     // 通知结束时间
	TimeRange []string   `json:"time_range"`                                         // 通知结束时间
	Channels  []string   `json:"channels"`
}

type ApiNotify struct {
	Url       string   `json:"url"`
	Period    []string `json:"period"`
	TimeStart uint64   `json:"time_start"` // 通知开始时间
	TimeEnd   uint64   `json:"time_end"`   // 通知结束时间
	TimeRange []string `json:"time_range"` // 通知结束时间
}

// AlertAction 告警行为
type AlertAction struct {
	ID             int64        `gorm:"column:id;primaryKey;autoIncrement:true;comment:主键" json:"id,string"` // 主键
	Name           string       `gorm:"column:name;not null;comment:告警策略名称" json:"name"`                     // 告警策略名称
	UserNotifyList []UserNotify `gorm:"-" json:"user_notify"`                                                // 通知周期，逗号分隔
	ApiNotifyList  []ApiNotify  `gorm:"-" json:"api_notify"`                                                 // 通知用户
	UserNotify     string       `gorm:"column:user_notify;comment:通知用户" json:"-"`                            // 通知用户
	ApiNotify      string       `gorm:"column:api_notify;comment:回调通知" json:"-"`                             // 回调通知
	UserCount      uint32       `gorm:"-" json:"user_count"`                                                 // 主键
	baize.BaseEntity
	DeptID int64 `gorm:"column:dept_id;comment:部门ID" json:"dept_id,string"` // 部门ID
	State  uint8 `gorm:"column:state;not null;default:0" json:"state"`
}

func FromSetAlertActionToData(action *SetAlertAction) *AlertAction {
	var err error
	var userNotify []byte
	var apiNotify []byte
	if len(action.UserNotifyList) != 0 {
		userNotify, err = json.Marshal(action.UserNotifyList)
		if err != nil {
			zap.L().Error("json marshal", zap.Error(err))
		}
	} else {
		userNotify = []byte("[]")
	}
	if len(action.ApiNotifyList) != 0 {
		apiNotify, err = json.Marshal(action.ApiNotifyList)
		if err != nil {
			zap.L().Error("json marshal", zap.Error(err))
		}
	} else {
		apiNotify = []byte("[]")
	}
	return &AlertAction{
		ID:         action.ID,
		Name:       action.Name,
		UserNotify: string(userNotify),
		ApiNotify:  string(apiNotify),
	}
}

type SetAlertAction struct {
	ID             int64        `gorm:"column:id;primaryKey;autoIncrement:true;comment:主键" json:"id,string"` // 主键
	Name           string       `gorm:"column:name;not null;comment:告警策略名称" json:"name"`                     // 告警策略名称
	UserNotifyList []UserNotify `json:"user_notify"`                                                         // 通知周期，逗号分隔
	ApiNotifyList  []ApiNotify  `json:"api_notify"`                                                          // 通知用户
}

type SysAlertActionListReq struct {
	Name string `form:"name"` // 告警策略名称
	baize.BaseEntityDQL
}

type SysAlertActionList struct {
	Rows  []*AlertAction `json:"rows"`
	Total uint64         `json:"total"`
}
