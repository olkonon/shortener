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
		storeByID: make(map[string]map[string]string),
	}
}

// InMemory простое птокобезопасное хранилище на map реализующее интерфейс Storage
type InMemory struct {
	storeByID map[string]map[string]string
	lock      sync.RWMutex
}

func (im *InMemory) GenIDByURL(_ context.Context, url string, user string) (string, error) {
	im.lock.Lock()
	defer im.lock.Unlock()

	if _, isExists := im.storeByID[user]; !isExists {
		im.storeByID[user] = make(map[string]string)
	}

	newID := common.GenHashedString(url)
	if val, IDIsExists := im.storeByID[user][newID]; IDIsExists {
		if val == url {
			return newID, storage.ErrDuplicateURL
		}
		return "", errors.New("can't generate new ID")
	}

	im.storeByID[user][newID] = url

	return newID, nil
}

func (im *InMemory) BatchSave(_ context.Context, data []storage.BatchSaveRequest, user string) ([]storage.BatchSaveResponse, error) {
	im.lock.Lock()
	defer im.lock.Unlock()

	batchUpdate := make(map[string]string)
	result := make([]storage.BatchSaveResponse, len(data))
	if _, isExists := im.storeByID[user]; !isExists {
		im.storeByID[user] = make(map[string]string)
	}

	for i, val := range data {
		newID := common.GenHashedString(val.OriginalURL)
		if existsURL, IDIsExists := im.storeByID[user][newID]; IDIsExists {
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
		im.storeByID[user][key] = val
	}

	return result, nil
}

func (im *InMemory) GetURLByID(_ context.Context, ID string) (string, error) {
	im.lock.RLock()
	defer im.lock.RUnlock()
	for _, userStore := range im.storeByID {
		if url, isExists := userStore[ID]; isExists {
			return url, nil
		}
	}

	return "", errors.New("unknown id")
}

func (im *InMemory) Close() error {
	return nil
}
