package domain

import "time"

type User struct {
	ID        string
	Name      string
	Avatar    string
	ExpiresAt time.Time
}
