package store

import (
	"database/sql"

	_ "github.com/lib/pq" // это анонимный импорт (чтобы не импортировались методы)
)

type Store struct {
	config         *Config
	db             *sql.DB
	userRepository *UserRepository
}

// функция New возвращает сконфигурированный интсатнс структуры Store
func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

// подключение к БД
func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL) // открытие базы происходит лениво, при следующем первом вызове
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

// отключение от БД
func (s *Store) Close() {
	s.db.Close()
}

// для того чтобы пользователи из внешнего мира могли обратиться к UserRepository и его методам, например так store.User().Create()
func (s *Store) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
