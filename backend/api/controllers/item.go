package controllers

import (
	"encoding/json"
	config "github.com/resssoft/mediaArchive/configuration"
	"github.com/resssoft/mediaArchive/models"
	"github.com/resssoft/mediaArchive/repositories"
	"github.com/resssoft/mediaArchive/services"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

var (
	letters    = "qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890-="
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type ItemRouter struct {
	repo repositories.ItemRepository
	app  services.ItemApplication
}

func NewItemRoute(repo repositories.ItemRepository, app services.ItemApplication) *ItemRouter {
	return &ItemRouter{repo: repo, app: app}
}

func (r *ItemRouter) ItemsList(ctx *fasthttp.RequestCtx) {
	items, err := r.app.List("", "")
	if err != nil {
		writeJsonResponse(ctx, http.StatusBadRequest, getError(err.Error(), 32))
		return
	}
	writeJsonResponse(ctx, http.StatusOK, models.Response{
		Data:  items,
		Count: len(items),
	})
}

func (r *ItemRouter) AddItem(ctx *fasthttp.RequestCtx) {
	newItem := new(models.Item)
	err := json.Unmarshal(ctx.PostBody(), newItem)
	if err != nil {
		writeJsonResponse(ctx, http.StatusBadRequest, getError(err.Error(), 31))
		return
	}
	err = r.app.AddItem(*newItem)
	if err != nil {
		writeJsonResponse(ctx, http.StatusBadRequest, getError(err.Error(), 32))
		return
	}
	writeJsonResponse(ctx, http.StatusOK, "OK")
}

func (r *ItemRouter) UpdateItem(ctx *fasthttp.RequestCtx) {
}

func (r *ItemRouter) GetItem(ctx *fasthttp.RequestCtx) {
}

func (r *ItemRouter) DeleteItem(ctx *fasthttp.RequestCtx) {
}

func (r *ItemRouter) ExportItems(ctx *fasthttp.RequestCtx) {
}

func (r *ItemRouter) ImportItems(ctx *fasthttp.RequestCtx) {
	params := new(models.ImportParams)
	err := json.Unmarshal(ctx.PostBody(), params)
	if err != nil {
		writeJsonResponse(ctx, http.StatusBadRequest, getError(err.Error(), 01))
		return
	}
	if params.FileName == "" {
		writeJsonResponse(ctx, http.StatusBadRequest, getError("File name is empty", 02))
		return
	}
	go func() {
		//TODO: move logic to task
		fileData, err := ioutil.ReadFile(config.ImportDir + params.FileName)
		googleBookmarks, err := r.app.ImportFromJson(fileData)
		if err != nil {
			log.Error().AnErr("import file error", err).Send()
			return
		} else {
			flatBookmarks := googleBookmarks.GetFlatList()
			for _, flatGB := range flatBookmarks {
				r.app.AddItem(models.Item{
					Title:       flatGB.Name,
					Assignment:  models.ItemForWebPage,
					URL:         flatGB.Url,
					Tags:        nil,
					Categories:  nil,
					Groups:      []string{"gb-" + flatGB.Folder},
					Service:     "googleBookmarks",
					ServiceData: flatGB,
				})
			}
		}
	}()
	writeJsonResponse(ctx, http.StatusOK, "OK")
}

func (r *ItemRouter) UploadFile(ctx *fasthttp.RequestCtx) {
	fh, err := ctx.FormFile("file")
	if err != nil {
		log.Error().AnErr("upload file error", err).Send()
		writeJsonResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	//TODO: add file name prefix with userId
	//TODO: save file info to user
	fileName := StringWithCharset(10, letters) + time.Now().Format(config.DateTimeFlatFormat) + ".import"
	if err := fasthttp.SaveMultipartFile(fh, config.ImportDir+fileName); err != nil {
		log.Error().AnErr("upload file error", err).Send()
		writeJsonResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	writeJsonResponse(ctx, http.StatusOK, models.Response{
		Data: fileName,
	})
	return
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
