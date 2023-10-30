package memory

import (
	"errors"
	"github.com/olkonon/shortener/internal/app/common"
	"sync"
)

const ShortURLLen = 10

func NewInMemory() *InMemory {
	return &InMemory{
		storeByID:  make(map[string]string),
		storeByURL: make(map[string]string),
	}
}

// InMemory простое птокобезопасное хранилище на map реализующее интерфейс Storage
type InMemory struct {
	storeByID  map[string]string
	storeByURL map[string]string
	lock       sync.RWMutex
}

func (im *InMemory) GenIDByURL(url string) (string, error) {
	im.lock.Lock()
	defer im.lock.Unlock()

	ID, isExists := im.storeByURL[url]
	if isExists {
		return ID, nil
	}

	newID := common.GenRandomString(ShortURLLen)
	_, IDIsExists := im.storeByID[newID]
	counter := 1

	for IDIsExists {
		//Исключения случая совпадения, он крайне маловероятен, но все же...
		newID = common.GenRandomString(ShortURLLen)
		_, IDIsExists = im.storeByID[newID]
		counter++
		if counter >= 128 {
			return "", errors.New("can't generate new ID")
		}
	}

	im.storeByURL[url] = newID
	im.storeByID[newID] = url

	return newID, nil
}

func (im *InMemory) GetURLByID(ID string) (string, error) {
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
