package result

type SyncProductResponse interface {
	GetCode() int
	GetMessage() string
	GetMetadata() map[string]any
}
