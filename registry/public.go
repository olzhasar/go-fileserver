package registry

import (
	"math/rand"
)

const TOKEN_LENGTH = 16
const TOKEN_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateUniqueToken() string {
	letters := make([]rune, TOKEN_LENGTH)

	for i := 0; i < TOKEN_LENGTH; i++ {
		letters[i] = rune(TOKEN_CHARS[rand.Intn(TOKEN_LENGTH)])
	}

	return string(letters)
}

func RecordFile(r Registry, fileName string, generateToken func() string) (token string, err error) {
	token = generateToken()

	for r.Has(token) {
		token = generateToken()
	}
	// todo - use mutex to make this thread-safe
	err = r.Record(token, fileName)
	if err != nil {
		return "", err
	}

	return token, err
}
