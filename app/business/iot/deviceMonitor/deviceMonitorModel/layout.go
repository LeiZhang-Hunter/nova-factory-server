package deviceMonitorModel

// DeviceLayoutRequest 设备布局请求
type DeviceLayoutRequest struct {
	FloorId int64 `form:"floorId,string" binding:"required"`
}
