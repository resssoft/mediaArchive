package controller

import (
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/resssoft/mediaArchive/app"
	config "github.com/resssoft/mediaArchive/configuration"
	"github.com/resssoft/mediaArchive/database"
	"github.com/resssoft/mediaArchive/model"
	"github.com/resssoft/mediaArchive/repository"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

func CORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next(ctx)
	}
}

func Routing(db database.MongoClientApplication, address string) error {
	router := fasthttprouter.New()
	router.GET("/api/version", version)

	itemRepo := repository.NewItemRepo(db)
	itemApp := app.NewItemApp(itemRepo)
	itemRouter := NewUserRoute(itemRepo, itemApp)
	router.POST("/api/item/", itemRouter.AddItem)
	router.GET("/api/item/:id", itemRouter.GetItem)
	router.PUT("/api/item/", itemRouter.UpdateItem)
	router.DELETE("/api/item/", itemRouter.DeleteItem)
	router.GET("/api/items/", itemRouter.ItemsList)
	router.GET("/api/items/export", itemRouter.ExportItems)

	log.Info().Msg("Launched under version: " + config.Version)
	log.Info().Msg("Start by address: " + address)
	return fasthttp.ListenAndServe(address, CORS(router.Handler))
}

func version(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	fmt.Fprintf(ctx, "{version:%s}", config.Version)
}

func getError(msg string, code int) model.RequestError {
	return model.RequestError{
		Error: msg,
		Code:  code,
	}
}

func writeJsonResponse(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(code)
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		log.Err(err).Send()
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}
