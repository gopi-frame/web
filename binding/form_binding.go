package binding

import (
	"mime/multipart"
	"net/http"
)

// DefaultMaxMemory is the default max memory to parse multipart form
var DefaultMaxMemory int64 = 32 << 20

// Form implements [Parser], it will parse postform/multipart-form into container
// it returns an error when request.[ParseForm]/[ParseMultipartForm] returns error
//
// Example:
//
//	var container = &struct {
//	    Name string `form:"name"` // string value
//	    Age  int `form:"age"` // int value
//	    Tags []string `form:"tags"` // string slice
//	    Photos *formdata.UploadedFiles `form:"photo[]"` // file with same field
//	    Avatar *formdata.UploadedFile `form:"avatar"` // single file
//	    Attachment *multipart.FileHeader `form:"attachment"` // single file with type *multipart.FileHeader
//	    VIP bool `form:"vip"` // bool value
//	}{}
//	err := Form(request, container)
//	if err != nil {
//	    panic(err)
//	}
var Form Binding = func(request *http.Request, container any) error {
	form := &multipart.Form{
		Value: request.Form,
	}
	if request.MultipartForm != nil {
		form.File = request.MultipartForm.File
	}
	return formtostruct(form, container, "form")
}
