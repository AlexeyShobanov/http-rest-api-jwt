// описываюстя интерфейсы для всех стров
package store

type Store interface {
	User() UserRepository
}
