package registry_test

import (
	"fmt"
	"testing"

	"github.com/olzhasar/go-fileserver/registry"
)

const TMP_DB_PATH = "../db.sqlite3"

type RegistryTestCase struct {
	name   string
	create func() registry.Registry
}

func NewSQLiteRegistry() registry.Registry {
	registry, _ := registry.NewSQLiteRegistry(TMP_DB_PATH)
	return registry
}

func TestRegistries(t *testing.T) {
	cases := []RegistryTestCase{
		{
			"InMemory",
			registry.NewInMemoryRegistry,
		},
		{
			"SQLite",
			NewSQLiteRegistry,
		},
	}

	for _, test := range cases {
		t.Run(fmt.Sprintf("%s:records filename to registry", test.name), func(t *testing.T) {
			reg := test.create()
			defer reg.Close()
			defer reg.Clear()

			fileName := "test.txt"
			token := "123456"

			err := reg.Record(token, fileName)

			if err != nil {
				t.Fatalf("Expected no error, got %q", err)
			}

			got, ok := reg.Get(token)

			if !ok {
				t.Errorf("Want %q to be in registry, but it's not", token)
			}

			if got != fileName {
				t.Errorf("Got %q, want %q", got, fileName)
			}
		})
		t.Run(fmt.Sprintf("%s:returns false for nonexistent keys", test.name), func(t *testing.T) {
			reg := test.create()
			defer reg.Close()
			defer reg.Clear()

			got, ok := reg.Get("123456")

			if ok {
				t.Error("Got ok true, want false")
			}

			if got != "" {
				t.Errorf("Got %q, want empty string", got)
			}
		})
		t.Run(fmt.Sprintf("%s:Has() returns proper values", test.name), func(t *testing.T) {
			reg := test.create()
			defer reg.Close()
			defer reg.Clear()

			existing_token := "123456"
			reg.Record(existing_token, "file.txt")

			nonexistent_token := "987654"

			if !reg.Has(existing_token) {
				t.Errorf("Want %q to be in registry, but it's not", existing_token)
			}

			if reg.Has(nonexistent_token) {
				t.Errorf("Token %q should not be in registry, but it is", nonexistent_token)
			}
		})
	}
}
