// хелпер для конфигуриции тестового стора и создания функции для сброса тестовых таблиц
package store

import (
	"fmt"
	"strings"
	"testing"
)

// создаем стор, и функцию для очистки базы после теста
func TestStore(t *testing.T, datanaseURL string) (*Store, func(...string)) {
	t.Helper() // этот метод говорит о том что данная функция вспомогательная, ее не нужно тестировать и учитывать

	config := NewConfig()
	config.DatabaseURL = datanaseURL
	s := New(config)
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}

	return s, func(tables ...string) {
		if len(tables) > 0 {
			if _, err := s.db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))); err != nil {
				t.Fatal(err)
			}
		}

		s.Close()
	}
}
