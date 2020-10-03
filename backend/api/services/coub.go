package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/resssoft/mediaArchive/interfaces"
	"github.com/resssoft/mediaArchive/models"
	"github.com/rs/zerolog/log"
	"regexp"
)

const (
	apiUrlTemplate = "https://coub.com/api/v2/coubs/%s.json"
	serviceName    = "coub.com"
)

var (
	permalinkRE = regexp.MustCompile(`.*?\/([^\/]*$)`)
	serviceRE   = regexp.MustCompile(`https?:\/\/coub\.com`)
)

type CoubHandler struct {
	interfaces.ItemHandler
}

func (handler *CoubHandler) Check(item models.Item) bool {
	if item.Service == serviceName {
		return true
	}
	log.Info().Msg("check url coube " + item.URL)
	log.Info().Interface("serviceRE", serviceRE.MatchString(item.URL)).Send()
	return serviceRE.MatchString(item.URL)
}

func (handler *CoubHandler) ProcessItem(item models.Item) models.Item {
	item.Service = serviceName
	apiResult, err := handler.itemRawApiInfo(item.URL)
	if len(apiResult) == 0 {
		item.Error = "cant coub processed: " + err.Error()
		log.Error().AnErr("download error", err).Send()
		return item
	}
	item.Error = ""
	item.Cache = string(apiResult)
	coubItem := models.CoubItem{}
	if err := json.Unmarshal(apiResult, &coubItem); err != nil {
		log.Info().
			Err(err).
			Msg("invalid request body for  partners count")
	} else {
		item.Title = coubItem.Title
		item.Categories = coubItem.GetCategoriesNames()
		item.Tags = coubItem.GetTagNames()
		item.Image = coubItem.Picture
		if item.Assignment == "" || item.Assignment == models.ItemForWebPage {
			item.Assignment = models.ItemForVideo
		}
		if item.Assignment == models.ItemForMusic {
			item.Media = []models.ItemMedia{
				{
					Type:      "music",
					RemoteUrl: coubItem.AudioFileUrl,
				},
			}
		}
		item.ServiceData = coubItem
	}
	return item
}

func containsKey(s []string, e int) bool {
	for key, _ := range s {
		if key == e {
			return true
		}
	}
	return false
}

func (handler *CoubHandler) itemRawApiInfo(url string) ([]byte, error) {
	permalinkMatch := permalinkRE.FindStringSubmatch(url)
	if containsKey(permalinkMatch, 1) {
		downloader := Downloader{}
		result, err := downloader.Download(fmt.Sprintf(apiUrlTemplate, permalinkMatch[1]))
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New("cant find permalink")
}
