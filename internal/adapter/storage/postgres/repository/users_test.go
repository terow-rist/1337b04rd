package repository

import (
	"1337bo4rd/internal/core/domain"
	"testing"
	"time"
)

// fakeUserRepo simulates the UserRepository behavior in-memory.
type fakeUserRepo struct {
	users map[string]domain.User
}

func (f *fakeUserRepo) GetUserById(id string) (*domain.User, error) {
	user, ok := f.users[id]
	if !ok {
		return nil, nil
	}
	return &user, nil
}

func (f *fakeUserRepo) SaveUser(user *domain.User) error {
	if f.users == nil {
		f.users = make(map[string]domain.User)
	}
	f.users[user.ID] = *user
	return nil
}

func TestSaveAndGetUser(t *testing.T) {
	repo := &fakeUserRepo{}

	expected := &domain.User{
		ID:        "123",
		Name:      "Rick Sanchez",
		Avatar:    "portal.png",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	err := repo.SaveUser(expected)
	if err != nil {
		t.Fatalf("SaveUser failed: %v", err)
	}

	result, err := repo.GetUserById("123")
	if err != nil {
		t.Fatalf("GetUserById failed: %v", err)
	}
	if result == nil {
		t.Fatal("expected user, got nil")
	}
	if result.Name != expected.Name || result.Avatar != expected.Avatar {
		t.Errorf("unexpected user data: got %+v, expected %+v", result, expected)
	}
}

func TestGetUserById_NotFound(t *testing.T) {
	repo := &fakeUserRepo{users: map[string]domain.User{}}

	result, err := repo.GetUserById("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil, got %+v", result)
	}
}
