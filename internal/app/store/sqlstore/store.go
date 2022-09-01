package sqlstore

import (
	"database/sql"

	"github.com/ash/http-rest-api/internal/app/store"
	_ "github.com/lib/pq" // это анонимный импорт (чтобы не импортировались методы)
)

type Store struct {
	db             *sql.DB
	userRepository *UserRepository
}

// функция New возвращает сконфигурированный интсатнс структуры Store
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// для того чтобы пользователи из внешнего мира могли обратиться к UserRepository и его методам, например так store.User().Create()
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
