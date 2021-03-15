package interfaces

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	"github.com/spf13/viper"
	"net/http"
	"proxy/internal/application"
	"proxy/internal/infrasctructure"
)

func CreateProxy() (*Proxy, http.Handler) {
	conf := pgx.ConnConfig{
		User:                 viper.GetString("postgres.user"),
		Database:             viper.GetString("postgres.db"),
		Password:             viper.GetString("postgres.password"),
		Port:                 uint16(viper.GetInt("postgres.port")),
		Host:                 viper.GetString("postgres.host"),
		PreferSimpleProtocol: false,
	}

	db := infrasctructure.CreateDatabaseConnection(conf)
	manager := application.NewDataManager(db)

	proxy := NewProxy(manager)

	repeaterRouter := mux.NewRouter()
	repeaterRouter.HandleFunc("/repeat/{id}", proxy.Repeat).
		Methods(http.MethodGet)
	repeaterRouter.HandleFunc("/requests", proxy.AllRequests).
		Methods(http.MethodGet)
	repeaterRouter.HandleFunc("/request/{id}", proxy.RequestById).
		Methods(http.MethodGet)
	repeaterRouter.HandleFunc("/scan/{id}", proxy.ScanRequest).
		Methods(http.MethodGet)

	return proxy, repeaterRouter
}
