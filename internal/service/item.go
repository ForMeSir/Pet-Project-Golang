package service

import (
	"pet/internal/model"
	"pet/internal/repository"

	"github.com/google/uuid"
)

type ItemService struct {
	repo repository.Item
}

func NewItemService(repo repository.Item) *ItemService {
	return &ItemService{repo: repo}
}

func (i *ItemService) CreateItem(item model.Item) (uuid.UUID, error) {
	return i.repo.Create(item)
}

func (i *ItemService) DeleteItem(id uuid.UUID) error {
	return i.repo.Delete(id)
}

func (i *ItemService) UpdateItem() {
}

func (i *ItemService) FindItemByTitle(title string, limit int, offset int) ([]model.Item, error) {
	return i.repo.FindByTitle(title, limit, offset)
}
