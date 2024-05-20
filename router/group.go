package router

import (
	"fmt"
	"go/ast"
	"net/http"
	"reflect"
	"strings"

	"github.com/gopi-frame/contract/web"
	"github.com/gopi-frame/exception"
	"github.com/gopi-frame/support/lists"
	"github.com/gopi-frame/web/controller"
	"github.com/gopi-frame/web/middleware"
	"github.com/gopi-frame/web/request"
)

// Group used to manage a group of [Route]
type Group struct {
	router             *Router
	Prefix             string
	Groups             []*Group
	Routes             []*Route
	Middlewares        *lists.List[middleware.Middleware]
	ControllerInstance controller.Interface
	ControllerType     reflect.Type
}

// List lists all routes in current group
func (group *Group) List() []*Route {
	routes := make([]*Route, 0)
	for _, route := range group.Routes {
		routes = append(routes, route)
	}
	for _, routeGroup := range group.Groups {
		routes = append(routes, routeGroup.List()...)
	}
	return routes
}

// Use sets middlewares to the group
func (group *Group) Use(middlewares ...middleware.Middleware) {
	group.Middlewares.Push(middlewares...)
}

// Group registers sub route group and returns the sub group instance
func (group *Group) Group(prefix string, callback func(group *Group)) *Group {
	routeGroup := &Group{
		router: group.router,
		Prefix: strings.Join([]string{
			strings.TrimRight(group.Prefix, "/"),
			strings.TrimLeft(prefix, "/"),
		}, "/"),
		Groups:      make([]*Group, 0),
		Routes:      make([]*Route, 0),
		Middlewares: group.Middlewares,
	}
	group.Groups = append(group.Groups, routeGroup)
	callback(routeGroup)
	return routeGroup
}

// Controller registers a sub route group with specific controller instance
func (group *Group) Controller(prefix string, controller controller.Interface, callback func(group *Group)) *Group {
	controllerType := reflect.TypeOf(controller)
	if controllerType.Kind() != reflect.Pointer {
		panic(exception.NewTypeException("controller must be a pointer"))
	}
	routeGroup := &Group{
		router: group.router,
		Prefix: strings.Join([]string{
			strings.TrimRight(group.Prefix, "/"),
			strings.TrimLeft(prefix, "/"),
		}, "/"),
		Routes:             make([]*Route, 0),
		Middlewares:        group.Middlewares,
		ControllerInstance: controller,
		ControllerType:     controllerType,
	}
	callback(routeGroup)
	group.Groups = append(group.Groups, routeGroup)
	return routeGroup
}

// Route registers a handler route to current group and it returns an instance of [RouteHandler]
func (group *Group) Route(method, path string, handler any) *Route {
	pathWithPrefix := strings.Join([]string{
		strings.TrimRight(group.Prefix, "/"),
		strings.TrimLeft(path, "/"),
	}, "/")
	if path == "" {
		pathWithPrefix = pathWithPrefix[:len(pathWithPrefix)-1]
	}
	route := &Route{
		router:      group.router,
		method:      method,
		path:        pathWithPrefix,
		middlewares: group.Middlewares,
	}
	switch v := handler.(type) {
	case func(*request.Request) web.Responser:
		route.handler = v
	case string:
		if group.ControllerType == nil {
			panic(exception.NewTypeException("handler can not be type string, no controller found on current group"))
		}
		if strings.TrimSpace(v) == "" {
			panic(exception.NewEmptyArgumentException("handler"))
		}
		if !ast.IsExported(v) {
			panic(exception.NewUnexportedMethodException(v))
		}
		methodType, ok := group.ControllerType.MethodByName(v)
		if !ok {
			panic(exception.NewNoSuchMethodException(group.ControllerType, v))
		}
		if numIn := methodType.Type.NumIn(); numIn != 1 {
			panic(exception.NewTypeException(fmt.Sprintf("invalid number of input, method %s type should be func() web.Responser", v)))
		}
		if numOut := methodType.Type.NumOut(); numOut != 1 {
			panic(exception.NewTypeException(fmt.Sprintf("invalid number of output, method %s type should be func() web.Responser", v)))
		}
		if outputType := methodType.Type.Out(0); !outputType.Implements(reflect.TypeFor[web.Responser]()) {
			panic(exception.NewTypeException(fmt.Sprintf("invalid type of output, method %s type should be func() web.Responser", v)))
		}
		route.handler = func(r *request.Request) web.Responser {
			var controllerValue = reflect.New(group.ControllerType.Elem())
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
	group.Routes = append(group.Routes, route)
	return route
}

// HEAD registers a handler route with method [http.MethodHead]
func (group *Group) HEAD(path string, handler any) *Route {
	return group.Route(http.MethodHead, path, handler)
}

// CONNECT registers a handler route with method [http.MethodConnect]
func (group *Group) CONNECT(path string, handler any) *Route {
	return group.Route(http.MethodConnect, path, handler)
}

// OPTIONS registers a handler route with method [http.MethodOptions]
func (group *Group) OPTIONS(path string, handler any) *Route {
	return group.Route(http.MethodOptions, path, handler)
}

// TRACE registers a handler route with method [http.MethodTrace]
func (group *Group) TRACE(path string, handler any) *Route {
	return group.Route(http.MethodTrace, path, handler)
}

// GET registers a handler route with method [http.MethodGet]
func (group *Group) GET(path string, handler any) *Route {
	return group.Route(http.MethodGet, path, handler)
}

// POST registers a handler route with method [http.MethodPost]
func (group *Group) POST(path string, handler any) *Route {
	return group.Route(http.MethodPost, path, handler)
}

// PUT registers a handler route with method [http.MethodPut]
func (group *Group) PUT(path string, handler any) *Route {
	return group.Route(http.MethodPut, path, handler)
}

// PATCH registers a handler route with method [http.MethodPatch]
func (group *Group) PATCH(path string, handler any) *Route {
	return group.Route(http.MethodPatch, path, handler)
}

// DELETE registers a handler route with method [http.MethodDelete]
func (group *Group) DELETE(path string, handler any) *Route {
	return group.Route(http.MethodDelete, path, handler)
}
