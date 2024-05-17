package middleware

import (
	"github.com/gopi-frame/contract/web"
	"github.com/gopi-frame/web/request"
)

// Middleware middleware
type Middleware interface {
	Handle(*request.Request, func(*request.Request) web.Responser) web.Responser
}
