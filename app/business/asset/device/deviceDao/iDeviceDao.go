package deviceDao

import (
	"nova-factory-server/app/business/asset/device/deviceModels"

	"github.com/gin-gonic/gin"
)

type IDeviceDao interface {
	InsertDevice(c *gin.Context, device *deviceModels.DeviceInfo) (*deviceModels.DeviceVO, error)
	UpdateDevice(c *gin.Context, device *deviceModels.DeviceInfo) (*deviceModels.DeviceVO, error)
	GetDeviceGroupByName(c *gin.Context, name string) (*deviceModels.DeviceVO, error)
	GetNoExitIdDeviceGroupByName(c *gin.Context, name string, id uint64) (*deviceModels.DeviceVO, error)
	SelectDeviceList(c *gin.Context, req *deviceModels.DeviceListReq) (*deviceModels.DeviceInfoListData, error)
	DeleteByDeviceIds(c *gin.Context, ids []int64) error
	GetLocalByGateWayId(c *gin.Context, id int64) ([]*deviceModels.DeviceVO, error)
	GetByIds(c *gin.Context, ids []int64) ([]*deviceModels.DeviceVO, error)
	GetByIdString(c *gin.Context, id string) (*deviceModels.DeviceVO, error)
	GetById(c *gin.Context, id int64) (*deviceModels.DeviceVO, error)
	GetByTag(c *gin.Context, number string) (*deviceModels.DeviceVO, error)
	// SelectPublicDeviceList 非登录情况下请求的接口
	SelectPublicDeviceList(c *gin.Context, req *deviceModels.DeviceListReq) (*deviceModels.DeviceInfoListData, error)
	// All 读取全部信息
	All(c *gin.Context) ([]*deviceModels.DeviceVO, error)
	// Count 设备数量统计
	Count(c *gin.Context) (int64, error)
}
