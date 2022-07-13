package validatorx

import (
	"encoding/json"
	"fmt"
	"reflect"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	validate     *validator.Validate
	translatorer *ut.UniversalTranslator
)

func init() {
	validate = New()
	translatorer, _ = newTranslator(validate)
}

func New() *validator.Validate {
	v := validator.New()
	v.SetTagName("binding")
	registerTagName(v)
	return v
}

func registerTagName(v *validator.Validate) {
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})
}

func GetValidate() *validator.Validate {
	return validate
}

func GetTranslator(lan string) ut.Translator {
	t, _ := translatorer.GetTranslator(lan)
	return t
}

func BindJson(target interface{}, dest []byte) error {
	err := json.Unmarshal(dest, target)
	if err != nil {
		return fmt.Errorf("json.Unmarshal (%s) error: %v", dest, err)
	}
	return validate.Struct(target)
}

func Struct(target interface{}) error {
	return validate.Struct(target)
}
