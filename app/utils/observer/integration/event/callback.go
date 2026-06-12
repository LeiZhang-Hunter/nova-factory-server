package event

// Callback 事件处理回调接口，观察者完成同步后调用相应方法通知结果
type Callback interface {
	// OnSuccess 同步处理成功时调用
	OnSuccess()
	// OnError 同步处理失败时调用
	OnError()
}
