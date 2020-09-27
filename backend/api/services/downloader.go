package services

import (
	"github.com/resssoft/mediaArchive/interfaces"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

type Downloader struct {
	LastStatus int
	LastSize   int
	interfaces.IDownLoader
}

func (d *Downloader) Download(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	d.LastStatus = response.StatusCode
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	d.LastSize = len(data)
	if response.StatusCode != 200 {
		log.
			Error().
			Int("statusCode", response.StatusCode).
			Str("URL", url).
			Str("body", string(data)).
			Send()
	}
	return data, err
}
