package controllers

import (
	"encoding/json"
	"github.com/resssoft/mediaArchive/models"
	"github.com/resssoft/mediaArchive/repositories"
	"github.com/resssoft/mediaArchive/services"
	"github.com/valyala/fasthttp"
	"net/http"
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
