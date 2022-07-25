package utils

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz123456789")

// RandomString generate random string.
func RandomString(n int) string {
	randSrc := rand.NewSource(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[randSrc.Int63()%int64(len(letters))]
	}
	return string(b)
}
