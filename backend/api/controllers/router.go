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
	"github.com/resssoft/mediaArchive/services/translation"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"time"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

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
	itemApp := services.NewItemApp(itemRepo)
	itemRouter := NewItemRoute(itemRepo, itemApp)
	router.POST("/api/item/", itemRouter.AddItem)
	router.GET("/api/item/:id", itemRouter.GetItem)
	router.PUT("/api/item/", itemRouter.UpdateItem)
	router.DELETE("/api/item/", itemRouter.DeleteItem)
	router.GET("/api/items/", itemRouter.ItemsList)
	router.GET("/api/items/export", itemRouter.ExportItems)
	router.POST("/api/items/import", itemRouter.ImportItems)
	router.POST("/api/items/upload", itemRouter.UploadFile)

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
