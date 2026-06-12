package kind

// Kind 集成系统类型标识，每个第三方系统使用唯一的字符串标识自己的类型
// 定义集成系统的类型标识。
// Kind 作为各第三方系统的枚举标识（如 "gjp_v1"、"Kingdee"），
// 用于注册表索引、日志区分及缓存键前缀生成。
type Kind string
