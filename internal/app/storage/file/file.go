package file

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"github.com/olkonon/shortener/internal/app/common"
	"github.com/olkonon/shortener/internal/app/storage"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

type Record struct {
	ID   string
	URL  string
	User string
}

func NewFileStorage(path string) *InFile {
	tmp := &InFile{
		storeByID: make(map[string]map[string]string),
		filePath:  path,
	}
	if err := tmp.loadCacheFromFile(); err != nil {
		//Данная ошибка фатальна, так как означает что данные повреждены или операция I/O вызывает ошибки!
		log.Fatal(err)
	}
	//Создание файла если нет или добавление в конец если есть
	if f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600); err != nil {
		//Данная ошибка фатальна, так как означает что операция I/O вызывает ошибки!
		log.Fatal(err)
	} else {
		tmp.f = f
	}
	return tmp
}

// InFile простое птокобезопасное хранилище на map реализующее интерфейс InFile, но хранящее свои данные в файле
type InFile struct {
	storeByID map[string]map[string]string
	filePath  string
	f         *os.File
	lock      sync.RWMutex
}

func (fs *InFile) GenIDByURL(_ context.Context, url string, user string) (string, error) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	if _, isExists := fs.storeByID[user]; !isExists {
		fs.storeByID[user] = make(map[string]string)
	}

	newID := common.GenHashedString(url)
	if val, IDIsExists := fs.storeByID[user][newID]; IDIsExists {
		if val == url {
			return newID, storage.ErrDuplicateURL
		}
		return "", errors.New("can't generate new ID")
	}

	fs.storeByID[user][newID] = url

	err := fs.appendToFile(Record{
		ID:   newID,
		URL:  url,
		User: user,
	})
	return newID, err
}

func (fs *InFile) BatchSave(_ context.Context, data []storage.BatchSaveRequest, user string) ([]storage.BatchSaveResponse, error) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	if _, isExists := fs.storeByID[user]; !isExists {
		fs.storeByID[user] = make(map[string]string)
	}

	result := make([]storage.BatchSaveResponse, len(data))
	for i, val := range data {
		newID := common.GenHashedString(val.OriginalURL)
		if existsURL, IDIsExists := fs.storeByID[user][newID]; IDIsExists {
			if existsURL == val.OriginalURL {
				result[i] = storage.BatchSaveResponse{
					CorrelationID: val.CorrelationID,
					ShortID:       newID,
				}
				continue
			}
			return result, errors.New("can't generate new ID")
		}

		err := fs.appendToFile(Record{
			ID:   newID,
			URL:  val.OriginalURL,
			User: user,
		})
		if err != nil {
			return nil, err
		}

		fs.storeByID[user][newID] = val.OriginalURL
		result[i] = storage.BatchSaveResponse{
			CorrelationID: val.CorrelationID,
			ShortID:       newID,
		}
	}
	return result, nil
}

func (fs *InFile) GetURLByID(_ context.Context, ID string) (string, error) {
	fs.lock.RLock()
	defer fs.lock.RUnlock()

	for _, userStore := range fs.storeByID {
		if url, isExists := userStore[ID]; isExists {
			return url, nil
		}
	}

	return "", errors.New("unknown id")
}

func (fs *InFile) Close() error {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	if fs.f != nil {
		return fs.f.Close()
	}
	return nil
}

func (fs *InFile) appendToFile(rec Record) error {
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}

	//Записываем данные
	_, err = fs.f.Write(data)
	if err != nil {
		return err
	}
	//Записываем разделитель
	_, err = fs.f.Write([]byte("\n"))
	if err != nil {
		return err
	}
	//Вызываем sync для гарантии не потери данных (замедлит работу, но существенно повысит надежность,
	//так количество операций записи много меньше количества операций чтения, существенного влияния на
	//производительность сервиса оказать не должно)
	err = fs.f.Sync()
	if err != nil {
		return err
	}
	return nil
}

func (fs *InFile) loadCacheFromFile() error {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	f, err := os.OpenFile(fs.filePath, os.O_CREATE|os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Error("Close file error:", err)
		}
	}()
	r := bufio.NewReader(f)
	//Построчное чтение и декодирование файла
	for {
		data, _, err := r.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		//Парсинг json строки
		var rec Record
		if err = json.Unmarshal(data, &rec); err != nil {
			return err
		}
		if _, isExists := fs.storeByID[rec.User]; !isExists {
			fs.storeByID[rec.User] = make(map[string]string)
		}

		fs.storeByID[rec.User][rec.ID] = rec.URL
	}
	return nil
}

func (fs *InFile) GetByUser(_ context.Context, user string) ([]storage.UserRecord, error) {
	fs.lock.RLock()
	defer fs.lock.RUnlock()

	if urlList, isExists := fs.storeByID[user]; isExists {
		if len(urlList) == 0 {
			return nil, storage.ErrUserURLListEmpty
		}
		result := make([]storage.UserRecord, 0)
		for short, original := range urlList {
			result = append(result, storage.UserRecord{
				OriginalURL: original,
				ShortID:     short,
			})
		}
		return result, nil
	}
	return nil, storage.ErrUserURLListEmpty
}
