package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/ash/http-rest-api/internal/app/apiserver"
)

var configPath string

// обарабатывает аргумент при вызовае бинарника (связанная переменная, название рагуменат, значение по умолчанию, док коммент)
func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse() //  парсим флаги в связанные переменные

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config) // (toml-файл с конфигом, переменная куда записать данные из файла)
	if err != nil {
		log.Fatal(err)
	}

	s := apiserver.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
