package service

import "pet/internal/repository"

type ItemService struct {
	repo repository.Item
}

func NewItemService(repo repository.Item) *ItemService {
	return &ItemService{repo: repo}
}

func (i *ItemService) CreateItem() {
}

func (i *ItemService) DeleteItem() {
}

func (i *ItemService) UpdateItem() {
}
