package systemdao

import (
	modelentity "nova-factory-server/app/business/admin/system/systemmodels/entity"
	modelquery "nova-factory-server/app/business/admin/system/systemmodels/query"
	modelrequest "nova-factory-server/app/business/admin/system/systemmodels/request"
	modelresponse "nova-factory-server/app/business/admin/system/systemmodels/response"

	"github.com/gin-gonic/gin"
)

type ISysShiftDao interface {
	Set(c *gin.Context, valueVO *modelrequest.SysWorkShiftSettingVO) (*modelentity.SysWorkShiftSetting, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *modelquery.SysWorkShiftSettingReq) (*modelresponse.SysWorkShiftSettingList, error)
	// Check 校验班次时间，防止重复
	Check(c *gin.Context, id int64, startTime int32, endTime int32) *modelentity.SysWorkShiftSetting
	// GetEnableShift 读取启用班次
	GetEnableShift(c *gin.Context) ([]*modelentity.SysWorkShiftSetting, error)
}
