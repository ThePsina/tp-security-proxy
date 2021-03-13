package interfaces

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"proxy/pkg/domain/entity"
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

	request, err := proxy.dm.GetRequestHeaders(int64(id))
	if err != nil {
		proxy.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	requests := proxy.scan(request.Headers)
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

func (proxy *Proxy) scan(headers http.Header) []string {
	result := make([]string, 0, 0)

	for _, paramVal := range entity.Params {
		for headerKey, headerValues := range headers {
			for _, headerVal := range headerValues {
				if strings.HasPrefix(headerKey, paramVal) {
					result = append(result, headerKey+" = "+headerVal)
				}
			}
		}
	}

	return result
}
