package common

import "math/rand"

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenRandomString генерирует строку случайных символов английского алфавита длинной length
func GenRandomString(length int) string {
	buf := make([]byte, length)
	for i := range buf {
		buf[i] = letters[rand.Intn(len(letters))]
	}
	return string(buf)
}
