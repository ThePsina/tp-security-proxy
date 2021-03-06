package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"proxy/src/interfaces"
)

func main() {
	interceptor := &http.Server{
		Addr: ":8001",
		Handler: http.HandlerFunc(interfaces.Intercept),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	repeater := &http.Server{
		Addr: ":8002",
		Handler: http.HandlerFunc(interfaces.Repeat),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	go func() {
		log.Fatal(repeater.ListenAndServe())
	}()

	fmt.Println("Server start at 8001")
	log.Fatal(interceptor.ListenAndServe())
}
