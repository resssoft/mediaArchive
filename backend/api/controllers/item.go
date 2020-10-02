package controllers

import (
	"bytes"
	"encoding/json"
	config "github.com/resssoft/mediaArchive/configuration"
	"github.com/resssoft/mediaArchive/models"
	"github.com/resssoft/mediaArchive/repositories"
	"github.com/resssoft/mediaArchive/services"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
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

func NewUserRoute(repo repositories.ItemRepository, app services.ItemApplication) *ItemRouter {
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
	fh, err := ctx.FormFile("file")
	if err != nil {
		log.Error().AnErr("upload file error", err).Send()
		writeJsonResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	fileName := config.ImportDir + StringWithCharset(10, letters) + ".import"
	if err := fasthttp.SaveMultipartFile(fh, fileName); err != nil {
		log.Error().AnErr("upload file error", err).Send()
		writeJsonResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	file, err := fh.Open()
	defer file.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		log.Error().AnErr("read file error", err).Send()
		writeJsonResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	go func() {
		googleBookmarks, err := r.app.ImportFromJson(buf.Bytes())
		if err != nil {
			log.Error().AnErr("read file error", err).Send()
			writeJsonResponse(ctx, http.StatusBadRequest, err.Error())
			return
		} else {
			flatBookmarks := googleBookmarks.GetFlatList()
			log.Info().Interface("googleBookmarks", flatBookmarks).Send()
			println(len(flatBookmarks))
		}
	}()
	writeJsonResponse(ctx, http.StatusOK, "OK")
}

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
