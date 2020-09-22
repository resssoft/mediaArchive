package app

import (
	"github.com/resssoft/mediaArchive/model"
	"github.com/resssoft/mediaArchive/repository"
)

type ItemApplication interface {
	AddItem(model.Item) error
	List(string, interface{}) ([]*model.Item, error)
}

type itemApp struct {
	repo repository.ItemRepository
}

func NewItemApp(repo repository.ItemRepository) ItemApplication {
	app := &itemApp{
		repo: repo,
	}
	return app
}

func (r *itemApp) AddItem(item model.Item) error {
	var err error
	err = r.repo.Add(item)
	if err != nil {
		return err
	}
	return err
}

func (r *itemApp) List(name string, value interface{}) ([]*model.Item, error) {
	items, err := r.repo.List(name, value)
	if err != nil {
		return items, err
	}
	return items, err
}
