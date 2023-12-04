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
		storeByID: make(map[string]map[string]Record),
	}
}

type Record struct {
	OriginalURL string
	IsDeleted   bool
}

// InMemory простое птокобезопасное хранилище на map реализующее интерфейс Storage
type InMemory struct {
	storeByID map[string]map[string]Record
	lock      sync.RWMutex
}

func (im *InMemory) GenIDByURL(_ context.Context, url string, user string) (string, error) {
	im.lock.Lock()
	defer im.lock.Unlock()

	if _, isExists := im.storeByID[user]; !isExists {
		im.storeByID[user] = make(map[string]Record)
	}

	newID := common.GenHashedString(url)
	userStore := im.storeByID[user]
	if val, IDIsExists := userStore[newID]; IDIsExists {
		if val.OriginalURL == url {
			return newID, storage.ErrDuplicateURL
		}
		return "", errors.New("can't generate new ID")
	}

	im.storeByID[user][newID] = Record{OriginalURL: url,
		IsDeleted: false,
	}

	return newID, nil
}

func (im *InMemory) BatchSave(_ context.Context, data []storage.BatchSaveRequest, user string) ([]storage.BatchSaveResponse, error) {
	im.lock.Lock()
	defer im.lock.Unlock()

	batchUpdate := make(map[string]string)
	result := make([]storage.BatchSaveResponse, len(data))
	if _, isExists := im.storeByID[user]; !isExists {
		im.storeByID[user] = make(map[string]Record)
	}

	for i, val := range data {
		newID := common.GenHashedString(val.OriginalURL)
		if existsURL, IDIsExists := im.storeByID[user][newID]; IDIsExists {
			if existsURL.OriginalURL == val.OriginalURL {
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
		im.storeByID[user][key] = Record{
			OriginalURL: val,
			IsDeleted:   false,
		}
	}

	return result, nil
}

func (im *InMemory) GetURLByID(_ context.Context, ID string) (string, error) {
	im.lock.RLock()
	defer im.lock.RUnlock()
	for _, userStore := range im.storeByID {
		if url, isExists := userStore[ID]; isExists {
			if url.IsDeleted {
				return "", storage.ErrDeletedURL
			}
			return url.OriginalURL, nil
		}
	}

	return "", errors.New("unknown id")
}

func (im *InMemory) GetByUser(_ context.Context, user string) ([]storage.UserRecord, error) {
	im.lock.RLock()
	defer im.lock.RUnlock()

	if urlList, isExists := im.storeByID[user]; isExists {
		if len(urlList) == 0 {
			return nil, storage.ErrUserURLListEmpty
		}
		result := make([]storage.UserRecord, 0)
		for short, original := range urlList {
			if !original.IsDeleted {
				result = append(result, storage.UserRecord{
					OriginalURL: original.OriginalURL,
					ShortID:     short,
				})
			}
		}
		return result, nil
	}
	return nil, storage.ErrUserURLListEmpty
}

func (im *InMemory) BatchDelete(_ context.Context, data []string, user string) {
	go func() {
		//Async
		im.lock.Lock()
		defer im.lock.Unlock()

		for _, shortURL := range data {
			if val, ok := im.storeByID[user]; ok {
				if _, exists := val[shortURL]; exists {
					original := im.storeByID[user][shortURL]
					original.IsDeleted = true
					im.storeByID[user][shortURL] = original
				}
			}
		}
	}()
}

func (im *InMemory) Close() error {
	return nil
}
