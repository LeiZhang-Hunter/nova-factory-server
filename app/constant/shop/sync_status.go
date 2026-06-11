package shop

const (
	OrderSyncStatusPending int32 = 0 // 待同步
	OrderSyncStatusSuccess int32 = 1 // 同步成功
	OrderSyncStatusFailed  int32 = 2 // 同步失败
)
