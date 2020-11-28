package services

import (
	"encoding/json"
	"github.com/resssoft/mediaArchive/models"
	"github.com/resssoft/mediaArchive/repositories"
	"github.com/rs/zerolog/log"
)

const groupDefaultLevel = 100

type ItemApplication interface {
	AddItem(models.Item) error
	AddItemGroup(models.ItemGroup) error
	List(models.DataFilter) ([]*models.Item, error)
	ImportFromJson([]byte) (models.GoogleBookmarks, error)
	GroupList(models.DataFilter) ([]*models.ItemGroup, error)
	AddSimpleGroups([]string, string)
}

var coubService = CoubHandler{}

type itemApp struct {
	repo      repositories.ItemRepository
	groupRepo repositories.ItemGroupRepository
}

func NewItemApp(repo repositories.ItemRepository, groupRepo repositories.ItemGroupRepository) ItemApplication {
	app := &itemApp{
		repo:      repo,
		groupRepo: groupRepo,
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
	// can be added if doesn't exist
	go r.AddSimpleGroups(item.Groups, item.UserID)
	return err
}

func (r *itemApp) List(filter models.DataFilter) ([]*models.Item, error) {
	items, err := r.repo.List(filter)
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

func (r *itemApp) AddItemGroup(itemGroup models.ItemGroup) error {
	var err error
	if itemGroup.ParentCode != "" {
		itemGroup.Code = itemGroup.ParentCode + "#" + itemGroup.Name
	}
	if itemGroup.Code == "" {
		itemGroup.Code = itemGroup.Name
	}
	err = r.groupRepo.Add(itemGroup)
	if err != nil {
		return err
	}
	return err
}

func (r *itemApp) GroupList(filter models.DataFilter) ([]*models.ItemGroup, error) {
	items, err := r.groupRepo.List(filter)
	if err != nil {
		return items, err
	}
	return items, err
}

func (r *itemApp) AddSimpleGroups(codes []string, userID string) {
	for _, code := range codes {
		r.groupRepo.Add(models.ItemGroup{
			Code:   code,
			Name:   code,
			UserID: userID,
			Sort:   groupDefaultLevel,
		})
	}
}
