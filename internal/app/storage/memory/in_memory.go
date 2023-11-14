package memory

import (
	"context"
	"errors"
	"github.com/olkonon/shortener/internal/app/common"
	"github.com/olkonon/shortener/internal/app/storage"
	"sync"
)

func NewInMemory() *InMemory {
	return &InMemory{
		storeByID: make(map[string]string),
	}
}

// InMemory простое птокобезопасное хранилище на map реализующее интерфейс Storage
type InMemory struct {
	storeByID map[string]string
	lock      sync.RWMutex
}

func (im *InMemory) GenIDByURL(_ context.Context, url string) (string, error) {
	im.lock.Lock()
	defer im.lock.Unlock()

	newID := common.GenHashedString(url)
	if val, IDIsExists := im.storeByID[newID]; IDIsExists {
		if val == url {
			return newID, nil
		}
		return "", errors.New("can't generate new ID")
	}

	im.storeByID[newID] = url

	return newID, nil
}

func (im *InMemory) BatchSave(_ context.Context, data []storage.BatchSaveRequest) ([]storage.BatchSaveResponse, error) {
	im.lock.Lock()
	defer im.lock.Unlock()

	batchUpdate := make(map[string]string)
	result := make([]storage.BatchSaveResponse, len(data))
	for i, val := range data {
		newID := common.GenHashedString(val.OriginalURL)
		if existsURL, IDIsExists := im.storeByID[newID]; IDIsExists {
			if existsURL == val.OriginalURL {
				result[i] = storage.BatchSaveResponse{
					CorrelationID: val.CorrelationID,
					ShortID:       newID,
				}
				continue
			}
			return result, errors.New("can't generate new ID")
		}

		batchUpdate[newID] = val.OriginalURL
		result[i] = storage.BatchSaveResponse{
			CorrelationID: val.CorrelationID,
			ShortID:       newID,
		}
	}

	//Это нужно для атомарности, чтобы если возникнет ошибка данные не изменились
	for key, val := range batchUpdate {
		im.storeByID[key] = val
	}

	return result, nil
}

func (im *InMemory) GetURLByID(_ context.Context, ID string) (string, error) {
	im.lock.RLock()
	defer im.lock.RUnlock()

	url, isExists := im.storeByID[ID]
	if !isExists {
		return "", errors.New("unknown id")
	}
	return url, nil
}

func (im *InMemory) Close() error {
	return nil
}
