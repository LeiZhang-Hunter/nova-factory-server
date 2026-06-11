package result

type base interface {
	Ptr() any
	RawStr() (string, error)
	MetaData() map[string]any
}
