package router

import (
	libctx "context"
	"errors"
	"fmt"
	"go/ast"
	"net/http"
	"reflect"
	"strings"

	pipelinecontract "github.com/gopi-frame/contract/pipeline"
	"github.com/gopi-frame/contract/validation"
	"github.com/gopi-frame/contract/web"
	"github.com/gopi-frame/exception"
	"github.com/gopi-frame/pipeline"
	"github.com/gopi-frame/support/lists"
	"github.com/gopi-frame/web/controller"
	"github.com/gopi-frame/web/middleware"
	"github.com/gopi-frame/web/middleware/cors"
	"github.com/gopi-frame/web/request"
	"github.com/julienschmidt/httprouter"
)

// ErrValidateEngineEmpty empty validation engine error
var ErrValidateEngineEmpty = errors.New("validate engine is nil, please call SetValidateEngine to set it first")

// NewRouter creates a new [Router] instance
func NewRouter() *Router {
	router := &Router{
		root: &Group{
			Prefix:      "/",
			Middlewares: lists.NewList[middleware.Middleware](),
			Groups:      make([]*Group, 0),
			Routes:      make([]*Route, 0),
		},
		HTTPRouter: httprouter.New(),
	}
	router.root.router = router
	router.HTTPRouter.PanicHandler = func(w http.ResponseWriter, r *http.Request, i interface{}) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	return router
}

// Router http router
type Router struct {
	root           *Group
	HTTPRouter     *httprouter.Router
	routes         []*Route
	validateEngine validation.Engine
	cors           *cors.CORS
}

// SetCORS sets cors configs
func (router *Router) SetCORS(options ...cors.Option) *Router {
	router.cors = cors.New(options...)
	return router
}

// SetValidateEngine sets custom validate engine
func (router *Router) SetValidateEngine(ve validation.Engine) *Router {
	router.validateEngine = ve
	return router
}

// SetErrorHandler sets custom error handler
func (router *Router) SetErrorHandler(handler func(*http.Request, error) web.Responser) *Router {
	router.HTTPRouter.PanicHandler = func(w http.ResponseWriter, r *http.Request, i interface{}) {
		switch v := i.(type) {
		case error:
			handler(r, v).Send(w, r)
		default:
			handler(r, exception.NewValueException(v)).Send(w, r)
		}
	}
	return router
}

// Register used to registe routes by custom callback
func (router *Router) Register(register func(router *Router)) {
	register(router)
}

// Run starts the http server
func (router *Router) Run(addr string) error {
	if router.routes == nil {
		router.routes = router.List()
	}
	for _, route := range router.routes {
		router.HTTPRouter.Handle(route.Method(), route.Path(), func(route *Route) httprouter.Handle {
			return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
				ctx := r.Context()
				ctx = libctx.WithValue(ctx, httprouter.ParamsKey, p)
				req := request.NewRequest(r, p)
				lines := pipeline.New[*request.Request, web.Responser]()
				lines = lines.Send(req)
				stops := make([]pipelinecontract.Stop[*request.Request, web.Responser], 0)
				route.middlewares.Each(func(_ int, middleware middleware.Middleware) bool {
					stops = append(stops, pipeline.Stop(middleware.Handle))
					return true
				})
				if route.HasValidation() {
					stops = append(stops, pipeline.Stop(route.validation.Handle))
				}
				resp := lines.Then(route.handler)
				resp.Send(w, r)
			}
		}(route))
	}
	if router.cors != nil {
		router.HTTPRouter.HandleOPTIONS = true
		router.HTTPRouter.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cors := router.cors
			w.Header().Set("Access-Control-Allow-Credentials", fmt.Sprintf("%v", cors.AllowCredentials))
			if len(cors.AllowHeaders) > 0 {
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(cors.AllowHeaders, ","))
			}
			if len(cors.AllowMethods) > 0 {
				w.Header().Set("Access-Control-Request-Method", strings.Join(cors.AllowMethods, ","))
			} else {
				w.Header().Set("Access-Control-Request-Method", w.Header().Get("Allow"))
			}
			if len(cors.AllowOrigins) > 0 {
				w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
			} else {
				w.Header().Set("Access-Control-Allow-Origin", strings.Join(cors.AllowOrigins, ","))
			}
			if len(cors.ExposeHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(cors.ExposeHeaders, ","))
			}
			w.WriteHeader(http.StatusNoContent)
		})
	}
	return http.ListenAndServe(addr, router.HTTPRouter)
}

// List lists all routes in current group
func (router *Router) List() []*Route {
	routes := make([]*Route, 0)
	for _, route := range router.root.Routes {
		routes = append(routes, route)
	}
	for _, routeGroup := range router.root.Groups {
		routes = append(routes, routeGroup.List()...)
	}
	return routes
}

// Use sets middlewares to the group
func (router *Router) Use(middlewares ...middleware.Middleware) {
	router.root.Middlewares.Push(middlewares...)
}

// Group registers sub route group and returns the sub group instance
func (router *Router) Group(prefix string, callback func(group *Group)) *Group {
	routeGroup := &Group{
		router: router.root.router,
		Prefix: strings.Join([]string{
			strings.TrimRight(router.root.Prefix, "/"),
			strings.TrimLeft(prefix, "/"),
		}, "/"),
		Groups:      make([]*Group, 0),
		Routes:      make([]*Route, 0),
		Middlewares: router.root.Middlewares,
	}
	router.root.Groups = append(router.root.Groups, routeGroup)
	callback(routeGroup)
	return routeGroup
}

// Controller registers a sub route group with specific controller instance
func (router *Router) Controller(prefix string, controller controller.Interface, callback func(group *Group)) *Group {
	controllerType := reflect.TypeOf(controller)
	if controllerType.Kind() != reflect.Pointer {
		panic(exception.NewTypeException("controller must be a pointer"))
	}
	routeGroup := &Group{
		router: router.root.router,
		Prefix: strings.Join([]string{
			strings.TrimRight(router.root.Prefix, "/"),
			strings.TrimLeft(prefix, "/"),
		}, "/"),
		Routes:             make([]*Route, 0),
		Middlewares:        router.root.Middlewares,
		ControllerInstance: controller,
		ControllerType:     controllerType,
	}
	callback(routeGroup)
	router.root.Groups = append(router.root.Groups, routeGroup)
	return routeGroup
}

// Route registers a handler route to current group and it returns an instance of [RouteHandler]
func (router *Router) Route(method, path string, handler any) *Route {
	pathWithPrefix := strings.Join([]string{
		strings.TrimRight(router.root.Prefix, "/"),
		strings.TrimLeft(path, "/"),
	}, "/")
	if path == "" {
		pathWithPrefix = pathWithPrefix[:len(pathWithPrefix)-1]
	}
	route := &Route{
		router:      router.root.router,
		method:      method,
		path:        pathWithPrefix,
		middlewares: router.root.Middlewares,
	}
	switch v := handler.(type) {
	case func(*request.Request) web.Responser:
		route.handler = v
	case string:
		if router.root.ControllerType == nil {
			panic(exception.NewTypeException("handler can not be type string, no controller found on current group"))
		}
		if strings.TrimSpace(v) == "" {
			panic(exception.NewEmptyArgumentException("handler"))
		}
		if !ast.IsExported(v) {
			panic(exception.NewUnexportedMethodException(v))
		}
		methodType, ok := router.root.ControllerType.MethodByName(v)
		if !ok {
			panic(exception.NewNoSuchMethodException(router.root.ControllerType, v))
		}
		if numIn := methodType.Type.NumIn(); numIn != 1 {
			panic(exception.NewTypeException(fmt.Sprintf("invalid number of input, method %s type should be func(*context.Request) web.Responser", v)))
		}
		if numOut := methodType.Type.NumOut(); numOut != 1 {
			panic(exception.NewTypeException(fmt.Sprintf("invalid number of output, method %s type should be func(*context.Request) web.Responser", v)))
		}
		if outputType := methodType.Type.Out(0); !outputType.Implements(reflect.TypeFor[web.Responser]()) {
			panic(exception.NewTypeException(fmt.Sprintf("invalid type of output, method %s type should be func(*context.Request) web.Responser", v)))
		}
		route.handler = func(r *request.Request) web.Responser {
			var controllerValue = reflect.New(router.root.ControllerType.Elem())
			controllerValue.MethodByName("Init").Call([]reflect.Value{
				reflect.ValueOf(r),
			})
			outputs := controllerValue.MethodByName(v).Call([]reflect.Value{})
			resp := outputs[0].Interface().(web.Responser)
			return resp
		}
	default:
		panic(exception.NewTypeException("invalid handler type, only string and func(*context.Request) web.Responser are allowed"))
	}
	router.root.Routes = append(router.root.Routes, route)
	return route
}

// HEAD registers a handler route with method [http.MethodHead]
func (router *Router) HEAD(path string, handler any) *Route {
	return router.root.Route(http.MethodHead, path, handler)
}

// CONNECT registers a handler route with method [http.MethodConnect]
func (router *Router) CONNECT(path string, handler any) *Route {
	return router.root.Route(http.MethodConnect, path, handler)
}

// OPTIONS registers a handler route with method [http.MethodOptions]
func (router *Router) OPTIONS(path string, handler any) *Route {
	return router.root.Route(http.MethodOptions, path, handler)
}

// TRACE registers a handler route with method [http.MethodTrace]
func (router *Router) TRACE(path string, handler any) *Route {
	return router.root.Route(http.MethodTrace, path, handler)
}

// GET registers a handler route with method [http.MethodGet]
func (router *Router) GET(path string, handler any) *Route {
	return router.root.Route(http.MethodGet, path, handler)
}

// POST registers a handler route with method [http.MethodPost]
func (router *Router) POST(path string, handler any) *Route {
	return router.root.Route(http.MethodPost, path, handler)
}

// PUT registers a handler route with method [http.MethodPut]
func (router *Router) PUT(path string, handler any) *Route {
	return router.root.Route(http.MethodPut, path, handler)
}

// PATCH registers a handler route with method [http.MethodPatch]
func (router *Router) PATCH(path string, handler any) *Route {
	return router.root.Route(http.MethodPatch, path, handler)
}

// DELETE registers a handler route with method [http.MethodDelete]
func (router *Router) DELETE(path string, handler any) *Route {
	return router.root.Route(http.MethodDelete, path, handler)
}
