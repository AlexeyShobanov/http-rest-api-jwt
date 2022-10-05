package main

import (
	"flag"

	"github.com/BurntSushi/toml"
	"github.com/ash/http-rest-api/internal/app/apiserver"
	"github.com/sirupsen/logrus"
)

var configPath string

// обарабатывает аргумент при вызовае бинарника (связанная переменная, название рагуменат, значение по умолчанию, док коммент)
func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse() //  парсим флаги в связанные переменные

	config := apiserver.NewConfig()               // создаем переменную config (с некоторыми значениями по умолчанию)
	_, err := toml.DecodeFile(configPath, config) // (декодируем toml-файл с конфигом и записываем его в переменную config)
	if err != nil {
		logrus.Fatal(err)
	}

	// стартуем сервер
	if err := apiserver.Start(config); err != nil {
		logrus.Fatal(err)
	}
}
