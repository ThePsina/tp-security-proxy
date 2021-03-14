package interfaces

import (
	"github.com/sirupsen/logrus"
	"proxy/internal/domain/repository"
)

type Proxy struct {
	dm     repository.ReqRepo
	logger *logrus.Logger
}

func NewProxy(dm repository.ReqRepo) *Proxy {
	return &Proxy{dm, logrus.New()}
}
