package systemdao

import (
	"nova-factory-server/app/business/admin/system/systemmodels"

	"github.com/gin-gonic/gin"
)

type ISysShiftDao interface {
	Set(c *gin.Context, valueVO *systemmodels.SysWorkShiftSettingVO) (*systemmodels.SysWorkShiftSetting, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *systemmodels.SysWorkShiftSettingReq) (*systemmodels.SysWorkShiftSettingList, error)
	// Check 校验班次时间，防止重复
	Check(c *gin.Context, id int64, startTime int32, endTime int32) *systemmodels.SysWorkShiftSetting
	// GetEnableShift 读取启用班次
	GetEnableShift(c *gin.Context) ([]*systemmodels.SysWorkShiftSetting, error)
}
