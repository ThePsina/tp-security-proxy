package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"proxy/internal/interfaces"
	"time"
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
	proxy, repeaterRouter := interfaces.CreateProxy()

	interceptor := &http.Server{
		ReadTimeout:  viper.GetDuration("server.timeout.read") * time.Second,
		WriteTimeout: viper.GetDuration("server.timeout.write") * time.Second,
		Addr:         ":" + viper.GetString("server.port.interceptor"),
		Handler:      http.HandlerFunc(proxy.Intercept),
	}

	repeater := &http.Server{
		ReadTimeout:  viper.GetDuration("server.timeout.read") * time.Second,
		WriteTimeout: viper.GetDuration("server.timeout.write") * time.Second,
		Addr:         ":" + viper.GetString("server.port.repeater"),
		Handler:      repeaterRouter,
	}

	go func() {
		log.Fatal(repeater.ListenAndServe())
	}()

	fmt.Printf("Server start at %s", viper.GetString("server.interceptor_port"))
	log.Fatal(interceptor.ListenAndServe())
}
