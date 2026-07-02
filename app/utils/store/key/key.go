package key

type keys interface {
	GetUserId(key string) int64
}
