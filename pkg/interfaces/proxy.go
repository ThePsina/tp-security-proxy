package interfaces

import (
	"proxy/pkg/domain/repository"
)

type Proxy struct {
	dm repository.ReqRepo
}

func NewProxy(dm repository.ReqRepo) *Proxy {
	return &Proxy{dm}
}
