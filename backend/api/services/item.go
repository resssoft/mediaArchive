package services

import (
	"github.com/resssoft/mediaArchive/models"
	"github.com/resssoft/mediaArchive/repositories"
)

type ItemApplication interface {
	AddItem(models.Item) error
	List(string, interface{}) ([]*models.Item, error)
}

type itemApp struct {
	repo repositories.ItemRepository
}

func NewItemApp(repo repositories.ItemRepository) ItemApplication {
	app := &itemApp{
		repo: repo,
	}
	return app
}

func (r *itemApp) AddItem(item models.Item) error {
	var err error
	err = r.repo.Add(item)
	if err != nil {
		return err
	}
	return err
}

func (r *itemApp) List(name string, value interface{}) ([]*models.Item, error) {
	items, err := r.repo.List(name, value)
	if err != nil {
		return items, err
	}
	return items, err
}
