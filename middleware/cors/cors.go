package cors

import (
	"fmt"
	"strings"

	"github.com/gopi-frame/web/contract"
	"github.com/gopi-frame/web/request"
)

// New new cors middleware
func New(options ...Option) *CORS {
	cors := new(CORS)
	for _, option := range options {
		option(cors)
	}
	return cors
}

// CORS cors middleware
type CORS struct {
	AllowCredentials bool
	AllowHeaders     []string
	AllowOrigins     []string
	AllowMethods     []string
	ExposeHeaders    []string
}

// Handle handle
func (cors *CORS) Handle(request *request.Request, next func(*request.Request) contract.Responser) contract.Responser {
	resp := next(request)
	resp.Headers().Set("Access-Control-Allow-Credentials", fmt.Sprintf("%v", cors.AllowCredentials))
	if len(cors.AllowHeaders) > 0 {
		resp.Headers().Set("Access-Control-Allow-Headers", strings.Join(cors.AllowHeaders, ","))
	}
	if len(cors.AllowOrigins) > 0 {
		resp.Headers().Set("Access-Control-Allow-Origin", strings.Join(cors.AllowOrigins, ","))
	} else {
		resp.Headers().Set("Access-Control-Allow-Origin", request.Header("Origin").String())
	}
	if len(cors.AllowMethods) > 0 {
		resp.Headers().Set("Access-Control-Request-Method", strings.Join(cors.AllowMethods, ","))
	} else {
		resp.Headers().Set("Access-Control-Request-Method", resp.Headers().Get("Allow"))
	}
	if len(cors.ExposeHeaders) > 0 {
		resp.Headers().Set("Access-Control-Expose-Headers", strings.Join(cors.ExposeHeaders, ","))
	}
	return resp
}
