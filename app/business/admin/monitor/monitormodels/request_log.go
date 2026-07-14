package monitormodels

import "time"

type RequestLog struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id,string"`
	LogType       string    `gorm:"column:log_type" json:"logType" form:"logType"`
	SourceModule  string    `gorm:"column:source_module" json:"sourceModule" form:"sourceModule"`
	Status        int32     `gorm:"column:status" json:"status" form:"status"`
	RequestMethod string    `gorm:"column:request_method" json:"requestMethod"`
	RequestPath   string    `gorm:"column:request_path" json:"requestPath"`
	QueryString   string    `gorm:"column:query_string" json:"queryString"`
	HeadersJSON   string    `gorm:"column:headers_json" json:"headersJson"`
	BodyText      string    `gorm:"column:body_text" json:"bodyText"`
	ClientIP      string    `gorm:"column:client_ip" json:"clientIp"`
	UserAgent     string    `gorm:"column:user_agent" json:"userAgent"`
	ErrorMessage  string    `gorm:"column:error_message" json:"errorMessage"`
	CreateTime    time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime    time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

type RequestLogQuery struct {
	PageNum   int    `form:"pageNum"`
	PageSize  int    `form:"pageSize"`
	LogType   string `form:"logType"`
	Status    *int32 `form:"status"`
	BeginTime string `form:"beginTime"`
	EndTime   string `form:"endTime"`
}

type RequestLogListData struct {
	Rows  []*RequestLog `json:"rows"`
	Total int64         `json:"total"`
}
