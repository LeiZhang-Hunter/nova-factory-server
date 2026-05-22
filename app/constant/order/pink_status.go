package order

const (
	PinkStatusWaiting  int32 = 0 // 待支付（团长刚创建团）
	PinkStatusProgress int32 = 1 // 进行中（已支付，等待成团）
	PinkStatusSuccess  int32 = 2 // 成团成功
	PinkStatusFailed   int32 = 3 // 拼团失败（超时未成团）
)
