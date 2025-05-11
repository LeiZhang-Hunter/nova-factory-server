package deviceDaoImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type sysDeviceGroupDataDao struct {
	ms        *gorm.DB
	tableName string
}

func NewSysDeviceGroupDaoImpl(ms *gorm.DB) deviceDao.IDeviceGroupDao {
	return &sysDeviceGroupDataDao{
		ms:        ms,
		tableName: "sys_device_group",
	}
}

func (s *sysDeviceGroupDataDao) InsertDeviceGroup(c *gin.Context, group *deviceModels.DeviceGroup) (*deviceModels.DeviceGroupVO, error) {
	if group == nil {
		return nil, errors.New("device is nil")
	}
	vo := deviceModels.NewDeviceGroupVO(group)
	vo.GroupId = uint64(snowflake.GenID())
	vo.SetCreateBy(baizeContext.GetUserId(c))
	vo.DeptId = uint64(baizeContext.GetDeptId(c))
	ret := s.ms.Table(s.tableName).Create(vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return vo, nil
}
func (s *sysDeviceGroupDataDao) UpdateDeviceGroup(c *gin.Context, group *deviceModels.DeviceGroup) (*deviceModels.DeviceGroupVO, error) {
	if group == nil {
		return nil, errors.New("group is nil")
	}
	//deptId := baizeContext.GetDeptId(c)
	ret := s.ms.Table(s.tableName).Where("group_id = ?", group.GroupId).Updates(group)
	if ret.Error != nil {
		return nil, ret.Error
	}
	var vo deviceModels.DeviceGroupVO
	ret = s.ms.Table(s.tableName).Where("group_id = ?", group.GroupId).Find(&vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &vo, nil
}

func (s *sysDeviceGroupDataDao) GetDeviceGroupByName(c *gin.Context, name string) (*deviceModels.DeviceGroupVO, error) {
	var vo deviceModels.DeviceGroupVO
	ret := s.ms.Table(s.tableName).Where("name = ?", name).First(&vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &vo, nil
}

func (s *sysDeviceGroupDataDao) GetNoExitIdDeviceGroupByName(c *gin.Context, name string, id uint64) (*deviceModels.DeviceGroupVO, error) {
	var vo deviceModels.DeviceGroupVO
	ret := s.ms.Table(s.tableName).Where("group_id != ?", id).Where("name = ?", name).First(&vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &vo, nil
}

func (s *sysDeviceGroupDataDao) SelectDeviceGroupList(c *gin.Context, req *deviceModels.DeviceGroupDQL) (*deviceModels.DeviceGroupListData, error) {
	db := s.ms.Table(s.tableName)
	if req == nil {
		req = &deviceModels.DeviceGroupDQL{}
	}
	if req != nil && req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	size := 0
	if req == nil || req.Size <= 0 {
		size = 20
	} else {
		size = int(req.Size)
	}
	offset := 0
	if req == nil || req.Page <= 0 {
		req.Page = 1
	} else {
		offset = int((req.Page - 1) * req.Size)
	}

	db = baizeContext.GetGormDataScope(c, db)
	db = db.Where("state = ?", 0)

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &deviceModels.DeviceGroupListData{
			Rows:  make([]*deviceModels.DeviceGroupVO, 0),
			Total: 0,
		}, ret.Error
	}
	var dto []*deviceModels.DeviceGroupVO
	ret = db.Offset(offset).Limit(size).Order("create_time desc").Find(&dto)
	if ret.Error != nil {
		return &deviceModels.DeviceGroupListData{
			Rows:  make([]*deviceModels.DeviceGroupVO, 0),
			Total: 0,
		}, ret.Error
	}
	return &deviceModels.DeviceGroupListData{
		Rows:  dto,
		Total: total,
	}, nil
}

func (s *sysDeviceGroupDataDao) GetDeviceGroupByIds(c *gin.Context, ids []uint64) ([]*deviceModels.DeviceGroupVO, error) {
	if ids == nil || len(ids) == 0 {
		return make([]*deviceModels.DeviceGroupVO, 0), nil
	}

	var array []*deviceModels.DeviceGroupVO
	ret := s.ms.Table(s.tableName).Where("group_id in (?)", ids).Find(&array)
	if ret.Error != nil {
		zap.L().Error("GetDeviceGroupByIds error", zap.Error(ret.Error))
		return make([]*deviceModels.DeviceGroupVO, 0), ret.Error
	}
	return array, nil
}

func (s *sysDeviceGroupDataDao) DeleteByGroupIds(c *gin.Context, ids []int64) error {
	ret := s.ms.Table(s.tableName).Where("group_id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}
