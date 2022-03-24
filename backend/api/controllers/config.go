package controllers

import (
	"github.com/resssoft/mediaArchive/models"

	"fmt"
	"github.com/resssoft/mediaArchive/pkg/requestFilter"
	"github.com/resssoft/mediaArchive/repositories"
	configApp "github.com/resssoft/mediaArchive/services/config"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

type ConfigRouter struct {
	repo repositories.ConfigRepository
	app  configApp.ConfigApplication
}

func NewConfigRoute(repo repositories.ConfigRepository, app configApp.ConfigApplication) *ConfigRouter {
	return &ConfigRouter{repo: repo, app: app}
}

func (r *ConfigRouter) ConfigData(ctx *fasthttp.RequestCtx) {
	var filter models.DataFilter
	entityType := fmt.Sprintf("%v", ctx.UserValue("group"))
	entityName := fmt.Sprintf("%v", ctx.UserValue("name"))

	if entityName == "" || ctx.UserValue("name") == nil {

		filter.Append(
			"id",
			entityType,
		)
	} else {
		filter.Append(
			"group",
			entityType,
		)
		filter.Append(
			"name",
			entityName,
		)
	}

	log.Info().Interface("filter", filter).Send()
	items, err := r.app.List(filter)
	if err != nil {
		writeJsonResponse(ctx, http.StatusBadRequest, getError("config error: "+err.Error(), 1))
		return
	}
	if len(items) == 0 {
		writeJsonResponse(ctx, http.StatusBadRequest, getError("config not found", 1))
		return
	}

	log.Info().Msg(ctx.Request.String())
	log.Info().Msg(entityType + entityName)

	ctx.SetContentType("application/json; charset=utf8")
	ctx.Response.Header.SetCanonical([]byte("Cache-Control"), []byte("max-age=3600"))
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.SetBody([]byte(items[0].Data))
}

/*
func (r *ConfigRouter) ConfigDatax(ctx *fasthttp.RequestCtx) {
	entityType := fmt.Sprintf("%v", ctx.UserValue("group"))
	entityName := fmt.Sprintf("%v", ctx.UserValue("name"))
	id, err := strconv.Atoi(fmt.Sprintf("%v", ctx.UserValue("group")))
	if (entityName == "" || ctx.UserValue("name") == nil) && err == nil && id > 0 {
		item, err := r.systemService.Config().GetById(id)
		if err != nil {
			writeJsonResponse(ctx, http.StatusBadRequest, getError("error: no config by id found", 1))
			return
		}
		writeJsonResponse(ctx, http.StatusOK, dataResponse(1, 1, 0, item))
		return
	}
	json := fmt.Sprintf("[{\"Condition\":\"ilike\",\"Data\":{\"Name\":\"%v\"}},{\"Condition\":\"ilike\",\"Data\":{\"Group\":\"%v\"}}]", entityName, entityType)
	data, err := r.systemService.Config().GetList("id", 1, 1, []byte(json))
	if err != nil {
		writeJsonResponse(ctx, http.StatusBadRequest, getError("config error: "+err.Error(), 1))
		return
	}
	if len(data) == 0 {
		writeJsonResponse(ctx, http.StatusBadRequest, getError("error: no config found", 2))
		return
	}
	ctx.SetContentType("application/json; charset=utf8")
	ctx.Response.Header.SetCanonical([]byte("Cache-Control"), []byte("max-age=3600"))
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.SetBody([]byte(data[0].Data))
}
*/

func (r *ConfigRouter) ConfigList(ctx *fasthttp.RequestCtx) {
	var filter models.DataFilter
	reqfilters, err := requestFilter.BuildFilter(argListLowCase(ctx), ctx.PostBody())
	if err != nil {
		writeJsonResponse(ctx, http.StatusBadRequest, getError("filter error: "+err.Error(), 1))
		return
	}
	log.Info().Interface("filter by list", reqfilters).Send()
	filter.AppendFromRequestFilter(reqfilters)
	items, err := r.app.List(filter)
	if err != nil {
		writeJsonResponse(ctx, http.StatusBadRequest, getError("list error: "+err.Error(), 1))
		return
	}
	writeJsonResponse(ctx, http.StatusOK, dataResponse(len(items), len(items), 0, items))
}

func (r *ConfigRouter) ConfigGet(ctx *fasthttp.RequestCtx) {
	id, err := strconv.Atoi(fmt.Sprintf("%v", ctx.UserValue("id")))
	var filter models.DataFilter
	filter.Append("id", id)
	items, err := r.app.List(filter)
	if err != nil {
		writeJsonResponse(ctx, http.StatusBadRequest, getError("error: "+err.Error(), 1))
		return
	}
	writeJsonResponse(ctx, http.StatusOK, dataResponse(1, 1, 0, items[0]))
}

func (r *ConfigRouter) ConfigUpdate(ctx *fasthttp.RequestCtx) {
	item, err := r.app.Validate(ctx.PostBody())
	if err != nil {
		writeJsonResponse(ctx, http.StatusBadRequest, getError(err.Error(), 1))
		return
	}
	err = r.app.Update(item)
	if err != nil {
		writeJsonResponse(ctx, http.StatusBadRequest, getError(err.Error(), 1))
		return
	}
	writeJsonResponse(ctx, http.StatusNoContent, nil)
}
