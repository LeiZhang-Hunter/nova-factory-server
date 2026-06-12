package result

// base 同步结果基础接口，所有同步响应的响应体均需实现。
// 提供指针获取、原始字符串序列化及元数据访问能力。
// 定义同步结果的基础接口与通用类型。
// 各类同步响应（OAuthTokenResponse、OrderSyncResponse 等）
// 均需实现 base 接口，以便上层统一获取原始数据与元数据。
type base interface {
	// Ptr 返回当前响应对象的指针
	Ptr() any
	// RawStr 将响应序列化为原始 JSON 字符串
	RawStr() (string, error)
	// MetaData 返回响应的扩展元数据键值对
	MetaData() map[string]any
}
