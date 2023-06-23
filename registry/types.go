package registry

type Registry interface {
	Record(token, fileName string) error
	Get(token string) (fileName string, ok bool)
	Has(token string) bool
	Clear()
	Close()
}
