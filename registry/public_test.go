package registry_test

import (
	"testing"

	"github.com/olzhasar/go-fileserver/registry"
)

func TestGenerateUniqueToken(t *testing.T) {
	t.Run("returns specified length", func(t *testing.T) {
		token := registry.GenerateUniqueToken()

		if len(token) != registry.TOKEN_LENGTH {
			t.Fatalf("Got token %q of length %d, want %d", token, len(token), registry.TOKEN_LENGTH)
		}
	})

	t.Run("returns unique values", func(t *testing.T) {
		tokens := make(map[string]bool)

		for i := 0; i < 5; i++ {
			token := registry.GenerateUniqueToken()
			if tokens[token] {
				t.Fatalf("Repeated token value %q", token)
			}
			tokens[token] = true
		}
	})
}

func TestRecordFile(t *testing.T) {
	t.Run("records fileName to registry", func(t *testing.T) {
		reg := registry.NewInMemoryRegistry()

		fileName := "test.txt"

		token, err := registry.RecordFile(reg, fileName, registry.GenerateUniqueToken)

		if err != nil {
			t.Fatalf("Error returned while trying to record file\n%q", err)
		}

		assertFileSavedUnderToken(t, reg, token, fileName)
	})
	t.Run("invokes generate func in case of duplicates", func(t *testing.T) {
		reg := registry.NewInMemoryRegistry()

		fileName := "test.txt"
		fileNameExisting := "existing.txt"

		count := 1
		token_a := "aaaaaaaaaaaaaaaa"
		token_b := "bbbbbbbbbbbbbbbb"

		gen := func() string {
			if count <= 4 {
				count++
				return token_a
			}
			return token_b
		}

		registry.RecordFile(reg, fileNameExisting, gen)
		token, err := registry.RecordFile(reg, fileName, gen)

		if err != nil {
			t.Fatalf("Error returned while trying to record file\n%q", err)
		}

		if token != token_b {
			t.Fatalf("Got token %q, want %q", token, token_b)
		}

		assertFileSavedUnderToken(t, reg, token, fileName)
		assertFileSavedUnderToken(t, reg, token_a, fileNameExisting)
	})
}

func assertFileSavedUnderToken(t testing.TB, r registry.Registry, token, fileName string) {
	got, ok := r.Get(token)

	if !ok {
		t.Fatalf("Token %q has not been saved to registry", token)
	}

	if got != fileName {
		t.Errorf("Got %q, want %q", got, fileName)
	}
}
