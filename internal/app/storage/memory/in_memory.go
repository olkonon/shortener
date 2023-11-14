package memory

import (
	"context"
	"errors"
	"github.com/olkonon/shortener/internal/app/common"
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
