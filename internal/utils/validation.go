package utils

import (
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	log "github.com/sirupsen/logrus"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func NewValidator() (*validator.Validate, *ut.Translator) {
	log.Debug("Initializing validator")
	validate := validator.New(validator.WithRequiredStructEnabled())

	translator := newTranslator(validate)

	return validate, translator
}

func newTranslator(validate *validator.Validate) *ut.Translator {
	log.Debug("Initializing validator error translations")

	fieldNameExceptions := map[string]string{
		"indexlimit": "index_limit",
	}

	en := en.New()
	utTranslation := ut.New(en, en)
	translator, _ := utTranslation.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, translator)

	validate.RegisterTranslation("required", translator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})

	validate.RegisterTranslation("required", translator, func(ut ut.Translator) error {
		return ut.Add("mongodb", "{0} must be a mongodb ObjectId", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		name := strings.ToLower(fe.Field())

		changedName, ok := fieldNameExceptions[name]
		if ok {
			name = changedName
		}

		t, _ := ut.T("required", name)
		return t
	})

	return &translator
}