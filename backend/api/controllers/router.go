package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/resssoft/mediaArchive/auth"
	config "github.com/resssoft/mediaArchive/configuration"
	"github.com/resssoft/mediaArchive/database"
	"github.com/resssoft/mediaArchive/models"
	"github.com/resssoft/mediaArchive/repositories"
	"github.com/resssoft/mediaArchive/services"
	configApp "github.com/resssoft/mediaArchive/services/config"
	"github.com/resssoft/mediaArchive/services/translation"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

type DataResponse struct {
	Total int         `json:"total"`
	Count int         `json:"count"`
	Page  int         `json:"page"`
	Data  interface{} `json:"data"`
}

func CORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		next(ctx)
	}
}

func Routing(db database.MongoClientApplication, tr translation.TranslatorApplication, address string) error {
	authMiddleware := AuthMiddleware(WithExpirationValidator())
	router := fasthttprouter.New()
	router.GET("/api/version", version)

	itemRepo := repositories.NewItemRepo(db)
	itemGroupRepo := repositories.NewItemGroupRepo(db)
	configRepo := repositories.NewConfigRepo(db)

	itemApp := services.NewItemApp(itemRepo, itemGroupRepo)
	configApp := configApp.NewItemApp(configRepo)

	itemRouter := NewItemRoute(itemRepo, itemApp)
	configRouter := NewConfigRoute(configRepo, configApp)

	router.POST("/api/item/", authMiddleware.Wrap(itemRouter.AddItem))
	router.GET("/api/item/:id", authMiddleware.Wrap(itemRouter.GetItem))
	router.PUT("/api/item/:id", authMiddleware.Wrap(itemRouter.UpdateItem))
	router.PATCH("/api/item/:id", authMiddleware.Wrap(itemRouter.UpdateItem))
	router.DELETE("/api/item/", authMiddleware.Wrap(itemRouter.DeleteItem))
	router.GET("/api/items/", authMiddleware.Wrap(itemRouter.ItemsList))
	router.POST("/api/items", authMiddleware.Wrap(itemRouter.ItemsList))
	router.GET("/api/items/export", authMiddleware.Wrap(itemRouter.ExportItems))
	router.POST("/api/items/import", authMiddleware.Wrap(itemRouter.ImportItems))
	router.POST("/api/items/upload", authMiddleware.Wrap(itemRouter.UploadFile))

	router.POST("/api/item-group/", authMiddleware.Wrap(itemRouter.AddItemGroup))
	router.GET("/api/item-groups/", authMiddleware.Wrap(itemRouter.ItemsGroups))
	router.POST("/api/item-groups/", authMiddleware.Wrap(itemRouter.ItemsGroups))

	//router.GET("/api/config/:group", authMiddleware.Wrap(configRouter.ConfigData))
	router.GET("/api/config/:group", authMiddleware.Wrap(configRouter.ConfigData))
	router.GET("/api/config/:group/:name", authMiddleware.Wrap(configRouter.ConfigData))
	router.POST("/api/config/list", authMiddleware.Wrap(configRouter.ConfigList))
	//router.GET("/api/config/:id", authMiddleware.Wrap(configRouter.ConfigGet))
	//router.POST("/api/config/", authMiddleware.Wrap(configRouter.ConfigAdd))
	//router.PATCH("/api/config/", authMiddleware.Wrap(configRouter.ConfigUpdate))
	//router.DELETE("/api/config/:id", authMiddleware.Wrap(configRouter.ConfigDelete))
	router.GET("/api/user", authMiddleware.Wrap(usettest))

	userRepo := repositories.NewUserRepo(db)
	userApp := services.NewUserApp(userRepo)
	userApp.AddOwner()
	userRouter := NewUserRoute(userRepo, userApp)
	router.POST("/api/auth/sign-in", userRouter.Login)
	router.POST("/api/auth/sign-out", authMiddleware.Wrap(userRouter.Logout))
	router.GET("/api/auth/refresh-token", AuthMiddleware().Wrap(userRouter.RefreshToken))
	router.GET("/api/user/info/", authMiddleware.Wrap(userRouter.UserInfo))
	router.POST("/api/user/", authMiddleware.Wrap(userRouter.AddUser))

	transactionRouter := NewTranslationRoute(tr)
	router.GET("/api/languages/", transactionRouter.LanguagesList)

	log.Info().Msg("Launched under version: " + config.Version)
	log.Info().Msg("Start by address: " + address)
	return fasthttp.ListenAndServe(address, CORS(router.Handler))
}

func version(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	fmt.Fprintf(ctx, "{version:%s}", config.Version)
}

func usettest(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	fmt.Fprint(ctx, `{"total":1,"count":1,"page":0,"data":{"Id":634,"Email":"admin@admin.net","Password":"","Role":"admin"}}`)
}

func getError(msg string, code int) models.RequestError {
	return models.RequestError{
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

type Middleware interface {
	Wrap(fasthttp.RequestHandler) fasthttp.RequestHandler
}

type SimpleMiddleware func(fasthttp.RequestHandler) fasthttp.RequestHandler

func (m SimpleMiddleware) Wrap(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return m(next)
}

func AuthMiddleware(opts ...jwt.Validator) SimpleMiddleware {
	handler := func(ctx *fasthttp.RequestCtx) error {
		token := auth.ExtractBearerToken(ctx)
		payload, err := auth.VerifyToken(token, opts...)
		if err != nil {
			return err
		}
		log.Info().Msg(ctx.Request.String())
		ctx.SetUserValue("userRole", payload.Perms)
		ctx.SetUserValue("userId", payload.UserId)
		ctx.SetUserValue("userLang", payload.UserLang)
		ctx.SetUserValue("session", payload.Payload.JWTID)
		return nil
	}
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			if err := handler(ctx); err != nil {
				log.Debug().Err(err).Msg("Token verification error")
				ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			} else {
				next(ctx)
			}
		}
	}
}

func WithExpirationValidator() jwt.Validator {
	return jwt.ExpirationTimeValidator(time.Now())
}

func dataResponse(total, count, page int, data interface{}) DataResponse {
	return DataResponse{
		Total: total,
		Count: count,
		Page:  page,
		Data:  data,
	}
}

func argListLowCase(ctx *fasthttp.RequestCtx) map[string]string {
	args := make(map[string]string, 0)
	ctx.QueryArgs().VisitAll(func(key, value []byte) {
		key = []byte(strings.ToLower(string(key)))
		args[strings.ToLower(string(key))] = string(value)
	})
	return args
}
