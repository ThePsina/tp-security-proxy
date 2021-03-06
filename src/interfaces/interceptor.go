package interfaces

import (
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

func Intercept(w http.ResponseWriter, r *http.Request) {
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
