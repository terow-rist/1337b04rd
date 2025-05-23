package repository

import (
	"1337bo4rd/internal/core/domain"
	"database/sql"
	"errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUserById(id string) (*domain.User, error) {
	var user domain.User
	query := `
	SELECT * FROM users
	WHERE user_id = $1
	`

	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Avatar, &user.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, nil
}

func (r *UserRepository) SaveUser(user *domain.User) error {
	query := `
	INSERT INTO users (user_id, name, avatar, expires_at)
	VALUES ($1,$2,$3,$4)
	ON CONFLICT(user_id) DO UPDATE SET
		name=excluded.name,
		avatar=excluded.avatar,
		expires_at=excluded.expires_at
	`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Avatar, user.ExpiresAt)
	return err
}
