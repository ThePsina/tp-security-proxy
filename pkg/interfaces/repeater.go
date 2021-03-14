package interfaces

import (
	"bufio"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"proxy/pkg/domain/entity"
	"strconv"
	"strings"
)

func (proxy *Proxy) Repeat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	storedRequest, err := proxy.dm.GetRequestById(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	defer client.CloseIdleConnections()

	newRequest, err := createNewRequest(storedRequest)
	if err != nil {
		proxy.logger.Error(err)
		return
	}

	resp, err := client.Do(newRequest)
	if err != nil {
		proxy.logger.Error(err)
		return
	}
	defer resp.Body.Close()

	copyHeaders(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func createNewRequest(storedRequest entity.Req) (*http.Request, error) {
	requestReader := bufio.NewReader(strings.NewReader(storedRequest.Request))
	buffer, err := http.ReadRequest(requestReader)
	if err != nil {
		return nil, err
	}

	newRequest, err := http.NewRequest(buffer.Method, storedRequest.Host, buffer.Body)
	if err != nil {
		return nil, err
	}

	copyHeaders(buffer.Header, newRequest.Header)
	newRequest.Header.Del("Proxy-Connection")

	return newRequest, nil
}
