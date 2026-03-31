package systemServiceImpl

import (
	"nova-factory-server/app/business/admin/system/systemdao"
	"nova-factory-server/app/business/admin/system/systemmodels"
	"nova-factory-server/app/business/admin/system/systemservice"

	"github.com/gin-gonic/gin"
)

type ISysShiftServiceImpl struct {
	dao systemdao.ISysShiftDao
}

func NewISysShiftServiceImpl(dao systemdao.ISysShiftDao) systemservice.ISysShiftService {
	return &ISysShiftServiceImpl{
		dao: dao,
	}
}

func (i *ISysShiftServiceImpl) Set(c *gin.Context, valueVO *systemmodels.SysWorkShiftSettingVO) (*systemmodels.SysWorkShiftSetting, error) {
	return i.dao.Set(c, valueVO)
}
func (i *ISysShiftServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
func (i *ISysShiftServiceImpl) List(c *gin.Context, req *systemmodels.SysWorkShiftSettingReq) (*systemmodels.SysWorkShiftSettingList, error) {
	return i.dao.List(c, req)
}

func (i *ISysShiftServiceImpl) Check(c *gin.Context, id int64, startTime int32, endTime int32) *systemmodels.SysWorkShiftSetting {
	return i.dao.Check(c, id, startTime, endTime)
}
