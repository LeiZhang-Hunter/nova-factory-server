package deviceServiceImpl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/asset/device/deviceService"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/business/system/systemModels"
)

type DeviceGroupService struct {
	iDeviceGroupDao deviceDao.IDeviceGroupDao
	iUserDao        systemDao.IUserDao
}

func NewDeviceGroupService(iDeviceGroupDao deviceDao.IDeviceGroupDao, iUserDao systemDao.IUserDao) deviceService.IDeviceGroupService {
	return &DeviceGroupService{
		iDeviceGroupDao: iDeviceGroupDao,
		iUserDao:        iUserDao,
	}
}

func (d *DeviceGroupService) InsertDeviceGroup(c *gin.Context, group *deviceModels.DeviceGroup) (*deviceModels.DeviceGroupVO, error) {

	info, err := d.iDeviceGroupDao.GetDeviceGroupByName(c, *group.Name)
	if info != nil {
		return nil, errors.New("设备分组名已经存在")
	}
	if err != nil {
		zap.L().Error("读取设备分组名字失败", zap.Error(err))
	}

	vo, err := d.iDeviceGroupDao.InsertDeviceGroup(c, group)
	if err != nil {
		zap.L().Error("读取设备分组名字失败", zap.Error(err))
		return nil, errors.New("添加设备分组名字失败")
	}
	return vo, nil
}

func (d *DeviceGroupService) UpdateDeviceGroup(c *gin.Context, group *deviceModels.DeviceGroup) (*deviceModels.DeviceGroupVO, error) {
	info, err := d.iDeviceGroupDao.GetNoExitIdDeviceGroupByName(c, *group.Name, group.GroupId)
	if info != nil {
		return nil, errors.New("设备分组名已经存在")
	}
	if err != nil {
		zap.L().Error("读取设备分组名字失败", zap.Error(err))
	}
	return d.iDeviceGroupDao.UpdateDeviceGroup(c, group)
}

func (d *DeviceGroupService) SelectDeviceGroupList(c *gin.Context, req *deviceModels.DeviceGroupDQL) (*deviceModels.DeviceGroupListData, error) {
	ret, err := d.iDeviceGroupDao.SelectDeviceGroupList(c, req)
	if err != nil {
		zap.L().Error("查询设备分组列表失败", zap.Error(err))
		return ret, err
	}

	//  读取用户id集合
	userIdMap := make(map[int64]bool)
	for _, v := range ret.Rows {
		if v.CreateBy > 0 {
			userIdMap[v.CreateBy] = true
		}

		if v.UpdateBy > 0 {
			userIdMap[v.UpdateBy] = true
		}
	}

	// 格式化服务id
	userIds := make([]int64, 0)
	for k, _ := range userIdMap {
		if k > 0 {
			userIds = append(userIds, k)
		}
	}

	users := d.iUserDao.SelectByUserIds(c, userIds)
	userVoMap := make(map[int64]*systemModels.SysUserDML)
	for _, v := range users {
		userVoMap[v.UserId] = v
	}

	for k, v := range ret.Rows {
		var createUserName string
		var updateUserName string
		userVo, ok := userVoMap[v.CreateBy]
		if ok {
			createUserName = userVo.UserName
		}

		userVo, ok = userVoMap[v.UpdateBy]
		if ok {
			updateUserName = userVo.UserName
		}
		ret.Rows[k].CreateUserName = createUserName
		ret.Rows[k].UpdateUserName = updateUserName
	}
	return ret, nil
}

func (d *DeviceGroupService) DeleteByGroupIds(c *gin.Context, ids []int64) error {
	return d.iDeviceGroupDao.DeleteByGroupIds(c, ids)
}
