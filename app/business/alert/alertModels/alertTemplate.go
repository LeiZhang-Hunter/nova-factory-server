package alertModels

import (
	"encoding/json"
	"go.uber.org/zap"
	"nova-factory-server/app/baize"
)

// SysAlertSinkTemplate 告警模板
type SysAlertSinkTemplate struct {
	ID        int64  `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增标识" json:"id,string"`                      // 自增标识
	Name      string `gorm:"column:name;not null;comment:告警模板名称" json:"name"`                                            // 告警模板名称
	Addr      string `gorm:"column:addr;not null;comment:发送alert的http地址，若为空，则不会发送" json:"addr"`                          // 发送alert的http地址，若为空，则不会发送
	Template  string `gorm:"column:template;comment:用来渲染的模板" json:"template"`                                            // 用来渲染的模板
	Timeout   int32  `gorm:"column:timeout;not null;comment:发送alert的http timeout" json:"timeout"`                        // 发送alert的http timeout
	Headers   string `gorm:"column:headers;comment:发送alert的http header" json:"headers"`                                  // 发送alert的http header
	Method    string `gorm:"column:method;not null;comment:发送alert的http method, 如果不填put(不区分大小写)，都认为是POST" json:"method"` // 发送alert的http method, 如果不填put(不区分大小写)，都认为是POST
	Extension string `gorm:"column:extension;comment:扩展信息" json:"extension"`                                             // 扩展信息
	DeptID    int64  `gorm:"column:dept_id;comment:部门ID" json:"dept_id"`                                                 // 部门ID
	baize.BaseEntity
	State bool `gorm:"column:state;comment:操作状态（0正常 -1删除）" json:"state"` // 操作状态（0正常 -1删除）
}

type SetSysAlertSinkTemplate struct {
	ID        int64             `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增标识" json:"id,string"`                      // 自增标识
	Name      string            `gorm:"column:name;not null;comment:告警模板名称" json:"name" binding:"required"`                         // 告警模板名称
	Addr      string            `gorm:"column:addr;not null;comment:发送alert的http地址，若为空，则不会发送" binding:"required" json:"addr"`       // 发送alert的http地址，若为空，则不会发送
	Template  string            `gorm:"column:template;comment:用来渲染的模板" json:"template" binding:"required"`                         // 用来渲染的模板
	Timeout   int32             `gorm:"column:timeout;not null;comment:发送alert的http timeout" json:"timeout"`                        // 发送alert的http timeout
	Headers   map[string]string `gorm:"column:headers;comment:发送alert的http header" json:"headers"`                                  // 发送alert的http header
	Method    string            `gorm:"column:method;not null;comment:发送alert的http method, 如果不填put(不区分大小写)，都认为是POST" json:"method"` // 发送alert的http method, 如果不填put(不区分大小写)，都认为是POST
	Extension string            `gorm:"column:extension;comment:扩展信息" json:"extension"`                                             // 扩展信息
}

func ToSysAlertSinkTemplate(data *SetSysAlertSinkTemplate) *SysAlertSinkTemplate {
	var headerStr []byte
	var err error
	if len(data.Headers) != 0 {
		headerStr, err = json.Marshal(data.Headers)
		if err != nil {
			zap.L().Error("json marshal error", zap.Error(err))
		}
	}
	return &SysAlertSinkTemplate{
		ID:        data.ID,
		Name:      data.Name,
		Addr:      data.Addr,
		Template:  data.Template,
		Timeout:   data.Timeout,
		Headers:   string(headerStr),
		Method:    data.Method,
		Extension: data.Extension,
	}
}

type SysAlertSinkTemplateReq struct {
	GatewayID int64  `form:"gateway_id"` // 网关id
	Name      string `form:"name"`       // 告警模板名称
	Addr      string `form:"addr"`
	baize.BaseEntityDQL
}

type SysAlertSinkTemplateListData struct {
	Rows  []*SysAlertSinkTemplate `json:"rows"`
	Total uint64                  `json:"total"`
}
