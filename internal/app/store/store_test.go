// здесь создается (конфигурируется) datanaseURL
package store_test

import (
	"os"
	"testing"
)

var databaseURL string

// эта функция вызвается одни раз перед всеми тестами в конкретном пакете
func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost port=25432 dbname=test user=test password=test sslmode=disable"
	}

	os.Exit(m.Run()) //нужно выйти с правильным кодом
}
