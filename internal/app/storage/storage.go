package storage

import (
	"context"
	"errors"
)

// ErrDuplicateURL говорит о том что пытаются добавить уже существующий URL
var ErrDuplicateURL = errors.New("duplicate! URL is exists")

// Storage интерфейс для хранилища данных
type Storage interface {
	//GenIDByURL генерирует ID сокращенной ссылки из полученного URL
	GenIDByURL(ctx context.Context, url string) (string, error)
	//GetURLByID возвращает URL соответствующий ID сокращенной ссылки
	GetURLByID(ctx context.Context, id string) (string, error)
	//BatchSave сохраняет пачку запросов
	BatchSave(ctx context.Context, data []BatchSaveRequest) ([]BatchSaveResponse, error)
	//Close корректно завершает работу любого Storage
	Close() error
}

type BatchSaveRequest struct {
	CorrelationID string
	OriginalURL   string
}

type BatchSaveResponse struct {
	CorrelationID string
	ShortID       string
}
