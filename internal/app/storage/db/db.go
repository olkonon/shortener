package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/olkonon/shortener/internal/app/common"
	"github.com/olkonon/shortener/internal/app/storage"
	log "github.com/sirupsen/logrus"
	"time"
)

const CreateTable = `CREATE TABLE IF NOT EXISTS urls (
    	user_id varchar(36) NOT NULL,
    	original_url varchar(256) NOT NULL,
    	short_url varchar(10) NOT NULL,
    	is_deleted boolean NOT NULL,
    	PRIMARY KEY (user_id,original_url)
)`
const SelectURLByID = `SELECT original_url,is_deleted FROM urls WHERE short_url=$1;`
const SelectURLByUser = `SELECT original_url,short_url FROM urls WHERE user_id=$1 AND NOT is_deleted;`
const InsertToTable = `INSERT INTO urls (short_url,original_url,user_id,is_deleted) VALUES ($1,$2,$3,false)`
const DeleteURLByID = `UPDATE urls SET is_deleted=TRUE WHERE user_id=$1 AND short_url = any($2);`

type ChanMsg struct {
	User      string
	ShortURLs []string
}

func NewDatabaseStore(dsn string) *DatabaseStore {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		//Фатальная ошибка с базой что-то явно не так
		log.Fatal("DB connect error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		//Фатальная ошибка с базой что-то явно не так если она за 5с не ответила на пинг
		log.Fatal("DB Ping error", err)
	}

	_, err = db.Exec(CreateTable)
	if err != nil {
		//Фатальная ошибка с базой что-то явно не так
		log.Fatal("DB init tables error", err)
	}

	tmp := &DatabaseStore{
		db:               db,
		deletedChan:      make(chan ChanMsg, 32),
		stopChan:         make(chan bool),
		stopFinishedChan: make(chan bool),
	}

	go tmp.deleteWorker()
	return tmp
}

type DatabaseStore struct {
	db               *sql.DB
	deletedChan      chan ChanMsg
	stopChan         chan bool
	stopFinishedChan chan bool
}

func (dbs *DatabaseStore) GenIDByURL(ctx context.Context, url string, user string) (string, error) {
	newID := common.GenHashedString(url)
	_, err := dbs.db.ExecContext(ctx, InsertToTable, newID, url, user)
	var pgError *pq.Error
	if err == nil {
		return newID, nil
	}
	if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
		return newID, storage.ErrDuplicateURL
	}
	return "", err
}

func (dbs *DatabaseStore) BatchSave(ctx context.Context, data []storage.BatchSaveRequest, user string) ([]storage.BatchSaveResponse, error) {
	result := make([]storage.BatchSaveResponse, len(data))
	tx, err := dbs.db.Begin()
	if err != nil {
		log.Error(err)
		return result, err
	}
	//Откат транзакции если Commit не прошел
	defer tx.Rollback()

	insertStmt, err := tx.PrepareContext(ctx, InsertToTable)
	if err != nil {
		log.Error(err)
		return result, err
	}

	txStmt := tx.StmtContext(ctx, insertStmt)

	for i, val := range data {
		newID := common.GenHashedString(val.OriginalURL)
		if _, err = txStmt.ExecContext(ctx, newID, val.OriginalURL, user); err != nil {
			return result, err
		}

		result[i].CorrelationID = val.CorrelationID
		result[i].ShortID = newID
	}

	return result, tx.Commit()
}

func (dbs *DatabaseStore) GetURLByID(ctx context.Context, ID string) (string, error) {
	rowURL := dbs.db.QueryRowContext(ctx, SelectURLByID, ID)
	var url string
	var isDeleted bool
	err := rowURL.Scan(&url, &isDeleted)
	if err != nil {
		return "", err
	}
	if isDeleted {
		return "", storage.ErrDeletedURL
	}
	return url, nil
}

func (dbs *DatabaseStore) GetByUser(ctx context.Context, user string) ([]storage.UserRecord, error) {
	result := make([]storage.UserRecord, 0)

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := dbs.db.QueryContext(dbCtx, SelectURLByUser, user)

	if err != nil {
		return result, err
	}
	if rows.Err() != nil {
		return result, rows.Err()
	}

	for rows.Next() {
		record := storage.UserRecord{}
		err = rows.Scan(&record.OriginalURL, &record.ShortID)
		if err != nil {
			return result, err
		}

		result = append(result, record)
	}

	if len(result) == 0 {
		return result, storage.ErrUserURLListEmpty
	}

	return result, nil
}

func (dbs *DatabaseStore) Close() error {
	dbs.stopChan <- true
	//Ждем пока воркер закончит работу
	<-dbs.stopFinishedChan

	if dbs.db != nil {
		return dbs.db.Close()
	}

	return nil
}

func (dbs *DatabaseStore) BatchDelete(_ context.Context, data []string, user string) {
	go func() {
		// async
		dbs.deletedChan <- ChanMsg{
			User:      user,
			ShortURLs: data,
		}
	}()
}

func (dbs *DatabaseStore) deleteWorker() {
	defer func() {
		dbs.stopFinishedChan <- true
	}()
	for {
		select {
		case data := <-dbs.deletedChan:
			dbs.deleteRecords(data)
		case <-dbs.stopChan:
			return
		}
	}
}

func (dbs *DatabaseStore) deleteRecords(data ChanMsg) {
	//Async
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := dbs.db.QueryContext(timeoutCtx, DeleteURLByID, data.User, pq.Array(data.ShortURLs))
	if err != nil {
		log.Error(err)
		return
	}
	if rows.Err() != nil {
		log.Error(rows.Err())
		return
	}
}
