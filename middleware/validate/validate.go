package validate

import (
	"net/http"
	"reflect"

	validationcontract "github.com/gopi-frame/validation/contract"
	"github.com/gopi-frame/web/contract"
	"github.com/gopi-frame/web/request"
	"github.com/gopi-frame/web/response"
)

// New new validation middleware
func New[F validationcontract.Form](validator validationcontract.Engine, bindings ...contract.Resolver) *Validate {
	v := new(Validate)
	v.validator = validator
	v.formType = reflect.TypeFor[F]()
	v.bindings = bindings
	return v
}

// For validate middleware for
func For(form validationcontract.Form, validator validationcontract.Engine, bindings ...contract.Resolver) *Validate {
	v := new(Validate)
	v.validator = validator
	v.formType = reflect.TypeOf(form)
	if v.formType.Kind() == reflect.Pointer {
		v.formType = v.formType.Elem()
	}
	v.bindings = bindings
	return v
}

// Validate validation middleware
type Validate struct {
	validator validationcontract.Engine
	bindings  []contract.Resolver
	formType  reflect.Type
}

// Handle handle
func (v *Validate) Handle(request *request.Request, next func(*request.Request) contract.Responser) contract.Responser {
	form := reflect.New(v.formType).Interface().(validationcontract.Form)
	if err := request.Bind(form, v.bindings...); err != nil {
		return response.New(http.StatusBadRequest, err.Error())
	}
	locale := request.Locale()
	if locale == "" {
		locale = "en"
	}
	form.SetLocale(locale)
	v.validator.ValidateForm(form)
	return next(request)
}
