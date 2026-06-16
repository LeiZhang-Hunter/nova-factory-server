package guanjiapo

import "nova-factory-server/app/utils/observer/integration/adapter/guanjiapo/client"

// KindGuanJiaPo 管家婆集成类型标识，保持 gjp_v1 稳定。
const KindGuanJiaPo = client.KindGuanJiaPo

// 以下类型别名保持外部调用方兼容，实际定义在 client 包。
type (
	ConfigSnapshot = client.ConfigSnapshot
	Credentials    = client.Credentials
)
