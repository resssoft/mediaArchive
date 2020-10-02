package services

import (
	"encoding/json"
	"github.com/resssoft/mediaArchive/models"
	"github.com/resssoft/mediaArchive/repositories"
	"github.com/rs/zerolog/log"
)

type ItemApplication interface {
	AddItem(models.Item) error
	List(string, interface{}) ([]*models.Item, error)
	ImportFromJson([]byte) (models.GoogleBookmarks, error)
}

var coubService = CoubHandler{}

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
	if item.URL != "" {
		switch {
		case coubService.Check(item):
			item = coubService.ProcessItem(item)
		}
	}
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

func (r *itemApp) ImportFromJson(data []byte) (models.GoogleBookmarks, error) {
	googleBookmarks := models.GoogleBookmarks{}
	if err := json.Unmarshal(data, &googleBookmarks); err != nil {
		log.Info().
			Err(err).
			Msg("invalid request body for  partners count")
		return googleBookmarks, err
	} else {
		return googleBookmarks, nil
	}
}
