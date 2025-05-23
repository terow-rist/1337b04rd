package port

import "1337bo4rd/internal/core/domain"

type UserRepository interface {
	GetUserById(id string) (*domain.User, error)
	SaveUser(user *domain.User) error
}
