package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/olkonon/shortener/internal/app/common"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

const ShortURLLen = 10

type Record struct {
	ID  string
	URL string
}

func NewFileStorage(path string) *InFile {
	tmp := &InFile{
		storeByID:  make(map[string]string),
		storeByURL: make(map[string]string),
		filePath:   path,
	}
	if err := tmp.loadCacheFromFile(); err != nil {
		//Данная ошибка фатальна, так как означает что данные повреждены или операция I/O вызывает ошибки!
		log.Fatal(err)
	}
	return tmp
}

// InFile простое птокобезопасное хранилище на map реализующее интерфейс InFile, но хранящее свои данные в файле
type InFile struct {
	storeByID  map[string]string
	storeByURL map[string]string
	filePath   string
	lock       sync.RWMutex
}

func (fs *InFile) GenIDByURL(url string) (string, error) {
	fs.lock.Lock()
	defer fs.lock.Unlock()

	ID, isExists := fs.storeByURL[url]
	if isExists {
		return ID, nil
	}

	newID := common.GenRandomString(ShortURLLen)
	_, IDIsExists := fs.storeByID[newID]
	counter := 1

	for IDIsExists {
		//Исключения случая совпадения, он крайне маловероятен, но все же...
		newID = common.GenRandomString(ShortURLLen)
		_, IDIsExists = fs.storeByID[newID]
		counter++
		if counter >= 128 {
			return "", errors.New("can't generate new ID")
		}
	}

	fs.storeByURL[url] = newID
	fs.storeByID[newID] = url

	err := fs.appendToFile(Record{
		ID:  newID,
		URL: url,
	})
	return newID, err
}

func (fs *InFile) GetURLByID(ID string) (string, error) {
	fs.lock.RLock()
	defer fs.lock.RUnlock()

	url, isExists := fs.storeByID[ID]
	if !isExists {
		return "", errors.New("unknown id")
	}
	return url, nil
}

func (fs *InFile) appendToFile(rec Record) error {
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}

	//Создание файла если нет или добавление в конец если есть
	f, err := os.OpenFile(fs.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	//Авто закрытие файла при выходе из функции
	defer func() {
		if err := f.Close(); err != nil {
			log.Error("Close file error:", err)
		}
	}()

	//Записываем данные
	_, err = f.Write(data)
	if err != nil {
		return err
	}
	//Записываем разделитель
	_, err = f.Write([]byte("\n"))
	if err != nil {
		return err
	}
	//Вызываем sync для записи файла на диск
	err = f.Sync()
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
		fs.storeByID[rec.ID] = rec.URL
		fs.storeByURL[rec.URL] = rec.ID
	}
	return nil
}
