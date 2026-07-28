package query

import "nova-factory-server/app/baize"

// SysWorkShiftSettingReq 班次配置列表
type SysWorkShiftSettingReq struct {
	Name string `form:"name" ` //排序规则  降序desc   asc升序
	baize.BaseEntityDQL
}
