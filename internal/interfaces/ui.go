package interfaces

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"proxy/internal/domain/entity"
	"strconv"
	"strings"
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

	requests := proxy.scan(request.URL)
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

func (proxy *Proxy) scan(rawUrl string) []string {
	result := make([]string, 0, 0)

	urlToScan, _ := url.Parse(rawUrl)
	for headerKey, headerValues := range urlToScan.Query() {
		for _, headerVal := range headerValues {
			for _, paramVal := range entity.Params {
				if strings.Contains(headerKey, paramVal) && len(headerKey) == len(paramVal) {
					result = append(result, paramVal+" = "+headerVal)
				}
			}
		}
	}

	return result
}
