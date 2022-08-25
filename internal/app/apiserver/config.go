// конфигурируем сервер
package apiserver

import "github.com/ash/http-rest-api/internal/app/store"

// структура описывающая конфиг для сервера
type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Store    *store.Config
}

// создаем конфиг с дефолтными параметрами
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
		Store:    store.NewConfig(),
	}
}
