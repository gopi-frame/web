package validate

import (
	"net/http"
	"reflect"

	"github.com/gopi-frame/contract/validation"
	"github.com/gopi-frame/contract/web"
	"github.com/gopi-frame/web/request"
	"github.com/gopi-frame/web/response"
)

// New new validation middleware
func New[F validation.Form](validator validation.Engine, bindings ...web.Resolver) *Validate {
	v := new(Validate)
	v.validator = validator
	v.formType = reflect.TypeFor[F]()
	v.bindings = bindings
	return v
}

// For validate middleware for
func For(form validation.Form, validator validation.Engine, bindings ...web.Resolver) *Validate {
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
	validator validation.Engine
	bindings  []web.Resolver
	formType  reflect.Type
}

// Handle handle
func (v *Validate) Handle(request *request.Request, next func(*request.Request) web.Responser) web.Responser {
	form := reflect.New(v.formType).Interface().(validation.Form)
	if err := request.Bind(form, v.bindings...); err != nil {
		return response.NewResponse(http.StatusBadRequest, err.Error())
	}
	locale := request.Locale()
	if locale == "" {
		locale = "en"
	}
	form.SetLocale(locale)
	v.validator.ValidateForm(form)
	return next(request)
}
