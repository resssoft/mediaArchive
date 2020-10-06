package controllers

import (
	"github.com/resssoft/mediaArchive/models"
	"github.com/resssoft/mediaArchive/services/translation"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"net/http"
)

type TranslationRouter struct {
	app translation.TranslatorApplication
}

func NewTranslationRoute(app translation.TranslatorApplication) *TranslationRouter {
	return &TranslationRouter{app: app}
}

func (r *TranslationRouter) LanguagesList(ctx *fasthttp.RequestCtx) {
	languages := r.app.AvailableLanguages()
	log.Info().Interface("lngs", languages).Send()
	writeJsonResponse(ctx, http.StatusOK, models.Response{
		Data:  languages,
		Count: len(languages),
	})
	log.Info().Msg(r.app.Msg("en", "HelloPerson", "Hello", map[string]interface{}{
		"Name": "username",
	}))
	log.Info().Msg(r.app.Msg("ru", "HelloPerson", "Hello", map[string]interface{}{
		"Name": "username",
	}))
}
