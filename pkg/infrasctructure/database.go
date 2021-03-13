package infrasctructure

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/spf13/viper"
	"net/http"
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
		if err = tx.Rollback(); err != nil {
		}
	}()

	var id int64
	err = tx.QueryRow("insert into requests (host, request) values ($1, $2) returning id", req.Host, req.Request).
		Scan(&id)

	fmt.Println("\n", id, "\n\n")

	for key, vval := range req.Headers {
		for _, val := range vval {
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
	ctx, cancel := context.WithTimeout(context.Background(), db.timing)
	defer cancel()

	tx, err := db.connection.BeginEx(ctx, nil)
	if err != nil {
		return []entity.Req{}, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
		}
	}()

	out, err := tx.Query("select id, host, request from requests")
	if err != nil {
		return []entity.Req{}, err
	}
	defer out.Close()

	requests := make([]entity.Req, 0, 0)
	for out.Next() {
		var request entity.Req

		if err = out.Scan(&request.Id, &request.Host, &request.Request); err != nil {
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
	ctx, cancel := context.WithTimeout(context.Background(), db.timing)
	defer cancel()

	tx, err := db.connection.BeginEx(ctx, nil)
	if err != nil {
		return entity.Req{}, err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
		}
	}()

	out := tx.QueryRow("select request, host from requests where id=$1", id)

	request := entity.Req{Id: id}
	if err = out.Scan(&request.Request, &request.Host); err != nil {
		return entity.Req{}, err
	}

	if err = tx.Commit(); err != nil {
		return entity.Req{}, err
	}

	return request, nil
}

func (db *Database) GetRequestHeaders(id int64) (entity.Req, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.timing)
	defer cancel()

	tx, err := db.connection.BeginEx(ctx, nil)
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
