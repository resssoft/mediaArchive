package translation

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	config "github.com/resssoft/mediaArchive/configuration"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

var (
	onceTranslationAction sync.Once
)

type TranslatorApplication interface {
	Msg(string, string, string, map[string]interface{}) string
	AvailableLanguages() []string
}

type TranslatorApp struct {
	bundle     *i18n.Bundle
	localizers map[string]*i18n.Localizer
}

func ProvideTranslator() (TranslatorApplication, error) {
	var err error
	app := TranslatorApp{}
	onceTranslationAction.Do(func() {
		err = configureTranslation(&app, config.TranslationDir)
	})
	if err != nil {
		return &app, err
	}
	return &app, nil
}

func configureTranslation(app *TranslatorApp, dir string) error {
	app.bundle = i18n.NewBundle(language.English)
	app.bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	fileList, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Err(err).Send()
		return err
	}
	app.localizers = make(map[string]*i18n.Localizer, 0)

	for _, val := range fileList {
		if !val.IsDir() {
			_, fileName := filepath.Split(val.Name())
			if filepath.Ext(fileName) == ".toml" {
				_, err := app.bundle.LoadMessageFile(dir + val.Name())
				if err != nil {
					log.Err(err).Send()
				}
			}
			lang := strings.TrimSuffix(fileName, filepath.Ext(fileName))
			app.localizers[lang] = i18n.NewLocalizer(app.bundle, lang)
		}
	}
	return nil
}

func (t *TranslatorApp) Msg(lang, id, defaultMsg string, data map[string]interface{}) string {
	loc, ok := t.localizers[lang]
	if ok {
		return loc.MustLocalize(&i18n.LocalizeConfig{
			TemplateData:   data,
			DefaultMessage: &i18n.Message{ID: id, Other: defaultMsg},
		})
	}
	t.bundle.LanguageTags()
	return defaultMsg
}

func (t *TranslatorApp) AvailableLanguages() []string {
	langList := make([]string, 0)
	for _, langTag := range t.bundle.LanguageTags() {
		langList = append(langList, langTag.String())
	}
	return langList
}
