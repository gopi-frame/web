package reqid

import (
	"github.com/google/uuid"
	"github.com/gopi-frame/web/contract"
	"github.com/gopi-frame/web/request"
)

const defaultHeaderKey = "X-Request-ID"

// New new request id middleware
func New(options ...Option) *RequestID {
	r := new(RequestID)
	r.Header = defaultHeaderKey
	r.Generator = uuid.NewString
	for _, option := range options {
		option(r)
	}
	return r
}

// RequestID request id generator
type RequestID struct {
	Header    string
	Generator func() string
}

// Handle handle
func (r *RequestID) Handle(request *request.Request, next func(*request.Request) contract.Responser) contract.Responser {
	var id = r.Generator()
	resp := next(request)
	if resp != nil && !resp.HasHeader(r.Header) {
		resp.SetHeader(r.Header, id)
	}
	return resp
}
