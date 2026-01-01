package deviceDaoImpl

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/constant/device"
	"nova-factory-server/app/constant/protocols"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type sysDeviceDataDao struct {
	ms        *gorm.DB
	tableName string
}

func NewSysDeviceDaoImpl(ms *gorm.DB) deviceDao.IDeviceDao {
	return &sysDeviceDataDao{
		ms:        ms,
		tableName: "sys_device",
	}
}

func (s *sysDeviceDataDao) InsertDevice(c *gin.Context, device *deviceModels.DeviceInfo) (*deviceModels.DeviceVO, error) {
	if device == nil {
		return nil, errors.New("device is nil")
	}
	vo := deviceModels.NewDeviceVO(device)
	vo.DeviceId = uint64(snowflake.GenID())
	vo.SetCreateBy(baizeContext.GetUserId(c))
	ret := s.ms.Table(s.tableName).Create(vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return vo, nil
}
func (s *sysDeviceDataDao) UpdateDevice(c *gin.Context, device *deviceModels.DeviceInfo) (*deviceModels.DeviceVO, error) {
	if device == nil {
		return nil, errors.New("device is nil")
	}
	//deptId := baizeContext.GetDeptId(c)
	ret := s.ms.Table(s.tableName).Where("device_id = ?", device.DeviceId).Updates(device)
	if ret.Error != nil {
		return nil, ret.Error
	}
	var vo deviceModels.DeviceVO
	ret = s.ms.Table(s.tableName).Where("device_id = ?", device.DeviceId).Find(&vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &vo, nil
}

func (s *sysDeviceDataDao) GetDeviceGroupByName(c *gin.Context, name string) (*deviceModels.DeviceVO, error) {
	var vo deviceModels.DeviceVO
	ret := s.ms.Table(s.tableName).Where("name = ?", name).First(&vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &vo, nil
}

func (s *sysDeviceDataDao) GetByIds(c *gin.Context, ids []int64) ([]*deviceModels.DeviceVO, error) {
	var list []*deviceModels.DeviceVO
	ret := s.ms.Table(s.tableName).Where("device_id in (?)", ids).Find(&list)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return list, nil
}

func (s *sysDeviceDataDao) GetNoExitIdDeviceGroupByName(c *gin.Context, name string, id uint64) (*deviceModels.DeviceVO, error) {
	var vo deviceModels.DeviceVO
	ret := s.ms.Table(s.tableName).Where("group_id != ?", id).Where("name = ?", name).First(&vo)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &vo, nil
}

func (s *sysDeviceDataDao) SelectDeviceList(c *gin.Context, req *deviceModels.DeviceListReq) (*deviceModels.DeviceInfoListData, error) {
	db := s.ms.Table(s.tableName)

	if req != nil && req.Name != nil && *req.Name != "" {
		db = db.Where("name LIKE ?", "%"+*req.Name+"%")
	}
	if req != nil && req.DeviceGroupId > 0 {
		db = db.Where("device_group_id = ?", req.DeviceGroupId)
	}
	if req != nil && req.Number != nil && *req.Number != "" {
		db = db.Where("number = ?", req.Number)
	}
	if req != nil && req.ControlType != nil {
		db = db.Where("control_type = ?", *req.ControlType)
	}
	if req != nil && req.Type != nil && *req.Type != "" {
		db = db.Where("type = ?", *req.Type)
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
	db = db.Where("state", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(c, db)
	var dto []*deviceModels.DeviceVO

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &deviceModels.DeviceInfoListData{
			Rows:  make([]*deviceModels.DeviceVO, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Limit(size).Find(&dto).Order("create_time desc")
	if ret.Error != nil {
		return &deviceModels.DeviceInfoListData{
			Rows:  make([]*deviceModels.DeviceVO, 0),
			Total: 0,
		}, ret.Error
	}
	return &deviceModels.DeviceInfoListData{
		Rows:  dto,
		Total: total,
	}, nil
}

func (s *sysDeviceDataDao) DeleteByDeviceIds(c *gin.Context, ids []int64) error {
	ret := s.ms.Table(s.tableName).Where("device_id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (s *sysDeviceDataDao) GetLocalByGateWayId(c *gin.Context, id int64) ([]*deviceModels.DeviceVO, error) {
	var dto []*deviceModels.DeviceVO
	ret := s.ms.Table(s.tableName).Where("device_gateway_id in (?)", id).Where("status = ?", true).Where("communication_type = ?", protocols.LOCAL).Where("state = ?", commonStatus.NORMAL).Find(&dto)

	var res []*deviceModels.DeviceVO = make([]*deviceModels.DeviceVO, 0)
	for k, v := range dto {
		if v.Extension == "" {
			continue
		}
		if v.ProtocolType == device.MQTT {
			var ext deviceModels.ExtensionInfo
			err := json.Unmarshal([]byte(v.Extension), &ext)
			if err != nil {
				zap.L().Error("get extension info error", zap.Error(err))
				continue
			}
			dto[k].ExtensionInfo = &ext
			res = append(res, dto[k])
		} else if v.ProtocolType == device.MODBUS_TCP {
			var ext deviceModels.ExtensionInfo
			err := json.Unmarshal([]byte(v.Extension), &ext)
			if err != nil {
				zap.L().Error("get extension info error", zap.Error(err))
				continue
			}
			dto[k].ExtensionInfo = &ext
			res = append(res, dto[k])
		}

	}
	return res, ret.Error
}

func (s *sysDeviceDataDao) GetByIdString(c *gin.Context, id string) (*deviceModels.DeviceVO, error) {
	var info *deviceModels.DeviceVO
	ret := s.ms.Table(s.tableName).Where("device_id = ?", id).Where("state = ?", commonStatus.NORMAL).First(&info)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return info, nil
}

func (s *sysDeviceDataDao) GetById(c *gin.Context, id int64) (*deviceModels.DeviceVO, error) {
	var info *deviceModels.DeviceVO
	ret := s.ms.Table(s.tableName).Where("device_id = ?", id).Where("state = ?", commonStatus.NORMAL).First(&info)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return info, nil
}

func (s *sysDeviceDataDao) GetByTag(c *gin.Context, number string) (*deviceModels.DeviceVO, error) {
	var info *deviceModels.DeviceVO
	ret := s.ms.Table(s.tableName).Where("number = ?", number).Where("state = ?", commonStatus.NORMAL).First(&info)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return info, nil
}

// SelectPublicDeviceList 非登录情况下请求的接口
func (s *sysDeviceDataDao) SelectPublicDeviceList(c *gin.Context, req *deviceModels.DeviceListReq) (*deviceModels.DeviceInfoListData, error) {
	db := s.ms.Table(s.tableName)

	if req != nil && req.Name != nil && *req.Name != "" {
		db = db.Where("name LIKE ?", "%"+*req.Name+"%")
	}
	if req != nil && req.DeviceGroupId > 0 {
		db = db.Where("device_group_id = ?", req.DeviceGroupId)
	}
	if req != nil && req.Number != nil && *req.Number != "" {
		db = db.Where("number = ?", req.Number)
	}
	if req != nil && req.ControlType != nil {
		db = db.Where("control_type = ?", *req.ControlType)
	}
	if req != nil && req.Type != nil && *req.Type != "" {
		db = db.Where("type = ?", *req.Type)
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
	db = db.Where("state", commonStatus.NORMAL)
	var dto []*deviceModels.DeviceVO

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &deviceModels.DeviceInfoListData{
			Rows:  make([]*deviceModels.DeviceVO, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Limit(size).Find(&dto).Order("create_time desc")
	if ret.Error != nil {
		return &deviceModels.DeviceInfoListData{
			Rows:  make([]*deviceModels.DeviceVO, 0),
			Total: 0,
		}, ret.Error
	}
	return &deviceModels.DeviceInfoListData{
		Rows:  dto,
		Total: total,
	}, nil
}
