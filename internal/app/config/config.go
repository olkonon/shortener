package config

import (
	"flag"
	"github.com/olkonon/shortener/internal/app/common"
	"os"
)

type Config struct {
	BaseURL       string
	ListenAddress string
}

func Parse() Config {
	address := flag.String("a", common.DefaultListenAddress, "Listen server address, default localhost:8080")
	baseURL := flag.String("b", common.DefaultBaseURL, "Short URL base address default http://localhost:8080") // делаем разбор командной строки
	flag.Parse()

	return Config{
		BaseURL:       mergeSetting(*baseURL, "BASE_URL"),
		ListenAddress: mergeSetting(*address, "SERVER_ADDRESS"),
	}
}

func mergeSetting(flagSetting, envSettingName string) string {
	envSetting := os.Getenv(envSettingName)
	if envSetting == "" {
		return flagSetting
	}
	return envSetting
}
