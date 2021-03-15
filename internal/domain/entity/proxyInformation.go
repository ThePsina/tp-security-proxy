package entity

import (
	"crypto/tls"
	"net/http"
)

type ProxyInformation struct {
	InterceptedHttpsRequest *http.Request
	ForwardedHttpsRequest   *http.Request
	Scheme                  string
	Config                  *tls.Config
}
