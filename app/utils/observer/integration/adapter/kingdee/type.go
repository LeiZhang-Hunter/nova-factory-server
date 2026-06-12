// 金蝶（Kingdee）集成适配器，定义金蝶系统的类型标识常量。
// 当前仅包含 Kind 定义与占位 Observer，具体同步逻辑待实现。
package kingdee

import (
	"nova-factory-server/app/utils/observer/integration/kind"
)

const (
	// KindKingdee 金蝶集成系统的类型标识，注册与路由时使用
	KindKingdee kind.Kind = "Kingdee"
)
