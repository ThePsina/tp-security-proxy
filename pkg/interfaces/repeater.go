package interfaces

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (proxy *Proxy) Repeat(w http.ResponseWriter, r *http.Request) {
	idUrl := r.URL.Query()["id"]
	if len(idUrl) == 0 {
		return
	}

	id, err := strconv.Atoi(idUrl[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	req, err := proxy.dm.GetRequestById(int64(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	reqReader := bufio.NewReader(strings.NewReader(req.Request))
	buffer, err := http.ReadRequest(reqReader)
	if err != nil {
		log.Fatal(err)
	}

	httpReq, err := http.NewRequest(buffer.Method, req.Host, buffer.Body)
	if err != nil {
		log.Fatal(err)
	}

	copyHeaders(buffer.Header, httpReq.Header)

	resp, err := client.Do(httpReq)
	if err != nil {
		log.Fatal(err)
	}

	copyHeaders(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)
	_, _ =io.Copy(w, resp.Body)
}
