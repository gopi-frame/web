package binding

import (
	"mime/multipart"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// URI implements [Parser], it parses httprouter.[Params] into container
var URI Binding = func(request *http.Request, container any) error {
	form := &multipart.Form{
		Value: make(map[string][]string),
		File:  make(map[string][]*multipart.FileHeader),
	}
	params := httprouter.ParamsFromContext(request.Context())
	for _, param := range params {
		form.Value[param.Key] = []string{param.Value}
	}
	return formtostruct(form, container, "param")
}
