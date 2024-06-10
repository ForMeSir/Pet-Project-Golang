package service

import (
	"pet/internal/model"
	"pet/internal/repository"

	"github.com/google/uuid"
)

type Authorization interface {
	CreateUser(name string, username string, password string) (id uuid.UUID, err error)
	CreateAdmin(name string, username string, password string) (id uuid.UUID, err error)
	GetUser(username string, password string) (user model.User, err error)
	GenerateToken(userId uuid.UUID, userRole string) (refreshtoken string, accesstoken string, err error)
	ParseToken(token string) (ptoken model.Token, err error)
	Refresh(refreshtoken string, accesstoken string) (actoken string, err error)
}

type Item interface {
	CreateItem(item model.Item) (uuid.UUID, error)
	FindItemByTitle(title string)([]model.Item, error)
	DeleteItem()
	UpdateItem()
}

type Service struct {
	Authorization
	Item
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Item:          NewItemService(repos.Item),
	}
}
