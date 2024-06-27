package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"       db:"id"`
	Name         string    `json:"name"     db:"name"`
	Username     string    `json:"username" db:"username"`
	PasswordHasn string    `json:"password" db:"password_hash"`
	Role         string    `json:"role"     db:"user_role"`
}

type Token struct {
	SessionID   uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	UserRole    string    `json:"user_role"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiratedAt time.Time `json:"expirated_at"`
}

type Session struct {
	Token        Token
	IsBlocked    bool   `db:"is_blocked"`
	Refreshtoken string `db:"refreshtoken"`
}

type Item struct {
	Id          uuid.UUID `json:"id,omitempty" db:"id"`
	Title       string    `json:"title"       binding:"required" db:"title"`
	Description string    `json:"description" binding:"required" db:"description"`
	Price       int       `json:"price"       binding:"required" db:"price"`
	Image       string    `json:"image"       binding:"required" db:"image"`
}

type UpdateItem struct {
	Title       *string `json:"title" db:"title"`
	Description *string `json:"description" db:"description"`
	Price       *int    `json:"price" db:"price"`
	Image       *string `json:"image" db:"image"`
}
