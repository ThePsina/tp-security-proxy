package interfaces

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"proxy/pkg/domain/entity"
	"time"
)

func (proxy *Proxy) Intercept(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		proxy.tunnel(w, r)
		return
	}

	proxy.proxy(w, r)
}

func (proxy *Proxy) proxy(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
	}

	req := entity.Req{Host: r.Host, Request: string(dump)}
	if err = proxy.dm.Insert(req); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	// resp.Header returns map[string][]string
	// key is header key
	// keyValues is all values in header key
	for key, keyValues := range resp.Header {
		for _, val := range keyValues {
			w.Header().Add(key, val)
		}
	}

	w.WriteHeader(resp.StatusCode)
	if _, err = io.Copy(w, resp.Body); err != nil {
		log.Println(err)
	}
}

func (proxy *Proxy) tunnel(w http.ResponseWriter, r *http.Request) {
	dstConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	cliConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	go transfer(dstConn, cliConn)
	go transfer(cliConn, dstConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer func() {
		if err := destination.Close(); err != nil {
			log.Println(err)
		}
	}()
	defer func() {
		if err := source.Close(); err != nil {
			log.Println(err)
		}
	}()

	if _, err := io.Copy(destination, source); err != nil {
		log.Println(err)
	}
}
