package port

import "1337bo4rd/internal/core/domain"

type UserRepository interface {
	GetUserById(id string) (*domain.User, error)
	SaveUser(user *domain.User) error
	UpdateUser(user *domain.User) error
}

type UserService interface {
	UpdateUser(user *domain.User) error
}
