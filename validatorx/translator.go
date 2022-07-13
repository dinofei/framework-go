package validatorx

import (
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/gotomicro/ego/core/elog"
)

const (
	EN = "en"
	ZH = "zh"
)

type (
	translator struct {
		translatorFunc func() locales.Translator
		registerFunc   func(v *validator.Validate, trans ut.Translator) error
	}
	rule struct {
		tag         string
		translation string
		override    bool
	}
)

var (
	translators = map[string]translator{
		EN: {en.New, entranslations.RegisterDefaultTranslations},
		ZH: {zh.New, zhtranslations.RegisterDefaultTranslations},
	}
	rules = map[string][]rule{
		EN: {
			{"required_if", "{0} is a required field", false},
		},
		ZH: {
			{"required_if", "{0}为必填字段", false},
		},
	}
)

func newTranslator(v *validator.Validate) (*ut.UniversalTranslator, error) {
	var uni *ut.UniversalTranslator
	for lan, t := range translators {
		vt := t.translatorFunc()
		if uni == nil {
			uni = ut.New(vt)
		} else {
			err := uni.AddTranslator(vt, true)
			if err != nil {
				return nil, err
			}
		}
		tt, _ := uni.GetTranslator(lan)
		err := t.registerFunc(v, tt)
		if err != nil {
			return nil, err
		}
		if rs, ok := rules[lan]; ok {
			for _, r := range rs {
				_ = v.RegisterTranslation(
					r.tag,
					tt,
					registerFunc(r.tag, r.translation, r.override),
					translateFunc)
			}
		}
	}
	return uni, nil
}

func registerFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, override); err != nil {
			return
		}

		return
	}
}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		elog.Error("warning: error translating FieldError", elog.FieldErr(fe))
		return fe.(error).Error()
	}

	return t
}
