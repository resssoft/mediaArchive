package controller

import (
	"encoding/json"
	"github.com/resssoft/mediaArchive/app"
	"github.com/resssoft/mediaArchive/model"
	"github.com/resssoft/mediaArchive/repository"
	"github.com/valyala/fasthttp"
	"net/http"
)

type ItemRouter struct {
	repo repository.ItemRepository
	app  app.ItemApplication
}

func NewUserRoute(repo repository.ItemRepository, app app.ItemApplication) *ItemRouter {
	return &ItemRouter{repo: repo, app: app}
}

func (r *ItemRouter) ItemsList(ctx *fasthttp.RequestCtx) {
	items, err := r.app.List("", "")
	if err != nil {
		writeJsonResponse(ctx, http.StatusBadRequest, getError(err.Error(), 32))
		return
	}
	writeJsonResponse(ctx, http.StatusOK, model.Response{
		Data:  items,
		Count: len(items),
	})
}

func (r *ItemRouter) AddItem(ctx *fasthttp.RequestCtx) {
	newItem := new(model.Item)
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
