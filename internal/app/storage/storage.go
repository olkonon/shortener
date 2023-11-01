package storage

// Storage интерфейс для хранилища данных
type Storage interface {
	//GenIDByURL генерирует ID сокращенной ссылки из полученного URL
	GenIDByURL(url string) (string, error)
	//GetURLByID возвращает URL соответствующий ID сокращенной ссылки
	GetURLByID(id string) (string, error)
	//Close корректно завершает работу любого Storage
	Close() error
}
