package main

import (
	"crypto/tls"
	"fmt"
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

	db := infrasctructure.CreateDatabaseConnection(conn)
	manager := application.NewDataManager(db)

	proxy := interfaces.NewProxy(manager)
	interceptor := &http.Server{
		Addr:         ":" + viper.GetString("server.interceptor_port"),
		Handler:      http.HandlerFunc(proxy.Intercept),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	repeater := &http.Server{
		Addr:         ":" + viper.GetString("server.repeater_port"),
		Handler:      http.HandlerFunc(proxy.Repeat),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	go func() {
		log.Fatal(repeater.ListenAndServe())
	}()

	fmt.Printf("Server start at %s", viper.GetString("server.interceptor_port"))
	log.Fatal(interceptor.ListenAndServe())
}
