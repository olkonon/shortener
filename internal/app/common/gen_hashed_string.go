package common

import (
	"crypto/sha256"
	"encoding/base64"
)

func GenHashedString(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	//Вероятность совпадения достаточно мала
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))[:10]
}
