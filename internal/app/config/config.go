package config

import (
	"flag"
	"github.com/olkonon/shortener/internal/app/common"
	"os"
)

type Config struct {
	BaseURL         string
	ListenAddress   string
	StorageFilePath string
	DSN             string
}

func Parse() Config {
	address := flag.String("a", common.DefaultListenAddress, "Listen server address, default "+common.DefaultListenAddress)
	baseURL := flag.String("b", common.DefaultBaseURL, "Short URL base address, default "+common.DefaultBaseURL)
	filePath := flag.String("f", common.DefaultStorageFilePath, "File path for base file storage, default "+common.DefaultStorageFilePath)
	dsn := flag.String("d", common.DefaultDBDSN, "DB connection URL "+common.DefaultDBDSN)
	// делаем разбор командной строки
	flag.Parse()

	return Config{
		BaseURL:         mergeSetting(*baseURL, "BASE_URL"),
		ListenAddress:   mergeSetting(*address, "SERVER_ADDRESS"),
		StorageFilePath: mergeSetting(*filePath, "FILE_STORAGE_PATH"),
		DSN:             mergeSetting(*dsn, "DATABASE_DSN"),
	}
}

func mergeSetting(flagSetting, envSettingName string) string {
	envSetting := os.Getenv(envSettingName)
	if envSetting == "" {
		return flagSetting
	}
	return envSetting
}
