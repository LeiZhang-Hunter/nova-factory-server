package deviceDaoImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type IDeviceSubjectDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIDeviceSubjectDaoImpl(db *gorm.DB) deviceDao.IDeviceSubjectDao {
	return &IDeviceSubjectDaoImpl{
		db:    db,
		table: "sys_device_subject",
	}
}

func (i *IDeviceSubjectDaoImpl) Set(c *gin.Context, data *deviceModels.SysDeviceSubjectVO) (*deviceModels.SysDeviceSubject, error) {
	if data == nil {
		return nil, errors.New("data should not be nil")
	}
	value := deviceModels.ToSysDeviceSubject(data)
	if data.ID != 0 {
		value.SetUpdateBy(baizeContext.GetUserId(c))
		ret := i.db.Table(i.table).Updates(data)
		return value, ret.Error
	} else {
		value.ID = snowflake.GenID()
		value.DeptID = baizeContext.GetDeptId(c)
		value.SetCreateBy(baizeContext.GetUserId(c))
		ret := i.db.Table(i.table).Where("id = ?", data.ID).Create(value)
		return value, ret.Error
	}
}

func (i *IDeviceSubjectDaoImpl) List(c *gin.Context, req *deviceModels.SysDeviceSubjectReq) {

}

func (i *IDeviceSubjectDaoImpl) Remove(c *gin.Context, ids []string) {

}
