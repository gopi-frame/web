package middleware

import (
	"github.com/gopi-frame/web/contract"
	"github.com/gopi-frame/web/request"
)

// Middleware middleware
type Middleware interface {
	Handle(*request.Request, func(*request.Request) contract.Responser) contract.Responser
}
