package systemServiceImpl

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/business/system/systemModels"
	"nova-factory-server/app/business/system/systemService"
)

type ISysShiftServiceImpl struct {
	dao systemDao.ISysShiftDao
}

func NewISysShiftServiceImpl(dao systemDao.ISysShiftDao) systemService.ISysShiftService {
	return &ISysShiftServiceImpl{
		dao: dao,
	}
}

func (i *ISysShiftServiceImpl) Set(c *gin.Context, valueVO *systemModels.SysWorkShiftSettingVO) (*systemModels.SysWorkShiftSetting, error) {
	return i.dao.Set(c, valueVO)
}
func (i *ISysShiftServiceImpl) Remove(c *gin.Context, ids []string) error {
	return i.dao.Remove(c, ids)
}
func (i *ISysShiftServiceImpl) List(c *gin.Context, req *systemModels.SysWorkShiftSettingReq) (*systemModels.SysWorkShiftSettingList, error) {
	return i.dao.List(c, req)
}

func (i *ISysShiftServiceImpl) Check(c *gin.Context, id int64, startTime int32, endTime int32) *systemModels.SysWorkShiftSetting {
	return i.dao.Check(c, id, startTime, endTime)
}
