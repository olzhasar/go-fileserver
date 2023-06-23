package registry

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

func (r *InMemoryRegistry) Clear() {
	for key := range r.data {
		delete(r.data, key)
	}
}

func (r *InMemoryRegistry) Close() {}

func NewInMemoryRegistry() Registry {
	data := make(map[string]string)
	return &InMemoryRegistry{data}
}
