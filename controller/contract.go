package controller

import "github.com/gopi-frame/web/request"

// Interface controller interface
type Interface interface {
	Init(request *request.Request)
}
