package binding

import (
	"encoding/xml"
	"net/http"

	"github.com/gopi-frame/exception"
)

// XML implements [Parser], it parses request body into container
// Make sure your request body is not nil and is XML-Encoded, or it will returns an error
//
// Example:
//
//	var container = &struct{
//	    Name string `xml:"name"`
//	    Age int `xml:"age"`
//	    Tags []string `xml:"tags"`
//	    VIP bool `xml:"vip"`
//	}{}
//
//	err := XML(request, container)
//	if err != nil {
//	    panic(err)
//	}
var XML Binding = func(request *http.Request, container any) error {
	if request == nil || request.Body == nil {
		return exception.NewEmptyArgumentException("request")
	}
	decoder := xml.NewDecoder(request.Body)
	return decoder.Decode(container)
}
