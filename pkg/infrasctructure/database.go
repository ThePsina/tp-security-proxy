package infrasctructure

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/spf13/viper"
	"proxy/pkg/domain/entity"
	"time"
)

type Database struct {
	connection *pgx.Conn
	timing     time.Duration
}

func CreateDatabaseConnection(conn *pgx.Conn) *Database {
	return &Database{connection: conn, timing: viper.GetDuration("db_connection.timing") * time.Second}
}

func (db *Database) Insert(req entity.Req) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.timing)
	defer cancel()

	tx, err := db.connection.BeginEx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {}
	}()

	_, err = tx.Exec("insert into requests (host, request) values ($1, $2)", req.Host, req.Request)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *Database) GetRequestList() ([]entity.Req, error) {
	return nil, nil
}

func (db *Database) GetRequestById(id int64) (entity.Req, error) {
	return entity.Req{}, nil
}
