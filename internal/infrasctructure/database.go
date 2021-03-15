package infrasctructure

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/spf13/viper"
	"net/http"
	"proxy/internal/domain/entity"
	"time"
)

type Database struct {
	config pgx.ConnConfig
	timing     time.Duration
}

func CreateDatabaseConnection(conf pgx.ConnConfig) *Database {
	return &Database{config: conf, timing: viper.GetDuration("db_connection.timing") * time.Second}
}

func (db *Database) Insert(req entity.Req) error {
	connection, err := pgx.Connect(db.config)
	if err != nil {
		return err
	}
	defer connection.Close()

	ctx, cancel := context.WithTimeout(context.Background(), db.timing)
	defer cancel()

	tx, err := connection.BeginEx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
		}
	}()

	var id int64
	err = tx.QueryRow("insert into requests (host, request) values ($1, $2) returning id", req.URL, req.Request).
		Scan(&id)

	for key, vval := range req.Headers {
		for _, val := range vval {
			if key == "Proxy-Connection" {
				continue
			}
			_, err = tx.Exec("insert into headers (req_id, key, val) values ($1, $2, $3)", id, key, val)
			if err != nil {
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *Database) GetRequestList() ([]entity.Req, error) {
	connection, err := pgx.Connect(db.config)
	if err != nil {
		return []entity.Req{}, err
	}
	defer connection.Close()

	ctx, cancel := context.WithTimeout(context.Background(), db.timing)
	defer cancel()

	tx, err := connection.BeginEx(ctx, nil)
	if err != nil {
		return []entity.Req{}, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
		}
	}()

	out, err := tx.Query("select id, host from requests")
	if err != nil {
		return []entity.Req{}, err
	}
	defer out.Close()

	requests := make([]entity.Req, 0, 0)
	for out.Next() {
		var request entity.Req

		if err = out.Scan(&request.Id, &request.URL); err != nil {
			return []entity.Req{}, err
		}

		requests = append(requests, request)
	}

	if err = out.Err(); err != nil {
		return []entity.Req{}, err
	}
	if err = tx.Commit(); err != nil {
		return []entity.Req{}, err
	}

	return requests, nil
}

func (db *Database) GetRequestById(id int64) (entity.Req, error) {
	connection, err := pgx.Connect(db.config)
	if err != nil {
		return entity.Req{}, err
	}
	defer connection.Close()

	ctx, cancel := context.WithTimeout(context.Background(), db.timing)
	defer cancel()

	tx, err := connection.BeginEx(ctx, nil)
	if err != nil {
		return entity.Req{}, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
		}
	}()

	out := tx.QueryRow("select request, host from requests where id=$1", id)

	request := entity.Req{Id: id}
	if err = out.Scan(&request.Request, &request.URL); err != nil {
		return entity.Req{}, err
	}

	if err = tx.Commit(); err != nil {
		return entity.Req{}, err
	}

	return request, nil
}

func (db *Database) GetRequestHeaders(id int64) (entity.Req, error) {
	connection, err := pgx.Connect(db.config)
	if err != nil {
		return entity.Req{}, err
	}
	defer connection.Close()

	ctx, cancel := context.WithTimeout(context.Background(), db.timing)
	defer cancel()

	tx, err := connection.BeginEx(ctx, nil)
	if err != nil {
		return entity.Req{}, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
		}
	}()

	out, err := tx.Query("select key, val from headers where req_id=$1", id)
	if err != nil {
		return entity.Req{}, err
	}

	request := entity.Req{Id: id, Headers: http.Header{}}
	for out.Next() {
		var key string
		var val string

		if err = out.Scan(&key, &val); err != nil {
			return entity.Req{}, err
		}

		request.Headers.Add(key, val)
	}

	if err = out.Err(); err != nil {
		return entity.Req{}, err
	}
	if err = tx.Commit(); err != nil {
		return entity.Req{}, err
	}

	return request, nil
}
