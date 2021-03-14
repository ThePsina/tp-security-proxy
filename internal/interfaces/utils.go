package interfaces

import (
	"github.com/spf13/viper"
	"net/http"
	"path/filepath"
)

func copyHeaders(from, to http.Header) {
	for header, values := range from {
		for _, value := range values {
			to.Add(header, value)
		}
	}
}

func ParseConfig(filename string) error {
	fullPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(fullPath)
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}
