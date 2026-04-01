package devicedaoImpl

import (
	"errors"
	"nova-factory-server/app/business/iot/asset/device/devicedao"
	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type sysDeviceGroupDataDao struct {
	ms        *gorm.DB
	tableName string
}

func NewSysDeviceGroupDaoImpl(ms *gorm.DB) devicedao.IDeviceGroupDao {
	return &sysDeviceGroupDataDao{
		ms:        ms,
		tableName: "sys_device_group",
	}
}

func (s *sysDeviceGroupDataDao) InsertDeviceGroup(c *gin.Context, group *devicemodels.DeviceGroup) (*devicemodels.DeviceGroupVO, error) {
	if group == nil {
		return nil, errors.New("device is nil")
	}
	vo := devicemodels.NewDeviceGroupVO(group)
	vo.GroupId = uint64(snowflake.GenID())
	vo.SetCreateBy(baizeContext.GetUserId(c))
	vo.DeptId = uint64(baizeContext.GetDeptId(c))
	ret := s.ms.Table(s.tableName).Create(vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return vo, nil
}
func (s *sysDeviceGroupDataDao) UpdateDeviceGroup(c *gin.Context, group *devicemodels.DeviceGroup) (*devicemodels.DeviceGroupVO, error) {
	if group == nil {
		return nil, errors.New("group is nil")
	}
	//deptId := baizeContext.GetDeptId(c)
	ret := s.ms.Table(s.tableName).Where("group_id = ?", group.GroupId).Updates(group)
	if ret.Error != nil {
		return nil, ret.Error
	}
	var vo devicemodels.DeviceGroupVO
	ret = s.ms.Table(s.tableName).Where("group_id = ?", group.GroupId).Find(&vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &vo, nil
}

func (s *sysDeviceGroupDataDao) GetDeviceGroupByName(c *gin.Context, name string) (*devicemodels.DeviceGroupVO, error) {
	var vo devicemodels.DeviceGroupVO
	ret := s.ms.Table(s.tableName).Where("name = ?", name).First(&vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &vo, nil
}

func (s *sysDeviceGroupDataDao) GetNoExitIdDeviceGroupByName(c *gin.Context, name string, id uint64) (*devicemodels.DeviceGroupVO, error) {
	var vo devicemodels.DeviceGroupVO
	ret := s.ms.Table(s.tableName).Where("group_id != ?", id).Where("name = ?", name).First(&vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &vo, nil
}

func (s *sysDeviceGroupDataDao) SelectDeviceGroupList(c *gin.Context, req *devicemodels.DeviceGroupDQL) (*devicemodels.DeviceGroupListData, error) {
	db := s.ms.Table(s.tableName)
	if req == nil {
		req = &devicemodels.DeviceGroupDQL{}
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
		return &devicemodels.DeviceGroupListData{
			Rows:  make([]*devicemodels.DeviceGroupVO, 0),
			Total: 0,
		}, ret.Error
	}
	var dto []*devicemodels.DeviceGroupVO
	ret = db.Offset(offset).Limit(size).Order("create_time desc").Find(&dto)
	if ret.Error != nil {
		return &devicemodels.DeviceGroupListData{
			Rows:  make([]*devicemodels.DeviceGroupVO, 0),
			Total: 0,
		}, ret.Error
	}
	return &devicemodels.DeviceGroupListData{
		Rows:  dto,
		Total: total,
	}, nil
}

func (s *sysDeviceGroupDataDao) GetDeviceGroupByIds(c *gin.Context, ids []uint64) ([]*devicemodels.DeviceGroupVO, error) {
	if ids == nil || len(ids) == 0 {
		return make([]*devicemodels.DeviceGroupVO, 0), nil
	}

	var array []*devicemodels.DeviceGroupVO
	ret := s.ms.Table(s.tableName).Where("group_id in (?)", ids).Find(&array)
	if ret.Error != nil {
		zap.L().Error("GetDeviceGroupByIds error", zap.Error(ret.Error))
		return make([]*devicemodels.DeviceGroupVO, 0), ret.Error
	}
	return array, nil
}

func (s *sysDeviceGroupDataDao) DeleteByGroupIds(c *gin.Context, ids []int64) error {
	ret := s.ms.Table(s.tableName).Where("group_id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}
