package store

// описываются интерфейсы для всех репозиториев
import "github.com/ash/http-rest-api/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}
