package entity

import "net/http"

type Req struct {
	Id      int64
	Headers http.Header
	URL     string
	Request string
}
