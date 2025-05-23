package service

import (
	"1337bo4rd/internal/core/domain"
	"1337bo4rd/internal/core/port"
	"crypto/rand"
	"encoding/hex"
	"time"
)

const SessionTTL = 7 * 24 * time.Hour

type UserService struct {
	repo   port.UserRepository
	avatar port.AvatarProvider
}

func NewUserService(repo port.UserRepository, avatar port.AvatarProvider) *UserService {
	return &UserService{
		repo:   repo,
		avatar: avatar,
	}
}

func (s *UserService) FindOrCreate(id string) (*domain.User, bool, error) {
	now := time.Now()
	if id != "" {
		if ex, err := s.repo.GetUserById(id); err != nil {
			return nil, false, err
		} else if ex != nil && ex.ExpiresAt.After(now) {
			return ex, false, nil
		}

	}
	name, avatarURL, err := s.avatar.Next()
	if err != nil {
		return nil, false, err
	}
	b := make([]byte, 16)
	rand.Read(b)
	nid := hex.EncodeToString(b)
	user := &domain.User{
		ID:        nid,
		Name:      name,
		Avatar:    avatarURL,
		ExpiresAt: now.Add(SessionTTL),
	}

	if err := s.repo.SaveUser(user); err != nil {
		return nil, false, err
	}
	return user, true, nil
}
