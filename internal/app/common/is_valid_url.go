package common

import "net/url"

// IsValidURL Проверка, что data валидный URL
func IsValidURL(data string) bool {
	u, err := url.ParseRequestURI(data)
	return err == nil && u.Scheme != "" && u.Host != ""
}
