package store

// структура описывающая конфигурацию БД
type Config struct {
	DatabaseURL string `toml:"database_url"` //в конце добавлен тег для ключа в toml-файле
}

// функция возвращающая указатель на конфиг структуру для БД
func NewConfig() *Config {
	return &Config{}
}
