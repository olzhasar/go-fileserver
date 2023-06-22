package registry

import (
	"math/rand"
)

const TOKEN_LENGTH = 16
const TOKEN_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateUniqueToken() string {
	letters := make([]rune, TOKEN_LENGTH)

	for i := 0; i < TOKEN_LENGTH; i++ {
		letters[i] = rune(TOKEN_CHARS[rand.Intn(TOKEN_LENGTH)])
	}

	return string(letters)
}
