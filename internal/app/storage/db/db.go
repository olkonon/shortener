package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/olkonon/shortener/internal/app/common"
	log "github.com/sirupsen/logrus"
	"time"
)

const CreateTable = `CREATE TABLE IF NOT EXISTS urls (
		short_url varchar(10) NOT NULL,
    	original_url varchar(256) NOT NULL,
    	PRIMARY KEY (short_url,original_url)
)`
const SelectURLByID = `SELECT original_url FROM urls WHERE short_url=$1;`
const InsertToTable = `INSERT INTO urls (short_url,original_url) VALUES ($1,$2) ON CONFLICT DO NOTHING`

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

	return &DatabaseStore{
		db: db,
	}
}

type DatabaseStore struct {
	db *sql.DB
}

func (dbs *DatabaseStore) GenIDByURL(ctx context.Context, url string) (string, error) {
	newID := common.GenHashedString(url)
	_, err := dbs.db.ExecContext(ctx, InsertToTable, newID, url)
	var pgError *pq.Error
	if err == nil {
		return newID, nil
	}
	if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
		return newID, nil
	}
	return "", err
}

func (dbs *DatabaseStore) GetURLByID(ctx context.Context, ID string) (string, error) {
	rowURL := dbs.db.QueryRowContext(ctx, SelectURLByID, ID)
	var url string
	err := rowURL.Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (dbs *DatabaseStore) Close() error {
	return dbs.db.Close()
}
