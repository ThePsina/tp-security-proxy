package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"proxy/pkg/application"
	"proxy/pkg/infrasctructure"
	"proxy/pkg/interfaces"
)

func init() {
	pflag.StringP("config", "c", "", "path to config file")
	pflag.BoolP("help", "h", false, "usage info")

	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatalln(err)
	}

	if viper.GetBool("help") {
		pflag.Usage()
		os.Exit(0)
	}

	if err := interfaces.ParseConfig(viper.GetString("config")); err != nil {
		log.Fatalln("Error during parse defaults", err)
	}
}

func main() {
	conf := pgx.ConnConfig{
		User:                 viper.GetString("postgres.user"),
		Database:             viper.GetString("postgres.db"),
		Password:             viper.GetString("postgres.password"),
		Port:                 uint16(viper.GetInt("postgres.port")),
		Host:                 viper.GetString("postgres.host"),
		PreferSimpleProtocol: false,
	}

	conn, err := pgx.Connect(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	db := infrasctructure.CreateDatabaseConnection(conn)
	manager := application.NewDataManager(db)

	proxy := interfaces.NewProxy(manager)
	interceptorRouter := mux.NewRouter()
	interceptorRouter.HandleFunc("/", proxy.Intercept)
	interceptorRouter.HandleFunc("/requests", proxy.AllRequests).
		Methods(http.MethodGet)
	interceptorRouter.HandleFunc("/request/{id}", proxy.RequestById).
		Methods(http.MethodGet)
	interceptorRouter.HandleFunc("/scan/{id}", proxy.ScanRequest).
		Methods(http.MethodGet)

	interceptor := &http.Server{
		Addr:         ":" + viper.GetString("server.interceptor_port"),
		Handler:      interceptorRouter,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	repeaterRouter := mux.NewRouter()
	repeaterRouter.HandleFunc("/repeat/{id}", proxy.Repeat).
		Methods(http.MethodGet)

	repeater := &http.Server{
		Addr:         ":" + viper.GetString("server.repeater_port"),
		Handler:      repeaterRouter,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	go func() {
		log.Fatal(repeater.ListenAndServe())
	}()

	fmt.Printf("Server start at %s", viper.GetString("server.interceptor_port"))
	log.Fatal(interceptor.ListenAndServe())
}
