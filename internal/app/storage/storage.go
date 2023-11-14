package storage

import "context"

// Storage интерфейс для хранилища данных
type Storage interface {
	//GenIDByURL генерирует ID сокращенной ссылки из полученного URL
	GenIDByURL(ctx context.Context, url string) (string, error)
	//GetURLByID возвращает URL соответствующий ID сокращенной ссылки
	GetURLByID(ctx context.Context, id string) (string, error)
	//Close корректно завершает работу любого Storage
	Close() error
}
