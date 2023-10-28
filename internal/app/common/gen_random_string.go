package common

import (
	"math/rand"
	"time"
)

var seed *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenRandomString генерирует строку случайных символов английского алфавита длинной length
func GenRandomString(length int) string {
	buf := make([]byte, length)
	for i := range buf {
		buf[i] = letters[seed.Intn(len(letters))]
	}
	return string(buf)
}
