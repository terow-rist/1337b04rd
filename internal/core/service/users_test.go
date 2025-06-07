package service

import (
	"1337bo4rd/internal/core/domain"
	"testing"
	"time"
)

// --- Mock Implementations ---

type mockUserRepo struct {
	savedUser  *domain.User
	existing   *domain.User
	shouldFail bool
}

func (m *mockUserRepo) GetUserById(id string) (*domain.User, error) {
	if m.shouldFail {
		return nil, errMock
	}
	return m.existing, nil
}

func (m *mockUserRepo) SaveUser(user *domain.User) error {
	m.savedUser = user
	return nil
}

func (m *mockUserRepo) UpdateUser(user *domain.User) error {
	return nil
}

type mockAvatarProvider struct {
	name   string
	avatar string
	err    error
}

func (m *mockAvatarProvider) Next() (string, string, error) {
	return m.name, m.avatar, m.err
}

// --- Dummy error for simulation ---
var errMock = &mockError{"mock error"}

type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}

// --- Test ---

func TestFindOrCreate_NewUser(t *testing.T) {
	repo := &mockUserRepo{}
	avatar := &mockAvatarProvider{name: "Rick", avatar: "url/to/rick.png"}

	service := NewUserService(repo, avatar)

	user, created, err := service.FindOrCreate("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !created {
		t.Errorf("expected user to be created")
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if repo.savedUser == nil {
		t.Fatal("expected user to be saved")
	}
	if repo.savedUser.Name != "Rick" {
		t.Errorf("expected name to be Rick; got %s", repo.savedUser.Name)
	}
	if repo.savedUser.Avatar != "url/to/rick.png" {
		t.Errorf("expected avatar to be url/to/rick.png; got %s", repo.savedUser.Avatar)
	}
}

func TestFindOrCreate_ExistingValidUser(t *testing.T) {
	now := time.Now()
	existing := &domain.User{
		ID:        "user-1",
		Name:      "Morty",
		Avatar:    "url/to/morty.png",
		ExpiresAt: now.Add(1 * time.Hour),
	}

	repo := &mockUserRepo{existing: existing}
	avatar := &mockAvatarProvider{name: "should-not-be-used", avatar: "nope"}

	service := NewUserService(repo, avatar)

	user, created, err := service.FindOrCreate("user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created {
		t.Errorf("expected user to be reused, not created")
	}
	if user.ID != "user-1" {
		t.Errorf("expected ID user-1; got %s", user.ID)
	}
}
