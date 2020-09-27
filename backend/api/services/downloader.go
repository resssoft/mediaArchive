package services

import (
	"github.com/resssoft/mediaArchive/interfaces"
	"io/ioutil"
	"log"
	"net/http"
)

type Downloader struct {
	status int
	interfaces.IDownLoader
}

func (d *Downloader) Download(url string) ([]byte, error) {
	log.Printf("Download of %s\n", url)
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Status: %v Error: %v \n", err.Error())
		return nil, err
	}
	d.status = response.StatusCode
	log.Println(response.StatusCode)
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Read body error: %v \n", err.Error())
	}
	return data, err
}
