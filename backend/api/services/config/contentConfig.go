package config

import (
	"encoding/json"
	"github.com/resssoft/mediaArchive/models"
	"github.com/resssoft/mediaArchive/repositories"
)

const groupDefaultLevel = 100

type ConfigApplication interface {
	AddItem(models.Config) error
	Validate([]byte) (models.Config, error)
	List(models.DataFilter) ([]models.Config, error)
	//Count([]byte) (int, error)
	//Validate([]byte) (models.Config, error)
	//Add(models.Config) (models.Config, error)
	//GetById(int) (models.Config, error)
	//Get(string, string) (models.Config, error)
	//Delete(int) error
	Update(models.Config) error
}

type configApp struct {
	repo repositories.ConfigRepository
}

func NewItemApp(repo repositories.ConfigRepository) ConfigApplication {
	app := &configApp{
		repo: repo,
	}
	return app
}

func (r *configApp) Validate(data []byte) (models.Config, error) {
	configData := models.Config{}
	err := json.Unmarshal(data, &configData)
	return configData, err
}

func (r *configApp) AddItem(item models.Config) error {
	return r.repo.Add(item)
}

func (r *configApp) List(filter models.DataFilter) ([]models.Config, error) {
	items, err := r.repo.List(filter)
	if err != nil {
		return items, err
	}
	return items, err
}

func (r *configApp) Update(configData models.Config) error {
	return r.repo.Update(configData.ID.String(), configData)
}
