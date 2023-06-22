package registry

type Registry interface {
	Record(token, fileName string) error
	Get(token string) (fileName string, ok bool)
	Has(token string) bool
}

type InMemoryRegistry struct {
	data map[string]string
}

func (r *InMemoryRegistry) Record(token, fileName string) error {
	r.data[token] = fileName
	return nil
}

func (r *InMemoryRegistry) Get(token string) (fileName string, ok bool) {
	val, ok := r.data[token]
	return val, ok
}

func (r *InMemoryRegistry) Has(fileName string) bool {
	_, ok := r.data[fileName]
	return ok
}

func NewInMemoryRegistry() Registry {
	data := make(map[string]string)
	return &InMemoryRegistry{data}
}
