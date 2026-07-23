package response

import modelentity "nova-factory-server/app/business/admin/system/systemmodels/entity"

type SysWorkShiftSettingList struct {
	Rows  []*modelentity.SysWorkShiftSetting `json:"rows"`
	Total int64                              `json:"total"`
}
