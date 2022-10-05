// конфигурируем сервер
package apiserver

// структура описывающая конфиг для сервера
type Config struct {
	BindAddr      string `toml:"bind_addr"`
	LogLevel      string `toml:"log_level"`
	LogPath       string `toml:"log_path"`
	DatabaseURL   string `toml:"database_url"`
	ConfTokenPath string `toml:"conf_token_path"`
}

// создаем конфиг с дефолтными параметрами
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}
