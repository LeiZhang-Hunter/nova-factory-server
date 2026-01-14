package homeModels

type HomeStats struct {
	OnlineDevices    int64 `json:"onlineDevices"`
	OfflineDevices   int64 `json:"offlineDevices"`
	BuildingCount    int64 `json:"buildingCount"`
	ExceptionCount   int64 `json:"exceptionCount"`
	MaintenanceCount int64 `json:"maintenanceCount"`
}
