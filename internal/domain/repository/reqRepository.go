package repository

import "proxy/internal/domain/entity"

type ReqRepo interface {
	Insert(entity.Req) error
	GetRequestList() ([]entity.Req, error)
	GetRequestById(int64) (entity.Req, error)
	GetRequestHeaders(int64) (entity.Req, error)
}
