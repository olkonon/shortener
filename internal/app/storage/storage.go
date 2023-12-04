package storage

import (
	"context"
	"errors"
)

// ErrDuplicateURL говорит о том что пытаются добавить уже существующий URL
var ErrDuplicateURL = errors.New("duplicate! URL is exists")
var ErrUserURLListEmpty = errors.New("user no URL")
var ErrDeletedURL = errors.New("url is deleted")

// Storage интерфейс для хранилища данных
type Storage interface {
	//GenIDByURL генерирует ID сокращенной ссылки из полученного URL
	GenIDByURL(ctx context.Context, url string, user string) (string, error)
	//GetURLByID возвращает URL соответствующий ID сокращенной ссылки
	GetURLByID(ctx context.Context, id string) (string, error)
	//GetByUser возвращает все сохраненные URL для пользователя
	GetByUser(ctx context.Context, user string) ([]UserRecord, error)
	//BatchSave сохраняет пачку запросов
	BatchSave(ctx context.Context, data []BatchSaveRequest, user string) ([]BatchSaveResponse, error)
	//BatchDelete асинхронно удаляет пачку url у пользователя
	BatchDelete(ctx context.Context, data []string, user string)
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

type UserRecord struct {
	OriginalURL string
	ShortID     string
}
