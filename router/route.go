package router

import (
	"reflect"
	"runtime"

	"github.com/gopi-frame/support/lists"
	validationcontract "github.com/gopi-frame/validation/contract"
	"github.com/gopi-frame/web/contract"
	"github.com/gopi-frame/web/middleware"
	"github.com/gopi-frame/web/middleware/validate"
	"github.com/gopi-frame/web/request"
)

// Route basic route struct
type Route struct {
	router      *Router
	name        string
	method      string
	path        string
	middlewares *lists.List[middleware.Middleware]
	validation  middleware.Middleware
	handler     func(*request.Request) contract.Responser
}

// AS sets the name
func (route *Route) AS(name string) *Route {
	route.name = name
	return route
}

// Name returns the name of route
func (route *Route) Name() string {
	return route.name
}

// Method returns the method of route
func (route *Route) Method() string {
	return route.method
}

// Path returns the path of route
func (route *Route) Path() string {
	return route.path
}

// Middlewares returns the middlewares of the route
func (route *Route) Middlewares() []middleware.Middleware {
	return route.middlewares.ToArray()
}

// HasValidation returns whether the route has a binded validation
func (route *Route) HasValidation() bool {
	return route.validation != nil
}

// Use sets middlewares
func (route *Route) Use(middlewares ...middleware.Middleware) *Route {
	route.middlewares.Push(middlewares...)
	return route
}

// Handler returns the handler's name
func (route *Route) Handler() string {
	return runtime.FuncForPC(reflect.ValueOf(route.handler).Pointer()).Name()
}

// Validate binds validation form to the route
func (route *Route) Validate(form validationcontract.Form, bindings ...contract.Resolver) *Route {
	if route.router.validateEngine == nil {
		panic(ErrValidateEngineEmpty)
	}
	route.validation = validate.For(form, route.router.validateEngine, bindings...)
	return route
}
