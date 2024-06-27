package repository

import (
	"pet/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Item interface {
	Create(item model.Item) (uuid.UUID, error)
	FindByTitle(title string, limit int, offset int) (items []model.Item, err error)
	// FindAllByFilter
	Update(id uuid.UUID, item model.UpdateItem) (err error)
	Delete(id uuid.UUID) (err error)
}

type Authorization interface {
	CreateUser(name string, username string, password string, userRole string) (id uuid.UUID, err error)
	GetUser(username string, password string) (user model.User, err error)
	CreateSession(userId uuid.UUID, userRole string, sessionid uuid.UUID, refreshtoken string) (err error)
	FindSession(sessionid uuid.UUID) (session model.Session, err error)
}

type Repository struct {
	Authorization
	Item
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthService(db),
		Item:          NewItemService(db),
	}
}
