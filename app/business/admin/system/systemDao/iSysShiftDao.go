package systemDao

import (
	"nova-factory-server/app/business/admin/system/systemModels"

	"github.com/gin-gonic/gin"
)

type ISysShiftDao interface {
	Set(c *gin.Context, valueVO *systemModels.SysWorkShiftSettingVO) (*systemModels.SysWorkShiftSetting, error)
	Remove(c *gin.Context, ids []string) error
	List(c *gin.Context, req *systemModels.SysWorkShiftSettingReq) (*systemModels.SysWorkShiftSettingList, error)
	// Check 校验班次时间，防止重复
	Check(c *gin.Context, id int64, startTime int32, endTime int32) *systemModels.SysWorkShiftSetting
	// GetEnableShift 读取启用班次
	GetEnableShift(c *gin.Context) ([]*systemModels.SysWorkShiftSetting, error)
}
