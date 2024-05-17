package binding

import (
	"net/http"

	"github.com/gopi-frame/exception"
	"github.com/pelletier/go-toml/v2"
)

// TOML implements [Parser], it will parse request body into container
// Make sure your request body is not nil and is TOML-encoded, or it will return an error
//
// Example:
//
//	var container = &struct{
//	    Name string `toml:"name"`
//	    Age int `toml:"age"`
//	    Tags []string `toml:"tags"`
//	    VIP bool `toml:"vip"`
//	}{}
//
//	err := TOML(request, container)
//	if err != nil {
//	    panic(err)
//	}
var TOML Binding = func(request *http.Request, container any) error {
	if request == nil || request.Body == nil {
		return exception.NewEmptyArgumentException("request")
	}
	decoder := toml.NewDecoder(request.Body)
	return decoder.Decode(container)
}
