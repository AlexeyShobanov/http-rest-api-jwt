// хелпер для конфигуриции тестового стора и создания функции для сброса тестовых таблиц
package sqlstore

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

// создаем стор, и функцию для очистки базы после теста
func TestDB(t *testing.T, datanaseURL string) (*sql.DB, func(...string)) {
	t.Helper() // этот метод говорит о том что данная функция вспомогательная, ее не нужно тестировать и учитывать

	db, err := sql.Open("postgres", datanaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
		}

		db.Close()
	}
}
