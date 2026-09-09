package domain

import "time"

type User struct {
	ID           UserID
	Email        string
	Role         Role
	PasswordHash string
	CreatedAt    time.Time
}
