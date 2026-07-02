package key

type emptyKeys struct{}

func newEmptyPermissions() keys {
	return &emptyKeys{}
}

func (e *emptyKeys) GetUserId(key string) int64 {
	return 0
}
