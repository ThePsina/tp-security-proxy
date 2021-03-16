package interfaces

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"proxy/internal/domain/entity"
	"strconv"
)

func (proxy *Proxy) AllRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := proxy.dm.GetRequestList()
	if err != nil {
		proxy.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(requests)
	if err != nil {
		proxy.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(res); err != nil {
		proxy.logger.Debug(err)
	}
}

func (proxy *Proxy) RequestById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		proxy.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	request, err := proxy.dm.GetRequestById(int64(id))
	if err != nil {
		proxy.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(request)
	if err != nil {
		proxy.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(res); err != nil {
		proxy.logger.Debug(err)
	}
}

func (proxy *Proxy) ScanRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		proxy.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	request, err := proxy.dm.GetRequestById(int64(id))
	if err != nil {
		proxy.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	requests := proxy.scan(request)
	var res []byte
	if len(requests) == 0 {
		res, err = json.Marshal("no param-miner in request")
		if err != nil {
			proxy.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		res, err = json.Marshal(requests)
		if err != nil {
			proxy.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(res); err != nil {
		proxy.logger.Debug(err)
	}
}

// diff between longest and shortest query-param in request
const queryDiff = 32

func (proxy *Proxy) scan(storedRequest entity.Req) []string {
	result := make([]string, 0, 0)

	originalLength, err := proxy.repeatForScan(storedRequest)
	if err != nil {
		proxy.logger.Error(err)
		return nil
	}

	for i, param := range entity.Params {
		if i % 500 == 0 {
			proxy.logger.Infof("working on: %d\n", i)
		}
		modifiedRequest := storedRequest
		modifiedRequest.URL = modifiedRequest.URL + "?" + param + "=1"

		modifiedLength, err := proxy.repeatForScan(modifiedRequest)
		if err != nil {
			proxy.logger.Error(err)
			return nil
		}

		if modifiedLength-originalLength > queryDiff {
			proxy.logger.Infof("add to response: %s\n", param)
			result = append(result, param)
		}
	}
	proxy.logger.Info("\nwork is done")

	return result
}

func (proxy *Proxy) repeatForScan(storedRequest entity.Req) (int, error) {
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	defer client.CloseIdleConnections()

	newRequest, err := createNewRequest(storedRequest)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(newRequest)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	resp.Header.Del("Date")
	length := resp.Header.Get("Content-Length")

	return strconv.Atoi(length)
}
